package types

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewKycEvent(address string, did string, action string, seq uint64) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "sequence", Value: fmt.Sprintf("%d", seq)},
		{Key: "address", Value: address},
		{Key: "did", Value: did},
		{Key: "action", Value: action},
	}
	return sdkTypes.NewEvent("kyc_event", attributes...)
}
