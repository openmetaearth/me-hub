package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/did/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group did queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryDid())
	cmd.AddCommand(CmdQueryDidInfo())
	cmd.AddCommand(CmdQueryDidInfos())
	cmd.AddCommand(CmdQueryDidDocument())
	cmd.AddCommand(CmdQueryService())
	cmd.AddCommand(CmdQueryServices())
	cmd.AddCommand(CmdQueryCredential())

	return cmd
}

func CmdQueryDid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did [address]",
		Short: "query did",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr := args[0]
			res, err := queryClient.Did(cmd.Context(), &types.QueryDid{Address: addr})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryDidInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did-info [did]",
		Short: "Query did base information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			res, err := queryClient.DidInfo(cmd.Context(), &types.QueryDidInfo{Did: did})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryDidInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did-infos",
		Short: "Query did base information list",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, _ := client.ReadPageRequest(cmd.Flags())
			res, err := queryClient.DidInfos(cmd.Context(), &types.QueryDidInfos{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdQueryDidDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "document [did]",
		Short: "query did document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			res, err := queryClient.DidDocument(cmd.Context(), &types.QueryDidDocument{Did: did})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryService() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service [sid]",
		Short: "query credential service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			sid := args[0]
			res, err := queryClient.Service(cmd.Context(), &types.QueryService{Sid: sid})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryServices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "Query did service list",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, _ := client.ReadPageRequest(cmd.Flags())
			res, err := queryClient.Services(cmd.Context(), &types.QueryServices{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdQueryCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credential [did] [sid]",
		Short: "query verifiable credential",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			sid := args[1]
			res, err := queryClient.Credential(cmd.Context(), &types.QueryCredential{Did: did, Sid: sid})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
