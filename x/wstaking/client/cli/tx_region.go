package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/st-chain/me-hub/utils"
	"strings"

	"github.com/st-chain/me-hub/x/wstaking/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdNewRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-region [name] [validator]",
		Short: "Broadcast message new-region",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := strings.ToLower(args[0])
			argName := args[0]
			argVal := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := strings.Trim(argName, " ")
			if name != "" {
				name = strings.ToUpper(name)
				_, err = utils.CheckRegionName(name)
				if err != nil {
					return types.ErrRegionName.Wrap(err.Error())
				}
			}

			msg := types.NewMsgNewRegion(
				clientCtx.GetFromAddress().String(),
				argRegionId,
				name,
				argVal,
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

func CmdRemoveRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-region [region-id]",
		Short: "Broadcast message remove-region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveRegion(
				clientCtx.GetFromAddress().String(),
				argRegionId,
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

func CmdWithdrawFromRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-from-region [region-id] [receiver] [amount]",
		Short: "Send coins from region-treasury to receiver by global admin",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Send coins from region-treasury to receiver by global admin.
Example:
$ %s tx staking withdraw-from-region me_earth me1h47kmp4q5vkwjw350y5v5ecuzjtmct4zmrlhwf 100mec --from global-admin
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := args[0]
			argsReceiver := args[1]
			argsAmount := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(argsAmount)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawFromRegion(
				clientCtx.GetFromAddress().String(),
				argRegionId,
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
