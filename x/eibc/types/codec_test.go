package types_test

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/eibc/types"
)

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
