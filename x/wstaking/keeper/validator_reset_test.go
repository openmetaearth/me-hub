package keeper_test

import (
	"errors"
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
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/openmetaearth/me-hub/app/params"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
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
	pubkey, err := s.App.SlashingKeeper.GetPubkey(s.Ctx, consAddr.Bytes())
	s.Require().NoError(err)
	expectedPubkey, err := old.ConsPubKey()
	s.Require().NoError(err)
	s.Require().Equal(expectedPubkey, pubkey)
	// Renaming the operator does not change the consensus validator or its power.
	updates, err := k.ApplyAndReturnValidatorSetUpdates(s.Ctx)
	s.Require().NoError(err)
	s.Require().Empty(updates)
}

func (s *KeeperTestSuite) setOperatorReplacementKYC(addr sdk.AccAddress, regionID string) {
	did := "did:operator-replacement:" + addr.String()
	s.App.KycKeeper.SetDID(s.Ctx, addr, did)
	s.App.KycKeeper.SetDidInfo(s.Ctx, did, didtypes.DidInfo{
		Did: did, Address: addr.String(), RegionId: regionID,
		KycLevel: 2, Status: didtypes.DID_STATUS_ACTIVE,
	})
	vc := didtypes.Credential{Did: did, Data: []byte(regionID)}
	s.App.KycKeeper.SetKYC(s.Ctx, did, vc)
	s.App.KycKeeper.AddFilters(s.Ctx, did, [][]byte{[]byte(regionID)}, vc)
}

func (s *KeeperTestSuite) TestSetRegionOperatorUpdatesKycDelegations() {
	s.SetupTest()
	k := s.Keeper()
	old := s.usaValidator
	k.SetRegion(s.Ctx, types.Region{RegionId: old.Description.RegionID, OperatorAddress: old.OperatorAddress, RegionShare: old.Tokens})
	newAddr := sdk.ValAddress(s.TestAccs[0])
	delegator := s.TestAccs[1]
	s.setOperatorReplacementKYC(delegator, old.Description.RegionID)
	delegation := stakingtypes.NewDelegation(delegator, old.GetOperator(), sdk.NewDec(3))
	delegation.Amount = sdk.NewInt(31)
	delegation.Unmovable = sdk.NewInt(7)
	delegation.StartHeight = 42
	k.SetDelegation(s.Ctx, delegation)
	// Delegations in other regions must not be changed.
	other := stakingtypes.NewDelegation(s.TestAccs[2], s.meEarthValidator.GetOperator(), sdk.NewDec(5))
	s.setOperatorReplacementKYC(s.TestAccs[2], s.meEarthValidator.Description.RegionID)
	k.SetDelegation(s.Ctx, other)

	s.Require().NoError(k.SetRegionOperator(s.Ctx, old.Description.RegionID, newAddr.String()))
	delegation.ValidatorAddress = newAddr.String()
	actual, found := k.GetDelegation(s.Ctx, delegator, sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(delegation, actual)
	actual, found = k.GetDelegation(s.Ctx, s.TestAccs[2], sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal(other, actual)
	// Replacement completes immediately; it does not leave a next-block marker.
	s.Require().Empty(k.GetAllChangeDelegationValidator(s.Ctx))
}

func (s *KeeperTestSuite) TestUpdateValidatorUsesSourceRegionDelegations() {
	for _, regionID := range []string{types.ExperienceRegionId, "usa"} {
		s.Run(regionID, func() {
			s.SetupTest()
			k := s.Keeper()
			old := s.experienceValidator
			if regionID != types.ExperienceRegionId {
				old = s.usaValidator
				s.setOperatorReplacementKYC(s.TestAccs[1], regionID)
			}
			k.SetRegion(s.Ctx, types.Region{RegionId: regionID, OperatorAddress: old.OperatorAddress, RegionShare: old.Tokens})
			delegation := stakingtypes.NewDelegation(s.TestAccs[1], old.GetOperator(), sdk.OneDec())
			k.SetDelegation(s.Ctx, delegation)
			newAddr := sdk.ValAddress(s.TestAccs[0])
			_, err := s.msgServer.UpdateValidator(s.Ctx, &types.MsgUpdateValidator{
				StakerAddress: s.Dao.GlobalDao, OperatorAddress: old.OperatorAddress, OwnerAddress: s.TestAccs[0].String(),
				Description: stakingtypes.Description{RegionID: "ind"},
			})
			s.Require().NoError(err)
			delegation.ValidatorAddress = newAddr.String()
			actual, found := k.GetDelegation(s.Ctx, s.TestAccs[1], sdk.ValAddress{})
			s.Require().True(found)
			s.Require().Equal(delegation, actual)
			s.Require().Equal("ind", k.MustGetValidator(s.Ctx, newAddr).Description.RegionID)
			s.Require().Empty(k.GetAllChangeDelegationValidator(s.Ctx))
		})
	}
}

func (s *KeeperTestSuite) TestResetValidatorReinitializesDistribution() {
	s.SetupTest()
	k := s.Keeper()
	region, _ := k.GetRegion(s.Ctx, types.ExperienceRegionId)
	oldAddr := mustValidatorAddress(region.OperatorAddress)
	newAddr := sdk.ValAddress(s.TestAccs[0])
	d := s.App.DistrKeeper
	// The product uses region treasury rewards. Restore the original validator
	// lifecycle, not reward-period preservation under a different operator key.
	d.SetValidatorCurrentRewards(s.Ctx, oldAddr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 9))
	d.SetValidatorHistoricalRewards(s.Ctx, oldAddr, 8, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))
	d.SetValidatorSlashEvent(s.Ctx, oldAddr, 10, 8, distrtypes.NewValidatorSlashEvent(8, sdk.ZeroDec()))
	feePool := d.GetFeePool(s.Ctx)

	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	// Protobuf decodes empty repeated fields as nil, not sdk.DecCoins{}.
	current := d.GetValidatorCurrentRewards(s.Ctx, newAddr)
	s.Require().Equal(uint64(1), current.Period)
	s.Require().Empty(current.Rewards)
	historical := d.GetValidatorHistoricalRewards(s.Ctx, newAddr, 0)
	s.Require().Equal(uint32(1), historical.ReferenceCount)
	s.Require().Empty(historical.CumulativeRewardRatio)
	s.Require().True(d.GetValidatorOutstandingRewardsCoins(s.Ctx, newAddr).IsZero())
	s.Require().True(d.GetValidatorAccumulatedCommission(s.Ctx, newAddr).Commission.IsZero())
	s.Require().Equal(feePool, d.GetFeePool(s.Ctx))
	store := s.Ctx.KVStore(s.App.GetKey(distrtypes.StoreKey))
	for _, key := range [][]byte{
		distrtypes.GetValidatorCurrentRewardsKey(oldAddr),
		distrtypes.GetValidatorHistoricalRewardsKey(oldAddr, 0),
		distrtypes.GetValidatorHistoricalRewardsKey(oldAddr, 8),
		distrtypes.GetValidatorOutstandingRewardsKey(oldAddr),
		distrtypes.GetValidatorAccumulatedCommissionKey(oldAddr),
	} {
		s.Require().False(store.Has(key))
	}
	_, found := d.GetValidatorSlashEvent(s.Ctx, oldAddr, 10, 8)
	s.Require().False(found)
	// The recreated distribution records remain usable by registered hooks.
	s.Require().NotPanics(func() {
		s.Require().NoError(k.Hooks().BeforeValidatorSlashed(s.Ctx, newAddr, sdk.ZeroDec()))
	})
}

func (s *KeeperTestSuite) TestSetRegionOperatorAllowsOwnerWithDelegation() {
	for _, name := range []string{"old operator", "another operator"} {
		s.Run(name, func() {
			s.SetupTest()
			k := s.Keeper()
			region, found := k.GetRegion(s.Ctx, types.ExperienceRegionId)
			s.Require().True(found)
			oldAddr := mustValidatorAddress(region.OperatorAddress)
			newOwner := s.TestAccs[0]
			newAddr := sdk.ValAddress(newOwner)
			delegationOperator := oldAddr
			if name == "another operator" {
				for _, validator := range k.GetAllValidators(s.Ctx) {
					if validator.OperatorAddress != oldAddr.String() {
						delegationOperator = validator.GetOperator()
						break
					}
				}
				s.Require().NotEqual(oldAddr, delegationOperator)
			}
			// Demand staking does not prevent this account from becoming an operator.
			delegation := stakingtypes.NewDelegation(newOwner, delegationOperator, sdk.NewDec(3))
			delegation.Amount = sdk.NewInt(31)
			delegation.UnMeidAmount = sdk.NewInt(55)
			delegation.Unmovable = sdk.NewInt(7)
			delegation.StartHeight = 42
			k.SetDelegation(s.Ctx, delegation)

			s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
			s.Require().Equal(newOwner.String(), k.MustGetValidator(s.Ctx, newAddr).OwnerAddress)
			if name == "old operator" {
				delegation.ValidatorAddress = newAddr.String()
			}
			actual, found := k.GetDelegation(s.Ctx, newOwner, sdk.ValAddress{})
			s.Require().True(found)
			s.Require().Equal(delegation, actual)
		})
	}
}

func (s *KeeperTestSuite) TestSetRegionOperatorRejectsInvalidReplacement() {
	for _, name := range []string{"same operator", "existing validator", "module account", "invalid address", "missing region", "missing bound validator", "mismatched region", "target unbonding delegation", "target last power"} {
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
				s.setOccupiedOperatorTarget(name, mustValidatorAddress(target))
			}
			before := s.App.ExportState(s.Ctx)
			events := len(s.Ctx.EventManager().Events())
			err := k.SetRegionOperator(s.Ctx, regionID, target)
			s.Require().Error(err)
			s.Require().Equal(before, s.App.ExportState(s.Ctx))
			s.Require().Len(s.Ctx.EventManager().Events(), events)
			switch name {
			case "existing validator", "target unbonding delegation", "target last power":
				s.Require().ErrorIs(err, types.ErrValidatorExist)
			}
		})
	}
}

func (s *KeeperTestSuite) setOccupiedOperatorTarget(name string, target sdk.ValAddress) {
	k := s.Keeper()
	end := s.Ctx.BlockTime().Add(time.Hour)
	switch name {
	case "target unbonding delegation":
		k.SetUnbondingDelegationEntry(s.Ctx, s.TestAccs[1], target, 1, end, sdk.OneInt())
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
	delegator := s.TestAccs[2]
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
	ubd := k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 10, end, sdk.NewInt(13))
	k.InsertUBDQueue(s.Ctx, ubd, end)

	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	moved := k.MustGetValidator(s.Ctx, newAddr)
	s.Require().Equal(old.UnbondingIds, moved.UnbondingIds)
	s.Require().Equal(old.UnbondingOnHoldRefCount, moved.UnbondingOnHoldRefCount)
	s.Require().False(s.Ctx.KVStore(k.GetStoreKey()).Has(stakingtypes.GetLastValidatorPowerKey(newAddr)))
	byID, found := k.GetValidatorByUnbondingID(s.Ctx, old.UnbondingIds[0])
	s.Require().True(found)
	s.Require().Equal(moved, byID)
	_, found = k.GetUnbondingDelegation(s.Ctx, delegator, oldAddr)
	s.Require().False(found)
	movedUBD, found := k.GetUnbondingDelegationByUnbondingID(s.Ctx, ubd.Entries[0].UnbondingId)
	s.Require().True(found)
	ubd.ValidatorAddress = newAddr.String()
	s.Require().Equal(ubd, movedUBD)
	s.Require().Empty(k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr))
	s.Require().Equal([]stakingtypes.UnbondingDelegation{ubd}, k.GetUnbondingDelegationsFromValidator(s.Ctx, newAddr))
	s.Require().Equal([]stakingtypes.DVPair{{DelegatorAddress: delegator.String(), ValidatorAddress: newAddr.String()}}, k.GetUBDQueueTimeSlice(s.Ctx, end))
	// Queue references must still complete normally after the address change.
	coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 13))
	s.App.AccountKeeper.SetAccount(s.Ctx, s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, delegator))
	s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, coins))
	s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, minttypes.ModuleName, stakingtypes.NotBondedPoolName, coins))
	s.Ctx = s.Ctx.WithBlockTime(end).WithBlockHeight(100)
	pairs := k.DequeueAllMatureUBDQueue(s.Ctx, end)
	s.Require().Len(pairs, 1)
	paid, err := k.CompleteUnbonding(s.Ctx, delegator, mustValidatorAddress(pairs[0].ValidatorAddress))
	s.Require().NoError(err)
	s.Require().Equal(coins, paid)
	// On-hold validator remains queued; releasing the hold permits completion.
	k.UnbondAllMatureValidators(s.Ctx)
	s.Require().True(k.MustGetValidator(s.Ctx, newAddr).IsUnbonding())
	s.Require().NoError(k.UnbondingCanComplete(s.Ctx, old.UnbondingIds[0]))
	k.UnbondAllMatureValidators(s.Ctx)
	s.Require().True(k.MustGetValidator(s.Ctx, newAddr).IsUnbonded())
}

func (s *KeeperTestSuite) TestResetValidatorBatchesSharedUnbondingQueue() {
	s.SetupTest()
	k := s.Keeper()
	s.Ctx = s.Ctx.WithBlockTime(time.Unix(1700000000, 0).UTC()).WithBlockHeight(20)
	region, _ := k.GetRegion(s.Ctx, types.ExperienceRegionId)
	oldAddr := mustValidatorAddress(region.OperatorAddress)
	newAddr := sdk.ValAddress(s.TestAccs[0])
	other := sdk.ValAddress(s.TestAccs[1])
	end := s.Ctx.BlockTime().Add(time.Hour)
	later := end.Add(time.Nanosecond) // Distinct encoded keys must not be coalesced.
	missing := later.Add(time.Hour)
	const count = 1000
	delegators := s.NewAccounts(count)
	unbondings := make([]stakingtypes.UnbondingDelegation, 0, count)
	queue := make([]stakingtypes.DVPair, 0, count*2+2)
	for _, delegator := range delegators {
		k.SetStake(s.Ctx, types.NewStake(delegator, oldAddr, sdk.OneDec()))
		k.SetDelegation(s.Ctx, stakingtypes.NewDelegation(delegator, oldAddr, sdk.OneDec()))
		k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 10, end, sdk.OneInt())
		k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 11, end, sdk.OneInt())
		k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 12, later, sdk.OneInt())
		ubd := k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 13, missing, sdk.OneInt())
		unbondings = append(unbondings, ubd)
		pair := stakingtypes.DVPair{DelegatorAddress: delegator.String(), ValidatorAddress: oldAddr.String()}
		queue = append(queue, pair, pair) // Preserve duplicate references and ordering.
	}
	// Do not change unrelated operators or an unselected delegator's queue entry.
	queue = append(queue,
		stakingtypes.DVPair{DelegatorAddress: delegators[0].String(), ValidatorAddress: other.String()},
		stakingtypes.DVPair{DelegatorAddress: s.TestAccs[2].String(), ValidatorAddress: oldAddr.String()},
	)
	k.SetUBDQueueTimeSlice(s.Ctx, end, queue)
	k.SetUBDQueueTimeSlice(s.Ctx, later, queue)
	unrelatedStake := types.NewStake(s.TestAccs[2], other, sdk.NewDec(3))
	unrelatedDelegation := stakingtypes.NewDelegation(s.TestAccs[2], other, sdk.NewDec(4))
	k.SetStake(s.Ctx, unrelatedStake)
	k.SetDelegation(s.Ctx, unrelatedDelegation)
	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	for i := range queue[:count*2] {
		queue[i].ValidatorAddress = newAddr.String()
	}
	s.Require().Equal(queue, k.GetUBDQueueTimeSlice(s.Ctx, end))
	s.Require().Equal(queue, k.GetUBDQueueTimeSlice(s.Ctx, later))
	s.Require().False(s.Ctx.KVStore(k.GetStoreKey()).Has(stakingtypes.GetUnbondingDelegationTimeKey(missing)))
	for _, ubd := range unbondings {
		ubd.ValidatorAddress = newAddr.String()
		for _, entry := range ubd.Entries {
			moved, found := k.GetUnbondingDelegationByUnbondingID(s.Ctx, entry.UnbondingId)
			s.Require().True(found)
			s.Require().Equal(ubd, moved)
		}
		delegator := sdk.MustAccAddressFromBech32(ubd.DelegatorAddress)
		_, found := k.GetStake(s.Ctx, delegator, oldAddr)
		s.Require().False(found)
		stake, found := k.GetStake(s.Ctx, delegator, newAddr)
		s.Require().True(found)
		s.Require().Equal(types.NewStake(delegator, newAddr, sdk.OneDec()), stake)
		delegation, found := k.GetDelegation(s.Ctx, delegator, newAddr)
		s.Require().True(found)
		s.Require().Equal(stakingtypes.NewDelegation(delegator, newAddr, sdk.OneDec()), delegation)
	}
	s.Require().Empty(k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr))
	stake, found := k.GetStake(s.Ctx, s.TestAccs[2], other)
	s.Require().True(found)
	s.Require().Equal(unrelatedStake, stake)
	delegation, found := k.GetDelegation(s.Ctx, s.TestAccs[2], other)
	s.Require().True(found)
	s.Require().Equal(unrelatedDelegation, delegation)
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

type failingValidatorCreationHook struct{ stakingtypes.MultiStakingHooks }

func (failingValidatorCreationHook) AfterValidatorCreated(sdk.Context, sdk.ValAddress) error {
	return errors.New("injected validator creation failure")
}

func (s *KeeperTestSuite) TestResetValidatorRollsBackPositionMigration() {
	s.SetupTest()
	k := s.Keeper()
	region, _ := k.GetRegion(s.Ctx, types.ExperienceRegionId)
	target := sdk.ValAddress(s.TestAccs[0])
	oldAddr := mustValidatorAddress(region.OperatorAddress)
	end := s.Ctx.BlockTime().Add(time.Hour)
	for _, delegator := range s.TestAccs[1:] {
		ubd := k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 10, end, sdk.OneInt())
		k.InsertUBDQueue(s.Ctx, ubd, end)
	}
	queueBefore := k.GetUBDQueueTimeSlice(s.Ctx, end)
	unbondingsBefore := k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr)
	// Genesis export omits queue keys and secondary indexes; check them explicitly.
	assertUnbondingRollback := func() {
		s.Require().Equal(queueBefore, k.GetUBDQueueTimeSlice(s.Ctx, end))
		s.Require().Equal(unbondingsBefore, k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr))
		s.Require().Empty(k.GetUnbondingDelegationsFromValidator(s.Ctx, target))
		for _, ubd := range unbondingsBefore {
			for _, entry := range ubd.Entries {
				byID, found := k.GetUnbondingDelegationByUnbondingID(s.Ctx, entry.UnbondingId)
				s.Require().True(found)
				s.Require().Equal(ubd, byID)
			}
		}
	}
	// Fail creation after position/queue writes and the normal lifecycle hooks ran.
	replacement := stakingkeeper.NewKeeper(k.GetCdc(), k.GetStoreKey(), s.App.AccountKeeper, s.App.BankKeeper, k.GetAuthority())
	replacement.SetHooks(stakingtypes.NewMultiStakingHooks(k.Hooks(), failingValidatorCreationHook{}))
	originalKeeper := k.Keeper
	k.Keeper = replacement
	defer func() { k.Keeper = originalKeeper }()
	before := s.App.ExportState(s.Ctx)
	events := len(s.Ctx.EventManager().Events())
	s.Require().ErrorContains(k.SetRegionOperator(s.Ctx, region.RegionId, target.String()), "injected validator creation failure")
	s.Require().Equal(before, s.App.ExportState(s.Ctx))
	s.Require().Len(s.Ctx.EventManager().Events(), events)
	assertUnbondingRollback()
	// UpdateValidator also rolls back description/region edits made before reset.
	_, err := s.msgServer.UpdateValidator(s.Ctx, &types.MsgUpdateValidator{
		StakerAddress: s.Dao.GlobalDao, OperatorAddress: region.OperatorAddress, OwnerAddress: sdk.AccAddress(target).String(),
		Description: stakingtypes.Description{RegionID: "ind"},
	})
	s.Require().Error(err)
	s.Require().Equal(before, s.App.ExportState(s.Ctx))
	s.Require().Len(s.Ctx.EventManager().Events(), events)
	assertUnbondingRollback()
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
