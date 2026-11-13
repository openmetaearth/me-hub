package v2_0_18_rc2_test //nolint:revive

import (
	"errors"
	"testing"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/upgrades/v2_0_18_rc2"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	regiongovtypes "github.com/openmetaearth/me-hub/x/regiongov/types"
)

func TestRegisteredUpgradePreservesExistingState(t *testing.T) {
	app := apptesting.Setup(t, false)
	now := time.Unix(1_700_000_000, 0).UTC()
	ctx := app.NewContext(false, cometbftproto.Header{Height: 10, Time: now})
	require.True(t, app.UpgradeKeeper.HasHandler("v2.0.18.rc2"))
	require.True(t, app.UpgradeKeeper.HasHandler("v2.0.18"))
	require.Empty(t, v2_0_18_rc2.Upgrade.StoreUpgrades.Added)
	require.Empty(t, v2_0_18_rc2.Upgrade.StoreUpgrades.Deleted)
	require.Empty(t, v2_0_18_rc2.Upgrade.StoreUpgrades.Renamed)

	proposer := sample.AccAddress()
	genesis := regiongovtypes.DefaultGenesis()
	// An incremental upgrade must not replace old or customized parameters
	// with the new one-hour initialization default.
	genesis.Params.MinVotingPeriod = 120
	genesis.Params.PassThresholdPercent = 67
	genesis.Params.EquityWithdrawPercent = 5
	genesis.ProposalSeq = 41
	genesis.Proposals = []regiongovtypes.Proposal{{
		Id: 41, Kind: regiongovtypes.ProposalKindTreasuryWithdraw,
		Proposer: proposer, Title: "existing proposal", Status: regiongovtypes.ProposalStatusVoting,
		SubmitTime: now, VotingEndTime: now.Add(time.Minute), VotingPeriod: 60,
		Electorate: regiongovtypes.ElectorateKindCountryRegionOperators, DaoVote: regiongovtypes.DaoVoteYes,
		Treasury: &regiongovtypes.TreasuryWithdrawPayload{Withdrawals: []regiongovtypes.Withdrawal{{
			RegionId: "usa", Receiver: proposer, Amount: sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
		}}},
	}}
	genesis.Electors = []regiongovtypes.Elector{{ProposalId: 41, Operator: proposer, RegionId: "usa"}}
	genesis.Votes = []regiongovtypes.Vote{{ProposalId: 41, Voter: proposer, Option: regiongovtypes.VoteOptionYes}}
	genesis.Equities = []regiongovtypes.RegionEquity{{RegionId: "usa", Addresses: []string{proposer}}}
	genesis.RegionFreezes = []regiongovtypes.RegionFreeze{{RegionId: "usa", Frozen: true}}
	genesis.PendingActionSeq = 7
	genesis.PendingActions = []regiongovtypes.PendingAction{{
		Id: 7, Kind: regiongovtypes.PendingActionKindSetEquity, Status: regiongovtypes.PendingActionStatusPending,
		RegionId: "usa", Proposer: proposer, Addresses: []string{proposer}, SubmitTime: now,
	}}
	regiongov.InitGenesis(ctx, app.RegionGovKeeper, *genesis)
	before := regiongov.ExportGenesis(ctx, app.RegionGovKeeper)
	versionsBefore := app.UpgradeKeeper.GetModuleVersionMap(ctx)
	stakingParamsBefore := app.StakingKeeper.GetParams(ctx)
	validatorsBefore := app.StakingKeeper.GetAllValidators(ctx)
	regionsBefore := app.StakingKeeper.GetAllRegion(ctx)

	plan := upgradetypes.Plan{Name: v2_0_18_rc2.UpgradeName, Height: ctx.BlockHeight()}
	require.NotPanics(t, func() { app.UpgradeKeeper.ApplyUpgrade(ctx, plan) })
	require.Equal(t, ctx.BlockHeight(), app.UpgradeKeeper.GetDoneHeight(ctx, plan.Name))
	require.Equal(t, versionsBefore, app.UpgradeKeeper.GetModuleVersionMap(ctx))
	require.Equal(t, before, regiongov.ExportGenesis(ctx, app.RegionGovKeeper))
	require.Equal(t, stakingParamsBefore, app.StakingKeeper.GetParams(ctx))
	require.Equal(t, validatorsBefore, app.StakingKeeper.GetAllValidators(ctx))
	require.Equal(t, regionsBefore, app.StakingKeeper.GetAllRegion(ctx))
	require.NoError(t, regiongov.ExportGenesis(ctx, app.RegionGovKeeper).Validate())
	queueKey := regiongovtypes.VotingQueueKey(genesis.Proposals[0].VotingEndTime, 41)
	require.Equal(t, regiongovtypes.Uint64ToBytes(41), ctx.KVStore(app.GetKey(regiongovtypes.StoreKey)).Get(queueKey))
}

type migrationTestModule struct {
	regiongov.AppModule
}

func (migrationTestModule) ConsensusVersion() uint64 { return 2 }

func TestUpgradeRunsMigrationsAndPropagatesErrors(t *testing.T) {
	for _, fail := range []bool{false, true} {
		name := "success"
		if fail {
			name = "migration error"
		}
		t.Run(name, func(t *testing.T) {
			mm := module.NewManager(migrationTestModule{})
			cfg := module.NewConfigurator(nil, nil, nil)
			migrationErr := errors.New("migration failed")
			called := false
			require.NoError(t, cfg.RegisterMigration(regiongovtypes.ModuleName, 1, func(sdk.Context) error {
				called = true
				if fail {
					return migrationErr
				}
				return nil
			}))
			handler := v2_0_18_rc2.CreateUpgradeHandler(mm, cfg, nil, nil)
			ctx := sdk.Context{}.WithLogger(log.NewNopLogger())
			versions, err := handler(ctx, upgradetypes.Plan{Name: v2_0_18_rc2.UpgradeName}, module.VersionMap{regiongovtypes.ModuleName: 1})
			require.True(t, called)
			if fail {
				require.ErrorIs(t, err, migrationErr)
				require.Nil(t, versions)
			} else {
				require.NoError(t, err)
				require.Equal(t, module.VersionMap{regiongovtypes.ModuleName: 2}, versions)
			}
		})
	}
}
