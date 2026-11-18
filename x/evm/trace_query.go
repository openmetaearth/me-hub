package evm

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/gogoproto/grpc"
	gogoproto "github.com/cosmos/gogoproto/proto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// traceQueryRegistrar carries v2.0.18's TraceTx input checks across the SDK
// upgrade without downgrading Ethermint. rc4's request UnpackInterfaces dereferences
// a missing Msg before the keeper is called. Decode this one RPC as a standard
// protobuf message first, then unpack only non-nil Ethereum messages. Other RPCs
// retain their upstream handlers and the wire descriptor is unchanged.
// Remove this adapter once the v0.50-compatible Ethermint release includes
// upstream fix 08a9eaaee040 (both keeper validation and nil-safe unpacking).
type traceQueryRegistrar struct {
	grpc.Server
}

func (r traceQueryRegistrar) RegisterService(sd *googlegrpc.ServiceDesc, server interface{}) {
	desc, err := gogoproto.HybridResolver.FindDescriptorByName("ethermint.evm.v1.Query.TraceTx")
	if err != nil {
		panic(err)
	}
	registry := types.NewInterfaceRegistry()
	evmtypes.RegisterInterfaces(registry)

	service := *sd
	service.Methods = append([]googlegrpc.MethodDesc(nil), sd.Methods...)
	for i := range service.Methods {
		if service.Methods[i].MethodName == "TraceTx" {
			service.Methods[i].Handler = traceTxHandler(desc.(protoreflect.MethodDescriptor), registry)
		}
	}
	r.Server.RegisterService(&service, server)
}

func traceTxHandler(desc protoreflect.MethodDescriptor, unpacker types.AnyUnpacker) func(interface{}, context.Context, func(interface{}) error, googlegrpc.UnaryServerInterceptor) (interface{}, error) {
	return func(server interface{}, ctx context.Context, decode func(interface{}) error, interceptor googlegrpc.UnaryServerInterceptor) (interface{}, error) {
		wire := dynamicpb.NewMessage(desc.Input())
		if err := decode(wire); err != nil {
			return nil, err
		}
		// Respect the interceptor during SDK handler-type discovery as well as
		// normal gRPC calls; decoding/unpacking happens inside the handler.
		handler := func(ctx context.Context, input interface{}) (interface{}, error) {
			bz, err := proto.Marshal(input.(proto.Message))
			if err != nil {
				return nil, err
			}
			req := new(evmtypes.QueryTraceTxRequest)
			if err := req.Unmarshal(bz); err != nil {
				return nil, err
			}
			if err := unpackTraceTx(req, unpacker); err != nil {
				return nil, err
			}
			resp, err := server.(evmtypes.QueryServer).TraceTx(ctx, req)
			if err != nil {
				return nil, err
			}
			bz, err = resp.Marshal()
			if err != nil {
				return nil, err
			}
			output := dynamicpb.NewMessage(desc.Output())
			if err := proto.Unmarshal(bz, output); err != nil {
				return nil, err
			}
			return output, nil
		}
		if interceptor == nil {
			return handler(ctx, wire)
		}
		return interceptor(ctx, wire, &googlegrpc.UnaryServerInfo{
			Server: server, FullMethod: "/ethermint.evm.v1.Query/TraceTx",
		}, handler)
	}
}

func unpackTraceTx(req *evmtypes.QueryTraceTxRequest, unpacker types.AnyUnpacker) error {
	if req.Msg == nil {
		return status.Error(codes.InvalidArgument, "msg is required")
	}
	if err := req.Msg.UnpackInterfaces(unpacker); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid msg data: %v", err)
	}
	if _, err := evmtypes.UnpackTxData(req.Msg.Data); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid msg data: %v", err)
	}
	for i, msg := range req.Predecessors {
		if msg == nil {
			return status.Errorf(codes.InvalidArgument, "predecessors[%d] is required", i)
		}
		if err := msg.UnpackInterfaces(unpacker); err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid predecessors[%d] data: %v", i, err)
		}
		if _, err := evmtypes.UnpackTxData(msg.Data); err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid predecessors[%d] data: %v", i, err)
		}
	}
	return nil
}
