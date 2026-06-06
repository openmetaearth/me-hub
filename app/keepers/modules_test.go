package keepers

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func TestModuleAccountAddrsAllowsBridgeFeePoolReceives(t *testing.T) {
	blockedAddrs := (&AppKeepers{}).ModuleAccountAddrs()
	bridgeFeePoolAddr := authtypes.NewModuleAddress(types.BridgeFeePool).String()

	if blockedAddrs[bridgeFeePoolAddr] {
		t.Fatalf("bridge fee pool must be able to receive configured bridge fees")
	}

	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName).String()
	if !blockedAddrs[feeCollectorAddr] {
		t.Fatalf("unrelated module accounts should remain blocked")
	}

	if _, ok := MaccPerms[types.BridgeFeePool]; !ok {
		t.Fatalf("bridge fee pool should remain a module account")
	}
}
