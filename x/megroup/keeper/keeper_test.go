package keeper_test

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/megroup/keeper"
	"github.com/openmetaearth/me-hub/x/megroup/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer types.MsgServer
	testAccs  []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	s.App = app
	s.Ctx = ctx
	s.msgServer = keeper.NewMsgServerImpl(*app.GroupKeeper)
	s.testAccs = s.NewAccounts(2)
}

func (s *KeeperTestSuite) setActiveKyc(address sdk.AccAddress, did, regionID string) {
	s.App.KycKeeper.SetDID(s.Ctx, address, did)
	s.App.KycKeeper.SetDidInfo(s.Ctx, did, didtypes.DidInfo{
		Did:      did,
		Address:  address.String(),
		RegionId: regionID,
		KycLevel: didtypes.KYC_LEVEL_TWO,
		Status:   didtypes.DID_STATUS_ACTIVE,
	})
}

func (s *KeeperTestSuite) TestJoinGroupRejectsEmptyAdminBeforeMutatingState() {
	s.SetupTest()

	const (
		groupID  uint64 = 1
		regionID        = "test-region"
	)
	applicant := s.testAccs[0]
	treasure := s.testAccs[1]
	s.setActiveKyc(applicant, "did-join-001", regionID)

	s.App.StakingKeeper.SetRegion(s.Ctx, wstakingtypes.Region{
		RegionId:           regionID,
		RegionTreasureAddr: treasure.String(),
		RegionShare:        sdk.ZeroInt(),
		DelegateInterest:   sdk.ZeroDec(),
		DelegateAmount:     sdk.ZeroInt(),
		FixedDepositAmount: sdk.ZeroInt(),
	})
	s.App.GroupKeeper.SetGroupInfo(s.Ctx, types.GroupInfo{
		Id:          groupID,
		Admin:       "",
		TotalWeight: sdk.ZeroInt().String(),
		RegionID:    regionID,
	})
	s.App.GroupKeeper.SetGroupToRegion(s.Ctx, regionID, groupID)
	s.App.GroupKeeper.SetGroupMemberCount(s.Ctx, groupID, 0)

	_, err := s.msgServer.JoinGroup(s.Ctx, &types.MsgJoinGroup{
		Creator:          applicant.String(),
		GroupId:          groupID,
		ApplicantAddress: applicant.String(),
	})
	s.Require().ErrorIs(err, types.ErrProcData)
	s.Require().ErrorContains(err, "group admin address is invalid")

	_, joined := s.App.GroupKeeper.GetMemberJoined(s.Ctx, applicant.String())
	s.Require().False(joined)

	groupMemberCount, found := s.App.GroupKeeper.GetGroupMemberCount(s.Ctx, groupID)
	s.Require().True(found)
	s.Require().Equal(uint64(0), groupMemberCount)
}
