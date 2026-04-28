package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/spf13/cobra"
	"strings"
)

func NewUndelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undelegate [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "undelegate liquid tokens from a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`undelegate an amount of liquid coins from a validator to your wallet.
Example:
$ %s tx staking undelegate 1000mec true --from mykey
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			//isMeid := strings.ToLower(args[1]) == "true"
			delAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgUndelegate(delAddr, sdk.ValAddress{}, amount)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FlagSetValidatorAddress())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
