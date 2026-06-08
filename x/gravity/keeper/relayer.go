package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (k Keeper) SlashRelayer(ctx sdk.Context, relayerAddr sdk.AccAddress) error {
	// Get the relayer's bonded tokens
	bondedTokens := k.GetBondedTokens(ctx, relayerAddr)

	// Burn the bonded tokens
	if err := k.BurnTokens(ctx, relayerAddr, bondedTokens); err != nil {
		return err
	}

	// Increment the slash times and set online to false
	k.SetSlashTimes(ctx, relayerAddr, k.GetSlashTimes(ctx, relayerAddr)+1)
	k.SetOnline(ctx, relayerAddr, false)

	return nil
}