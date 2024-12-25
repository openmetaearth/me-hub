package types

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	didTypes "github.com/st-chain/me-hub/x/did/types"
)

const (
	EventTypeApprove   = "approve"
	EventTypeUpdate    = "update"
	EventTypeRemove    = "remove"
	EventTypeCreateSBT = "create_sbt"
	EventTypeUpdateSBT = "create_sbt"
	EventTypeDeleteSBT = "delete_sbt"
)

const (
	AttributeKeyRegionId        = "region_id"
	AttributeKeyRegionIdChanged = "region_id_changed"
	AttributeKeyLevel           = "level"
	AttributeKeyLevelChanged    = "level_changed"
)

func NewKycEvent(address string, did string, level didTypes.KycLevel, action string, seq uint64) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "sequence", Value: fmt.Sprintf("%d", seq)},
		{Key: "address", Value: address},
		{Key: "did", Value: did},
		{Key: "level", Value: level.String()},
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
