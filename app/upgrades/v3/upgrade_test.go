package v3

import (
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

func TestConvertOldRollappToNew_ChannelEnablesTransfers(t *testing.T) {
	owner := sample.AccAddress()
	old := rollapptypes.Rollapp{
		RollappId: "mecheckin_100-1",
		Owner:     owner,
		ChannelId: "channel-0",
	}

	got := ConvertOldRollappToNew(old, 42, rollapptypes.DefaultMinSequencerBondGlobalCoin)

	require.True(t, got.Launched)
	require.Equal(t, uint64(42), got.GenesisState.TransferProofHeight)
	require.True(t, got.GenesisState.IsTransferEnabled())
	require.True(t, got.GenesisInfo.Sealed)
	require.Equal(t, "me", got.GenesisInfo.Bech32Prefix)
	require.NotEmpty(t, got.GenesisInfo.GenesisChecksum)
	require.True(t, got.GenesisInfo.InitialSupply.IsZero())
	require.Equal(t, "*", got.InitialSequencer)
	require.Equal(t, rollapptypes.Rollapp_EVM, got.VmType)
	require.Equal(t, sdk.NewCoins(rollapptypes.DefaultMinSequencerBondGlobalCoin), got.MinSequencerBond)
	require.Equal(t, []rollapptypes.Revision{{Number: 0, StartHeight: 0}}, got.Revisions)
	require.NoError(t, got.ValidateBasic())
}

func TestConvertOldRollappToNew_ZeroLatestHeightUsesOne(t *testing.T) {
	old := rollapptypes.Rollapp{
		RollappId: "mecheckin_100-1",
		Owner:     sample.AccAddress(),
		ChannelId: "channel-0",
	}

	got := ConvertOldRollappToNew(old, 0, rollapptypes.DefaultMinSequencerBondGlobalCoin)
	require.Equal(t, uint64(1), got.GenesisState.TransferProofHeight)
}

func TestConvertOldRollappToNew_PreservesExistingProofHeight(t *testing.T) {
	old := rollapptypes.Rollapp{
		RollappId:    "mecheckin_100-1",
		Owner:        sample.AccAddress(),
		GenesisState: rollapptypes.RollappGenesisState{TransferProofHeight: 7},
		GenesisInfo: rollapptypes.GenesisInfo{
			GenesisChecksum: "abc",
			Bech32Prefix:    "me",
			InitialSupply:   math.ZeroInt(),
			Sealed:          true,
		},
		InitialSequencer: "*",
		MinSequencerBond: sdk.NewCoins(rollapptypes.DefaultMinSequencerBondGlobalCoin),
		VmType:           rollapptypes.Rollapp_WASM,
		Revisions:        []rollapptypes.Revision{{Number: 1, StartHeight: 10}},
	}

	got := ConvertOldRollappToNew(old, 99, rollapptypes.DefaultMinSequencerBondGlobalCoin)
	require.Equal(t, uint64(7), got.GenesisState.TransferProofHeight)
	require.True(t, got.Launched)
	require.Equal(t, rollapptypes.Rollapp_WASM, got.VmType)
	require.Equal(t, []rollapptypes.Revision{{Number: 1, StartHeight: 10}}, got.Revisions)
	require.NoError(t, got.ValidateBasic())
}

func TestConvertOldSequencerToNew_OptsIn(t *testing.T) {
	pk := ed25519.GenPrivKey().PubKey()
	pkAny, err := codectypes.NewAnyWithValue(pk)
	require.NoError(t, err)

	old := sequencertypes.Sequencer{
		Address:      sdk.AccAddress(pk.Address()).String(),
		DymintPubKey: pkAny,
		RollappId:    "mecheckin_100-1",
		Status:       sequencertypes.Bonded,
		Tokens:       sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(1))),
		OptedIn:      false,
		Metadata: sequencertypes.SequencerMetadata{
			Moniker: "old-seq",
			Details: "from v2",
		},
	}

	got := ConvertOldSequencerToNew(old)
	require.True(t, got.OptedIn)
	require.Equal(t, old.Address, got.Address)
	require.Equal(t, old.RollappId, got.RollappId)
	require.Equal(t, old.Status, got.Status)
	require.Equal(t, "old-seq", got.Metadata.Moniker)
	require.Equal(t, "10000000000", got.Metadata.GasPrice)
}
