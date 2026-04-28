package keeper

import (
	"encoding/hex"
	"fmt"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// BlockValidatorUpdates calculates the ValidatorUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	replacePubKey, err := k.UpdateValidatorPubKey(ctx)
	if err != nil {
		updateInfo, errP := k.GetReplaceConsensusPubKeyInfo(ctx)
		if errP != nil {
			panic(fmt.Sprintf("GetReplaceConsensusPubKeyInfo error,err = %s ", errP.Error()))
		}
		k.DeleteReplaceConsensusPubKey(ctx)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeReplacePubKeyFailed,
				sdk.NewAttribute(types.AttributeKeyOperatorAddress, updateInfo.OperatorAddress),
				sdk.NewAttribute(types.AttributeKeyOldConsAddr, sdk.ConsAddress(updateInfo.OldConsAddress).String()),
				sdk.NewAttribute(types.AttributeKeyPubKey, hex.EncodeToString(updateInfo.PubKey)),
				sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
				sdk.NewAttribute(types.AttributeKeyFailedReason, err.Error())))

		k.Logger(ctx).Error("failed to replace validator pubkey", "error", err.Error(),
			"block height", ctx.BlockHeight())
	} else {
		if replacePubKey != nil {
			newPubkey, errP := cryptocodec.ToTmProtoPublicKey(replacePubKey.NewPubKey)
			if errP != nil {
				panic(errP)
			}
			oldPubkey, errP := cryptocodec.ToTmProtoPublicKey(replacePubKey.OldPubKey)
			if errP != nil {
				panic(errP)
			}
			validatorUpdates = append(validatorUpdates, abci.ValidatorUpdate{
				PubKey: oldPubkey,
				Power:  0,
			})
			valAddr, errP := sdk.ValAddressFromBech32(replacePubKey.OperatorAddress)
			if errP != nil {
				panic(fmt.Sprintf("invalid validator address %s,err = %s", replacePubKey.OperatorAddress, errP.Error()))
			}
			validator, found := k.GetValidator(ctx, valAddr)
			if !found {
				panic(fmt.Sprintf("validator not found for address %s", replacePubKey.OperatorAddress))
			}
			power := validator.ConsensusPower(k.PowerReduction(ctx))
			validatorUpdates = append(validatorUpdates, abci.ValidatorUpdate{
				PubKey: newPubkey,
				Power:  power,
			})
			// Log the removal
			k.Logger(ctx).Info("completed pubb key replaced in validatorUpdates ", "validator", valAddr.String(), "block height", ctx.BlockHeight())
		}
	}

	// unbond all mature validators from the unbonding queue
	k.UnbondAllMatureValidators(ctx)

	// Remove all mature unbonding stakes from the ubs queue.
	matureUnBondStakes := k.SequeueAllMatureUBSQueue(ctx, ctx.BlockHeader().Time)
	for _, svPair := range matureUnBondStakes {
		addr, err := sdk.ValAddressFromBech32(svPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		stakerAddress := sdk.MustAccAddressFromBech32(svPair.StakerAddress)

		balances, err := k.CompleteStakeUnBonding(ctx, stakerAddress, addr)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnStakeBonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, svPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyStaker, svPair.StakerAddress),
			),
		)
	}

	// Remove all mature unbonding delegations from the ubd queue.
	matureDelUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureDelUnbonds {
		addr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)

		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnDelBonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyRegionId, strings.ToLower(types.ExperienceRegionName)),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress := sdk.MustAccAddressFromBech32(dvvTriplet.DelegatorAddress)

		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				stakingtypes.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(stakingtypes.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(stakingtypes.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return validatorUpdates
}
