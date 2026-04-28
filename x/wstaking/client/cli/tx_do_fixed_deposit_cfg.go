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
	"strconv"
	"strings"
)

func CmdNewFixedDepositCfg() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new-fixed-deposit-cfg [regionId] [term] [rate]",
		Short:   "Broadcast message new fixed deposit config",
		Example: fmt.Sprintf("%s tx staking new-fixed-deposit-cfg me_earth 1 0.1", version.AppName),
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := args[0]
			argTerm := args[1]
			argRate := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			term, err := strconv.ParseInt(argTerm, 10, 64)
			if err != nil {
				return types.ErrParameter.Wrap("term error")
			}

			rate, err := sdk.NewDecFromStr(argRate)
			if err != nil {
				return types.ErrParameter.Wrap("rate error")
			}

			msg := types.NewMsgNewFixedDepositCfg(
				clientCtx.GetFromAddress().String(),
				argRegionId,
				term,
				rate,
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

func CmdRemoveFixedDepositCfg() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove-fixed-deposit-cfg [regionId] [term]",
		Short:   "Broadcast message remove-fixed-deposit-cfg",
		Example: fmt.Sprintf("%s tx staking remove-fixed-deposit-cfg me_earth 1", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := args[0]
			argTerm := args[1]

			term, err := strconv.ParseInt(argTerm, 10, 64)
			if err != nil {
				return types.ErrParameter.Wrap("term error")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveFixedDepositCfg(
				clientCtx.GetFromAddress().String(),
				argRegionId,
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

func CmdSetFixedDepositCfgStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-fixed-deposit-cfg-status [regionId] [term] [status]",
		Short:   "Broadcast message set-fixed-deposit-cfg-status",
		Example: fmt.Sprintf("%s tx staking set-fixed-deposit-cfg-status me_earth 1 FIXED_DEPOSIT_CFG_INACTIVE", version.AppName),
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := args[0]
			argTerm := args[1]
			argStatus := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			term, err := strconv.ParseInt(argTerm, 10, 64)
			if err != nil {
				return types.ErrParameter.Wrap("term error")
			}

			status, ok := types.FIXED_DEPOSIT_CFG_STATUS_value[strings.ToUpper(strings.Trim(argStatus, " "))]
			if !ok {
				return types.ErrParameter.Wrap("period error")
			}

			msg := types.NewMsgSetFixedDepositCfgStatus(
				clientCtx.GetFromAddress().String(),
				argRegionId,
				term,
				types.FIXED_DEPOSIT_CFG_STATUS(status),
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

func CmdSetFixedDepositCfgRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-fixed-deposit-cfg-rate [regionId] [term] [rate]",
		Short:   "Broadcast message set-fixed-deposit-cfg-rate",
		Example: fmt.Sprintf("%s tx staking set-fixed-deposit-cfg-rate me_earth 1 0.1", version.AppName),
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRegionId := args[0]
			argTerm := args[1]
			argRate := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			term, err := strconv.ParseInt(argTerm, 10, 64)
			if err != nil {
				return types.ErrParameter.Wrapf("period error: %v", err)
			}

			rate, err := sdk.NewDecFromStr(argRate)
			if err != nil {
				return types.ErrParameter.Wrap("rate error")
			}

			msg := types.NewMsgSetFixedDepositCfgRate(
				clientCtx.GetFromAddress().String(),
				argRegionId,
				term,
				rate,
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
