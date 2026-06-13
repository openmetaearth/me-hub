(Go)

## Solution
The error is fixed by modifying the `Delegate` function in `x/wstaking/keeper/delegation.go`. The original implementation returns `newShares` as the correct amount for experience region delegations, which was causing the incorrect event emission. By ensuring that the bond amount is properly accounted for in the `delegation.Amount`, the newShares variable now correctly reflects the intended delegation value.

package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/openmetaearth/me-hub/x/wstaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Delegate returns the newShares for a delegation event.
func (k Keeper) Delegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int,
	tokenSrc stakingtypes.BondStatus,
	valAddr sdk.ValAddress, delegator validator,
) (newShares sdk.Dec, err error) {
	regionId := k.GetRegionByAccount(ctx, delAddr)
	if !regionId.IsValid() {
		return nil, types.ErrRegionNotExists
	}
	region, isFound := k.GetRegion(ctx, regionId)
	if !isFound {
		return nil, 

ACTUAL REPO CODE (use these exact function names, imports, and patterns):
// FILE: app/ante/ante.go
package ante

import (
	"fmt"
	"runtime/debug"

	sdk