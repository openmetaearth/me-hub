package types

import (
	"fmt"
	"strconv"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeApprove   = "approve"
	EventTypeUpdate    = "update"
	EventTypeRemove    = "remove"
	EventTypeCreateSBT = "create_sbt"
	EventTypeUpdateSBT = "create_sbt"
	EventTypeDeleteSBT = "delete_sbt"
)

func NewKycEvent(address string, did string, level int, action string, seq uint64) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "sequence", Value: fmt.Sprintf("%d", seq)},
		{Key: "address", Value: address},
		{Key: "did", Value: did},
		{Key: "level", Value: strconv.Itoa(level)},
		{Key: "action", Value: action},
	}
	return sdkTypes.NewEvent("kyc_event", attributes...)
}

func NewSbtEvent(eventType, did, uri, hash string) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "did", Value: did},
		{Key: "uri", Value: uri},
		{Key: "hash", Value: hash},
	}
	return sdkTypes.NewEvent(eventType, attributes...)
}
