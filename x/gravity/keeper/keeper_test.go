package keeper_test

import (
	"crypto/ecdsa"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	minttypes "github.com/openmetaearth/me-hub/x/wmint/types"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

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
		// trontypes.ModuleName,
	}
	for _, moduleName := range subModules {
		suite.Run(t, &KeeperTestSuite{
			chainName: moduleName,
		})
	}
}

func (s *KeeperTestSuite) MsgServer() types.MsgServer {
	// if suite.chainName == trontypes.ModuleName {
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
	// case trontypes.ModuleName:
	//	return s.App.TronKeeper.Keeper
	default:
		panic(fmt.Sprintf("invalid chain name:%s", s.chainName))
	}
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	s.Ctx = app.NewContext(false).WithBlockHeight(0).WithChainID(apptesting.TestChainID)
	s.App = app

	err := app.BankKeeper.SetParams(s.Ctx, banktypes.DefaultParams())
	s.Require().NoError(err)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	err = app.StakingKeeper.SetParams(s.Ctx, stakingParams)
	s.Require().NoError(err)

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	stakingMsgServer := wstakingkeeper.NewMsgServerImpl(app.StakingKeeper, app.TransferKeeper, stakingKeeperMsgSrv)

	s.InitializeDao()

	validators, err := s.App.StakingKeeper.GetValidators(s.Ctx, 10)
	s.Require().NoError(err)
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
		err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, s.relayerAddrs[i], sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 10000000000)})
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
	return externalAddress, signature
}

func (s *KeeperTestSuite) SendClaim(externalClaim types.ExternalClaim) {
	if claim, ok := externalClaim.(*types.MsgSendToMeClaim); ok {
		_, err := s.MsgServer().SendToMeClaim(s.Ctx, claim)
		s.Require().NoError(err)
	}
}

func (s *KeeperTestSuite) PubKeyToExternalAddr(publicKey ecdsa.PublicKey) string {
	address := crypto.PubkeyToAddress(publicKey)
	return types.ExternalAddrToStr(s.chainName, address.Bytes())
}

// nextEventNonce is the nonce Attest will accept from this relayer.
// New relayers with no personal nonce are treated as caught up to lastObserved.
func (s *KeeperTestSuite) nextEventNonce(relayer sdk.AccAddress) uint64 {
	return s.Keeper().GetLastEventNonceByRelayer(s.Ctx, relayer) + 1
}

func (s *KeeperTestSuite) Commit(block ...int64) { //nolint:stylecheck // Existing suite tests use both conventional receiver names.
	s.Ctx = apptesting.MintBlock(s.App, s.Ctx, block...)
}

func (s *KeeperTestSuite) NewBridgeToken(receiver sdk.AccAddress, initAmount sdk.Coin) (bridgeToken types.BridgeToken) {
	s.EqualValues(sdk.NewCoin(initAmount.Denom, sdkmath.ZeroInt()), s.App.BankKeeper.GetSupply(s.Ctx, initAmount.Denom))

	err := s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(initAmount))
	s.NoError(err)

	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, receiver, sdk.NewCoins(initAmount))
	s.NoError(err)
	s.EqualValues(initAmount, s.App.BankKeeper.GetSupply(s.Ctx, initAmount.Denom))

	tokenContract := helpers.GenerateAddress().Hex()
	bridgeToken = types.BridgeToken{
		ContractAddress: tokenContract,
		Denom:           initAmount.Denom,
		Supply:          initAmount.Amount,
	}
	s.Keeper().SetBridgeToken(s.Ctx, &bridgeToken)
	s.Equal(s.App.BankKeeper.GetAllBalances(s.Ctx, receiver).AmountOf(initAmount.Denom).String(), initAmount.Amount.String())

	return bridgeToken
}
