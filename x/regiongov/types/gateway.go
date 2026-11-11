package types

import (
	// Register standard protobuf descriptors so grpc-gateway can discover all
	// regiongov enums without maintaining a separate registration list.
	_ "github.com/openmetaearth/me-hub/api/metaearth/regiongov"
)
