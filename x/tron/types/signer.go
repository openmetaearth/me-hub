package types

import (
	"crypto/ecdsa"

	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

const tronSignaturePrefix = "\x19TRON Signed Message:\n32"

// NewTronSignature creates a new signature over a given byte array
func NewTronSignature(hash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, "private key")
	}
	protectedHash := crypto.Keccak256Hash(append([]byte(tronSignaturePrefix), hash...))
	return crypto.Sign(protectedHash.Bytes(), privateKey)
}

// TronAddressFromSignature recovers the Tron address from a signature and hash.
// It handles Ethereum-style V value normalization (27/28 -> 0/1) for backwards compatibility.
// Returns the Tron address as a string or an error if signature recovery fails.
// Note: This function modifies the signature slice in place for V value normalization.
func TronAddressFromSignature(hash []byte, signature []byte) (string, error) {
	if len(signature) < 65 {
		return "", errorsmod.Wrap(types.ErrInvalid, "signature too short")
	}
	// To verify signature
	// - use crypto.SigToPub to get the public key
	// - use crypto.PubkeyToAddress to get the address
	// - compare this to the address given.

	// for backwards compatibility reasons  the V value of an Ethereum sig is presented
	// as 27 or 28, internally though it should be a 0-3 value due to changed formats.
	// It seems that go-ethereum expects this to be done before sigs actually reach it's
	// internal validation functions. In order to comply with this requirement we check
	// the sig and if it's in standard format we correct it. If it's in go-ethereum's expected
	// format already we make no changes.
	//
	// We could attempt to break or otherwise exit early on obviously invalid values for this
	// byte, but that's a task best left to go-ethereum
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] = signature[64] - 27
	}

	protectedHash := crypto.Keccak256Hash(append([]byte(tronSignaturePrefix), hash...))
	pubkey, err := crypto.SigToPub(protectedHash.Bytes(), signature)
	if err != nil {
		return "", errorsmod.Wrap(err, "signature to public key")
	}

	addr := address.PubkeyToAddress(*pubkey)
	return addr.String(), nil
}

// ValidateTronSignature validates that the given signature was created by the private key
// corresponding to the provided Tron address. It recovers the address from the signature
// and hash, then compares it to the expected address. Returns an error if validation fails.
func ValidateTronSignature(hash []byte, signature []byte, tronAddress string) error {
	addr, err := TronAddressFromSignature(hash, signature)
	if err != nil {
		return err
	}
	if addr != tronAddress {
		return errorsmod.Wrap(types.ErrInvalid, "signature not matching")
	}
	return nil
}
