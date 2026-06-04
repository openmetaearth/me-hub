package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestAttestationThresholdDoesNotFloorBelowConfiguredRatio() {
	k := s.Keeper()

	highPowerRelayer := s.relayerAddrs[0]
	lowPowerRelayer := s.relayerAddrs[1]
	k.SetRelayer(s.Ctx, highPowerRelayer, types.Relayer{
		RelayerAddress:  highPowerRelayer.String(),
		ExternalAddress: helpers.GenExternalAddr(s.chainName),
		DelegateAmount:  sdkmath.NewInt(199).Mul(sdk.DefaultPowerReduction),
		StartHeight:     s.Ctx.BlockHeight(),
		Online:          true,
	})
	k.SetRelayer(s.Ctx, lowPowerRelayer, types.Relayer{
		RelayerAddress:  lowPowerRelayer.String(),
		ExternalAddress: helpers.GenExternalAddr(s.chainName),
		DelegateAmount:  sdkmath.NewInt(100).Mul(sdk.DefaultPowerReduction),
		StartHeight:     s.Ctx.BlockHeight(),
		Online:          true,
	})
	k.SetLastTotalPower(s.Ctx)
	s.Require().Equal(sdkmath.NewInt(299), k.GetLastTotalPower(s.Ctx))

	tokenContract := helpers.GenHexAddress().String()
	bridgeToken := &types.BridgeToken{
		ContractAddress: tokenContract,
		Denom:           "rounding",
		Name:            "Rounding Token",
		Symbol:          "ROUNDING",
		Decimal:         6,
		Supply:          sdkmath.ZeroInt(),
	}
	k.SetBridgeToken(s.Ctx, bridgeToken)

	receiver := s.relayerAddrs[2]
	claim := &types.MsgSendToMeClaim{
		EventNonce:     1,
		BlockHeight:    1,
		TokenContract:  tokenContract,
		Amount:         sdkmath.NewInt(1234),
		Sender:         helpers.GenExternalAddr(s.chainName),
		Receiver:       receiver.String(),
		RelayerAddress: highPowerRelayer.String(),
		ChainName:      s.chainName,
	}

	_, err := s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(s.Ctx), claim)
	s.Require().NoError(err)

	attestation := k.GetAttestation(s.Ctx, claim.EventNonce, claim.ClaimHash())
	s.Require().NotNil(attestation)
	s.Require().False(attestation.Observed, "199/299 power is below the configured 6666/10000 threshold")
	s.Require().Equal(sdkmath.ZeroInt(), s.App.BankKeeper.GetBalance(s.Ctx, receiver, bridgeToken.Denom).Amount)

	storedBridgeToken, err := k.GetBridgeTokenByContract(s.Ctx, tokenContract)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.ZeroInt(), storedBridgeToken.Supply)
}
