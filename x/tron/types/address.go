package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"

	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
)

var _ gravitytypes.ExternalAddress = TronAddress{}

type TronAddress struct{}

func (b TronAddress) ValidateExternalAddr(addr string) error {
	return ValidateTronAddress(addr)
}

func (b TronAddress) ExternalAddrToAccAddr(addr string) sdk.AccAddress {
	tronAddr, err := tronaddress.Base58ToAddress(addr)
	if err != nil {
		panic(err)
	}
	return tronAddr.Bytes()[1:]
}

func (b TronAddress) ExternalAddrToHexAddr(addr string) gethcommon.Address {
	tronAddr, err := tronaddress.Base58ToAddress(addr)
	if err != nil {
		panic(err)
	}
	return gethcommon.BytesToAddress(tronAddr.Bytes()[1:])
}

func (b TronAddress) ExternalAddrToStr(bz []byte) string {
	if len(bz) == gethcommon.AddressLength {
		bz = append([]byte{tronaddress.TronBytePrefix}, bz...)
	}
	return tronaddress.Address(bz).String()
}

// ValidateTronAddress validates the ethereum address strings
func ValidateTronAddress(address string) error {
	if address == "" {
		return errors.New("empty")
	}
	if len(address) != tronaddress.AddressLengthBase58 {
		return fmt.Errorf("invalid address length: expected %d chars, got %d", tronaddress.AddressLengthBase58, len(address))
	}
	tronAddr, err := common.DecodeCheck(address)
	if err != nil {
		return errors.New("doesn't pass format validation")
	}
	if len(tronAddr) != tronaddress.AddressLength {
		return fmt.Errorf("invalid address length: expected decoded %d bytes, got %d", tronaddress.AddressLength, len(tronAddr))
	}
	if tronAddr[0] != tronaddress.TronBytePrefix {
		return errors.New("invalid tron prefix")
	}
	expectAddress := common.EncodeCheck(tronAddr)
	if expectAddress != address {
		return fmt.Errorf("mismatch expected: %s, got: %s", expectAddress, address)
	}
	return nil
}
