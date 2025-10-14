package types

// Module name and store keys
const (
	// ModuleName is the name of the module
	ModuleName = "blacklist"

	// StoreKey is the store key string for the module
	StoreKey = ModuleName

	// RouterKey is the message route for the module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute = ModuleName
)

// Key prefixes
var (
	// BlacklistKeyPrefix is the prefix for storing blacklist
	BlacklistKeyPrefix = []byte{0x01}
)

// Keys for store prefixes
func KeyPrefix(p string) []byte {
	return []byte(p)
}
