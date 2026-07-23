package types

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
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
	AttributeKeyCreator         = "creator"
	AttributeKeyAccount         = "account"
	AttributeKeySubAccount      = "sub_account"
	AttributeKeyDid             = "did"
)

func NewKycEvent(
	address string,
	did string,
	level didtypes.KycLevel,
	action string,
	seq uint64,
) sdktypes.Event {
	attributes := []sdktypes.Attribute{
		{Key: "sequence", Value: fmt.Sprintf("%d", seq)},
		{Key: "address", Value: address},
		{Key: "did", Value: did},
		{Key: "level", Value: level.String()},
		{Key: "action", Value: action},
	}
	return sdktypes.NewEvent("kyc_event", attributes...)
}

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
