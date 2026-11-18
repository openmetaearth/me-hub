package keeper_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/runtime"
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
	power := mustResetValue(k.GetLastValidatorPower(s.Ctx, oldAddr))
	s.Require().Positive(power)

	// The v3 fixture has no genesis stakes. Seed the original owner explicitly.
	k.SetStake(s.Ctx, types.NewStake(sdk.AccAddress(oldAddr), oldAddr, sdkmath.LegacyOneDec()))
	// A DAO rotation may leave stakes belonging to earlier DAOs. All must move.
	extra := types.NewStake(s.TestAccs[1], oldAddr, sdkmath.LegacyNewDec(7))
	extra.StartHeight = 123
	k.SetStake(s.Ctx, extra)
	stakes, err := k.GetStakesByValidator(s.Ctx, oldAddr)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(len(stakes), 2)
	// Experience-region delegators need not have KYC credentials.
	delegation := stakingtypes.NewDelegation(s.TestAccs[2].String(), oldAddr.String(), sdkmath.LegacyZeroDec())
	delegation.StartHeight = 42
	delegation.UnMeidAmount = sdkmath.NewInt(55)
	k.SetDelegation(s.Ctx, delegation)
	k.SetGroupKeeper(s.App.GroupKeeper)
	s.App.GroupKeeper.SetGroupToRegion(s.Ctx, region.RegionId, 42)
	s.App.GroupKeeper.SetGroupInfo(s.Ctx, grouptypes.GroupInfo{Id: 42, Admin: old.OwnerAddress})
	signing := slashingtypes.NewValidatorSigningInfo(consAddr, 17, 9, time.Unix(100, 0).UTC(), false, 3)
	s.Require().NoError(s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, consAddr, signing))
	s.Require().NoError(s.App.SlashingKeeper.SetMissedBlockBitmapValue(s.Ctx, consAddr, 2, true))

	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	_, err = k.GetValidator(s.Ctx, oldAddr)
	s.Require().Error(err)
	updated := k.MustGetValidator(s.Ctx, newAddr)
	expected := old
	expected.OperatorAddress = newAddr.String()
	expected.OwnerAddress = sdk.AccAddress(newAddr).String()
	s.Require().Equal(expected, updated)
	byCons, err := k.GetValidatorByConsAddr(s.Ctx, consAddr)
	s.Require().NoError(err)
	s.Require().Equal(newAddr.String(), byCons.OperatorAddress)
	s.Require().Zero(mustResetValue(k.GetLastValidatorPower(s.Ctx, oldAddr)))
	s.Require().Equal(power, mustResetValue(k.GetLastValidatorPower(s.Ctx, newAddr)))
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
	movedDelegation, err := k.GetDelegation(s.Ctx, s.TestAccs[2], newAddr)
	s.Require().NoError(err)
	delegation.ValidatorAddress = newAddr.String()
	s.Require().Equal(delegation, movedDelegation)
	gotSigning, err := s.App.SlashingKeeper.GetValidatorSigningInfo(s.Ctx, consAddr)
	s.Require().NoError(err)
	s.Require().Equal(signing, gotSigning)
	s.Require().True(mustResetValue(s.App.SlashingKeeper.GetMissedBlockBitmapValue(s.Ctx, consAddr, 2)))
	pubkey, err := s.App.SlashingKeeper.GetPubkey(s.Ctx, consAddr)
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
	delegation := stakingtypes.NewDelegation(delegator.String(), old.OperatorAddress, sdkmath.LegacyNewDec(3))
	delegation.Amount = sdkmath.NewInt(31)
	delegation.Unmovable = sdkmath.NewInt(7)
	delegation.StartHeight = 42
	k.SetDelegation(s.Ctx, delegation)
	// Delegations in other regions must not be changed.
	other := stakingtypes.NewDelegation(s.TestAccs[2].String(), s.meEarthValidator.OperatorAddress, sdkmath.LegacyNewDec(5))
	s.setOperatorReplacementKYC(s.TestAccs[2], s.meEarthValidator.Description.RegionID)
	k.SetDelegation(s.Ctx, other)

	s.Require().NoError(k.SetRegionOperator(s.Ctx, old.Description.RegionID, newAddr.String()))
	delegation.ValidatorAddress = newAddr.String()
	actual, err := k.GetDelegation(s.Ctx, delegator, sdk.ValAddress{})
	s.Require().NoError(err)
	s.Require().Equal(delegation, actual)
	actual, err = k.GetDelegation(s.Ctx, s.TestAccs[2], sdk.ValAddress{})
	s.Require().NoError(err)
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
			delegation := stakingtypes.NewDelegation(s.TestAccs[1].String(), old.OperatorAddress, sdkmath.LegacyOneDec())
			k.SetDelegation(s.Ctx, delegation)
			newAddr := sdk.ValAddress(s.TestAccs[0])
			_, err := s.msgServer.UpdateValidator(s.Ctx, &types.MsgUpdateValidator{
				StakerAddress: s.Dao.GlobalDao, OperatorAddress: old.OperatorAddress, OwnerAddress: s.TestAccs[0].String(),
				Description: stakingtypes.Description{RegionID: "ind"},
			})
			s.Require().NoError(err)
			delegation.ValidatorAddress = newAddr.String()
			actual, err := k.GetDelegation(s.Ctx, s.TestAccs[1], sdk.ValAddress{})
			s.Require().NoError(err)
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
	s.Require().NoError(d.SetValidatorCurrentRewards(s.Ctx, oldAddr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 9)))
	s.Require().NoError(d.SetValidatorHistoricalRewards(s.Ctx, oldAddr, 8, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1)))
	s.Require().NoError(d.SetValidatorSlashEvent(s.Ctx, oldAddr, 10, 8, distrtypes.NewValidatorSlashEvent(8, sdkmath.LegacyZeroDec())))
	feePool := mustResetValue(d.FeePool.Get(s.Ctx))

	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	// Protobuf decodes empty repeated fields as nil, not sdk.DecCoins{}.
	current := mustResetValue(d.GetValidatorCurrentRewards(s.Ctx, newAddr))
	s.Require().Equal(uint64(1), current.Period)
	s.Require().Empty(current.Rewards)
	historical := mustResetValue(d.GetValidatorHistoricalRewards(s.Ctx, newAddr, 0))
	s.Require().Equal(uint32(1), historical.ReferenceCount)
	s.Require().Empty(historical.CumulativeRewardRatio)
	s.Require().True(mustResetValue(d.GetValidatorOutstandingRewardsCoins(s.Ctx, newAddr)).IsZero())
	s.Require().True(mustResetValue(d.GetValidatorAccumulatedCommission(s.Ctx, newAddr)).Commission.IsZero())
	s.Require().Equal(feePool, mustResetValue(d.FeePool.Get(s.Ctx)))
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
	_, found, err := d.GetValidatorSlashEvent(s.Ctx, oldAddr, 10, 8)
	s.Require().NoError(err)
	s.Require().False(found)
	// The recreated distribution records remain usable by registered hooks.
	s.Require().NotPanics(func() {
		s.Require().NoError(k.Hooks().BeforeValidatorSlashed(s.Ctx, newAddr, sdkmath.LegacyZeroDec()))
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
				for _, validator := range mustResetValue(k.GetAllValidators(s.Ctx)) {
					if validator.OperatorAddress != oldAddr.String() {
						delegationOperator = mustValidatorAddress(validator.OperatorAddress)
						break
					}
				}
				s.Require().NotEqual(oldAddr, delegationOperator)
			}
			// Demand staking does not prevent this account from becoming an operator.
			delegation := stakingtypes.NewDelegation(newOwner.String(), delegationOperator.String(), sdkmath.LegacyNewDec(3))
			delegation.Amount = sdkmath.NewInt(31)
			delegation.UnMeidAmount = sdkmath.NewInt(55)
			delegation.Unmovable = sdkmath.NewInt(7)
			delegation.StartHeight = 42
			k.SetDelegation(s.Ctx, delegation)

			s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
			s.Require().Equal(newOwner.String(), k.MustGetValidator(s.Ctx, newAddr).OwnerAddress)
			if name == "old operator" {
				delegation.ValidatorAddress = newAddr.String()
			}
			actual, err := k.GetDelegation(s.Ctx, newOwner, sdk.ValAddress{})
			s.Require().NoError(err)
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
				for _, v := range mustResetValue(k.GetAllValidators(s.Ctx)) {
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
				s.Require().NoError(k.SetValidator(s.Ctx, old))
			default:
				s.setOccupiedOperatorTarget(name, mustValidatorAddress(target))
			}
			before := s.operatorReplacementSnapshot()
			events := len(s.Ctx.EventManager().Events())
			err := k.SetRegionOperator(s.Ctx, regionID, target)
			s.Require().Error(err)
			s.Require().Equal(before, s.operatorReplacementSnapshot())
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
		mustResetValue(k.SetUnbondingDelegationEntry(s.Ctx, s.TestAccs[1], target, 1, end, sdkmath.OneInt()))
	case "target last power":
		s.Require().NoError(k.SetLastValidatorPower(s.Ctx, target, 1))
	default:
		s.FailNow("unknown occupied target test", name)
	}
}

func (s *KeeperTestSuite) TestResetValidatorKeepsJailedConsensusRemoval() {
	s.SetupTest()
	k := s.Keeper()
	region, _ := k.GetRegion(s.Ctx, types.ExperienceRegionId)
	old := k.MustGetValidator(s.Ctx, mustValidatorAddress(region.OperatorAddress))
	s.Require().NoError(k.JailValidator(s.Ctx, old))
	power := mustResetValue(k.GetLastValidatorPower(s.Ctx, mustValidatorAddress(old.OperatorAddress)))
	s.Require().Positive(power)
	newAddr := sdk.ValAddress(s.TestAccs[0])
	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	renamed := k.MustGetValidator(s.Ctx, newAddr)
	s.Require().True(renamed.Jailed)
	s.Require().Equal(power, mustResetValue(k.GetLastValidatorPower(s.Ctx, newAddr)))
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
	old.UnbondingIds = []uint64{mustResetValue(k.IncrementUnbondingID(s.Ctx))}
	s.Require().NoError(k.SetValidator(s.Ctx, old))
	s.Require().NoError(k.DeleteLastValidatorPower(s.Ctx, oldAddr))
	s.Require().NoError(k.InsertUnbondingValidatorQueue(s.Ctx, old))
	s.Require().NoError(k.SetValidatorByUnbondingID(s.Ctx, old, old.UnbondingIds[0]))
	ubd := mustResetValue(k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 10, end, sdkmath.NewInt(13)))
	s.Require().NoError(k.InsertUBDQueue(s.Ctx, ubd, end))

	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	moved := k.MustGetValidator(s.Ctx, newAddr)
	s.Require().Equal(old.UnbondingIds, moved.UnbondingIds)
	s.Require().Equal(old.UnbondingOnHoldRefCount, moved.UnbondingOnHoldRefCount)
	s.Require().False(s.Ctx.KVStore(k.GetStoreKey()).Has(stakingtypes.GetLastValidatorPowerKey(newAddr)))
	byID, err := k.GetValidatorByUnbondingID(s.Ctx, old.UnbondingIds[0])
	s.Require().NoError(err)
	s.Require().Equal(moved, byID)
	_, err = k.GetUnbondingDelegation(s.Ctx, delegator, oldAddr)
	s.Require().Error(err)
	movedUBD, err := k.GetUnbondingDelegationByUnbondingID(s.Ctx, ubd.Entries[0].UnbondingId)
	s.Require().NoError(err)
	ubd.ValidatorAddress = newAddr.String()
	s.Require().Equal(ubd, movedUBD)
	s.Require().Empty(mustResetValue(k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr)))
	s.Require().Equal([]stakingtypes.UnbondingDelegation{ubd}, mustResetValue(k.GetUnbondingDelegationsFromValidator(s.Ctx, newAddr)))
	s.Require().Equal([]stakingtypes.DVPair{{DelegatorAddress: delegator.String(), ValidatorAddress: newAddr.String()}}, mustResetValue(k.GetUBDQueueTimeSlice(s.Ctx, end)))
	// Queue references must still complete normally after the address change.
	coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 13))
	s.App.AccountKeeper.SetAccount(s.Ctx, s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, delegator))
	s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, coins))
	s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, minttypes.ModuleName, stakingtypes.NotBondedPoolName, coins))
	s.Ctx = s.Ctx.WithBlockTime(end).WithBlockHeight(100)
	pairs := mustResetValue(k.DequeueAllMatureUBDQueue(s.Ctx, end))
	s.Require().Len(pairs, 1)
	paid, err := k.CompleteUnbonding(s.Ctx, delegator, mustValidatorAddress(pairs[0].ValidatorAddress))
	s.Require().NoError(err)
	s.Require().Equal(coins, paid)
	// On-hold validator remains queued; releasing the hold permits completion.
	s.Require().NoError(k.UnbondAllMatureValidators(s.Ctx))
	s.Require().True(k.MustGetValidator(s.Ctx, newAddr).IsUnbonding())
	s.Require().NoError(k.UnbondingCanComplete(s.Ctx, old.UnbondingIds[0]))
	s.Require().NoError(k.UnbondAllMatureValidators(s.Ctx))
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
		k.SetStake(s.Ctx, types.NewStake(delegator, oldAddr, sdkmath.LegacyOneDec()))
		k.SetDelegation(s.Ctx, stakingtypes.NewDelegation(delegator.String(), oldAddr.String(), sdkmath.LegacyOneDec()))
		mustResetValue(k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 10, end, sdkmath.OneInt()))
		mustResetValue(k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 11, end, sdkmath.OneInt()))
		mustResetValue(k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 12, later, sdkmath.OneInt()))
		ubd := mustResetValue(k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 13, missing, sdkmath.OneInt()))
		unbondings = append(unbondings, ubd)
		pair := stakingtypes.DVPair{DelegatorAddress: delegator.String(), ValidatorAddress: oldAddr.String()}
		queue = append(queue, pair, pair) // Preserve duplicate references and ordering.
	}
	// Do not change unrelated operators or an unselected delegator's queue entry.
	queue = append(queue,
		stakingtypes.DVPair{DelegatorAddress: delegators[0].String(), ValidatorAddress: other.String()},
		stakingtypes.DVPair{DelegatorAddress: s.TestAccs[2].String(), ValidatorAddress: oldAddr.String()},
	)
	s.Require().NoError(k.SetUBDQueueTimeSlice(s.Ctx, end, queue))
	s.Require().NoError(k.SetUBDQueueTimeSlice(s.Ctx, later, queue))
	unrelatedStake := types.NewStake(s.TestAccs[2], other, sdkmath.LegacyNewDec(3))
	unrelatedDelegation := stakingtypes.NewDelegation(s.TestAccs[2].String(), other.String(), sdkmath.LegacyNewDec(4))
	k.SetStake(s.Ctx, unrelatedStake)
	k.SetDelegation(s.Ctx, unrelatedDelegation)
	s.Require().NoError(k.SetRegionOperator(s.Ctx, region.RegionId, newAddr.String()))
	for i := range queue[:count*2] {
		queue[i].ValidatorAddress = newAddr.String()
	}
	s.Require().Equal(queue, mustResetValue(k.GetUBDQueueTimeSlice(s.Ctx, end)))
	s.Require().Equal(queue, mustResetValue(k.GetUBDQueueTimeSlice(s.Ctx, later)))
	s.Require().False(s.Ctx.KVStore(k.GetStoreKey()).Has(stakingtypes.GetUnbondingDelegationTimeKey(missing)))
	for _, ubd := range unbondings {
		ubd.ValidatorAddress = newAddr.String()
		for _, entry := range ubd.Entries {
			moved, err := k.GetUnbondingDelegationByUnbondingID(s.Ctx, entry.UnbondingId)
			s.Require().NoError(err)
			s.Require().Equal(ubd, moved)
		}
		delegator := sdk.MustAccAddressFromBech32(ubd.DelegatorAddress)
		_, found := k.GetStake(s.Ctx, delegator, oldAddr)
		s.Require().False(found)
		stake, found := k.GetStake(s.Ctx, delegator, newAddr)
		s.Require().True(found)
		s.Require().Equal(types.NewStake(delegator, newAddr, sdkmath.LegacyOneDec()), stake)
		delegation, err := k.GetDelegation(s.Ctx, delegator, newAddr)
		s.Require().NoError(err)
		s.Require().Equal(stakingtypes.NewDelegation(delegator.String(), newAddr.String(), sdkmath.LegacyOneDec()), delegation)
	}
	s.Require().Empty(mustResetValue(k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr)))
	stake, found := k.GetStake(s.Ctx, s.TestAccs[2], other)
	s.Require().True(found)
	s.Require().Equal(unrelatedStake, stake)
	delegation, err := k.GetDelegation(s.Ctx, s.TestAccs[2], other)
	s.Require().NoError(err)
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
			before := s.operatorReplacementSnapshot()
			events := len(s.Ctx.EventManager().Events())
			err = k.SetRegionOperator(s.Ctx, region.RegionId, sdk.ValAddress(s.TestAccs[0]).String())
			s.Require().ErrorIs(err, types.ErrExistingReplaceValidator)
			s.Require().Equal(before, s.operatorReplacementSnapshot())
			s.Require().Len(s.Ctx.EventManager().Events(), events)
		})
	}
}

type failingValidatorCreationHook struct{ stakingtypes.MultiStakingHooks }

func (failingValidatorCreationHook) AfterValidatorCreated(context.Context, sdk.ValAddress) error {
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
		ubd := mustResetValue(k.SetUnbondingDelegationEntry(s.Ctx, delegator, oldAddr, 10, end, sdkmath.OneInt()))
		s.Require().NoError(k.InsertUBDQueue(s.Ctx, ubd, end))
	}
	queueBefore := mustResetValue(k.GetUBDQueueTimeSlice(s.Ctx, end))
	unbondingsBefore := mustResetValue(k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr))
	// Genesis export omits queue keys and secondary indexes; check them explicitly.
	assertUnbondingRollback := func() {
		s.Require().Equal(queueBefore, mustResetValue(k.GetUBDQueueTimeSlice(s.Ctx, end)))
		s.Require().Equal(unbondingsBefore, mustResetValue(k.GetUnbondingDelegationsFromValidator(s.Ctx, oldAddr)))
		s.Require().Empty(mustResetValue(k.GetUnbondingDelegationsFromValidator(s.Ctx, target)))
		for _, ubd := range unbondingsBefore {
			for _, entry := range ubd.Entries {
				byID, err := k.GetUnbondingDelegationByUnbondingID(s.Ctx, entry.UnbondingId)
				s.Require().NoError(err)
				s.Require().Equal(ubd, byID)
			}
		}
	}
	// Fail creation after position/queue writes and the normal lifecycle hooks ran.
	replacement := stakingkeeper.NewKeeper(k.GetCdc(), runtime.NewKVStoreService(k.GetStoreKey().(*storetypes.KVStoreKey)), s.App.AccountKeeper, s.App.BankKeeper, k.GetAuthority(), k.ValidatorAddressCodec(), k.ConsensusAddressCodec())
	replacement.SetHooks(stakingtypes.NewMultiStakingHooks(k.Hooks(), failingValidatorCreationHook{}))
	originalKeeper := k.Keeper
	k.Keeper = replacement
	defer func() { k.Keeper = originalKeeper }()
	before := s.operatorReplacementSnapshot()
	events := len(s.Ctx.EventManager().Events())
	s.Require().ErrorContains(k.SetRegionOperator(s.Ctx, region.RegionId, target.String()), "injected validator creation failure")
	s.Require().Equal(before, s.operatorReplacementSnapshot())
	s.Require().Len(s.Ctx.EventManager().Events(), events)
	assertUnbondingRollback()
	// UpdateValidator also rolls back description/region edits made before reset.
	_, err := s.msgServer.UpdateValidator(s.Ctx, &types.MsgUpdateValidator{
		StakerAddress: s.Dao.GlobalDao, OperatorAddress: region.OperatorAddress, OwnerAddress: sdk.AccAddress(target).String(),
		Description: stakingtypes.Description{RegionID: "ind"},
	})
	s.Require().Error(err)
	s.Require().Equal(before, s.operatorReplacementSnapshot())
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
	_, err = s.Keeper().GetValidator(s.Ctx, oldAddr)
	s.Require().Error(err)
	validator := s.Keeper().MustGetValidator(s.Ctx, sdk.ValAddress(s.TestAccs[0]))
	s.Require().Equal(region.RegionId, validator.Description.RegionID)
}

// Snapshot raw stores, including queues and secondary indexes. Genesis export
// depends on unrelated module initialization and omits several migration keys.
func (s *KeeperTestSuite) operatorReplacementSnapshot() map[string]map[string]string {
	result := make(map[string]map[string]string)
	for name, key := range s.App.GetKVStoreKeys() {
		values := make(map[string]string)
		iterator := s.Ctx.KVStore(key).Iterator(nil, nil)
		for ; iterator.Valid(); iterator.Next() {
			values[string(iterator.Key())] = string(iterator.Value())
		}
		s.Require().NoError(iterator.Close())
		result[name] = values
	}
	return result
}

// mustResetValue unwraps SDK v0.50 keeper results in assertions and setup.
func mustResetValue[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
