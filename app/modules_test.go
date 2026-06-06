package app

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
)

func TestGenTxMessageValidator(t *testing.T) {
	createValidator := validCreateValidatorMsg(t)
	newRegion := wstakingtypes.NewMsgNewRegion(createValidator.DelegatorAddress, "usa", createValidator.ValidatorAddress)

	tests := []struct {
		name    string
		msgs    []sdk.Msg
		wantErr bool
	}{
		{
			name:    "empty gentx",
			msgs:    nil,
			wantErr: true,
		},
		{
			name:    "single valid create validator",
			msgs:    []sdk.Msg{createValidator},
			wantErr: false,
		},
		{
			name:    "valid create validator and region",
			msgs:    []sdk.Msg{createValidator, newRegion},
			wantErr: false,
		},
		{
			name:    "malformed create validator",
			msgs:    []sdk.Msg{&stakingtypes.MsgCreateValidator{}},
			wantErr: true,
		},
		{
			name:    "first message is not create validator",
			msgs:    []sdk.Msg{validMsgSend()},
			wantErr: true,
		},
		{
			name:    "arbitrary second message",
			msgs:    []sdk.Msg{createValidator, validMsgSend()},
			wantErr: true,
		},
		{
			name:    "malformed second region message",
			msgs:    []sdk.Msg{createValidator, wstakingtypes.NewMsgNewRegion("", "usa", createValidator.ValidatorAddress)},
			wantErr: true,
		},
		{
			name:    "too many gentx messages",
			msgs:    []sdk.Msg{createValidator, newRegion, validMsgSend()},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := GenTxMessageValidator(tc.msgs)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func validCreateValidatorMsg(t *testing.T) *stakingtypes.MsgCreateValidator {
	t.Helper()

	privKey := secp256k1.GenPrivKey()
	valAddr := sdk.ValAddress(privKey.PubKey().Address())

	msg, err := stakingtypes.NewMsgCreateValidator(
		valAddr,
		privKey.PubKey(),
		sdk.NewCoin(params.BaseDenom, math.NewInt(100_000_000)),
		stakingtypes.Description{Moniker: "validator"},
		stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		math.OneInt(),
	)
	require.NoError(t, err)
	require.NoError(t, msg.ValidateBasic())

	return msg
}

func validMsgSend() *banktypes.MsgSend {
	from := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	to := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	return banktypes.NewMsgSend(from, to, sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1)))
}
