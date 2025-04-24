package v2_0_10

import (
	"encoding/json"
	"fmt"
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
	mockData := map[string]map[string]NftUri{
		"collection1": {
			"nft1": {URI: "https://example.com/nft1", URIHash: "hash1"},
			"nft2": {URI: "https://example.com/nft2", URIHash: "hash2"},
		},
		"collection2": {
			"nft1": {URI: "https://example.com/nft1", URIHash: "hash1"},
			"nft2": {URI: "https://example.com/nft2", URIHash: "hash2"},
		},
		"collection3": {
			"nft1": {URI: "https://example.com/nft1", URIHash: "hash1"},
			"nft2": {URI: "https://example.com/nft2", URIHash: "hash2"},
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
