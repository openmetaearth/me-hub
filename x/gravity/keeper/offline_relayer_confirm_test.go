package keeper_test

import (
	"encoding/hex"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
)

func (s *KeeperTestSuite) TestOfflineRelayerCannotSubmitBatchConfirm() {
	relayerMsg := s.bondRelayerForOfflineConfirmTest(0)
	tokenContract := helpers.GenerateAddress().Hex()
	batch := &types.OutgoingTxBatch{
		BatchNonce:    1,
		BatchTimeout:  100,
		TokenContract: tokenContract,
		FeeReceive:    helpers.GenerateAddress().Hex(),
		Transactions: []*types.OutgoingTransferTx{
			{
				Id:          1,
				Sender:      s.relayerAddrs[1].String(),
				DestAddress: helpers.GenerateAddress().Hex(),
				Token: types.ERC20Token{
					Contract: tokenContract,
					Amount:   sdkmath.NewInt(100),
				},
				Fee: types.ERC20Token{
					Contract: tokenContract,
					Amount:   sdkmath.NewInt(1),
				},
			},
		},
	}
	s.Require().NoError(s.Keeper().StoreBatch(s.Ctx, batch))

	checkpoint, err := batch.GetCheckpoint(s.Keeper().GetGravityID(s.Ctx))
	s.Require().NoError(err)
	signature, err := types.NewEthereumSignature(checkpoint, s.externalPris[0])
	if trontypes.ModuleName == s.chainName {
		signature, err = trontypes.NewTronSignature(checkpoint, s.externalPris[0])
	}
	s.Require().NoError(err)

	s.setRelayerOffline(0)

	_, err = s.MsgServer().ConfirmBatch(sdk.WrapSDKContext(s.Ctx), &types.MsgConfirmBatch{
		Nonce:           batch.BatchNonce,
		TokenContract:   tokenContract,
		RelayerAddress:  relayerMsg.RelayerAddress,
		ExternalAddress: relayerMsg.ExternalAddress,
		Signature:       hex.EncodeToString(signature),
		ChainName:       s.chainName,
	})

	s.Require().ErrorIs(err, types.ErrRelayerNotOnLine)
	s.Require().Nil(s.Keeper().GetBatchConfirm(s.Ctx, tokenContract, batch.BatchNonce, s.relayerAddrs[0]))
}

func (s *KeeperTestSuite) TestOfflineRelayerCannotSubmitRelayerSetConfirm() {
	relayerMsg := s.bondRelayerForOfflineConfirmTest(0)
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.Keeper().EndBlocker(s.Ctx)

	relayerSet := s.Keeper().GetRelayerSet(s.Ctx, 1)
	s.Require().NotNil(relayerSet)
	externalAddress, signature := s.SignRelayerSetConfirm(s.externalPris[0], relayerSet)

	s.setRelayerOffline(0)

	_, err := s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), &types.MsgRelayerSetConfirm{
		Nonce:           relayerSet.Nonce,
		RelayerAddress:  relayerMsg.RelayerAddress,
		ExternalAddress: externalAddress,
		Signature:       hex.EncodeToString(signature),
		ChainName:       s.chainName,
	})

	s.Require().ErrorIs(err, types.ErrRelayerNotOnLine)
	s.Require().Nil(s.Keeper().GetRelayerSetConfirm(s.Ctx, relayerSet.Nonce, s.relayerAddrs[0]))
}

func (s *KeeperTestSuite) bondRelayerForOfflineConfirmTest(index int) *types.MsgBondedRelayer {
	msg := &types.MsgBondedRelayer{
		RelayerAddress:  s.relayerAddrs[index].String(),
		ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[index].PublicKey),
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(10*1e8)),
		ChainName:       s.chainName,
	}
	_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msg)
	s.Require().NoError(err)
	return msg
}

func (s *KeeperTestSuite) setRelayerOffline(index int) {
	relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[index])
	s.Require().True(found)
	relayer.Online = false
	s.Keeper().SetRelayer(s.Ctx, s.relayerAddrs[index], relayer)
}
