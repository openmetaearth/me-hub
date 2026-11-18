package v3_test

import (
	"testing"
	"time"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/app/upgrades"
	"github.com/openmetaearth/me-hub/app/upgrades/v3"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	regiongovtypes "github.com/openmetaearth/me-hub/x/regiongov/types"
)

func TestV3UpgradePreservesV2RegionGovState(t *testing.T) {
	meApp := apptesting.Setup(t)
	ctx := meApp.NewContext(false)
	operator, equity := sample.AccAddress(), sample.AccAddress()
	now := time.Unix(1_700_000_000, 0).UTC()

	// Regiongov already existed at consensus version 1 in v2.0.18, with the
	// same KV schema. Nondefault params and nonzero sequences detect an
	// accidental InitGenesis call for a supposedly new module during v3.
	state := regiongovtypes.GenesisState{
		Params: regiongovtypes.Params{
			MinVotingPeriod:       7200,
			PassThresholdPercent:  67,
			EquityWithdrawPercent: 3,
		},
		ProposalSeq:      42,
		PendingActionSeq: 17,
		Proposals: []regiongovtypes.Proposal{{
			Id:            41,
			Kind:          regiongovtypes.ProposalKindTreasuryWithdraw,
			Proposer:      operator,
			Title:         "existing v2 treasury proposal",
			Status:        regiongovtypes.ProposalStatusVoting,
			SubmitTime:    now,
			VotingEndTime: now.Add(2 * time.Hour),
			VotingPeriod:  7200,
			Electorate:    regiongovtypes.ElectorateKindCountryRegionOperators,
			DaoVote:       regiongovtypes.DaoVoteYes,
			Treasury: &regiongovtypes.TreasuryWithdrawPayload{Withdrawals: []regiongovtypes.Withdrawal{{
				RegionId: "usa",
				Receiver: equity,
				Amount:   sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100)),
			}}},
		}},
		Electors: []regiongovtypes.Elector{{ProposalId: 41, Operator: operator, RegionId: "usa"}},
		Votes:    []regiongovtypes.Vote{{ProposalId: 41, Voter: operator, Option: regiongovtypes.VoteOptionYes}},
		Equities: []regiongovtypes.RegionEquity{{RegionId: "usa", Addresses: []string{equity}}},
		RegionFreezes: []regiongovtypes.RegionFreeze{{
			RegionId: "chn", Frozen: true,
		}},
		EquityFreezes: []regiongovtypes.EquityFreeze{{
			RegionId: "usa", Address: equity, Frozen: true,
		}},
		PendingActions: []regiongovtypes.PendingAction{{
			Id:         13,
			Kind:       regiongovtypes.PendingActionKindSetEquity,
			Status:     regiongovtypes.PendingActionStatusPending,
			RegionId:   "usa",
			Proposer:   operator,
			Addresses:  []string{equity},
			SubmitTime: now,
		}},
	}
	require.NoError(t, state.Validate())
	regiongov.InitGenesis(ctx, meApp.RegionGovKeeper, state)
	store := ctx.KVStore(meApp.GetKey(regiongovtypes.StoreKey))
	require.True(t, store.Has(regiongovtypes.VotingQueueKey(state.Proposals[0].VotingEndTime, 41)))
	before := snapshotRegionGovStore(t, store)
	require.NotEmpty(t, before)

	fromVM, err := meApp.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), fromVM[regiongovtypes.ModuleName], "v2.0.18 regiongov must be an existing module")

	handler := v3.Upgrade.CreateHandler(meApp.ModuleManager(), meApp.Configurator(), &upgrades.UpgradeKeepers{
		AccountKeeper:     &meApp.AccountKeeper,
		GovKeeper:         meApp.GovKeeper,
		RollappKeeper:     meApp.RollappKeeper,
		SequencerKeeper:   meApp.SequencerKeeper,
		ParamsKeeper:      &meApp.ParamsKeeper,
		DelayedAckKeeper:  &meApp.DelayedAckKeeper,
		EIBCKeeper:        &meApp.EIBCKeeper,
		LightClientKeeper: &meApp.LightClientKeeper,
		IBCKeeper:         meApp.IBCKeeper,
		MintKeeper:        &meApp.MintKeeper,
		SlashingKeeper:    &meApp.SlashingKeeper,
		ConsensusKeeper:   &meApp.ConsensusParamsKeeper,
		StakingKeeper:     meApp.StakingKeeper,
	})
	toVM, err := handler(ctx, upgradetypes.Plan{Name: v3.UpgradeName, Height: ctx.BlockHeight()}, fromVM)
	require.NoError(t, err)
	require.Equal(t, uint64(1), fromVM[regiongovtypes.ModuleName])
	require.Equal(t, fromVM[regiongovtypes.ModuleName], toVM[regiongovtypes.ModuleName])
	require.Equal(t, before, snapshotRegionGovStore(t, store), "v3 must preserve every regiongov KV entry, including the voting queue")
	require.Equal(t, &state, regiongov.ExportGenesis(ctx, meApp.RegionGovKeeper))
}

func TestV3StoreUpgradesLeaveExistingRegionGovStoreUntouched(t *testing.T) {
	require.NotContains(t, v3.Upgrade.StoreUpgrades.Added, regiongovtypes.StoreKey)
	require.NotContains(t, v3.Upgrade.StoreUpgrades.Deleted, regiongovtypes.StoreKey)
	for _, renamed := range v3.Upgrade.StoreUpgrades.Renamed {
		require.NotEqual(t, regiongovtypes.StoreKey, renamed.OldKey)
		require.NotEqual(t, regiongovtypes.StoreKey, renamed.NewKey)
	}
}

func snapshotRegionGovStore(t *testing.T, store storetypes.KVStore) map[string][]byte {
	t.Helper()
	state := make(map[string][]byte)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		state[string(iterator.Key())] = append([]byte(nil), iterator.Value()...)
	}
	require.NoError(t, iterator.Close())
	return state
}
