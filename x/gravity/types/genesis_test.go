package types_test

import (
	"bytes"
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	_ "github.com/openmetaearth/me-hub/app/params"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
)

func TestGenesisStateValidateBasicAcceptsMaxUint64RelayerPower(t *testing.T) {
	delegateAmount := sdkmath.NewIntFromUint64(math.MaxUint64).Mul(sdk.DefaultPowerReduction)
	state := genesisWithRelayers(delegateAmount, []gravitytypes.Relayer{
		genesisRelayer(1, delegateAmount, true),
	})

	require.NoError(t, state.ValidateBasic())
}

func TestGenesisStateValidateBasicRejectsRelayerPowerOverflow(t *testing.T) {
	delegateAmount := sdkmath.NewIntFromUint64(math.MaxUint64).Add(sdkmath.OneInt()).Mul(sdk.DefaultPowerReduction)
	state := genesisWithRelayers(delegateAmount, []gravitytypes.Relayer{
		genesisRelayer(1, delegateAmount, true),
	})

	err := state.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "relayer power must fit uint64")
}

func TestGenesisStateValidateBasicRejectsTotalOnlineRelayerPowerOverflow(t *testing.T) {
	delegateAmount := sdkmath.NewIntFromUint64(math.MaxUint64).Mul(sdk.DefaultPowerReduction)
	state := genesisWithRelayers(delegateAmount, []gravitytypes.Relayer{
		genesisRelayer(1, delegateAmount, true),
		genesisRelayer(2, delegateAmount, true),
	})

	err := state.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "total online relayer power exceeds uint64")
}

func genesisWithRelayers(maxDelegate sdkmath.Int, relayers []gravitytypes.Relayer) *gravitytypes.GenesisState {
	params := gravitytypes.DefaultParams()
	params.MaxDelegate = maxDelegate
	return &gravitytypes.GenesisState{
		Params:   params,
		Relayers: relayers,
	}
}

func genesisRelayer(seed byte, delegateAmount sdkmath.Int, online bool) gravitytypes.Relayer {
	return gravitytypes.Relayer{
		RelayerAddress:  sdk.AccAddress(bytes.Repeat([]byte{seed}, 20)).String(),
		ExternalAddress: "0x0000000000000000000000000000000000000001",
		DelegateAmount:  delegateAmount,
		Online:          online,
	}
}
