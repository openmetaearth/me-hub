package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestSlashRelayerOfflineBehavior() {
	k := s.Keeper()

	// 1. Setup initial set of 10 bonded relayers
	s.setupBondedRelayerSetForAuditTest()

	// Verify initial total power is set and greater than zero
	initialPower := k.GetLastTotalPower(s.Ctx)
	s.Require().True(initialPower.GT(sdkmath.ZeroInt()), "initial total power must be greater than zero")

	// Get one relayer to slash
	slashedRelayerAddr := s.relayerAddrs[0]
	relayer, found := k.GetRelayer(s.Ctx, slashedRelayerAddr)
	s.Require().True(found, "relayer must be found")
	s.Require().True(relayer.Online, "relayer must be online initially")

	relayerPower := relayer.GetPower()

	// 2. Slash the relayer until they go offline
	maxSlash := int(k.GetParams(s.Ctx).MaxSlashTimes)
	for i := 0; i < maxSlash+5; i++ {
		err := k.SlashRelayer(s.Ctx, slashedRelayerAddr.String())
		s.Require().NoError(err)

		relayer, found = k.GetRelayer(s.Ctx, slashedRelayerAddr)
		s.Require().True(found)
		if !relayer.Online {
			break
		}
	}
	s.Require().False(relayer.Online, "relayer must be offline after slashing")

	// 3. Verify LastTotalPower is recalculated and reduced by the offline relayer's power
	postSlashPower := k.GetLastTotalPower(s.Ctx)
	expectedPower := initialPower.Sub(relayerPower)
	s.Require().Equal(expectedPower.String(), postSlashPower.String(), "LastTotalPower must be reduced by the offline relayer's power")

	// 4. Verify that the offline relayer cannot submit claims via Attest
	tokenContract := "0x0000000000000000000000000000000000000009"
	k.SetBridgeToken(s.Ctx, &types.BridgeToken{
		ContractAddress: tokenContract,
		Denom:           "test",
		Name:            "Test Token",
	})

	claim := &types.MsgSendToMeClaim{
		EventNonce:     k.GetLastEventNonceByRelayer(s.Ctx, slashedRelayerAddr) + 1,
		BlockHeight:    1,
		TokenContract:  tokenContract,
		Amount:         sdkmath.NewInt(100),
		Sender:         s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
		Receiver:       s.relayerAddrs[1].String(),
		RelayerAddress: slashedRelayerAddr.String(),
		ChainName:      s.chainName,
	}

	_, err := k.Attest(s.Ctx, slashedRelayerAddr, claim)
	s.Require().Error(err, "Attest must reject claims from offline relayer")
	s.Require().ErrorIs(err, types.ErrRelayerNotOnLine, "error must be ErrRelayerNotOnLine")
}
