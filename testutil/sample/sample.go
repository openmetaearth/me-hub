package sample

import (
	"crypto/sha256"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccAddress returns a sample account address
func AccAddress() string {
	return Acc().String()
}

func Acc() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr)
}

// AccAddressFromSecret returns a deterministic account address for a secret.
func AccAddressFromSecret(secret string) string {
	seed := sha256.Sum256([]byte(secret))
	pk := &ed25519.PrivKey{Key: seed[:]}
	return sdk.AccAddress(pk.PubKey().Address()).String()
}

// GenerateAddresses generates numOfAddresses bech32 address
func GenerateAddresses(numOfAddresses int) []string {
	addresses := []string{}
	for i := 0; i < numOfAddresses; i++ {
		addresses = append(addresses, AccAddress())
	}
	return addresses
}
