package keeper_test

import (
	"crypto/ecdsa"
	"fmt"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/st-chain/me-hub/app/apptesting"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/testutil/helpers"
	bsctypes "github.com/st-chain/me-hub/x/bsc/types"
	"github.com/st-chain/me-hub/x/gravity/keeper"
	"github.com/st-chain/me-hub/x/gravity/types"
	trontypes "github.com/st-chain/me-hub/x/tron/types"
	minttypes "github.com/st-chain/me-hub/x/wmint/types"
	wstakingkeeper "github.com/st-chain/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer           types.MsgServer
	queryClient         types.QueryClient
	meEarthValidator    stakingtypes.Validator
	experienceValidator stakingtypes.Validator
	usaValidator        stakingtypes.Validator

	relayerAddrs  []sdk.AccAddress
	relayerNumber int
	externalPris  []*ecdsa.PrivateKey
	chainName     string
}

func TestGravityKeeperTestSuite(t *testing.T) {
	subModules := []string{
		bsctypes.ModuleName,
		//trontypes.ModuleName,
	}
	for _, moduleName := range subModules {
		suite.Run(t, &KeeperTestSuite{
			chainName: moduleName,
		})
	}
}

func (s *KeeperTestSuite) MsgServer() types.MsgServer {
	//if suite.chainName == trontypes.ModuleName {
	//	return tronkeeper.NewMsgServerImpl(suite.app.TronKeeper)
	//}
	return keeper.NewMsgServerImpl(s.Keeper())
}

func (s *KeeperTestSuite) QueryClient() types.QueryClient {
	queryHelper := baseapp.NewQueryServerTestHelper(s.Ctx, s.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(s.Keeper()))
	return types.NewQueryClient(queryHelper)
}

func (s *KeeperTestSuite) Keeper() keeper.Keeper {
	switch s.chainName {
	case bsctypes.ModuleName:
		return s.App.BscKeeper
	//case trontypes.ModuleName:
	//	return s.App.TronKeeper.Keeper
	default:
		panic(fmt.Sprintf("invalid chain name:%s", s.chainName))
	}
}

func (s *KeeperTestSuite) SetupTest() {

	app := apptesting.Setup(s.T(), false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	err := app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	s.Require().NoError(err)

	err = app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	s.Require().NoError(err)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	err = app.StakingKeeper.SetParams(ctx, stakingParams)
	s.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	queryClient := types.NewQueryClient(queryHelper)

	s.App = app
	s.Ctx = ctx
	s.queryClient = queryClient

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	stakingMsgServer := wstakingkeeper.NewMsgServerImpl(app.StakingKeeper, app.TransferKeeper, stakingKeeperMsgSrv)

	s.InitializeDao()

	validators := s.App.StakingKeeper.GetValidators(s.Ctx, 10)
	s.Require().True(len(validators) >= 3)
	s.meEarthValidator = validators[0]
	s.experienceValidator = validators[1]
	s.usaValidator = validators[2]

	newRegion := wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            wstakingtypes.ExperienceRegionName,
		OperatorAddress: s.experienceValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newRegion = wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newRegion = wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            wstakingtypes.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err = stakingMsgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	s.relayerNumber = 10
	s.relayerAddrs = s.NewAccounts(s.relayerNumber)
	s.Require().EqualValues(s.relayerNumber, len(s.relayerAddrs))
	s.externalPris = helpers.CreateMultiECDSA(s.relayerNumber)

	proposalRelayer := &types.ProposalRelayer{}
	for i := 0; i < s.relayerNumber; i++ {
		err = s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, s.relayerAddrs[i], sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 10000000000)})
		s.Require().NoError(err)
		proposalRelayer.Relayers = append(proposalRelayer.Relayers, s.relayerAddrs[i].String())
	}
	s.Keeper().SetProposalRelayer(s.Ctx, proposalRelayer)
}

func (s *KeeperTestSuite) SetupSubTest() {
	s.SetupTest()
}

func (s *KeeperTestSuite) SignRelayerSetConfirm(external *ecdsa.PrivateKey, relayerSet *types.RelayerSet) (string, []byte) {
	externalAddress := crypto.PubkeyToAddress(external.PublicKey).String()
	gravityId := s.Keeper().GetGravityID(s.Ctx)
	checkpoint, err := relayerSet.GetCheckpoint(gravityId)
	s.NoError(err)
	signature, err := types.NewEthereumSignature(checkpoint, external)
	s.NoError(err)
	if trontypes.ModuleName == s.chainName {
		//externalAddress = tronaddress.PubkeyToAddress(external.PublicKey).String()
		//
		//checkpoint, err = trontypes.GetCheckpointRelayerSet(relayerSet, gravityId)
		//s.Require().NoError(err)
		//
		//signature, err = trontypes.NewTronSignature(checkpoint, external)
		//s.Require().NoError(err)
	}
	return externalAddress, signature
}

func (s *KeeperTestSuite) SendClaim(externalClaim types.ExternalClaim) {
	var err error
	switch claim := externalClaim.(type) {
	case *types.MsgSendToMeClaim:
		_, err = s.MsgServer().SendToMeClaim(s.Ctx, claim)
		s.NoError(err)
		s.Require().NoError(err)
	}
}

func (s *KeeperTestSuite) PubKeyToExternalAddr(publicKey ecdsa.PublicKey) string {
	address := crypto.PubkeyToAddress(publicKey)
	return types.ExternalAddrToStr(s.chainName, address.Bytes())
}
