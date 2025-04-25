package cli

import (
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/rollapp/types"
	"strconv"
)

func CmdChallengeDaFraud() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "challengeDaFraud [RollappId] [StartHeight] [NumBlocks] [Namespace] [Commitment] [daRoot] [DaBlockHeight]",
		Short: "challengeDaFraud",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rollappID := args[0]
			startHeight, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			numberBlk, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}
			namespace, err := hex.DecodeString(args[3])
			if err != nil {
				return err
			}

			commitment, err := hex.DecodeString(args[4])
			if err != nil {
				return err
			}

			daRoot, err := hex.DecodeString(args[5])
			if err != nil {
				return err
			}

			daBlkHeight, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgSubmitDaFraudRequest{
				Creator:       clientCtx.GetFromAddress().String(),
				RollappId:     rollappID,
				StartHeight:   startHeight,
				NumBlocks:     uint32(numberBlk),
				Namespace:     namespace,
				Commitment:    commitment,
				DaRoot:        daRoot,
				DaBlockHeight: daBlkHeight,
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSubmitDaFraudVerifyData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submitDaFraudVerifyData [RollappId] [DaPath] [StartHeight] [NumBlocks] [Result]",
		Short: "SubmitDaFraudVerifyData",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rollappID := args[0]
			daPath := args[1]
			startHeight, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			numberBlk, err := strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				return err
			}
			result, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgDaFraudVerifyResult{
				Creator:     clientCtx.GetFromAddress().String(),
				DaPath:      daPath,
				RollappId:   rollappID,
				StartHeight: startHeight,
				NumBlocks:   uint32(numberBlk),
				Result:      int32(result),
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
