package apptesting

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/rand"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app"
	"github.com/openmetaearth/me-hub/app/params"
	daotypes "github.com/openmetaearth/me-hub/x/dao/types"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencerkeeper "github.com/openmetaearth/me-hub/x/sequencer/keeper"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

var Alice = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"

var defaultMinSequencerBond = sdk.NewCoin(params.BaseDenom, math.NewInt(1000000))

func init() {
	config := sdk.GetConfig()
	params.SetAddressPrefixes(config)
	config.Seal()
}

type KeeperTestHelper struct {
	suite.Suite
	App      *app.App
	Ctx      sdk.Context
	Dao      daotypes.DaoAddresses
	TestAccs []sdk.AccAddress
}

// InitializeDao creates random addresses for all DAO roles and stores them in the DAO keeper.
// This must be called in SetupTest after s.App and s.Ctx are initialized.
func (s *KeeperTestHelper) InitializeDao() {
	globalDao := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	meidDao := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	devOperator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	airdropAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	s.Dao = daotypes.DaoAddresses{
		GlobalDao:      globalDao.String(),
		MeidDao:        meidDao.String(),
		DevOperator:    devOperator.String(),
		AirdropAddress: airdropAddr.String(),
	}
	s.App.DaoKeeper.SetDaoAddresses(s.Ctx, s.Dao)
}

// NewAccounts creates n funded test accounts and returns their addresses.
func (s *KeeperTestHelper) NewAccounts(n int) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		addrs[i] = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
		FundAccount(s.App, s.Ctx, addrs[i], sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(1_000_000_000))))
	}
	return addrs
}

// NewAccount creates a single funded test account and returns its address and private key.
func (s *KeeperTestHelper) NewAccount() (sdk.AccAddress, *ed25519.PrivKey) {
	privKey := ed25519.GenPrivKey()
	addr := sdk.AccAddress(privKey.PubKey().Address())
	FundAccount(s.App, s.Ctx, addr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(1_000_000_000))))
	return addr, privKey
}

func (s *KeeperTestHelper) CreateDefaultRollappAndProposer() (string, string) {
	rollappId := s.CreateDefaultRollapp()
	proposer := s.CreateDefaultSequencer(s.Ctx, rollappId)
	return rollappId, proposer
}

func (s *KeeperTestHelper) CreateDefaultRollapp() string {
	rollappId := fmt.Sprintf("testrollapp%d_1-1", rand.Int63()) //nolint:gosec // this is for a test
	s.CreateRollappByName(rollappId)
	return rollappId
}

func (s *KeeperTestHelper) CreateRollappByName(name string) {
	msgCreateRollapp := rollapptypes.MsgCreateRollapp{
		Creator:       Alice,
		RollappId:     name,
		MaxSequencers: 10,
	}

	msgServer := rollappkeeper.NewMsgServerImpl(*s.App.RollappKeeper)
	_, err := msgServer.CreateRollapp(s.Ctx, &msgCreateRollapp)
	s.Require().NoError(err)
}

func (s *KeeperTestHelper) CreateDefaultSequencer(ctx sdk.Context, rollappId string) string {
	pubkey := ed25519.GenPrivKey().PubKey()
	err := s.CreateSequencerByPubkey(ctx, rollappId, pubkey)
	s.Require().NoError(err)
	return sdk.AccAddress(pubkey.Address()).String()
}

func (s *KeeperTestHelper) CreateSequencerByPubkey(ctx sdk.Context, rollappId string, pubKey types.PubKey) error {
	addr := sdk.AccAddress(pubKey.Address())
	FundAccount(s.App, ctx, addr, sdk.NewCoins(defaultMinSequencerBond))

	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	s.Require().Nil(err)

	sequencerMsg1 := sequencertypes.MsgCreateSequencer{
		Creator:      addr.String(),
		DymintPubKey: pkAny,
		Bond:         defaultMinSequencerBond,
		RollappId:    rollappId,
		Description:  sequencertypes.Description{},
	}

	msgServer := sequencerkeeper.NewMsgServerImpl(*s.App.SequencerKeeper)
	_, err = msgServer.CreateSequencer(ctx, &sequencerMsg1)
	return err
}

func (s *KeeperTestHelper) PostStateUpdate(ctx sdk.Context, rollappId, seqAddr string, startHeight, numOfBlocks uint64) (lastHeight uint64, err error) {
	var bds rollapptypes.BlockDescriptors
	bds.BD = make([]rollapptypes.BlockDescriptor, numOfBlocks)
	for k := uint64(0); k < numOfBlocks; k++ {
		bds.BD[k] = rollapptypes.BlockDescriptor{Height: startHeight + k}
	}

	updateState := rollapptypes.MsgUpdateState{
		Creator:     seqAddr,
		RollappId:   rollappId,
		StartHeight: startHeight,
		NumBlocks:   numOfBlocks,
		DAPath:      "",
		BDs:         bds,
	}
	msgServer := rollappkeeper.NewMsgServerImpl(*s.App.RollappKeeper)
	_, err = msgServer.UpdateState(ctx, &updateState)
	return startHeight + numOfBlocks, err
}

func (s *KeeperTestHelper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := bankutil.FundAccount(s.Ctx, s.App.BankKeeper, acc, amounts)
	s.Require().NoError(err)
}

func (s *KeeperTestHelper) FundModuleAcc(moduleName string, amounts sdk.Coins) {
	err := bankutil.FundModuleAccount(s.Ctx, s.App.BankKeeper, moduleName, amounts)
	s.Require().NoError(err)
}

func (s *KeeperTestHelper) FundForAliasRegistration(msgCreateRollApp rollapptypes.MsgCreateRollapp) {
	// no-op: alias registration not supported in me-hub
}

func (s *KeeperTestHelper) FinalizeAllPendingPackets(address string) int {
	// no-op: MsgFinalizePacket not supported in me-hub
	return 0
}

func (s *KeeperTestHelper) StateNotAltered() {
	oldState := s.App.ExportState(s.Ctx)
	_, err := s.App.Commit()
	s.Require().NoError(err)
	newState := s.App.ExportState(s.Ctx)
	s.Require().Equal(oldState, newState)
}
