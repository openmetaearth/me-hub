package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MigrateRequiredOperatorAccounts re-encodes existing proposal approvers without
// consulting staking or changing snapshots, votes, queues, or proposal status.
func (k Keeper) MigrateRequiredOperatorAccounts(ctx sdk.Context) error {
	proposals := k.GetAllProposals(ctx)
	// Validate every address before writing any proposal.
	for i := range proposals {
		if err := proposals[i].NormalizeRequiredOperator(); err != nil {
			return fmt.Errorf("proposal %d required operator: %w", proposals[i].Id, err)
		}
	}
	for _, proposal := range proposals {
		if proposal.RequiredOperator != "" {
			k.SetProposal(ctx, proposal)
		}
	}
	return nil
}
