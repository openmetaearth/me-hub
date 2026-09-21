package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/did/types"
)

func (k Keeper) SetZkCoreIDBinding(ctx sdk.Context, did string, coreInfo types.ZkAuthClaimCoreIdInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetZkCoreIDInfoKey(did), k.cdc.MustMarshal(&coreInfo))
	store.Set(types.GetZkCoreIDDidKey(coreInfo.CoreId), []byte(did))
}

func (k Keeper) GetZkCoreIDInfo(ctx sdk.Context, did string) (types.ZkAuthClaimCoreIdInfo, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetZkCoreIDInfoKey(did))
	if bz == nil {
		return types.ZkAuthClaimCoreIdInfo{}, false
	}

	var info types.ZkAuthClaimCoreIdInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

func (k Keeper) GetDIDByZkCoreID(ctx sdk.Context, coreID []byte) (string, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetZkCoreIDDidKey(coreID))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

func (k Keeper) SetZkContractAddressInfo(ctx sdk.Context, contractInfo types.ZkContractAddressInfo) {
	//contractInfo.ContractAddress = common.HexToAddress(contractInfo.ContractAddress).Hex()
	ctx.KVStore(k.storeKey).Set(types.GetZkContractAddressKey(contractInfo.Type),
		k.cdc.MustMarshal(&contractInfo))
}

func (k Keeper) GetZkContractAddress(ctx sdk.Context, contractType types.ZkContractAddressType) *types.ZkContractAddressInfo {
	bz := ctx.KVStore(k.storeKey).Get(types.GetZkContractAddressKey(contractType))
	if bz == nil {
		return nil
	}

	var info types.ZkContractAddressInfo
	k.cdc.MustUnmarshal(bz, &info)
	return &info
}
