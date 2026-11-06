package keeper_test

import (
	"fmt"
	"strings"
	"time"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/openmetaearth/me-hub/app/params"
	grouptypes "github.com/openmetaearth/me-hub/x/megroup/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func mustValidatorAddress(address string) sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return addr
}

func (s *KeeperTestSuite) TestSetRegionOperatorResetsValidator() {
	s.SetupTest()
	k := s.Keeper()
	region, found := k.GetRegion(s.Ctx, types.ExperienceRegionId)
	s.Require().True(found)
	oldAddr := mustValidatorAddress(region.OperatorAddress)
	old := k.MustGetValidator(s.Ctx, oldAddr)
	newAddr := sdk.ValAddress(s.TestAccs[0])
	consAddr, err := old.GetConsAddr()
	s.Require().NoError(err)
	power := k.GetLastValidatorPower(s.Ctx, oldAddr)
	s.Require().Positive(power)

	// A DAO rotation may leave stakes belonging to earlier DAOs. All must move.
	extra := types.NewStake(s.TestAccs[1], oldAddr, sdk.NewDec(7))
	extra.StartHeight = 123
	k.SetStake(s.Ctx, extra)
	stakes, err := k.GetStakesByValidator(s.Ctx, oldAddr)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(len(stakes), 2)
	// Experience-region delegators need not have KYC credentials.
	delegation := stakingtypes.NewDelegation(s.TestAccs[2], oldAddr, sdk.ZeroDec())
	delegation.StartHeight = 42
	delegation.UnMeidAmount = sdk.NewInt(55)
	k.SetDelegation(s.Ctx, delegation)
	k.SetGroupKeeper(s.App.GroupKeeper)
	s.App.GroupKeeper.SetGroupToRegion(s.Ctx, region.RegionId, 42)
	s.App.GroupKeeper.SetGroupInfo(s.Ctx, grouptypes.GroupInfo{Id: 42, Admin: old.OwnerAddress})
	signing := slashingtypes.NewValidatorSigningInfo(consAddr, 17, 9, time.Unix(100, 0).UTC(), false, 3)
	s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, consAddr, signing)
	s.App.SlashingKeeper.SetValidatorMissedBlockBitArray(s.Ctx, consAddr, 2, true)

	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	_, found = k.GetValidator(s.Ctx, oldAddr)
	s.Require().False(found)
	updated := k.MustGetValidator(s.Ctx, newAddr)
	expected := old
	expected.OperatorAddress = newAddr.String()
	expected.OwnerAddress = sdk.AccAddress(newAddr).String()
	s.Require().Equal(expected, updated)
	byCons, found := k.GetValidatorByConsAddr(s.Ctx, consAddr)
	s.Require().True(found)
	s.Require().Equal(newAddr.String(), byCons.OperatorAddress)
	s.Require().Zero(k.GetLastValidatorPower(s.Ctx, oldAddr))
	s.Require().Equal(power, k.GetLastValidatorPower(s.Ctx, newAddr))
	updatedRegion, found := k.GetRegion(s.Ctx, region.RegionId)
	s.Require().True(found)
	region.OperatorAddress = newAddr.String()
	s.Require().Equal(region, updatedRegion)
	group, found := s.App.GroupKeeper.GetGroupInfo(s.Ctx, 42)
	s.Require().True(found)
	s.Require().Equal(expected.OwnerAddress, group.Admin)
	for _, stake := range stakes {
		staker := sdk.MustAccAddressFromBech32(stake.StakerAddress)
		_, found = k.GetStake(s.Ctx, staker, oldAddr)
		s.Require().False(found)
		moved, found := k.GetStake(s.Ctx, staker, newAddr)
		s.Require().True(found)
		stake.ValidatorAddress = newAddr.String()
		s.Require().Equal(*stake, moved)
	}
	movedDelegation, found := k.GetDelegation(s.Ctx, s.TestAccs[2], newAddr)
	s.Require().True(found)
	delegation.ValidatorAddress = newAddr.String()
	s.Require().Equal(delegation, movedDelegation)
	gotSigning, found := s.App.SlashingKeeper.GetValidatorSigningInfo(s.Ctx, consAddr)
	s.Require().True(found)
	s.Require().Equal(signing, gotSigning)
	s.Require().True(s.App.SlashingKeeper.GetValidatorMissedBlockBitArray(s.Ctx, consAddr, 2))
	// Renaming the operator does not change the consensus validator or its power.
	updates, err := k.ApplyAndReturnValidatorSetUpdates(s.Ctx)
	s.Require().NoError(err)
	s.Require().Empty(updates)
}

func (s *KeeperTestSuite) TestSetRegionOperatorRejectsInvalidReplacement() {
	for _, name := range []string{"same operator", "existing validator", "module account", "invalid address", "missing region", "missing bound validator", "mismatched region", "target stake", "target delegation", "target unbonding stake", "target unbonding delegation", "target redelegation", "target last power"} {
		s.Run(name, func() {
			s.SetupTest()
			k := s.Keeper()
			region, found := k.GetRegion(s.Ctx, types.ExperienceRegionId)
			s.Require().True(found)
			regionID := region.RegionId
			old := k.MustGetValidator(s.Ctx, mustValidatorAddress(region.OperatorAddress))
			target := sdk.ValAddress(s.TestAccs[0]).String()
			switch name {
			case "same operator":
				target = old.OperatorAddress
			case "existing validator":
				for _, v := range k.GetAllValidators(s.Ctx) {
					if v.OperatorAddress != old.OperatorAddress {
						target = v.OperatorAddress
						break
					}
				}
			case "module account":
				target = sdk.ValAddress(authtypes.NewModuleAddress(types.BondedStakePoolName)).String()
			case "invalid address":
				target = "invalid"
			case "missing region":
				regionID = "missing"
			case "missing bound validator":
				region.OperatorAddress = sdk.ValAddress(s.TestAccs[1]).String()
				k.SetRegion(s.Ctx, region)
			case "mismatched region":
				old.Description.RegionID = "other"
				k.SetValidator(s.Ctx, old)
			default:
				s.setOccupiedOperatorTarget(name, old.GetOperator(), mustValidatorAddress(target))
			}
			before := s.App.ExportState(s.Ctx)
			events := len(s.Ctx.EventManager().Events())
			s.Require().Error(k.SetRegionOperator(s.Ctx, regionID, target))
			s.Require().Equal(before, s.App.ExportState(s.Ctx))
			s.Require().Len(s.Ctx.EventManager().Events(), events)
		})
	}
}

func (s *KeeperTestSuite) setOccupiedOperatorTarget(name string, oldAddr, target sdk.ValAddress) {
	k := s.Keeper()
	end := s.Ctx.BlockTime().Add(time.Hour)
	switch name {
	case "target stake":
		k.SetStake(s.Ctx, types.NewStake(s.TestAccs[1], target, sdk.OneDec()))
	case "target delegation":
		k.SetDelegation(s.Ctx, stakingtypes.NewDelegation(s.TestAccs[1], target, sdk.OneDec()))
	case "target unbonding stake":
		k.SetUnbondingStakeEntry(s.Ctx, s.TestAccs[1], target, 1, end, sdk.OneInt())
	case "target unbonding delegation":
		k.SetUnbondingDelegationEntry(s.Ctx, s.TestAccs[1], target, 1, end, sdk.OneInt())
	case "target redelegation":
		k.SetRedelegationEntry(s.Ctx, s.TestAccs[1], oldAddr, target, 1, end, sdk.OneInt(), sdk.OneDec(), sdk.OneDec())
	case "target last power":
		k.SetLastValidatorPower(s.Ctx, target, 1)
	default:
		s.FailNow("unknown occupied target test", name)
	}
}

func (s *KeeperTestSuite) TestResetValidatorKeepsJailedConsensusRemoval() {
	s.SetupTest()
	k := s.Keeper()
	region, _ := k.GetRegion(s.Ctx, types.ExperienceRegionId)
	old := k.MustGetValidator(s.Ctx, mustValidatorAddress(region.OperatorAddress))
	k.JailValidator(s.Ctx, old)
	power := k.GetLastValidatorPower(s.Ctx, old.GetOperator())
	s.Require().Positive(power)
	newAddr := sdk.ValAddress(s.TestAccs[0])
	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	renamed := k.MustGetValidator(s.Ctx, newAddr)
	s.Require().True(renamed.Jailed)
	s.Require().Equal(power, k.GetLastValidatorPower(s.Ctx, newAddr))
	updates, err := k.ApplyAndReturnValidatorSetUpdates(s.Ctx)
	s.Require().NoError(err)
	s.Require().Len(updates, 1)
	pubkey, err := old.ConsPubKey()
	s.Require().NoError(err)
	expectedKey, err := cryptocodec.ToTmProtoPublicKey(pubkey)
	s.Require().NoError(err)
	s.Require().Equal(expectedKey, updates[0].PubKey)
	s.Require().Zero(updates[0].Power)
	s.Require().True(k.MustGetValidator(s.Ctx, newAddr).IsUnbonding())
}

func (s *KeeperTestSuite) TestResetValidatorMigratesUnbondingState() {
	s.SetupTest()
	k := s.Keeper()
	s.Ctx = s.Ctx.WithBlockTime(time.Unix(1700000000, 0).UTC()).WithBlockHeight(20)
	region, _ := k.GetRegion(s.Ctx, types.ExperienceRegionId)
	oldAddr := mustValidatorAddress(region.OperatorAddress)
	newAddr := sdk.ValAddress(s.TestAccs[0])
	other := s.meEarthValidator.GetOperator()
	if other.Equals(oldAddr) {
		other = s.usaValidator.GetOperator()
	}
	staker, delegator := s.TestAccs[1], s.TestAccs[2]
	end := s.Ctx.BlockTime().Add(time.Hour)
	old := k.MustGetValidator(s.Ctx, oldAddr)
	old.Status = stakingtypes.Unbonding
	old.UnbondingHeight = 10
	old.UnbondingTime = end
	old.UnbondingOnHoldRefCount = 1
	old.UnbondingIds = []uint64{k.IncrementUnbondingID(s.Ctx)}
	k.SetValidator(s.Ctx, old)
	k.DeleteLastValidatorPower(s.Ctx, oldAddr)
	k.InsertUnbondingValidatorQueue(s.Ctx, old)
	k.SetValidatorByUnbondingID(s.Ctx, old, old.UnbondingIds[0])
	ubs := k.SetUnbondingStakeEntry(s.Ctx, staker, oldAddr, 10, end, sdk.NewInt(11))
	k.InsertUBSQueue(s.Ctx, ubs, end)
	ubd := k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 10, end, sdk.NewInt(13))
	k.InsertUBDQueue(s.Ctx, ubd, end)
	red := k.SetRedelegationEntry(s.Ctx, delegator, oldAddr, other, 10, end, sdk.NewInt(17), sdk.OneDec(), sdk.OneDec())
	k.InsertRedelegationQueue(s.Ctx, red, end)
	incoming := k.SetRedelegationEntry(s.Ctx, staker, other, oldAddr, 10, end, sdk.NewInt(19), sdk.OneDec(), sdk.OneDec())
	k.InsertRedelegationQueue(s.Ctx, incoming, end)

	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	moved := k.MustGetValidator(s.Ctx, newAddr)
	s.Require().Equal(old.UnbondingIds, moved.UnbondingIds)
	s.Require().Equal(old.UnbondingOnHoldRefCount, moved.UnbondingOnHoldRefCount)
	s.Require().False(s.Ctx.KVStore(k.GetStoreKey()).Has(stakingtypes.GetLastValidatorPowerKey(newAddr)))
	byID, found := k.GetValidatorByUnbondingID(s.Ctx, old.UnbondingIds[0])
	s.Require().True(found)
	s.Require().Equal(moved, byID)
	_, found = k.GetUnbondingStake(s.Ctx, staker, oldAddr)
	s.Require().False(found)
	movedUBS, found := k.GetUnbondingStake(s.Ctx, staker, newAddr)
	s.Require().True(found)
	ubs.ValidatorAddress = newAddr.String()
	s.Require().Equal(ubs, movedUBS)
	s.Require().Equal([]types.SVPair{{StakerAddress: staker.String(), ValidatorAddress: newAddr.String()}}, k.GetUBSQueueTimeSlice(s.Ctx, end))
	_, found = k.GetUnbondingDelegation(s.Ctx, delegator, oldAddr)
	s.Require().False(found)
	movedUBD, found := k.GetUnbondingDelegationByUnbondingID(s.Ctx, ubd.Entries[0].UnbondingId)
	s.Require().True(found)
	ubd.ValidatorAddress = newAddr.String()
	s.Require().Equal(ubd, movedUBD)
	s.Require().Empty(k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr))
	s.Require().Equal([]stakingtypes.UnbondingDelegation{ubd}, k.GetUnbondingDelegationsFromValidator(s.Ctx, newAddr))
	s.Require().Equal([]stakingtypes.DVPair{{DelegatorAddress: delegator.String(), ValidatorAddress: newAddr.String()}}, k.GetUBDQueueTimeSlice(s.Ctx, end))
	_, found = k.GetRedelegation(s.Ctx, delegator, oldAddr, other)
	s.Require().False(found)
	movedRED, found := k.GetRedelegationByUnbondingID(s.Ctx, red.Entries[0].UnbondingId)
	s.Require().True(found)
	red.ValidatorSrcAddress = newAddr.String()
	s.Require().Equal(red, movedRED)
	s.Require().Empty(k.GetRedelegationsFromSrcValidator(s.Ctx, oldAddr))
	s.Require().Equal([]stakingtypes.Redelegation{red}, k.GetRedelegationsFromSrcValidator(s.Ctx, newAddr))
	movedIncoming, found := k.GetRedelegationByUnbondingID(s.Ctx, incoming.Entries[0].UnbondingId)
	s.Require().True(found)
	incoming.ValidatorDstAddress = newAddr.String()
	s.Require().Equal(incoming, movedIncoming)
	s.Require().False(k.HasReceivingRedelegation(s.Ctx, staker, oldAddr))
	s.Require().True(k.HasReceivingRedelegation(s.Ctx, staker, newAddr))
	s.Require().Equal([]stakingtypes.DVVTriplet{
		{DelegatorAddress: delegator.String(), ValidatorSrcAddress: newAddr.String(), ValidatorDstAddress: other.String()},
		{DelegatorAddress: staker.String(), ValidatorSrcAddress: other.String(), ValidatorDstAddress: newAddr.String()},
	}, k.GetRedelegationQueueTimeSlice(s.Ctx, end))

	// Queue references must still complete normally after the address change.
	coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 11))
	s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, coins))
	s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, minttypes.ModuleName, types.NotBondedStakePoolName, coins))
	s.Ctx = s.Ctx.WithBlockTime(end).WithBlockHeight(100)
	pairs := k.SequeueAllMatureUBSQueue(s.Ctx, end)
	s.Require().Len(pairs, 1)
	paid, err := k.CompleteStakeUnBonding(s.Ctx, staker, mustValidatorAddress(pairs[0].ValidatorAddress))
	s.Require().NoError(err)
	s.Require().Equal(coins, paid)
	// On-hold validator remains queued; releasing the hold permits completion.
	k.UnbondAllMatureValidators(s.Ctx)
	s.Require().True(k.MustGetValidator(s.Ctx, newAddr).IsUnbonding())
	s.Require().NoError(k.UnbondingCanComplete(s.Ctx, old.UnbondingIds[0]))
	k.UnbondAllMatureValidators(s.Ctx)
	s.Require().True(k.MustGetValidator(s.Ctx, newAddr).IsUnbonded())
}

func (s *KeeperTestSuite) TestResetValidatorRejectsPendingConsensusRotation() {
	for _, offset := range []int64{-1, 0, 1, 2} {
		s.Run(fmt.Sprint(offset), func() {
			s.SetupTest()
			k := s.Keeper()
			region, _ := k.GetRegion(s.Ctx, types.ExperienceRegionId)
			pubkey := ed25519.GenPrivKey().PubKey()
			request, err := types.NewMsgReplaceConsensusPubKeyRequest(s.Dao.GlobalDao, region.OperatorAddress, pubkey, types.MinReplacePubKeyBlockNumber)
			s.Require().NoError(err)
			_, err = s.msgServer.ReplaceConsensusPubKey(s.Ctx, request)
			s.Require().NoError(err)
			info, err := k.GetReplaceConsensusPubKeyInfo(s.Ctx)
			s.Require().NoError(err)
			s.Ctx = s.Ctx.WithBlockHeight(info.UpdateAtHeight + offset)
			if offset >= 0 {
				_, err = k.UpdateValidatorPubKey(s.Ctx.WithBlockHeight(info.UpdateAtHeight))
				s.Require().NoError(err)
			}
			before := s.App.ExportState(s.Ctx)
			events := len(s.Ctx.EventManager().Events())
			err = k.SetRegionOperator(s.Ctx, region.RegionId, sdk.ValAddress(s.TestAccs[0]).String())
			s.Require().ErrorIs(err, types.ErrExistingReplaceValidator)
			s.Require().Equal(before, s.App.ExportState(s.Ctx))
			s.Require().Len(s.Ctx.EventManager().Events(), events)
		})
	}
}

func (s *KeeperTestSuite) TestResetValidatorRollsBackPositionMigration() {
	s.SetupTest()
	k := s.Keeper()
	region, _ := k.GetRegion(s.Ctx, types.ExperienceRegionId)
	target := sdk.ValAddress(s.TestAccs[0])
	// Force an error after positions have been re-keyed in the cache.
	s.App.DistrKeeper.SetValidatorCurrentRewards(s.Ctx, target, distrtypes.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}, Period: 99})
	before := s.App.ExportState(s.Ctx)
	events := len(s.Ctx.EventManager().Events())
	s.Require().Error(k.SetRegionOperator(s.Ctx, region.RegionId, target.String()))
	s.Require().Equal(before, s.App.ExportState(s.Ctx))
	s.Require().Len(s.Ctx.EventManager().Events(), events)
	// UpdateValidator also rolls back description/region edits made before reset.
	_, err := s.msgServer.UpdateValidator(s.Ctx, &types.MsgUpdateValidator{
		StakerAddress: s.Dao.GlobalDao, OperatorAddress: region.OperatorAddress, OwnerAddress: sdk.AccAddress(target).String(),
		Description: stakingtypes.Description{RegionID: "ind"},
	})
	s.Require().Error(err)
	s.Require().Equal(before, s.App.ExportState(s.Ctx))
	s.Require().Len(s.Ctx.EventManager().Events(), events)
}

func (s *KeeperTestSuite) TestResetValidatorCanonicalAddress() {
	s.SetupTest()
	region, _ := s.Keeper().GetRegion(s.Ctx, types.ExperienceRegionId)
	target := sdk.ValAddress(s.TestAccs[0])
	s.Require().NoError(s.Keeper().SetRegionOperator(s.Ctx, region.RegionId, strings.ToUpper(target.String())))
	s.Require().Equal(target.String(), s.Keeper().MustGetValidator(s.Ctx, target).OperatorAddress)
}

func (s *KeeperTestSuite) TestUpdateValidatorResetsAllStakes() {
	s.SetupTest()
	region, found := s.Keeper().GetRegion(s.Ctx, types.ExperienceRegionId)
	s.Require().True(found)
	oldAddr := mustValidatorAddress(region.OperatorAddress)
	// The authorized DAO is intentionally not the original genesis staker.
	_, found = s.Keeper().GetStake(s.Ctx, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), oldAddr)
	s.Require().False(found)
	response, err := s.msgServer.UpdateValidator(s.Ctx, &types.MsgUpdateValidator{
		StakerAddress: s.Dao.GlobalDao, OperatorAddress: oldAddr.String(), OwnerAddress: s.TestAccs[0].String(),
		Description: stakingtypes.Description{RegionID: stakingtypes.DoNotModifyDesc},
	})
	s.Require().NoError(err)
	s.Require().NotNil(response)
	_, found = s.Keeper().GetValidator(s.Ctx, oldAddr)
	s.Require().False(found)
	validator := s.Keeper().MustGetValidator(s.Ctx, sdk.ValAddress(s.TestAccs[0]))
	s.Require().Equal(region.RegionId, validator.Description.RegionID)
}
