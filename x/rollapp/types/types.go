package types

import (
	"encoding/json"

	common "github.com/openmetaearth/me-hub/x/common/types"
)

type StateStatus common.Status

func (md *RollappMetadata) IsEmpty() bool {
	bz, _ := json.Marshal(md)
	return string(bz) == "{}"
}
