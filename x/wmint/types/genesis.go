package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

const (
	GenesisMintedCoinAmountField       = "minted_coin_amount"
	GenesisPerBlockMintCoinAmountField = "per_block_mint_coin_amount"
)

type GenesisState struct {
	Minter                 minttypes.Minter `json:"minter"`
	Params                 minttypes.Params `json:"params"`
	MintedCoinAmount       string           `json:"minted_coin_amount,omitempty"`
	PerBlockMintCoinAmount string           `json:"per_block_mint_coin_amount,omitempty"`
}

func DefaultGenesisState() GenesisState {
	mintGenesis := minttypes.DefaultGenesisState()
	return NewGenesisState(mintGenesis, *big.NewInt(0), *big.NewInt(0))
}

func NewGenesisState(mintGenesis *minttypes.GenesisState, mintedAmount, perBlockAmount big.Int) GenesisState {
	return GenesisState{
		Minter:                 mintGenesis.Minter,
		Params:                 mintGenesis.Params,
		MintedCoinAmount:       mintedAmount.String(),
		PerBlockMintCoinAmount: perBlockAmount.String(),
	}
}

func (g GenesisState) MintGenesisState() minttypes.GenesisState {
	return minttypes.GenesisState{
		Minter: g.Minter,
		Params: g.Params,
	}
}

func ParseGenesisAmount(field, value string) (big.Int, error) {
	if value == "" {
		return *big.NewInt(0), nil
	}

	amount, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return big.Int{}, fmt.Errorf("%s must be a base-10 integer", field)
	}
	if amount.Sign() < 0 {
		return big.Int{}, fmt.Errorf("%s must not be negative", field)
	}

	return *amount, nil
}

func ValidateGenesis(data GenesisState) error {
	if err := minttypes.ValidateGenesis(data.MintGenesisState()); err != nil {
		return err
	}

	if _, err := ParseGenesisAmount(GenesisMintedCoinAmountField, data.MintedCoinAmount); err != nil {
		return err
	}
	if _, err := ParseGenesisAmount(GenesisPerBlockMintCoinAmountField, data.PerBlockMintCoinAmount); err != nil {
		return err
	}

	return nil
}

func MarshalGenesis(data GenesisState) (json.RawMessage, error) {
	return json.Marshal(data)
}

func MustMarshalGenesis(data GenesisState) json.RawMessage {
	bz, err := MarshalGenesis(data)
	if err != nil {
		panic(err)
	}
	return bz
}

func UnmarshalGenesis(bz json.RawMessage) (GenesisState, error) {
	var data GenesisState
	if err := json.Unmarshal(bz, &data); err != nil {
		return GenesisState{}, err
	}
	return data, nil
}
