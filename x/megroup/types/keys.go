package types

const (
	// ModuleName defines the module name
	ModuleName = "megroup"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_megroup"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	GroupKey       = "Group/value/"
	GroupLastIDKey = "Group/LastID/"
	GroupRegionKey = "Group/RegionID/"
)

const (
	GroupMemberKey = "GroupMember/value/"
)
