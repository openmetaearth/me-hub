package app

import (
	"fmt"

	"cosmossdk.io/core/address"
	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/gogoproto/proto"
	protov2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

// registerWstakingCustomGetSigners fills in SIGN_MODE_DIRECT signers for wstaking
// msgs that do not yet have (cosmos.msg.v1.signer) in the generated file
// descriptor. Remove these entries after proto-gen picks up the proto options.
func registerWstakingCustomGetSigners(opts *signing.Options, addrCodec address.Codec) {
	byField := func(field protoreflect.Name) signing.GetSignersFunc {
		return protoFieldGetSigners(field, addrCodec)
	}

	for _, item := range []struct {
		msg   proto.Message
		field protoreflect.Name
	}{
		{(*wstakingtypes.MsgNewRegion)(nil), "creator"},
		{(*wstakingtypes.MsgRemoveRegion)(nil), "creator"},
		{(*wstakingtypes.MsgWithdrawFromRegion)(nil), "withdrawer"},
		{(*wstakingtypes.MsgWithdrawFromGlobalDaoFeePool)(nil), "withdrawer"},
		{(*wstakingtypes.MsgNewFixedDepositCfg)(nil), "dao"},
		{(*wstakingtypes.MsgSetFixedDepositCfgStatus)(nil), "admin"},
		{(*wstakingtypes.MsgSetFixedDepositCfgRate)(nil), "admin"},
		{(*wstakingtypes.MsgRemoveFixedDepositCfg)(nil), "admin"},
		{(*wstakingtypes.MsgDoFixedDeposit)(nil), "account"},
		{(*wstakingtypes.MsgWithdrawFixedDeposit)(nil), "account"},
		{(*wstakingtypes.MsgNewRecord)(nil), "from"},
		{(*wstakingtypes.MsgReviewRecord)(nil), "from"},
		{(*wstakingtypes.MsgTransferRegion)(nil), "creator"},
		{(*wstakingtypes.MsgIbcTransferFromRegionTreasure)(nil), "creator"},
		{(*wstakingtypes.MsgReplaceConsensusPubKeyRequest)(nil), "creator"},
		{(*wstakingtypes.MsgSendToModule)(nil), "sender"},
	} {
		opts.DefineCustomGetSigners(protoreflect.FullName(proto.MessageName(item.msg)), byField(item.field))
	}
}

func protoFieldGetSigners(field protoreflect.Name, addrCodec address.Codec) signing.GetSignersFunc {
	return func(msg protov2.Message) ([][]byte, error) {
		m := msg.ProtoReflect()
		fd := m.Descriptor().Fields().ByName(field)
		if fd == nil {
			return nil, fmt.Errorf("message %s has no %s field", m.Descriptor().FullName(), field)
		}
		bz, err := addrCodec.StringToBytes(m.Get(fd).String())
		if err != nil {
			return nil, fmt.Errorf("invalid %s on %s: %w", field, m.Descriptor().FullName(), err)
		}
		return [][]byte{bz}, nil
	}
}
