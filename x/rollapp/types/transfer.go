package types

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

type TransferData struct {
	transfertypes.FungibleTokenPacketData
	// Rollapp will be the nil if the packet is not to/from a registered rollapp
	Rollapp *Rollapp
}

// IsRollapp returns whether the transfer came from a rollapp or was sent to a rollapp
func (d TransferData) IsRollapp() bool {
	return d.Rollapp != nil
}

func (d TransferData) RollappId() string {
	if d.Rollapp != nil {
		return d.Rollapp.RollappId
	}
	return ""
}

// MustAmountInt returns the int amount. Should call validateBasic first!
func (d TransferData) MustAmountInt() math.Int {
	x, ok := sdk.NewIntFromString(d.Amount)
	if !ok {
		panic(fmt.Sprintf("parse transfer amount to Int: %s", d.Amount))
	}
	return x
}
