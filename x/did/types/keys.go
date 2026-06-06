package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "did"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_did"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	ParamsKey     = []byte{0x01}
	DIDPrefix     = []byte{0x10}
	DidInfoPrefix = []byte{0x11}
	IssuerPrefix  = []byte{0x20}
	ServicePrefix = []byte{0x30}
	//ServiceIssuerPrefix   = []byte{0x31}
	CredentialPrefix   = []byte{0x40}
	FilterLoggerPrefix = []byte{0x50}
	FilterPrefix       = []byte{0x51}
)

func GetDIDKey(addr sdk.AccAddress) []byte {
	return append(DIDPrefix, addr.Bytes()...)
}

func GetDidInfoKey(did string) []byte {
	return append(DidInfoPrefix, []byte(did)...)
}

func GetIssuerKey(did string) []byte {
	return append(IssuerPrefix, []byte(did)...)
}

func GetServiceKey(sid string) []byte {
	return append(ServicePrefix, []byte(sid)...)
}

//func GetServiceIssuerKey(sid string) []byte {
//	return append(ServiceIssuerPrefix, []byte(sid)...)
//}

func GetCredentialPrefixByDid(did string) []byte {
	return append(CredentialPrefix, []byte(did)...)
}

func GetCredentialKey(did, sid string) []byte {
	return append(GetCredentialPrefixByDid(did), sid...)
}

func GetFilterLoggerKey(did, sid string) []byte {
	return append(append(FilterLoggerPrefix, did...), sid...)
}

func appendLengthPrefixed(dst, value []byte) []byte {
	var lenBz [8]byte
	binary.BigEndian.PutUint64(lenBz[:], uint64(len(value)))
	dst = append(dst, lenBz[:]...)
	return append(dst, value...)
}

func GetLegacyFilterPrefixBySidAndFilter(sid string, filter []byte) []byte {
	return append(append(append([]byte{}, FilterPrefix...), sid...), filter...)
}

func GetLegacyFilterKey(sid, did string, filter []byte) []byte {
	return append(GetLegacyFilterPrefixBySidAndFilter(sid, filter), did...)
}

func GetFilterPrefixBySidAndFilter(sid string, filter []byte) []byte {
	key := append([]byte{}, FilterPrefix...)
	key = appendLengthPrefixed(key, []byte(sid))
	return appendLengthPrefixed(key, filter)
}

func GetFilterKey(sid, did string, filter []byte) []byte {
	return append(GetFilterPrefixBySidAndFilter(sid, filter), did...)
}
