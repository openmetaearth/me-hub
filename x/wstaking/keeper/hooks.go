package keeper

import "github.com/openmetaearth/me-hub/x/wstaking/types"

func (k *Keeper) SetWstakingHooks(hk types.WstakingHooks) {
	if k.wstakingHooks != nil {
		panic("cannot set wstaking hooks twice")
	}
	k.wstakingHooks = hk
}

// Hooks gets the hooks for wstaking *Keeper {
func (k *Keeper) WstakingHooks() types.WstakingHooks {
	if k.wstakingHooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiWstakingHooks{}
	}

	return k.wstakingHooks
}
