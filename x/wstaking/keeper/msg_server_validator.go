package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// CreateValidator defines wrapped method for creating a new validator.
func (k MsgServer) CreateValidator(
	goCtx context.Context, msg *stakingtypes.MsgCreateValidator,
) (*stakingtypes.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.DelegatorAddress) {
		return nil, types.ErrCheckGlobalDao
	}

	_, err := utils.CheckRegionName(strings.ToUpper(msg.Description.RegionID))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRegionName, msg.Description.RegionID)
	}
	msg.Description.RegionID = strings.ToLower(msg.Description.RegionID)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	if msg.Commission.Rate.LT(k.MinCommissionRate(ctx)) {
		return nil, sdkerrors.Wrapf(stakingtypes.ErrCommissionLTMinRate, "cannot set validator commission to less than minimum rate of %s", k.MinCommissionRate(ctx))
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, valAddr); found {
		return nil, stakingtypes.ErrValidatorOwnerExists
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", msg.Pubkey.GetCachedValue())
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, stakingtypes.ErrValidatorPubKeyExists
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Value.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Value.Denom, bondDenom,
		)
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		pkType := pk.Type()
		hasKeyType := false
		for _, keyType := range cp.Validator.PubKeyTypes {
			if pkType == keyType {
				hasKeyType = true
				break
			}
		}
		if !hasKeyType {
			return nil, sdkerrors.Wrapf(
				stakingtypes.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	validator, err := stakingtypes.NewValidator(valAddr, pk, msg.Description)
	if err != nil {
		return nil, err
	}

	commission := stakingtypes.NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)

	validator, err = validator.SetInitialCommission(commission)
	if err != nil {
		return nil, err
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	validator.MinSelfDelegation = msg.MinSelfDelegation
	validator.OwnerAddress = sdk.AccAddress(valAddr).String()

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)
	// call the after-creation hook
	if err := k.Hooks().AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
		return nil, err
	}

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	_, err = k.Keeper.Stake(ctx, delegatorAddress, msg.Value.Amount, stakingtypes.Unbonded, validator, true, "create_validator_"+msg.Description.RegionID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeCreateValidator,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.String()),
		),
	})

	return &stakingtypes.MsgCreateValidatorResponse{}, nil
}

func (k MsgServer) EditValidator(context.Context, *stakingtypes.MsgEditValidator) (*stakingtypes.MsgEditValidatorResponse, error) {
	return &stakingtypes.MsgEditValidatorResponse{}, fmt.Errorf("not implemented, please use UpdateValidator instead")
}

// 1. only perform the node replacement when the target block height is reached.
// 2. handle the transfer of power, staking information, and status information.
// 3. check if the new pubkey is already bound to an existing validator
func (k MsgServer) ReplaceConsensusPubKey(goCtx context.Context, req *types.MsgReplaceConsensusPubKeyRequest) (*types.MsgReplaceConsensusPubKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.daoKeeper.IsGlobalDao(ctx, req.Creator) {
		return nil, types.ErrCheckGlobalDao
	}
	// Check if any validator replacement is already in progress
	if k.IsHasReplaceConsensusPubKey(ctx) {
		return nil, types.ErrExistingReplaceValidator
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ReplacePubKey.OperatorAddress)
	if err != nil {
		return nil, err
	}

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}
	if validator.IsJailed() {
		return nil, stakingtypes.ErrValidatorJailed
	}

	if !validator.IsBonded() {
		return nil, types.ErrValidatorNotBonded
	}

	pk, ok := req.ReplacePubKey.PubKey.GetCachedValue().(*ed25519.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(stakingtypes.ErrValidatorPubKeyTypeNotSupported, "Expecting ed25519.PubKey, got %T", req.ReplacePubKey.PubKey.GetCachedValue())
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, stakingtypes.ErrValidatorPubKeyExists
	}

	pubKeyData, err := pk.Marshal()
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrProtoProc, "marshal pubkey error: %v", err)
	}
	needReplaceConsAddr, err := validator.GetConsAddr()
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInterProc, "GetConsAddr from validator error: %v", err)
	}

	update := &types.UpdatePubKeyInfo{
		OperatorAddress: req.ReplacePubKey.OperatorAddress,
		OldConsAddress:  needReplaceConsAddr.Bytes(),
		PubKey:          pubKeyData,
		UpdateAtHeight:  ctx.BlockHeight() + req.ReplacePubKey.BlockNumber,
	}

	if err = k.SetReplacePubKeyInfo(ctx, update); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReplacePubKey,
			sdk.NewAttribute(types.AttributeKeyOperatorAddress, update.OperatorAddress),
			sdk.NewAttribute(types.AttributeKeyPubKey, hex.EncodeToString(update.PubKey)),
			sdk.NewAttribute(types.AttributeKeyOldConsAddr, needReplaceConsAddr.String()),
			sdk.NewAttribute(types.AttributeKeyUpdateAtHeight, fmt.Sprintf("%d", update.UpdateAtHeight)),
		),
	})

	return &types.MsgReplaceConsensusPubKeyResponse{}, nil
}
