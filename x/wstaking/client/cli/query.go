package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	stakingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the staking module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingQueryCmd.AddCommand(
		GetCmdQueryValidators(),
		GetCmdQueryValidator(),
		GetCmdQueryValidatorDelegations(),
		GetCmdQueryPool(),
		GetCmdQueryStakingParams(),
		GetCmdQueryRegion(),
		GetCmdQueryAllRegion(),
		GetCmdQueryDelegatorRewards(),
		GetCmdQueryDelegation(),
		CmdQueryAllDelegations(),
		GetCmdQueryStakes(),
		CmdListFixedDeposit(),
		CmdShowFixedDeposit(),
		CmdFixedDepositByRegion(),
		CmdFixedDepositByAcct(),
		CmdListFixedDepositCfg(),
		CmdQueryFixedDepositCfg(),
		CmdShowFixedDepositTotalAmount(),
		CmdShowFixedDepositAmountByAcct(),
		CmdShowAllRecord(),
		CmdShowRecordByAddress(),
		CmdShowReviewRecord(),
	)
	return stakingQueryCmd
}
