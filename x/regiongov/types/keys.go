package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "regiongov"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

var (
	ParamsKey           = []byte{0x01}
	ProposalSeqKey      = []byte{0x02}
	PendingActionSeqKey = []byte{0x03}

	ProposalPrefix      = []byte{0x10}
	ElectorPrefix       = []byte{0x11}
	VotePrefix          = []byte{0x12}
	EquityPrefix        = []byte{0x13}
	RegionFreezePrefix  = []byte{0x14}
	EquityFreezePrefix  = []byte{0x15}
	PendingActionPrefix = []byte{0x16}
	// VotingQueuePrefix indexes open proposals by voting end time.
	VotingQueuePrefix = []byte{0x17}
)

// Uint64ToBytes encodes an id so store iteration follows numeric order.
func Uint64ToBytes(v uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, v)
	return bz
}

// ProposalKey returns the store key of one proposal.
func ProposalKey(id uint64) []byte {
	return append(append([]byte{}, ProposalPrefix...), Uint64ToBytes(id)...)
}

// ElectorKey returns the store key of one snapshotted elector.
func ElectorKey(proposalID uint64, address string) []byte {
	return compositeKey(ElectorPrefix, proposalID, address)
}

// ElectorPrefixByProposal returns the prefix of every elector on a proposal.
func ElectorPrefixByProposal(proposalID uint64) []byte {
	return append(append([]byte{}, ElectorPrefix...), Uint64ToBytes(proposalID)...)
}

// VoteKey returns the store key of one ballot.
func VoteKey(proposalID uint64, voter string) []byte {
	return compositeKey(VotePrefix, proposalID, voter)
}

// VotePrefixByProposal returns the prefix of every ballot on a proposal.
func VotePrefixByProposal(proposalID uint64) []byte {
	return append(append([]byte{}, VotePrefix...), Uint64ToBytes(proposalID)...)
}

// EquityKey returns the store key of one region's equity list.
func EquityKey(regionID string) []byte {
	return append(append([]byte{}, EquityPrefix...), []byte(regionID)...)
}

// RegionFreezeKey returns the store key of one region freeze.
func RegionFreezeKey(regionID string) []byte {
	return append(append([]byte{}, RegionFreezePrefix...), []byte(regionID)...)
}

// EquityFreezeKey returns the store key of one equity-address freeze.
func EquityFreezeKey(regionID, address string) []byte {
	key := append(append([]byte{}, EquityFreezePrefix...), []byte(regionID)...)
	key = append(key, 0)
	return append(key, []byte(address)...)
}

// EquityFreezePrefixByRegion returns the prefix of equity freezes in one region.
func EquityFreezePrefixByRegion(regionID string) []byte {
	key := append(append([]byte{}, EquityFreezePrefix...), []byte(regionID)...)
	return append(key, 0)
}

// PendingActionKey returns the store key of one pending action.
func PendingActionKey(id uint64) []byte {
	return append(append([]byte{}, PendingActionPrefix...), Uint64ToBytes(id)...)
}

// VotingQueueKey returns the store key of one proposal in the end-time queue.
func VotingQueueKey(end time.Time, id uint64) []byte {
	key := append(append([]byte{}, VotingQueuePrefix...), sdk.FormatTimeBytes(end)...)
	return append(key, Uint64ToBytes(id)...)
}

// VotingQueueTimePrefix returns the prefix of every queue entry ending at end.
func VotingQueueTimePrefix(end time.Time) []byte {
	return append(append([]byte{}, VotingQueuePrefix...), sdk.FormatTimeBytes(end)...)
}

func compositeKey(prefix []byte, id uint64, address string) []byte {
	key := append(append([]byte{}, prefix...), Uint64ToBytes(id)...)
	return append(key, []byte(address)...)
}
