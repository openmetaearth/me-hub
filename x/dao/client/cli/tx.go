package cli

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/dao/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdUpdateGlobalDao())
	cmd.AddCommand(CmdFreeGasAccount())
	return cmd
}

func CmdUpdateGlobalDao() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-dao [GlobalDao] [MeidDao] [DevOperator] [AirdropAddress]",
		Short: "Broadcast message",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			daoAddresses := types.DaoAddresses{
				GlobalDao:      args[0],
				MeidDao:        args[1],
				DevOperator:    args[2],
				AirdropAddress: args[3],
			}
			msg := types.NewMsgUpdateDao(
				clientCtx.GetFromAddress(),
				daoAddresses,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdFreeGasAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "free-gas-account [accounts-json]",
		Short: "Broadcast message to set free gas accounts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse the JSON input into a slice of FreeGasAccount structs
			var accounts []types.FreeGasAccount
			if err := json.Unmarshal([]byte(args[0]), &accounts); err != nil {
				return fmt.Errorf("invalid JSON input: %w", err)
			}

			msg := types.NewMsgFreeGasAccount(
				clientCtx.GetFromAddress(),
				accounts,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
