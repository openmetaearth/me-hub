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
	//cmd.AddCommand(CmdRegisterRollAppIDTest())
	return cmd
}

func CmdStakeForSequencer() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stakeForSequencer  [rollappId] [amount]",
		Short:   "stakeForSequencer",
		Example: "dymd tx hubRollUp stakeForSequencer  <rollappId> <amount>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			//creator := args[0]
			rollappID := args[0]
			amount := args[1]
			val, err := strconv.ParseUint(amount, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgSeqStaking(clientCtx.GetFromAddress().String(), rollappID, 0, val)
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
		Use:     "unStake  [rollappId] [amount]",
		Short:   "unstake mec",
		Example: "dymd tx hubRollUp unStake <rollappId> <amount>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			rollappID := args[0]
			amount := args[1]
			val, err := strconv.ParseUint(amount, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgSeqUnStaking(clientCtx.GetFromAddress().String(), rollappID, 0, val)
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
		Use:     "setParams [rollappId] [electionPeriod] [seqNumber] [backupNumber] [minStake] [firstElectTime] [allowApplyTime] [electInterim] [daFraudChallengeStake]",
		Short:   "set rollup Params",
		Example: "med tx hubRollUp setParams ",
		Args:    cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			rollappID := args[0]
			electionPeriod, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			seqNumber, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}
			backupNumber, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}
			minStakeAmount, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}
			firstElectTime, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}
			allowApplyTime, err := strconv.Atoi(args[6])
			if err != nil {
				return err
			}
			electInterim, err := strconv.Atoi(args[7])
			if err != nil {
				return err
			}
			daChallengeStake, err := strconv.Atoi(args[8])
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
				DaFraudChallengeStake:  uint32(daChallengeStake),
			}

			req := &types.MsgSetRollupParamsRequest{
				Creator:   clientCtx.GetFromAddress().String(),
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
