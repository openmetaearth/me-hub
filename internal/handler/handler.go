package handler

import (
	"context"
	"fmt"

	"github.com/openmetaearth/me-hub/types"
)

// HandleReward handles a received reward
func (h *Handler) HandleReward(ctx context.Context, reward types.Reward) error {
	// Check if the reward has already been claimed
	if reward.Claimed {
		// Do not process a previously claimed reward
		return nil
	}

	// Process the reward
	err := h.processReward(ctx, reward)
	if err != nil {
		return fmt.Errorf("failed to process reward: %w", err)
	}

	return nil
}

// processReward processes a received reward
func (h *Handler) processReward(ctx context.Context, reward types.Reward) error {
	// Implement reward processing logic here
	// ...
}