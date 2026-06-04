package keeper_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app/apptesting"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	kyckeeper "github.com/openmetaearth/me-hub/x/kyc/keeper"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/megroup/keeper"
	"github.com/openmetaearth/me-hub/x/megroup/types"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wminttypes "github.com/openmetaearth/me-hub/x/wmint/types"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	groupMsgServer types.MsgServer
	kycMsgServer   kyctypes.MsgServer
	queryClient    types.QueryClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	ctx := app.GetBaseApp().NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, *app.GroupKeeper)

	s.App = app
	s.Ctx = ctx
	s.groupMsgServer = keeper.NewMsgServerImpl(*app.GroupKeeper)
	s.kycMsgServer = kyckeeper.NewMsgServerImpl(*app.KycKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)

	s.InitializeDao()

	stakingMsgServer := wstakingkeeper.NewMsgServerImpl(
		app.StakingKeeper,
		app.TransferKeeper,
		stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper),
	)
	validators := app.StakingKeeper.GetValidators(ctx, 10)
	s.Require().GreaterOrEqual(len(validators), 3)
	_, err := stakingMsgServer.NewRegion(ctx, &wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            wstakingtypes.ExperienceRegionName,
		OperatorAddress: validators[1].OperatorAddress,
	})
	s.Require().NoError(err)
	_, err = stakingMsgServer.NewRegion(ctx, &wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: validators[2].OperatorAddress,
	})
	s.Require().NoError(err)
	_, err = stakingMsgServer.NewRegion(ctx, &wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            wstakingtypes.MeEarthRegionName,
		OperatorAddress: validators[0].OperatorAddress,
	})
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestKycLevelDowngradeRemovesJoinedMember() {
	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).
		WithBlockHeight(wminttypes.OneDayTotalBlocks).
		WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	member, pubKey := s.NewAccount()
	const did = "did-kyc-downgrade"
	regionID := wstakingtypes.MeEarthRegionId

	_, err := s.kycMsgServer.Approve(s.Ctx, &kyctypes.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: regionID,
		Address:  member.String(),
		Pubkey:   pubKey,
		Uri:      "http://127.0.0.1/kyc",
		Hash:     "hash-before-downgrade",
		Level:    didtypes.KYC_LEVEL_TWO,
	})
	s.Require().NoError(err)

	groupID, found := s.App.GroupKeeper.GetGroupIdByRegion(s.Ctx, regionID)
	s.Require().True(found)
	_, err = s.groupMsgServer.JoinGroup(s.Ctx, &types.MsgJoinGroup{
		Creator:          member.String(),
		GroupId:          groupID,
		ApplicantAddress: member.String(),
	})
	s.Require().NoError(err)

	_, err = s.kycMsgServer.Update(s.Ctx, &kyctypes.MsgUpdate{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: regionID,
		Uri:      "http://127.0.0.1/kyc",
		Hash:     "hash-after-downgrade",
		Level:    didtypes.KYC_LEVEL_ONE,
	})
	s.Require().NoError(err)

	joined, found := s.App.GroupKeeper.GetMemberJoined(s.Ctx, member.String())
	s.Require().True(found)
	s.Require().Zero(joined.GroupId)

	count, found := s.App.GroupKeeper.GetGroupMemberCount(s.Ctx, groupID)
	s.Require().True(found)
	s.Require().Zero(count)

	_, err = s.App.GroupKeeper.GroupMember(s.Ctx, &types.QueryGetGroupMemberRequest{
		Address: member.String(),
	})
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestKycRegionUpdateDoesNotMigrateDowngradedMember() {
	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).
		WithBlockHeight(wminttypes.OneDayTotalBlocks).
		WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	member, pubKey := s.NewAccount()
	const did = "did-kyc-region-downgrade"
	regionID := wstakingtypes.MeEarthRegionId
	newRegionID := "usa"

	_, err := s.kycMsgServer.Approve(s.Ctx, &kyctypes.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: regionID,
		Address:  member.String(),
		Pubkey:   pubKey,
		Uri:      "http://127.0.0.1/kyc",
		Hash:     "hash-before-region-downgrade",
		Level:    didtypes.KYC_LEVEL_TWO,
	})
	s.Require().NoError(err)

	oldGroupID, found := s.App.GroupKeeper.GetGroupIdByRegion(s.Ctx, regionID)
	s.Require().True(found)
	newGroupID, found := s.App.GroupKeeper.GetGroupIdByRegion(s.Ctx, newRegionID)
	s.Require().True(found)

	_, err = s.groupMsgServer.JoinGroup(s.Ctx, &types.MsgJoinGroup{
		Creator:          member.String(),
		GroupId:          oldGroupID,
		ApplicantAddress: member.String(),
	})
	s.Require().NoError(err)

	_, err = s.kycMsgServer.Update(s.Ctx, &kyctypes.MsgUpdate{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: newRegionID,
		Uri:      "http://127.0.0.1/kyc",
		Hash:     "hash-after-region-downgrade",
		Level:    didtypes.KYC_LEVEL_ONE,
	})
	s.Require().NoError(err)

	joined, found := s.App.GroupKeeper.GetMemberJoined(s.Ctx, member.String())
	s.Require().True(found)
	s.Require().Zero(joined.GroupId)

	oldCount, found := s.App.GroupKeeper.GetGroupMemberCount(s.Ctx, oldGroupID)
	s.Require().True(found)
	s.Require().Zero(oldCount)

	newCount, found := s.App.GroupKeeper.GetGroupMemberCount(s.Ctx, newGroupID)
	s.Require().True(found)
	s.Require().Zero(newCount)
}
