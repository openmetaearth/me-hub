package cli

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/openmetaearth/me-hub/utils"

	"github.com/openmetaearth/me-hub/x/wstaking/types"

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

func CmdTransferRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-region  [from-region] [to-region] [address]",
		Short: "transfer meid to new region.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := args[0]
			to := args[1]
			addr := args[2]
			address := make([]string, 0, 1)
			address = append(address, addr)
			creatorAddr := clientCtx.GetFromAddress().String()
			msg := &types.MsgTransferRegion{FromRegion: from, ToRegion: to, Address: address, Creator: creatorAddr}

			//if err = msg.ValidateBasic(); err != nil {
			//	return err
			//}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdGrantRegionWithdraw returns the CLI command for granting (or
// updating) withdraw for a region treasury.
func CmdGrantRegionWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant-region-withdraw [region-id] [address]",
		Short: "Grant (or update) an address to withdraw from a region treasury",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Grant an address to withdraw from a specific region treasury.
If the region already has a granted address, it will be overwritten.
Only GlobalDao can call this.

Example:
$ %s tx staking grant-region-withdraw usa me1abc...xyz --from global_dao
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			regionId := args[0]
			address := args[1]

			msg := types.NewMsgGrantRegionWithdraw(
				clientCtx.GetFromAddress().String(),
				regionId,
				address,
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

// CmdRevokeRegionWithdraw returns the CLI command for revoking
// withdraw for a region treasury.
func CmdRevokeRegionWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-region-withdraw [region-id]",
		Short: "Revoke the withdrawer for a region treasury",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Revoke the withdrawer previously granted for a region treasury.
Only GlobalDao can call this.

Example:
$ %s tx staking revoke-region-withdraw usa --from global_dao
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			regionId := args[0]

			msg := types.NewMsgRevokeRegionWithdraw(
				clientCtx.GetFromAddress().String(),
				regionId,
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
