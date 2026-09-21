package types

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	iden3utils "github.com/iden3/go-iden3-crypto/v2/utils"
)

const (
	ZkCoreIDLength                 = 31
	MaxZkAuthProofLength           = 64 * 1024
	ZkCoreIDBindingChallengeDomain = "me-hub:bond-zk-core-id:v1"
)

// ValidateZkCoreID validates the canonical 31-byte iden3 ID and checksum.
func ValidateZkCoreID(coreID []byte) error {
	if len(coreID) != ZkCoreIDLength {
		return fmt.Errorf("core ID must be %d bytes", ZkCoreIDLength)
	}

	var checksum uint16
	for _, b := range coreID[:ZkCoreIDLength-2] {
		checksum += uint16(b)
	}
	provided := binary.LittleEndian.Uint16(coreID[ZkCoreIDLength-2:])
	if provided == 0 || provided != checksum {
		return fmt.Errorf("core ID checksum mismatch")
	}
	return nil
}

// ZkCoreIDToInt converts the canonical little-endian iden3 ID to its circuit scalar.
func ZkCoreIDToInt(coreID []byte) (*big.Int, error) {
	if err := ValidateZkCoreID(coreID); err != nil {
		return nil, err
	}
	return iden3utils.SetBigIntFromLEBytes(new(big.Int), coreID), nil
}

// CalculateZkCoreIDBindingChallenge returns the AuthV3 challenge for a binding request.
// It is keccak256(abi.encode(domain, chainID, creator, did, coreID, nonce, expiresAt))
// with the high nibble cleared so the result is safely inside the circuit field.
func CalculateZkCoreIDBindingChallenge(
	chainID string,
	creator string,
	did string,
	coreID []byte,
	nonce uint64,
	expiresAt int64,
) (*big.Int, error) {
	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, err
	}
	bytesType, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return nil, err
	}
	uint64Type, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return nil, err
	}
	int64Type, err := abi.NewType("int64", "", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments{
		{Type: stringType}, {Type: stringType}, {Type: stringType}, {Type: stringType},
		{Type: bytesType}, {Type: uint64Type}, {Type: int64Type},
	}
	encoded, err := args.Pack(
		ZkCoreIDBindingChallengeDomain,
		chainID,
		creator,
		did,
		coreID,
		nonce,
		expiresAt,
	)
	if err != nil {
		return nil, err
	}

	hash := crypto.Keccak256(encoded)
	hash[0] &= 0x0f
	return new(big.Int).SetBytes(hash), nil
}
