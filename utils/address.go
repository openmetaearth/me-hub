package utils

import (
	"crypto/sha256"
	"errors"
	"github.com/btcsuite/btcutil/base58"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
)

func ParseAddress(addr string) (accAddr sdk.AccAddress, isEvmAddr bool, err error) {
	_, bytes, decodeErr := bech32.DecodeAndConvert(addr)
	if decodeErr == nil {
		return sdk.AccAddress(bytes), false, nil
	}
	ethAddrError := ValidateEthereumAddress(addr)
	if ethAddrError == nil {
		return sdk.AccAddress(common.HexToAddress(addr).Bytes()), true, nil
	}
	return nil, false, errors.Join(decodeErr, ethAddrError)
}

func MeBech32ToEth(addr string) (string, error) {
	_, bytes, err := bech32.DecodeAndConvert(addr)
	if err != nil {
		return "", err
	}
	if len(bytes) != 20 {
		return "", errors.New("address length must be 20 bytes")
	}
	return ToChecksummed(bytes), nil
}

// MeBech32ToTron me1... Bech32 to Tron(Base58Check)
// Tron: 0x41 + 20 byte + SHA256
func MeBech32ToTron(addr string) (string, error) {
	_, raw, err := bech32.DecodeAndConvert(addr)
	if err != nil {
		return "", err
	}
	if len(raw) != 20 {
		return "", errors.New("address length must be 20 bytes")
	}

	payload := append([]byte{0x41}, raw...)

	h1 := sha256.Sum256(payload)
	h2 := sha256.Sum256(h1[:])
	full := append(payload, h2[0:4]...)

	return base58.Encode(full), nil
}
