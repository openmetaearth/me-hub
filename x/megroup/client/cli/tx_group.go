package cli

import (
	"strconv"
	"time"

	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/openmetaearth/me-hub/x/megroup/types"
	"github.com/spf13/cobra"
)

func CmdCreateGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-group [admin] [regionID]",
		Short: "Create a new group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			argGroupInfo := &types.GroupInfo{
				Id:          0,
				Admin:       args[0],
				Metadata:    "",
				Version:     1,
				TotalWeight: "",
				CreatedAt:   time.Time{},
				RegionID:    args[1],
			}

			msg := types.NewMsgCreateGroup(clientCtx.GetFromAddress().String(), argGroupInfo)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			//clientCtx.PrintProto(argGroupInfo)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-group [id] [group-info]",
		Short: "Update a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argGroupInfo := new(types.GroupInfo)
			err = json.Unmarshal([]byte(args[1]), argGroupInfo)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateGroup(clientCtx.GetFromAddress().String(), id, argGroupInfo)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-group [id]",
		Short: "Delete a group by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteGroup(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
