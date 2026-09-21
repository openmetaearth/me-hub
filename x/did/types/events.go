package types

import (
	"encoding/hex"
	"fmt"
	"strconv"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeCreateDid           = "create_did"
	EventTypeUpdateDidStatus     = "update_did_status"
	EventTypeCreateService       = "create_service"
	EventTypeUpdateServiceStatus = "update_service_status"
	EventTypeCreateVC            = "create_vc"
	EventTypeUpdateVC            = "update_vc"
	EventTypeRemoveVC            = "remove_vc"
	EventTypeBondZkCoreID        = "bond_zk_core_id"
	EventTypeSetZkContract       = "set_zk_contract_address"
)

func NewDidEvent(eventType, did, address, status string) sdktypes.Event {
	attributes := []sdktypes.Attribute{
		{Key: "did", Value: did},
		{Key: "address", Value: address},
		{Key: "status", Value: status},
	}
	return sdktypes.NewEvent(eventType, attributes...)
}

func NewServiceEvent(eventType, sid, name, status string, issuers []string) sdktypes.Event {
	attributes := []sdktypes.Attribute{
		{Key: "sid", Value: sid},
		{Key: "name", Value: name},
		{Key: "status", Value: status},
		{Key: "issuers", Value: fmt.Sprintf("%v", issuers)},
	}
	return sdktypes.NewEvent(eventType, attributes...)
}

func NewVcEvent(eventType, sid, did, hash, uri string) sdktypes.Event {
	attributes := []sdktypes.Attribute{
		{Key: "sid", Value: sid},
		{Key: "did", Value: did},
		{Key: "hash", Value: hash},
		{Key: "uri", Value: uri},
	}
	return sdktypes.NewEvent(eventType, attributes...)
}

func NewBondZkCoreIDEvent(did, creator string, coreID []byte,blockHeight int64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeBondZkCoreID,
		sdktypes.NewAttribute("did", did),
		sdktypes.NewAttribute("creator", creator),
		sdktypes.NewAttribute("core_id", hex.EncodeToString(coreID)),
		sdktypes.NewAttribute("block_height", strconv.FormatInt(blockHeight,10)),
	)
}

func NewSetZkContractAddressEvent(info ZkContractAddressInfo,blockHeight int64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeSetZkContract,
		sdktypes.NewAttribute("contract_type", info.Type.String()),
		sdktypes.NewAttribute("contract_address", info.ContractAddress),
		sdktypes.NewAttribute("block_height", strconv.FormatInt(blockHeight,10)),
	)
}
