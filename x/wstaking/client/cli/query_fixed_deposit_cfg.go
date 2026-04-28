package cli

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"strings"

	"github.com/spf13/cobra"
	"strconv"
)

func CmdListFixedDepositCfg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fixed-deposit-cfg [region-id,region-id,...]",
		Short: "show some regions fixed deposit config by region ids",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			regionIdStr := args[0]
			newRegionIdStr := strings.Trim(regionIdStr, " ")

			var regionIds []string
			if newRegionIdStr == "" {
				regionIds = []string{}
			} else {
				regionIds = strings.Split(newRegionIdStr, ",")
				for _, regionId := range regionIds {
					_, err := utils.CheckRegionName(strings.ToUpper(regionId))
					if err != nil {
						return sdkerrors.Wrap(types.ErrRegionName, err.Error())
					}
				}
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFixedDepositCfgRequest{
				RegionIds: regionIds,
			}

			res, err := queryClient.FixedDepositCfg(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryFixedDepositCfg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fixed-deposit-cfg-by-term",
		Short: "show fixed deposit config by term",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argRegionId := args[0]
			argTerm := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			term, err := strconv.ParseInt(argTerm, 10, 64)
			if err != nil {
				return types.ErrParameter.Wrap("term error")
			}

			params := &types.QueryFixedDepositCfgByTermRequest{
				RegionId: argRegionId,
				Term:     term,
			}

			res, err := queryClient.FixedDepositCfgByTerm(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
