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
	"github.com/st-chain/me-hub/utils"
	"github.com/st-chain/me-hub/x/wstaking/types"
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

// todo:
// 1. 增加区块高度到达时，节点才进行替换
// 2. 增加power的转移，以及质押信息以及状态信息的转移
// 4. 检查是否有业务质押和验证者关系的绑定，以及是否需要转移
func (k MsgServer) ReplaceConsensusPubKey(goCtx context.Context, req *types.MsgReplaceConsensusPubKeyRequest) (*types.MsgReplaceConsensusPubKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.daoKeeper.IsGlobalDao(ctx, req.Creator) {
		return nil, types.ErrCheckGlobalDao
	}
	// Check if any validator replacement is already in progress
	if k.IsHasRepalceConsensusPubKey(ctx) {
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

	tmpPk, ok := req.ReplacePubKey.PubKey.GetCachedValue().(cryptotypes.PubKey)
    if !ok {
        return nil, sdkerrors.Wrapf(stakingtypes.ErrValidatorPubKeyTypeNotSupported, "Expecting cryptotypes.PubKey, got %T", req.ReplacePubKey.PubKey.GetCachedValue())
    }

   	pk, ok := tmpPk.(*ed25519.PubKey)
	if !ok {
        return nil, sdkerrors.Wrapf(stakingtypes.ErrValidatorPubKeyTypeNotSupported, "Expecting ed25519.PubKey, got %T", pk)
    }

	/*
	pk, ok := req.ReplacePubKey.PubKey.GetCachedValue().(*ed25519.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(stakingtypes.ErrValidatorPubKeyTypeNotSupported, "Expecting ed25519.PubKey, got %T", pk)
	}
		*/

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, stakingtypes.ErrValidatorPubKeyExists
	}

	/*
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

	*/
	pubKeyData, err := pk.Marshal()
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrProtoProc, "marshal pubkey error: %v", err)
	}

	update := &types.UpdatePubKeyInfo{
		OperatorAddress: req.ReplacePubKey.OperatorAddress,
		PubKey:          pubKeyData,
		UpdateAtHeight:  ctx.BlockHeight() + req.ReplacePubKey.BlockNumber,
	}

	if err = k.SetRepalcePubKeyInfo(ctx, update); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReplacePubKey,
			sdk.NewAttribute(types.AttributeKeyOperatorAddress, update.OperatorAddress),
			sdk.NewAttribute(types.AttributeKeyPubKey, hex.EncodeToString(update.PubKey)),
			sdk.NewAttribute(types.AttributeKeyUpdateAtHeight, fmt.Sprintf("%d", update.UpdateAtHeight)),
		),
	})

	/*
		newValidatorValAddr := sdk.ValAddress(pk.Address())
		if k.MoveStakesToAnotherVal(ctx, oldValAddr, newValidatorValAddr) != nil {
			return nil, err
		}
		//replace validator pubkey
		newValidator := oldValidator
		newValidator.ConsensusPubkey = req.ReplaceValidator.NewValidatorPubKey

		//todo
		/*
			1. 查询旧validator所关联的regiond_id
			2. 解绑旧validator和region_id的绑定关系
			3. 绑定新validator和region_id的绑定关系
	*/

	/*
		// Remove old ConsAddr index
		oldConsensAddr, err := validator.GetConsAddr()
		if err != nil {
			return nil, err
		}
		k.RemoveValidatorByConsAddr(ctx, oldConsensAddr)

		// Update Validator with new PubKey
		validator.ConsensusPubkey = req.UpdatePubKey.PubKey
		// Update Validator in store
		k.SetValidator(ctx, validator)
		// Set new ConsAddr index
		k.SetValidatorByConsAddr(ctx, validator)

		// Store pending removal for old key
		k.SetPendingKeyRemoval(ctx, valAddr, oldPk)
	*/

	return &types.MsgReplaceConsensusPubKeyResponse{}, nil
}
