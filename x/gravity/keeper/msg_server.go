package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

var _ types.MsgServer = MsgServer{}

type MsgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &MsgServer{Keeper: keeper}
}

func (s MsgServer) BondedRelayer(c context.Context, msg *types.MsgBondedRelayer) (*types.MsgBondedRelayerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	if !s.IsProposalRelayer(ctx, msg.RelayerAddress) {
		return nil, types.ErrNotProposedRelayer
	}
	if msg.DelegateAmount.Denom != params.BaseDenom {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "delegate denom got %s, expected %s", msg.DelegateAmount.Denom, params.BaseDenom)
	}
	// check relayer address is not existed
	if _, found := s.GetRelayer(ctx, relayerAddress); found {
		return nil, errorsmod.Wrap(types.ErrInvalid, "relayer already bonded")
	}
	// check external address is bound to relayer
	if _, found := s.GetRelayerByExternalAddress(ctx, msg.ExternalAddress); found {
		return nil, errorsmod.Wrap(types.ErrInvalid, "external already bonded")
	}
	minThreshold := s.GetGravityMinDelegate(ctx)
	relayer := types.Relayer{
		RelayerAddress:  msg.RelayerAddress,
		ExternalAddress: msg.ExternalAddress,
		DelegateAmount:  msg.DelegateAmount.Amount,
		StartHeight:     ctx.BlockHeight(),
		Online:          true,
		SlashTimes:      0,
	}
	if msg.DelegateAmount.Amount.LT(minThreshold) {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if msg.DelegateAmount.Amount.GT(s.GetGravityMaxDelegate(ctx)) {
		return nil, types.ErrDelegateAmountAboveMaximum
	}

	maxBondedAmount := s.GetAllBondedAmount(ctx).Mul(types.AttestationProposalRelayerChangePowerThreshold).Quo(sdk.NewInt(100))
	if !maxBondedAmount.IsZero() && msg.DelegateAmount.Amount.GT(maxBondedAmount) {
		return nil, errorsmod.Wrapf(types.ErrMaxChangePowerLimitExceeded,
			"max bond amount: %s, bond amont: %s", maxBondedAmount, msg.DelegateAmount.Amount)
	}

	if err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, relayerAddress, s.moduleName, sdk.NewCoins(msg.DelegateAmount)); err != nil {
		return nil, err
	}
	s.SetRelayer(ctx, relayerAddress, relayer)
	s.SetRelayerByExternalAddress(ctx, msg.ExternalAddress, relayerAddress)
	s.SetLastTotalPower(ctx)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBondedRelayer,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.RelayerAddress),
		sdk.NewAttribute(types.AttributeKeyReceiver, authtypes.NewModuleAddress(s.moduleName).String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.DelegateAmount.String()),
		sdk.NewAttribute(types.AttributeKeyExternalAddress, msg.ExternalAddress),
	))
	return &types.MsgBondedRelayerResponse{}, nil
}

func (s MsgServer) AddDelegate(c context.Context, msg *types.MsgAddDelegate) (*types.MsgAddDelegateResponse, error) {
	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	ctx := sdk.UnwrapSDKContext(c)
	if !s.IsProposalRelayer(ctx, msg.RelayerAddress) {
		return nil, types.ErrNotProposedRelayer
	}
	if msg.Amount.Denom != params.BaseDenom {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "delegate denom got %s, expected %s", msg.Amount.Denom, params.BaseDenom)
	}
	relayer, found := s.GetRelayer(ctx, relayerAddress)
	if !found {
		return nil, types.ErrNotFoundRelayer
	}

	clearSlashAmount := relayer.GetSlashAmount(s.GetSlashFraction(ctx))
	if clearSlashAmount.IsPositive() && msg.Amount.Amount.LT(clearSlashAmount.Amount) {
		return nil, errorsmod.Wrap(types.ErrInvalid, "not sufficient slash amount")
	}
	if clearSlashAmount.IsPositive() {
		if err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, relayerAddress, types.SlashingModuleAccount, sdk.NewCoins(clearSlashAmount)); err != nil {
			return nil, err
		}
	}

	delegateCoin := msg.Amount.Sub(clearSlashAmount)
	minThreshold := s.GetGravityMinDelegate(ctx)
	relayer.DelegateAmount = relayer.DelegateAmount.Add(delegateCoin.Amount)
	if relayer.DelegateAmount.LT(minThreshold) {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if relayer.DelegateAmount.GT(s.GetGravityMaxDelegate(ctx)) {
		return nil, types.ErrDelegateAmountAboveMaximum
	}

	maxBondedAmount := s.GetAllBondedAmount(ctx).Mul(types.AttestationProposalRelayerChangePowerThreshold).Quo(sdk.NewInt(100))
	if relayer.DelegateAmount.GT(maxBondedAmount) {
		return nil, errorsmod.Wrapf(types.ErrMaxChangePowerLimitExceeded,
			"max bond amount: %s, bond amont: %s", maxBondedAmount, msg.Amount)
	}

	if delegateCoin.IsPositive() {
		if err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, relayerAddress, s.moduleName, sdk.NewCoins(delegateCoin)); err != nil {
			return nil, err
		}
	}

	if !relayer.Online {
		relayer.Online = true
		relayer.StartHeight = ctx.BlockHeight()
	}

	relayer.SlashTimes = 0
	s.SetRelayer(ctx, relayerAddress, relayer)
	s.SetLastTotalPower(ctx)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAddDelegate,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.RelayerAddress),
		sdk.NewAttribute(types.AttributeKeyReceiver, authtypes.NewModuleAddress(s.moduleName).String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
	))
	return &types.MsgAddDelegateResponse{}, nil
}

func (s MsgServer) UnbondedRelayer(c context.Context, msg *types.MsgUnbondedRelayer) (*types.MsgUnbondedRelayerResponse, error) {
	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	ctx := sdk.UnwrapSDKContext(c)

	//if ctx.BlockHeight() > 10167300 {
	//	if s.IsProposalRelayer(ctx, msg.RelayerAddress) {
	//		return nil, errorsmod.Wrap(types.ErrInvalid, "need to pass a proposal to unbond")
	//	}
	//}

	relayer, found := s.GetRelayer(ctx, relayerAddress)
	if !found {
		return nil, types.ErrNotFoundRelayer
	}

	if relayer.Online {
		return nil, errorsmod.Wrap(types.ErrInvalid, "relayer on line")
	}

	slashAmount := relayer.GetSlashAmount(s.GetSlashFraction(ctx))
	if slashAmount.IsPositive() {
		if err := s.bankKeeper.SendCoinsFromModuleToModule(ctx, s.moduleName, types.SlashingModuleAccount, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
	}

	unbondAmount := relayer.DelegateAmount.Sub(slashAmount.Amount)
	if unbondAmount.IsPositive() {
		if err := s.bankKeeper.SendCoinsFromModuleToAccount(ctx, s.moduleName, relayerAddress, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, unbondAmount))); err != nil {
			return nil, err
		}
	}

	s.DelRelayerByExternalAddress(ctx, relayer.ExternalAddress)
	s.DelRelayer(ctx, relayerAddress)
	s.DelLastEventNonceByRelayer(ctx, relayerAddress)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUnBondedRelayer,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.RelayerAddress),
		sdk.NewAttribute(types.AttributeKeyReceiver, authtypes.NewModuleAddress(s.moduleName).String()),
		sdk.NewAttribute(types.AttributeKeySlashAmount, slashAmount.String()),
		sdk.NewAttribute(types.AttributeKeyUnbondAmount, unbondAmount.String()),
	))
	return &types.MsgUnbondedRelayerResponse{}, nil
}

func (s MsgServer) RelayerSetConfirm(c context.Context, msg *types.MsgRelayerSetConfirm) (*types.MsgRelayerSetConfirmResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	relayerSet := s.GetRelayerSet(ctx, msg.Nonce)
	if relayerSet == nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, "couldn't find relayerSet")
	}

	checkpoint, err := relayerSet.GetCheckpoint(s.GetGravityID(ctx))
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, err.Error())
	}

	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	if err = s.confirmHandlerCommon(ctx, relayerAddress, msg.ExternalAddress, msg.Signature, checkpoint); err != nil {
		return nil, err
	}

	// check if we already have this confirm
	if s.GetRelayerSetConfirm(ctx, msg.Nonce, relayerAddress) != nil {
		return nil, types.ErrDuplicateRelayerConfirms
	}

	s.SetRelayerSetConfirm(ctx, relayerAddress, msg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRelayerSetConfirm,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.RelayerAddress),
	))
	return &types.MsgRelayerSetConfirmResponse{}, nil
}

// RelayerSetUpdateClaim handles claims for executing a relayer set update on Ethereum
func (s MsgServer) RelayerSetUpdateClaim(c context.Context, msg *types.MsgRelayerSetUpdateClaim) (*types.MsgRelayerSetUpdateClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	err := s.checkIsRelayer(ctx, relayerAddress)
	if err != nil {
		return nil, err
	}

	for _, member := range msg.Members {
		if _, found := s.GetRelayerByExternalAddress(ctx, member.ExternalAddress); !found {
			return nil, errorsmod.Wrapf(types.ErrInvalid, "external address not exist %s", member.ExternalAddress)
		}
	}

	// Add the claim to the store
	if _, err := s.Attest(ctx, relayerAddress, msg); err != nil {
		return nil, err
	}
	return &types.MsgRelayerSetUpdateClaimResponse{}, nil
}

func (s MsgServer) BridgeTokenClaim(c context.Context, msg *types.MsgBridgeTokenClaim) (*types.MsgBridgeTokenClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.claimHandlerCommon(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgBridgeTokenClaimResponse{}, nil
}

func (s MsgServer) SendToMeClaim(c context.Context, msg *types.MsgSendToMeClaim) (*types.MsgSendToMeClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.claimHandlerCommon(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgSendToMeClaimResponse{}, nil
}

func (s MsgServer) SendToExternal(c context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	txID, err := s.AddToOutgoingPool(ctx, sender, msg.Dest, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendToExternalResponse{OutgoingTxId: txID}, nil
}

func (s MsgServer) CancelSendToExternal(c context.Context, msg *types.MsgCancelSendToExternal) (*types.MsgCancelSendToExternalResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if _, err := s.RemoveFromOutgoingPoolAndRefund(ctx, msg.TransactionId, sdk.MustAccAddressFromBech32(msg.Sender)); err != nil {
		return nil, err
	}
	return &types.MsgCancelSendToExternalResponse{}, nil
}

func (s MsgServer) SendToExternalClaim(c context.Context, msg *types.MsgSendToExternalClaim) (*types.MsgSendToExternalClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.claimHandlerCommon(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgSendToExternalClaimResponse{}, nil
}

func (s MsgServer) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, "sender address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	bridgeToken, err := s.GetBridgeTokenByDenom(ctx, msg.Denom)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "get bridge token: %v", err)
	}

	if err := s.checkIsRelayer(ctx, sender); err != nil {
		return nil, err
	}

	batch, err := s.BuildOutgoingTxBatch(ctx, bridgeToken.ContractAddress, msg.FeeReceive, types.OutgoingTxBatchSize, msg.MinimumFee, msg.BaseFee)
	if err != nil {
		return nil, err
	}

	return &types.MsgRequestBatchResponse{
		BatchNonce: batch.BatchNonce,
	}, nil
}

func (s MsgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := s.GetOutgoingTxBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, "couldn't find batch")
	}

	checkpoint, err := batch.GetCheckpoint(s.GetGravityID(ctx))
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, err.Error())
	}

	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	if err = s.confirmHandlerCommon(ctx, relayerAddress, msg.ExternalAddress, msg.Signature, checkpoint); err != nil {
		return nil, err
	}

	// check if we already have this confirm
	if s.GetBatchConfirm(ctx, msg.TokenContract, msg.Nonce, relayerAddress) != nil {
		return nil, types.ErrDuplicateRelayerConfirms
	}

	s.SetBatchConfirm(ctx, relayerAddress, msg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatchConfirm,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.RelayerAddress),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprintf("%d", msg.Nonce)),
		sdk.NewAttribute(types.AttributeKeyTokenContract, msg.TokenContract),
	))
	return &types.MsgConfirmBatchResponse{}, nil
}

func (s MsgServer) IncreaseBridgeFee(c context.Context, msg *types.MsgIncreaseBridgeFee) (*types.MsgIncreaseBridgeFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.AddUnbatchedTxBridgeFee(ctx, msg.TransactionId, sdk.MustAccAddressFromBech32(msg.Sender), msg.AddBridgeFee); err != nil {
		return nil, err
	}
	return &types.MsgIncreaseBridgeFeeResponse{}, nil
}

func (s MsgServer) ProposalRelayers(c context.Context, msg *types.MsgProposalRelayers) (*types.MsgProposalRelayersResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if !s.daoKeeper.IsDao(ctx, msg.Authority) {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority")
	}
	if err := s.UpdateProposalRelayers(ctx, msg.Relayers); err != nil {
		return nil, err
	}
	return &types.MsgProposalRelayersResponse{}, nil
}

func (s MsgServer) UpdateParams(c context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if s.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", s.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.SetParams(ctx, &req.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}
