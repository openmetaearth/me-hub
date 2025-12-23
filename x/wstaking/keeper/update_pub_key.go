package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k Keeper) UpdateValidatorPubKey(ctx sdk.Context) (*types.ReplaceNodePubKey,error) {
	updateInfo, err := k.GetRepalceConsensusPubKeyInfo(ctx)
	if err != nil {
		return nil,err
	}
	if updateInfo == nil {
		return nil,nil
	}
	if ctx.BlockHeight() < updateInfo.UpdateAtHeight {
		return nil,nil
	} else if ctx.BlockHeight() > updateInfo.UpdateAtHeight {
		//delete record
		k.Logger(ctx).Error("delete replace pubkey info delayed.", "need delete at %d", updateInfo.UpdateAtHeight,
			"now height is %d", ctx.BlockHeight())
		k.DeleteRepalceConsensusPubKey(ctx)
		return nil,nil
	} else {
		//do update
		k.Logger(ctx).Info("start to replace validator pubkey", "operator_address", updateInfo.OperatorAddress,
			"pub_key", hex.EncodeToString(updateInfo.PubKey), "height", ctx.BlockHeight())
		pk := new(ed25519.PubKey)
		if err = pk.Unmarshal(updateInfo.PubKey); err != nil {
			return nil, sdkerrors.Wrapf(types.ErrProtoProc, "unmarshal pubkey error: %v,inputKey = %s",
				err, hex.EncodeToString(updateInfo.PubKey))
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
		oldConAddr, err := validator.GetConsAddr()
		if err != nil {
			return nil, err
		}
		oldPubKey,ok := validator.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
		if !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T",  validator.ConsensusPubkey.GetCachedValue())
		}
		k.RemoveValidatorByConsAddr(ctx, oldConAddr)
		anyPk, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, err
		}
		validator.ConsensusPubkey = anyPk
		k.SetValidator(ctx, validator)
		k.SetValidatorByConsAddr(ctx, validator)
		k.DeleteRepalceConsensusPubKey(ctx)
	
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeStartReplacePubKey,
				sdk.NewAttribute(types.AttributeKeyOperatorAddress, updateInfo.OperatorAddress),
				sdk.NewAttribute(types.AttributeKeyOldConsAddr, oldConAddr.String()),
				sdk.NewAttribute(types.AttributeKeyNowConsAddr, sdk.GetConsAddress(pk).String()),
				sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight()))),
		)
		return &types.ReplaceNodePubKey{
			OperatorAddress: updateInfo.OperatorAddress,
			OldPubKey: oldPubKey,
			NewPubKey: pk,
		},nil

	}

}

func (k Keeper) SetRepalcePubKeyInfo(ctx sdk.Context, data *types.UpdatePubKeyInfo) error {
	bz, err := json.Marshal(data)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "marshal repalce pubkey info error: %v", err)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	store.Set(types.KeyPrefix(types.ReplaceConsensusPubKey), bz)
	return nil
}

// GetGroup returns a group from its id
func (k Keeper) GetRepalceConsensusPubKeyInfo(ctx sdk.Context) (*types.UpdatePubKeyInfo, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	data := store.Get(types.KeyPrefix(types.ReplaceConsensusPubKey))
	if data == nil {
		return nil, nil
	}
	val := types.UpdatePubKeyInfo{}
	err := json.Unmarshal(data, &val)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "unmarshal repalce pubkey info error: %v", err)
	}
	return &val, nil
}

func (k Keeper) DeleteRepalceConsensusPubKey(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	store.Delete(types.KeyPrefix(types.ReplaceConsensusPubKey))
}

func (k Keeper) IsHasRepalceConsensusPubKey(ctx sdk.Context) bool {
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
