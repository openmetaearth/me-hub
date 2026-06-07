package types

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestPowerReduction(t *testing.T) {
	t.Log(sdk.DefaultPowerReduction.String())

}

func TestValidateGenesisRegionWithdrawGrants(t *testing.T) {
	validAddr := sdk.AccAddress([]byte("region_withdrawer_01")).String()

	validGenesis := func() *GenesisState {
		genesis := DefaultGenesisState()
		genesis.Regions = []Region{{RegionId: ExperienceRegionId}}
		genesis.RegionWithdrawGrants = []RegionWithdrawGrant{{
			RegionId: ExperienceRegionId,
			Address:  validAddr,
		}}
		return genesis
	}

	require.NoError(t, ValidateGenesis(validGenesis()))
	encoded, err := validGenesis().Marshal()
	require.NoError(t, err)
	var decoded GenesisState
	require.NoError(t, decoded.Unmarshal(encoded))
	require.Equal(t, validGenesis().RegionWithdrawGrants, decoded.RegionWithdrawGrants)

	registry := codectypes.NewInterfaceRegistry()
	RegisterInterfaces(registry)
	stakingtypes.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	jsonGenesis := cdc.MustMarshalJSON(validGenesis())
	require.True(t, strings.Contains(string(jsonGenesis), "regionWithdrawGrants"))
	var decodedJSON GenesisState
	cdc.MustUnmarshalJSON(jsonGenesis, &decodedJSON)
	require.Equal(t, validGenesis().RegionWithdrawGrants, decodedJSON.RegionWithdrawGrants)

	duplicate := validGenesis()
	duplicate.RegionWithdrawGrants = append(duplicate.RegionWithdrawGrants, duplicate.RegionWithdrawGrants[0])
	require.ErrorContains(t, ValidateGenesis(duplicate), "duplicate region withdraw grant")

	unknownRegion := validGenesis()
	unknownRegion.RegionWithdrawGrants[0].RegionId = "unknown"
	require.ErrorContains(t, ValidateGenesis(unknownRegion), "unknown region")

	invalidAddress := validGenesis()
	invalidAddress.RegionWithdrawGrants[0].Address = "not-a-valid-address"
	require.ErrorContains(t, ValidateGenesis(invalidAddress), "invalid region withdraw grant address")
}
