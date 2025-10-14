package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/st-chain/me-hub/x/blacklist/types"
)

// CmdUpdateBlacklist implements the update-blacklist command
func CmdUpdateBlacklist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-blacklist",
		Short: "Update blacklist addresses (add/remove addresses)",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Get flag values
			removedAddrs, err := cmd.Flags().GetString("remove")
			if err != nil {
				return err
			}

			addedAddrs, err := cmd.Flags().GetString("add")
			if err != nil {
				return err
			}

			if removedAddrs == "" && addedAddrs == "" {
				return fmt.Errorf("at least one of --remove or --add must be provided")
			}

			// Parse addresses to remove
			var addressesToRemove []string
			if removedAddrs != "" {
				addressesToRemove = strings.Split(removedAddrs, ",")
				// Validate each address
				for _, addr := range addressesToRemove {
					if _, err := sdk.AccAddressFromBech32(addr); err != nil {
						return fmt.Errorf("address to remove %s is not a valid address: %w", addr, err)
					}
				}
			}

			// Parse addresses to add
			var addressesToAdd []string
			if addedAddrs != "" {
				addressesToAdd = strings.Split(addedAddrs, ",")
				// Validate each address
				for _, addr := range addressesToAdd {
					if _, err := sdk.AccAddressFromBech32(addr); err != nil {
						return fmt.Errorf("address to add %s is not a valid address: %w", addr, err)
					}
				}
			}

			msg := types.NewMsgUpdateBlacklist(
				clientCtx.GetFromAddress().String(),
				addressesToRemove,
				addressesToAdd,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("remove", "", "Comma-separated list of addresses to remove from blacklist")
	cmd.Flags().String("add", "", "Comma-separated list of addresses to add to blacklist")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
