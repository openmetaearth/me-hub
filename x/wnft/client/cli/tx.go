package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/nft"
	wnfttypes "github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:                        nft.ModuleName,
		Short:                      "nft transactions subcommands",
		Long:                       "Provides the most common nft logic for upper-level applications, compatible with Ethereum's erc721 contract",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nftTxCmd.AddCommand(
		NewCmdNewClass(),
		NewCmdMintNFT(),
		NewCmdSend(),
	)

	return nftTxCmd
}

// NewCmdNewClass creates a CLI command for MsgNewClass.
func NewCmdNewClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-class [class-id] [name] [symbol] [description] [uri] [uri_hash] [total_supply]",
		Args:  cobra.ExactArgs(7),
		Short: "create a class",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s new-class [class-id] [name] [symbol] [description] [uri] [uri_hash] [total_supply] --from [sender] --chain-id <chain-id>`, version.AppName, nft.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			classId := args[0]
			name := args[1]
			symbol := args[2]
			description := args[3]
			uri := args[4]
			uriHash := args[5]

			argTotalSupply := args[6]
			totalSupply, err := strconv.ParseUint(argTotalSupply, 10, 64)
			if err != nil {
				return err
			}

			msg := wnfttypes.NewMsgNewClass(classId, clientCtx.GetFromAddress().String(), name, symbol, description, uri, uriHash, totalSupply)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewCmdMintNFT creates a CLI command for NewCmdMintNFT.
func NewCmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [class-id] [token-id] [uri] [uri-hash] [receiver] --from [sender]",
		Args:  cobra.ExactArgs(5),
		Short: "create a nft",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s mint [class-id] [token-id] [uri] [uri-hash] [receiver] --from [sender] --chain-id <chain-id>`, version.AppName, nft.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			classId := args[0]

			tokenId := args[1]

			url := args[2]
			urlHash := args[3]
			receiver := args[4]

			if err != nil {
				return types.ErrParameter.Wrap("term error")
			}

			msg := wnfttypes.NewMsgMintNFT(classId, tokenId, url, urlHash, clientCtx.GetFromAddress().String(), receiver)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewCmdSend creates a CLI command for MsgSend.
func NewCmdSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [class-id] [nft-id] [receiver] --from [sender]",
		Args:  cobra.ExactArgs(3),
		Short: "transfer ownership of nft",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s send <class-id> <nft-id> <receiver> --from <sender> --chain-id <chain-id>`, version.AppName, nft.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := wnfttypes.MsgSend{
				ClassId:  args[0],
				Id:       args[1],
				Sender:   clientCtx.GetFromAddress().String(),
				Receiver: args[2],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
