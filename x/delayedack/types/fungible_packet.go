package types

import (
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

type TransferDataWithFinalization struct {
	rollapptypes.TransferData
	// Proof height is only be populated if and only if the rollappID is not empty
	ProofHeight uint64
	// Finalized is only be populated if and only if the rollappID is not empty
	Finalized bool
}
