package keeper_test

import (
	"encoding/hex"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestNewRelayerVotesDoNotCountBeforeExternalSetObservesIt() {
	k := s.Keeper()

	// 1. Setup initial set of 10 bonded relayers (total power 10 * 1e9 = 10e9)
	s.setupBondedRelayerSetForAuditTest()

	// Register a bridge token
	tokenContract := "0x0000000000000000000000000000000000000009"
	k.SetBridgeToken(s.Ctx, &types.BridgeToken{
		ContractAddress: tokenContract,
		Denom:           "test",
		Name:            "Test Token",
	})

	// 2. Bond an 11th relayer locally (power = 2e9)
	newRelayers := s.NewAccounts(1)
	newRelayer := newRelayers[0]
	newRelayerKey := helpers.CreateMultiECDSA(1)[0]

	// Fund new relayer
	err := s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, newRelayer, sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 10000000000)})
	s.Require().NoError(err)

	msg := &types.MsgBondedRelayer{
		RelayerAddress:  newRelayer.String(),
		ExternalAddress: s.PubKeyToExternalAddr(newRelayerKey.PublicKey),
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(int64(2e9))),
		ChainName:       s.chainName,
	}
	_, err = s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msg)
	s.Require().NoError(err)
	k.EndBlocker(s.Ctx)

	// Confirm that the new relayer's power is indeed set, but the last observed relayer set has not been updated yet
	lastObserved := k.GetLastObservedRelayerSet(s.Ctx)
	s.Require().NotNil(lastObserved)
	s.Require().Equal(10, len(lastObserved.Members), "observed relayer set must still have 10 members")

	// 3. Construct a claim (nonce 2)
	receiver := s.NewAccounts(1)[0]
	senderExternal := s.PubKeyToExternalAddr(s.externalPris[0].PublicKey)
	claimAmount := sdkmath.NewInt(12345)

	buildClaimForRelayer := func(relayerAddr sdk.AccAddress) *types.MsgSendToMeClaim {
		return &types.MsgSendToMeClaim{
			EventNonce:     k.GetLastEventNonceByRelayer(s.Ctx, relayerAddr) + 1,
			BlockHeight:    1,
			TokenContract:  tokenContract,
			Amount:         claimAmount,
			Sender:         senderExternal,
			Receiver:       receiver.String(),
			RelayerAddress: relayerAddr.String(),
			ChainName:      s.chainName,
		}
	}

	// 4. Have 6 of the old relayers vote (6 * 1e9 = 6e9 power)
	for i := 0; i < 6; i++ {
		claim := buildClaimForRelayer(s.relayerAddrs[i])
		_, err = s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(s.Ctx), claim)
		s.Require().NoError(err)
	}

	// At this point, vote power is 6e9.
	// Old threshold: 10e9 * 2 / 3 = 6.67e9. Not met.
	// New threshold (if new relayer is counted): 12e9 * 2 / 3 = 8e9.
	// Check that the claim is not yet observed.
	testClaim := buildClaimForRelayer(s.relayerAddrs[0])
	att := k.GetAttestation(s.Ctx, testClaim.GetEventNonce(), testClaim.ClaimHash())
	if att != nil {
		s.Require().False(att.Observed, "claim must not be observed yet")
	}

	// 5. Have the new (11th) relayer vote (adds 2e9 power)
	newRelayerClaim := buildClaimForRelayer(newRelayer)
	_, err = s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(s.Ctx), newRelayerClaim)
	s.Require().NoError(err)

	// If the new relayer's vote is counted: 6e9 + 2e9 = 8e9 power.
	// If it is incorrectly counted, total vote power is 8e9 out of 12e9, which meets the 2/3 threshold, and the claim would be observed.
	// If it is correctly ignored (since it's not in the LastObservedRelayerSet), the vote power is just 6e9, which is below the threshold, and the claim must NOT be observed.
	att = k.GetAttestation(s.Ctx, testClaim.GetEventNonce(), testClaim.ClaimHash())
	if att != nil {
		s.Require().False(att.Observed, "new relayer vote must not mint before the external bridge observes that relayer set")
	}

	// Assert that no vouchers were minted yet
	balance := s.App.BankKeeper.GetBalance(s.Ctx, receiver, "test")
	s.Require().Equal(sdkmath.ZeroInt(), balance.Amount, "new relayer vote must not mint before the external bridge observes that relayer set")
}
