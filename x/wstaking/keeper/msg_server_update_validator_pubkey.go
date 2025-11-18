package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k MsgServer) UpdateValidatorPubkey(goCtx context.Context, msg *types.MsgUpdateValidatorPubkey) (*types.MsgUpdateValidatorPubkeyResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.daoKeeper.IsGlobalDao(ctx, msg.StakerAddress) {
		return nil, types.ErrCheckGlobalDao
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	// Get old consensus pubkey before update - CRITICAL for removing old pubkey from validator set
	oldPk, err := validator.ConsPubKey()
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "error getting old consensus pubkey: %s", err)
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	newConsAddr := sdk.GetConsAddress(pk)

	// Check if the new pubkey is already in use by another validator
	if _, found := k.GetValidatorByConsAddr(ctx, newConsAddr); found {
		return nil, stakingtypes.ErrValidatorPubKeyExists
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

	k.RemoveValidator(ctx, validator.GetOperator())
	k.DeleteLastValidatorPower(ctx, validator.GetOperator())
	if validator.Status == stakingtypes.Unbonding {
		k.DeleteValidatorQueue(ctx, validator)
	}

	// Update the validator's consensus pubkey
	validator.ConsensusPubkey = msg.Pubkey

	// Set the validator with new pubkey
	k.SetValidator(ctx, validator)

	// Set the new consensus address mapping
	if err := k.SetValidatorByConsAddr(ctx, validator); err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set validator by consensus address")
	}

	// Update validator by power index
	k.SetValidatorByPowerIndex(ctx, validator)

	// If validator is bonded, we need to queue validator updates for EndBlock
	// CRITICAL: Per Cosmos SDK ADR-016 and CometBFT requirements, we MUST generate TWO updates:
	// 1. Remove the old pubkey (power = 0) - This tells CometBFT to stop expecting signatures from old key
	// 2. Add the new pubkey (power = current) - This tells CometBFT to start expecting signatures from new key
	// Without removing the old pubkey, the validator set will contain BOTH keys, causing consensus failures
	if validator.IsBonded() {
		oldConsAddr := sdk.ConsAddress(oldPk.Address())
		newConsAddr := sdk.ConsAddress(pk.Address())

		ctx.Logger().Info("Preparing validator pubkey rotation",
			"validator", validator.GetOperator().String(),
			"old_consensus_address", oldConsAddr.String(),
			"new_consensus_address", newConsAddr.String(),
			"power", validator.GetConsensusPower(k.PowerReduction(ctx)),
		)

		// Convert old pubkey to Tendermint proto format for removal
		oldTmProtoPk, err := codec.ToTmProtoPublicKey(oldPk)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to convert old pubkey to tendermint proto")
		}

		// Get current validator power for the new pubkey
		newPower := validator.GetConsensusPower(k.PowerReduction(ctx))

		// Get current validator updates (if any)
		updates := k.GetValidatorUpdates(ctx)

		// CRITICAL: Append TWO updates in correct order
		// First: Remove old pubkey by setting power to 0
		updates = append(updates, abci.ValidatorUpdate{
			PubKey: oldTmProtoPk,
			Power:  0,
		})
		// Second: Add new pubkey with current voting power
		updates = append(updates, validator.ABCIValidatorUpdate(k.PowerReduction(ctx)))

		ctx.Logger().Info("Queued validator updates for EndBlock",
			"validator", validator.GetOperator().String(),
			"total_updates", len(updates),
			"old_pubkey_removed", true,
			"new_pubkey_added", true,
		)

		// Store the validator updates to be returned in EndBlock
		k.SetValidatorUpdates(ctx, updates)

		// Keep last validator power consistent
		k.SetLastValidatorPower(ctx, validator.GetOperator(), newPower)
	}

	if validator.Status == stakingtypes.Unbonding {
		k.InsertUnbondingValidatorQueue(ctx, validator)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateValidatorPubkey,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.OperatorAddress),
			sdk.NewAttribute(types.AttributeKeyValidatorPubKey, pk.String()),
		),
	})

	return &types.MsgUpdateValidatorPubkeyResponse{}, nil
}
