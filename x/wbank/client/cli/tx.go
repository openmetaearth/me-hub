package cli

import (
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/wbank/types"
)

const (
	FlagRegionID = "region-id"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        banktypes.ModuleName,
		Short:                      "Bank transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		bankcli.NewSendTxCmd(),
		bankcli.NewMultiSendTxCmd(),
		NewWithdrawTreasuryTxCmd(),
		NewSendToTreasuryTxCmd(),
		NewSendToAirdropTxCmd(),
	)

	return txCmd
}

// NewSendToGlobalDaoTxCmd returns a CLI command handler for creating a MsgSend transaction.
func NewWithdrawTreasuryTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-treasury [receiver] [amount]",
		Short: "Send funds from treasury to global dao.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Send funds from treasury to global dao.
Example:
$ %s tx bank send-global-dao 1000mec --from global-admin
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

			receiver, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawTreasury(clientCtx.GetFromAddress(), receiver, coins)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewSendToTreasuryTxCmd returns a CLI command handler for creating a MsgSendToTreasury transaction.
func NewSendToTreasuryTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-to-treasury [amount]",
		Short: "Send funds from global dao to treasury.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Send funds from global dao to treasury.
Example:
$ %s tx bank send-to-treasury 1000mec --from global-dao
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

			coins, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSendToTreasury(clientCtx.GetFromAddress(), coins)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewSendToAirdropTxCmd returns a CLI command handler for creating a MsgSendToAirdrop transaction.
func NewSendToAirdropTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendToAirdrop [amount]",
		Short: "Send funds from region treasury to airdrop address.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Send funds from region treasury to airdrop address.

Example:
$ %s tx bank sendToAirdrop 1000mec --from global-admin --region-id me_earth
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

			coins, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			regionID, _ := cmd.Flags().GetString(FlagRegionID)
			msg := types.NewMsgSendToAirdrop(clientCtx.GetFromAddress(), regionID, coins)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagRegionID, "", "example: --region-id me_earth")
	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(FlagRegionID)
	return cmd
}
