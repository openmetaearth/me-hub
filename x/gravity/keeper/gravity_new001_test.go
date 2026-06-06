package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// TestGravityNew001_FixSlashTimesReset verifies the fix for:
// AddDelegate unconditionally resetting SlashTimes to 0.
//
// Before fix: AddDelegate always reset SlashTimes=0, allowing infinite bypass.
// After fix:  SlashTimes only resets when relayer was fully offline (>= MaxSlashTimes).
func (s *KeeperTestSuite) TestGravityNew001_FixSlashTimesResetViaAddDelegate() {
	k := s.Keeper()
	relayer := s.relayerAddrs[0]

	// Set MaxSlashTimes to 3
	params2 := k.GetParams(s.Ctx)
	params2.MaxSlashTimes = 3
	k.SetParams(s.Ctx, &params2)

	// Bond relayer
	minDelegate := k.GetGravityMinDelegate(s.Ctx)
	bondMsg := &types.MsgBondedRelayer{
		RelayerAddress:  relayer.String(),
		ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, minDelegate),
		ChainName:       s.chainName,
	}
	_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), bondMsg)
	s.Require().NoError(err)

	// === CYCLE 1: Slash 2 times (below MaxSlashTimes), AddDelegate should NOT reset ===
	for i := uint64(0); i < 2; i++ {
		s.Require().NoError(k.SlashRelayer(s.Ctx, relayer.String()))
	}
	relayerData, _ := k.GetRelayer(s.Ctx, relayer)
	s.Require().True(relayerData.Online, "cycle 1: still online after 2 slashes")
	s.Require().Equal(int64(2), relayerData.SlashTimes, "cycle 1: SlashTimes=2")

	// AddDelegate — fix should NOT reset SlashTimes (2 < MaxSlashTimes=3)
	addMsg := &types.MsgAddDelegate{
		RelayerAddress: relayer.String(),
		Amount:         sdk.NewCoin(params.BaseDenom, minDelegate),
		ChainName:      s.chainName,
	}
	_, err = s.MsgServer().AddDelegate(sdk.WrapSDKContext(s.Ctx), addMsg)
	s.Require().NoError(err)

	relayerData, _ = k.GetRelayer(s.Ctx, relayer)
	s.Require().Equal(int64(2), relayerData.SlashTimes,
		"FIXED: AddDelegate did NOT reset SlashTimes when < MaxSlashTimes")
	s.Require().True(relayerData.Online)

	// === CYCLE 2: 1 more slash → hits MaxSlashTimes → goes offline ===
	s.Require().NoError(k.SlashRelayer(s.Ctx, relayer.String()))

	relayerData, _ = k.GetRelayer(s.Ctx, relayer)
	s.Require().False(relayerData.Online, "cycle 2: offline after MaxSlashTimes (3)")
	s.Require().Equal(int64(3), relayerData.SlashTimes, "cycle 2: SlashTimes=3")

	// === RECOVERY: AddDelegate when fully offline → SHOULD reset SlashTimes ===
	_, err = s.MsgServer().AddDelegate(sdk.WrapSDKContext(s.Ctx), addMsg)
	s.Require().NoError(err)

	relayerData, _ = k.GetRelayer(s.Ctx, relayer)
	s.Require().Equal(int64(0), relayerData.SlashTimes,
		"FIXED: AddDelegate DID reset SlashTimes when >= MaxSlashTimes (recovery)")
	s.Require().True(relayerData.Online, "FIXED: relayer back online after recovery")
}
