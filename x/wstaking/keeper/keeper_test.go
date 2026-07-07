package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	testutilstypes "github.com/openmetaearth/me-hub/testutil/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer           wstakingkeeper.MsgServer
	queryClient         types.QueryClient
	meEarthValidator    stakingtypes.Validator
	experienceValidator stakingtypes.Validator
	usaValidator        stakingtypes.Validator
	TestAccs            []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) Keeper() *wstakingkeeper.Keeper {
	return s.App.StakingKeeper
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T())
	ctx := app.GetBaseApp().NewContext(false)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	err := app.StakingKeeper.SetParams(ctx, stakingParams)
	s.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	nativeQuerier := wstakingkeeper.Querier{Keeper: app.StakingKeeper}
	types.RegisterQueryServer(queryHelper, nativeQuerier)
	queryClient := types.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.App = app
	s.Ctx = ctx

	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper.Keeper)
	s.msgServer = wstakingkeeper.NewMsgServerImpl(app.StakingKeeper, app.TransferKeeper, stakingKeeperMsgSrv)

	s.InitializeDao()

	validators, err := s.Keeper().GetValidators(s.Ctx, 10)
	s.Require().NoError(err)
	s.Require().True(len(validators) >= 3)
	s.meEarthValidator = validators[0]
	s.experienceValidator = validators[1]
	s.usaValidator = validators[2]

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.ExperienceRegionName,
		OperatorAddress: s.experienceValidator.OperatorAddress,
	}
	_, err = s.msgServer.NewRegion(s.Ctx, &newRegion)

	s.Require().NoError(err)

	s.TestAccs = s.NewAccounts(3)
}

// InitKyc sets up KYC for the given account with a DID and region ID.
func (s *KeeperTestSuite) InitKyc(account sdk.AccAddress, did string, regionId string) {
	s.App.KycKeeper.SetDID(s.Ctx, account, did)
	s.App.KycKeeper.SetKYC(s.Ctx, did, didtypes.Credential{
		Did:  did,
		Sid:  "kyc",
		Data: []byte(regionId),
	})
}

func SetValidatorV1(ctx sdk.Context, k *wstakingkeeper.Keeper, validator testutilstypes.ValidatorV1) {
	store := ctx.KVStore(k.GetStoreKey())
	bz := k.GetCdc().MustMarshal(&validator)
	addr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
	if err != nil {
		panic(err)
	}
	store.Set(stakingtypes.GetValidatorKey(addr), bz)
}

func GetValidatorV2(ctx sdk.Context, k *wstakingkeeper.Keeper, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.GetStoreKey())
	value := store.Get(stakingtypes.GetValidatorKey(addr))
	if value == nil {
		return validator, false
	}
	err := k.GetCdc().Unmarshal(value, &validator)
	if err != nil {
		return validator, false
	}
	return validator, true
}

func (s *KeeperTestSuite) TestMigrateValidator() {
	val1 := testutilstypes.ValidatorV1{
		OperatorAddress: "mevaloper139mq752delxv78jvtmwxhasyrycufsvr707ate",
		ConsensusPubkey: nil,
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          sdkmath.NewInt(100),
		StakerShares:    sdkmath.LegacyNewDec(100),
		Description: stakingtypes.Description{
			Moniker:         "node1",
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
			RegionID:        "usa",
		},
		UnbondingHeight:         0,
		UnbondingTime:           time.Time{},
		Commission:              stakingtypes.Commission{},
		MinSelfStake:            sdkmath.Int{},
		DelegationAmount:        sdkmath.Int{},
		MeidAmount:              sdkmath.Int{},
		OwnerAddress:            "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u",
		UnbondingIds:            nil,
		UnbondingOnHoldRefCount: 0,
	}
	SetValidatorV1(s.Ctx, s.App.StakingKeeper, val1)
	s.T().Log(val1.String())

	addr, err := sdk.ValAddressFromBech32(val1.OperatorAddress)
	if err != nil {
		panic(err)
	}
	// run migration to convert V1 validator format to current V2 format
	err = s.App.StakingKeeper.MigrateValidatorsFromV1(s.Ctx)
	require.NoError(s.T(), err)

	validator, found := GetValidatorV2(s.Ctx, s.App.StakingKeeper, addr)
	require.True(s.T(), found)

	validators, err := s.App.StakingKeeper.GetAllValidators(s.Ctx)
	require.NoError(s.T(), err)
	require.Equal(s.T(), len(validators), 4)
	for _, v := range validators {
		if v.OperatorAddress == validator.OperatorAddress {
			s.T().Log(validator.String())
			require.Equal(s.T(), validator.String(), v.String())
		}
	}
}
