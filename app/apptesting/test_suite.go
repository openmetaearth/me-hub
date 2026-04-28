package apptesting

import (
	"fmt"
	"github.com/cometbft/cometbft/libs/rand"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/openmetaearth/me-hub/app"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/dao/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/suite"

	bankutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"

	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencerkeeper "github.com/openmetaearth/me-hub/x/sequencer/keeper"

	daotypes "github.com/openmetaearth/me-hub/x/dao/types"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

var (
	alice = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"
	bond  = sequencertypes.DefaultParams().MinBond
)

type KeeperTestHelper struct {
	suite.Suite
	App *app.App
	Ctx sdk.Context
	Dao daotypes.DaoAddresses
}

func (s *KeeperTestHelper) CreateDefaultRollapp() string {
	return s.CreateRollappWithName(rand.Str(8))
}

func (s *KeeperTestHelper) CreateRollappWithName(name string) string {
	msgCreateRollapp := rollapptypes.MsgCreateRollapp{
		Creator:       alice,
		RollappId:     name,
		MaxSequencers: 5,
	}

	msgServer := rollappkeeper.NewMsgServerImpl(*s.App.RollappKeeper)
	_, err := msgServer.CreateRollapp(s.Ctx, &msgCreateRollapp)
	s.Require().NoError(err)
	return name
}

func (s *KeeperTestHelper) CreateDefaultSequencer(ctx sdk.Context, rollappId string) string {
	pubkey1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pubkey1.Address())
	pkAny1, err := codectypes.NewAnyWithValue(pubkey1)
	s.Require().Nil(err)

	// fund account
	err = bankutil.FundAccount(s.App.BankKeeper, ctx, addr1, sdk.NewCoins(bond))
	s.Require().Nil(err)

	sequencerMsg1 := sequencertypes.MsgCreateSequencer{
		Creator:      addr1.String(),
		DymintPubKey: pkAny1,
		Bond:         bond,
		RollappId:    rollappId,
		Description:  sequencertypes.Description{},
	}

	msgServer := sequencerkeeper.NewMsgServerImpl(s.App.SequencerKeeper)
	_, err = msgServer.CreateSequencer(ctx, &sequencerMsg1)
	s.Require().Nil(err)
	return addr1.String()
}

func (s *KeeperTestHelper) PostStateUpdate(ctx sdk.Context, rollappId, seqAddr string, startHeight, numOfBlocks uint64) (lastHeight uint64, err error) {
	var bds rollapptypes.BlockDescriptors
	bds.BD = make([]rollapptypes.BlockDescriptor, numOfBlocks)
	for k := 0; k < int(numOfBlocks); k++ {
		bds.BD[k] = rollapptypes.BlockDescriptor{Height: startHeight + uint64(k)}
	}

	updateState := rollapptypes.MsgUpdateState{
		Creator:     seqAddr,
		RollappId:   rollappId,
		StartHeight: startHeight,
		NumBlocks:   numOfBlocks,
		DAPath:      "",
		Version:     0,
		BDs:         bds,
	}
	msgServer := rollappkeeper.NewMsgServerImpl(*s.App.RollappKeeper)
	_, err = msgServer.UpdateState(ctx, &updateState)
	return startHeight + numOfBlocks, err
}

// FundAcc funds target address with specified amount.
func (s *KeeperTestHelper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := bankutil.FundAccount(s.App.BankKeeper, s.Ctx, acc, amounts)
	s.Require().NoError(err)
}

// FundModuleAcc funds target modules with specified amount.
func (suite *KeeperTestHelper) FundModuleAcc(moduleName string, amounts sdk.Coins) {
	err := bankutil.FundModuleAccount(suite.App.BankKeeper, suite.Ctx, moduleName, amounts)
	suite.Require().NoError(err)
}

// StateNotAltered validates that app state is not altered. Fails if it is.
func (suite *KeeperTestHelper) StateNotAltered() {
	oldState := suite.App.ExportState(suite.Ctx)
	suite.App.Commit()
	newState := suite.App.ExportState(suite.Ctx)
	suite.Require().Equal(oldState, newState)
}

func (s *KeeperTestHelper) InitializeDao() {
	globalDaoPrivKey, _ := ethsecp256k1.GenerateKey()
	globalDaoAcc := authtypes.NewBaseAccount(globalDaoPrivKey.PubKey().Address().Bytes(), globalDaoPrivKey.PubKey(), 1, 0)
	//globalOutput, _ := keyring.NewKeyOutput("global_dao", keyring.TypeLocal, globalDaoAddress, globalDaoPrivKey.PubKey())

	meidDao, _ := ethsecp256k1.GenerateKey()
	meidDaoAcc := authtypes.NewBaseAccount(meidDao.PubKey().Address().Bytes(), meidDao.PubKey(), 1, 0)

	devOperator, _ := ethsecp256k1.GenerateKey()
	devOperatorAcc := authtypes.NewBaseAccount(devOperator.PubKey().Address().Bytes(), devOperator.PubKey(), 2, 0)

	airdrop, _ := ethsecp256k1.GenerateKey()
	airdropAcc := authtypes.NewBaseAccount(airdrop.PubKey().Address().Bytes(), airdrop.PubKey(), 3, 0)
	airdropAddress := sdk.AccAddress(airdrop.PubKey().Address().Bytes())

	s.App.DaoKeeper.SetDaoAddresses(s.Ctx, types.DaoAddresses{
		GlobalDao:      globalDaoAcc.Address,
		MeidDao:        meidDaoAcc.Address,
		DevOperator:    devOperatorAcc.Address,
		AirdropAddress: airdropAcc.Address,
	})

	dao, found := s.App.DaoKeeper.GetDaoAddresses(s.Ctx)
	s.Require().True(found)
	s.Dao = dao

	s.InitKyc(globalDaoAcc.GetAddress(), "0000000000001", wstakingtypes.MeEarthRegionId)

	_ = s.App.BankKeeper.MintCoins(s.Ctx, mintypes.ModuleName, sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000000000)})
	_ = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, mintypes.ModuleName, globalDaoAcc.GetAddress(), sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000)})
	_ = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, mintypes.ModuleName, airdropAddress, sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 1000000000000)})
}

func (s *KeeperTestHelper) InitKyc(address sdk.AccAddress, did string, regionId string) {
	//address, _ := s.App.KycKeeper.MustAccAddressFromPubkeyString(pubkey)
	if _, found := s.App.KycKeeper.GetDID(s.Ctx, address); found {
		panic(fmt.Errorf("issuer %s already exists", address))
	}

	s.App.DidKeeper.SetDID(s.Ctx, address, did)
	s.App.DidKeeper.SetDidInfo(s.Ctx, did, didtypes.DidInfo{
		Did:     did,
		Address: address.String(),
		Pubkey:  "",
		Status:  didtypes.DID_STATUS_ACTIVE,
	})

	service := didtypes.Service{
		Sid:         kyctypes.ModuleName,
		Name:        kyctypes.ModuleName,
		Description: "The KYC verifiable credential issuer based The DID(Decentralized Identity).",
		Issuers:     []string{did},
		Status:      didtypes.SERVICE_STATUS_ACTIVE,
	}
	s.App.DidKeeper.SetService(s.Ctx, service.Sid, service)

	kyc := didtypes.NewCredential(did, service.Sid, "", "", []byte(regionId))
	s.App.KycKeeper.SetKYC(s.Ctx, did, kyc)
	s.App.KycKeeper.AddFilters(s.Ctx, did, [][]byte{[]byte(regionId)}, kyc)
}

func (s *KeeperTestHelper) NewAccount() (sdk.AccAddress, string) {
	globalDaoPrivKey, _ := ethsecp256k1.GenerateKey()
	globalDaoAddress := sdk.AccAddress(globalDaoPrivKey.PubKey().Address().Bytes())
	globalOutput, _ := keyring.NewKeyOutput("global_dao", keyring.TypeLocal, globalDaoAddress, globalDaoPrivKey.PubKey())
	return globalDaoAddress, globalOutput.PubKey
}

func (s *KeeperTestHelper) NewAccounts(count int) []sdk.AccAddress {
	accounts := make([]sdk.AccAddress, count)
	for i := 0; i < count; i++ {
		key, _ := ethsecp256k1.GenerateKey()
		address := sdk.AccAddress(key.PubKey().Address().Bytes())
		//account := authtypes.NewBaseAccount(airdrop.PubKey().Address().Bytes(), airdrop.PubKey(), 3, 0)
		accounts[i] = address
	}
	return accounts
}
