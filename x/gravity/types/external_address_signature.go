package types

import (
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

const bondedRelayerExternalAddressDomain = "metaearth.gravity.v1.MsgBondedRelayer.ExternalAddress"

func GetBondedRelayerExternalAddressCheckpoint(gravityID, chainName, relayerAddress, externalAddress string) []byte {
	fields := []string{
		bondedRelayerExternalAddressDomain,
		"gravity_id=" + gravityID,
		"chain_name=" + chainName,
		"relayer_address=" + relayerAddress,
		"external_address=" + externalAddress,
	}
	return crypto.Keccak256([]byte(strings.Join(fields, "\x00")))
}
