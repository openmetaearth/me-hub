package v2_0_10

import (
	"encoding/json"
	"io/ioutil"
)

type NftReader interface {
	ReadNft(filePath string) (map[string]map[string]NftUri, error)
}

type RealNftReader struct{}

type NftUri struct {
	URI     string `json:"uri"`
	URIHash string `json:"uri_hash"`
}

func (r RealNftReader) ReadNft(filePath string) (map[string]map[string]NftUri, error) {
	data := map[string]map[string]NftUri{}
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
