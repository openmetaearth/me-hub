package ibctesting_test

import (
	"bytes"
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	"github.com/cosmos/ibc-go/v8/testing/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app"
	"github.com/openmetaearth/me-hub/app/apptesting"
	common "github.com/openmetaearth/me-hub/x/common/types"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

// chainIDPrefix defines the default chain ID prefix for Evmos test chains
var chainIDPrefix = "evmos_9000-"

func init() {
	ibctesting.ChainIDPrefix = chainIDPrefix
	ibctesting.ChainIDSuffix = ""
	ibctesting.DefaultTestingAppInit = func() (ibctesting.TestingApp, map[string]json.RawMessage) {
		a, gs := apptesting.SetupTestingApp()
		return &apptesting.IBCTestApp{App: a}, gs
	}
}

func convertToApp(chain *ibctesting.TestChain) *app.App {
	if wrapper, ok := chain.App.(*apptesting.IBCTestApp); ok {
		return wrapper.App
	}
	a, ok := chain.App.(*app.App)
	require.True(chain.Coordinator.T, ok)

	return a
}

// utilSuite is a testing suite to test keeper functions.
type utilSuite struct {
	suite.Suite
	coordinator *ibctesting.Coordinator
}

func hubChainID() string {
	return ibctesting.GetChainID(1)
}

func cosmosChainID() string {
	return ibctesting.GetChainID(2)
}

func rollappChainID() string {
	return ibctesting.GetChainID(3)
}

func (s *utilSuite) hubChain() *ibctesting.TestChain {
	return s.coordinator.GetChain(hubChainID())
}

func (s *utilSuite) cosmosChain() *ibctesting.TestChain {
	return s.coordinator.GetChain(cosmosChainID())
}

func (s *utilSuite) rollappChain() *ibctesting.TestChain {
	return s.coordinator.GetChain(rollappChainID())
}

func (s *utilSuite) hubApp() *app.App {
	return convertToApp(s.hubChain())
}

func (s *utilSuite) rollappApp() *app.App {
	return convertToApp(s.rollappChain())
}

func (s *utilSuite) hubCtx() sdk.Context {
	return s.hubChain().GetContext()
}

func (s *utilSuite) cosmosCtx() sdk.Context {
	return s.cosmosChain().GetContext()
}

func (s *utilSuite) rollappCtx() sdk.Context {
	return s.rollappChain().GetContext()
}

func (s *utilSuite) rollappMsgServer() rollapptypes.MsgServer {
	return rollappkeeper.NewMsgServerImpl(s.hubApp().RollappKeeper)
}

// SetupTest creates a coordinator with 2 test chains.
func (s *utilSuite) SetupTest() {
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2) // initializes test chains
	s.coordinator.Chains[rollappChainID()] = s.newTestChainWithSingleValidator(s.T(), s.coordinator, rollappChainID())
}

// CreateRollappWithFinishedGenesis creates a rollapp whose 'genesis' protocol is complete:
// that is, they have finished all genesis transfers and their bridge is enabled.
func (s *utilSuite) createRollappWithFinishedGenesis(canonicalChannelID string) {
	s.createRollapp(true, &canonicalChannelID)
}

func (s *utilSuite) createRollapp(transfersEnabled bool, channelID *string) {
	creator := s.hubChain().SenderAccount.GetAddress().String()
	msgCreateRollapp := rollapptypes.NewMsgCreateRollapp(creator, rollappChainID(), creator,
		rollapptypes.DefaultMinSequencerBondGlobalCoin, "", rollapptypes.Rollapp_EVM, nil, nil, sdk.DefaultBondDenom)
	_, err := s.hubChain().SendMsgs(msgCreateRollapp)
	s.Require().NoError(err) // message committed
	if channelID != nil {
		a := s.hubApp()
		ra := a.RollappKeeper.MustGetRollapp(s.hubCtx(), rollappChainID())
		ra.ChannelId = *channelID
		if transfersEnabled {
			ra.GenesisState.TransferProofHeight = uint64(s.hubCtx().BlockHeight())
		}
		a.RollappKeeper.SetRollapp(s.hubCtx(), ra)
	}
}

func (s *utilSuite) registerSequencer() {
	bond := rollapptypes.DefaultMinSequencerBondGlobalCoin
	// fund account
	err := bankutil.FundAccount(s.hubCtx(), s.hubApp().BankKeeper, s.hubChain().SenderAccount.GetAddress(), sdk.NewCoins(bond))
	s.Require().Nil(err)

	// using validator pubkey as the dymint pubkey
	pk, err := cryptocodec.FromTmPubKeyInterface(s.rollappChain().Vals.Validators[0].PubKey)
	s.Require().Nil(err)

	msgCreateSequencer, err := sequencertypes.NewMsgCreateSequencer(
		s.hubChain().SenderAccount.GetAddress().String(),
		pk,
		rollappChainID(),
		&sequencertypes.SequencerMetadata{},
		bond,
		s.hubChain().SenderAccount.GetAddress().String(),
		nil,
	)
	s.Require().NoError(err) // message committed
	_, err = s.hubChain().SendMsgs(msgCreateSequencer)
	s.Require().NoError(err) // message committed
}

func (s *utilSuite) updateRollappState(endHeight uint64) {
	// Get the start index and start height based on the latest state info
	rollappKeeper := s.hubApp().RollappKeeper
	latestStateInfoIndex, _ := rollappKeeper.GetLatestStateInfoIndex(s.hubCtx(), rollappChainID())
	stateInfo, found := rollappKeeper.GetStateInfo(s.hubCtx(), rollappChainID(), latestStateInfoIndex.Index)
	startHeight := uint64(1)
	if found {
		startHeight = stateInfo.StartHeight + stateInfo.NumBlocks
	}
	numBlocks := endHeight - startHeight + 1
	// populate the block descriptors
	blockDescriptors := &rollapptypes.BlockDescriptors{BD: make([]rollapptypes.BlockDescriptor, numBlocks)}
	for i := 0; i < int(numBlocks); i++ {
		blockDescriptors.BD[i] = rollapptypes.BlockDescriptor{
			Height:    startHeight + uint64(i),
			StateRoot: bytes.Repeat([]byte{byte(startHeight) + byte(i)}, 32),
		}
	}
	// Update the state
	msgUpdateState := rollapptypes.NewMsgUpdateState(
		s.hubChain().SenderAccount.GetAddress().String(),
		rollappChainID(),
		"mock-da-path",
		startHeight,
		endHeight-startHeight+1, // numBlocks
		0,
		blockDescriptors,
	)
	err := msgUpdateState.ValidateBasic()
	s.Require().NoError(err)
	_, err = s.rollappMsgServer().UpdateState(s.hubCtx(), msgUpdateState)
	s.Require().NoError(err)
}

func (s *utilSuite) finalizeRollappState(index, endHeight uint64) (sdk.Events, error) {
	rollappKeeper := s.hubApp().RollappKeeper
	ctx := s.hubCtx()

	stateInfoIdx := rollapptypes.StateInfoIndex{RollappId: rollappChainID(), Index: index}
	stateInfo, found := rollappKeeper.GetStateInfo(ctx, rollappChainID(), stateInfoIdx.Index)
	s.Require().True(found)
	stateInfo.NumBlocks = endHeight - stateInfo.StartHeight + 1
	stateInfo.Status = common.Status_FINALIZED
	// update the status of the stateInfo
	rollappKeeper.SetStateInfo(ctx, stateInfo)
	// update the LatestStateInfoIndex of the rollapp
	rollappKeeper.SetLatestFinalizedStateIndex(ctx, stateInfoIdx)
	err := rollappKeeper.GetHooks().AfterStateFinalized(
		ctx,
		rollappChainID(),
		&stateInfo,
	)

	return ctx.EventManager().Events(), err
}

func (s *utilSuite) newTransferPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	path.EndpointA.ChannelConfig.Version = transfertypes.Version
	path.EndpointB.ChannelConfig.Version = transfertypes.Version

	return path
}

func (s *utilSuite) getRollappToHubIBCDenomFromPacket(packet channeltypes.Packet) string {
	var data transfertypes.FungibleTokenPacketData
	err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data)
	s.Require().NoError(err)
	return s.getIBCDenomForChannel(packet.GetDestChannel(), data.Denom)
}

func (s *utilSuite) getIBCDenomForChannel(channel, denom string) string {
	// since SendPacket did not prefix the denomination, we must prefix denomination here
	sourcePrefix := transfertypes.GetDenomPrefix("transfer", channel)
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + denom
	// construct the denomination trace from the full raw denomination
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
	return denomTrace.IBCDenom()
}

func (s *utilSuite) newTestChainWithSingleValidator(t *testing.T, coord *ibctesting.Coordinator, chainID string) *ibctesting.TestChain {
	genAccs := []authtypes.GenesisAccount{}
	genBals := []banktypes.Balance{}
	senderAccs := []ibctesting.SenderAccount{}

	// generate genesis accounts

	valPrivKey := mock.NewPV()
	valPubKey, err := valPrivKey.GetPubKey()
	s.Require().NoError(err)

	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)

	amount, ok := sdkmath.NewIntFromString("10000000000000000000")
	s.Require().True(ok)

	// add sender account
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amount)),
	}

	genAccs = append(genAccs, acc)
	genBals = append(genBals, balance)

	senderAcc := ibctesting.SenderAccount{
		SenderAccount: acc,
		SenderPrivKey: senderPrivKey,
	}

	senderAccs = append(senderAccs, senderAcc)

	var validators []*cometbfttypes.Validator
	signersByAddress := make(map[string]cometbfttypes.PrivValidator, 1)

	validators = append(validators, cometbfttypes.NewValidator(valPubKey, 1))

	signersByAddress[valPubKey.Address().String()] = valPrivKey
	valSet := cometbfttypes.NewValidatorSet(validators)

	app := ibctesting.SetupWithGenesisValSet(t, valSet, genAccs, chainID, sdk.DefaultPowerReduction, genBals...)

	// create current header and call begin block
	header := cometbftproto.Header{
		ChainID: chainID,
		Height:  1,
		Time:    coord.CurrentTime.UTC(),
	}

	txConfig := app.GetTxConfig()

	// create an account to send transactions from
	chain := &ibctesting.TestChain{
		Coordinator:    coord,
		ChainID:        chainID,
		App:            app,
		CurrentHeader:  header,
		QueryServer:    app.GetIBCKeeper(),
		TxConfig:       txConfig,
		Codec:          app.AppCodec(),
		Vals:           valSet,
		NextVals:       valSet,
		Signers:        signersByAddress,
		SenderPrivKey:  senderAcc.SenderPrivKey,
		SenderAccount:  senderAcc.SenderAccount,
		SenderAccounts: senderAccs,
	}

	coord.CommitBlock(chain)

	return chain
}
