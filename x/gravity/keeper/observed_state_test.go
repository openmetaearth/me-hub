package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestTryAttestationRejectsExternalBlockHeightRegression() {
	k := s.Keeper()
	relayer := s.relayerAddrs[0]
	k.SetRelayer(s.Ctx, relayer, types.Relayer{
		RelayerAddress: relayer.String(),
		DelegateAmount: sdk.DefaultPowerReduction,
		Online:         true,
	})
	k.SetLastTotalPower(s.Ctx)
	k.SetLastObservedEventNonce(s.Ctx, 1)
	k.SetLastObservedBlockHeight(s.Ctx, 1000, uint64(s.Ctx.BlockHeight()))

	claim := &types.MsgBridgeTokenClaim{
		EventNonce:     2,
		BlockHeight:    1,
		TokenContract:  helpers.GenerateAddress().Hex(),
		Name:           "LowHeight",
		Symbol:         "LOW",
		Decimals:       6,
		RelayerAddress: relayer.String(),
		ChainName:      s.chainName,
	}
	att := &types.Attestation{
		Observed: false,
		Votes:    []string{relayer.String()},
	}

	k.TryAttestation(s.Ctx, att, claim)

	lastObservedHeight := k.GetLastObservedBlockHeight(s.Ctx)
	s.Require().EqualValues(1000, lastObservedHeight.ExternalBlockHeight)
	s.Require().False(k.HasBridgeToken(s.Ctx, claim.TokenContract))

	storedAtt := k.GetAttestation(s.Ctx, claim.GetEventNonce(), claim.ClaimHash())
	s.Require().NotNil(storedAtt)
	s.Require().True(storedAtt.Observed)
}

func (s *KeeperTestSuite) TestAttestationHandlerRejectsRelayerSetNonceRegression() {
	k := s.Keeper()
	members := types.BridgeValidators{
		{
			Power:           types.PowerBase,
			ExternalAddress: helpers.GenerateAddress().Hex(),
		},
	}
	nonce1Set := &types.RelayerSet{
		Nonce:   1,
		Height:  10,
		Members: members,
	}
	nonce2Set := &types.RelayerSet{
		Nonce:   2,
		Height:  20,
		Members: members,
	}
	k.StoreRelayerSet(s.Ctx, nonce1Set)
	k.StoreRelayerSet(s.Ctx, nonce2Set)
	k.SetLastObservedRelayerSet(s.Ctx, nonce2Set)

	err := k.AttestationHandler(s.Ctx, &types.MsgRelayerSetUpdateClaim{
		EventNonce:      3,
		BlockHeight:     30,
		RelayerSetNonce: 1,
		Members:         members,
		RelayerAddress:  s.relayerAddrs[0].String(),
		ChainName:       s.chainName,
	})

	s.Require().ErrorContains(err, "relayer set nonce regression")
	s.Require().EqualValues(2, k.GetLastObservedRelayerSet(s.Ctx).Nonce)
}

func (s *KeeperTestSuite) TestAttestationHandlerRejectsZeroRelayerSetNonce() {
	k := s.Keeper()
	err := k.AttestationHandler(s.Ctx, &types.MsgRelayerSetUpdateClaim{
		EventNonce:      2,
		BlockHeight:     20,
		RelayerSetNonce: 0,
		Members: types.BridgeValidators{
			{
				Power:           1,
				ExternalAddress: helpers.GenerateAddress().Hex(),
			},
		},
		RelayerAddress: s.relayerAddrs[0].String(),
		ChainName:      s.chainName,
	})

	s.Require().ErrorContains(err, "zero relayer set nonce")
}
