package v2_0_10

import (
	"encoding/json"
	"io/ioutil"
)

type NftReader interface {
	ReadNft(filePath string) (map[string]ClassUri, error)
}

type RealNftReader struct{}

type ClassUri struct {
	ClassURI     string            `json:"class_uri"`
	ClassURIHash string            `json:"class_uri_hash"`
	NftData      map[string]NftUri `json:"nft_data"`
}

type NftUri struct {
	URI     string `json:"nft_uri"`
	URIHash string `json:"nft_uri_hash"`
}

func (r RealNftReader) ReadNft(filePath string) (map[string]ClassUri, error) {
	data := map[string]ClassUri{}
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
