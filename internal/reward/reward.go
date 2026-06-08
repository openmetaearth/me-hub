package reward

import (
	"context"
	"fmt"

	"github.com/openmetaearth/me-hub/types"
)

// GetReward retrieves a reward by ID
func (r *RewardManager) GetReward(ctx context.Context, rewardID string) (*types.Reward, error) {
	// Retrieve the reward from the database
	reward, err := r.db.GetReward(ctx, rewardID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve reward: %w", err)
	}

	// Check if the reward has already been claimed
	if reward.Claimed {
		// Return a nil reward to indicate that the reward has already been claimed
		return nil, nil
	}

	return reward, nil
}