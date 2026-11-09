package types

import (
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
)

func TestElectorAccountAddress(t *testing.T) {
	account := sdk.MustAccAddressFromBech32(sample.AccAddress())
	for _, address := range []string{account.String(), strings.ToUpper(account.String()), sdk.ValAddress(account).String()} {
		actual, err := (Elector{Operator: address}).AccountAddress()
		require.NoError(t, err)
		require.True(t, account.Equals(actual))
	}
	_, err := (Elector{Operator: "invalid"}).AccountAddress()
	require.Error(t, err)
}

func TestGenesisEquityElectorRules(t *testing.T) {
	account := sample.AccAddress()
	operator := sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String()
	now := time.Unix(1_700_000_000, 0).UTC()
	newGenesis := func() *GenesisState {
		gs := DefaultGenesis()
		gs.ProposalSeq = 1
		gs.Proposals = []Proposal{{
			Id: 1, Kind: ProposalKindTreasuryWithdraw, Proposer: account,
			Title: "equity withdrawal", Status: ProposalStatusVoting,
			SubmitTime: now, VotingEndTime: now.Add(time.Hour), VotingPeriod: 3600,
			Electorate: ElectorateKindRegionEquity, RequiredOperator: operator,
			Treasury: &TreasuryWithdrawPayload{Withdrawals: []Withdrawal{{
				RegionId: "chn", Receiver: account, Amount: sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
			}}},
		}}
		gs.Electors = []Elector{{ProposalId: 1, Operator: account, RegionId: "chn"}}
		return gs
	}
	require.NoError(t, newGenesis().Validate())
	for _, tc := range []struct {
		name   string
		modify func(*GenesisState)
	}{
		{"missing operator", func(gs *GenesisState) { gs.Proposals[0].RequiredOperator = "" }},
		{"operator is not valoper", func(gs *GenesisState) { gs.Proposals[0].RequiredOperator = account }},
		{"equity is not account", func(gs *GenesisState) { gs.Electors[0].Operator = operator }},
		{"foreign region", func(gs *GenesisState) { gs.Electors[0].RegionId = "usa" }},
		{"operator in counted electorate", func(gs *GenesisState) {
			val, err := sdk.ValAddressFromBech32(operator)
			require.NoError(t, err)
			gs.Electors[0].Operator = sdk.AccAddress(val).String()
		}},
		{"duplicate canonical elector", func(gs *GenesisState) {
			duplicate := gs.Electors[0]
			duplicate.Operator = strings.ToUpper(duplicate.Operator)
			gs.Electors = append(gs.Electors, duplicate)
		}},
		{"multiple withdrawal regions", func(gs *GenesisState) {
			withdrawal := gs.Proposals[0].Treasury.Withdrawals[0]
			withdrawal.RegionId = "usa"
			gs.Proposals[0].Treasury.Withdrawals = append(gs.Proposals[0].Treasury.Withdrawals, withdrawal)
		}},
		{"non treasury electorate", func(gs *GenesisState) {
			gs.Proposals[0].Kind = ProposalKindUpdateMinVotingPeriod
			gs.Proposals[0].Treasury = nil
			gs.Proposals[0].UpdateMinVotingPeriod = &UpdateMinVotingPeriodPayload{MinVotingPeriod: 120}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gs := newGenesis()
			tc.modify(gs)
			require.Error(t, gs.Validate())
		})
	}
}
