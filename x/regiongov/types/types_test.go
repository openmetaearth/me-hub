package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
)

func TestParamsValidate(t *testing.T) {
	require.Equal(t, uint64(3600), DefaultParams().MinVotingPeriod)
	require.NoError(t, DefaultParams().Validate())

	params := DefaultParams()
	params.MinVotingPeriod = 0
	require.Error(t, params.Validate())

	params = DefaultParams()
	params.PassThresholdPercent = 0
	require.Error(t, params.Validate())

	params = DefaultParams()
	params.PassThresholdPercent = 101
	require.Error(t, params.Validate())

	params = DefaultParams()
	params.EquityWithdrawPercent = 0
	require.Error(t, params.Validate())
}

func TestDefaultGenesis(t *testing.T) {
	require.NoError(t, DefaultGenesis().Validate())
}

func TestGenesisProposalRoundTripRules(t *testing.T) {
	proposer := sample.AccAddress()
	receiver := sample.AccAddress()
	operator := sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String()
	now := time.Unix(1_700_000_000, 0).UTC()

	gs := DefaultGenesis()
	gs.Params.MinVotingPeriod = 1
	gs.ProposalSeq = 1
	gs.Proposals = []Proposal{{
		Id:            1,
		Kind:          ProposalKindTreasuryWithdraw,
		Proposer:      proposer,
		Title:         "withdraw",
		Status:        ProposalStatusVoting,
		SubmitTime:    now,
		VotingEndTime: now.Add(time.Hour),
		VotingPeriod:  1,
		Electorate:    ElectorateKindCountryRegionOperators,
		Treasury: &TreasuryWithdrawPayload{Withdrawals: []Withdrawal{{
			RegionId: "me_earth",
			Receiver: receiver,
			Amount:   sdk.NewCoins(sdk.NewInt64Coin("umec", 10)),
		}}},
	}}
	gs.Electors = []Elector{{ProposalId: 1, Operator: operator, RegionId: "cn"}}
	gs.Votes = []Vote{{ProposalId: 1, Voter: sample.AccAddress(), Option: VoteOptionYes}}
	require.NoError(t, gs.Validate())

	gs.ProposalSeq = 0
	require.Error(t, gs.Validate())
}

func TestMsgSubmitProposalValidateBasic(t *testing.T) {
	proposer := sample.AccAddress()
	receiver := sample.AccAddress()
	msg := &MsgSubmitProposal{
		Proposer:     proposer,
		Kind:         ProposalKindTreasuryWithdraw,
		Title:        "withdraw",
		VotingPeriod: 1,
		Treasury: &TreasuryWithdrawPayload{Withdrawals: []Withdrawal{{
			RegionId: "cn",
			Receiver: receiver,
			Amount:   sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
		}}},
	}
	require.NoError(t, msg.ValidateBasic())
	require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(proposer)}, msg.GetSigners())
	require.NotEmpty(t, msg.GetSignBytes())

	msg.VotingPeriod = 0
	require.Error(t, msg.ValidateBasic())
}

func TestMsgVoteValidateBasic(t *testing.T) {
	account := sdk.MustAccAddressFromBech32(sample.AccAddress())
	msg := &MsgVote{Voter: account.String(), ProposalId: 1, Option: VoteOptionYes}
	require.NoError(t, msg.ValidateBasic())
	require.Equal(t, []sdk.AccAddress{account}, msg.GetSigners())
	require.Contains(t, string(msg.GetSignBytes()), account.String())

	// Compatibility with old Elector snapshots must not allow validator
	// addresses in transactions: the voter and signer are ordinary accounts.
	msg.Voter = sdk.ValAddress(account).String()
	require.Error(t, msg.ValidateBasic())
	msg.Voter = account.String()
	msg.ProposalId = 0
	require.Error(t, msg.ValidateBasic())
	msg.ProposalId = 1
	msg.Option = VoteOptionUnspecified
	require.Error(t, msg.ValidateBasic())
}
