package keeper_test

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestReplaceConsensusPubKeyPreservesDowntimeSigningInfo() {
	validator := s.experienceValidator
	oldConsAddr, err := validator.GetConsAddr()
	s.Require().NoError(err)

	oldSigningInfo := slashingtypes.NewValidatorSigningInfo(
		oldConsAddr,
		1,
		99,
		time.Unix(123, 0).UTC(),
		false,
		50,
	)
	s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, oldConsAddr, oldSigningInfo)
	s.App.SlashingKeeper.SetValidatorMissedBlockBitArray(s.Ctx, oldConsAddr, 3, true)
	s.App.SlashingKeeper.SetValidatorMissedBlockBitArray(s.Ctx, oldConsAddr, 17, true)

	newPubKey := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
	newPubKeyData, err := newPubKey.Marshal()
	s.Require().NoError(err)
	updateHeight := int64(111)

	s.Ctx = s.Ctx.WithBlockHeight(updateHeight)
	err = s.Keeper().SetReplacePubKeyInfo(s.Ctx, &types.UpdatePubKeyInfo{
		OperatorAddress: validator.OperatorAddress,
		OldConsAddress:  oldConsAddr.Bytes(),
		PubKey:          newPubKeyData,
		UpdateAtHeight:  updateHeight,
	})
	s.Require().NoError(err)

	replacePubKey, err := s.Keeper().UpdateValidatorPubKey(s.Ctx)
	s.Require().NoError(err)
	s.Require().NotNil(replacePubKey)

	newConsAddr := sdk.GetConsAddress(newPubKey)
	newSigningInfo, found := s.App.SlashingKeeper.GetValidatorSigningInfo(s.Ctx, newConsAddr)
	s.Require().True(found)
	s.Require().Equal(newConsAddr.String(), newSigningInfo.Address)
	s.Require().Equal(oldSigningInfo.StartHeight, newSigningInfo.StartHeight)
	s.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)
	s.Require().Equal(oldSigningInfo.JailedUntil, newSigningInfo.JailedUntil)
	s.Require().Equal(oldSigningInfo.Tombstoned, newSigningInfo.Tombstoned)
	s.Require().Equal(oldSigningInfo.MissedBlocksCounter, newSigningInfo.MissedBlocksCounter)
	s.Require().True(s.App.SlashingKeeper.GetValidatorMissedBlockBitArray(s.Ctx, newConsAddr, 3))
	s.Require().True(s.App.SlashingKeeper.GetValidatorMissedBlockBitArray(s.Ctx, newConsAddr, 17))
}
