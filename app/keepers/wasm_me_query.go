package keepers

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func setupStargateQueryPlugin(appCodec codec.Codec, queryRouter wasmkeeper.GRPCQueryRouter) wasmkeeper.Option {
	return wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Stargate: stargateQuerier(appCodec, queryRouter),
	})
}

func stargateQuerier(appCodec codec.Codec, queryRouter wasmkeeper.GRPCQueryRouter) func(sdk.Context, *wasmvmtypes.StargateQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error) {
		// The wasmd querier mutates response instances, so allocate them per request.
		acceptedQueries := wasmkeeper.AcceptedStargateQueries{
			"/metaearth.wstaking.Query/FixedDeposit": &wstakingtypes.QueryGetFixedDepositResponse{},
		}
		return wasmkeeper.AcceptListStargateQuerier(acceptedQueries, queryRouter, appCodec)(ctx, request)
	}
}
