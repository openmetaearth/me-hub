package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func TestTxCommands(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	txConfig := tx.NewTxConfig(cdc, tx.DefaultSignModes)
	signer := sdk.AccAddress([]byte("signer-address-12345")).String()
	target := sdk.AccAddress([]byte("target-address-12345")).String()
	operator := sdk.ValAddress([]byte("operator-address-123")).String()
	proposalFile := func(name, body string) string {
		path := filepath.Join(t.TempDir(), name)
		require.NoError(t, os.WriteFile(path, []byte(body), 0600))
		return path
	}
	treasury := proposalFile("treasury.json", fmt.Sprintf(`{"kind":"PROPOSAL_KIND_TREASURY_WITHDRAW","title":"Withdraw","summary":"Budget","voting_period":"3600","treasury":{"withdrawals":[{"region_id":"region-1","receiver":%q,"amount":[{"denom":"umec","amount":"100"}]}]}}`, target))
	replacement := proposalFile("operator.json", fmt.Sprintf(`{"proposer":%q,"kind":"PROPOSAL_KIND_REPLACE_REGION_OPERATOR","title":"Replace","voting_period":"3600","replace_operator":{"region_id":"region-1","new_operator":%q}}`, target, operator))
	invalidJSON := proposalFile("invalid.json", `{`)
	invalidPayload := proposalFile("invalid-payload.json", `{"kind":"PROPOSAL_KIND_TREASURY_WITHDRAW","title":"Withdraw","voting_period":"3600"}`)

	cases := []struct {
		name string
		args []string
		want sdk.Msg
	}{
		{"treasury", []string{"submit-proposal", treasury}, &types.MsgSubmitProposal{Proposer: signer, Kind: types.ProposalKindTreasuryWithdraw, Title: "Withdraw", Summary: "Budget", VotingPeriod: 3600, Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{{RegionId: "region-1", Receiver: target, Amount: sdk.NewCoins(sdk.NewInt64Coin("umec", 100))}}}}},
		{"operator", []string{"submit-proposal", replacement}, &types.MsgSubmitProposal{Proposer: signer, Kind: types.ProposalKindReplaceRegionOperator, Title: "Replace", VotingPeriod: 3600, ReplaceOperator: &types.ReplaceRegionOperatorPayload{RegionId: "region-1", NewOperator: operator}}},
		{"yes", []string{"vote", "1", "yes"}, &types.MsgVote{Voter: signer, ProposalId: 1, Option: types.VoteOptionYes}},
		{"no", []string{"vote", "2", "NO"}, &types.MsgVote{Voter: signer, ProposalId: 2, Option: types.VoteOptionNo}},
		{"update", []string{"submit-equity-update", "region-1", target, signer}, &types.MsgSubmitEquityUpdate{Proposer: signer, RegionId: "region-1", Addresses: []string{target, signer}}},
		{"clear", []string{"submit-equity-update", "region-1"}, &types.MsgSubmitEquityUpdate{Proposer: signer, RegionId: "region-1"}},
		{"submit-freeze", []string{"submit-freeze-equity", "region-1", target}, &types.MsgSubmitFreezeEquity{Proposer: signer, RegionId: "region-1", Addresses: []string{target}}},
		{"submit-unfreeze", []string{"submit-unfreeze-equity", "region-1", target}, &types.MsgSubmitUnfreezeEquity{Proposer: signer, RegionId: "region-1", Addresses: []string{target}}},
		{"approve", []string{"approve-pending-action", "3"}, &types.MsgApprovePendingAction{Authority: signer, ActionId: 3}},
		{"reject", []string{"reject-pending-action", "4"}, &types.MsgRejectPendingAction{Authority: signer, ActionId: 4}},
		{"freeze-region", []string{"freeze-region", "region-1"}, &types.MsgFreezeRegion{Authority: signer, RegionIds: []string{"region-1"}}},
		{"unfreeze-region", []string{"unfreeze-region", "region-1"}, &types.MsgUnfreezeRegion{Authority: signer, RegionIds: []string{"region-1"}}},
		{"freeze-regions", []string{"freeze-region", "region-1", "region-2"}, &types.MsgFreezeRegion{Authority: signer, RegionIds: []string{"region-1", "region-2"}}},
		{"unfreeze-regions", []string{"unfreeze-region", "region-1", "region-2"}, &types.MsgUnfreezeRegion{Authority: signer, RegionIds: []string{"region-1", "region-2"}}},
		{"freeze-equity", []string{"freeze-equity", "region-1", target}, &types.MsgFreezeEquity{Authority: signer, RegionId: "region-1", Address: target}},
		{"unfreeze-equity", []string{"unfreeze-equity", "region-1", target}, &types.MsgUnfreezeEquity{Authority: signer, RegionId: "region-1", Address: target}},
		{"submit-freeze-many", []string{"submit-freeze-equity", "region-1", target, signer}, &types.MsgSubmitFreezeEquity{Proposer: signer, RegionId: "region-1", Addresses: []string{target, signer}}},
		{"submit-unfreeze-many", []string{"submit-unfreeze-equity", "region-1", target, signer}, &types.MsgSubmitUnfreezeEquity{Proposer: signer, RegionId: "region-1", Addresses: []string{target, signer}}},
		{"submit-freeze-empty", []string{"submit-freeze-equity", "region-1"}, nil},
		{"submit-unfreeze-empty", []string{"submit-unfreeze-equity", "region-1"}, nil},
		{"submit-freeze-duplicate", []string{"submit-freeze-equity", "region-1", target, target}, nil},
		{"submit-unfreeze-invalid", []string{"submit-unfreeze-equity", "region-1", target, "invalid"}, nil},
		{"zero-id", []string{"vote", "0", "yes"}, nil},
		{"negative-id", []string{"approve-pending-action", "-1"}, nil},
		{"overflow-id", []string{"reject-pending-action", "18446744073709551616"}, nil},
		{"invalid-id", []string{"vote", "abc", "yes"}, nil},
		{"invalid-vote", []string{"vote", "1", "abstain"}, nil},
		{"invalid-address", []string{"freeze-equity", "region-1", "invalid"}, nil},
		{"duplicate-address", []string{"submit-equity-update", "region-1", target, target}, nil},
		{"empty-region", []string{"freeze-region", ""}, nil},
		{"missing-args", []string{"vote", "1"}, nil},
		{"missing-regions", []string{"unfreeze-region"}, nil},
		{"duplicate-regions", []string{"freeze-region", "usa", "USA"}, nil},
		{"missing-file", []string{"submit-proposal", filepath.Join(t.TempDir(), "missing.json")}, nil},
		{"invalid-json", []string{"submit-proposal", invalidJSON}, nil},
		{"invalid-payload", []string{"submit-proposal", invalidPayload}, nil},
	}

	require.Len(t, GetTxCmd().Commands(), 11)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := GetTxCmd()
			cmd, args, err := root.Find(tc.args)
			require.NoError(t, err)
			ctx := client.Context{}.WithCodec(cdc).WithInterfaceRegistry(registry).WithTxConfig(txConfig).WithGenerateOnly(true)
			cmd.SetContext(context.WithValue(context.Background(), client.ClientContextKey, &ctx))
			require.NoError(t, cmd.Flags().Set(flags.FlagFrom, signer))
			require.NoError(t, cmd.Flags().Set(flags.FlagGenerateOnly, "true"))
			require.NoError(t, cmd.Flags().Set(flags.FlagChainID, "regiongov-test"))
			var output bytes.Buffer
			ctx = ctx.WithOutput(&output)
			err = cmd.ValidateArgs(args)
			if err == nil {
				err = cmd.RunE(cmd, args)
			}
			if tc.want == nil {
				require.Error(t, err)
				require.Empty(t, output.Bytes())
				return
			}
			require.NoError(t, err)
			generated, err := txConfig.TxJSONDecoder()(output.Bytes())
			require.NoError(t, err)
			require.Len(t, generated.GetMsgs(), 1)
			require.Equal(t, tc.want, generated.GetMsgs()[0])
		})
	}
}
