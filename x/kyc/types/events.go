package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeApprove          = "approve"
	EventTypeUpdate           = "update"
	EventTypeRemove           = "remove"
	EventTypeCreateSBT        = "create_sbt"
	EventTypeUpdateSBT        = "update_sbt"
	EventTypeDeleteSBT        = "delete_sbt"
	EventTypeCreateSubAccount = "create_sub_account"
)

const (
	AttributeKeyAddress         = "address"
	AttributeKeyRegionId        = "region_id"
	AttributeKeyRegionIdChanged = "region_id_changed"
	AttributeKeyLevel           = "level"
	AttributeKeyLevelChanged    = "level_changed"
	AttributeKeyInviter         = "inviter"
	AttributeKeyAccount         = "account"
	AttributeKeySubAccount      = "sub_account"
	AttributeKeyDid             = "did"
	AttributeKeyCreator         = "creator"
)

func NewSbtEvent(eventType, did, uri, hash, regionId, kycLevel, meIdAddress string) sdktypes.Event {
	attributes := []sdktypes.Attribute{
		{Key: "did", Value: did},
		{Key: "uri", Value: uri},
		{Key: "hash", Value: hash},
		{Key: "regionId", Value: regionId},
		{Key: "kycLevel", Value: kycLevel},
		{Key: "meIdAddress", Value: meIdAddress},
		{Key: "class_id", Value: ModuleName},
	}
	return sdktypes.NewEvent(eventType, attributes...)
}
