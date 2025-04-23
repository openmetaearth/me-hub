package v2_0_10_test

type MockKycPubkeyReader struct {
	Data map[string]string
	Err  error
}

func (m MockKycPubkeyReader) ReadKycPubkey(filePath string) (map[string]string, error) {
	return m.Data, m.Err
}
