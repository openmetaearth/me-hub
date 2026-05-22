package types

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeApprove   = "approve"
	EventTypeUpdate    = "update"
	EventTypeRemove    = "remove"
	EventTypeCreateSBT = "create_sbt"
	EventTypeUpdateSBT = "update_sbt"
	EventTypeDeleteSBT = "delete_sbt"
)

const (
	AttributeKeyAddress         = "address"
	AttributeKeyRegionId        = "region_id"
	AttributeKeyRegionIdChanged = "region_id_changed"
	AttributeKeyLevel           = "level"
	AttributeKeyLevelChanged    = "level_changed"
	AttributeKeyInviter         = "inviter"
)

func NewSbtEvent(eventType, did, uri, hash, regionId, kycLevel, meIdAddress string) sdkTypes.Event {
	attributes := []sdkTypes.Attribute{
		{Key: "did", Value: did},
		{Key: "uri", Value: uri},
		{Key: "hash", Value: hash},
		{Key: "regionId", Value: regionId},
		{Key: "kycLevel", Value: kycLevel},
		{Key: "meIdAddress", Value: meIdAddress},
		{Key: "class_id", Value: ModuleName},
	}
	return sdkTypes.NewEvent(eventType, attributes...)
}
