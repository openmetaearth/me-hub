package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
)

func CmdWithdrawFromTreasury() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-from-treasury [receiver] [amount]",
		Short: "Send coins from treasury to receiver by global dao",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Send coins from treasury to receiver by global dao.
Example:
$ %s tx staking withdraw-from-treasury me1h47kmp4q5vkwjw350y5v5ecuzjtmct4zmrlhwf 100mec --from global-dao
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argsReceiver := args[0]
			argsAmount := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(argsAmount)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawFromTreasury(
				clientCtx.GetFromAddress().String(),
				argsReceiver,
				amount,
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
