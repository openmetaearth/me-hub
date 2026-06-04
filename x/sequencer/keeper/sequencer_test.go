package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/openmetaearth/me-hub/testutil/nullify"
	"github.com/openmetaearth/me-hub/x/sequencer/keeper"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNSequencer(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Sequencer {
	items := make([]types.Sequencer, n)
	for i := range items {
		seq := types.Sequencer{
			SequencerAddress: strconv.Itoa(i),
			Status:           types.Bonded,
		}
		items[i] = seq

		keeper.SetSequencer(ctx, items[i])
	}
	return items
}

func TestSequencerGet(t *testing.T) {
	keeper, ctx := keepertest.SequencerKeeper(t)
	items := createNSequencer(keeper, ctx, 10)
	for _, item := range items {
		item := item
		rst, found := keeper.GetSequencer(ctx,
			item.SequencerAddress,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestSequencerGetAll(t *testing.T) {
	keeper, ctx := keepertest.SequencerKeeper(t)
	items := createNSequencer(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllSequencers(ctx)),
	)
}

func TestSequencersByRollappGet(t *testing.T) {
	keeper, ctx := keepertest.SequencerKeeper(t)
	items := createNSequencer(keeper, ctx, 10)
	rst := keeper.GetSequencersByRollapp(ctx,
		items[0].RollappId,
	)

	require.Equal(t, len(rst), len(items))
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(rst),
	)
}

func (suite *SequencerTestSuite) TestRotatingSequencerByBond() {
	suite.SetupTest()
	rollappId := suite.CreateDefaultRollapp()
	params := suite.App.SequencerKeeper.GetParams(suite.Ctx)
	params.MinBond = sdk.NewCoin(bond.Denom, sdk.NewInt(100))
	suite.App.SequencerKeeper.SetParams(suite.Ctx, params)

	numOfSequencers := 5

	// create sequencers
	seqAddrs := make([]string, numOfSequencers)
	for j := 0; j < len(seqAddrs)-1; j++ {
		seqAddrs[j] = suite.CreateSequencerWithBond(suite.Ctx, rollappId, params.MinBond)
	}
	// last one with high bond is the expected new proposer
	seqAddrs[len(seqAddrs)-1] = suite.CreateSequencerWithBond(
		suite.Ctx,
		rollappId,
		sdk.NewCoin(params.MinBond.Denom, params.MinBond.Amount.MulRaw(2)),
	)
	expectedProposer := seqAddrs[len(seqAddrs)-1]

	// check starting proposer and unbond
	sequencer, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, seqAddrs[0])
	suite.Require().True(found)
	suite.Require().True(sequencer.Proposer)

	suite.App.SequencerKeeper.RotateProposer(suite.Ctx, rollappId)

	// check proposer rotation
	newProposer, _ := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, expectedProposer)
	suite.Equal(types.Bonded, newProposer.Status)
	suite.True(newProposer.Proposer)
}

func (suite *SequencerTestSuite) TestRotateProposerUsesMinBondDenom() {
	suite.SetupTest()
	rollappId := suite.CreateDefaultRollapp()

	bondDenom := bond.Denom
	otherDenom := "otherdenom"
	if otherDenom == bondDenom {
		otherDenom = "otherdenom2"
	}

	lowMinBondWithLargeOtherDenom := types.Sequencer{
		SequencerAddress: "a-low-min-bond",
		RollappId:        rollappId,
		Status:           types.Bonded,
		Tokens: sdk.NewCoins(
			sdk.NewCoin(bondDenom, sdk.NewInt(90)),
			sdk.NewCoin(otherDenom, sdk.NewInt(1000)),
		),
	}
	expectedProposer := types.Sequencer{
		SequencerAddress: "b-high-min-bond",
		RollappId:        rollappId,
		Status:           types.Bonded,
		Tokens: sdk.NewCoins(
			sdk.NewCoin(bondDenom, sdk.NewInt(100)),
			sdk.NewCoin(otherDenom, sdk.NewInt(1)),
		),
	}
	lowerBond := types.Sequencer{
		SequencerAddress: "c-lower-bond",
		RollappId:        rollappId,
		Status:           types.Bonded,
		Tokens:           sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(80))),
	}

	suite.App.SequencerKeeper.SetSequencer(suite.Ctx, lowMinBondWithLargeOtherDenom)
	suite.App.SequencerKeeper.SetSequencer(suite.Ctx, expectedProposer)
	suite.App.SequencerKeeper.SetSequencer(suite.Ctx, lowerBond)

	suite.App.SequencerKeeper.RotateProposer(suite.Ctx, rollappId)

	selected, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, expectedProposer.SequencerAddress)
	suite.Require().True(found)
	suite.True(selected.Proposer)

	lowBond, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, lowMinBondWithLargeOtherDenom.SequencerAddress)
	suite.Require().True(found)
	suite.False(lowBond.Proposer)
}
