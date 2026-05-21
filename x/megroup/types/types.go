package types

import (
	"encoding/binary"
)

const (
	EvtGroupCreated     string = "EventGroupCreated"
	EvtJoinGroup        string = "EventJoinGroup"
	EvtGrpMigrateByKyc  string = "EventGroupMigrateByKycChange"
	EvtJoinGroupReward  string = "EventJoinGroupReward"
	EvtUpdateGroupAdmin string = "EventUpdateGroupAdmin"
)

type GroupMetaData struct {
	SubmitMetaData string `json:"submit_meta_data"`
	AdminMeid      string `json:"admin_meid_data"`
}

// GetGroupMemberIDBytes returns the byte representation of the ID
func GetBytesFromUint64(val uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, val)
	return bz
}

// GetGroupMemberIDFromBytes returns ID in uint64 format from a byte array
func GetUint64FromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
