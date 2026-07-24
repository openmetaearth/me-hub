package keeper_test

import (
	"encoding/hex"
	"math/big"

	sdkmath "cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"

	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) newCreateValidatorMsgForTest(operatorAcc sdk.AccAddress, pubKey *ed25519.PubKey) *stakingtypes.MsgCreateValidator {
	validatorPubKeyAny, err := codectypes.NewAnyWithValue(pubKey)
	s.Require().NoError(err)

	createAmount := sdk.NewCoin(
		params.BaseDenom,
		sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit), nil)),
	)

	minCommissionRate, err := s.Keeper().MinCommissionRate(s.Ctx)
	s.Require().NoError(err)

	return &stakingtypes.MsgCreateValidator{
		Description: stakingtypes.Description{
			Moniker:  "replacement-conflict",
			RegionID: wstakingtypes.MeEarthRegionName,
		},
		Commission: stakingtypes.NewCommissionRates(
			minCommissionRate,
			sdkmath.LegacyOneDec(),
			sdkmath.LegacyZeroDec(),
		),
		MinSelfDelegation: sdkmath.OneInt(),
		DelegatorAddress:  s.Dao.GlobalDao,
		ValidatorAddress:  sdk.ValAddress(operatorAcc).String(),
		Pubkey:            validatorPubKeyAny,
		Value:             createAmount,
	}
}

func (s *KeeperTestSuite) TestMsgServerCreateValidatorRejectsPendingReplacementPubKey() {
	s.SetupTest()

	pendingPubKey, ok := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
	s.Require().True(ok)

	replaceReq, err := wstakingtypes.NewMsgReplaceConsensusPubKeyRequest(
		s.Dao.GlobalDao,
		s.meEarthValidator.OperatorAddress,
		pendingPubKey,
		wstakingtypes.MinReplacePubKeyBlockNumber,
	)
	s.Require().NoError(err)

	_, err = s.msgServer.ReplaceConsensusPubKey(s.Ctx, replaceReq)
	s.Require().NoError(err)

	pubKeyData, err := pendingPubKey.Marshal()
	s.Require().NoError(err)

	msg := s.newCreateValidatorMsgForTest(s.TestAccs[0], pendingPubKey)
	_, err = s.msgServer.CreateValidator(s.Ctx, msg)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), hex.EncodeToString(pubKeyData))
}

func (s *KeeperTestSuite) TestMsgServerCreateValidatorAllowsDifferentPendingReplacementPubKey() {
	s.SetupTest()

	pendingPubKey, ok := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
	s.Require().True(ok)

	replaceReq, err := wstakingtypes.NewMsgReplaceConsensusPubKeyRequest(
		s.Dao.GlobalDao,
		s.meEarthValidator.OperatorAddress,
		pendingPubKey,
		wstakingtypes.MinReplacePubKeyBlockNumber,
	)
	s.Require().NoError(err)

	_, err = s.msgServer.ReplaceConsensusPubKey(s.Ctx, replaceReq)
	s.Require().NoError(err)

	createPubKey, ok := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
	s.Require().True(ok)

	msg := s.newCreateValidatorMsgForTest(s.TestAccs[1], createPubKey)
	_, err = s.msgServer.CreateValidator(s.Ctx, msg)
	s.Require().NoError(err)

	validatorAddr := sdk.ValAddress(s.TestAccs[1])
	validator, err := s.Keeper().GetValidator(s.Ctx, validatorAddr)
	s.Require().NoError(err)

	validatorConsAddr, err := validator.GetConsAddr()
	s.Require().NoError(err)
	s.Require().Equal(sdk.GetConsAddress(createPubKey).String(), sdk.ConsAddress(validatorConsAddr).String())
}
