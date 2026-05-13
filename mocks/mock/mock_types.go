package mock

import (
	sdkmath "cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Delegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Delegation struct {
	// delegator_address is the bech32-encoded address of the delegator.
	stakingtypes.Delegation
	StartHeight  int64
	Amount       sdkmath.Int `protobuf:"bytes,5,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount"`
	Unmovable    sdkmath.Int `protobuf:"bytes,6,opt,name=unmovable,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"unmovable"`
	UnMeidAmount sdkmath.Int `protobuf:"bytes,7,opt,name=unMeidAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"unMeidAmount"`
}
