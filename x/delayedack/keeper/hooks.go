package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	eibctypes "github.com/openmetaearth/me-hub/x/eibc/types"
)

/* -------------------------------------------------------------------------- */
/*                                 eIBC Hooks                                 */
/* -------------------------------------------------------------------------- */
var _ eibctypes.EIBCHooks = eibcHooks{}

const (
	deletePacketsBatchSize = 1000
)

type eibcHooks struct {
	eibctypes.BaseEIBCHook
	Keeper
}

func (k Keeper) GetEIBCHooks() eibctypes.EIBCHooks {
	return eibcHooks{
		BaseEIBCHook: eibctypes.BaseEIBCHook{},
		Keeper:       k,
	}
}

// AfterDemandOrderFulfilled is called every time a demand order is fulfilled.
// Once it is fulfilled the underlying packet recipient should be updated to the fulfiller.
func (k eibcHooks) AfterDemandOrderFulfilled(ctx sdk.Context, demandOrder *eibctypes.DemandOrder, fulfillerAddress string) error {
	err := k.UpdateRollappPacketTransferAddress(ctx, demandOrder.TrackingPacketKey, fulfillerAddress)
	if err != nil {
		return err
	}
	return nil
}
