package types

const (
	// ModuleName defines the module name
	ModuleName = "kyc"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_kyc"
)

const RegionIdLength = 2

var (
	ParamsKey    = []byte{0x01}
	RegionPrefix = []byte{0x10}
)

func GetRegionKey(regionId string) []byte {
	if len(regionId) != RegionIdLength {
		panic("Invalid region code")
	}
	return append(RegionPrefix, []byte(regionId)...)
}

const (
	KycEventSeqKey = "KycEventSeq/value/"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
