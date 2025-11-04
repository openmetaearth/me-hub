package keepers

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type MeMsg struct {
	MultiSend *MultiSend `json:"multi_send,omitempty"`
}

type MultiSend struct {
	Amount wasmvmtypes.Coins `json:"amount,omitempty"`
	Output []Output          `json:"output,omitempty"`
}

type Output struct {
	Address string            `json:"address,omitempty"`
	Amount  wasmvmtypes.Coins `json:"amount,omitempty"`
}

func EncodeMeMsg(sender sdk.AccAddress, msg *MeMsg) ([]sdk.Msg, error) {
	if msg.MultiSend == nil {
		return nil, errorsmod.Wrap(types.ErrUnknownMsg, "unknown variant of Bank")
	}
	if len(msg.MultiSend.Output) == 0 {
		return nil, nil
	}
	toSend, err := wasmkeeper.ConvertWasmCoinsToSdkCoins(msg.MultiSend.Amount)
	if err != nil {
		return nil, err
	}

	var outputs []banktypes.Output
	for _, o := range msg.MultiSend.Output {
		amt, err := wasmkeeper.ConvertWasmCoinsToSdkCoins(o.Amount)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, banktypes.Output{
			Address: o.Address,
			Coins:   amt,
		})
	}
	sdkMsg := banktypes.MsgMultiSend{
		Inputs: []banktypes.Input{
			{Address: sender.String(), Coins: toSend},
		},
		Outputs: outputs,
	}
	return []sdk.Msg{&sdkMsg}, nil
}

// SetupCustomMsgs sets up the custom message handlers for the app
func (a *AppKeepers) SetupCustomMsgs() wasmkeeper.Option {
	// Create KYC custom querier

	return wasmkeeper.WithMessageEncoders(&wasmkeeper.MessageEncoders{
		Custom: func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
			var meMsg MeMsg
			if err := json.Unmarshal(msg, &meMsg); err != nil {
				return nil, err
			}
			switch {
			case meMsg.MultiSend != nil:
				return EncodeMeMsg(sender, &meMsg)
			default:
				return nil, errorsmod.Wrap(types.ErrUnknownMsg, "unknown variant of MeMsg")
			}
		},
	})
}
