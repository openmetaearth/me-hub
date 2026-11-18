package types

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"

	regiongovapi "github.com/openmetaearth/me-hub/api/metaearth/regiongov"
)

func TestRemovedProposalKindRejected(t *testing.T) {
	msg := minVotingPeriodTestMessage()
	msg.Kind = ProposalKind(3)
	msg.Treasury = nil
	wire, err := proto.Marshal(&msg)
	require.NoError(t, err)
	// Reproduce the former field 8 payload without depending on its removed Go type.
	wire = protowire.AppendTag(wire, 8, protowire.BytesType)
	wire = protowire.AppendBytes(wire, []byte{0x08, 0x78})
	var decoded MsgSubmitProposal
	require.NoError(t, proto.Unmarshal(wire, &decoded))
	require.Equal(t, ProposalKind(3), decoded.Kind)
	require.ErrorContains(t, decoded.ValidateBasic(), "unknown proposal kind")
}

func TestRemovedProposalWireNumbersReserved(t *testing.T) {
	kind := regiongovapi.ProposalKind(0).Descriptor()
	require.EqualValues(t, 1, ProposalKindTreasuryWithdraw)
	require.EqualValues(t, 2, ProposalKindReplaceRegionOperator)
	require.Equal(t, "PROPOSAL_KIND_TREASURY_WITHDRAW", string(kind.Values().ByNumber(1).Name()))
	require.Equal(t, "PROPOSAL_KIND_REPLACE_REGION_OPERATOR", string(kind.Values().ByNumber(2).Name()))
	require.Nil(t, kind.ParentFile().Messages().ByName("UpdateMinVotingPeriodPayload"))
	require.Nil(t, kind.Values().ByNumber(3))
	require.True(t, kind.ReservedRanges().Has(3))
	require.True(t, kind.ReservedNames().Has("PROPOSAL_KIND_UPDATE_MIN_VOTING_PERIOD"))
	proposal := (&regiongovapi.Proposal{}).ProtoReflect().Descriptor()
	require.Nil(t, proposal.Fields().ByNumber(15))
	require.True(t, proposal.ReservedRanges().Has(15))
	require.True(t, proposal.ReservedNames().Has("update_min_voting_period"))
	msg := (&regiongovapi.MsgSubmitProposal{}).ProtoReflect().Descriptor()
	require.Nil(t, msg.Fields().ByNumber(8))
	require.True(t, msg.ReservedRanges().Has(8))
	require.True(t, msg.ReservedNames().Has("update_min_voting_period"))
}

func TestRemovedProposalHistoricalGenesis(t *testing.T) {
	for _, status := range []ProposalStatus{ProposalStatusVoting, ProposalStatusPassed, ProposalStatusRejected, ProposalStatusFailed} {
		t.Run(status.String(), func(t *testing.T) {
			gs := minVotingPeriodTestGenesis(minVotingPeriodTestMessage())
			gs.Proposals[0].Kind = ProposalKind(3)
			gs.Proposals[0].Treasury = nil
			gs.Proposals[0].Status = status
			require.NoError(t, gs.Validate(), "historical records remain exportable/importable")
			gs.Proposals[0].Kind = ProposalKind(4)
			require.ErrorContains(t, gs.Validate(), "unknown proposal kind")
		})
	}
}
