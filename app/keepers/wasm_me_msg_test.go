package keepers

import (
	"bytes"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	_ "github.com/openmetaearth/me-hub/app/params"
)

func TestEncodeMeMsgMultiSendValidation(t *testing.T) {
	sender := testAccAddress(1)
	recipientA := testAccAddress(2)
	recipientB := testAccAddress(3)

	tests := []struct {
		name    string
		msg     *MeMsg
		wantErr string
	}{
		{
			name: "valid multisend",
			msg: &MeMsg{MultiSend: &MultiSend{
				Amount: wasmvmtypes.Coins{testWasmCoin("15")},
				Output: []Output{
					{Address: recipientA.String(), Amount: wasmvmtypes.Coins{testWasmCoin("10")}},
					{Address: recipientB.String(), Amount: wasmvmtypes.Coins{testWasmCoin("5")}},
				},
			}},
		},
		{
			name: "invalid output address",
			msg: &MeMsg{MultiSend: &MultiSend{
				Amount: wasmvmtypes.Coins{testWasmCoin("10")},
				Output: []Output{
					{Address: "not-an-address", Amount: wasmvmtypes.Coins{testWasmCoin("10")}},
				},
			}},
			wantErr: "invalid multisend output address",
		},
		{
			name: "module account output address",
			msg: &MeMsg{MultiSend: &MultiSend{
				Amount: wasmvmtypes.Coins{testWasmCoin("10")},
				Output: []Output{
					{Address: authtypes.NewModuleAddress(govtypes.ModuleName).String(), Amount: wasmvmtypes.Coins{testWasmCoin("10")}},
				},
			}},
			wantErr: "module account",
		},
		{
			name: "non-positive output amount",
			msg: &MeMsg{MultiSend: &MultiSend{
				Amount: wasmvmtypes.Coins{testWasmCoin("10")},
				Output: []Output{
					{Address: recipientA.String(), Amount: wasmvmtypes.Coins{testWasmCoin("0")}},
				},
			}},
			wantErr: "multisend output amount must be positive",
		},
		{
			name: "input output mismatch",
			msg: &MeMsg{MultiSend: &MultiSend{
				Amount: wasmvmtypes.Coins{testWasmCoin("10")},
				Output: []Output{
					{Address: recipientA.String(), Amount: wasmvmtypes.Coins{testWasmCoin("9")}},
				},
			}},
			wantErr: "does not match output amount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgs, err := EncodeMeMsg(sender, tt.msg)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				require.Nil(t, msgs)
				return
			}

			require.NoError(t, err)
			require.Len(t, msgs, 1)

			multiSend, ok := msgs[0].(*banktypes.MsgMultiSend)
			require.True(t, ok)
			require.NoError(t, multiSend.ValidateBasic())
			require.Equal(t, sender.String(), multiSend.Inputs[0].Address)
			require.Equal(t, recipientA.String(), multiSend.Outputs[0].Address)
			require.Equal(t, recipientB.String(), multiSend.Outputs[1].Address)
		})
	}
}

func TestEncodeMeMsgKeepsEmptyOutputNoop(t *testing.T) {
	msgs, err := EncodeMeMsg(testAccAddress(1), &MeMsg{MultiSend: &MultiSend{}})
	require.NoError(t, err)
	require.Nil(t, msgs)
}

func testAccAddress(fill byte) sdk.AccAddress {
	return sdk.AccAddress(bytes.Repeat([]byte{fill}, 20))
}

func testWasmCoin(amount string) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{Denom: "umec", Amount: amount}
}
