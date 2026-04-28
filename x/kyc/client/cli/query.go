package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/spf13/cobra"
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

	cmd.AddCommand(CmdQueryProtocol())
	cmd.AddCommand(CmdQueryDid())
	cmd.AddCommand(CmdQueryKYC())
	cmd.AddCommand(CmdQueryKYCs())
	cmd.AddCommand(CmdQuerySBT())
	return cmd
}

func CmdQueryProtocol() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "protocol",
		Short: "shows the KYC protocol",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Protocol(cmd.Context(), &types.QueryProtocol{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryDid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did [address]",
		Short: "query DID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr := args[0]
			res, err := queryClient.DID(cmd.Context(), &types.QueryDID{Address: addr})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryKYC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kyc [DID]",
		Short: "Query the KYC information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			res, err := queryClient.KYC(cmd.Context(), &types.QueryKYC{Did: did})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryKYCs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "KYCs",
		Short: "Query the KYCs information",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			region := ""
			if f := cmd.Flag("region"); f != nil {
				region = f.Value.String()
			}

			pageReq, _ := client.ReadPageRequest(cmd.Flags())
			res, err := queryClient.KYCs(cmd.Context(), &types.QueryKYCs{
				RegionId:   region,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	f := cmd.Flags()
	f.String("region", "", "filter by region_id ,example: me_earth")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	return cmd
}

func CmdQuerySBT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sbt [did]",
		Short: "Query the SBT information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			res, err := queryClient.SBT(cmd.Context(), &types.QuerySBT{Did: did})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
