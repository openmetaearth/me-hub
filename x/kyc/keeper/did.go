package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

/*
These methods wrap the did methods without any other logic
*/

func (k *Keeper) HasDID(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.didKeeper.HasDID(ctx, addr)
}

func (k *Keeper) GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool) {
	return k.didKeeper.GetDID(ctx, addr)
}

func (k *Keeper) SetDID(ctx sdk.Context, addr sdk.AccAddress, did string) {
	k.didKeeper.SetDID(ctx, addr, did)
}

func (k *Keeper) HasDidInfo(ctx sdk.Context, did string) bool {
	return k.didKeeper.HasDidInfo(ctx, did)
}

func (k *Keeper) GetDidInfo(ctx sdk.Context, did string) (didtypes.DidInfo, bool) {
	return k.didKeeper.GetDidInfo(ctx, did)
}

func (k *Keeper) SetDidInfo(ctx sdk.Context, did string, info didtypes.DidInfo) {
	k.didKeeper.SetDidInfo(ctx, did, info)
}

func (k Keeper) SetKycIssers(ctx sdk.Context, oldDaoAddress, newDaoAddress []string) error {
	if len(oldDaoAddress) != len(newDaoAddress) {
		return nil
	}

	service, ok := k.GetService(ctx)
	if !ok {
		return fmt.Errorf("kyc service not found")
	}

	dids := []string{}
	for i, dao := range oldDaoAddress {
		did, found := k.GetDID(ctx, sdk.MustAccAddressFromBech32(dao))
		if !found {
			return fmt.Errorf("old address %s did not exists, please create did before repalce did", dao)
		}

		k.didKeeper.DeleteDID(ctx, sdk.MustAccAddressFromBech32(dao))
		k.didKeeper.SetDID(ctx, sdk.MustAccAddressFromBech32(newDaoAddress[i]), did)

		didInfo, found := k.GetDidInfo(ctx, did)
		if !found {
			return fmt.Errorf("old address %s did info not exists, please create did info before replace did info", dao)
		}

		didInfo.Address = newDaoAddress[i]
		k.SetDidInfo(ctx, did, didInfo)
		dids = append(dids, did)
	}

	service.Issuers = dids
	k.didKeeper.SetService(ctx, types.ModuleName, service)
	return nil
}
