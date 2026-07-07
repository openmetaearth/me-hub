package did

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/did/keeper"
	"github.com/openmetaearth/me-hub/x/did/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	for _, info := range genState.Infos {
		addr := sdk.MustAccAddressFromBech32(info.Address)
		k.SetDID(ctx, addr, info.Did)
		k.SetDidInfo(ctx, info.Did, info)
	}

	for _, svc := range genState.Svcs {
		k.SetService(ctx, svc.Sid, svc)
	}

	for _, vc := range genState.Vcs {
		k.SetCredential(ctx, vc.Did, vc.Sid, vc)
	}

	// set filter
	for _, flog := range genState.Flogs {
		vc, found := k.GetCredential(ctx, flog.Did, flog.Sid)
		if !found {
			panic(fmt.Errorf("credential[did:%s, sid:%s] not found", flog.Did, flog.Sid))
		}
		k.AddFilters(ctx, flog.Did, flog.Sid, flog.Filters, vc)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Infos = k.GetDidInfos(ctx)
	genesis.Svcs = k.GetServices(ctx)
	genesis.Vcs = k.GetCredentials(ctx)
	genesis.Flogs = k.GetFilterLoggers(ctx)
	return genesis
}
