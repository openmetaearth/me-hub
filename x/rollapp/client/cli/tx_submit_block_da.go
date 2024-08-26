package cli

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dymensionxyz/dymension/v3/x/rollapp/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdSubmitBlockDa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-block-da [rollapp-id] [start-height] [num-blocks] [da-path] [version] [blks] [commitment] [daroot]",
		Short: "submit block and da commitment-proof",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRollappId := args[0]
			argStartHeight, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argNumBlocks, err := cast.ToUint32E(args[2])
			if err != nil {
				return err
			}
			argDAPath := args[3]
			argVersion, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}
			argBlks := new(types.MsgLightBlkInfos)
			err = json.Unmarshal([]byte(args[5]), argBlks)
			if err != nil {
				return err
			}

			daCommitment, err := hex.DecodeString(args[6])
			if err != nil {
				return err
			}

			daRoot, err := hex.DecodeString(args[7])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgBlkDAInfo{
				Creator:         clientCtx.GetFromAddress().String(),
				RollappId:       argRollappId,
				StartHeight:     argStartHeight,
				NumBlocks:       argNumBlocks,
				DAPath:          argDAPath,
				Version:         argVersion,
				Blocks:          *argBlks,
				CommitmentProof: daCommitment,
				DaRoot:          daRoot,
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

func CmdGetLastSubmitBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-last-submit-block [rollapp-id] [start-height] [num-blocks] [da-path] [version] [blks] [commitment] [daroot]",
		Short: "get last submit block info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewMsgClient(clientCtx)
			argRollappId := args[0]
			if argRollappId == "" {
				return fmt.Errorf("rollappID can not be empty")
			}
			req := &types.MsgLastSubmitBlkRequest{
				RollappId: argRollappId,
			}

			res, err := queryClient.GetLastSubmitBlockInfo(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
