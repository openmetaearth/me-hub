package types

import (
	"fmt"
	"regexp"
	"sync"
)

const (
	MainnetV1ChainId  = "me-chain"
	MainnetV2ChainId  = "mechain_2404-1"
	MainnetEvmChainID = 2404
)

var (
	chainId         = MainnetV1ChainId
	once            sync.Once
	eip155SuffixReg = regexp.MustCompile(`_[1-9][0-9]*-[1-9][0-9]*$`)
)

func SetChainId(id string) {
	once.Do(func() {
		chainId = id
	})
}

func ChainId() string {
	return chainId
}

func ChainIdWithEIP155() string {
	return ChainIdWithEIP155From(ChainId())
}

// ChainIdWithEIP155From maps a Cosmos chain-id to the EIP-155 format expected by Ethermint.
// Legacy mainnet chains use the plain "me-chain" id; EVM uses chain id 2404 (mechain_2404-1).
func ChainIdWithEIP155From(id string) string {
	if eip155SuffixReg.MatchString(id) {
		return id
	}
	if id == MainnetV1ChainId {
		return MainnetV2ChainId
	}
	return fmt.Sprintf("%s_%d-1", id, MainnetEvmChainID)
}
