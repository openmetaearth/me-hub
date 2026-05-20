package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	"github.com/openmetaearth/me-hub/x/delayedack/types"
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

/* -------------------------------------------------------------------------- */
/*                               Epoch Hooks                                  */
/* -------------------------------------------------------------------------- */

type epochHooks struct {
	Keeper
}

// GetEpochHooks returns an epochHooks instance wrapping the keeper.
func (k Keeper) GetEpochHooks() epochHooks {
	return epochHooks{Keeper: k}
}

// AfterEpochEnd deletes all finalized and reverted rollapp packets for the matching epoch identifier.
func (e epochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) error {
	params := e.GetParams(ctx)
	if params.EpochIdentifier != epochIdentifier {
		return nil
	}

	limit := int(params.DeletePacketsEpochLimit)
	filter := types.ByStatus(commontypes.Status_FINALIZED, commontypes.Status_REVERTED)
	if limit > 0 {
		filter = filter.Take(limit)
	}

	packets := e.ListRollappPackets(ctx, filter)
	for _, packet := range packets {
		p := packet
		if err := e.deleteRollappPacket(ctx, &p); err != nil {
			return err
		}
	}

	return nil
}
