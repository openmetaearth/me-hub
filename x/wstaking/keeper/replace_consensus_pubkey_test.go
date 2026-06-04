package keeper_test

import (
	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestReplaceConsensusPubKeyKeepsOldConsAddrForHistoricalEvidence() {
	validator := s.meEarthValidator
	oldConsAddr, err := validator.GetConsAddr()
	s.Require().NoError(err)
	s.Require().NoError(s.Keeper().SetValidatorByConsAddr(s.Ctx, validator))

	_, found := s.Keeper().GetValidatorByConsAddr(s.Ctx, oldConsAddr)
	s.Require().True(found, "test precondition: old consensus address should resolve before replacement")

	newPubKey := ed25519.GenPrivKey().PubKey()
	newConsAddr := sdk.GetConsAddress(newPubKey)
	msg, err := types.NewMsgReplaceConsensusPubKeyRequest(s.Dao.GlobalDao, validator.OperatorAddress, newPubKey, 1)
	s.Require().NoError(err)

	_, err = s.msgServer.ReplaceConsensusPubKey(s.Ctx, msg)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithBlockHeight(1)
	updates := s.Keeper().BlockValidatorUpdates(s.Ctx)
	s.Require().Len(updates, 2, "replacement block should remove old key power and add new key power")

	_, found = s.Keeper().GetValidatorByConsAddr(s.Ctx, oldConsAddr)
	s.Require().True(found, "old consensus address should still resolve immediately after replacement")
	_, found = s.Keeper().GetValidatorByConsAddr(s.Ctx, newConsAddr)
	s.Require().True(found, "new consensus address should resolve after replacement")

	s.Ctx = s.Ctx.WithBlockHeight(2)
	_ = s.Keeper().BlockValidatorUpdates(s.Ctx)
	_, found = s.Keeper().GetValidatorByConsAddr(s.Ctx, oldConsAddr)
	s.Require().True(found, "old consensus address should still resolve during replacement cleanup delay")

	s.Ctx = s.Ctx.WithBlockHeight(3)
	_ = s.Keeper().BlockValidatorUpdates(s.Ctx)
	_, found = s.Keeper().GetValidatorByConsAddr(s.Ctx, oldConsAddr)
	s.Require().True(found, "historical evidence for the old consensus key still needs old consAddr -> validator lookup")
	s.Require().False(s.Keeper().IsHasReplaceConsensusPubKey(s.Ctx), "replacement state should be cleaned after the delay")
}
