package keeper_test

import (
	"crypto/ecdsa"
	"fmt"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/suite"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	meEarthValidator stakingtypes.Validator

	relayerAddrs  []sdk.AccAddress
	relayerNumber int
	externalPris  []*ecdsa.PrivateKey
	chainName     string

	queryServer gravitytypes.QueryClient
	msgServer   gravitytypes.MsgServer
	signer      *helpers.Signer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	s.Ctx = app.NewContext(false, cometbftproto.Header{Height: 0, ChainID: apptesting.TestChainID})
	s.App = app

	err := app.AccountKeeper.SetParams(s.Ctx, authtypes.DefaultParams())
	s.Require().NoError(err)

	err = app.BankKeeper.SetParams(s.Ctx, banktypes.DefaultParams())
	s.Require().NoError(err)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	err = app.StakingKeeper.SetParams(s.Ctx, stakingParams)
	s.Require().NoError(err)

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	stakingMsgServer := wstakingkeeper.NewMsgServerImpl(app.StakingKeeper, app.TransferKeeper, stakingKeeperMsgSrv)

	s.InitializeDao()

	validators := s.App.StakingKeeper.GetValidators(s.Ctx, 10)
	s.Require().True(len(validators) >= 3)
	s.meEarthValidator = validators[0]

	newRegion := wstakingtypes.MsgNewRegion{
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

	err = s.App.TronKeeper.SetParams(s.Ctx, &gravitytypes.Params{
		GravityId:                          "me-tron-bridge",
		AverageBlockTime:                   5000,
		ExternalBatchTimeout:               24 * 3600 * 1000, // 24 hours
		AverageExternalBlockTime:           3000,
		SignedWindow:                       30_000,
		SlashFraction:                      sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		RelayerSetUpdatePowerChangePercent: sdk.MustNewDecFromStr("0.2"),
		MaxRelayers:                        10,
		MinDelegate:                        sdk.NewInt(1000000000),
		MaxDelegate:                        sdk.NewInt(100000000000),
	})
	s.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(s.Ctx, s.App.InterfaceRegistry())
	gravitytypes.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(s.App.TronKeeper))
	s.queryServer = gravitytypes.NewQueryClient(queryHelper)

	s.msgServer = keeper.NewMsgServerImpl(s.App.TronKeeper)
	s.signer = helpers.NewSigner(helpers.NewEthPrivKey())
	apptesting.AddTestAddr(s.App, s.Ctx, s.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e8)))))
}

func (s *KeeperTestSuite) NewOutgoingTxBatch() *gravitytypes.OutgoingTxBatch {
	batchNonce := tmrand.Uint64()
	tokenContract := helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex())
	newOutgoingTx := &gravitytypes.OutgoingTxBatch{
		BatchNonce: batchNonce,
		Transactions: []*gravitytypes.OutgoingTransferTx{
			{
				Sender:      s.signer.AccAddress().String(),
				DestAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
				Token: gravitytypes.ERC20Token{
					Contract: tokenContract,
					Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
				},
				Fee: gravitytypes.ERC20Token{
					Contract: tokenContract,
					Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
				},
			},
		},
		TokenContract: tokenContract,
		FeeReceive:    helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
		Block:         batchNonce,
	}
	err := s.App.TronKeeper.StoreBatch(s.Ctx, newOutgoingTx)
	s.Require().NoError(err)
	return newOutgoingTx
}

func (s *KeeperTestSuite) NewRelayer() (sdk.AccAddress, cryptotypes.PrivKey) {
	relayer := helpers.GenAccAddress()
	externalKey := helpers.NewEthPrivKey()
	externalAddress := helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String())
	newRelayer := gravitytypes.Relayer{
		RelayerAddress:  relayer.String(),
		ExternalAddress: externalAddress,
	}
	s.App.TronKeeper.SetRelayer(s.Ctx, relayer, newRelayer)
	s.App.TronKeeper.SetRelayerByExternalAddress(s.Ctx, externalAddress, relayer)
	return relayer, externalKey
}

func (s *KeeperTestSuite) CurrentRelayerSet(externalKey cryptotypes.PrivKey) *gravitytypes.RelayerSet {
	currentRelayerSet := gravitytypes.CurrentRelayerSet(tmrand.Uint64(), tmrand.Uint64(), gravitytypes.BridgeValidators{
		{
			Power:           tmrand.Uint64(),
			ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
		},
	})
	s.App.TronKeeper.StoreRelayerSet(s.Ctx, currentRelayerSet)
	return currentRelayerSet
}

func (s *KeeperTestSuite) NewBridgeToken(bridger sdk.AccAddress) []gravitytypes.BridgeToken {
	bridgeTokens := make([]gravitytypes.BridgeToken, 0)
	for i := 0; i < 3; i++ {
		bt := gravitytypes.BridgeToken{
			ContractAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
			Denom:           fmt.Sprintf("test%d", i),
			Name:            "",
			Symbol:          fmt.Sprintf("test%d", i),
			Decimal:         0,
			Supply:          sdk.NewInt(0),
		}
		err := s.App.TronKeeper.AttestationHandler(s.Ctx, &gravitytypes.MsgBridgeTokenClaim{
			TokenContract:  bt.ContractAddress,
			Symbol:         bt.Denom,
			RelayerAddress: bridger.String(),
		})
		s.Require().NoError(err)
		_, err = s.App.TronKeeper.GetBridgeTokenByContract(s.Ctx, bt.ContractAddress)
		s.Require().NoError(err)
		bt.Denom = utils.GetDenom(bt.Denom)
		bridgeTokens = append(bridgeTokens, bt)
	}
	return bridgeTokens
}
