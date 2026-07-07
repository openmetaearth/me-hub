package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

func (k *Keeper) GetService(ctx sdk.Context) (svc didtypes.Service, found bool) {
	return k.didKeeper.GetService(ctx, types.ModuleName)
}

func (k *Keeper) SetService(ctx sdk.Context, svc didtypes.Service) {
	k.didKeeper.SetService(ctx, types.ModuleName, svc)
}

func (k *Keeper) GetRegion(ctx sdk.Context, regionId string) (reg types.Region, found bool) {
	raw, found := k.stkKeeper.GetRegion(ctx, regionId)
	if !found {
		return types.Region{}, false
	}

	return types.Region{Id: raw.RegionId, Name: raw.Name}, true
}

func (k *Keeper) GetAllRegions(ctx sdk.Context) (regs []types.Region) {
	raws := k.stkKeeper.GetAllRegion(ctx)
	regions := make([]types.Region, len(raws))
	for _, rew := range raws {
		regions = append(regions, types.Region{Id: rew.RegionId, Name: rew.Name})
	}

	return regions
}

func (k *Keeper) GetProtocol(ctx sdk.Context) (types.Protocol, bool) {
	svc, found := k.GetService(ctx)
	if !found {
		return types.Protocol{}, false
	}

	regions := k.GetAllRegions(ctx)
	return types.Protocol{Service: svc, Regions: regions}, true
}

func (k *Keeper) SetProtocol(ctx sdk.Context, proto types.Protocol) {
	k.SetService(ctx, proto.Service)
}
