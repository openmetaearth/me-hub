package types

import (
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeCreateDid           = "create_did"
	EventTypeUpdateDidStatus     = "update_did_status"
	EventTypeCreateService       = "create_service"
	EventTypeUpdateServiceStatus = "update_service_status"
	EventTypeCreateVC            = "create_vc"
	EventTypeUpdateVC            = "update_vc"
	EventTypeRemoveVC            = "remove_vc"
)

func NewDidEvent(eventType, did, address, status string) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "did", Value: did},
		{Key: "address", Value: address},
		{Key: "status", Value: status},
	}
	return sdkTypes.NewEvent(eventType, attributes...)
}

func NewServiceEvent(eventType, sid, name, status string, issuers []string) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "sid", Value: sid},
		{Key: "name", Value: name},
		{Key: "status", Value: status},
		{Key: "issuers", Value: fmt.Sprintf("%v", issuers)},
	}
	return sdkTypes.NewEvent(eventType, attributes...)
}

func NewVcEvent(eventType, sid, did, hash, uri string) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "sid", Value: sid},
		{Key: "did", Value: did},
		{Key: "hash", Value: hash},
		{Key: "uri", Value: uri},
	}
	return sdkTypes.NewEvent(eventType, attributes...)
}
