package regiongov_test

import (
	"testing"
	"time"

	cmtbftdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	"github.com/openmetaearth/me-hub/x/regiongov/keeper"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func TestInitExportGenesis(t *testing.T) {
	proposer := sample.AccAddress()
	receiver := sample.AccAddress()
	operatorAcc := sample.AccAddress()
	operator := sdk.ValAddress(sdk.MustAccAddressFromBech32(operatorAcc)).String()
	equity := sample.AccAddress()
	now := time.Unix(1_700_000_000, 0).UTC()

	genesis := types.GenesisState{
		Params:           types.DefaultParams(),
		ProposalSeq:      1,
		PendingActionSeq: 1,
		Proposals: []types.Proposal{{
			Id:            1,
			Kind:          types.ProposalKindTreasuryWithdraw,
			Proposer:      proposer,
			Title:         "global treasury",
			Summary:       "public",
			Status:        types.ProposalStatusVoting,
			SubmitTime:    now,
			VotingEndTime: now.Add(2 * time.Hour),
			VotingPeriod:  7200,
			Electorate:    types.ElectorateKindCountryRegionOperators,
			Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{{
				RegionId: "me_earth",
				Receiver: receiver,
				Amount:   sdk.NewCoins(sdk.NewInt64Coin("umec", 100)),
			}}},
		}},
		Electors: []types.Elector{{ProposalId: 1, Operator: operator, RegionId: "cn"}},
		Votes:    []types.Vote{{ProposalId: 1, Voter: operatorAcc, Option: types.VoteOptionYes}},
		Equities: []types.RegionEquity{{RegionId: "cn", Addresses: []string{equity}}},
		RegionFreezes: []types.RegionFreeze{{
			RegionId: "us",
			Frozen:   true,
		}},
		EquityFreezes: []types.EquityFreeze{{
			RegionId: "cn",
			Address:  equity,
			Frozen:   true,
		}},
		PendingActions: []types.PendingAction{{
			Id:         1,
			Kind:       types.PendingActionKindSetEquity,
			Status:     types.PendingActionStatusPending,
			RegionId:   "cn",
			Proposer:   operatorAcc,
			Addresses:  []string{equity},
			SubmitTime: now,
		}},
	}
	require.NoError(t, genesis.Validate())

	k, ctx := newKeeper(t)
	regiongov.InitGenesis(ctx, k, genesis)
	exported := regiongov.ExportGenesis(ctx, k)

	require.Equal(t, genesis.Params, exported.Params)
	require.Equal(t, genesis.ProposalSeq, exported.ProposalSeq)
	require.Equal(t, genesis.PendingActionSeq, exported.PendingActionSeq)
	require.Equal(t, genesis.Proposals, exported.Proposals)
	require.Equal(t, genesis.Electors, exported.Electors)
	require.Equal(t, genesis.Votes, exported.Votes)
	require.Equal(t, genesis.Equities, exported.Equities)
	require.Equal(t, genesis.RegionFreezes, exported.RegionFreezes)
	require.Equal(t, genesis.EquityFreezes, exported.EquityFreezes)
	require.Equal(t, genesis.PendingActions, exported.PendingActions)
}

func TestQueryParams(t *testing.T) {
	k, ctx := newKeeper(t)
	regiongov.InitGenesis(ctx, k, *types.DefaultGenesis())

	res, err := k.Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, uint64(3600), res.Params.MinVotingPeriod)
	require.Equal(t, uint32(51), res.Params.PassThresholdPercent)
	require.Equal(t, uint32(1), res.Params.EquityWithdrawPercent)
}

func newKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	t.Helper()
	return newKeeperWith(t, nil, nil)
}

func newKeeperWith(t *testing.T, staking types.StakingKeeper, dao types.DaoKeeper) (keeper.Keeper, sdk.Context) {
	t.Helper()
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	db := cmtbftdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	k := keeper.NewKeeper(cdc, storeKey, staking, dao)
	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	return k, ctx
}
