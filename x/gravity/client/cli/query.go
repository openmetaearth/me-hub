package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/gravity/types"
)

func GetQueryCmd(moduleName string, subNames ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", moduleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	for _, chainName := range subNames {
		cmd.AddCommand(GetQueryCmd(chainName))
	}
	if len(subNames) == 0 {
		cmd.AddCommand(getQuerySubCmds(moduleName)...)
	}
	return cmd
}

func getQuerySubCmds(chainName string) []*cobra.Command {
	cmds := []*cobra.Command{
		// query module params
		CmdGetParams(chainName),
	}

	for _, command := range cmds {
		flags.AddQueryFlagsToCmd(command)
	}
	return cmds
}

func CmdGetParams(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current parameters information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(&res.Params)
		},
	}
	return cmd
}
