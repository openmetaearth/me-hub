package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/distribution/types"
)

// ClaimRewards allows a validator to claim their outstanding rewards
func (k Keeper) ClaimRewards(ctx context.Context, validatorAddr string) (*types.Reward, error) {
	// Get the validator's rewards
	rewards, found := k.GetValidatorRewards(ctx, validatorAddr)
	if !found {
		return nil, fmt.Errorf("validator %s has no rewards", validatorAddr)
	}

	// Claim the rewards
	err := k.ClaimValidatorRewards(ctx, validatorAddr, rewards)
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

// ClaimValidatorRewards claims the rewards for a validator
func (k Keeper) ClaimValidatorRewards(ctx context.Context, validatorAddr string, rewards *types.Reward) error {
	// Update the validator's rewards
	k.SetValidatorRewards(ctx, validatorAddr, nil)

	// Send the rewards to the validator
	err := k.SendRewards(ctx, validatorAddr, rewards)
	if err != nil {
		return err
	}

	return nil
}

// SendRewards sends the rewards to the validator
func (k Keeper) SendRewards(ctx context.Context, validatorAddr string, rewards *types.Reward) error {
	// Send the rewards
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, validatorAddr, rewards.Coins)
	if err != nil {
		return err
	}

	return nil
}