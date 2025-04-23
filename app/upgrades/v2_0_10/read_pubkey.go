package v2_0_10

import (
	"encoding/json"
	"io/ioutil"
)

type KycPubkeyReader interface {
	ReadKycPubkey(filePath string) (map[string]string, error)
}

type RealKycPubkeyReader struct{}

func (r RealKycPubkeyReader) ReadKycPubkey(filePath string) (map[string]string, error) {
	data := make(map[string]string)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
