package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Msg represents a cross-chain message
type Msg struct {
	Amount    sdk.Int    `json:"amount"`
	Recipient string    `json:"recipient"`
	Nonce     uint64    `json:"nonce"`
	SourceChain string  `json:"source_chain"`
	Hash      [32]byte `json:"hash"`
}

// String returns a string representation of the message
func (m Msg) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Hash returns the hash of the message
func (m Msg) Hash() [32]byte {
	return sha256.Sum256([]byte(m.String()))
}