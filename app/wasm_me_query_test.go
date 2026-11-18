package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type wasmQueryRouterStub struct {
	route baseapp.GRPCQueryHandler
}

func (r wasmQueryRouterStub) Route(string) baseapp.GRPCQueryHandler { return r.route }

func TestWasmMEQueryAllowlist(t *testing.T) {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	for _, tc := range []struct {
		path     string
		response proto.Message
	}{
		{"/metaearth.did.Query/Did", &didtypes.QueryDidResponse{Info: didtypes.DidInfo{Did: "did:me:test", Address: "test-account"}}},
		{"/metaearth.wstaking.Query/Region", &wstakingtypes.QueryRegionResponse{Region: wstakingtypes.Region{
			RegionId: "chn", RegionShare: sdkmath.NewInt(1), DelegateInterest: sdkmath.LegacyZeroDec(),
			DelegateAmount: sdkmath.ZeroInt(), FixedDepositAmount: sdkmath.ZeroInt(),
		}}},
		{"/metaearth.wstaking.Query/FixedDeposit", &wstakingtypes.QueryGetFixedDepositResponse{FixedDeposit: wstakingtypes.FixedDeposit{
			Id: 42, Principal: sdk.NewInt64Coin("umec", 100), Interest: sdk.NewInt64Coin("umec", 1), Rate: sdkmath.LegacyZeroDec(),
		}}},
	} {
		t.Run(tc.path, func(t *testing.T) {
			calls := 0
			query := stargateQuerier(cdc, wasmQueryRouterStub{route: func(_ sdk.Context, req *abci.RequestQuery) (*abci.ResponseQuery, error) {
				calls++
				require.Equal(t, tc.path, req.Path)
				require.Equal(t, []byte{1, 2, 3}, req.Data)
				return &abci.ResponseQuery{Value: cdc.MustMarshal(tc.response)}, nil
			}})
			actual, err := query(sdk.Context{}, &wasmvmtypes.StargateQuery{Path: tc.path, Data: []byte{1, 2, 3}})
			require.NoError(t, err)
			require.JSONEq(t, string(cdc.MustMarshalJSON(tc.response)), string(actual))
			require.Equal(t, 1, calls)
		})
	}
	query := stargateQuerier(cdc, wasmQueryRouterStub{route: func(sdk.Context, *abci.RequestQuery) (*abci.ResponseQuery, error) {
		t.Fatal("unapproved query must not reach the router")
		return nil, nil
	}})
	_, err := query(sdk.Context{}, &wasmvmtypes.StargateQuery{Path: "/cosmos.bank.v1beta1.Query/AllBalances"})
	require.ErrorContains(t, err, "path is not allowed")
	_, err = stargateQuerier(cdc, wasmQueryRouterStub{})(sdk.Context{}, &wasmvmtypes.StargateQuery{Path: "/metaearth.did.Query/Did"})
	require.ErrorContains(t, err, "No route")
}

func TestWasmMEQueriesDoNotShareResponseState(t *testing.T) {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	calls := 0
	query := stargateQuerier(cdc, wasmQueryRouterStub{route: func(sdk.Context, *abci.RequestQuery) (*abci.ResponseQuery, error) {
		calls++
		response := &didtypes.QueryDidResponse{}
		if calls == 1 {
			response.Info = didtypes.DidInfo{Did: "first-request-only", Address: "first-account"}
		}
		return &abci.ResponseQuery{Value: cdc.MustMarshal(response)}, nil
	}})
	first, err := query(sdk.Context{}, &wasmvmtypes.StargateQuery{Path: "/metaearth.did.Query/Did"})
	require.NoError(t, err)
	require.Contains(t, string(first), "first-request-only")
	second, err := query(sdk.Context{}, &wasmvmtypes.StargateQuery{Path: "/metaearth.did.Query/Did"})
	require.NoError(t, err)
	require.NotContains(t, string(second), "first-request-only")
	require.NotContains(t, string(second), "first-account")
}
