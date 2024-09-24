package types

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewKycEvent(address string, did string, action string) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "address", Value: address},
		{Key: "did", Value: did},
		{Key: "action", Value: action},
	}
	return sdkTypes.NewEvent("kyc_event", attributes...)
}
