package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestAttestationThresholdDoesNotFloorBelowConfiguredRatio() {
	k := s.Keeper()

	relayerPower := []struct {
		addr  string
		power int64
	}{
		{addr: s.relayerAddrs[0].String(), power: 199},
		{addr: s.relayerAddrs[1].String(), power: 100},
	}

	for _, rp := range relayerPower {
		k.SetRelayer(s.Ctx, sdk.MustAccAddressFromBech32(rp.addr), gravitytypes.Relayer{
			RelayerAddress: rp.addr,
			DelegateAmount: sdkmath.NewInt(rp.power).Mul(sdk.DefaultPowerReduction),
			Online:         true,
		})
	}
	k.SetLastTotalPower(s.Ctx)
	s.Require().EqualValues(299, k.GetLastTotalPower(s.Ctx).Int64())

	bridgeHolder, _ := s.NewAccount()
	receiver, _ := s.NewAccount()

	bridgeToken := s.NewBridgeToken(bridgeHolder, sdk.NewCoin("rounding", sdkmath.NewInt(1000)))
	s.Require().EqualValues(sdkmath.NewInt(1000), bridgeToken.Supply)

	claim := &gravitytypes.MsgSendToMeClaim{
		EventNonce:     1,
		BlockHeight:    1,
		TokenContract:  bridgeToken.ContractAddress,
		Amount:         sdkmath.NewInt(1234),
		Sender:         helpers.GenExternalAddr(s.chainName),
		Receiver:       receiver.String(),
		RelayerAddress: s.relayerAddrs[0].String(),
		ChainName:      s.chainName,
	}

	_, err := s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(s.Ctx), claim)
	s.Require().NoError(err)

	att := k.GetAttestation(s.Ctx, claim.GetEventNonce(), claim.ClaimHash())
	s.Require().NotNil(att)
	s.Require().False(att.Observed)
	s.Require().EqualValues(uint64(0), k.GetLastObservedEventNonce(s.Ctx))

	receiverBalance := s.App.BankKeeper.GetBalance(s.Ctx, receiver, bridgeToken.Denom)
	s.Require().True(receiverBalance.Amount.IsZero())

	storedBridgeToken, err := k.GetBridgeTokenByContract(s.Ctx, bridgeToken.ContractAddress)
	s.Require().NoError(err)
	s.Require().EqualValues(bridgeToken.Supply, storedBridgeToken.Supply)
}
