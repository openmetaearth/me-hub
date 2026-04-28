package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) UpdateValidatorPubKey(ctx sdk.Context) (*types.ReplaceNodePubKey, error) {
	updateInfo, err := k.GetReplaceConsensusPubKeyInfo(ctx)
	if err != nil {
		panic(fmt.Sprintf("GetReplaceConsensusPubKeyInfo error,err = %s ", err.Error()))
	}
	if updateInfo == nil {
		return nil, nil
	}
	if ctx.BlockHeight() < updateInfo.UpdateAtHeight {
		return nil, nil
	} else {
		if ctx.BlockHeight() == updateInfo.UpdateAtHeight {
			//do update
			k.Logger(ctx).Info("start to replace validator pubkey", "operator_address", updateInfo.OperatorAddress,
				"pub_key", hex.EncodeToString(updateInfo.PubKey), "height", ctx.BlockHeight())
			pk := new(ed25519.PubKey)
			if err = pk.Unmarshal(updateInfo.PubKey); err != nil {
				//return nil, sdkerrors.Wrapf(types.ErrProtoProc, "unmarshal pubkey error: %v,inputKey = %s",
				//	err, hex.EncodeToString(updateInfo.PubKey))
				panic(fmt.Sprintf("unmarshal pubkey error: %s ,inputKey = %s", err.Error(), hex.EncodeToString(updateInfo.PubKey)))
			}
			valAddr, err := sdk.ValAddressFromBech32(updateInfo.OperatorAddress)
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

			if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
				return nil, stakingtypes.ErrValidatorPubKeyExists
			}

			oldPubKey, ok := validator.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
			if !ok {
				panic(fmt.Sprintf("parser validator's old pub key error.expecting cryptotypes.PubKey, got %T", validator.ConsensusPubkey.GetCachedValue()))
			}

			anyPk, err := codectypes.NewAnyWithValue(pk)
			if err != nil {
				panic(fmt.Sprintf("codectypes.NewAnyWithValue ConsensusPubkey error.err = %s ", err.Error()))
			}
			validator.ConsensusPubkey = anyPk
			k.SetValidator(ctx, validator)
			k.SetValidatorByConsAddr(ctx, validator)
			if err = k.Hooks().AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
				k.Logger(ctx).Info("AfterValidatorCreated hook ", "err", err.Error())
				return nil, sdkerrors.Wrapf(types.ErrInterProc, "AfterValidatorCreated hook error: %v", err)
			}
			//directly set new signing info for new cons addr
			newConsAddr := sdk.GetConsAddress(pk)
			newSigningInfo := slashingtypes.NewValidatorSigningInfo(
				newConsAddr,
				ctx.BlockHeight(),
				0,
				time.Unix(0, 0),
				false,
				0,
			)
			k.slashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newSigningInfo)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.EventTypeStartReplacePubKey,
					sdk.NewAttribute(types.AttributeKeyOperatorAddress, updateInfo.OperatorAddress),
					sdk.NewAttribute(types.AttributeKeyOldConsAddr, sdk.ConsAddress(updateInfo.OldConsAddress).String()),
					sdk.NewAttribute(types.AttributeKeyNowConsAddr, sdk.GetConsAddress(pk).String()),
					sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight()))),
			)
			return &types.ReplaceNodePubKey{
				OperatorAddress: updateInfo.OperatorAddress,
				OldPubKey:       oldPubKey,
				NewPubKey:       pk,
			}, nil

		} else if ctx.BlockHeight() == (updateInfo.UpdateAtHeight + 2) { //delay remove old cons addr because of distribution rewards delayed by one block
			k.RemoveValidatorByConsAddr(ctx, sdk.ConsAddress(updateInfo.OldConsAddress))
			k.DeleteReplaceConsensusPubKey(ctx)
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.EventTypeDelayRemoveOldConsAddr,
					sdk.NewAttribute(types.AttributeKeyOperatorAddress, updateInfo.OperatorAddress),
					sdk.NewAttribute(types.AttributeKeyOldConsAddr, sdk.ConsAddress(updateInfo.OldConsAddress).String()),
					sdk.NewAttribute(types.AttributeKeyUpdateAtHeight, fmt.Sprintf("%d", updateInfo.UpdateAtHeight)),
					sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight()))),
			)
			k.Logger(ctx).Info("completed delayed removed old cons addr from index", "old_cons_addr",
				sdk.ConsAddress(updateInfo.OldConsAddress).String(), "height", ctx.BlockHeight())
			return nil, nil

		} else if ctx.BlockHeight() == (updateInfo.UpdateAtHeight + 1) {
			//do nothing, wait for next block to remove old cons addr
			return nil, nil
		} else {
			return nil, sdkerrors.Wrapf(types.ErrInterProc, "ReplaceConsensusPubKeyInfo is still exist when block height greater than update_at_height. now height = %d, update at height = %d",
				ctx.BlockHeight(), updateInfo.UpdateAtHeight)
		}
	}
}

func (k Keeper) SetReplacePubKeyInfo(ctx sdk.Context, data *types.UpdatePubKeyInfo) error {
	bz, err := json.Marshal(data)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "marshal replace pubkey info error: %v", err)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	store.Set(types.KeyPrefix(types.ReplaceConsensusPubKey), bz)
	return nil
}

// GetGroup returns a group from its id
func (k Keeper) GetReplaceConsensusPubKeyInfo(ctx sdk.Context) (*types.UpdatePubKeyInfo, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	data := store.Get(types.KeyPrefix(types.ReplaceConsensusPubKey))
	if data == nil {
		return nil, nil
	}
	val := types.UpdatePubKeyInfo{}
	err := json.Unmarshal(data, &val)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "unmarshal replace pubkey info error: %v", err)
	}
	return &val, nil
}

func (k Keeper) DeleteReplaceConsensusPubKey(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	store.Delete(types.KeyPrefix(types.ReplaceConsensusPubKey))
}

func (k Keeper) IsHasReplaceConsensusPubKey(ctx sdk.Context) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	data := store.Get(types.KeyPrefix(types.ReplaceConsensusPubKey))
	if data == nil {
		return false
	}
	return true
}

func (k Keeper) RemoveValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(stakingtypes.GetValidatorByConsAddrKey(consAddr))
}

func (k Keeper) MoveStakesToAnotherVal(ctx sdk.Context, fromValAddr, toValAddr sdk.ValAddress) error {
	stakes, err := k.GetStakesByValidator(ctx, fromValAddr)
	if err != nil {
		return err
	}
	if 0 == len(stakes) {
		return sdkerrors.Wrapf(types.ErrStakeOnValidatorIsEmpty, "old validatorAddr =%s", fromValAddr.String())
	}
	for _, stake := range stakes {
		//remove old stake record
		k.RemoveStake(ctx, *stake)
		//create new stake record
		stake.ValidatorAddress = toValAddr.String()
		stake.StartHeight = ctx.BlockHeight()
		k.SetStake(ctx, *stake)
	}
	return nil

}
