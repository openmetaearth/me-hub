package keeper_test

// Audit regression tests for findings from the Kimi 02_gravity_mint audit
// that require a live keeper context.
//
// These are "red" tests: they encode the CORRECT behavior from the audit.
// On current buggy code they FAIL; once the bug is fixed they PASS. Do not
// weaken them to make them pass - fix the underlying code instead.
//
// Coverage:
//   GRAV-001  Attest must reject historical-nonce claims (no backward rewind)
//   GRAV-004  Failed AttestationHandler must not advance lastObservedEventNonce

import (
	"encoding/hex"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
)

// ---------------------------------------------------------------------------
// GRAV-001: Historical-nonce rewind via Attest
// ---------------------------------------------------------------------------
//
// Attest at attestation.go:30-37 only rejects a claim when
//   nonce != expectedNonce && nonce > lastObservedNonce
// This silently accepts any historical nonce in (lastEventNonce,
// lastObservedNonce], advances the relayer's personal lastEventNonce
// backwards to that historical value, and corrupts downstream slashing
// and continuity accounting.
//
// Correct behavior: reject any claim whose nonce is not strictly the next
// expected nonce (and certainly any nonce <= lastObservedNonce).

func (s *KeeperTestSuite) TestGrav001_AttestRejectsHistoricalNonceRewind() {
	k := s.Keeper()
	relayer := s.relayerAddrs[0]

	// Prime state: global observed nonce = 10, relayer's personal nonce = 10.
	// A well-behaved relayer would submit nonce 11 next.
	k.SetLastObservedEventNonce(s.Ctx, 10)
	k.SetLastEventNonceByRelayer(s.Ctx, relayer, 10)

	// Construct a historical claim at nonce 5 - already observed globally.
	historicalClaim := &types.MsgSendToMeClaim{
		EventNonce:     5,
		BlockHeight:    1,
		TokenContract:  "0x0000000000000000000000000000000000000001",
		Amount:         sdkmath.NewInt(1000),
		Sender:         "0x0000000000000000000000000000000000000002",
		Receiver:       s.relayerAddrs[1].String(),
		RelayerAddress: relayer.String(),
		ChainName:      s.chainName,
	}

	_, err := k.Attest(s.Ctx, relayer, historicalClaim)
	s.Require().Errorf(err,
		"GRAV-001: Attest must reject a historical claim at nonce 5 when lastObservedNonce is 10")
	s.Require().ErrorIs(err, types.ErrNonContinuousEventNonce,
		"GRAV-001: rejection should use ErrNonContinuousEventNonce, got %v", err)

	// Regardless of error path, the relayer's own lastEventNonce must NOT
	// regress. The exploit is specifically about walking this value backwards.
	s.Require().EqualValues(uint64(10), k.GetLastEventNonceByRelayer(s.Ctx, relayer),
		"GRAV-001: relayer lastEventNonce must not regress after rejected historical claim")
}

// A second variant: nonce equal to lastObservedNonce exactly. The existing
// explicit check at attestation.go:31-33 handles the case where
// nonce == lastEventNonce, but a relayer whose own nonce is ahead of the
// global observed can still submit nonce == lastObservedNonce. Verify that
// is also rejected.
func (s *KeeperTestSuite) TestGrav001_AttestRejectsNonceEqualToLastObserved() {
	k := s.Keeper()
	relayer := s.relayerAddrs[0]

	// Global observed = 10, relayer's personal = 12 (ahead of global).
	// This is an unusual but reachable state if relayer participation drifts.
	k.SetLastObservedEventNonce(s.Ctx, 10)
	k.SetLastEventNonceByRelayer(s.Ctx, relayer, 12)

	historicalClaim := &types.MsgSendToMeClaim{
		EventNonce:     10, // == lastObservedNonce, != lastEventNonce
		BlockHeight:    1,
		TokenContract:  "0x0000000000000000000000000000000000000001",
		Amount:         sdkmath.NewInt(1000),
		Sender:         "0x0000000000000000000000000000000000000002",
		Receiver:       s.relayerAddrs[1].String(),
		RelayerAddress: relayer.String(),
		ChainName:      s.chainName,
	}

	_, err := k.Attest(s.Ctx, relayer, historicalClaim)
	s.Require().Errorf(err,
		"GRAV-001: Attest must reject a claim at nonce == lastObservedNonce when != lastEventNonce")

	s.Require().EqualValues(uint64(12), k.GetLastEventNonceByRelayer(s.Ctx, relayer),
		"GRAV-001: relayer lastEventNonce must not regress")
}

// ---------------------------------------------------------------------------
// GRAV-004: Failed AttestationHandler must not lock the nonce
// ---------------------------------------------------------------------------
//
// TryAttestation calls SetLastObservedEventNonce and sets att.Observed = true
// BEFORE calling processAttestation. If the handler returns an error, the
// error is logged and returned upward, but the state changes in the parent
// context (nonce advance, Observed flag, SetAttestation) remain committed.
// This means a transient handler failure permanently locks that nonce,
// bricking the bridge for all future events via the continuity requirement.
//
// This test triggers a natural handler failure by submitting
// MsgSendToMeClaim for a TokenContract that has never been registered.
// attestation_handler.go:19-22 returns ErrInvalid immediately.
//
// Correct behavior: lastObservedEventNonce stays at its pre-call value,
// and the attestation's Observed flag remains false, so a retry after
// the underlying condition is fixed can still succeed.

func (s *KeeperTestSuite) TestGrav004_FailedAttestationMustNotLockNonce() {
	k := s.Keeper()

	// ---- Setup: bond all relayers, confirm and observe the initial relayer set. ----
	s.setupBondedRelayerSetForAuditTest()

	// After RelayerSetUpdateClaim, lastObservedEventNonce == 1. Capture it.
	preObserved := k.GetLastObservedEventNonce(s.Ctx)
	s.Require().EqualValues(uint64(1), preObserved,
		"test precondition: relayer set update claim should be observed at nonce 1")

	// Use a bridge-token contract that has NEVER been registered. The handler
	// path for MsgSendToMeClaim calls GetBridgeTokenByContract first and
	// returns ErrInvalid before touching any state.
	unknownContract := "0xdeadbeefcafebabe00000000000000000000dead"
	s.Require().False(k.HasBridgeToken(s.Ctx, unknownContract),
		"test precondition: the attacker contract must not be a registered bridge token")

	// Every relayer submits SendToMeClaim for this nonexistent token, at
	// their next expected nonce (2). The msg_server swallows the handler
	// error and returns OK; we assert on the persisted state, not on err.
	receiver := s.relayerAddrs[0]
	sender := s.PubKeyToExternalAddr(s.externalPris[0].PublicKey)
	for i, relayer := range s.relayerAddrs {
		claim := &types.MsgSendToMeClaim{
			EventNonce:     k.GetLastEventNonceByRelayer(s.Ctx, relayer) + 1,
			BlockHeight:    1,
			TokenContract:  unknownContract,
			Amount:         sdkmath.NewIntWithDecimal(1, 6),
			Sender:         sender,
			Receiver:       receiver.String(),
			RelayerAddress: relayer.String(),
			ChainName:      s.chainName,
		}
		_, err := s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(s.Ctx), claim)
		// msg_server wraps handler errors; don't require NoError, but the
		// tx must not have changed global observed state in a way that locks
		// the nonce regardless.
		_ = i
		_ = err
	}

	// Assertion 1: lastObservedEventNonce must NOT have advanced past the
	// pre-call value. If it has, the failing mint has permanently locked
	// the nonce and all subsequent legitimate mints will fail continuity.
	postObserved := k.GetLastObservedEventNonce(s.Ctx)
	s.Require().Equalf(preObserved, postObserved,
		"GRAV-004: lastObservedEventNonce advanced from %d to %d even though "+
			"AttestationHandler returned an error (bridge token does not exist). "+
			"This permanently locks the nonce and bricks the bridge for retries.",
		preObserved, postObserved)

	// Assertion 2: the attestation record for the failed claim should NOT
	// be marked Observed, so a future retry (after the underlying condition
	// is fixed) can still reach quorum. We recreate the claim hash here.
	probeClaim := &types.MsgSendToMeClaim{
		EventNonce:     preObserved + 1,
		BlockHeight:    1,
		TokenContract:  unknownContract,
		Amount:         sdkmath.NewIntWithDecimal(1, 6),
		Sender:         sender,
		Receiver:       receiver.String(),
		RelayerAddress: s.relayerAddrs[0].String(),
		ChainName:      s.chainName,
	}
	att := k.GetAttestation(s.Ctx, probeClaim.GetEventNonce(), probeClaim.ClaimHash())
	if att != nil {
		s.Require().Falsef(att.Observed,
			"GRAV-004: attestation at nonce %d is marked Observed even though its "+
				"handler failed; retries are now impossible",
			probeClaim.GetEventNonce())
	}
}

// ---------------------------------------------------------------------------
// Helper: bond relayers, confirm and observe the initial relayer set.
// Mirrors the opening of TestRequestBatchBaseFee in msg_server_test.go, but
// stops before BridgeTokenClaim so the state is a clean "bridge initialized,
// no tokens registered" baseline for audit tests.
// ---------------------------------------------------------------------------

func (s *KeeperTestSuite) setupBondedRelayerSetForAuditTest() {
	totalPower := sdkmath.ZeroInt()
	delegateAmounts := make([]sdkmath.Int, 0, len(s.relayerAddrs))
	for i, relayer := range s.relayerAddrs {
		msg := &types.MsgBondedRelayer{
			RelayerAddress:  relayer.String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(int64(1e9))),
			ChainName:       s.chainName,
		}
		delegateAmounts = append(delegateAmounts, msg.DelegateAmount.Amount)
		totalPower = totalPower.Add(msg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction))
		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msg)
		s.Require().NoError(err)
	}
	s.Keeper().EndBlocker(s.Ctx)

	var externalRelayerMembers types.BridgeValidators
	for i, key := range s.externalPris {
		power := delegateAmounts[i].Quo(sdk.DefaultPowerReduction).
			MulRaw(int64(types.PowerBase)).Quo(totalPower)
		externalRelayerMembers = append(externalRelayerMembers, types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: s.PubKeyToExternalAddr(key.PublicKey),
		})
	}
	sort.Sort(externalRelayerMembers)

	latestRelayerSetNonce := s.Keeper().GetLastRelayerSetNonce(s.Ctx)
	s.Require().EqualValues(1, latestRelayerSetNonce)
	nonce1RelayerSet := s.Keeper().GetRelayerSet(s.Ctx, 1)
	gravityId := s.Keeper().GetGravityID(s.Ctx)
	checkpoint, err := nonce1RelayerSet.GetCheckpoint(gravityId)
	s.Require().NoError(err)
	if trontypes.ModuleName == s.chainName {
		checkpoint, err = trontypes.GetCheckpointRelayerSet(nonce1RelayerSet, gravityId)
		s.Require().NoError(err)
	}

	for i := range s.relayerAddrs {
		external2Signature, err := types.NewEthereumSignature(checkpoint, s.externalPris[i])
		s.Require().NoError(err)
		if trontypes.ModuleName == s.chainName {
			external2Signature, err = trontypes.NewTronSignature(checkpoint, s.externalPris[i])
			s.Require().NoError(err)
		}
		msg := &types.MsgRelayerSetConfirm{
			Nonce:           nonce1RelayerSet.Nonce,
			RelayerAddress:  s.relayerAddrs[i].String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			Signature:       hex.EncodeToString(external2Signature),
			ChainName:       s.chainName,
		}
		_, err = s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), msg)
		s.Require().NoError(err)
	}

	for i := range s.relayerAddrs {
		msg := &types.MsgRelayerSetUpdateClaim{
			EventNonce:      1,
			BlockHeight:     1,
			RelayerSetNonce: 1,
			Members:         externalRelayerMembers,
			RelayerAddress:  s.relayerAddrs[i].String(),
			ChainName:       s.chainName,
		}
		_, err := s.MsgServer().RelayerSetUpdateClaim(sdk.WrapSDKContext(s.Ctx), msg)
		s.Require().NoError(err)
	}
	s.Keeper().EndBlocker(s.Ctx)
}
