package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type WstakingHooks interface {
	BeforeValidatorStakingModified(ctx sdk.Context, val sdk.ValAddress) error
}
type MultiWstakingHooks []WstakingHooks

func NewMultiWstakingHooks(hooks ...WstakingHooks) MultiWstakingHooks {
	return hooks
}

func (h MultiWstakingHooks) BeforeValidatorStakingModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].BeforeValidatorStakingModified(ctx, valAddr); err != nil {
			return err
		}
	}
	return nil
}
