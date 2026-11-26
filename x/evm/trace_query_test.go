package evm

import (
	"context"
	"math/big"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogoproto "github.com/cosmos/gogoproto/proto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type traceQueryStub struct {
	evmtypes.QueryServer
	calls int
}

func (s *traceQueryStub) TraceTx(_ context.Context, req *evmtypes.QueryTraceTxRequest) (*evmtypes.QueryTraceTxResponse, error) {
	s.calls++
	if _, err := evmtypes.UnpackTxData(req.Msg.Data); err != nil {
		return nil, err
	}
	for _, msg := range req.Predecessors {
		if _, err := evmtypes.UnpackTxData(msg.Data); err != nil {
			return nil, err
		}
	}
	return &evmtypes.QueryTraceTxResponse{Data: []byte(`{"ok":true}`)}, nil
}

func TestTraceTxABCIValidation(t *testing.T) {
	valid := evmtypes.NewTxContract(big.NewInt(100), 0, big.NewInt(0), 21000, big.NewInt(1), nil, nil, nil, nil)
	for _, tc := range []struct {
		name string
		req  *evmtypes.QueryTraceTxRequest
		want string
	}{
		{"missing message", &evmtypes.QueryTraceTxRequest{}, "msg is required"},
		{"missing data", &evmtypes.QueryTraceTxRequest{Msg: &evmtypes.MsgEthereumTx{}}, "invalid msg data"},
		{"unknown data", &evmtypes.QueryTraceTxRequest{Msg: &evmtypes.MsgEthereumTx{Data: &codectypes.Any{TypeUrl: "/unknown.Tx"}}}, "invalid msg data"},
		{"empty predecessor", &evmtypes.QueryTraceTxRequest{Msg: valid, Predecessors: []*evmtypes.MsgEthereumTx{{}}}, "invalid predecessors[0] data"},
		{"valid with predecessor", &evmtypes.QueryTraceTxRequest{Msg: valid, Predecessors: []*evmtypes.MsgEthereumTx{valid}}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			registry := codectypes.NewInterfaceRegistry()
			evmtypes.RegisterInterfaces(registry)
			router := baseapp.NewGRPCQueryRouter()
			router.SetInterfaceRegistry(registry)
			server := &traceQueryStub{}
			evmtypes.RegisterQueryServer(traceQueryRegistrar{Server: router}, server)
			bz, err := gogoproto.Marshal(tc.req)
			require.NoError(t, err)
			handler := router.Route("/ethermint.evm.v1.Query/TraceTx")
			require.NotNil(t, handler)
			var resp *abci.ResponseQuery
			require.NotPanics(t, func() { resp, err = handler(sdk.Context{}, &abci.RequestQuery{Data: bz}) })
			if tc.want != "" {
				require.Equal(t, codes.InvalidArgument, status.Code(err))
				require.ErrorContains(t, err, tc.want)
				require.Zero(t, server.calls)
				return
			}
			require.NoError(t, err)
			require.Equal(t, 1, server.calls)
			var result evmtypes.QueryTraceTxResponse
			require.NoError(t, codec.NewProtoCodec(registry).Unmarshal(resp.Value, &result))
			require.JSONEq(t, `{"ok":true}`, string(result.Data))
		})
	}
}

type traceServiceCapture struct {
	desc   *grpc.ServiceDesc
	server interface{}
}

func (s *traceServiceCapture) RegisterService(desc *grpc.ServiceDesc, server interface{}) {
	s.desc, s.server = desc, server
}

func TestTraceTxDirectGRPCAndInterceptor(t *testing.T) {
	capture := new(traceServiceCapture)
	server := &traceQueryStub{}
	evmtypes.RegisterQueryServer(traceQueryRegistrar{Server: capture}, server)
	for _, method := range capture.desc.Methods {
		if method.MethodName != "TraceTx" {
			continue
		}
		intercepted := false
		_, err := method.Handler(server, context.Background(), func(input interface{}) error {
			return proto.Unmarshal(nil, input.(proto.Message))
		}, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			intercepted = true
			require.Equal(t, "/ethermint.evm.v1.Query/TraceTx", info.FullMethod)
			return handler(ctx, req)
		})
		require.True(t, intercepted)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
		require.ErrorContains(t, err, "msg is required")
		require.Zero(t, server.calls)
		return
	}
	t.Fatal("TraceTx handler not registered")
}

func TestTraceTxHybridRoute(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	evmtypes.RegisterInterfaces(registry)
	router := baseapp.NewGRPCQueryRouter()
	router.SetInterfaceRegistry(registry)
	server := &traceQueryStub{}
	evmtypes.RegisterQueryServer(traceQueryRegistrar{Server: router}, server)
	handlers := router.HybridHandlerByRequestName("ethermint.evm.v1.QueryTraceTxRequest")
	require.Len(t, handlers, 1)
	valid := evmtypes.NewTxContract(big.NewInt(100), 0, big.NewInt(0), 21000, big.NewInt(1), nil, nil, nil, nil)
	var response evmtypes.QueryTraceTxResponse
	require.NoError(t, handlers[0](context.Background(), &evmtypes.QueryTraceTxRequest{Msg: valid}, &response))
	require.JSONEq(t, `{"ok":true}`, string(response.Data))
	err := handlers[0](context.Background(), &evmtypes.QueryTraceTxRequest{}, &response)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Equal(t, 1, server.calls)
}

func TestTraceTxNilPredecessor(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	evmtypes.RegisterInterfaces(registry)
	valid := evmtypes.NewTxContract(big.NewInt(100), 0, big.NewInt(0), 21000, big.NewInt(1), nil, nil, nil, nil)
	req := &evmtypes.QueryTraceTxRequest{Msg: valid, Predecessors: []*evmtypes.MsgEthereumTx{nil}}
	err := unpackTraceTx(req, registry)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.ErrorContains(t, err, "predecessors[0] is required")
}
