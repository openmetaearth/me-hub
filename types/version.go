package types

import (
	"fmt"
	"strings"
	"sync"
)

const (
	MainnetV1ChainId  = "mechain"
	MainnetV2ChainId  = "mechain_100-1"
	MainnetEvmChainID = 100

	TestnetV1ChainId  = "mechain_testnet"
	TestnetV2ChainId  = "mechain_testnet_101-1"
	TestnetEvmChainID = 101
)

var (
	chainId = MainnetV1ChainId
	once    sync.Once
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
	if strings.Contains(ChainId(), "testnet") {
		return fmt.Sprintf("%s_%d-1", ChainId(), TestnetEvmChainID)
	}
	return fmt.Sprintf("%s_%d-1", ChainId(), MainnetEvmChainID)
}
