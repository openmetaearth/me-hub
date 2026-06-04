package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type DidHooks interface {
	AfterDidStatusUpdated(ctx sdk.Context, info DidInfo) error
}

type MultiDidHooks []DidHooks

func NewMultiDidHooks(hooks ...DidHooks) MultiDidHooks {
	return hooks
}

func (h MultiDidHooks) AfterDidStatusUpdated(ctx sdk.Context, info DidInfo) error {
	for i := range h {
		if err := h[i].AfterDidStatusUpdated(ctx, info); err != nil {
			return err
		}
	}
	return nil
}
