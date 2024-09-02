package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dymensionxyz/dymension/v3/x/rollapp/types"
	"github.com/spf13/cobra"
	"strconv"
)

func CmdGetLastSubmitBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-last-submit-block [rollapp-id]",
		Short: "get last submit block info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
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

func CmdGetSubmitterBlockStatics() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-submitter-block-statics [rollapp-id] [startHeight] [endHeight]",
		Short: "get submitter block statics data",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			argRollappId := args[0]
			if argRollappId == "" {
				return fmt.Errorf("rollappID can not be empty")
			}
			startHeight, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			endHeight, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			req := &types.MsgSubmitBlockStaticsRequest{
				RollappId:   argRollappId,
				StartHeight: startHeight,
				EndHeight:   endHeight,
			}

			res, err := queryClient.GetSubmitterBlockStatics(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
