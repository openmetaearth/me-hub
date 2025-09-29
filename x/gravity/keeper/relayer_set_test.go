package keeper_test

//
//import (
//	"encoding/hex"
//	"fmt"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/stretchr/testify/require"
//	tmrand "github.com/tendermint/tendermint/libs/rand"
//
//	"github.com/st-chain/me-hub/testutil/helpers"
//	"github.com/st-chain/me-hub/x/gravity/types"
//)
//
//func (s *KeeperTestSuite) TestLastPendingRelayerSetRequestByAddr() {
//	testCases := []struct {
//		RelayerAddress sdk.AccAddress
//		StartHeight    int64
//
//		ExpectGravitySetSize int
//	}{
//		{
//			RelayerAddress:       s.relayerAddrs[0],
//			StartHeight:          1,
//			ExpectGravitySetSize: 3,
//		},
//		{
//			RelayerAddress:       s.relayerAddrs[1],
//			StartHeight:          2,
//			ExpectGravitySetSize: 2,
//		},
//		{
//			RelayerAddress:       s.relayerAddrs[2],
//			StartHeight:          3,
//			ExpectGravitySetSize: 1,
//		},
//	}
//
//	for i := 1; i <= 3; i++ {
//		s.Keeper().SetRelayer(s.Ctx, &types.GravitySet{
//			Nonce: uint64(i),
//			Members: types.BridgeValidators{{
//				Power:           uint64(i),
//				ExternalAddress: fmt.Sprintf("0x%d", i),
//			}},
//			Height: uint64(i),
//		})
//	}
//
//	wrapSDKContext := sdk.WrapSDKContext(s.Ctx)
//	for _, testCase := range testCases {
//		relayer := types.Relayer{
//			RelayerAddress: testCase.RelayerAddress.String(),
//			StartHeight:    testCase.StartHeight,
//		}
//		s.Keeper().SetRelayer(s.Ctx, relayer)
//		s.Keeper().SetGravityByBridger(s.Ctx, testCase.BridgerAddress, relayer.GetGravity())
//
//		response, err := s.Keeper().LastPendingGravitySetRequestByAddr(wrapSDKContext,
//			&types.QueryLastPendingGravitySetRequestByAddrRequest{
//				BridgerAddress: testCase.BridgerAddress.String(),
//			})
//		require.NoError(s.T(), err)
//		require.EqualValues(s.T(), testCase.ExpectGravitySetSize, len(response.GravitySets))
//	}
//}
//
//func (s *KeeperTestSuite) TestGetUnSlashedGravitySets() {
//	height := 100
//	index := 10
//	for i := 1; i <= index; i++ {
//		s.Keeper().StoreGravitySet(s.Ctx, &types.GravitySet{
//			Nonce: uint64(i),
//			Members: types.BridgeValidators{{
//				Power:           tmrand.Uint64(),
//				ExternalAddress: helpers.GenerateAddress().Hex(),
//			}},
//			Height: uint64(height + i),
//		})
//	}
//	s.Equal(uint64(0), s.Keeper().GetLastSlashedGravitySetNonce(s.Ctx))
//
//	sets := s.Keeper().GetUnSlashedGravitySets(s.Ctx, uint64(height+index))
//	require.EqualValues(s.T(), index-1, sets.Len())
//
//	s.Keeper().SetLastSlashedGravitySetNonce(s.Ctx, 1)
//	sets = s.Keeper().GetUnSlashedGravitySets(s.Ctx, uint64(height+index))
//	require.EqualValues(s.T(), index-2, sets.Len())
//
//	sets = s.Keeper().GetUnSlashedGravitySets(s.Ctx, uint64(height+index+1))
//	require.EqualValues(s.T(), index-1, sets.Len())
//}
//
//func (s *KeeperTestSuite) TestKeeper_IterateGravitySetConfirmByNonce() {
//	index := tmrand.Intn(20) + 1
//	for i := uint64(1); i <= uint64(index); i++ {
//		for _, relayer := range s.relayerAddrs {
//			s.Keeper().SetGravitySetConfirm(s.Ctx, relayer,
//				&types.MsgGravitySetConfirm{
//					Nonce:           i,
//					BridgerAddress:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
//					ExternalAddress: helpers.GenerateAddress().Hex(),
//					Signature:       "",
//					ChainName:       s.chainName,
//				},
//			)
//		}
//	}
//
//	index = tmrand.Intn(index) + 1
//	var confirms []*types.MsgGravitySetConfirm
//	s.Keeper().IterateGravitySetConfirmByNonce(s.Ctx, uint64(index), func(confirm *types.MsgGravitySetConfirm) bool {
//		confirms = append(confirms, confirm)
//		return false
//	})
//	s.Equal(len(confirms), len(s.relayerAddrs), index)
//}
//
//func (s *KeeperTestSuite) TestKeeper_DeleteGravitySetConfirm() {
//	member := make([]types.BridgeValidator, 0, len(s.externalPris))
//	for i, external := range s.externalPris {
//		externalAddr := s.PubKeyToExternalAddr(external.PublicKey)
//		member = append(member, types.BridgeValidator{
//			Power:           uint64(i),
//			ExternalAddress: externalAddr,
//		})
//	}
//	relayerSet := &types.GravitySet{
//		Nonce:   1,
//		Members: member,
//		Height:  uint64(s.Ctx.BlockHeight()),
//	}
//	s.Keeper().StoreGravitySet(s.Ctx, relayerSet)
//
//	for i, external := range s.externalPris {
//		externalAddress, signature := s.SignGravitySetConfirm(external, relayerSet)
//		s.Keeper().SetGravitySetConfirm(s.Ctx, s.relayerAddrs[i],
//			&types.MsgGravitySetConfirm{
//				Nonce:           relayerSet.Nonce,
//				BridgerAddress:  s.bridgerAddrs[i].String(),
//				ExternalAddress: externalAddress,
//				Signature:       hex.EncodeToString(signature),
//				ChainName:       s.chainName,
//			},
//		)
//	}
//	s.Keeper().SetLastObservedGravitySet(s.Ctx,
//		&types.GravitySet{
//			Nonce: relayerSet.Nonce + 1,
//		},
//	)
//
//	params := s.Keeper().GetParams(s.Ctx)
//	params.SignedWindow = 10
//	err := s.Keeper().SetParams(s.Ctx, &params)
//	s.Require().NoError(err)
//	s.Commit()
//	for _, relayer := range s.relayerAddrs {
//		s.NotNil(s.Keeper().GetGravitySetConfirm(s.Ctx, relayerSet.Nonce, relayer))
//	}
//
//	s.Commit(int64(params.SignedWindow + 1))
//	for _, relayer := range s.relayerAddrs {
//		s.Nil(s.Keeper().GetGravitySetConfirm(s.Ctx, relayerSet.Nonce, relayer))
//	}
//}
//
//func (s *KeeperTestSuite) TestKeeper_IterateGravitySet() {
//	member := make([]types.BridgeValidator, 0, len(s.externalPris))
//	for i, external := range s.externalPris {
//		member = append(member, types.BridgeValidator{
//			Power:           uint64(i),
//			ExternalAddress: crypto.PubkeyToAddress(external.PublicKey).String(),
//		})
//	}
//	for i := 1; i <= 10; i++ {
//		s.Keeper().StoreGravitySet(s.Ctx, &types.GravitySet{
//			Nonce:   uint64(i),
//			Members: member,
//			Height:  uint64(i + 100),
//		})
//	}
//	i := uint64(0)
//	relayerSets := types.GravitySets{}
//	s.Keeper().IterateGravitySetByNonce(s.Ctx, 0, func(relayerSet *types.GravitySet) bool {
//		i = i + 1
//		s.Equal(i, relayerSet.Nonce)
//		relayerSets = append(relayerSets, relayerSet)
//		return false
//	})
//	s.Equal(len(relayerSets), 10)
//
//	relayerSets = types.GravitySets{}
//	s.Keeper().IterateGravitySetByNonce(s.Ctx, 1, func(relayerSet *types.GravitySet) bool {
//		relayerSets = append(relayerSets, relayerSet)
//		return false
//	})
//	s.Equal(len(relayerSets), 10)
//
//	relayerSets = types.GravitySets{}
//	s.Keeper().IterateGravitySetByNonce(s.Ctx, 2, func(relayerSet *types.GravitySet) bool {
//		relayerSets = append(relayerSets, relayerSet)
//		return false
//	})
//	s.Equal(len(relayerSets), 9)
//
//	s.Keeper().IterateGravitySets(s.Ctx, true, func(relayerSet *types.GravitySet) bool {
//		s.Equal(i, relayerSet.Nonce, relayerSet.Nonce)
//		i = i - 1
//		return false
//	})
//
//	s.Keeper().IterateGravitySets(s.Ctx, false, func(relayerSet *types.GravitySet) bool {
//		i = i + 1
//		s.Equal(i, relayerSet.Nonce)
//		return false
//	})
//}
