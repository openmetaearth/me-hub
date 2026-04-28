package cli

import (
	"fmt"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/nft"

	nftcli "github.com/cosmos/cosmos-sdk/x/nft/client/cli"
)

// Flag names and values
const (
	FlagOwner   = "owner"
	FlagClassID = "class-id"
	FlagTokenId = "token-id"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	nftQueryCmd := &cobra.Command{
		Use:                        nft.ModuleName,
		Short:                      "Querying commands for the nft module",
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nftQueryCmd.AddCommand(
		nftcli.GetCmdQueryClass(),
		nftcli.GetCmdQueryClasses(),
		nftcli.GetCmdQueryNFT(),
		nftcli.GetCmdQueryNFTs(),
		nftcli.GetCmdQueryOwner(),
		nftcli.GetCmdQueryBalance(),
		nftcli.GetCmdQuerySupply(),
		GetCmdQueryClassAddress(),
		GetCmdQueryNftFilter(),
	)
	return nftQueryCmd
}

// GetCmdQueryClassAddress implements the query class by address command.
func GetCmdQueryClassAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "class-address [class-id] [owner]",
		Args:    cobra.ExactArgs(2),
		Short:   "query information and status related to the specified NFT class (album).",
		Example: fmt.Sprintf(`$ %s query %s class-address [class-id] [owner]`, version.AppName, nft.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			owner := args[1]

			if len(owner) > 0 {
				if _, err := sdk.AccAddressFromBech32(owner); err != nil {
					return err
				}
			}

			res, err := queryClient.ClassAddress(cmd.Context(), &types.QueryClassAddressRequest{
				ClassId: args[0],
				Address: owner,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryNftFilter implements the query by filter command.
func GetCmdQueryNftFilter() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nft-filter --class-id [class-id] --token-id [token-id] --owner [owner]",
		Short:   "nft query filter",
		Example: fmt.Sprintf(`$ %s query %s nft-filter --class-id [class-id] --token-id [token-id] --owner [owner]`, version.AppName, nft.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			owner, err := cmd.Flags().GetString(FlagOwner)

			if len(owner) > 0 {
				if _, err := sdk.AccAddressFromBech32(owner); err != nil {
					return err
				}
			}

			classId, _ := cmd.Flags().GetString(FlagClassID)

			tokenId, _ := cmd.Flags().GetString(FlagTokenId)

			res, err := queryClient.NftFilter(cmd.Context(), &types.QueryNftFilterRequest{
				Owner:   owner,
				ClassId: classId,
				TokenId: tokenId,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagOwner, "", "The owner of the nft")
	cmd.Flags().String(FlagClassID, "", "The class-id of the nft")
	cmd.Flags().String(FlagTokenId, "", "nft token id")

	return cmd
}
