package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"strings"
)

func (s *KeeperTestSuite) TestProtocol() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	query := &types.QueryProtocol{}
	res, err := s.queryClient.Protocol(s.Ctx, query)
	s.Require().NoError(err)

	genesis := types.DefaultGenesis()

	s.Require().Equal(res.Protocol.Service.Sid, types.ModuleName)
	s.Require().Equal(res.Protocol.Service.Name, types.ModuleName)
	s.Require().Equal(res.Protocol.Service.Description, "The KYC verifiable credential issuer based The DID(Decentralized Identity).")
	s.Require().Equal(len(res.Protocol.Service.Issuers), len(genesis.Issuers))
	s.Require().Equal(res.Protocol.Service.Status, didtypes.SERVICE_STATUS_ACTIVE)
	s.Require().Equal(len(res.Protocol.Regions), 0)
}

func (s *KeeperTestSuite) TestDID() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	did := "1111111111111111"
	kycAccount, newUserPubkey := s.NewAccount()
	inviter, _ := s.NewAccount()
	msg := &types.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
		Address:  kycAccount.String(),
		Pubkey:   newUserPubkey,
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
		Inviter:  inviter.String(),
	}
	_, err := s.msgServer.Approve(s.Ctx, msg)
	s.Require().NoError(err)

	query := &types.QueryDID{
		Address: kycAccount.String(),
	}
	res, err := s.queryClient.DID(s.Ctx, query)
	s.Require().NoError(err)

	s.Require().Equal(res.Info.Did, did)
	s.Require().Equal(res.Info.Address, kycAccount.String())
	s.Require().Equal(res.Info.Pubkey, newUserPubkey)
	s.Require().Equal(res.Info.Status, didtypes.DID_STATUS_ACTIVE)
}

func (s *KeeperTestSuite) TestDIDs() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	did := "1111111111111111"
	kycAccount, newUserPubkey := s.NewAccount()
	inviter, _ := s.NewAccount()
	msg := &types.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
		Address:  kycAccount.String(),
		Pubkey:   newUserPubkey,
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
		Inviter:  inviter.String(),
	}
	_, err := s.msgServer.Approve(s.Ctx, msg)
	s.Require().NoError(err)

	query := &types.QueryDIDs{
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
	}
	res, err := s.queryClient.DIDs(s.Ctx, query)
	s.Require().NoError(err)
	s.Require().Equal(len(res.Infos), 1)

	info := res.Infos[0]
	s.Require().Equal(info.Did, did)
	s.Require().Equal(info.Address, kycAccount.String())
	s.Require().Equal(info.Pubkey, newUserPubkey)
	s.Require().Equal(info.Status, didtypes.DID_STATUS_ACTIVE)
}

func (s *KeeperTestSuite) TestKYC() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	did := "1111111111111111"
	kycAccount, newUserPubkey := s.NewAccount()
	inviter, _ := s.NewAccount()
	msg := &types.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
		Address:  kycAccount.String(),
		Pubkey:   newUserPubkey,
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
		Inviter:  inviter.String(),
	}
	_, err := s.msgServer.Approve(s.Ctx, msg)
	s.Require().NoError(err)

	query := &types.QueryKYC{
		Did: did,
	}
	res, err := s.queryClient.KYC(s.Ctx, query)
	s.Require().NoError(err)

	s.Require().Equal(res.Kyc.Did, did)
	s.Require().Equal(res.Kyc.Sid, types.ModuleName)
	s.Require().Equal(res.Kyc.Hash, "aaaa")
	s.Require().Equal(res.Kyc.Uri, "http://127.0.0.1/8001")
	s.Require().Equal(res.Kyc.Data, []byte(strings.ToLower(wstakingtypes.MeEarthRegionName)))
}

func (s *KeeperTestSuite) TestKYCs() {
	s.SetupTest()

	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	did := "1111111111111111"
	kycAccount, newUserPubkey := s.NewAccount()
	inviter, _ := s.NewAccount()
	msg := &types.MsgApprove{
		Issuer:   s.Dao.GlobalDao,
		Did:      did,
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
		Address:  kycAccount.String(),
		Pubkey:   newUserPubkey,
		Uri:      "http://127.0.0.1/8001",
		Hash:     "aaaa",
		Inviter:  inviter.String(),
	}
	_, err := s.msgServer.Approve(s.Ctx, msg)
	s.Require().NoError(err)

	query := &types.QueryKYCs{
		RegionId: strings.ToLower(wstakingtypes.MeEarthRegionName),
	}
	res, err := s.queryClient.KYCs(s.Ctx, query)
	s.Require().NoError(err)
	s.Require().Equal(len(res.KYCs), 1)

	kyc := res.KYCs[0]
	s.Require().Equal(kyc.Did, did)
	s.Require().Equal(kyc.Sid, types.ModuleName)
	s.Require().Equal(kyc.Hash, "aaaa")
	s.Require().Equal(kyc.Uri, "http://127.0.0.1/8001")
	s.Require().Equal(kyc.Data, []byte(strings.ToLower(wstakingtypes.MeEarthRegionName)))
}
