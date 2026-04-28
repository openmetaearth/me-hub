package types

import gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"

const (
	// ModuleName is the name of the module
	ModuleName = "bsc"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

func init() {
	gravitytypes.RegisterExternalAddress(ModuleName, gravitytypes.EthereumAddress{})
}
