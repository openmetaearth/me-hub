package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const IROTokenPrefix = "IRO/"

func IRODenom(rollappID string) string {
	return fmt.Sprintf("%s%s", IROTokenPrefix, rollappID)
}

func RollappIDFromIRODenom(denom string) (string, bool) {
	return strings.CutPrefix(denom, IROTokenPrefix)
}

// Plan is a minimal stub so rollapp can compile without the full IRO module.
// Phase 1 does not register x/iro.
type Plan struct {
	Id        uint64
	RollappId string
	// Allocation kept for API compatibility with SetIROPlanToRollapp signature.
	TotalAllocation sdk.Coin
}
