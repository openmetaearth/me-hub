package keeper_test

import (
	"time"

	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestUpdateValidatorPubKeyPreservesDowntimeSigningInfo() {
	s.SetupTest()

	ctx := s.Ctx.WithBlockHeight(111)
	validator := s.meEarthValidator
	s.Require().True(validator.IsBonded(), "test setup needs a bonded validator")

	oldConsAddr, err := validator.GetConsAddr()
	s.Require().NoError(err)

	oldSigningInfo := slashingtypes.NewValidatorSigningInfo(
		oldConsAddr,
		1,
		99,
		time.Unix(123, 0).UTC(),
		false,
		2,
	)
	s.App.SlashingKeeper.SetValidatorSigningInfo(ctx, oldConsAddr, oldSigningInfo)
	s.App.SlashingKeeper.SetValidatorMissedBlockBitArray(ctx, oldConsAddr, 0, true)
	s.App.SlashingKeeper.SetValidatorMissedBlockBitArray(ctx, oldConsAddr, 7, true)

	newPubKey := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
	pubKeyData, err := newPubKey.Marshal()
	s.Require().NoError(err)

	err = s.Keeper().SetReplacePubKeyInfo(ctx, &types.UpdatePubKeyInfo{
		OperatorAddress: validator.OperatorAddress,
		OldConsAddress:  oldConsAddr.Bytes(),
		PubKey:          pubKeyData,
		UpdateAtHeight:  ctx.BlockHeight(),
	})
	s.Require().NoError(err)

	_, err = s.Keeper().UpdateValidatorPubKey(ctx)
	s.Require().NoError(err)

	newConsAddr := sdk.GetConsAddress(newPubKey)
	newSigningInfo, found := s.App.SlashingKeeper.GetValidatorSigningInfo(ctx, newConsAddr)
	s.Require().True(found, "new consensus address must have signing info")
	s.Require().Equal(newConsAddr.String(), newSigningInfo.Address)
	s.Require().Equal(oldSigningInfo.StartHeight, newSigningInfo.StartHeight)
	s.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)
	s.Require().Equal(oldSigningInfo.JailedUntil, newSigningInfo.JailedUntil)
	s.Require().Equal(oldSigningInfo.Tombstoned, newSigningInfo.Tombstoned)
	s.Require().Equal(oldSigningInfo.MissedBlocksCounter, newSigningInfo.MissedBlocksCounter)
	s.Require().True(s.App.SlashingKeeper.GetValidatorMissedBlockBitArray(ctx, newConsAddr, 0))
	s.Require().True(s.App.SlashingKeeper.GetValidatorMissedBlockBitArray(ctx, newConsAddr, 7))
	s.Require().False(s.App.SlashingKeeper.GetValidatorMissedBlockBitArray(ctx, newConsAddr, 8))
}
