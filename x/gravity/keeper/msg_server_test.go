package keeper_test

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"encoding/hex"
	"fmt"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
	"sort"
)

func (s *KeeperTestSuite) TestMsgBondedRelayer() {
	testCases := []struct {
		name   string
		pass   bool
		err    string
		preRun func(msg *types.MsgBondedRelayer)
	}{
		{
			name: "error - not a proposed relayer",
			preRun: func(msg *types.MsgBondedRelayer) {
				newAddress, _ := s.NewAccount()
				msg.RelayerAddress = newAddress.String()
			},
			pass: false,
			err:  types.ErrNotProposedRelayer.Error(),
		},
		{
			name: "error - relayer existed",
			preRun: func(msg *types.MsgBondedRelayer) {
				s.Keeper().SetRelayer(s.Ctx, sdk.MustAccAddressFromBech32(msg.RelayerAddress), types.Relayer{RelayerAddress: msg.RelayerAddress})
			},
			pass: false,
			err:  "relayer existed bridger address: invalid",
		},
		{
			name: "error - external address is bound",
			preRun: func(msg *types.MsgBondedRelayer) {
				s.Keeper().SetRelayerByExternalAddress(s.Ctx, msg.ExternalAddress, sdk.MustAccAddressFromBech32(msg.RelayerAddress))
			},
			pass: false,
			err:  "external address is bound to relayer: invalid",
		},
		{
			name: "error - stake denom not match chain params stake denom",
			preRun: func(msg *types.MsgBondedRelayer) {
				msg.DelegateAmount.Denom = "stake"
			},
			pass: false,
			err:  fmt.Sprintf("delegate denom got %s, expected %s: invalid", "stake", "umec"),
		},
		{
			name: "error - delegate amount less than threshold amount",
			preRun: func(msg *types.MsgBondedRelayer) {
				delegateThreshold := s.Keeper().GetGravityMinDelegate(s.Ctx)
				msg.DelegateAmount.Amount = delegateThreshold.Sub(sdk.NewInt(1))
			},
			pass: false,
			err:  types.ErrDelegateAmountBelowMinimum.Error(),
		},
		{
			name: "error - delegate amount grate than threshold amount",
			preRun: func(msg *types.MsgBondedRelayer) {
				maxDelegateAmount := s.Keeper().GetGravityMaxDelegate(s.Ctx)
				msg.DelegateAmount.Amount = maxDelegateAmount.Add(sdk.NewInt(1))
			},
			pass: false,
			err:  types.ErrDelegateAmountAboveMaximum.Error(),
		},
		{
			name: "pass",
			preRun: func(msg *types.MsgBondedRelayer) {
			},
			pass: true,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			relayerIndex := tmrand.Intn(len(s.relayerAddrs))
			msg := &types.MsgBondedRelayer{
				RelayerAddress:  s.relayerAddrs[relayerIndex].String(),
				ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[relayerIndex].PublicKey),
				DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
				ChainName:       s.chainName,
			}

			testCase.preRun(msg)

			_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msg)
			if !testCase.pass {
				s.Require().Error(err)
				s.Require().EqualValues(testCase.err, err.Error())
				return
			}

			// success check
			s.Require().NoError(err)

			// check relayer
			relayer, found := s.Keeper().GetRelayer(s.Ctx, sdk.MustAccAddressFromBech32(msg.RelayerAddress))
			s.Require().True(found)
			s.Require().NotNil(relayer)
			s.Require().EqualValues(msg.RelayerAddress, relayer.RelayerAddress)
			s.Require().EqualValues(msg.ExternalAddress, relayer.ExternalAddress)
			s.Require().True(relayer.Online)
			s.Require().EqualValues(int64(0), relayer.SlashTimes)
			s.Require().EqualValues(msg.RelayerAddress, relayer.RelayerAddress)

			relayerAddr, found := s.Keeper().GetRelayerByExternalAddress(s.Ctx, msg.ExternalAddress)
			s.True(found)
			s.Require().EqualValues(msg.RelayerAddress, relayerAddr.String())

			// check power
			totalPower := s.Keeper().GetLastTotalPower(s.Ctx)
			s.Require().EqualValues(msg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction).Int64(), totalPower.Int64())
		})
	}
}

func (s *KeeperTestSuite) TestMsgAddDelegate() {
	initDelegateAmount := s.Keeper().GetGravityMinDelegate(s.Ctx)
	testCases := []struct {
		name                 string
		pass                 bool
		err                  string
		preRun               func(msg *types.MsgAddDelegate)
		expectDelegateAmount func(msg *types.MsgAddDelegate) sdk.Int
	}{
		{
			name: "error - sender not relayer",
			preRun: func(msg *types.MsgAddDelegate) {
				msg.RelayerAddress = sdk.AccAddress(tmrand.Bytes(20)).String()
			},
			pass: false,
			err:  types.ErrNotProposedRelayer.Error(),
		},
		{
			name: "error - stake denom not match chain params stake denom",
			preRun: func(msg *types.MsgAddDelegate) {
				msg.Amount.Denom = "stake"
			},
			pass: false,
			err:  fmt.Sprintf("delegate denom got %s, expected %s: invalid", "stake", "umec"),
		},
		{
			name: "error - delegate amount less than threshold amount",
			preRun: func(msg *types.MsgAddDelegate) {
				params := s.Keeper().GetParams(s.Ctx)
				addDelegateThreshold := tmrand.Int63n(100000) + 1
				params.MinDelegate = initDelegateAmount.Add(sdk.NewInt(addDelegateThreshold))
				err := s.Keeper().SetParams(s.Ctx, &params)
				s.Require().NoError(err)
				msg.Amount.Amount = sdk.NewInt(tmrand.Int63n(addDelegateThreshold))
			},
			pass: false,
			err:  types.ErrDelegateAmountBelowMinimum.Error(),
		},
		{
			name: "error - delegate amount greater than threshold amount",
			preRun: func(msg *types.MsgAddDelegate) {
				maxDelegateAmount := s.Keeper().GetGravityMaxDelegate(s.Ctx)
				msg.Amount.Amount = maxDelegateAmount.Add(sdk.NewInt(1))
			},
			pass: false,
			err:  types.ErrDelegateAmountAboveMaximum.Error(),
		},
		{
			name: "pass",
			preRun: func(msg *types.MsgAddDelegate) {
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdk.Int {
				return initDelegateAmount.Add(msg.Amount.Amount)
			},
		},
		{
			name: "error - not sufficient slash amount",
			preRun: func(msg *types.MsgAddDelegate) {
				relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
				relayer, _ := s.Keeper().GetRelayer(s.Ctx, relayerAddress)
				relayer.SlashTimes = 1
				s.Keeper().SetRelayer(s.Ctx, relayerAddress, relayer)
				slashFraction := s.Keeper().GetSlashFraction(s.Ctx)
				slashAmount := sdk.NewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(relayer.SlashTimes).TruncateInt()
				randomAmount := tmrand.Int63n(slashAmount.Int64()) + 1
				msg.Amount.Amount = sdk.NewInt(randomAmount)
			},
			pass: false,
			err:  "not sufficient slash amount: invalid",
		},
		{
			name: "pass - add slash amount",
			preRun: func(msg *types.MsgAddDelegate) {
				relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
				relayer, _ := s.Keeper().GetRelayer(s.Ctx, relayerAddress)
				relayer.SlashTimes = 1
				relayer.Online = false
				s.Keeper().SetRelayer(s.Ctx, relayerAddress, relayer)

				slashFraction := s.Keeper().GetSlashFraction(s.Ctx)
				slashAmount := sdk.NewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(relayer.SlashTimes).TruncateInt()
				msg.Amount.Amount = slashAmount
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdk.Int {
				return initDelegateAmount
			},
		},
		{
			name: "pass - add more slash amount",
			preRun: func(msg *types.MsgAddDelegate) {
				relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
				relayer, _ := s.Keeper().GetRelayer(s.Ctx, relayerAddress)
				relayer.SlashTimes = 1
				relayer.Online = false
				s.Keeper().SetRelayer(s.Ctx, relayerAddress, relayer)

				slashFraction := s.Keeper().GetSlashFraction(s.Ctx)
				slashAmount := sdk.NewDecFromInt(initDelegateAmount).Mul(slashFraction).MulInt64(relayer.SlashTimes).TruncateInt()
				msg.Amount.Amount = slashAmount.Add(sdk.NewInt(1000))
			},
			pass: true,
			expectDelegateAmount: func(msg *types.MsgAddDelegate) sdk.Int {
				return initDelegateAmount.Add(sdk.NewInt(1000))
			},
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			relayerIndex := tmrand.Intn(len(s.relayerAddrs))

			// init bonded relayer
			_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), &types.MsgBondedRelayer{
				RelayerAddress:  s.relayerAddrs[relayerIndex].String(),
				ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[relayerIndex].PublicKey),
				DelegateAmount:  sdk.NewCoin(params.BaseDenom, initDelegateAmount),
				ChainName:       s.chainName,
			})
			s.Require().NoError(err)

			msg := &types.MsgAddDelegate{
				ChainName:      s.chainName,
				RelayerAddress: s.relayerAddrs[relayerIndex].String(),
				Amount:         sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
			}
			testCase.preRun(msg)

			_, err = s.MsgServer().AddDelegate(sdk.WrapSDKContext(s.Ctx), msg)
			if !testCase.pass {
				s.Require().Error(err)
				s.Require().EqualValues(testCase.err, err.Error())
				return
			}
			s.Require().NoError(err)

			// check relayer
			relayer, found := s.Keeper().GetRelayer(s.Ctx, sdk.MustAccAddressFromBech32(msg.RelayerAddress))
			s.Require().True(found)
			s.Require().NotNil(relayer)
			s.Require().EqualValues(msg.RelayerAddress, relayer.RelayerAddress)
			s.Require().True(relayer.Online)
			s.Require().EqualValues(0, relayer.SlashTimes)

			// check power
			totalPower := s.Keeper().GetLastTotalPower(s.Ctx)
			expectDelegateAmount := testCase.expectDelegateAmount(msg)
			s.Require().EqualValues(expectDelegateAmount.Quo(sdk.DefaultPowerReduction).Int64(), totalPower.Int64())
		})
	}
}

func (s *KeeperTestSuite) TestMsgSetRelayerSetConfirm() {
	normalMsg := &types.MsgBondedRelayer{
		RelayerAddress:  s.relayerAddrs[0].String(),
		ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
		ChainName:       s.chainName,
	}
	_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), normalMsg)
	s.Require().NoError(err)

	latestRelayerSetNonce := s.Keeper().GetLastRelayerSetNonce(s.Ctx)
	s.Require().EqualValues(0, latestRelayerSetNonce)

	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.Keeper().EndBlocker(s.Ctx)

	latestRelayerSetNonce = s.Keeper().GetLastRelayerSetNonce(s.Ctx)
	s.Require().EqualValues(1, latestRelayerSetNonce)

	s.Require().True(s.Keeper().HasRelayerSetRequest(s.Ctx, 1))
	s.Require().False(s.Keeper().HasRelayerSetRequest(s.Ctx, 2))

	nonce1RelayerSet := s.Keeper().GetRelayerSet(s.Ctx, 1)
	s.Require().EqualValues(uint64(1), nonce1RelayerSet.Nonce)
	s.Require().EqualValues(1, nonce1RelayerSet.Height)
	s.Require().EqualValues(1, len(nonce1RelayerSet.Members))
	s.Require().EqualValues(normalMsg.ExternalAddress, nonce1RelayerSet.Members[0].ExternalAddress)
	s.Require().EqualValues(types.PowerBase, nonce1RelayerSet.Members[0].Power)

	gravityId := s.Keeper().GetGravityID(s.Ctx)
	checkpoint, err := nonce1RelayerSet.GetCheckpoint(gravityId)
	if trontypes.ModuleName == s.chainName {
		checkpoint, err = trontypes.GetCheckpointRelayerSet(nonce1RelayerSet, gravityId)
	}
	s.Require().NoError(err)

	external1Signature, err := types.NewEthereumSignature(checkpoint, s.externalPris[0])
	if trontypes.ModuleName == s.chainName {
		external1Signature, err = trontypes.NewTronSignature(checkpoint, s.externalPris[0])
	}
	s.Require().NoError(err)
	external2Signature, err := types.NewEthereumSignature(checkpoint, s.externalPris[1])
	if trontypes.ModuleName == s.chainName {
		external2Signature, err = trontypes.NewTronSignature(checkpoint, s.externalPris[1])
	}
	s.Require().NoError(err)

	errMsgData := []struct {
		name      string
		msg       *types.MsgRelayerSetConfirm
		err       error
		errReason string
	}{
		{
			name: "Error relayerSet nonce",
			msg: &types.MsgRelayerSetConfirm{
				Nonce:           0,
				RelayerAddress:  s.relayerAddrs[0].String(),
				ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
				Signature:       hex.EncodeToString(external1Signature),
				ChainName:       s.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("couldn't find relayerSet: %s", types.ErrInvalid),
		},
		{
			name: "not relayer external address",
			msg: &types.MsgRelayerSetConfirm{
				Nonce:           nonce1RelayerSet.Nonce,
				RelayerAddress:  s.relayerAddrs[0].String(),
				ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[1].PublicKey),
				Signature:       hex.EncodeToString(external2Signature),
				ChainName:       s.chainName,
			},
			err:       types.ErrExternalAddressNotMatch,
			errReason: fmt.Sprintf("got %s, expected %s: %s", s.PubKeyToExternalAddr(s.externalPris[1].PublicKey), s.PubKeyToExternalAddr(s.externalPris[0].PublicKey), types.ErrExternalAddressNotMatch),
		},
		{
			name: "sign not match external-1  external-sign-2",
			msg: &types.MsgRelayerSetConfirm{
				Nonce:           nonce1RelayerSet.Nonce,
				RelayerAddress:  s.relayerAddrs[0].String(),
				ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
				Signature:       hex.EncodeToString(external2Signature),
				ChainName:       s.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s and sig %s: %s", normalMsg.ExternalAddress, hex.EncodeToString(checkpoint), hex.EncodeToString(external2Signature), types.ErrInvalid),
		},
		{
			name: "relayer address not match",
			msg: &types.MsgRelayerSetConfirm{
				Nonce:           nonce1RelayerSet.Nonce,
				RelayerAddress:  s.relayerAddrs[1].String(),
				ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
				Signature:       hex.EncodeToString(external1Signature),
				ChainName:       s.chainName,
			},
			err:       types.ErrNotFoundRelayer,
			errReason: fmt.Sprintf("%s", types.ErrNotFoundRelayer),
		},
	}

	for _, testData := range errMsgData {
		_, err = s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), testData.msg)
		s.Require().ErrorIs(err, testData.err, testData.name)
		s.Require().EqualValues(testData.errReason, err.Error(), testData.name)
	}

	normalRelayerSetConfirmMsg := &types.MsgRelayerSetConfirm{
		Nonce:           nonce1RelayerSet.Nonce,
		RelayerAddress:  s.relayerAddrs[0].String(),
		ExternalAddress: normalMsg.ExternalAddress,
		Signature:       hex.EncodeToString(external1Signature),
		ChainName:       s.chainName,
	}
	_, err = s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), normalRelayerSetConfirmMsg)
	s.Require().NoError(err)

	endBlockBeforeLatestRelayerSet := s.Keeper().GetLastRelayerSet(s.Ctx)
	s.Require().NotNil(endBlockBeforeLatestRelayerSet)
}

func (s *KeeperTestSuite) TestClaimWithRelayerOnline() {
	normalMsg := &types.MsgBondedRelayer{
		RelayerAddress:  s.relayerAddrs[0].String(),
		ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
		ChainName:       s.chainName,
	}
	_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), normalMsg)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.Keeper().EndBlocker(s.Ctx)

	latestRelayerSetNonce := s.Keeper().GetLastRelayerSetNonce(s.Ctx)
	s.Require().EqualValues(1, latestRelayerSetNonce)

	nonce1RelayerSet := s.Keeper().GetRelayerSet(s.Ctx, latestRelayerSetNonce)
	s.Require().EqualValues(uint64(1), nonce1RelayerSet.Nonce)
	s.Require().EqualValues(uint64(1), nonce1RelayerSet.Height)

	var gravityId string
	s.Require().NotPanics(func() {
		gravityId = s.Keeper().GetGravityID(s.Ctx)
	})
	s.Require().EqualValues(fmt.Sprintf("me-%s-bridge", s.chainName), gravityId)
	checkpoint, err := nonce1RelayerSet.GetCheckpoint(gravityId)
	if trontypes.ModuleName == s.chainName {
		checkpoint, err = trontypes.GetCheckpointRelayerSet(nonce1RelayerSet, gravityId)
	}
	s.Require().NoError(err)

	relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[0])
	s.Require().True(found)
	relayer.Online = true
	s.Keeper().SetRelayer(s.Ctx, s.relayerAddrs[0], relayer)

	external1Signature, err := types.NewEthereumSignature(checkpoint, s.externalPris[0])
	if trontypes.ModuleName == s.chainName {
		external1Signature, err = trontypes.NewTronSignature(checkpoint, s.externalPris[0])
	}
	s.Require().NoError(err)

	normalRelayerSetConfirmMsg := &types.MsgRelayerSetConfirm{
		Nonce:           latestRelayerSetNonce,
		RelayerAddress:  s.relayerAddrs[0].String(),
		ExternalAddress: normalMsg.ExternalAddress,
		Signature:       hex.EncodeToString(external1Signature),
		ChainName:       s.chainName,
	}
	_, err = s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), normalRelayerSetConfirmMsg)
	s.Require().Nil(err)
}

func (s *KeeperTestSuite) TestClaimMsgGasConsumed() {
	gasStatics := func(gasConsumed, maxGas uint64, minGas uint64, avgGas uint64) (uint64, uint64, uint64) {
		if gasConsumed > maxGas {
			maxGas = gasConsumed
		}
		if minGas == 0 || gasConsumed < minGas {
			minGas = gasConsumed
		}
		if avgGas == 0 {
			avgGas = gasConsumed
		} else {
			avgGas = (avgGas + gasConsumed) / 2
		}
		return maxGas, minGas, avgGas
	}
	testCases := []struct {
		name     string
		buildMsg func() types.ExternalClaim
		execute  func(claim types.ExternalClaim) (minGas, maxGas, avgGas uint64)
	}{
		{
			name: "MsgSendToMe",
			buildMsg: func() types.ExternalClaim {
				return &types.MsgBridgeTokenClaim{
					BlockHeight:   tmrand.Uint64(),
					TokenContract: helpers.GenHexAddress().String(),
					Name:          "Test Token",
					Symbol:        "TEST",
					Decimals:      uint64(tmrand.Int63n(18) + 1),
					ChainName:     s.chainName,
				}
			},
			execute: func(claimMsg types.ExternalClaim) (minGas, maxGas, avgGas uint64) {
				msg, ok := claimMsg.(*types.MsgBridgeTokenClaim)
				s.True(ok)
				for i, relayer := range s.relayerAddrs {
					eventNonce := s.Keeper().GetLastEventNonceByRelayer(s.Ctx, relayer)
					msg.EventNonce = eventNonce + 1
					msg.RelayerAddress = s.relayerAddrs[i].String()
					ctxWithGasMeter := s.Ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
					_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(ctxWithGasMeter), msg)
					s.Require().NoError(err)
					maxGas, minGas, avgGas = gasStatics(ctxWithGasMeter.GasMeter().GasConsumed(), maxGas, minGas, avgGas)
				}
				return
			},
		},
		{
			name: "MsgSendToMeClaim",
			buildMsg: func() types.ExternalClaim {
				return &types.MsgSendToMeClaim{
					BlockHeight:   tmrand.Uint64(),
					TokenContract: helpers.GenHexAddress().String(),
					Amount:        sdk.NewInt(tmrand.Int63n(100000) + 1).MulRaw(1e18),
					Sender:        helpers.GenExternalAddr(s.chainName),
					Receiver:      sdk.AccAddress(tmrand.Bytes(20)).String(),
					ChainName:     s.chainName,
				}
			},
			execute: func(claimMsg types.ExternalClaim) (minGas, maxGas, avgGas uint64) {
				msg, ok := claimMsg.(*types.MsgSendToMeClaim)
				s.True(ok)
				s.Keeper().SetBridgeToken(s.Ctx, &types.BridgeToken{
					ContractAddress: msg.TokenContract,
					Denom:           "test",
					Name:            "Test Token",
					Symbol:          "TEST",
					Decimal:         6,
					Supply:          sdk.NewInt(0),
				})
				for i, relayer := range s.relayerAddrs {
					eventNonce := s.Keeper().GetLastEventNonceByRelayer(s.Ctx, relayer)
					msg.EventNonce = eventNonce + 1
					msg.RelayerAddress = s.relayerAddrs[i].String()
					ctxWithGasMeter := s.Ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
					_, err := s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(ctxWithGasMeter), msg)
					s.Require().NoError(err)
					maxGas, minGas, avgGas = gasStatics(ctxWithGasMeter.GasMeter().GasConsumed(), maxGas, minGas, avgGas)
				}
				return
			},
		},
		{
			name: "RelayerSetUpdateClaim",
			buildMsg: func() types.ExternalClaim {
				var externalRelayerMembers types.BridgeValidators
				for _, key := range s.externalPris {
					bridgeVal := types.BridgeValidator{
						Power:           tmrand.Uint64(),
						ExternalAddress: s.PubKeyToExternalAddr(key.PublicKey),
					}
					externalRelayerMembers = append(externalRelayerMembers, bridgeVal)
				}
				return &types.MsgRelayerSetUpdateClaim{
					BlockHeight:     tmrand.Uint64(),
					RelayerSetNonce: tmrand.Uint64(),
					Members:         externalRelayerMembers,
					ChainName:       s.chainName,
				}
			},
			execute: func(claimMsg types.ExternalClaim) (minGas, maxGas, avgGas uint64) {
				msg, ok := claimMsg.(*types.MsgRelayerSetUpdateClaim)
				s.True(ok)
				s.Keeper().StoreRelayerSet(s.Ctx, &types.RelayerSet{
					Nonce:   msg.RelayerSetNonce,
					Height:  msg.BlockHeight,
					Members: msg.Members,
				})
				for i, relayer := range s.relayerAddrs {
					eventNonce := s.Keeper().GetLastEventNonceByRelayer(s.Ctx, relayer)
					msg.EventNonce = eventNonce + 1
					msg.RelayerAddress = s.relayerAddrs[i].String()
					ctxWithGasMeter := s.Ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
					_, err := s.MsgServer().RelayerSetUpdateClaim(sdk.WrapSDKContext(ctxWithGasMeter), msg)
					s.Require().NoError(err)
					maxGas, minGas, avgGas = gasStatics(ctxWithGasMeter.GasMeter().GasConsumed(), maxGas, minGas, avgGas)
				}
				return
			},
		},
		{
			name: "MsgSendToExternalClaim",
			buildMsg: func() types.ExternalClaim {
				return &types.MsgSendToExternalClaim{
					BlockHeight:   tmrand.Uint64(),
					BatchNonce:    tmrand.Uint64(),
					TokenContract: helpers.GenHexAddress().String(),
					ChainName:     s.chainName,
				}
			},
			execute: func(claimMsg types.ExternalClaim) (minGas, maxGas, avgGas uint64) {
				msg, ok := claimMsg.(*types.MsgSendToExternalClaim)
				s.True(ok)
				s.Require().NoError(s.Keeper().StoreBatch(s.Ctx, &types.OutgoingTxBatch{
					BatchNonce:    msg.BatchNonce,
					TokenContract: msg.TokenContract,
				}))

				for i, relayer := range s.relayerAddrs {
					eventNonce := s.Keeper().GetLastEventNonceByRelayer(s.Ctx, relayer)
					msg.EventNonce = eventNonce + 1
					msg.RelayerAddress = s.relayerAddrs[i].String()
					ctxWithGasMeter := s.Ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
					_, err := s.MsgServer().SendToExternalClaim(sdk.WrapSDKContext(ctxWithGasMeter), msg)
					s.Require().NoError(err)
					maxGas, minGas, avgGas = gasStatics(ctxWithGasMeter.GasMeter().GasConsumed(), maxGas, minGas, avgGas)
				}
				return
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(fmt.Sprintf("%s-%s", s.chainName, testCase.name), func() {
			for i, relayer := range s.relayerAddrs {
				msg := &types.MsgBondedRelayer{
					RelayerAddress:  relayer.String(),
					ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
					DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(2*1e8)),
					ChainName:       s.chainName,
				}
				_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msg)
				s.Require().NoError(err)
			}

			claimMsg := testCase.buildMsg()
			minGas, maxGas, avgGas := testCase.execute(claimMsg)
			s.Require().EqualValuesf(minGas, maxGas, "expect equal min:%d, max:%d, diff:%d", minGas, maxGas, maxGas-minGas)
			s.Require().EqualValuesf(minGas, maxGas, "expect equal min:%d, avg:%d, diff:%d", minGas, avgGas, avgGas-minGas)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBridgeTokenClaim() {
	normalMsg := &types.MsgBondedRelayer{
		RelayerAddress:  s.relayerAddrs[0].String(),
		ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
		ChainName:       s.chainName,
	}
	_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), normalMsg)
	s.Require().NoError(err)

	relayerLastEventNonce := s.Keeper().GetLastEventNonceByRelayer(s.Ctx, s.relayerAddrs[0])
	s.Require().EqualValues(0, relayerLastEventNonce)

	randomPrivateKey, err := crypto.GenerateKey()
	s.Require().NoError(err)
	testMsgs := []struct {
		name      string
		msg       *types.MsgBridgeTokenClaim
		err       error
		errReason string
	}{
		{
			name: "error: event nonce 2",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     2,
				BlockHeight:    1,
				TokenContract:  s.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				RelayerAddress: s.relayerAddrs[0].String(),
				ChainName:      s.chainName,
			},
			err:       types.ErrNonContinuousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 2, 1, types.ErrNonContinuousEventNonce),
		},
		{
			name: "error: event nonce 3",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     3,
				BlockHeight:    1,
				TokenContract:  s.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				RelayerAddress: s.relayerAddrs[0].String(),
				ChainName:      s.chainName,
			},
			err:       types.ErrNonContinuousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 3, 1, types.ErrNonContinuousEventNonce),
		},
		{
			name: "normal: event nonce 1",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     1,
				BlockHeight:    1,
				TokenContract:  s.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				RelayerAddress: s.relayerAddrs[0].String(),
				ChainName:      s.chainName,
			},
			err:       nil,
			errReason: "",
		},
		{
			name: "error again: event nonce 1",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     1,
				BlockHeight:    2,
				TokenContract:  s.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				RelayerAddress: s.relayerAddrs[0].String(),
				ChainName:      s.chainName,
			},
			err:       types.ErrNonContinuousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 1, 2, types.ErrNonContinuousEventNonce),
		},
		{
			name: "error again: event nonce 3",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     3,
				BlockHeight:    2,
				TokenContract:  s.PubKeyToExternalAddr(randomPrivateKey.PublicKey),
				Name:           "Test Token",
				Symbol:         "TEST",
				Decimals:       18,
				RelayerAddress: s.relayerAddrs[0].String(),
				ChainName:      s.chainName,
			},
			err:       types.ErrNonContinuousEventNonce,
			errReason: fmt.Sprintf("got %v, expected %v: %s", 3, 2, types.ErrNonContinuousEventNonce),
		},
	}

	for _, testData := range testMsgs {
		err = testData.msg.ValidateBasic()
		s.Require().NoError(err)
		_, err = s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), testData.msg)
		s.Require().ErrorIs(err, testData.err, testData.name)
		if err == nil {
			continue
		}
		s.Require().EqualValues(testData.errReason, err.Error(), testData.name)
	}
}

func (s *KeeperTestSuite) TestRequestBatchBaseFee() {
	// 1. First sets up a valid relayer set
	totalPower := sdk.ZeroInt()
	delegateAmounts := make([]sdk.Int, 0, len(s.relayerAddrs))
	for i, relayer := range s.relayerAddrs {
		msg := &types.MsgBondedRelayer{
			RelayerAddress:  relayer.String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt((tmrand.Int63n(5)+1)*1e8)),
			ChainName:       s.chainName,
		}
		delegateAmounts = append(delegateAmounts, msg.DelegateAmount.Amount)
		totalPower = totalPower.Add(msg.DelegateAmount.Amount.Quo(sdk.DefaultPowerReduction))
		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msg)
		s.Require().NoError(err)
	}
	s.Keeper().EndBlocker(s.Ctx)

	var externalRelayerMembers types.BridgeValidators
	for i, key := range s.externalPris {
		power := delegateAmounts[i].Quo(sdk.DefaultPowerReduction).MulRaw(int64(types.PowerBase)).Quo(totalPower)
		bridgeVal := types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: s.PubKeyToExternalAddr(key.PublicKey),
		}
		externalRelayerMembers = append(externalRelayerMembers, bridgeVal)
	}
	sort.Sort(externalRelayerMembers)

	// 2. RelayerSetConfirm
	latestRelayerSetNonce := s.Keeper().GetLastRelayerSetNonce(s.Ctx)
	s.Require().EqualValues(1, latestRelayerSetNonce)
	nonce1RelayerSet := s.Keeper().GetRelayerSet(s.Ctx, 1)
	gravityId := s.Keeper().GetGravityID(s.Ctx)
	checkpoint, err := nonce1RelayerSet.GetCheckpoint(gravityId)
	if trontypes.ModuleName == s.chainName {
		checkpoint, err = trontypes.GetCheckpointRelayerSet(nonce1RelayerSet, gravityId)
	}
	for i := range s.relayerAddrs {
		external2Signature, err := types.NewEthereumSignature(checkpoint, s.externalPris[i])
		if trontypes.ModuleName == s.chainName {
			external2Signature, err = trontypes.NewTronSignature(checkpoint, s.externalPris[i])
		}

		msg := &types.MsgRelayerSetConfirm{
			Nonce:           nonce1RelayerSet.Nonce,
			RelayerAddress:  s.relayerAddrs[i].String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			Signature:       hex.EncodeToString(external2Signature),
			ChainName:       s.chainName,
		}
		_, err = s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), msg)
		s.Require().NoError(err)
	}

	// after RelayerSetConfirm, external members should send to ethereum, then we send to me MsgRelayerSetUpdateClaim
	for i := range s.relayerAddrs {
		msg := &types.MsgRelayerSetUpdateClaim{
			EventNonce:      1,
			BlockHeight:     1,
			RelayerSetNonce: 1,
			Members:         externalRelayerMembers,
			RelayerAddress:  s.relayerAddrs[i].String(),
			ChainName:       s.chainName,
		}
		_, err := s.MsgServer().RelayerSetUpdateClaim(sdk.WrapSDKContext(s.Ctx), msg)
		s.Require().NoError(err)
	}
	s.Keeper().EndBlocker(s.Ctx)

	// 3. add bridge token.
	randomPrivateKey, err := crypto.GenerateKey()
	s.Require().NoError(err)
	tokenContract := s.PubKeyToExternalAddr(randomPrivateKey.PublicKey)

	for i, relayer := range s.relayerAddrs {
		normalMsg := &types.MsgBridgeTokenClaim{
			EventNonce:     s.Keeper().GetLastEventNonceByRelayer(s.Ctx, relayer) + 1,
			BlockHeight:    1,
			TokenContract:  tokenContract,
			Name:           "Tether USD",
			Symbol:         "USDT",
			Decimals:       18,
			RelayerAddress: s.relayerAddrs[i].String(),
			ChainName:      s.chainName,
		}
		_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), normalMsg)
		s.Require().NoError(err)
	}
	s.Keeper().EndBlocker(s.Ctx)

	bridgeDenomData, err := s.Keeper().GetBridgeTokenByDenom(s.Ctx, "uusdt")
	s.Require().NoError(err)
	s.Require().NotNil(bridgeDenomData)
	s.Require().EqualValues(tokenContract, bridgeDenomData.ContractAddress)

	// 4. sendToMe.
	sendToMeSendAddr := s.PubKeyToExternalAddr(s.externalPris[0].PublicKey)
	sendToMeReceiveAddr := s.relayerAddrs[0]
	sendToMeClaim := new(types.MsgSendToMeClaim)
	sendToMeAmount := sdkmath.NewIntWithDecimal(1000, 18)
	for i, relayer := range s.relayerAddrs {
		sendToMeClaim = &types.MsgSendToMeClaim{
			EventNonce:     s.Keeper().GetLastEventNonceByRelayer(s.Ctx, relayer) + 1,
			BlockHeight:    1,
			TokenContract:  tokenContract,
			Amount:         sendToMeAmount,
			Sender:         sendToMeSendAddr,
			Receiver:       sendToMeReceiveAddr.String(),
			RelayerAddress: s.relayerAddrs[i].String(),
			ChainName:      s.chainName,
		}
		_, err := s.MsgServer().SendToMeClaim(sdk.WrapSDKContext(s.Ctx), sendToMeClaim)
		s.Require().NoError(err)
	}

	balance := s.App.BankKeeper.GetBalance(s.Ctx, sendToMeReceiveAddr, bridgeDenomData.Denom)
	s.Require().NotNil(balance)
	s.Require().EqualValues(balance.Denom, bridgeDenomData.Denom)
	receiveAmount := types.GetMintCoin(sendToMeClaim.Amount, sendToMeClaim.ChainName, bridgeDenomData)
	s.Require().True(balance.Amount.Equal(receiveAmount.Amount))

	sendToExternal := func(bridgeFees []sdk.Int) {
		for _, bridgeFee := range bridgeFees {
			msg := &types.MsgSendToExternal{
				Sender:    sendToMeReceiveAddr.String(),
				Dest:      sendToMeSendAddr,
				Amount:    sdk.NewCoin(bridgeDenomData.Denom, sdk.NewInt(3)),
				BridgeFee: sdk.NewCoin(bridgeDenomData.Denom, bridgeFee),
				ChainName: s.chainName,
			}
			_, err := s.MsgServer().SendToExternal(sdk.WrapSDKContext(s.Ctx), msg)
			s.Require().NoError(err)
		}
	}

	sendToExternal([]sdk.Int{sdk.NewInt(1), sdk.NewInt(2), sdk.NewInt(3)})
	usdtBatchFee := s.Keeper().GetBatchFeesByTokenType(s.Ctx, tokenContract, 100, sdk.NewInt(0))
	s.Require().EqualValues(tokenContract, usdtBatchFee.TokenContract)
	s.Require().EqualValues(3, usdtBatchFee.TotalTxs)
	s.Require().EqualValues(sdk.NewInt(6000000000000), usdtBatchFee.TotalFees)

	testCases := []struct {
		testName       string
		baseFee        sdk.Int
		pass           bool
		expectTotalTxs uint64
		err            error
	}{
		{
			testName:       "Support - baseFee 1000",
			baseFee:        sdk.NewInt(1000).MulRaw(1e12),
			pass:           false,
			expectTotalTxs: 3,
			err:            errorsmod.Wrap(types.ErrEmpty, "no batch tx selected"),
		},
		{
			testName:       "Support - baseFee 2",
			baseFee:        sdk.NewInt(2).MulRaw(1e12),
			pass:           true,
			expectTotalTxs: 1,
			err:            nil,
		},
		{
			testName:       "Support - baseFee 0",
			baseFee:        sdk.NewInt(0),
			pass:           true,
			expectTotalTxs: 0,
			err:            nil,
		},
	}

	for _, testCase := range testCases {
		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
		_, err := s.MsgServer().RequestBatch(sdk.WrapSDKContext(s.Ctx), &types.MsgRequestBatch{
			Sender:     s.relayerAddrs[0].String(),
			Denom:      bridgeDenomData.Denom,
			MinimumFee: sdk.NewInt(1).MulRaw(1e12),
			FeeReceive: "0x0000000000000000000000000000000000000000",
			ChainName:  s.chainName,
			BaseFee:    testCase.baseFee,
		})
		if testCase.pass {
			s.Require().NoError(err)
			usdtBatchFee = s.Keeper().GetBatchFeesByTokenType(s.Ctx, tokenContract, 100, sdk.NewInt(0))
			s.Require().EqualValues(testCase.expectTotalTxs, usdtBatchFee.TotalTxs)
		} else {
			s.Require().NotNil(err)
			s.Require().Equal(err.Error(), testCase.err.Error())
		}
	}
}

func (s *KeeperTestSuite) addBridgeToken(tokenContract string, md banktypes.Metadata) {
	relayerLastEventNonce := s.Keeper().GetLastEventNonceByRelayer(s.Ctx, s.relayerAddrs[0])
	ctx := sdk.WrapSDKContext(s.Ctx.WithEventManager(sdk.NewEventManager()))
	_, err := s.MsgServer().BridgeTokenClaim(ctx, &types.MsgBridgeTokenClaim{
		EventNonce:     relayerLastEventNonce + 1,
		BlockHeight:    uint64(s.Ctx.BlockHeight()),
		TokenContract:  tokenContract,
		Name:           md.Name,
		Symbol:         md.Symbol,
		Decimals:       18,
		RelayerAddress: s.relayerAddrs[0].String(),
		ChainName:      s.chainName,
	})
	s.Require().NoError(err)

	s.checkObservationState(ctx, true)

	newRelayerLastEventNonce := s.Keeper().GetLastEventNonceByRelayer(s.Ctx, s.relayerAddrs[0])
	s.Require().EqualValues(relayerLastEventNonce+1, newRelayerLastEventNonce)
}

func (s *KeeperTestSuite) checkObservationState(ctx context.Context, expect bool) {
	foundObservation := false
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, event := range sdkCtx.EventManager().Events() {
		if event.Type != types.EventTypeContractEvent {
			continue
		}
		s.Require().False(foundObservation, "found multiple observation event")
		for _, attr := range event.Attributes {
			if attr.Key != types.AttributeKeyStateSuccess {
				continue
			}
			s.Require().EqualValues(fmt.Sprintf("%v", expect), attr.Value)
			foundObservation = true
			break
		}
	}
	s.Require().True(foundObservation, "not found observation event")
	sdkCtx.WithEventManager(sdk.NewEventManager())
}
