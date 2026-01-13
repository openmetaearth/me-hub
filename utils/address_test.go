package utils

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestMeBech32ToEthAndTron(t *testing.T) {
	raw := common.HexToAddress("0x2Af9729ca13181E775Db6A9391d7166D58Cfc7b1").Bytes()

	meAddr := "me19tuh989pxxq7wawmd2fer4ckd4vvl3a3sepez8"
	ethAddr, err := MeBech32ToEth(meAddr)
	require.NoError(t, err)
	require.Equal(t, "0x2Af9729ca13181E775Db6A9391d7166D58Cfc7b1", ethAddr)

	tronAddr, err := MeBech32ToTron(meAddr)
	require.NoError(t, err)
	require.NotEmpty(t, tronAddr)

	decoded := base58.Decode(tronAddr)
	require.Len(t, decoded, 1+20+4)
	require.Equal(t, byte(0x41), decoded[0])
	require.True(t, bytes.Equal(raw, decoded[1:21]))

	payload := decoded[:21]
	h1 := sha256.Sum256(payload)
	h2 := sha256.Sum256(h1[:])
	require.True(t, bytes.Equal(h2[0:4], decoded[21:]))

	expectedTron := base58.Encode(append(payload, h2[0:4]...))
	require.Equal(t, expectedTron, tronAddr)
}

func TestMeBech32ToEth_Error(t *testing.T) {
	_, err := MeBech32ToEth("me1invalidxyz")
	require.Error(t, err)
}

func TestMeBech32ToTron_Error(t *testing.T) {
	_, err := MeBech32ToTron("me1badaddr")
	require.Error(t, err)
}
