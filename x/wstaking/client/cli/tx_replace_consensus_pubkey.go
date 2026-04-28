package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/spf13/cobra"
)

func CmdReplaceConsensusPubKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace-consensus-pubkey [operator] [new-pubkey] [block-number]",
		Short: "Broadcast message replace-consensus-pubkey",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			operator := args[0]
			blocl_number, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var pk cryptotypes.PubKey
			if err = clientCtx.Codec.UnmarshalInterfaceJSON([]byte(args[1]), &pk); err != nil {
				return err
			}
			if pk.Bytes() == nil {
				return fmt.Errorf("pubkey by UnmarshalInterfaceJSON cannot be nil")
			}
			codecPubKey, err := codectypes.NewAnyWithValue(pk)
			if err != nil {
				return err
			}

			msg := &types.MsgReplaceConsensusPubKeyRequest{
				Creator: clientCtx.GetFromAddress().String(),
				ReplacePubKey: &types.MsgReplaceConsensusPubKey{
					OperatorAddress: operator,
					PubKey:          codecPubKey,
					BlockNumber:     blocl_number,
				}}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
