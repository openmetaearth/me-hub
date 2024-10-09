package cli

import (
    "strconv"
	
	"github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/st-chain/me-hub/x/did/types"
)

var _ = strconv.Itoa(0)

func CmdDidInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did-infos",
		Short: "Query did_infos",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDidInfosRequest{
				
            }

            pageReq, err := client.ReadPageRequest(cmd.Flags())
            if err != nil {
                return err
            }
            params.Pagination = pageReq

			res, err := queryClient.DidInfos(cmd.Context(), params)
            if err != nil {
                return err
            }

            return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}
