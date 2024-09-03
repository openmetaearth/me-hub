package types

import (
	"fmt"
	"sync"
)

const (
	V1ChainId  = "mechain"
	V2ChainId  = "mechain_100-1"
	EvmChainID = 100
)

var (
	chainId = V1ChainId
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
	return fmt.Sprintf("%s_%d-1", ChainId(), EvmChainID)
}
