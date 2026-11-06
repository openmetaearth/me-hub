package types

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/testutil/sample"
)

func minVotingPeriodTestMessage() MsgSubmitProposal {
	return MsgSubmitProposal{
		Proposer: sample.AccAddress(), Kind: ProposalKindUpdateMinVotingPeriod,
		Title: "update minimum voting period", VotingPeriod: 60,
		UpdateMinVotingPeriod: &UpdateMinVotingPeriodPayload{MinVotingPeriod: 120},
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
		UpdateMinVotingPeriod: msg.UpdateMinVotingPeriod,
	}}
	return gs
}

func TestUpdateMinVotingPeriodBoundsAndRoundTrip(t *testing.T) {
	require.Equal(t, ProposalKind(3), ProposalKindUpdateMinVotingPeriod)
	require.Equal(t, "PROPOSAL_KIND_UPDATE_MIN_VOTING_PERIOD", ProposalKindUpdateMinVotingPeriod.String())
	for _, value := range []uint64{0, 1, 60, 120, MaxVotingPeriod, MaxVotingPeriod + 1} {
		t.Run(fmt.Sprint(value), func(t *testing.T) {
			msg := minVotingPeriodTestMessage()
			msg.UpdateMinVotingPeriod.MinVotingPeriod = value
			gs := minVotingPeriodTestGenesis(msg)
			if value == 0 || value > MaxVotingPeriod {
				require.Error(t, msg.ValidateBasic())
				require.Error(t, gs.Validate())
				return
			}
			require.NoError(t, msg.ValidateBasic())
			require.NoError(t, gs.Validate())
			for _, message := range []proto.Message{&msg, &gs.Proposals[0]} {
				wire, err := proto.Marshal(message)
				require.NoError(t, err)
				decoded := proto.Clone(message)
				decoded.Reset()
				require.NoError(t, proto.Unmarshal(wire, decoded))
				require.Equal(t, message, decoded)
			}
		})
	}
}

func TestUpdateMinVotingPeriodPayloadExclusivity(t *testing.T) {
	for _, kind := range []ProposalKind{ProposalKindTreasuryWithdraw, ProposalKindReplaceRegionOperator, ProposalKindUpdateMinVotingPeriod} {
		for mask := 0; mask < 8; mask++ {
			t.Run(fmt.Sprintf("%s/payloads-%03b", kind, mask), func(t *testing.T) {
				msg := minVotingPeriodTestMessage()
				msg.Kind = kind
				msg.UpdateMinVotingPeriod = nil
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
				if mask&4 != 0 {
					msg.UpdateMinVotingPeriod = &UpdateMinVotingPeriodPayload{MinVotingPeriod: 120}
				}
				valid := kind == ProposalKindTreasuryWithdraw && mask == 1 ||
					kind == ProposalKindReplaceRegionOperator && mask == 2 ||
					kind == ProposalKindUpdateMinVotingPeriod && mask == 4
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

func TestUpdateMinVotingPeriodGenesisHistoricalDuration(t *testing.T) {
	for _, status := range []ProposalStatus{ProposalStatusVoting, ProposalStatusPassed, ProposalStatusRejected} {
		for _, period := range []uint64{0, 1, 60, MaxVotingPeriod, MaxVotingPeriod + 1} {
			t.Run(fmt.Sprintf("%s/%d", status, period), func(t *testing.T) {
				gs := minVotingPeriodTestGenesis(minVotingPeriodTestMessage())
				gs.Proposals[0].Status = status
				gs.Proposals[0].VotingPeriod = period
				// Avoid overflow in the malformed >Max case; duration bounds
				// must reject it independently of this otherwise valid deadline.
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
