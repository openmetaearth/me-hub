package v2_0_10_test

import "github.com/st-chain/me-hub/app/upgrades/v2_0_10"

type MockNftReader struct {
	Data map[string]v2_0_10.ClassUri
	Err  error
}

func (m MockNftReader) ReadNft(filePath string) (map[string]v2_0_10.ClassUri, error) {
	return m.Data, m.Err
}
