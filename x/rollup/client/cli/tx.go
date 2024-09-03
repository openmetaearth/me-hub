package cli

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/st-chain/me-hub/x/rollup/types"
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
	cmd.AddCommand(CmdSetRollupParams())

	return cmd
}

func CmdStakeForSequencer() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stakeForSequencer [creator] [rollappId] [amount]",
		Short:   "stakeForSequencer",
		Example: "dymd tx hubRollUp stakeForSequencer <creator> <rollappId> <amount>",
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
		Use:     "unStake [creator] [rollappId] [amount]",
		Short:   "unstake mec",
		Example: "dymd tx hubRollUp unStake <creator> <rollappId> <amount>",
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

func CmdSetRollupParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "setParams [creator] [rollappId] [electionPeriod] [seqNumber] [backupNumber] [minStake] [firstElectTime] [allowApplyTime] [electInterim]",
		Short:   "set rollup Params",
		Example: "med tx hubRollUp setParams ",
		Args:    cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			creator := args[0]
			rollappID := args[1]
			electionPeriod, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			seqNumber, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}
			backupNumber, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}
			minStakeAmount, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}
			firstElectTime, err := strconv.Atoi(args[6])
			if err != nil {
				return err
			}
			allowApplyTime, err := strconv.Atoi(args[7])
			if err != nil {
				return err
			}
			electInterim, err := strconv.Atoi(args[8])
			if err != nil {
				return err
			}
			params := &types.Params{
				ElectionPeriod:         uint32(electionPeriod),
				SequencerNumber:        uint32(seqNumber),
				BackupSequencerNumber:  uint32(backupNumber),
				MinStakeAmount:         minStakeAmount,
				FirstElectionInterval:  uint32(firstElectTime),
				AllowApplyElectionTime: uint32(allowApplyTime),
				ElectionInterimTime:    uint32(electInterim),
			}

			req := &types.MsgSetRollupParamsRequest{
				Creator:   creator,
				RollappID: rollappID,
				NewParams: params,
			}

			if err = req.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), req)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
