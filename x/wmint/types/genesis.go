package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
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
	var zero big.Int
	return NewGenesisState(mintGenesis, zero, zero)
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

func MarshalGenesis(cdc codec.JSONCodec, data GenesisState) (json.RawMessage, error) {
	mintGenesis := data.MintGenesisState()
	mintGenesisBz, err := cdc.MarshalJSON(&mintGenesis)
	if err != nil {
		return nil, err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(mintGenesisBz, &raw); err != nil {
		return nil, err
	}
	if raw == nil {
		raw = make(map[string]json.RawMessage)
	}

	raw[GenesisMintedCoinAmountField], err = json.Marshal(data.MintedCoinAmount)
	if err != nil {
		return nil, err
	}
	raw[GenesisPerBlockMintCoinAmountField], err = json.Marshal(data.PerBlockMintCoinAmount)
	if err != nil {
		return nil, err
	}

	return json.Marshal(raw)
}

func MustMarshalGenesis(cdc codec.JSONCodec, data GenesisState) json.RawMessage {
	bz, err := MarshalGenesis(cdc, data)
	if err != nil {
		panic(err)
	}
	return bz
}

func UnmarshalGenesis(cdc codec.JSONCodec, bz json.RawMessage) (GenesisState, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(bz, &raw); err != nil {
		return GenesisState{}, err
	}

	mintGenesisBz, err := json.Marshal(map[string]json.RawMessage{
		"minter": raw["minter"],
		"params": raw["params"],
	})
	if err != nil {
		return GenesisState{}, err
	}

	var mintGenesis minttypes.GenesisState
	if err := cdc.UnmarshalJSON(mintGenesisBz, &mintGenesis); err != nil {
		return GenesisState{}, err
	}

	var data GenesisState
	data.Minter = mintGenesis.Minter
	data.Params = mintGenesis.Params
	if bz, ok := raw[GenesisMintedCoinAmountField]; ok {
		if err := json.Unmarshal(bz, &data.MintedCoinAmount); err != nil {
			return GenesisState{}, err
		}
	}
	if bz, ok := raw[GenesisPerBlockMintCoinAmountField]; ok {
		if err := json.Unmarshal(bz, &data.PerBlockMintCoinAmount); err != nil {
			return GenesisState{}, err
		}
	}

	return data, nil
}
