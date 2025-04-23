package v2_0_10

import (
	"encoding/json"
	"io/ioutil"
)

type DIDReader interface {
	ReadDID(filePath string) (map[string]DidData, error)
}

type RealDIDReader struct{}

type DidData struct {
	Did        string `json:"did"`
	Level      uint64 `json:"level"`
	Uri        string `json:"uri"`
	UriHash    string `json:"uri_hash"`
	KycUri     string `json:"kyc_uri"`
	KycUriHash string `json:"kyc_uri_hash"`
	PubKey     string `json:"pubkey"`
}

func (r RealDIDReader) ReadDID(filePath string) (map[string]DidData, error) {
	data := make(map[string]DidData)
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
