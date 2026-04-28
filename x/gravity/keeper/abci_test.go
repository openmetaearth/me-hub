package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"encoding/hex"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestDepositClaim() {
	proposalRelayers, found := s.Keeper().GetProposalRelayer(s.Ctx)
	s.Require().True(found)
	s.Require().EqualValues(s.relayerNumber, len(proposalRelayers.Relayers))

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
	s.Keeper().EndBlocker(s.Ctx)

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
	s.Require().EqualValues(sdk.Coin{Amount: sendToMeClaim.Amount, Denom: utils.GetDenom(addBridgeTokenClaim.GetSymbol())}.String(), allBalances.String())

	bridgeToken, err := s.Keeper().GetBridgeTokenByContract(s.Ctx, bridgeTokenContract)
	s.Require().NoError(err)
	s.Require().EqualValues(sendToMeClaim.Amount, bridgeToken.Supply)
}

func (s *KeeperTestSuite) TestProposalRelayers() {
	proposalRelayers, found := s.Keeper().GetProposalRelayer(s.Ctx)
	s.Require().True(found)
	s.Require().EqualValues(s.relayerNumber, len(proposalRelayers.Relayers))

	// init 10 relayers
	for i := 0; i < s.relayerNumber; i++ {
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
		expectPower := sdkmath.NewInt(10 * 1e8).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		s.Require().True(expectPower.Equal(power))
	}

	bridgeToken := helpers.GenExternalAddr(s.chainName)

	// 6/10 < 2/3
	for i := 0; i < 6; i++ {
		addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  bridgeToken,
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       6,
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

	// last relayer finish attestation, 7/10 > 2/3
	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       6,
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
	for i := 0; i < 6; i++ {
		newRelayerList = append(newRelayerList, s.relayerAddrs[i].String())
	}
	_, err = s.MsgServer().ProposalRelayers(sdk.WrapSDKContext(s.Ctx), &types.MsgProposalRelayers{
		ChainName: s.chainName,
		Authority: s.Dao.GlobalDao,
		Relayers:  newRelayerList,
	})
	s.Require().ErrorIs(types.ErrMaxChangePowerLimitExceeded, err) // try update chain relayer power >= 30%, expect error

	expectTotalPower := sdkmath.NewInt(10 * 1e8).Mul(sdkmath.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := s.Keeper().GetLastTotalPower(s.Ctx)
	s.Require().True(expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalRelayerChangePowerThreshold.Mul(expectTotalPower).Quo(sdkmath.NewInt(int64(types.PowerBase)))
	expectDeletePower := sdkmath.NewInt(10 * 1e8).Mul(sdkmath.NewInt(4)).Quo(sdk.DefaultPowerReduction)
	s.Require().EqualValues(fmt.Sprintf("maxChangePowerThreshold: %s, deleteTotalPower: %s: %s",
		expectMaxChangePower.String(), expectDeletePower.String(), types.ErrMaxChangePowerLimitExceeded), err.Error())

	var newRelayerList2 []string
	for i := 0; i < 7; i++ {
		newRelayerList2 = append(newRelayerList2, s.relayerAddrs[i].String())
	}
	_, err = s.MsgServer().ProposalRelayers(sdk.WrapSDKContext(s.Ctx), &types.MsgProposalRelayers{
		ChainName: s.chainName,
		Authority: s.Dao.GlobalDao,
		Relayers:  newRelayerList2,
	})
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestAttestationAfterRelayerUpdate() {
	proposalRelayers, found := s.Keeper().GetProposalRelayer(s.Ctx)
	s.Require().True(found)
	s.Require().EqualValues(s.relayerNumber, len(proposalRelayers.Relayers))

	for i := 0; i < s.relayerNumber; i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:  s.relayerAddrs[i].String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
			ChainName:       s.chainName,
		}
		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)
		s.Require().NoError(err)
		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
		relayerSets := s.Keeper().GetRelayerSets(s.Ctx)
		s.Require().NotNil(relayerSets)
		s.Require().EqualValues(i+1, len(relayerSets))

		power := s.Keeper().GetLastTotalPower(s.Ctx)
		expectPower := sdkmath.NewInt(10 * 1e8).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		s.Require().True(expectPower.Equal(power))
	}

	// case1: normal, 7/10 > 2/3
	{
		firstBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  helpers.GenExternalAddr(s.chainName),
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			RelayerAddress: "",
			ChainName:      s.chainName,
		}

		for i := 0; i < 6; i++ {
			firstBridgeTokenClaim.RelayerAddress = s.relayerAddrs[i].String()
			_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), firstBridgeTokenClaim)
			s.Require().NoError(err)
			endBlockBeforeAttestation := s.Keeper().GetAttestation(s.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			s.Require().NotNil(endBlockBeforeAttestation)
			s.Require().False(endBlockBeforeAttestation.Observed)
			s.Require().NotNil(endBlockBeforeAttestation.Votes)
			s.Require().EqualValues(i+1, len(endBlockBeforeAttestation.Votes))

			endBlockAfterAttestation := s.Keeper().GetAttestation(s.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			s.Require().NotNil(endBlockAfterAttestation)
			s.Require().False(endBlockAfterAttestation.Observed)
		}

		firstBridgeTokenClaim.RelayerAddress = s.relayerAddrs[6].String()
		_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), firstBridgeTokenClaim)
		s.Require().NoError(err)
		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)

		attestation := s.Keeper().GetAttestation(s.Ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
		s.Require().NotNil(attestation)
		s.Require().True(attestation.Observed)
	}

	// case2: 5/10 < 2/3
	approveNumber := 2
	{
		secondBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     2,
			BlockHeight:    1001,
			TokenContract:  helpers.GenExternalAddr(s.chainName),
			Name:           "Test Token2",
			Symbol:         "TEST2",
			Decimals:       18,
			RelayerAddress: "",
			ChainName:      s.chainName,
		}

		for i := 0; i < approveNumber; i++ {
			secondBridgeTokenClaim.RelayerAddress = s.relayerAddrs[i].String()
			_, err := s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), secondBridgeTokenClaim)
			s.Require().NoError(err)
			endBlockBeforeAttestation := s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			s.Require().NotNil(endBlockBeforeAttestation)
			s.Require().False(endBlockBeforeAttestation.Observed)
			s.Require().NotNil(endBlockBeforeAttestation.Votes)
			s.Require().EqualValues(i+1, len(endBlockBeforeAttestation.Votes))

			s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
			s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
			endBlockAfterAttestation := s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			s.Require().NotNil(endBlockAfterAttestation)
			s.Require().False(endBlockAfterAttestation.Observed)
		}

		secondClaimAttestation := s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		s.Require().NotNil(secondClaimAttestation)
		s.Require().False(secondClaimAttestation.Observed)
		s.Require().NotNil(secondClaimAttestation.Votes)
		s.Require().EqualValues(approveNumber, len(secondClaimAttestation.Votes))

		// from 10 change to 7
		var newRelayerList []string
		for i := 0; i < 7; i++ {
			newRelayerList = append(newRelayerList, s.relayerAddrs[i].String())
		}
		_, err := s.MsgServer().ProposalRelayers(s.Ctx, &types.MsgProposalRelayers{
			Relayers:  newRelayerList,
			Authority: s.Dao.GlobalDao,
			ChainName: s.chainName,
		})
		s.Require().NoError(err)
		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

		secondClaimAttestation = s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		s.Require().NotNil(secondClaimAttestation)
		s.Require().False(secondClaimAttestation.Observed)
		s.Require().NotNil(secondClaimAttestation.Votes)
		s.Require().EqualValues(approveNumber, len(secondClaimAttestation.Votes))

		activeRelayers := s.Keeper().GetAllRelayers(s.Ctx, true)
		s.Require().NotNil(activeRelayers)
		s.Require().EqualValues(7, len(activeRelayers))
		for i := 0; i < 7; i++ {
			s.Require().NotNil(newRelayerList[i], activeRelayers[i].RelayerAddress)
		}

		// from 7 change to 5
		var newRelayerList2 []string
		for i := 0; i < 5; i++ {
			newRelayerList2 = append(newRelayerList2, s.relayerAddrs[i].String())
		}
		_, err = s.MsgServer().ProposalRelayers(s.Ctx, &types.MsgProposalRelayers{
			Relayers:  newRelayerList2,
			Authority: s.Dao.GlobalDao,
			ChainName: s.chainName,
		})
		s.Require().NoError(err)
		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

		secondClaimAttestation = s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		s.Require().NotNil(secondClaimAttestation)
		s.Require().False(secondClaimAttestation.Observed)
		s.Require().NotNil(secondClaimAttestation.Votes)
		s.Require().EqualValues(approveNumber, len(secondClaimAttestation.Votes))

		activeRelayers = s.Keeper().GetAllRelayers(s.Ctx, true)
		s.Require().NotNil(activeRelayers)
		s.Require().EqualValues(5, len(activeRelayers))
		for i := 0; i < 5; i++ {
			s.Require().NotNil(newRelayerList2[i], activeRelayers[i].RelayerAddress)
		}

		// change from 5 to 4
		var newRelayerList3 []string
		for i := 0; i < 4; i++ {
			newRelayerList3 = append(newRelayerList3, s.relayerAddrs[i].String())
		}
		_, err = s.MsgServer().ProposalRelayers(s.Ctx, &types.MsgProposalRelayers{
			Relayers:  newRelayerList3,
			Authority: s.Dao.GlobalDao,
			ChainName: s.chainName,
		})
		s.Require().NoError(err)
		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

		secondClaimAttestation = s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		s.Require().NotNil(secondClaimAttestation)
		s.Require().False(secondClaimAttestation.Observed)
		s.Require().NotNil(secondClaimAttestation.Votes)
		s.Require().EqualValues(approveNumber, len(secondClaimAttestation.Votes))

		activeRelayers = s.Keeper().GetAllRelayers(s.Ctx, true)
		s.Require().NotNil(activeRelayers)
		s.Require().EqualValues(4, len(activeRelayers))
		for i := 0; i < 4; i++ {
			s.Require().NotNil(newRelayerList3[i], activeRelayers[i].RelayerAddress)
		}

		secondBridgeTokenClaim.RelayerAddress = s.relayerAddrs[2].String()
		_, err = s.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(s.Ctx), secondBridgeTokenClaim)
		s.Require().NoError(err)

		s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
		s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

		secondClaimAttestation = s.Keeper().GetAttestation(s.Ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		s.Require().NotNil(secondClaimAttestation)
		s.Require().True(secondClaimAttestation.Observed)
		s.Require().NotNil(secondClaimAttestation.Votes)
		s.Require().EqualValues(approveNumber+1, len(secondClaimAttestation.Votes))
	}
}

func (s *KeeperTestSuite) TestRelayerDelete() {
	proposalRelayers, found := s.Keeper().GetProposalRelayer(s.Ctx)
	s.Require().True(found)
	s.Require().EqualValues(s.relayerNumber, len(proposalRelayers.Relayers))
	nonce := s.Keeper().GetLastRelayerSetNonce(s.Ctx)
	s.Require().EqualValues(0, nonce)

	for i := 0; i < len(s.relayerAddrs); i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:  s.relayerAddrs[i].String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
			ChainName:       s.chainName,
		}
		s.Require().NoError(msgBondedRelayer.ValidateBasic())
		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)
		s.Require().NoError(err)
	}
	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	allRelayers := s.Keeper().GetAllRelayers(s.Ctx, false)
	s.Require().NotNil(allRelayers)
	s.Require().EqualValues(len(s.relayerAddrs), len(allRelayers))

	relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[0])
	s.Require().True(found)
	s.Require().EqualValues(s.relayerAddrs[0].String(), relayer.RelayerAddress)

	externalAddress := s.PubKeyToExternalAddr(s.externalPris[0].PublicKey)
	relayerAddr, found := s.Keeper().GetRelayerByExternalAddress(s.Ctx, externalAddress)
	s.Require().True(found)
	s.Require().EqualValues(relayer.RelayerAddress, relayerAddr.String())
	s.Require().True(sdkmath.NewInt(10 * 1e8).Equal(relayer.DelegateAmount))

	newRelayerAddressList := make([]string, 0, len(s.relayerAddrs)-1)
	for _, address := range s.relayerAddrs[1:] {
		newRelayerAddressList = append(newRelayerAddressList, address.String())
	}
	_, err := s.MsgServer().ProposalRelayers(s.Ctx, &types.MsgProposalRelayers{
		Relayers:  newRelayerAddressList,
		Authority: s.Dao.GlobalDao,
		ChainName: s.chainName,
	})
	s.Require().NoError(err)
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

	nonce = s.Keeper().GetLastRelayerSetNonce(s.Ctx)
	s.Require().EqualValues(2, nonce)
	relayerSet := s.Keeper().GetLastRelayerSet(s.Ctx)
	s.Require().EqualValues(nonce, relayerSet.Nonce)
	s.Require().EqualValues(9, len(relayerSet.Members))

	relayer, found = s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[0])
	s.Require().True(found)
	s.Require().Equal(relayerAddr.String(), relayer.RelayerAddress)
	s.Require().False(relayer.Online)

	relayerAddr, found = s.Keeper().GetRelayerByExternalAddress(s.Ctx, externalAddress)
	s.Require().True(found)
	s.Require().Equal(relayerAddr.String(), relayer.RelayerAddress)

	_, err = s.MsgServer().UnbondedRelayer(s.Ctx, &types.MsgUnbondedRelayer{
		ChainName:      s.chainName,
		RelayerAddress: s.relayerAddrs[0].String(),
	})
	s.Require().NoError(err)
	relayer, found = s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[0])
	s.Require().False(found)
	relayerAddr, found = s.Keeper().GetRelayerByExternalAddress(s.Ctx, externalAddress)
	s.Require().False(found)
}

func (s *KeeperTestSuite) TestRelayerSetSlash() {
	for i := 0; i < len(s.relayerAddrs); i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:  s.relayerAddrs[i].String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
			ChainName:       s.chainName,
		}
		s.Require().NoError(msgBondedRelayer.ValidateBasic())
		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)
		s.Require().NoError(err)
	}
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.Keeper().EndBlocker(s.Ctx)

	allRelayers := s.Keeper().GetAllRelayers(s.Ctx, false)
	s.Require().NotNil(allRelayers)
	s.Require().Equal(len(s.relayerAddrs), len(allRelayers))

	relayerSets := s.Keeper().GetRelayerSets(s.Ctx)
	s.Require().NotNil(relayerSets)
	s.Require().EqualValues(1, len(relayerSets))

	for i := 0; i < len(s.relayerAddrs)-1; i++ {
		externalAddress, signature := s.SignRelayerSetConfirm(s.externalPris[i], relayerSets[0])
		relayerSetConfirm := &types.MsgRelayerSetConfirm{
			Nonce:           relayerSets[0].Nonce,
			RelayerAddress:  s.relayerAddrs[i].String(),
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       s.chainName,
		}
		s.Require().NoError(relayerSetConfirm.ValidateBasic())
		_, err := s.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(s.Ctx), relayerSetConfirm)
		s.Require().NoError(err)
	}

	s.Keeper().EndBlocker(s.Ctx)
	relayerSetHeight := int64(relayerSets[0].Height)
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.App.EndBlock(abci.RequestEndBlock{Height: s.Ctx.BlockHeight()})

	relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[len(s.relayerAddrs)-1])
	s.Require().True(found)
	s.Require().True(relayer.Online)
	s.Require().Equal(int64(0), relayer.SlashTimes)

	s.Ctx = s.Ctx.WithBlockHeight(relayerSetHeight + int64(s.Keeper().GetParams(s.Ctx).SignedWindow) + 1)
	s.Keeper().EndBlocker(s.Ctx)

	relayer, found = s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[len(s.relayerAddrs)-1])
	s.Require().True(found)
	s.Require().False(relayer.Online)
	s.Require().Equal(int64(1), relayer.SlashTimes)
}

func (s *KeeperTestSuite) TestSlashRelayer() {
	for i := 0; i < len(s.relayerAddrs); i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:  s.relayerAddrs[i].String(),
			ExternalAddress: s.PubKeyToExternalAddr(s.externalPris[i].PublicKey),
			DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(10*1e8)),
			ChainName:       s.chainName,
		}
		s.Require().NoError(msgBondedRelayer.ValidateBasic())
		_, err := s.MsgServer().BondedRelayer(sdk.WrapSDKContext(s.Ctx), msgBondedRelayer)
		s.Require().NoError(err)
	}

	params := s.Keeper().GetParams(s.Ctx)
	err := s.Keeper().SetParams(s.Ctx, &params)
	s.Require().NoError(err)
	for i := 0; i < len(s.relayerAddrs); i++ {
		relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[i])
		s.Require().True(found)
		s.Require().True(relayer.Online)
		s.Require().Equal(int64(0), relayer.SlashTimes)

		s.Keeper().SlashRelayer(s.Ctx, relayer.RelayerAddress)

		relayer, found = s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[i])
		s.Require().True(found)
		s.Require().False(relayer.Online)
		s.Require().Equal(int64(1), relayer.SlashTimes)
	}

	// repeat slash test.
	for i := 0; i < len(s.relayerAddrs); i++ {
		relayer, found := s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[i])
		s.Require().True(found)
		s.Require().False(relayer.Online)
		s.Require().Equal(int64(1), relayer.SlashTimes)

		s.Keeper().SlashRelayer(s.Ctx, relayer.RelayerAddress)

		relayer, found = s.Keeper().GetRelayer(s.Ctx, s.relayerAddrs[i])
		s.Require().True(found)
		s.Require().False(relayer.Online)
		// not online, not change
		s.Require().Equal(int64(1), relayer.SlashTimes)
	}
}
