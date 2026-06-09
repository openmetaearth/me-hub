package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgWithdrawDelegatorReward defines the MsgWithdrawDelegatorReward request type.
type MsgWithdrawDelegatorReward struct {
	DelegatorAddress string         `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address"`
	ValidatorAddress string         `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address"`
	Reward           sdk.Coin        `protobuf:"bytes,3,opt,name=reward,proto3" json:"reward"`
}