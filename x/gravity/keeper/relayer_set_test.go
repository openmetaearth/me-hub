package keeper_test

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestLastPendingRelayerSetRequestByAddr() {
	testCases := []struct {
		RelayerAddress sdk.AccAddress
		StartHeight    int64

		ExpectRelayerSetSize int
	}{
		{
			RelayerAddress:       s.relayerAddrs[0],
			StartHeight:          1,
			ExpectRelayerSetSize: 3,
		},
		{
			RelayerAddress:       s.relayerAddrs[1],
			StartHeight:          2,
			ExpectRelayerSetSize: 2,
		},
		{
			RelayerAddress:       s.relayerAddrs[2],
			StartHeight:          3,
			ExpectRelayerSetSize: 1,
		},
	}

	for i := 1; i <= 3; i++ {
		s.Keeper().StoreRelayerSet(s.Ctx, &types.RelayerSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           uint64(i),
				ExternalAddress: fmt.Sprintf("0x%d", i),
			}},
			Height: uint64(i),
		})
	}

	wrapSDKContext := sdk.WrapSDKContext(s.Ctx)
	for _, testCase := range testCases {
		relayer := types.Relayer{
			RelayerAddress: testCase.RelayerAddress.String(),
			StartHeight:    testCase.StartHeight,
		}
		s.Keeper().SetRelayer(s.Ctx, testCase.RelayerAddress, relayer)
		response, err := s.QueryClient().LastPendingRelayerSetRequestByAddr(wrapSDKContext,
			&types.QueryLastPendingRelayerSetRequestByAddrRequest{
				RelayerAddress: testCase.RelayerAddress.String(),
			})
		require.NoError(s.T(), err)
		require.EqualValues(s.T(), testCase.ExpectRelayerSetSize, len(response.GetRelayerSets()))
	}
}

func (s *KeeperTestSuite) TestGetUnSlashedRelayerSets() {
	height := 100
	index := 10
	for i := 1; i <= index; i++ {
		s.Keeper().StoreRelayerSet(s.Ctx, &types.RelayerSet{
			Nonce: uint64(i),
			Members: types.BridgeValidators{{
				Power:           tmrand.Uint64(),
				ExternalAddress: helpers.GenerateAddress().Hex(),
			}},
			Height: uint64(height + i),
		})
	}
	s.Equal(uint64(0), s.Keeper().GetLastSlashedRelayerSetNonce(s.Ctx))

	sets := s.Keeper().GetUnSlashedRelayerSets(s.Ctx, uint64(height+index))
	require.EqualValues(s.T(), index-1, sets.Len())

	s.Keeper().SetLastSlashedRelayerSetNonce(s.Ctx, 1)
	sets = s.Keeper().GetUnSlashedRelayerSets(s.Ctx, uint64(height+index))
	require.EqualValues(s.T(), index-2, sets.Len())

	sets = s.Keeper().GetUnSlashedRelayerSets(s.Ctx, uint64(height+index+1))
	require.EqualValues(s.T(), index-1, sets.Len())
}

func (s *KeeperTestSuite) TestKeeper_IterateRelayerSetConfirmByNonce() {
	index := tmrand.Intn(20) + 1
	for i := uint64(1); i <= uint64(index); i++ {
		for _, relayer := range s.relayerAddrs {
			s.Keeper().SetRelayerSetConfirm(s.Ctx, relayer,
				&types.MsgRelayerSetConfirm{
					Nonce:           i,
					RelayerAddress:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
					ExternalAddress: helpers.GenerateAddress().Hex(),
					Signature:       "",
					ChainName:       s.chainName,
				},
			)
		}
	}

	index = tmrand.Intn(index) + 1
	var confirms []*types.MsgRelayerSetConfirm
	s.Keeper().IterateRelayerSetConfirmByNonce(s.Ctx, uint64(index), func(confirm *types.MsgRelayerSetConfirm) bool {
		confirms = append(confirms, confirm)
		return false
	})
	s.Equal(len(confirms), len(s.relayerAddrs), index)
}

func (s *KeeperTestSuite) TestKeeper_DeleteRelayerSetConfirm() {
	member := make([]types.BridgeValidator, 0, len(s.externalPris))
	for i, external := range s.externalPris {
		externalAddr := s.PubKeyToExternalAddr(external.PublicKey)
		member = append(member, types.BridgeValidator{
			Power:           uint64(i),
			ExternalAddress: externalAddr,
		})
	}
	relayerSet := &types.RelayerSet{
		Nonce:   1,
		Members: member,
		Height:  uint64(s.Ctx.BlockHeight()),
	}
	s.Keeper().StoreRelayerSet(s.Ctx, relayerSet)

	for i, external := range s.externalPris {
		externalAddress, signature := s.SignRelayerSetConfirm(external, relayerSet)
		s.Keeper().SetRelayerSetConfirm(s.Ctx, s.relayerAddrs[i],
			&types.MsgRelayerSetConfirm{
				Nonce:           relayerSet.Nonce,
				RelayerAddress:  s.relayerAddrs[i].String(),
				ExternalAddress: externalAddress,
				Signature:       hex.EncodeToString(signature),
				ChainName:       s.chainName,
			},
		)
	}
	s.Keeper().SetLastObservedRelayerSet(s.Ctx,
		&types.RelayerSet{
			Nonce: relayerSet.Nonce + 1,
		},
	)

	params := s.Keeper().GetParams(s.Ctx)
	params.SignedWindow = 10
	err := s.Keeper().SetParams(s.Ctx, &params)
	s.Require().NoError(err)
	s.Ctx = s.Ctx.WithBlockHeight(1)
	s.Commit()
	for _, relayer := range s.relayerAddrs {
		s.NotNil(s.Keeper().GetRelayerSetConfirm(s.Ctx, relayerSet.Nonce, relayer))
	}

	s.Commit(int64(params.SignedWindow + 1))
	for _, relayer := range s.relayerAddrs {
		s.Nil(s.Keeper().GetRelayerSetConfirm(s.Ctx, relayerSet.Nonce, relayer))
	}
}

func (s *KeeperTestSuite) TestKeeper_IterateRelayerSet() {
	member := make([]types.BridgeValidator, 0, len(s.externalPris))
	for i, external := range s.externalPris {
		member = append(member, types.BridgeValidator{
			Power:           uint64(i),
			ExternalAddress: crypto.PubkeyToAddress(external.PublicKey).String(),
		})
	}
	for i := 1; i <= 10; i++ {
		s.Keeper().StoreRelayerSet(s.Ctx, &types.RelayerSet{
			Nonce:   uint64(i),
			Members: member,
			Height:  uint64(i + 100),
		})
	}
	i := uint64(0)
	relayerSets := types.RelayerSets{}
	s.Keeper().IterateRelayerSetByNonce(s.Ctx, 0, func(relayerSet *types.RelayerSet) bool {
		i = i + 1
		s.Equal(i, relayerSet.Nonce)
		relayerSets = append(relayerSets, relayerSet)
		return false
	})
	s.Equal(len(relayerSets), 10)

	relayerSets = types.RelayerSets{}
	s.Keeper().IterateRelayerSetByNonce(s.Ctx, 1, func(relayerSet *types.RelayerSet) bool {
		relayerSets = append(relayerSets, relayerSet)
		return false
	})
	s.Equal(len(relayerSets), 10)

	relayerSets = types.RelayerSets{}
	s.Keeper().IterateRelayerSetByNonce(s.Ctx, 2, func(relayerSet *types.RelayerSet) bool {
		relayerSets = append(relayerSets, relayerSet)
		return false
	})
	s.Equal(len(relayerSets), 9)

	s.Keeper().IterateRelayerSets(s.Ctx, true, func(relayerSet *types.RelayerSet) bool {
		s.Equal(i, relayerSet.Nonce, relayerSet.Nonce)
		i = i - 1
		return false
	})

	s.Keeper().IterateRelayerSets(s.Ctx, false, func(relayerSet *types.RelayerSet) bool {
		i = i + 1
		s.Equal(i, relayerSet.Nonce)
		return false
	})
}
