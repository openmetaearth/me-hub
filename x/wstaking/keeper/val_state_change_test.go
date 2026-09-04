package keeper_test

import (
	"encoding/hex"
	"math/big"

	sdkmath "cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking"

	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func validatorUpdatePowersByConsAddr(s *KeeperTestSuite, updates []abci.ValidatorUpdate) map[string][]int64 {
	s.T().Helper()

	powersByAddr := make(map[string][]int64, len(updates))
	for _, update := range updates {
		pubKey, err := cryptocodec.FromTmProtoPublicKey(update.PubKey)
		s.Require().NoError(err)
		consAddr := sdk.GetConsAddress(pubKey).String()
		powersByAddr[consAddr] = append(powersByAddr[consAddr], update.Power)
	}

	return powersByAddr
}

func (s *KeeperTestSuite) TestBlockValidatorUpdatesDeduplicatesOldConsensusPubKeyOnSameBlockPowerChange() {
	s.SetupTest()

	_, err := s.msgServer.NewRegion(s.Ctx, &wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            wstakingtypes.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	})
	s.Require().NoError(err)

	validatorAddr, err := sdk.ValAddressFromBech32(s.meEarthValidator.OperatorAddress)
	s.Require().NoError(err)

	validatorBefore, err := s.Keeper().GetValidator(s.Ctx, validatorAddr)
	s.Require().NoError(err)

	oldConsAddr, err := validatorBefore.GetConsAddr()
	s.Require().NoError(err)

	newPubKey := ed25519.GenPrivKey().PubKey()
	replaceReq, err := wstakingtypes.NewMsgReplaceConsensusPubKeyRequest(
		s.Dao.GlobalDao,
		s.meEarthValidator.OperatorAddress,
		newPubKey,
		wstakingtypes.MinReplacePubKeyBlockNumber,
	)
	s.Require().NoError(err)

	_, err = s.msgServer.ReplaceConsensusPubKey(s.Ctx, replaceReq)
	s.Require().NoError(err)

	updateInfo, err := s.Keeper().GetReplaceConsensusPubKeyInfo(s.Ctx)
	s.Require().NoError(err)
	s.Require().NotNil(updateInfo)

	s.Ctx = s.Ctx.WithBlockHeight(updateInfo.UpdateAtHeight)

	stakeAmount := sdk.NewCoin(params.BaseDenom, sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit), nil)))
	_, err = s.msgServer.Stake(s.Ctx, &wstakingtypes.MsgStake{
		StakerAddress:    s.Dao.GlobalDao,
		ValidatorAddress: s.meEarthValidator.OperatorAddress,
		Amount:           stakeAmount,
	})
	s.Require().NoError(err)

	updates := wstaking.EndBlock(s.Ctx, s.Keeper())
	powersByAddr := validatorUpdatePowersByConsAddr(s, updates)

	s.Require().Len(powersByAddr[sdk.ConsAddress(oldConsAddr).String()], 1)
	s.Require().EqualValues(0, powersByAddr[sdk.ConsAddress(oldConsAddr).String()][0])

	newConsAddr := sdk.GetConsAddress(newPubKey).String()
	s.Require().Len(powersByAddr[newConsAddr], 1)

	validatorAfter, err := s.Keeper().GetValidator(s.Ctx, validatorAddr)
	s.Require().NoError(err)
	s.Require().Equal(validatorAfter.ConsensusPower(s.Keeper().PowerReduction(s.Ctx)), powersByAddr[newConsAddr][0])
}

func (s *KeeperTestSuite) TestCreateValidatorRejectsPendingReplacementPubKey() {
	s.SetupTest()

	newPubKey := ed25519.GenPrivKey().PubKey()
	replaceReq, err := wstakingtypes.NewMsgReplaceConsensusPubKeyRequest(
		s.Dao.GlobalDao,
		s.meEarthValidator.OperatorAddress,
		newPubKey,
		wstakingtypes.MinReplacePubKeyBlockNumber,
	)
	s.Require().NoError(err)

	_, err = s.msgServer.ReplaceConsensusPubKey(s.Ctx, replaceReq)
	s.Require().NoError(err)

	validatorPubKeyAny, err := codectypes.NewAnyWithValue(newPubKey)
	s.Require().NoError(err)

	newEd25519PubKey, ok := newPubKey.(*ed25519.PubKey)
	s.Require().True(ok)

	pubKeyData, err := newEd25519PubKey.Marshal()
	s.Require().NoError(err)

	createAmount := sdk.NewCoin(params.BaseDenom, sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(params.BaseDenomUnit), nil)))
	minCommissionRate, err := s.Keeper().MinCommissionRate(s.Ctx)
	s.Require().NoError(err)
	_, err = s.msgServer.CreateValidator(s.Ctx, &stakingtypes.MsgCreateValidator{
		Description: stakingtypes.Description{
			Moniker:  "replacement-conflict",
			RegionID: wstakingtypes.MeEarthRegionName,
		},
		Commission:        stakingtypes.NewCommissionRates(minCommissionRate, sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec()),
		MinSelfDelegation: sdkmath.OneInt(),
		DelegatorAddress:  s.Dao.GlobalDao,
		ValidatorAddress:  sdk.ValAddress(s.TestAccs[0]).String(),
		Pubkey:            validatorPubKeyAny,
		Value:             createAmount,
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), hex.EncodeToString(pubKeyData))
}
