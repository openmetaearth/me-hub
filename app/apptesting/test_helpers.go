package apptesting

import (
	"encoding/json"
	"testing"
	"time"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/math"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	usim "github.com/cosmos/cosmos-sdk/testutil/sims"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/stretchr/testify/require"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	app "github.com/openmetaearth/me-hub/app"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

var TestChainID = "mechain_100-1"

var DefaultConsensusParams = func() *cometbftproto.ConsensusParams {
	ret := usim.DefaultConsensusParams
	ret.Block.MaxGas = -1
	return ret
}()

// Passed into `Simapp` constructor.
type SetupOptions struct {
	Logger             log.Logger
	DB                 *dbm.MemDB
	InvCheckPeriod     uint
	HomePath           string
	SkipUpgradeHeights map[int64]bool
	EncConfig          params.EncodingConfig
	AppOpts            types.AppOptions
}

// Having this enabled led to some problems because some tests use intrusive methods to modify the state, which breaks invariants
var InvariantCheckInterval = uint(0) // disabled

func SetupTestingApp() (*app.App, app.GenesisState) {
	// Register base denom before creating app to ensure DefaultGenesis() has proper params
	params.RegisterDenomsIfNeeded()

	newApp := app.New(log.NewNopLogger(), dbm.NewMemDB(), nil, true, usim.AppOptionsMap{"skip-wasm-init": true}, bam.SetChainID(TestChainID))
	encCdc := newApp.AppCodec()
	// Use BasicModuleManager to get default genesis for all modules so that
	// module params (e.g. EVM EvmDenom, gravity MinDelegate, etc.) are initialized.
	defaultGenesisState := newApp.BasicModuleManager.DefaultGenesis(encCdc)
	// Skip wasm genesis since WasmKeeper is not initialized in test mode (skip-wasm-init)
	delete(defaultGenesisState, wasmtypes.ModuleName)
	// Skip crisis module genesis to avoid invariant checks during InitChain
	delete(defaultGenesisState, "crisis")

	// force disable EnableCreate of x/evm
	var evmGenesisState evmtypes.GenesisState
	evmGenesisStateJson := defaultGenesisState[evmtypes.ModuleName]
	if len(evmGenesisStateJson) > 0 {
		encCdc.MustUnmarshalJSON(evmGenesisStateJson, &evmGenesisState)
	} else {
		evmGenesisState = *evmtypes.DefaultGenesisState()
	}
	evmGenesisState.Params.EnableCreate = false
	defaultGenesisState[evmtypes.ModuleName] = encCdc.MustMarshalJSON(&evmGenesisState)

	return newApp, defaultGenesisState
}

// Setup initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func Setup(t *testing.T) *app.App {
	t.Helper()

	app, genesisState := SetupTestingApp()

	// create validator set with 3 validators for wstaking tests (meEarth, experience, usa)
	privVal1 := mock.NewPV()
	pubKey1, err := privVal1.GetPubKey()
	require.NoError(t, err)
	privVal2 := mock.NewPV()
	pubKey2, err := privVal2.GetPubKey()
	require.NoError(t, err)
	privVal3 := mock.NewPV()
	pubKey3, err := privVal3.GetPubKey()
	require.NoError(t, err)

	validator1 := cometbfttypes.NewValidator(pubKey1, 1)
	validator2 := cometbfttypes.NewValidator(pubKey2, 1)
	validator3 := cometbfttypes.NewValidator(pubKey3, 1)
	valSet := cometbfttypes.NewValidatorSet([]*cometbfttypes.Validator{validator1, validator2, validator3})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balances := []banktypes.Balance{
		{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, math.NewInt(1000000000000000000))),
		},
	}

	genesisState = genesisStateWithValSet(t, app, genesisState, valSet, []authtypes.GenesisAccount{acc}, balances...)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	_, err = app.InitChain(
		&abci.RequestInitChain{
			ChainId:         TestChainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	require.NoError(t, err)
	return app
}

func (s *KeeperTestHelper) Commit() {
	_, err := s.App.FinalizeBlock(&abci.RequestFinalizeBlock{Height: s.Ctx.BlockHeight(), Time: s.Ctx.BlockTime()})
	if err != nil {
		panic(err)
	}
	_, err = s.App.Commit()
	if err != nil {
		panic(err)
	}

	newBlockTime := s.Ctx.BlockTime().Add(time.Second)

	header := s.Ctx.BlockHeader()
	header.Time = newBlockTime
	header.Height++

	s.Ctx = s.App.BaseApp.NewUncachedContext(false, header).WithHeaderInfo(coreheader.Info{
		Height: header.Height,
		Time:   header.Time,
	})
}

func genesisStateWithValSet(t *testing.T,
	app *app.App, genesisState app.GenesisState,
	valSet *cometbfttypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) app.GenesisState {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	// Use larger bond amount to support wstaking tests (min delegation = 10^BaseDenomUnit)
	bondAmt := math.NewInt(1_000_000_000_000_000_000) // 10^18, same as stake pool initial balance

	// Pre-define region IDs for validators (used by wstaking tests)
	regionIDs := []string{
		wstakingtypes.MeEarthRegionId,
		wstakingtypes.ExperienceRegionId,
		"usa",
	}

	for i, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		regionID := ""
		if i < len(regionIDs) {
			regionID = regionIDs[i]
		}
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{RegionID: regionID},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), math.LegacyOneDec()))

	}
	// set validators and delegations
	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(params.BaseDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account (one per validator)
	totalBondAmt := bondAmt.MulRaw(int64(len(delegations)))
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(params.BaseDenom, totalBondAmt)},
	})
	// add bonded amount to wstaking bonded stake pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(wstakingtypes.BondedStakePoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(params.BaseDenom, totalBondAmt)},
	})
	// update total supply (add wstaking pool amount)
	totalSupply = totalSupply.Add(sdk.NewCoin(params.BaseDenom, totalBondAmt))

	// add initial balance to wstaking stake pool (stake_tokens_pool) for tests (1e18 umec)
	stakePoolAmt := math.NewInt(1_000_000_000_000_000_000) // 10^18 umec (BaseDenomUnit=18)
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(wstakingtypes.StakePoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(params.BaseDenom, stakePoolAmt)},
	})
	totalSupply = totalSupply.Add(sdk.NewCoin(params.BaseDenom, stakePoolAmt))

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	return genesisState
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// CreateRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func CreateRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *app.App, ctx sdk.Context, accNum int, accAmt math.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, CreateRandomAccounts)
}

func addTestAddrs(app *app.App, ctx sdk.Context, accNum int, accAmt math.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	denom, _ := app.StakingKeeper.BondDenom(ctx)
	initCoins := sdk.NewCoins(sdk.NewCoin(denom, accAmt))

	for _, addr := range testAddrs {
		FundAccount(app, ctx, addr, initCoins)
	}

	return testAddrs
}

func FundAccount(app *app.App, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

func FundForAliasRegistration(app *app.App, ctx sdk.Context, alias, creator string) {
	// no-op: alias registration not supported in me-hub
}

// MintBlock advances the chain by one or more blocks and returns the updated context.
func MintBlock(a *app.App, ctx sdk.Context, block ...int64) sdk.Context {
	numBlocks := int64(1)
	if len(block) > 0 && block[0] > 0 {
		numBlocks = block[0]
	}
	for i := int64(0); i < numBlocks; i++ {
		_, err := a.FinalizeBlock(&abci.RequestFinalizeBlock{Height: ctx.BlockHeight(), Time: ctx.BlockTime()})
		if err != nil {
			panic(err)
		}
		_, err = a.Commit()
		if err != nil {
			panic(err)
		}
		newBlockTime := ctx.BlockTime().Add(time.Second)
		header := ctx.BlockHeader()
		header.Time = newBlockTime
		header.Height++
		ctx = a.BaseApp.NewUncachedContext(false, header).WithHeaderInfo(coreheader.Info{
			Height: header.Height,
			Time:   header.Time,
		})
	}
	return ctx
}
