package v2_0_10_test

import (
	"github.com/st-chain/me-hub/app/upgrades/v2_0_10"
)

type MockDIDReader struct {
	Data map[string]v2_0_10.DidData
	Err  error
}

func (m MockDIDReader) ReadDID(filePath string) (map[string]v2_0_10.DidData, error) {
	return m.Data, m.Err
}
