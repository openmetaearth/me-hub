package cli

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dymensionxyz/dymension/v3/x/rollup/types"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.MODULE_NAME,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.MODULE_NAME),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdStakeForSequencer())
	cmd.AddCommand(CmdUnStake())

	return cmd
}

func CmdStakeForSequencer() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "StakeForSequencer [creator] [rollappId] [amount]",
		Short:   "StakeForSequencer",
		Example: "dymd tx HUB_ROLLUP StakeForSequencer <creator> <rollappId> <amount>",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			creator := args[0]
			rollappID := args[1]
			amount := args[2]
			val, err := strconv.ParseUint(amount, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgSeqStaking(creator, rollappID, 0, val)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUnStake() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "UnStake [creator] [rollappId] [amount]",
		Short:   "unstake mec",
		Example: "dymd tx HUB_ROLLUP UnStake <creator> <rollappId> <amount>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			creator := args[0]
			rollappID := args[1]
			amount := args[2]
			val, err := strconv.ParseUint(amount, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgSeqUnStaking(creator, rollappID, 0, val)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
