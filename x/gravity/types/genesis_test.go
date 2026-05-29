package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_ValidateBasic_BridgeTokenDuplicates(t *testing.T) {
	validToken := BridgeToken{
		ContractAddress: "0x0000000000000000000000000000000000000001",
		Denom:           "uusdt",
		Name:            "Tether USD",
		Symbol:          "USDT",
		Decimal:         6,
	}
	validToken2 := BridgeToken{
		ContractAddress: "0x0000000000000000000000000000000000000002",
		Denom:           "uusdc",
		Name:            "USDC Coin",
		Symbol:          "USDC",
		Decimal:         6,
	}

	baseGenesis := func(tokens []BridgeToken) *GenesisState {
		return &GenesisState{
			Params:       DefaultParams(),
			BridgeTokens: tokens,
		}
	}

	t.Run("valid: no duplicates", func(t *testing.T) {
		gs := baseGenesis([]BridgeToken{validToken, validToken2})
		require.NoError(t, gs.ValidateBasic())
	})

	t.Run("valid: single token", func(t *testing.T) {
		gs := baseGenesis([]BridgeToken{validToken})
		require.NoError(t, gs.ValidateBasic())
	})

	t.Run("valid: empty bridge tokens", func(t *testing.T) {
		gs := baseGenesis(nil)
		require.NoError(t, gs.ValidateBasic())
	})

	t.Run("invalid: duplicate denom", func(t *testing.T) {
		dup := BridgeToken{
			ContractAddress: "0x0000000000000000000000000000000000000003",
			Denom:           "uusdt", // same denom as validToken
			Name:            "Fake USD",
			Symbol:          "fUSDT",
			Decimal:         6,
		}
		gs := baseGenesis([]BridgeToken{validToken, dup})
		err := gs.ValidateBasic()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrDuplicate)
		require.Contains(t, err.Error(), "uusdt")
	})

	t.Run("invalid: duplicate contract address", func(t *testing.T) {
		dup := BridgeToken{
			ContractAddress: "0x0000000000000000000000000000000000000001", // same contract as validToken
			Denom:           "udai",
			Name:            "Dai",
			Symbol:          "DAI",
			Decimal:         18,
		}
		gs := baseGenesis([]BridgeToken{validToken, dup})
		err := gs.ValidateBasic()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrDuplicate)
		require.Contains(t, err.Error(), "0x0000000000000000000000000000000000000001")
	})

	t.Run("invalid: empty denom", func(t *testing.T) {
		dup := BridgeToken{
			ContractAddress: "0x0000000000000000000000000000000000000003",
			Denom:           "",
		}
		gs := baseGenesis([]BridgeToken{dup})
		err := gs.ValidateBasic()
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty denom")
	})

	t.Run("invalid: empty contract address", func(t *testing.T) {
		dup := BridgeToken{
			ContractAddress: "",
			Denom:           "uusdt",
		}
		gs := baseGenesis([]BridgeToken{dup})
		err := gs.ValidateBasic()
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty contract address")
	})

	t.Run("invalid: 3-way split - two denoms same, three unique contracts", func(t *testing.T) {
		// This is the exact attack: same denom, different contracts.
		// InitGenesis would overwrite the denom index for the second entry
		// but leave the first contract's index pointing to a now-orphaned denom.
		dup := BridgeToken{
			ContractAddress: "0x0000000000000000000000000000000000000003",
			Denom:           "uusdt", // same denom as validToken
		}
		gs := baseGenesis([]BridgeToken{validToken, validToken2, dup})
		err := gs.ValidateBasic()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrDuplicate)
	})
}
