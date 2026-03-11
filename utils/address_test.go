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

func TestMe2Eth(t *testing.T) {
	me2Eth := make(map[string]string)
	me2Eth["me18g0j259s5l8sd2kzfy790t83t8gf4xd0wxuzyt"] = "0x3a1f2550b0A7cf06AAC2493c57AcF159D09A99AF"
	me2Eth["me1fzztxku3kj5s7r00p0gnrr60qlq6vvq94n0tph"] = "0x4884b35b91B4a90F0dEf0BD1318F4f07c1a63005"
	me2Eth["me17p8sa9lwk8334eul80e2md0t7hxsdkjrk87slk"] = "0xF04f0e97eeB1E31AE79F3bF2ADB5ebF5Cd06da43"
	me2Eth["me1kw3l6mk8t4z638am7zmqw3atvw764kxa3vxtpy"] = "0xB3A3fD6ec75D45a89Fbbf0B60747Ab63BDAAD8Dd"
	me2Eth["me1ezcc7r3wf8x6lyxryslwq3qpjhn6ph2dhxrlz6"] = "0xc8B18F0E2E49cDaf90C3243eE0440195E7A0dd4D"
	for meAddr, expectedEthAddr := range me2Eth {
		ethAddr, err := MeBech32ToEth(meAddr)
		require.NoError(t, err)
		require.Equal(t, expectedEthAddr, ethAddr)
	}

	me2Tron := make(map[string]string)
	me2Tron["me18g0j259s5l8sd2kzfy790t83t8gf4xd0wxuzyt"] = "TFGXWELFxbabByrzSPDwr7QJj5mPBUXw9Z"
	me2Tron["me1fzztxku3kj5s7r00p0gnrr60qlq6vvq94n0tph"] = "TGaed55sL7X49DzfXbxQxvQcrhJ7Qi1h11"
	me2Tron["me17p8sa9lwk8334eul80e2md0t7hxsdkjrk87slk"] = "TXsqt5KmeMPViJoCyyntzkcNrULok17bYC"
	me2Tron["me1kw3l6mk8t4z638am7zmqw3atvw764kxa3vxtpy"] = "TSM4Qctg1pmTZitMkj2tobdpZvZXtNH9Z4"
	me2Tron["me1ezcc7r3wf8x6lyxryslwq3qpjhn6ph2dhxrlz6"] = "TUGNrRkLmY2bpF8F7ukRnDQRt3D5PsvoQ6"
	for meAddr, expectedTronAddr := range me2Tron {
		tronAddr, err := MeBech32ToTron(meAddr)
		require.NoError(t, err)
		require.Equal(t, expectedTronAddr, tronAddr)
	}
}

func TestMeBech32ToEth_Error(t *testing.T) {
	_, err := MeBech32ToEth("me1invalidxyz")
	require.Error(t, err)
}

func TestMeBech32ToTron_Error(t *testing.T) {
	_, err := MeBech32ToTron("me1badaddr")
	require.Error(t, err)
}
