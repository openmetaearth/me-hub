package cmd

import (
	"encoding/hex"
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/cosmos/gogoproto/proto"

	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

// GetEncodeCommand returns the encode command to take a JSONified transaction and turn it into
// Amino-serialized bytes
func GetEncodeToRawTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "encode-raw-tx [file]",

		Short: "Encode transactions generated offline with raw tx format",
		Long: `Encode transactions created with the --generate-only flag or signed with the sign command.
Read a transaction from <file>, serialize it to the Protobuf wire protocol, and output it as base64.
If you supply a dash (-) argument in place of an input filename, the command reads from standard input.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			tx, err := authclient.ReadTxFromFile(clientCtx, args[0])
			if err != nil {
				return err
			}

			// re-encode it
			txBytes, err := clientCtx.TxConfig.TxEncoder()(tx)
			if err != nil {
				return err
			}
			raw := &txtypes.TxRaw{}
			err = proto.Unmarshal(txBytes, raw)
			if err != nil {
				return err
			}
			encodeJson, err := json.Marshal(raw)
			if err != nil {
				return err
			}
			//if flag hex is true
			if useHex, _ := cmd.Flags().GetBool("hex"); useHex {
				return clientCtx.PrintString(hex.EncodeToString(encodeJson))
			}
			return clientCtx.PrintString(string(encodeJson) + "\n")
		},
	}
	cmd.Flags().Bool("hex", true, "output with hex format")
	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(flags.FlagOutput) // encoding makes sense to output only json

	return cmd
}
