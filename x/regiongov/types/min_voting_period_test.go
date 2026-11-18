package types

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/testutil/sample"
)

func minVotingPeriodTestMessage() MsgSubmitProposal {
	return MsgSubmitProposal{
		Proposer: sample.AccAddress(), Kind: ProposalKindTreasuryWithdraw,
		Title: "treasury withdrawal", VotingPeriod: 60,
		Treasury: &TreasuryWithdrawPayload{Withdrawals: []Withdrawal{{
			RegionId: "usa", Receiver: sample.AccAddress(), Amount: sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
		}}},
	}
}

func minVotingPeriodTestGenesis(msg MsgSubmitProposal) *GenesisState {
	gs := DefaultGenesis()
	// A current minimum must not be applied retroactively to stored proposals.
	gs.Params.MinVotingPeriod = 120
	gs.ProposalSeq = 1
	now := time.Unix(1_700_000_000, 0).UTC()
	gs.Proposals = []Proposal{{
		Id: 1, Proposer: msg.Proposer, Kind: msg.Kind, Title: msg.Title,
		Status: ProposalStatusVoting, SubmitTime: now,
		VotingEndTime: now.Add(time.Duration(msg.VotingPeriod) * time.Second),
		VotingPeriod:  msg.VotingPeriod, Electorate: ElectorateKindCountryRegionOperators,
		Treasury: msg.Treasury, ReplaceOperator: msg.ReplaceOperator,
	}}
	return gs
}

func TestMinVotingPeriodParameterBounds(t *testing.T) {
	for _, value := range []uint64{0, 1, 60, 120, MaxVotingPeriod, MaxVotingPeriod + 1} {
		t.Run(fmt.Sprint(value), func(t *testing.T) {
			params := DefaultParams()
			params.MinVotingPeriod = value
			if value == 0 || value > MaxVotingPeriod {
				require.Error(t, params.Validate())
			} else {
				require.NoError(t, params.Validate())
			}
		})
	}
}

func TestProposalPayloadExclusivity(t *testing.T) {
	for _, kind := range []ProposalKind{ProposalKindTreasuryWithdraw, ProposalKindReplaceRegionOperator} {
		for mask := 0; mask < 4; mask++ {
			t.Run(fmt.Sprintf("%s/payloads-%02b", kind, mask), func(t *testing.T) {
				msg := minVotingPeriodTestMessage()
				msg.Kind = kind
				msg.Treasury = nil
				if mask&1 != 0 {
					msg.Treasury = &TreasuryWithdrawPayload{Withdrawals: []Withdrawal{{
						RegionId: "usa", Receiver: sample.AccAddress(), Amount: sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
					}}}
				}
				if mask&2 != 0 {
					msg.ReplaceOperator = &ReplaceRegionOperatorPayload{
						RegionId: "usa", NewOperator: sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String(),
					}
				}
				valid := kind == ProposalKindTreasuryWithdraw && mask == 1 ||
					kind == ProposalKindReplaceRegionOperator && mask == 2
				if valid {
					require.NoError(t, msg.ValidateBasic())
					require.NoError(t, minVotingPeriodTestGenesis(msg).Validate())
				} else {
					require.Error(t, msg.ValidateBasic())
					require.Error(t, minVotingPeriodTestGenesis(msg).Validate())
				}
			})
		}
	}
}

func TestGenesisHistoricalVotingDuration(t *testing.T) {
	for _, status := range []ProposalStatus{ProposalStatusVoting, ProposalStatusPassed, ProposalStatusRejected} {
		for _, period := range []uint64{0, 1, 60, MaxVotingPeriod, MaxVotingPeriod + 1} {
			t.Run(fmt.Sprintf("%s/%d", status, period), func(t *testing.T) {
				gs := minVotingPeriodTestGenesis(minVotingPeriodTestMessage())
				gs.Proposals[0].Status = status
				gs.Proposals[0].VotingPeriod = period
				gs.Proposals[0].VotingEndTime = gs.Proposals[0].SubmitTime.Add(time.Minute)
				if period == 0 || period > MaxVotingPeriod {
					require.Error(t, gs.Validate())
				} else {
					require.NoError(t, gs.Validate())
				}
			})
		}
	}
}
