package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type KycHook interface {
	SetKycIssers(ctx sdk.Context, oldDaoAddress, newDaoAddress []string) error
}
