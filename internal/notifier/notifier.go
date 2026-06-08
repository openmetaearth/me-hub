package notifier

import (
	"context"
	"fmt"

	"github.com/openmetaearth/me-hub/types"
)

// NotifyReward sends a notification for a received reward
func (n *Notifier) NotifyReward(ctx context.Context, reward types.Reward) error {
	// Check if the reward has already been claimed
	if reward.Claimed {
		// Do not send a notification for a previously claimed reward
		return nil
	}

	// Send the notification
	err := n.sendNotification(ctx, reward)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

// sendNotification sends a notification for a received reward
func (n *Notifier) sendNotification(ctx context.Context, reward types.Reward) error {
	// Implement notification sending logic here
	// ...
}