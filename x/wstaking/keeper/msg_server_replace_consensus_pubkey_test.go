package keeper_test

import (
	"math"

	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestReplaceConsensusPubKeySchedulesFutureHeight() {
	s.SetupTest()
	s.Ctx = s.Ctx.WithBlockHeight(10)

	msg, err := types.NewMsgReplaceConsensusPubKeyRequest(
		s.Dao.GlobalDao,
		s.meEarthValidator.OperatorAddress,
		ed25519.GenPrivKey().PubKey(),
		10,
	)
	s.Require().NoError(err)
	s.Require().NoError(msg.ValidateBasic())

	_, err = s.msgServer.ReplaceConsensusPubKey(s.Ctx, msg)
	s.Require().NoError(err)

	pending, err := s.Keeper().GetReplaceConsensusPubKeyInfo(s.Ctx)
	s.Require().NoError(err)
	s.Require().NotNil(pending)
	s.Require().Equal(int64(20), pending.UpdateAtHeight)
}

func (s *KeeperTestSuite) TestReplaceConsensusPubKeyRejectsOverflowingDelay() {
	s.SetupTest()
	s.Ctx = s.Ctx.WithBlockHeight(10)
	overflowingDelay := int64(math.MaxInt64) - s.Ctx.BlockHeight() + 1

	msg, err := types.NewMsgReplaceConsensusPubKeyRequest(
		s.Dao.GlobalDao,
		s.meEarthValidator.OperatorAddress,
		ed25519.GenPrivKey().PubKey(),
		overflowingDelay,
	)
	s.Require().NoError(err)
	s.Require().NoError(msg.ValidateBasic())

	_, err = s.msgServer.ReplaceConsensusPubKey(s.Ctx, msg)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "update height overflows")

	pending, err := s.Keeper().GetReplaceConsensusPubKeyInfo(s.Ctx)
	s.Require().NoError(err)
	s.Require().Nil(pending)
}
