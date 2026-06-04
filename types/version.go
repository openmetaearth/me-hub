package types

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

const (
	MainnetV1ChainId  = "mechain"
	MainnetV2ChainId  = "mechain_202404-1"
	MainnetEvmChainID = 202404

	TestnetV1ChainId  = "mechain_testnet"
	TestnetV2ChainId  = "mechain_testnet_202405-1"
	TestnetEvmChainID = 202405
)

var (
	chainId            = MainnetV1ChainId
	once               sync.Once
	eip155SuffixRegexp = regexp.MustCompile(`_[0-9]+-[0-9]+$`)
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
	curr := ChainId()
	if curr == "" {
		return ""
	}
	if eip155SuffixRegexp.MatchString(curr) {
		return curr
	}
	if strings.Contains(curr, "testnet") {
		return fmt.Sprintf("%s_%d-1", curr, TestnetEvmChainID)
	}
	return fmt.Sprintf("%s_%d-1", curr, MainnetEvmChainID)
}
