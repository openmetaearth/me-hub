package keeper_test

import (
	"encoding/hex"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/openmetaearth/me-hub/testutil/helpers"

	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
)

func (s *KeeperTestSuite) Test_msgServer_ConfirmBatch() {
	var msg *gravitytypes.MsgConfirmBatch
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "couldn't find batch",
			malleate: func() {
				msg = &gravitytypes.MsgConfirmBatch{
					Nonce:          tmrand.Uint64(),
					TokenContract:  helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					RelayerAddress: helpers.GenAccAddress().String(),
				}
			},
			expPass: false,
		},
		{
			name: "not found relayer",
			malleate: func() {
				newOutgoingTx := s.NewOutgoingTxBatch()
				msg = &gravitytypes.MsgConfirmBatch{
					Nonce:          newOutgoingTx.BatchNonce,
					TokenContract:  newOutgoingTx.TokenContract,
					RelayerAddress: helpers.GenAccAddress().String(),
				}
			},
			expPass: false,
		},
		{
			name: "signature decoding failed",
			malleate: func() {
				newOutgoingTx := s.NewOutgoingTxBatch()
				relayer, externalKey := s.NewRelayer()
				msg = &gravitytypes.MsgConfirmBatch{
					Nonce:           newOutgoingTx.BatchNonce,
					TokenContract:   newOutgoingTx.TokenContract,
					RelayerAddress:  relayer.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "confirm batch",
			malleate: func() {
				newOutgoingTx := s.NewOutgoingTxBatch()
				relayer, externalKey := s.NewRelayer()
				params, err := s.queryServer.Params(sdk.WrapSDKContext(s.Ctx), &gravitytypes.QueryParamsRequest{ChainName: trontypes.ModuleName})
				s.Require().NoError(err)
				batchHash, err := trontypes.GetCheckpointConfirmBatch(newOutgoingTx, params.Params.GravityId)
				s.Require().NoError(err)
				key, err := externalKey.(*ethsecp256k1.PrivKey).ToECDSA()
				s.Require().NoError(err)
				signature, err := trontypes.NewTronSignature(batchHash, key)
				s.Require().NoError(err)
				msg = &gravitytypes.MsgConfirmBatch{
					Nonce:           newOutgoingTx.BatchNonce,
					TokenContract:   newOutgoingTx.TokenContract,
					RelayerAddress:  relayer.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       hex.EncodeToString(signature),
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			testCase.malleate()
			_, err := s.msgServer.ConfirmBatch(sdk.WrapSDKContext(s.Ctx), msg)
			if testCase.expPass {
				s.Require().NoError(err)
			} else {
				s.Require().ErrorContains(err, testCase.name)
			}
		})
	}
}

func (s *KeeperTestSuite) Test_msgServer_RelayerSetConfirm() {
	var msg *gravitytypes.MsgRelayerSetConfirm
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "couldn't find relayerSet",
			malleate: func() {
				msg = &gravitytypes.MsgRelayerSetConfirm{
					Nonce:          tmrand.Uint64(),
					RelayerAddress: helpers.GenAccAddress().String(),
				}
			},
			expPass: false,
		},
		{
			name: "not found relayer",
			malleate: func() {
				currentRelayerSet := s.CurrentRelayerSet(helpers.NewEthPrivKey())
				msg = &gravitytypes.MsgRelayerSetConfirm{
					Nonce:           currentRelayerSet.Nonce,
					RelayerAddress:  helpers.GenAccAddress().String(),
					ExternalAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "signature decoding",
			malleate: func() {
				relayer, externalKey := s.NewRelayer()
				currentRelayerSet := s.CurrentRelayerSet(externalKey)
				msg = &gravitytypes.MsgRelayerSetConfirm{
					Nonce:           currentRelayerSet.Nonce,
					RelayerAddress:  relayer.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					Signature:       helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "relayer set confirm",
			malleate: func() {
				relayer, externalKey := s.NewRelayer()
				currentRelayerSet := s.CurrentRelayerSet(externalKey)
				key, err := externalKey.(*ethsecp256k1.PrivKey).ToECDSA()
				s.Require().NoError(err)
				params, err := s.queryServer.Params(sdk.WrapSDKContext(s.Ctx), &gravitytypes.QueryParamsRequest{ChainName: trontypes.ModuleName})
				s.Require().NoError(err)
				relayerSetHash, err := trontypes.GetCheckpointRelayerSet(currentRelayerSet, params.Params.GravityId)
				s.Require().NoError(err)
				signature, err := trontypes.NewTronSignature(relayerSetHash, key)
				s.Require().NoError(err)
				msg = &gravitytypes.MsgRelayerSetConfirm{
					Nonce:           currentRelayerSet.Nonce,
					RelayerAddress:  relayer.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       hex.EncodeToString(signature),
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			s.SetupTest()
			testCase.malleate()
			_, err := s.msgServer.RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), msg)
			if testCase.expPass {
				s.Require().NoError(err)
			} else {
				s.Require().ErrorContains(err, testCase.name)
			}
		})
	}
}
