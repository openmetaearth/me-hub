package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestSendToMeClaimRejectsMixedChainNameVotes() {
	s.setupBondedRelayerSetForAuditTest()

	tokenContract := s.PubKeyToExternalAddr(s.externalPris[0].PublicKey)
	receiver := s.relayerAddrs[0]
	s.Keeper().SetBridgeToken(s.Ctx, &types.BridgeToken{
		ContractAddress: tokenContract,
		Denom:           "uusdt",
		Name:            "Tether USD",
		Symbol:          "USDT",
		Decimal:         18,
		Supply:          sdkmath.ZeroInt(),
	})

	eventNonce := s.Keeper().GetLastObservedEventNonce(s.Ctx) + 1
	blockHeight := uint64(100)
	amount := sdkmath.NewIntWithDecimal(1, 18)
	sender := s.PubKeyToExternalAddr(s.externalPris[1].PublicKey)

	var honestClaim *types.MsgSendToMeClaim
	for i := 0; i < 6; i++ {
		honestClaim = &types.MsgSendToMeClaim{
			EventNonce:     eventNonce,
			BlockHeight:    blockHeight,
			TokenContract:  tokenContract,
			Amount:         amount,
			Sender:         sender,
			Receiver:       receiver.String(),
			RelayerAddress: s.relayerAddrs[i].String(),
			ChainName:      s.chainName,
		}
		_, err := s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(s.Ctx), honestClaim)
		s.Require().NoError(err)
	}

	att := s.Keeper().GetAttestation(s.Ctx, eventNonce, honestClaim.ClaimHash())
	s.Require().NotNil(att)
	s.Require().Len(att.Votes, 6)
	s.Require().False(att.Observed)

	wrongChainClaim := *honestClaim
	wrongChainClaim.RelayerAddress = s.relayerAddrs[6].String()
	wrongChainClaim.ChainName = types.ModuleName

	_, err := s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(s.Ctx), &wrongChainClaim)
	s.Require().Error(err)

	balance := s.App.BankKeeper.GetBalance(s.Ctx, receiver, "uusdt")
	s.Require().True(balance.Amount.IsZero())
	s.Require().EqualValues(eventNonce-1, s.Keeper().GetLastObservedEventNonce(s.Ctx))

	att = s.Keeper().GetAttestation(s.Ctx, eventNonce, honestClaim.ClaimHash())
	s.Require().NotNil(att)
	s.Require().Len(att.Votes, 6)
	s.Require().False(att.Observed)
}
