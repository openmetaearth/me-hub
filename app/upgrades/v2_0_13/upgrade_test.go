package v2_0_13

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
	"github.com/stretchr/testify/require"
)

func TestFundGravityRelayerModuleAccounts(t *testing.T) {
	relayers := []string{
		"me1frjhlw9slyy7mrhmk0r4vytkyldxqtkf326amv",
		"me1c5zp26c0gq2klk87nrpff3y52u34zn4ydug2yd",
	}
	delegateCoins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1e8))
	bankKeeper := &recordingBankKeeper{}

	err := fundGravityRelayerModuleAccounts(sdk.Context{}, bankKeeper, relayers, delegateCoins)

	require.NoError(t, err)
	require.Equal(t, []recordedModuleTransfer{
		{
			sender: sdk.MustAccAddressFromBech32(relayers[0]),
			module: bsctypes.ModuleName,
			coins:  delegateCoins,
		},
		{
			sender: sdk.MustAccAddressFromBech32(relayers[0]),
			module: trontypes.ModuleName,
			coins:  delegateCoins,
		},
		{
			sender: sdk.MustAccAddressFromBech32(relayers[1]),
			module: bsctypes.ModuleName,
			coins:  delegateCoins,
		},
		{
			sender: sdk.MustAccAddressFromBech32(relayers[1]),
			module: trontypes.ModuleName,
			coins:  delegateCoins,
		},
	}, bankKeeper.transfers)
}

type recordingBankKeeper struct {
	transfers []recordedModuleTransfer
}

type recordedModuleTransfer struct {
	sender sdk.AccAddress
	module string
	coins  sdk.Coins
}

func (k *recordingBankKeeper) SendCoinsFromAccountToModule(_ sdk.Context, sender sdk.AccAddress, module string, coins sdk.Coins) error {
	k.transfers = append(k.transfers, recordedModuleTransfer{
		sender: sender,
		module: module,
		coins:  coins,
	})
	return nil
}
