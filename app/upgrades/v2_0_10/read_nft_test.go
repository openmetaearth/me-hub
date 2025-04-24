package v2_0_10

import (
	"encoding/json"
	"fmt"
	"github.com/st-chain/me-hub/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadNft(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test_nft_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up the file after the test

	// Prepare mock NFT data
	mockData := map[string]ClassUri{
		"class1": {
			ClassURI:     "https://example.com/class1",
			ClassURIHash: utils.CalculateUriHash("https://example.com/class1"),
			NftData: map[string]NftUri{
				"nft1": {URI: "https://example.com/nft1", URIHash: utils.CalculateUriHash("https://example.com/nft1")},
				"nft2": {URI: "https://example.com/nft2", URIHash: utils.CalculateUriHash("https://example.com/nft2")},
			},
		},
		"class2": {
			ClassURI:     "https://example.com/class2",
			ClassURIHash: utils.CalculateUriHash("https://example.com/class1"),
			NftData: map[string]NftUri{
				"nft1": {URI: "https://example.com/nft1", URIHash: utils.CalculateUriHash("https://example.com/nft1")},
				"nft2": {URI: "https://example.com/nft2", URIHash: utils.CalculateUriHash("https://example.com/nft2")},
			},
		},
		"class3": {
			ClassURI:     "https://example.com/class3",
			ClassURIHash: utils.CalculateUriHash("https://example.com/class1"),
			NftData: map[string]NftUri{
				"nft1": {URI: "https://example.com/nft1", URIHash: utils.CalculateUriHash("https://example.com/nft1")},
				"nft2": {URI: "https://example.com/nft2", URIHash: utils.CalculateUriHash("https://example.com/nft2")},
			},
		},
	}

	// Write mock data to the temporary file
	mockDataBytes, err := json.Marshal(mockData)
	require.NoError(t, err)
	_, err = tempFile.Write(mockDataBytes)
	require.NoError(t, err)
	require.NoError(t, tempFile.Close())
	t.Log(string(mockDataBytes))

	// Test the ReadNft function
	reader := RealNftReader{}
	result, err := reader.ReadNft(tempFile.Name())
	fmt.Println(tempFile.Name())
	require.NoError(t, err)
	require.Equal(t, mockData, result)
}
