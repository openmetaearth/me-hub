package app

import (
	packetforwardmiddleware "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcporttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/openmetaearth/me-hub/x/bridgingfee"
	delayedackmodule "github.com/openmetaearth/me-hub/x/delayedack"
	denommetadatamodule "github.com/openmetaearth/me-hub/x/denommetadata"
	"github.com/openmetaearth/me-hub/x/rollapp/genesisbridge"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (a *AppKeepers) InitTransferStack() {
	a.TransferStack = ibctransfer.NewIBCModule(a.TransferKeeper)

	a.TransferStack = bridgingfee.NewIBCModule(
		a.TransferStack.(ibctransfer.IBCModule),
		a.DelayedAckKeeper,
		a.TransferKeeper,
		a.AccountKeeper.GetModuleAddress(wstakingtypes.BridgeFeePool),
		*a.RollappKeeper,
	)

	a.TransferStack = packetforwardmiddleware.NewIBCMiddleware(
		a.TransferStack,
		a.PacketForwardMiddlewareKeeper,
		0,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
	)
	a.TransferStack = denommetadatamodule.NewIBCModule(a.TransferStack, a.DenomMetadataKeeper, a.RollappKeeper)

	// already instantiated in SetupHooks()
	a.DelayedAckMiddleware.Setup(
		delayedackmodule.WithIBCModule(a.TransferStack),
		delayedackmodule.WithKeeper(a.DelayedAckKeeper),
		delayedackmodule.WithRollappKeeper(a.RollappKeeper),
	)
	a.TransferStack = a.DelayedAckMiddleware
	a.TransferStack = genesisbridge.NewIBCModule(a.TransferStack, a.RollappKeeper, a.TransferKeeper, a.DenomMetadataKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, a.TransferStack)
	a.IBCKeeper.SetRouter(ibcRouter)
}
