package keeper_test

import (
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestAttestationHandlerRejectsOversizedBridgeTokenDecimals() {
	k := s.Keeper()
	tokenContract := helpers.GenerateAddress().Hex()

	err := k.AttestationHandler(s.Ctx, &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1,
		TokenContract:  tokenContract,
		Name:           "Oversized",
		Symbol:         "USDT",
		Decimals:       types.MaxBridgeTokenDecimals + 1,
		RelayerAddress: s.relayerAddrs[0].String(),
		ChainName:      s.chainName,
	})

	s.Require().ErrorContains(err, "bridge token decimals")
	s.Require().False(k.HasBridgeToken(s.Ctx, tokenContract))
}
