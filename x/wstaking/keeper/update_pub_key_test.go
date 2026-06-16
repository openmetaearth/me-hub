package keeper_test

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestUpdateValidatorPubKeyDeletesOldSigningInfo() {
	s.SetupTest()

	validatorAddr, err := sdk.ValAddressFromBech32(s.meEarthValidator.OperatorAddress)
	s.Require().NoError(err)

	validator, found := s.Keeper().GetValidator(s.Ctx, validatorAddr)
	s.Require().True(found)

	oldConsAddr, err := validator.GetConsAddr()
	s.Require().NoError(err)

	oldSigningInfo := slashingtypes.NewValidatorSigningInfo(
		oldConsAddr,
		1,
		7,
		time.Unix(100, 0),
		false,
		3,
	)
	s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, oldConsAddr, oldSigningInfo)

	newPubKey := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
	newPubKeyBytes, err := newPubKey.Marshal()
	s.Require().NoError(err)
	newConsAddr := sdk.GetConsAddress(newPubKey)
	s.Require().NotEqual(oldConsAddr.String(), newConsAddr.String())

	err = s.Keeper().SetReplacePubKeyInfo(s.Ctx, &wstakingtypes.UpdatePubKeyInfo{
		OperatorAddress: validator.OperatorAddress,
		OldConsAddress:  oldConsAddr.Bytes(),
		PubKey:          newPubKeyBytes,
		UpdateAtHeight:  20,
	})
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithBlockHeight(20)
	replaceInfo, err := s.Keeper().UpdateValidatorPubKey(s.Ctx)
	s.Require().NoError(err)
	s.Require().NotNil(replaceInfo)

	_, found = s.App.SlashingKeeper.GetValidatorSigningInfo(s.Ctx, oldConsAddr)
	s.Require().True(found)
	_, found = s.App.SlashingKeeper.GetValidatorSigningInfo(s.Ctx, newConsAddr)
	s.Require().True(found)
	_, found = s.Keeper().GetValidatorByConsAddr(s.Ctx, newConsAddr)
	s.Require().True(found)

	s.Ctx = s.Ctx.WithBlockHeight(22)
	replaceInfo, err = s.Keeper().UpdateValidatorPubKey(s.Ctx)
	s.Require().NoError(err)
	s.Require().Nil(replaceInfo)

	_, found = s.App.SlashingKeeper.GetValidatorSigningInfo(s.Ctx, oldConsAddr)
	s.Require().False(found)
	_, found = s.App.SlashingKeeper.GetValidatorSigningInfo(s.Ctx, newConsAddr)
	s.Require().True(found)
	_, found = s.Keeper().GetValidatorByConsAddr(s.Ctx, oldConsAddr)
	s.Require().False(found)

	updateInfo, err := s.Keeper().GetReplaceConsensusPubKeyInfo(s.Ctx)
	s.Require().NoError(err)
	s.Require().Nil(updateInfo)
}
