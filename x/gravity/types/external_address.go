package types

import (
	"fmt"
	"regexp"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

type ExternalAddress interface {
	ValidateExternalAddr(addr string) error
	ExternalAddrToAccAddr(addr string) sdk.AccAddress
	ExternalAddrToHexAddr(addr string) common.Address
	ExternalAddrToStr(bz []byte) string
}

var reModuleName *regexp.Regexp

func init() {
	reModuleNameString := `[a-zA-Z][a-zA-Z0-9/]{1,32}`
	reModuleName = regexp.MustCompile(fmt.Sprintf(`^%s$`, reModuleNameString))
}

// ValidateModuleName is the default validation function for crosschain moduleName.
func ValidateModuleName(moduleName string) error {
	if !reModuleName.MatchString(moduleName) {
		return fmt.Errorf("invalid module name: %s", moduleName)
	}
	return nil
}

var externalAddressRouter = make(map[string]ExternalAddress)

func GetSupportChains() []string {
	chains := make([]string, 0, len(externalAddressRouter))
	for chainName := range externalAddressRouter {
		chains = append(chains, chainName)
	}
	sort.SliceStable(chains, func(i, j int) bool {
		return chains[i] < chains[j]
	})
	return chains
}

func RegisterExternalAddress(chainName string, validate ExternalAddress) {
	if err := ValidateModuleName(chainName); err != nil {
		panic(errortypes.ErrInvalidRequest.Wrapf("invalid chain name: %s", chainName))
	}
	if _, ok := externalAddressRouter[chainName]; ok {
		panic(fmt.Sprintf("duplicate registry msg validateBasic! chainName: %s", chainName))
	}
	externalAddressRouter[chainName] = validate
}

func ValidateExternalAddr(chainName, addr string) error {
	router, ok := externalAddressRouter[chainName]
	if !ok {
		return fmt.Errorf("unrecognized cross chain name: %s", chainName)
	}
	return router.ValidateExternalAddr(addr)
}

func ExternalAddrToAccAddr(chainName, addr string) sdk.AccAddress {
	router, ok := externalAddressRouter[chainName]
	if !ok {
		panic("unrecognized cross chain name: " + chainName)
	}
	return router.ExternalAddrToAccAddr(addr)
}

func ExternalAddrToHexAddr(chainName, addr string) common.Address {
	router, ok := externalAddressRouter[chainName]
	if !ok {
		panic("unrecognized cross chain name: " + chainName)
	}
	return router.ExternalAddrToHexAddr(addr)
}

func ExternalAddrToStr(chainName string, bz []byte) string {
	router, ok := externalAddressRouter[chainName]
	if !ok {
		panic("unrecognized cross chain name: " + chainName)
	}
	return router.ExternalAddrToStr(bz)
}
