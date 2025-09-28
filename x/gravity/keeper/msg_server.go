package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/gravity/types"
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
	// check relayer address is not existed
	if _, found := s.GetRelayer(ctx, relayerAddress); found {
		return nil, errorsmod.Wrap(types.ErrInvalid, "oracle existed bridger address")
	}
	// check external address is bound to oracle
	if _, found := s.GetRelayerByExternalAddress(ctx, msg.ExternalAddress); found {
		return nil, errorsmod.Wrap(types.ErrInvalid, "external address is bound to oracle")
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
	if minThreshold.Denom != msg.DelegateAmount.Denom {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "delegate denom got %s, expected %s", msg.DelegateAmount.Denom, minThreshold.Denom)
	}
	if msg.DelegateAmount.IsLT(minThreshold) {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if msg.DelegateAmount.Amount.GT(s.GetGravityMaxDelegate(ctx).Amount) {
		return nil, types.ErrDelegateAmountAboveMaximum
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
	if _, found := s.GetRelayer(ctx, relayerAddress); !found {
		return nil, errorsmod.Wrap(types.ErrInvalid, "no found bridger address")
	}
	relayer, found := s.GetRelayer(ctx, relayerAddress)
	if !found {
		return nil, types.ErrNotFoundRelayer
	}

	minThreshold := s.GetGravityMinDelegate(ctx)
	if minThreshold.Denom != msg.Amount.Denom {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "delegate denom got %s, expected %s", msg.Amount.Denom, minThreshold.Denom)
	}

	relayer.DelegateAmount = relayer.DelegateAmount.Add(msg.Amount.Amount)
	if relayer.DelegateAmount.LT(minThreshold.Amount) {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if relayer.DelegateAmount.GT(s.GetGravityMaxDelegate(ctx).Amount) {
		return nil, types.ErrDelegateAmountAboveMaximum
	}

	if err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, relayerAddress, s.moduleName, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, err
	}

	if !relayer.Online {
		relayer.Online = true
		relayer.StartHeight = ctx.BlockHeight()
		relayer.SlashTimes = 0
	}

	s.SetRelayer(ctx, relayerAddress, relayer)
	s.SetLastTotalPower(ctx)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBondedRelayer,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.RelayerAddress),
		sdk.NewAttribute(types.AttributeKeyReceiver, authtypes.NewModuleAddress(s.moduleName).String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
	))
	return &types.MsgAddDelegateResponse{}, nil
}

func (s MsgServer) UnbondedRelayer(c context.Context, msg *types.MsgUnbondedRelayer) (*types.MsgUnbondedRelayerResponse, error) {
	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	ctx := sdk.UnwrapSDKContext(c)

	if s.IsProposalRelayer(ctx, msg.RelayerAddress) {
		return nil, errorsmod.Wrap(types.ErrInvalid, "need to pass a proposal to unbind")
	}

	relayer, found := s.GetRelayer(ctx, relayerAddress)
	if !found {
		return nil, types.ErrNotFoundRelayer
	}

	if relayer.Online {
		return nil, errorsmod.Wrap(types.ErrInvalid, "oracle on line")
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
		return nil, errorsmod.Wrap(types.ErrInvalid, "couldn't find oracleSet")
	}

	checkpoint, err := relayerSet.GetCheckpoint(s.GetGravityID(ctx))
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, err.Error())
	}

	relayerAddress, err := s.confirmHandlerCommon(ctx, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}

	// check if we already have this confirm
	if s.GetRelayerSetConfirm(ctx, msg.Nonce, relayerAddress) != nil {
		return nil, errorsmod.Wrap(types.ErrDuplicate, "signature")
	}

	s.SetRelayerSetConfirm(ctx, relayerAddress, msg)
	return &types.MsgRelayerSetConfirmResponse{}, nil
}

// RelayerSetUpdateClaim handles claims for executing a relayer set update on Ethereum
func (s MsgServer) RelayerSetUpdateClaim(c context.Context, msg *types.MsgRelayerSetUpdatedClaim) (*types.MsgRelayerSetUpdatedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	err := s.checkIsRelayer(ctx, relayerAddress)
	if err != nil {
		return nil, err
	}

	for _, member := range msg.Members {
		if _, found := s.GetRelayerByExternalAddress(ctx, member.ExternalAddress); !found {
			return nil, errorsmod.Wrap(types.ErrInvalid, "external address")
		}
	}

	// Add the claim to the store
	if _, err := s.Attest(ctx, relayerAddress, msg); err != nil {
		return nil, err
	}
	return &types.MsgRelayerSetUpdatedClaimResponse{}, nil
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

	bridgeToken, _ := s.GetBridgeTokenByDenom(ctx, msg.Denom)
	if bridgeToken == nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	if err := s.checkIsRelayer(ctx, sender); err != nil {
		return nil, err
	}

	batch, err := s.BuildOutgoingTxBatch(ctx, bridgeToken.Contract, msg.FeeReceive, types.OutgoingTxBatchSize, msg.MinimumFee, msg.BaseFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

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

	relayerAddr, err := s.confirmHandlerCommon(ctx, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}

	// check if we already have this confirm
	if s.GetBatchConfirm(ctx, msg.TokenContract, msg.Nonce, relayerAddr) != nil {
		return nil, errorsmod.Wrap(types.ErrDuplicate, "signature")
	}

	s.SetBatchConfirm(ctx, relayerAddr, msg)
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
