package types_test

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/eibc/types"
)

func TestRegisterCodecIncludesMsgUpdateDemandOrder(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	types.RegisterCodec(cdc)

	msg := sdk.Msg(types.NewMsgUpdateDemandOrder(
		"cosmos18wvvwfmq77a6d8tza4h5sfuy2yj3jj88yqg82a",
		strings.Repeat("a", 64),
		"1",
	))

	bz, err := cdc.MarshalJSON(msg)
	require.NoError(t, err)
	require.Contains(t, string(bz), "eibc/MsgUpdateDemandOrder")

	var decoded sdk.Msg
	require.NoError(t, cdc.UnmarshalJSON(bz, &decoded))
	require.IsType(t, &types.MsgUpdateDemandOrder{}, decoded)
}

func TestRegisterInterfacesIncludesMsgUpdateDemandOrder(t *testing.T) {
	registry := cdctypes.NewInterfaceRegistry()
	sdk.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)

	require.Contains(
		t,
		registry.ListImplementations(sdk.MsgInterfaceProtoName),
		sdk.MsgTypeURL(&types.MsgUpdateDemandOrder{}),
	)
}
