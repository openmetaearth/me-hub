package keeper_test

import (
	"fmt"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/testutil/helpers"
	"strings"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestEndBlockDepositClaim() {
	normalMsg := &types.MsgBondedRelayer{
		RelayerAddress:  s.relayerAddrs[0].String(),
		ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[0].PublicKey),
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
		ChainName:       s.chainName,
	}

	_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), normalMsg)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)

	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

	bridgeTokenContract := helpers.GenExternalAddr(s.chainName)
	sendToMeSendAddr := helpers.GenExternalAddr(s.chainName)
	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeTokenContract,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       6,
		RelayerAddress: s.relayerAddrs[0].String(),
		ChainName:      s.chainName,
	}
	_, err = s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), addBridgeTokenClaim)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

	sendToMeClaim := &types.MsgSendToMeClaim{
		EventNonce:     2,
		BlockHeight:    1001,
		TokenContract:  bridgeTokenContract,
		Amount:         sdkmath.NewInt(1234),
		Sender:         sendToMeSendAddr,
		Receiver:       helpers.GenAccAddress().String(),
		RelayerAddress: s.relayerAddrs[0].String(),
		ChainName:      s.chainName,
	}
	s.SendClaim(sendToMeClaim)

	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

	allBalances := s.App.BankKeeper.GetAllBalances(s.Ctx, sdk.MustAccAddressFromBech32(sendToMeClaim.Receiver))
	s.Require().EqualValues(sdk.Coin{Amount: sendToMeClaim.Amount, Denom: strings.ToLower(addBridgeTokenClaim.GetSymbol())}.String(), allBalances.String())

	bridgeToken, err := s.Keeper().GetBridgeTokenByContract(s.Ctx, bridgeTokenContract)
	s.Require().NoError(err)
	s.Require().EqualValues(sendToMeClaim.Amount, bridgeToken.Supply)
}

func (s *KeeperTestSuite) TestRelayerUpdate() {
	if len(s.relayerAddrs) < 5 {
		return
	}
	for i := 0; i < 5; i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:  s.relayerAddrs[i].String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
			ChainName:       s.chainName,
		}
		s.Require().NoError(msgBondedRelayer.ValidateBasic())
		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)

		s.Require().NoError(err)
		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
		relayerSets := s.Keeper().GetRelayerSets(s.Ctx)
		s.Require().NotNil(relayerSets)
		s.Require().EqualValues(i+1, len(relayerSets))

		power := s.Keeper().GetLastTotalPower(s.Ctx)
		expectPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		s.Require().True(expectPower.Equal(power))
	}

	bridgeToken := helpers.GenExternalAddr(s.chainName)

	for i := 0; i < 6; i++ {
		addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  bridgeToken,
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			RelayerAddress: s.relayerAddrs[i].String(),
			ChainName:      s.chainName,
		}
		_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), addBridgeTokenClaim)
		s.Require().NoError(err)
		endBlockBeforeAttestation := s.Keeper().GetAttestation(s.Ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		s.Require().NotNil(endBlockBeforeAttestation)
		s.Require().False(endBlockBeforeAttestation.Observed)
		s.Require().NotNil(endBlockBeforeAttestation.Votes)
		s.Require().EqualValues(i+1, len(endBlockBeforeAttestation.Votes))

		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
		endBlockAfterAttestation := s.Keeper().GetAttestation(s.Ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		s.Require().NotNil(endBlockAfterAttestation)
		s.Require().False(endBlockAfterAttestation.Observed)
	}

	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       18,
		RelayerAddress: s.relayerAddrs[6].String(),
		ChainName:      s.chainName,
	}
	_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), addBridgeTokenClaim)
	s.Require().NoError(err)
	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	attestation := s.Keeper().GetAttestation(s.Ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())

	s.Require().NotNil(attestation)
	s.Require().True(attestation.Observed)

	var newRelayerList []string
	for i := 0; i < 7; i++ {
		newRelayerList = append(newRelayerList, s.relayerAddrs[i].String())
	}
	_, err = s.MsgServer().ProposalRelayers(sdk.WrapSDKContext(s.Ctx), &types.MsgProposalRelayers{
		ChainName: s.chainName,
		Authority: s.Dao.GlobalDao,
		Relayers:  newRelayerList,
	})
	s.Require().NoError(err)
	s.Require().ErrorIs(types.ErrInvalid, err)

	expectTotalPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := s.Keeper().GetLastTotalPower(s.Ctx)
	s.Require().True(expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalRelayerChangePowerThreshold.Mul(expectTotalPower).Quo(sdkmath.NewInt(100))
	expectDeletePower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(3)).Quo(sdk.DefaultPowerReduction)
	s.Require().EqualValues(fmt.Sprintf("max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s: %s", expectMaxChangePower.String(), expectDeletePower.String(), types.ErrInvalid), err.Error())

	var newRelayerList2 []string
	for i := 0; i < 8; i++ {
		newRelayerList2 = append(newRelayerList2, s.relayerAddrs[i].String())
	}
	_, err = s.MsgServer().ProposalRelayers(sdk.WrapSDKContext(s.Ctx), &types.MsgProposalRelayers{
		ChainName: s.chainName,
		Authority: s.Dao.GlobalDao,
		Relayers:  newRelayerList,
	})
	s.Require().NoError(err)
}

//func (s *KeeperTestSuite) TestAttestationAfterRelayerUpdate() {
//	if len(s.bridgerAddrs) < 20 {
//		return
//	}
//	for i := 0; i < 20; i++ {
//		msgBondedRelayer := &types.MsgBondedRelayer{
//			RelayerAddress:   s.relayerAddrs[i].String(),
//			BridgerAddress:   s.bridgerAddrs[i].String(),
//			ExternalAddress:  s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
//			ValidatorAddress: s.valAddrs[i].String(),
//			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
//			ChainName:        s.chainName,
//		}
//		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)
//		s.Require().NoError(err)
//		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//		relayerSets := s.Keeper().GetRelayerSets(s.Ctx)
//		s.Require().NotNil(relayerSets)
//		s.Require().EqualValues(i+1, len(relayerSets))
//
//		power := s.Keeper().GetLastTotalPower(s.Ctx)
//		expectPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
//		s.Require().True(expectPower.Equal(power))
//	}
//
//	{
//		firstBridgeTokenClaim := &types.MsgBridgeTokenClaim{
//			EventNonce:     1,
//			BlockHeight:    1000,
//			TokenContract:  helpers.GenExternalAddr(s.chainName),
//			Name:           "Test Token",
//			Symbol:         "TEST",
//			Decimals:       18,
//			BridgerAddress: "",
//			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
//			ChainName:      s.chainName,
//		}
//
//		for i := 0; i < 13; i++ {
//			firstBridgeTokenClaim.BridgerAddress = s.bridgerAddrs[i].String()
//			_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), firstBridgeTokenClaim)
//			s.Require().NoError(err)
//			endBlockBeforeAttestation := s.Keeper().GetAttestation(s.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
//			s.Require().NotNil(endBlockBeforeAttestation)
//			s.Require().False(endBlockBeforeAttestation.Observed)
//			s.Require().NotNil(endBlockBeforeAttestation.Votes)
//			s.Require().EqualValues(i+1, len(endBlockBeforeAttestation.Votes))
//
//			endBlockAfterAttestation := s.Keeper().GetAttestation(s.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
//			s.Require().NotNil(endBlockAfterAttestation)
//			s.Require().False(endBlockAfterAttestation.Observed)
//		}
//
//		firstBridgeTokenClaim.BridgerAddress = s.bridgerAddrs[13].String()
//		_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), firstBridgeTokenClaim)
//		s.Require().NoError(err)
//		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//		attestation := s.Keeper().GetAttestation(s.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
//
//		s.Require().NotNil(attestation)
//		s.Require().True(attestation.Observed)
//	}
//
//	{
//		secondBridgeTokenClaim := &types.MsgBridgeTokenClaim{
//			EventNonce:     2,
//			BlockHeight:    1001,
//			TokenContract:  helpers.GenExternalAddr(s.chainName),
//			Name:           "Test Token2",
//			Symbol:         "TEST2",
//			Decimals:       18,
//			BridgerAddress: "",
//			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
//			ChainName:      s.chainName,
//		}
//
//		for i := 0; i < 6; i++ {
//			secondBridgeTokenClaim.BridgerAddress = s.bridgerAddrs[i].String()
//			_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), secondBridgeTokenClaim)
//			s.Require().NoError(err)
//			endBlockBeforeAttestation := s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
//			s.Require().NotNil(endBlockBeforeAttestation)
//			s.Require().False(endBlockBeforeAttestation.Observed)
//			s.Require().NotNil(endBlockBeforeAttestation.Votes)
//			s.Require().EqualValues(i+1, len(endBlockBeforeAttestation.Votes))
//
//			s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//			s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//			endBlockAfterAttestation := s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
//			s.Require().NotNil(endBlockAfterAttestation)
//			s.Require().False(endBlockAfterAttestation.Observed)
//		}
//
//		secondClaimAttestation := s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
//		s.Require().NotNil(secondClaimAttestation)
//		s.Require().False(secondClaimAttestation.Observed)
//		s.Require().NotNil(secondClaimAttestation.Votes)
//		s.Require().EqualValues(6, len(secondClaimAttestation.Votes))
//
//		var newRelayerList []string
//		for i := 0; i < 15; i++ {
//			newRelayerList = append(newRelayerList, s.relayerAddrs[i].String())
//		}
//		_, err := s.MsgServer().UpdateChainRelayers(s.Ctx, &types.MsgUpdateChainRelayers{
//			Relayers:  newRelayerList,
//			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
//			ChainName: s.chainName,
//		})
//		s.Require().NoError(err)
//		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//
//		secondClaimAttestation = s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
//		s.Require().NotNil(secondClaimAttestation)
//		s.Require().False(secondClaimAttestation.Observed)
//		s.Require().NotNil(secondClaimAttestation.Votes)
//		s.Require().EqualValues(6, len(secondClaimAttestation.Votes))
//
//		activeRelayers := s.Keeper().GetAllRelayers(s.Ctx, true)
//		s.Require().NotNil(activeRelayers)
//		s.Require().EqualValues(15, len(activeRelayers))
//		for i := 0; i < 15; i++ {
//			s.Require().NotNil(newRelayerList[i], activeRelayers[i].RelayerAddress)
//		}
//
//		var newRelayerList2 []string
//		for i := 0; i < 11; i++ {
//			newRelayerList2 = append(newRelayerList2, s.relayerAddrs[i].String())
//		}
//		_, err = s.MsgServer().UpdateChainRelayers(s.Ctx, &types.MsgUpdateChainRelayers{
//			Relayers:  newRelayerList2,
//			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
//			ChainName: s.chainName,
//		})
//		s.Require().NoError(err)
//		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//
//		secondClaimAttestation = s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
//		s.Require().NotNil(secondClaimAttestation)
//		s.Require().False(secondClaimAttestation.Observed)
//		s.Require().NotNil(secondClaimAttestation.Votes)
//		s.Require().EqualValues(6, len(secondClaimAttestation.Votes))
//
//		activeRelayers = s.Keeper().GetAllRelayers(s.Ctx, true)
//		s.Require().NotNil(activeRelayers)
//		s.Require().EqualValues(11, len(activeRelayers))
//		for i := 0; i < 11; i++ {
//			s.Require().NotNil(newRelayerList2[i], activeRelayers[i].RelayerAddress)
//		}
//
//		var newRelayerList3 []string
//		for i := 0; i < 10; i++ {
//			newRelayerList3 = append(newRelayerList3, s.relayerAddrs[i].String())
//		}
//		_, err = s.MsgServer().UpdateChainRelayers(s.Ctx, &types.MsgUpdateChainRelayers{
//			Relayers:  newRelayerList3,
//			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
//			ChainName: s.chainName,
//		})
//		s.Require().NoError(err)
//		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//
//		secondClaimAttestation = s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
//		s.Require().NotNil(secondClaimAttestation)
//		s.Require().False(secondClaimAttestation.Observed)
//		s.Require().NotNil(secondClaimAttestation.Votes)
//		s.Require().EqualValues(6, len(secondClaimAttestation.Votes))
//
//		activeRelayers = s.Keeper().GetAllRelayers(s.Ctx, true)
//		s.Require().NotNil(activeRelayers)
//		s.Require().EqualValues(10, len(activeRelayers))
//		for i := 0; i < 10; i++ {
//			s.Require().NotNil(newRelayerList3[i], activeRelayers[i].RelayerAddress)
//		}
//
//		secondBridgeTokenClaim.BridgerAddress = s.bridgerAddrs[6].String()
//		_, err = s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), secondBridgeTokenClaim)
//		s.Require().NoError(err)
//
//		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//
//		secondClaimAttestation = s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
//		s.Require().NotNil(secondClaimAttestation)
//		s.Require().True(secondClaimAttestation.Observed)
//		s.Require().NotNil(secondClaimAttestation.Votes)
//		s.Require().EqualValues(7, len(secondClaimAttestation.Votes))
//	}
//}
//
//func (s *KeeperTestSuite) TestRelayerDelete() {
//	for i := 0; i < len(s.relayerAddrs); i++ {
//		msgBondedRelayer := &types.MsgBondedRelayer{
//			RelayerAddress:   s.relayerAddrs[i].String(),
//			BridgerAddress:   s.bridgerAddrs[i].String(),
//			ExternalAddress:  s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
//			ValidatorAddress: s.valAddrs[i].String(),
//			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
//			ChainName:        s.chainName,
//		}
//		s.Require().NoError(msgBondedRelayer.ValidateBasic())
//		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)
//		s.Require().NoError(err)
//	}
//	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//	allRelayers := s.Keeper().GetAllRelayers(s.Ctx, false)
//	s.Require().NotNil(allRelayers)
//	s.Require().EqualValues(len(s.relayerAddrs), len(allRelayers))
//
//	relayer := s.relayerAddrs[0]
//	bridger := s.bridgerAddrs[0]
//	externalAddress := s.PubKeyToExternalAddr(s.externalPris[0].PublicKey)
//
//	relayerAddr, found := s.Keeper().GetRelayerAddrByBridgerAddr(s.Ctx, bridger)
//	s.Require().True(found)
//	s.Require().EqualValues(relayer.String(), relayerAddr.String())
//
//	relayerAddr, found = s.Keeper().GetRelayerAddrByExternalAddr(s.Ctx, externalAddress)
//	s.Require().True(found)
//	s.Require().EqualValues(relayer.String(), relayerAddr.String())
//
//	relayerData, found := s.Keeper().GetRelayer(s.Ctx, relayer)
//	s.Require().True(found)
//	s.Require().NotNil(relayerData)
//	s.Require().EqualValues(relayer.String(), relayerData.RelayerAddress)
//	s.Require().EqualValues(bridger.String(), relayerData.BridgerAddress)
//	s.Require().EqualValues(externalAddress, relayerData.ExternalAddress)
//
//	s.Require().True(sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Equal(relayerData.DelegateAmount))
//
//	newRelayerAddressList := make([]string, 0, len(s.relayerAddrs)-1)
//	for _, address := range s.relayerAddrs[1:] {
//		newRelayerAddressList = append(newRelayerAddressList, address.String())
//	}
//
//	_, err := s.MsgServer().UpdateChainRelayers(s.Ctx, &types.MsgUpdateChainRelayers{
//		Relayers:  newRelayerAddressList,
//		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
//		ChainName: s.chainName,
//	})
//	s.Require().NoError(err)
//	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//
//	relayerAddr, found = s.Keeper().GetRelayerAddrByBridgerAddr(s.Ctx, bridger)
//	s.Require().True(found)
//	s.Require().Equal(relayerAddr, relayer)
//
//	relayerAddr, found = s.Keeper().GetRelayerAddrByExternalAddr(s.Ctx, externalAddress)
//	s.Require().True(found)
//	s.Require().Equal(relayerAddr, relayer)
//
//	relayerData, found = s.Keeper().GetRelayer(s.Ctx, relayer)
//	s.Require().True(found)
//}
//
//func (s *KeeperTestSuite) TestRelayerSetSlash() {
//	for i := 0; i < len(s.relayerAddrs); i++ {
//		msgBondedRelayer := &types.MsgBondedRelayer{
//			RelayerAddress:   s.relayerAddrs[i].String(),
//			BridgerAddress:   s.bridgerAddrs[i].String(),
//			ExternalAddress:  s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
//			ValidatorAddress: s.valAddrs[i].String(),
//			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
//			ChainName:        s.chainName,
//		}
//		s.Require().NoError(msgBondedRelayer.ValidateBasic())
//		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)
//		s.Require().NoError(err)
//	}
//	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//	s.Keeper().EndBlocker(s.Ctx)
//
//	allRelayers := s.Keeper().GetAllRelayers(s.Ctx, false)
//	s.Require().NotNil(allRelayers)
//	s.Require().Equal(len(s.relayerAddrs), len(allRelayers))
//
//	relayerSets := s.Keeper().GetRelayerSets(s.Ctx)
//	s.Require().NotNil(relayerSets)
//	s.Require().EqualValues(1, len(relayerSets))
//
//	for i := 0; i < len(s.relayerAddrs)-1; i++ {
//		externalAddress, signature := s.SignRelayerSetConfirm(s.externalPris[i], relayerSets[0])
//		relayerSetConfirm := &types.MsgRelayerSetConfirm{
//			Nonce:           relayerSets[0].Nonce,
//			BridgerAddress:  s.bridgerAddrs[i].String(),
//			ExternalAddress: externalAddress,
//			Signature:       hex.EncodeToString(signature),
//			ChainName:       s.chainName,
//		}
//		s.Require().NoError(relayerSetConfirm.ValidateBasic())
//		_, err := s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), relayerSetConfirm)
//		s.Require().NoError(err)
//	}
//
//	s.Keeper().EndBlocker(s.Ctx)
//	relayerSetHeight := int64(relayerSets[0].Height)
//	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
//	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
//
//	relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[len(s.relayerAddrs)-1])
//	s.Require().True(found)
//	s.Require().True(relayer.Online)
//	s.Require().Equal(int64(0), relayer.SlashTimes)
//
//	s.Ctx = s.Ctx.WithBlockHeight(relayerSetHeight + int64(s.Keeper().GetParams(s.Ctx).SignedWindow) + 1)
//	s.Keeper().EndBlocker(s.Ctx)
//
//	relayer, found = s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[len(s.relayerAddrs)-1])
//	s.Require().True(found)
//	s.Require().False(relayer.Online)
//	s.Require().Equal(int64(1), relayer.SlashTimes)
//}
//
//func (s *KeeperTestSuite) TestSlashRelayer() {
//	for i := 0; i < len(s.relayerAddrs); i++ {
//		msgBondedRelayer := &types.MsgBondedRelayer{
//			RelayerAddress:   s.relayerAddrs[i].String(),
//			BridgerAddress:   s.bridgerAddrs[i].String(),
//			ExternalAddress:  s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
//			ValidatorAddress: s.valAddrs[i].String(),
//			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
//			ChainName:        s.chainName,
//		}
//		s.Require().NoError(msgBondedRelayer.ValidateBasic())
//		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)
//		s.Require().NoError(err)
//	}
//
//	params := s.Keeper().GetParams(s.Ctx)
//	err := s.Keeper().SetParams(s.Ctx, &params)
//	s.Require().NoError(err)
//	for i := 0; i < len(s.relayerAddrs); i++ {
//		relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[i])
//		s.Require().True(found)
//		s.Require().True(relayer.Online)
//		s.Require().Equal(int64(0), relayer.SlashTimes)
//
//		s.Keeper().SlashRelayer(s.Ctx, relayer.RelayerAddress)
//
//		relayer, found = s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[i])
//		s.Require().True(found)
//		s.Require().False(relayer.Online)
//		s.Require().Equal(int64(1), relayer.SlashTimes)
//	}
//
//	// repeat slash test.
//	for i := 0; i < len(s.relayerAddrs); i++ {
//		relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[i])
//		s.Require().True(found)
//		s.Require().False(relayer.Online)
//		s.Require().Equal(int64(1), relayer.SlashTimes)
//
//		s.Keeper().SlashRelayer(s.Ctx, relayer.RelayerAddress)
//
//		relayer, found = s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[i])
//		s.Require().True(found)
//		s.Require().False(relayer.Online)
//		s.Require().Equal(int64(1), relayer.SlashTimes)
//	}
//}
