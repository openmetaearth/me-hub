package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (k Keeper) Slashing(ctx sdk.Context, sequencerAddr sdk.AccAddress) error {
	// Get the sequencer's bonded tokens
	bondedTokens := k.GetBondedTokens(ctx, sequencerAddr)

	// Calculate the slashing amount (e.g., 10% of bonded tokens)
	slashAmount := bondedTokens.Quo(sdk.NewInt(10))

	// Burn the slashing amount
	if err := k.BurnTokens(ctx, sequencerAddr, slashAmount); err != nil {
		return err
	}

	return nil
}