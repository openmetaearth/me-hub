package keeper_test

import (
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
)

func (s *KeeperTestSuite) TestQuery_BatchFees() {
	var (
		request  *types.QueryBatchFeeRequest
		response *types.QueryBatchFeeResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"baseFee is negative",
			func() {
				request = &types.QueryBatchFeeRequest{
					ChainName: trontypes.ModuleName,
					MinBatchFees: []types.MinBatchFee{
						{
							TokenContract: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
							BaseFee:       sdkmath.NewInt(-1),
						},
					},
				}
			},
			false,
		},
		{
			"validate tron address error",
			func() {
				request = &types.QueryBatchFeeRequest{
					ChainName: trontypes.ModuleName,
					MinBatchFees: []types.MinBatchFee{
						{
							TokenContract: helpers.GenHexAddress().Hex(),
						},
					},
				}
			},
			false,
		},
		{
			name: "baseFee normal",
			malleate: func() {
				bridgeTokens := s.NewBridgeToken(helpers.GenHexAddress().Bytes())
				minBatchFee := []types.MinBatchFee{
					{
						TokenContract: bridgeTokens[0].ContractAddress,
						BaseFee:       sdk.NewInt(1e7),
					},
				}
				for i := uint64(1); i <= 3; i++ {
					err := s.App.TronKeeper.AttestationHandler(s.Ctx, &gravitytypes.MsgSendToMeClaim{
						TokenContract:  bridgeTokens[0].ContractAddress,
						RelayerAddress: s.signer.AccAddress().String(),
						Amount:         sdk.NewInt(1e8),
						Receiver:       s.signer.AccAddress().String(),
					})
					s.Require().NoError(err)

					_, err = s.App.TronKeeper.AddToOutgoingPool(
						s.Ctx,
						s.signer.AccAddress(),
						helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
						sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(1e7)),
						sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(1e7)),
					)
					s.Require().NoError(err)
				}
				for i := uint64(1); i <= 2; i++ {
					_, err := s.App.TronKeeper.AddToOutgoingPool(
						s.Ctx,
						s.signer.AccAddress(),
						helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
						sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(1e7)),
						sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(1e2)),
					)
					s.Require().NoError(err)
				}
				request = &types.QueryBatchFeeRequest{
					ChainName:    trontypes.ModuleName,
					MinBatchFees: minBatchFee,
				}
				response = &types.QueryBatchFeeResponse{BatchFees: []*types.BatchFees{
					{
						TokenContract: bridgeTokens[0].ContractAddress,
						TotalFees:     sdk.NewInt(1e7 * 3),
						TotalTxs:      3,
						TotalAmount:   sdk.NewInt(1e7 * 3),
					},
				}}
			},
			expPass: true,
		},
		{
			name: "batch fee mul normal",
			malleate: func() {
				bridgeTokens := s.NewBridgeToken(helpers.GenHexAddress().Bytes())
				minBatchFee := []types.MinBatchFee{
					{
						TokenContract: bridgeTokens[0].ContractAddress,
						BaseFee:       sdk.NewInt(1e6),
					},
					{
						TokenContract: bridgeTokens[1].ContractAddress,
						BaseFee:       sdk.NewInt(1e10),
					},
				}
				err := s.App.TronKeeper.AttestationHandler(s.Ctx, &gravitytypes.MsgSendToMeClaim{
					TokenContract:  bridgeTokens[0].ContractAddress,
					RelayerAddress: s.signer.AccAddress().String(),
					Amount:         sdk.NewInt(1e18),
					Receiver:       s.signer.AccAddress().String(),
				})
				s.Require().NoError(err)
				err = s.App.TronKeeper.AttestationHandler(s.Ctx, &gravitytypes.MsgSendToMeClaim{
					TokenContract:  bridgeTokens[1].ContractAddress,
					RelayerAddress: s.signer.AccAddress().String(),
					Amount:         sdk.NewInt(1e18),
					Receiver:       s.signer.AccAddress().String(),
				})
				s.Require().NoError(err)
				for i := 1; i <= 2; i++ {
					_, err = s.App.TronKeeper.AddToOutgoingPool(
						s.Ctx,
						s.signer.AccAddress(),
						helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
						sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(1e6)),
						sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(1e5)))
					s.Require().NoError(err)
				}
				_, err = s.App.TronKeeper.AddToOutgoingPool(
					s.Ctx,
					s.signer.AccAddress(),
					helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(1e6)),
					sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(1e6)))
				s.Require().NoError(err)

				for i := 1; i <= 3; i++ {
					_, err = s.App.TronKeeper.AddToOutgoingPool(
						s.Ctx,
						s.signer.AccAddress(),
						helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
						sdk.NewCoin(bridgeTokens[1].Denom, sdk.NewInt(1e10)),
						sdk.NewCoin(bridgeTokens[1].Denom, sdk.NewInt(1e10)))
					s.Require().NoError(err)
				}
				request = &types.QueryBatchFeeRequest{
					ChainName:    trontypes.ModuleName,
					MinBatchFees: minBatchFee,
				}
				response = &types.QueryBatchFeeResponse{BatchFees: []*types.BatchFees{
					{
						TokenContract: bridgeTokens[0].ContractAddress,
						TotalFees:     sdk.NewInt(1e6),
						TotalTxs:      1,
						TotalAmount:   sdk.NewInt(1e6),
					},
					{
						TokenContract: bridgeTokens[1].ContractAddress,
						TotalFees:     sdk.NewInt(1e10 * 3),
						TotalTxs:      3,
						TotalAmount:   sdk.NewInt(1e10 * 3),
					},
				}}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			s.SetupTest()
			ctx := sdk.WrapSDKContext(s.Ctx)
			testCase.malleate()
			res, err := s.queryServer.BatchFees(ctx, request)
			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().ElementsMatch(response.BatchFees, res.BatchFees)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuery_BatchRequestByNonce() {
	var (
		request  *types.QueryBatchRequestByNonceRequest
		response *types.QueryBatchRequestByNonceResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "store normal batch",
			malleate: func() {
				bridgeTokens := s.NewBridgeToken(helpers.GenHexAddress().Bytes())
				feeReceive := helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex())
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: bridgeTokens[0].ContractAddress,
					Nonce:         3,
				}
				err := s.App.TronKeeper.StoreBatch(s.Ctx, &types.OutgoingTxBatch{
					BatchNonce:   3,
					BatchTimeout: 10000,
					Transactions: []*types.OutgoingTransferTx{
						{
							Token: types.ERC20Token{
								Contract: bridgeTokens[0].ContractAddress,
								Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
							},
							Fee: types.ERC20Token{
								Contract: bridgeTokens[0].ContractAddress,
								Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
							},
						},
					},
					TokenContract: bridgeTokens[0].ContractAddress,
					FeeReceive:    feeReceive,
				})
				s.Require().NoError(err)
				response = &types.QueryBatchRequestByNonceResponse{
					Batch: &types.OutgoingTxBatch{
						BatchNonce:   3,
						BatchTimeout: 10000,
						Transactions: []*types.OutgoingTransferTx{
							{
								Token: types.ERC20Token{
									Contract: bridgeTokens[0].ContractAddress,
									Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
								},
								Fee: types.ERC20Token{
									Contract: bridgeTokens[0].ContractAddress,
									Amount:   sdkmath.NewIntFromBigInt(big.NewInt(1e18)),
								},
							},
						},
						TokenContract: bridgeTokens[0].ContractAddress,
						FeeReceive:    feeReceive,
					},
				}
			},
			expPass: true,
		},
		{
			name: "request error nonce",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					Nonce:         0,
				}
			},
			expPass: false,
		},
		{
			name: "request error token",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.GenHexAddress().Hex(),
					Nonce:         8,
				}
			},
			expPass: false,
		},
		{
			name: "request nonexistent nonce",
			malleate: func() {
				bridgeTokens := s.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: bridgeTokens[0].ContractAddress,
					Nonce:         8,
				}
			},
			expPass: false,
		},
		{
			name: "request nonexistent token",
			malleate: func() {
				request = &types.QueryBatchRequestByNonceRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					Nonce:         8,
				}
			},
			expPass: false,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			s.SetupTest()
			testCase.malleate()
			res, err := s.queryServer.BatchRequestByNonce(sdk.WrapSDKContext(s.Ctx), request)
			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().Equal(response.Batch, res.Batch)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuery_BatchConfirms() {
	var (
		request  *types.QueryBatchConfirmsRequest
		response *types.QueryBatchConfirmsResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"token address error",
			func() {
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.GenHexAddress().Hex(),
				}
			},
			false,
		},
		{
			"token nonce is zero",
			func() {
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
					Nonce:         0,
				}
			},
			false,
		},
		{
			name: "request confirm nonexistent nonce",
			malleate: func() {
				bridgeTokens := s.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: bridgeTokens[0].ContractAddress,
					Nonce:         2,
				}
				s.App.TronKeeper.SetBatchConfirm(s.Ctx, s.signer.AccAddress(), &types.MsgConfirmBatch{
					Nonce: 1,
				})
				response = &types.QueryBatchConfirmsResponse{}
			},
			expPass: true,
		},
		{
			"set correct batch confirm",
			func() {
				relayer, externalKey := s.NewRelayer()
				bridgeTokens := s.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryBatchConfirmsRequest{
					ChainName:     trontypes.ModuleName,
					TokenContract: bridgeTokens[0].ContractAddress,
					Nonce:         1,
				}
				newConfirmBatch := &types.MsgConfirmBatch{
					ChainName:       trontypes.ModuleName,
					Nonce:           1,
					TokenContract:   bridgeTokens[0].ContractAddress,
					RelayerAddress:  relayer.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					Signature:       helpers.GenHexAddress().Hex(),
				}
				s.App.TronKeeper.SetBatchConfirm(s.Ctx, s.signer.AccAddress(), newConfirmBatch)
				response = &types.QueryBatchConfirmsResponse{Confirms: []*types.MsgConfirmBatch{newConfirmBatch}}
			},
			true,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			s.SetupTest()

			ctx := sdk.WrapSDKContext(s.Ctx)
			testCase.malleate()

			res, err := s.queryServer.BatchConfirms(ctx, request)

			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().ElementsMatch(response.Confirms, res.Confirms)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuery_BridgeToken() {
	var (
		request  *types.QueryBridgeTokenRequest
		response types.QueryBridgeTokenResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "token address error",
			malleate: func() {
				request = &types.QueryBridgeTokenRequest{
					ChainName:       trontypes.ModuleName,
					ContractAddress: helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "token that does not exist",
			malleate: func() {
				request = &types.QueryBridgeTokenRequest{
					ChainName:       trontypes.ModuleName,
					ContractAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "token normal",
			malleate: func() {
				bridgeTokens := s.NewBridgeToken(helpers.GenHexAddress().Bytes())
				request = &types.QueryBridgeTokenRequest{
					ChainName:       trontypes.ModuleName,
					ContractAddress: bridgeTokens[0].ContractAddress,
				}
				response = types.QueryBridgeTokenResponse{
					BridgeToken: &bridgeTokens[0],
					TotalSupply: sdk.NewCoin(bridgeTokens[0].Denom, sdk.NewInt(0)),
				}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			s.SetupTest()
			ctx := sdk.WrapSDKContext(s.Ctx)
			testCase.malleate()
			res, err := s.queryServer.BridgeToken(ctx, request)
			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().EqualValues(response.BridgeToken.ContractAddress, res.BridgeToken.ContractAddress)
				s.Require().EqualValues(response.BridgeToken.Denom, res.BridgeToken.Denom)
				s.Require().EqualValues(response.BridgeToken.Name, res.BridgeToken.Name)
				s.Require().EqualValues(response.BridgeToken.Symbol, res.BridgeToken.Symbol)
				s.Require().EqualValues(response.BridgeToken.Decimal, res.BridgeToken.Decimal)
				s.Require().EqualValues(response.TotalSupply, res.TotalSupply)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuery_GetRelayerByExternalAddr() {
	var (
		request  *types.QueryRelayerRequest
		response *types.QueryRelayerResponse
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name: "external address is error",
			malleate: func() {
				request = &types.QueryRelayerRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.GenHexAddress().Hex(),
				}
			},
			expPass: false,
		},
		{
			name: "nonexistent external address",
			malleate: func() {
				request = &types.QueryRelayerRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.HexAddrToTronAddr(helpers.GenHexAddress().Hex()),
				}
			},
			expPass: false,
		},
		{
			name: "normal external address",
			malleate: func() {
				bridger, externalKey := s.NewRelayer()
				request = &types.QueryRelayerRequest{
					ChainName:       trontypes.ModuleName,
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
				}
				response = &types.QueryRelayerResponse{Relayer: &types.Relayer{
					RelayerAddress:  bridger.String(),
					ExternalAddress: helpers.HexAddrToTronAddr(externalKey.PubKey().Address().String()),
					DelegateAmount:  sdkmath.ZeroInt(),
				}}
			},
			expPass: true,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			s.SetupTest()
			testCase.malleate()
			res, err := s.queryServer.Relayer(sdk.WrapSDKContext(s.Ctx), request)
			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().EqualValues(response.Relayer.RelayerAddress, res.Relayer.RelayerAddress)
				s.Require().EqualValues(response.Relayer.ExternalAddress, res.Relayer.ExternalAddress)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
