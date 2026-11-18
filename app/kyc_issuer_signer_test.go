package app_test

import (
	"testing"

	"cosmossdk.io/api/amino"
	msgv1 "cosmossdk.io/api/cosmos/msg/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/openmetaearth/me-hub/app"
	"github.com/openmetaearth/me-hub/testutil/sample"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	regiongovtypes "github.com/openmetaearth/me-hub/x/regiongov/types"
)

func TestKycIssuerSDKSignerAnnotations(t *testing.T) {
	encoding := app.MakeEncodingConfig()
	creator := sdk.MustAccAddressFromBech32(sample.AccAddress())
	msg := &kyctypes.MsgUpdateKycIssuer{Creator: creator.String(), Addresses: []string{sample.AccAddress()}, OptionType: 1}
	signers, reflected, err := encoding.Codec.GetMsgV1Signers(msg)
	require.NoError(t, err)
	require.Equal(t, [][]byte{creator.Bytes()}, signers)
	desc := reflected.ProtoReflect().Descriptor()
	require.Equal(t, "kyc/UpdateKycIssuer", proto.GetExtension(desc.Options(), amino.E_Name))
	require.Equal(t, true, proto.GetExtension(desc.ParentFile().Services().ByName("Msg").Options(), msgv1.E_Service))
	builder := encoding.TxConfig.NewTxBuilder()
	require.NoError(t, builder.SetMsgs(msg))
	wire, err := encoding.TxConfig.TxEncoder()(builder.GetTx())
	require.NoError(t, err)
	decoded, err := encoding.TxConfig.TxDecoder()(wire)
	require.NoError(t, err)
	require.Equal(t, msg, decoded.GetMsgs()[0])
}

func TestRegionGovMessagesDecodeBeforeAppConstruction(t *testing.T) {
	// MakeEncodingConfig is also used by standalone clients. Successful signer
	// reflection must not hide an unregistered concrete message in a TxBody Any.
	encoding := app.MakeEncodingConfig()
	signer := sample.AccAddress()
	for _, msg := range []sdk.Msg{
		&regiongovtypes.MsgSubmitProposal{Proposer: signer},
		&regiongovtypes.MsgVote{Voter: signer, ProposalId: 1, Option: regiongovtypes.VoteOptionYes},
		&regiongovtypes.MsgSubmitEquityUpdate{Proposer: signer},
		&regiongovtypes.MsgSubmitFreezeEquity{Proposer: signer},
		&regiongovtypes.MsgSubmitUnfreezeEquity{Proposer: signer},
		&regiongovtypes.MsgApprovePendingAction{Authority: signer},
		&regiongovtypes.MsgRejectPendingAction{Authority: signer},
		&regiongovtypes.MsgFreezeRegion{Authority: signer},
		&regiongovtypes.MsgUnfreezeRegion{Authority: signer},
		&regiongovtypes.MsgFreezeEquity{Authority: signer},
		&regiongovtypes.MsgUnfreezeEquity{Authority: signer},
	} {
		t.Run(sdk.MsgTypeURL(msg), func(t *testing.T) {
			builder := encoding.TxConfig.NewTxBuilder()
			require.NoError(t, builder.SetMsgs(msg))
			wire, err := encoding.TxConfig.TxEncoder()(builder.GetTx())
			require.NoError(t, err)
			decoded, err := encoding.TxConfig.TxDecoder()(wire)
			require.NoError(t, err)
			require.Equal(t, []sdk.Msg{msg}, decoded.GetMsgs())
		})
	}
}
