package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestGetLastEventNonceByRelayer_EmptyStoreUsesLastObserved() {
	k := s.Keeper()
	newRelayer := s.relayerAddrs[0]

	s.Require().EqualValues(uint64(0), k.GetLastObservedEventNonce(s.Ctx))
	s.Require().EqualValues(uint64(0), k.GetLastEventNonceByRelayer(s.Ctx, newRelayer),
		"empty store with no observed events must start at 0")

	k.SetLastObservedEventNonce(s.Ctx, 266)
	s.Require().EqualValues(uint64(266), k.GetLastEventNonceByRelayer(s.Ctx, newRelayer),
		"new relayer must be treated as caught up to lastObserved, not lastObserved-1")

	k.SetLastEventNonceByRelayer(s.Ctx, newRelayer, 100)
	s.Require().EqualValues(uint64(100), k.GetLastEventNonceByRelayer(s.Ctx, newRelayer),
		"stored personal nonce must win over lastObserved")
}

func (s *KeeperTestSuite) TestAttest_NewRelayerSkipsAlreadyObservedNonce() {
	k := s.Keeper()
	relayer := s.relayerAddrs[0]
	k.SetLastObservedEventNonce(s.Ctx, 266)

	s.Require().EqualValues(uint64(266), k.GetLastEventNonceByRelayer(s.Ctx, relayer))

	observedClaim := sendToMeClaim(relayer, s.chainName, 266)
	_, err := k.Attest(s.Ctx, relayer, observedClaim)
	s.Require().ErrorIs(err, types.ErrNonContinuousEventNonce,
		"new relayer must not re-claim the already-observed tip")

	nextClaim := sendToMeClaim(relayer, s.chainName, 267)
	_, err = k.Attest(s.Ctx, relayer, nextClaim)
	s.Require().NoError(err)
	s.Require().EqualValues(uint64(267), k.GetLastEventNonceByRelayer(s.Ctx, relayer))
}

func sendToMeClaim(relayer sdk.AccAddress, chainName string, nonce uint64) *types.MsgSendToMeClaim {
	return &types.MsgSendToMeClaim{
		EventNonce:     nonce,
		BlockHeight:    1,
		TokenContract:  "0x0000000000000000000000000000000000000001",
		Amount:         sdkmath.NewInt(1000),
		Sender:         "0x0000000000000000000000000000000000000002",
		Receiver:       relayer.String(),
		RelayerAddress: relayer.String(),
		ChainName:      chainName,
	}
}
