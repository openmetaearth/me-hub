package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	testkeeper "github.com/openmetaearth/me-hub/testutil/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

func TestKeeper_IsZkContractAddress(t *testing.T) {
	k, ctx := testkeeper.DidKeeper(t)
	stateAddress := common.HexToAddress("0x0000000000000000000000000000000000001001")
	authAddress := common.HexToAddress("0x0000000000000000000000000000000000001002")
	otherAddress := common.HexToAddress("0x0000000000000000000000000000000000002001")

	require.False(t, k.IsZkContractAddress(ctx, stateAddress))

	k.SetZkContractAddressInfo(ctx, didtypes.ZkContractAddressInfo{
		Type:            didtypes.PROXY_STATE,
		ContractAddress: stateAddress.Hex(),
	})
	k.SetZkContractAddressInfo(ctx, didtypes.ZkContractAddressInfo{
		Type:            didtypes.PROXY_AUTH_V3_VALIDATOR,
		ContractAddress: authAddress.Hex(),
	})

	require.True(t, k.IsZkContractAddress(ctx, stateAddress))
	require.True(t, k.IsZkContractAddress(ctx, authAddress))
	require.False(t, k.IsZkContractAddress(ctx, otherAddress))
}
