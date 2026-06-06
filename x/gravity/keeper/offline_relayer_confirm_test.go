package keeper_test

import (
	"encoding/hex"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestOfflineRelayerCannotSubmitRelayerSetConfirm() {
	relayerAddr := s.relayerAddrs[0]
	externalAddr := s.PubKeyToExternalAddr(s.externalPris[0].PublicKey)
	msg := &types.MsgBondedRelayer{
		RelayerAddress:  relayerAddr.String(),
		ExternalAddress: externalAddr,
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, s.Keeper().GetGravityMinDelegate(s.Ctx)),
		ChainName:       s.chainName,
	}
	_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msg)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.Keeper().EndBlocker(s.Ctx)

	relayerSet := s.Keeper().GetRelayerSet(s.Ctx, 1)
	s.Require().NotNil(relayerSet)
	checkpoint, err := relayerSet.GetCheckpoint(s.Keeper().GetGravityID(s.Ctx))
	s.Require().NoError(err)
	signature, err := types.NewEthereumSignature(checkpoint, s.externalPris[0])
	s.Require().NoError(err)

	relayer, found := s.Keeper().GetRelayer(s.Ctx, relayerAddr)
	s.Require().True(found)
	relayer.Online = false
	s.Keeper().SetRelayer(s.Ctx, relayerAddr, relayer)

	confirm := &types.MsgRelayerSetConfirm{
		Nonce:           relayerSet.Nonce,
		RelayerAddress:  relayerAddr.String(),
		ExternalAddress: externalAddr,
		Signature:       hex.EncodeToString(signature),
		ChainName:       s.chainName,
	}
	_, err = s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), confirm)
	s.Require().ErrorIs(err, types.ErrRelayerNotOnLine)
	s.Require().Nil(s.Keeper().GetRelayerSetConfirm(s.Ctx, relayerSet.Nonce, relayerAddr))
}

func (s *KeeperTestSuite) TestOfflineRelayerCannotSubmitBatchConfirm() {
	s.setupBondedRelayerSetForAuditTest()

	relayerAddr := s.relayerAddrs[0]
	externalAddr := s.PubKeyToExternalAddr(s.externalPris[0].PublicKey)
	bridgeToken := s.NewBridgeToken(relayerAddr, sdk.NewCoin("uofflineconfirm", sdkmath.NewInt(100)))

	_, err := s.MsgServer().SendToExternal(sdk.WrapSDKContext(s.Ctx), &types.MsgSendToExternal{
		Sender:    relayerAddr.String(),
		Dest:      externalAddr,
		Amount:    sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(10)),
		BridgeFee: sdk.NewCoin(bridgeToken.Denom, sdkmath.NewInt(1)),
		ChainName: s.chainName,
	})
	s.Require().NoError(err)

	s.Keeper().SetLastObservedBlockHeight(s.Ctx, 1_000, uint64(s.Ctx.BlockHeight()))
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	batchResponse, err := s.MsgServer().RequestBatch(sdk.WrapSDKContext(s.Ctx), &types.MsgRequestBatch{
		Sender:     relayerAddr.String(),
		Denom:      bridgeToken.Denom,
		MinimumFee: sdkmath.NewInt(1),
		FeeReceive: helpers.GenerateAddress().Hex(),
		ChainName:  s.chainName,
		BaseFee:    sdkmath.ZeroInt(),
	})
	s.Require().NoError(err)

	batch := s.Keeper().GetOutgoingTxBatch(s.Ctx, bridgeToken.ContractAddress, batchResponse.BatchNonce)
	s.Require().NotNil(batch)
	checkpoint, err := batch.GetCheckpoint(s.Keeper().GetGravityID(s.Ctx))
	s.Require().NoError(err)
	signature, err := types.NewEthereumSignature(checkpoint, s.externalPris[0])
	s.Require().NoError(err)

	relayer, found := s.Keeper().GetRelayer(s.Ctx, relayerAddr)
	s.Require().True(found)
	relayer.Online = false
	s.Keeper().SetRelayer(s.Ctx, relayerAddr, relayer)

	confirm := &types.MsgConfirmBatch{
		Nonce:           batch.BatchNonce,
		TokenContract:   batch.TokenContract,
		RelayerAddress:  relayerAddr.String(),
		ExternalAddress: externalAddr,
		Signature:       hex.EncodeToString(signature),
		ChainName:       s.chainName,
	}
	_, err = s.MsgServer().ConfirmBatch(sdk.WrapSDKContext(s.Ctx), confirm)
	s.Require().ErrorIs(err, types.ErrRelayerNotOnLine)
	s.Require().Nil(s.Keeper().GetBatchConfirm(s.Ctx, batch.TokenContract, batch.BatchNonce, relayerAddr))
}

func (s *KeeperTestSuite) TestOfflineRelayerConfirmsAreFilteredFromQueries() {
	queryServer := keeper.NewQueryServerImpl(s.Keeper())
	onlineRelayer := s.relayerAddrs[0]
	offlineRelayer := s.relayerAddrs[1]
	tokenContract := helpers.GenerateAddress().Hex()

	s.Keeper().SetRelayer(s.Ctx, onlineRelayer, types.Relayer{
		RelayerAddress:  onlineRelayer.String(),
		ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
		Online:          true,
	})
	s.Keeper().SetRelayer(s.Ctx, offlineRelayer, types.Relayer{
		RelayerAddress:  offlineRelayer.String(),
		ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[1].PublicKey),
		Online:          false,
	})

	s.Keeper().SetRelayerSetConfirm(s.Ctx, onlineRelayer, &types.MsgRelayerSetConfirm{
		Nonce:          1,
		RelayerAddress: onlineRelayer.String(),
		ChainName:      s.chainName,
	})
	s.Keeper().SetRelayerSetConfirm(s.Ctx, offlineRelayer, &types.MsgRelayerSetConfirm{
		Nonce:          1,
		RelayerAddress: offlineRelayer.String(),
		ChainName:      s.chainName,
	})
	s.Keeper().SetBatchConfirm(s.Ctx, onlineRelayer, &types.MsgConfirmBatch{
		Nonce:          1,
		TokenContract:  tokenContract,
		RelayerAddress: onlineRelayer.String(),
		ChainName:      s.chainName,
	})
	s.Keeper().SetBatchConfirm(s.Ctx, offlineRelayer, &types.MsgConfirmBatch{
		Nonce:          1,
		TokenContract:  tokenContract,
		RelayerAddress: offlineRelayer.String(),
		ChainName:      s.chainName,
	})

	relayerSetConfirms, err := queryServer.RelayerSetConfirmsByNonce(s.Ctx, &types.QueryRelayerSetConfirmsByNonceRequest{
		ChainName: s.chainName,
		Nonce:     1,
	})
	s.Require().NoError(err)
	s.Require().Len(relayerSetConfirms.Confirms, 1)
	s.Require().Equal(onlineRelayer.String(), relayerSetConfirms.Confirms[0].RelayerAddress)

	onlineRelayerSetConfirm, err := queryServer.RelayerSetConfirm(s.Ctx, &types.QueryRelayerSetConfirmRequest{
		ChainName:      s.chainName,
		RelayerAddress: onlineRelayer.String(),
		Nonce:          1,
	})
	s.Require().NoError(err)
	s.Require().NotNil(onlineRelayerSetConfirm.Confirm)
	s.Require().Equal(onlineRelayer.String(), onlineRelayerSetConfirm.Confirm.RelayerAddress)

	offlineRelayerSetConfirm, err := queryServer.RelayerSetConfirm(s.Ctx, &types.QueryRelayerSetConfirmRequest{
		ChainName:      s.chainName,
		RelayerAddress: offlineRelayer.String(),
		Nonce:          1,
	})
	s.Require().NoError(err)
	s.Require().Nil(offlineRelayerSetConfirm.Confirm)

	batchConfirms, err := queryServer.BatchConfirms(s.Ctx, &types.QueryBatchConfirmsRequest{
		ChainName:     s.chainName,
		TokenContract: tokenContract,
		Nonce:         1,
	})
	s.Require().NoError(err)
	s.Require().Len(batchConfirms.Confirms, 1)
	s.Require().Equal(onlineRelayer.String(), batchConfirms.Confirms[0].RelayerAddress)

	onlineBatchConfirm, err := queryServer.BatchConfirm(s.Ctx, &types.QueryBatchConfirmRequest{
		ChainName:      s.chainName,
		TokenContract:  tokenContract,
		RelayerAddress: onlineRelayer.String(),
		Nonce:          1,
	})
	s.Require().NoError(err)
	s.Require().NotNil(onlineBatchConfirm.Confirm)
	s.Require().Equal(onlineRelayer.String(), onlineBatchConfirm.Confirm.RelayerAddress)

	offlineBatchConfirm, err := queryServer.BatchConfirm(s.Ctx, &types.QueryBatchConfirmRequest{
		ChainName:      s.chainName,
		TokenContract:  tokenContract,
		RelayerAddress: offlineRelayer.String(),
		Nonce:          1,
	})
	s.Require().NoError(err)
	s.Require().Nil(offlineBatchConfirm.Confirm)
}
