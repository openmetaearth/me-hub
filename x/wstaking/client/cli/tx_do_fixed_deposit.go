package cli

import (
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func NewFixedDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-fixed [principal] [term]",
		Short: "Broadcast message deposit-fixed",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPrincipal, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			argTerm := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			term, err := strconv.ParseInt(argTerm, 10, 64)
			if err != nil {
				return types.ErrParameter.Wrap("term error")
			}

			msg := types.NewMsgDoFixedDeposit(
				clientCtx.GetFromAddress().String(),
				argPrincipal,
				term,
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

var _ = strconv.Itoa(0)

func NewFixedWithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-fixed [id]",
		Short: "Broadcast message withdraw_fixed_deposit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawFixedDeposit(
				clientCtx.GetFromAddress().String(),
				argId,
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
