package commands

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/distribution/types"
	"github.com/spf13/cobra"
)

// GetCmdClaimRewards returns the command to claim rewards
func GetCmdClaimRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-rewards [validator-addr]",
		Short: "Claim rewards for a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the validator address
			validatorAddr := args[0]

			// Create the message
			msg := types.NewMsgClaimRewards(validatorAddr)

			// Send the message
			_, err := ctx.Client.Tx.BroadcastTx(ctx.Client.Tx.GenerateTx(ctx.Client.From, msg))
			if err != nil {
				return err
			}

			fmt.Printf("Rewards claimed for validator %s\n", validatorAddr)

			return nil
		},
	}

	return cmd
}