package v2_0_10_test

import (
	"encoding/json"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/app/upgrades/v2_0_10"
	"github.com/st-chain/me-hub/utils"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
	"testing"
	"time"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/st-chain/me-hub/app"
	"github.com/st-chain/me-hub/app/apptesting"
)

// UpgradeTestSuite defines the structure for the upgrade test suite
type UpgradeTestSuite struct {
	suite.Suite
	Ctx                 sdk.Context
	App                 *app.App
	DIDReader           v2_0_10.DIDReader
	KycPubkeyReader     v2_0_10.KycPubkeyReader
	NftReader           v2_0_10.NftReader
	mockAddress1        string
	mockAddress2        string
	mockAddress3        string
	mockDid1            string
	meEarthValidator    stakingtypes.Validator
	experienceValidator stakingtypes.Validator
	usaValidator        stakingtypes.Validator
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) SetupTest() {
	s.App = apptesting.Setup(s.T(), false)
	s.Ctx = s.App.BaseApp.NewContext(false, cometbftproto.Header{Height: 1, ChainID: "mechain_100-1", Time: time.Now().UTC()})
	s.mockAddress1 = "me18uungln5qxndqelavjnzny6z87v530hmm74dtm" // maple skin present glad name second struggle correct submit learn guitar refuse common become sphere output pattern annual riot master tent buddy aisle abuse
	s.mockAddress2 = "me1lepklhvskft6cr5e0lzce0vwgpjtq6esmkjdga" // upon gate call badge film access impact adjust slow uncle trust path remove drip pulp pact already grape mouse benefit era bridge annual frost
	s.mockAddress3 = "me1vlljzn8mhy68dgatgl8hn24a0sgx2vk8ffm8hd" // upon gate call badge film access impact adjust slow uncle trust path remove drip pulp pact already grape mouse benefit era bridge annual frost
	s.mockDid1 = "9998887776660"
	s.DIDReader = MockDIDReader{
		Data: map[string]v2_0_10.DidData{
			s.mockAddress1: {
				Did:        s.mockDid1,
				Uri:        "https://example.com/nft1/metadata.json",
				UriHash:    "e00d196344dbd54550dadeab1167302ef39fade96eb211e302693a512ef131e1",
				KycUri:     "https://example.com/kyc/metadata.json",
				KycUriHash: "e00d196344dbd54550dadeab1167302ef39fade96eb211e302693a512ef131e2",
				PubKey:     "{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"Anmi0DiLED1oGRiVIPO4n6HSnk7iArQBdeR1HnHxodmB\"}",
			},
		},
		Err: nil,
	}
	s.KycPubkeyReader = MockKycPubkeyReader{
		Data: map[string]string{
			s.mockAddress1: "{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"Anmi0DiLED1oGRiVIPO4n6HSnk7iArQBdeR1HnHxodmB\"}",
		},
	}

	validators := s.App.StakingKeeper.GetValidators(s.Ctx, 10)
	s.Require().True(len(validators) >= 3)
	s.meEarthValidator = validators[0]
	s.experienceValidator = validators[1]
	s.usaValidator = validators[2]
}

func (s *UpgradeTestSuite) TestReadDidData() {
	data, err := s.DIDReader.ReadDID("dummy_path")
	s.Require().NoError(err)
	first, ok := data[s.mockAddress1]
	s.Require().True(ok)
	s.Require().EqualValues(s.mockDid1, first.Did)
}

func (s *UpgradeTestSuite) TestReadKycPubkey() {
	data, err := s.KycPubkeyReader.ReadKycPubkey("dummy_path")
	s.Require().NoError(err)
	first, ok := data[s.mockAddress1]
	s.Require().True(ok)
	s.Require().EqualValues("{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"Anmi0DiLED1oGRiVIPO4n6HSnk7iArQBdeR1HnHxodmB\"}", first)
}

func (s *UpgradeTestSuite) TestUpgradeMigrateKycData() {
	s.App.StakingKeeper.SetMeid(s.Ctx, wstakingtypes.Meid{
		Account:    s.mockAddress1,
		Creator:    "",
		RegionId:   wstakingtypes.MeEarthRegionId,
		RegionName: wstakingtypes.MeEarthRegionName,
		RewardType: 0,
	})

	s.App.StakingKeeper.SetMeid(s.Ctx, wstakingtypes.Meid{
		Account:    s.mockAddress2,
		Creator:    "",
		RegionId:   wstakingtypes.MeEarthRegionId,
		RegionName: wstakingtypes.MeEarthRegionName,
		RewardType: 0,
	})

	s.App.StakingKeeper.SetMeid(s.Ctx, wstakingtypes.Meid{
		Account:    s.mockAddress3,
		Creator:    "",
		RegionId:   wstakingtypes.MeEarthRegionId,
		RegionName: wstakingtypes.MeEarthRegionName,
		RewardType: 0,
	})

	// Call the MigrateKycData function
	v2_0_10.MigrateKycData(s.Ctx,
		s.App.StakingKeeper,
		s.App.DidKeeper,
		s.App.KycKeeper,
		s.App.WNFTKeeper,
		s.App.GroupKeeper,
		"dummy_path",
		s.DIDReader,
		s.KycPubkeyReader)

	// Verify the DID data for mockAddress1
	did, found := s.App.DidKeeper.GetDID(s.Ctx, sdk.MustAccAddressFromBech32(s.mockAddress1))
	s.Require().True(found)
	s.Require().Equal(s.mockDid1, did)

	// Verify the KYC public key data for mockAddress1
	didInfo, found := s.App.DidKeeper.GetDidInfo(s.Ctx, s.mockDid1)
	s.Require().True(found)
	s.Require().Equal("{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"Anmi0DiLED1oGRiVIPO4n6HSnk7iArQBdeR1HnHxodmB\"}", didInfo.Pubkey)

	// Verify the NFT URI and URI hash for mockAddress1
	sbt, found := s.App.KycKeeper.GetSBT(s.Ctx, s.mockDid1)
	s.Require().True(found)
	s.Require().Equal("https://example.com/nft1/metadata.json", sbt.Uri)
	s.Require().Equal("e00d196344dbd54550dadeab1167302ef39fade96eb211e302693a512ef131e1", sbt.UriHash)

	// Verify the KYC URI and URI hash for mockAddress1
	kyc, found := s.App.KycKeeper.GetKYC(s.Ctx, s.mockDid1)
	s.Require().True(found)
	s.Require().Equal("https://example.com/kyc/metadata.json", kyc.Uri)
	s.Require().Equal("e00d196344dbd54550dadeab1167302ef39fade96eb211e302693a512ef131e2", kyc.Hash)

	// check kyc filters
	filters, f := s.App.KycKeeper.GetFilters(s.Ctx, s.mockDid1)
	s.Require().True(f)
	s.Require().Len(filters, 1)
	s.Require().EqualValues(wstakingtypes.MeEarthRegionId, string(filters[0]))

	// Verify that mockAddress2 has no DidData and pubkey data
	did2, did2Found := s.App.DidKeeper.GetDID(s.Ctx, sdk.MustAccAddressFromBech32(s.mockAddress2))
	s.Require().True(did2Found)
	//s.Require().EqualValues("9998887776660", did2)
	s.T().Log(did2)
	did2Info, did2Found := s.App.DidKeeper.GetDidInfo(s.Ctx, did2)
	s.Require().True(did2Found)
	s.Require().EqualValues(wstakingtypes.MeEarthRegionId, did2Info.RegionId)

	// Verify the NFT URI and URI hash for mockAddress1
	sbt, found = s.App.KycKeeper.GetSBT(s.Ctx, did2)
	s.Require().True(found)
	s.Require().EqualValues("kyc", sbt.GetClassId())

	// Verify the KYC URI and URI hash for mockAddress1
	kyc, found = s.App.KycKeeper.GetKYC(s.Ctx, did2)
	s.Require().True(found)
	s.Require().EqualValues(wstakingtypes.MeEarthRegionId, string(kyc.GetData()))

	// Verify that mockAddress3 has no DidData and pubkey data
	did3, did3Found := s.App.DidKeeper.GetDID(s.Ctx, sdk.MustAccAddressFromBech32(s.mockAddress3))
	s.Require().True(did3Found)
	s.T().Log(did3)

	did3Info, did3Found := s.App.DidKeeper.GetDidInfo(s.Ctx, did3)
	s.Require().True(did3Found)
	s.Require().EqualValues(wstakingtypes.MeEarthRegionId, did3Info.RegionId)

	_, meidFound := s.App.StakingKeeper.GetMeid(s.Ctx, s.mockAddress1)
	s.Require().False(meidFound)

	_, meidFound = s.App.StakingKeeper.GetMeid(s.Ctx, s.mockAddress2)
	s.Require().False(meidFound)

	_, meidFound = s.App.StakingKeeper.GetMeid(s.Ctx, s.mockAddress3)
	s.Require().False(meidFound)

	// check kyc filters
	filters, f = s.App.KycKeeper.GetFilters(s.Ctx, did2)
	s.Require().True(f)
	s.Require().Len(filters, 1)
	s.Require().EqualValues(wstakingtypes.MeEarthRegionId, string(filters[0]))
}

func (s *UpgradeTestSuite) TestDidData() {
	list := make(map[string]v2_0_10.DidData)
	list["me1ujufste3u23tpq3qhlq77u94nhw99emy3pr4p2"] = v2_0_10.DidData{
		Did:     "2405208027001",
		Uri:     "https://example.com/nft/metadata.json",
		UriHash: utils.CalculateUriHash("https://example.com/nft/metadata.json"),
	}
	list["me1phcakjkaf9vrn6jgttl3747dgnnpt88rt9440d"] = v2_0_10.DidData{
		Did:     "CHN2405204091002",
		Uri:     "https://example.com/nft1/metadata.json",
		UriHash: utils.CalculateUriHash("https://example.com/nft1/metadata.json"),
	}
	marshal, err := json.MarshalIndent(list, "", "  ")
	s.Require().NoError(err)
	s.T().Log(string(marshal))
}

func (s *UpgradeTestSuite) TestMigrateFixedDeposit() {
	s.App.StakingKeeper.SetMeid(s.Ctx, wstakingtypes.Meid{
		Account:    s.mockAddress1,
		Creator:    "",
		RegionId:   wstakingtypes.MeEarthRegionId,
		RegionName: wstakingtypes.MeEarthRegionName,
		RewardType: 0,
	})

	// Set up initial state with mock data
	s.App.StakingKeeper.SetRegion(s.Ctx, wstakingtypes.Region{
		RegionId: wstakingtypes.MeEarthRegionId,
	})

	s.App.StakingKeeper.SetFixedDeposit(s.Ctx, wstakingtypes.FixedDeposit{
		Account:   s.mockAddress1,
		Principal: sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000)),
	})

	s.App.BankKeeper.SendCoins(s.Ctx,
		authtypes.NewModuleAddress(wstakingtypes.StakePoolName),
		authtypes.NewModuleAddress(wstakingtypes.FixedDepositPrincipalPool),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(1000))))

	// Call the MigrateFixedDeposit function
	v2_0_10.MigrateFixedDeposit(s.Ctx, s.App.StakingKeeper, s.App.KycKeeper, s.App.BankKeeper)

	// Verify the region's fixed deposit amount is updated
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, wstakingtypes.MeEarthRegionId)
	s.Require().True(found)
	s.Require().Equal(sdk.NewInt(1000).String(), region.FixedDepositAmount.String())

	// Verify the total deposit amount is equal to the balance
	balance := s.App.BankKeeper.GetBalance(s.Ctx, authtypes.NewModuleAddress(wstakingtypes.FixedDepositPrincipalPool), params.BaseDenom)
	s.Require().Equal(sdk.NewInt(1000).String(), balance.Amount.String())
}

func (s *UpgradeTestSuite) TestMigrateDelegation() {
	// Set up the experience region
	expRegion := wstakingtypes.Region{
		RegionId:        wstakingtypes.ExperienceRegionId,
		OperatorAddress: s.experienceValidator.OperatorAddress,
	}
	s.App.StakingKeeper.SetRegion(s.Ctx, expRegion)

	testRegion := wstakingtypes.Region{
		RegionId:        wstakingtypes.MeEarthRegionId,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	s.App.StakingKeeper.SetRegion(s.Ctx, testRegion)

	// Set up a delegator and delegation
	delegator := sdk.MustAccAddressFromBech32(s.mockAddress1)
	testDid := "0000000000011"
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: delegator.String(),
		ValidatorAddress: "",
	})

	// Verify the initial state (no OperatorAddress)
	delegation, found := s.App.StakingKeeper.GetDelegation(s.Ctx, delegator, sdk.ValAddress{})
	s.Require().True(found)
	s.Require().Equal("", delegation.ValidatorAddress)

	// Set up KYC data for the delegator
	s.App.KycKeeper.SetDID(s.Ctx, delegator, testDid)
	s.App.KycKeeper.SetKYC(s.Ctx, testDid, didtypes.Credential{
		Did:  testDid,
		Sid:  "kyc",
		Data: []byte(wstakingtypes.MeEarthRegionId),
	})

	// Call the MigrateDelegation function
	v2_0_10.MigrateDelegation(s.Ctx, s.App.StakingKeeper, s.App.KycKeeper)

	delegation, found = s.App.StakingKeeper.GetDelegation(s.Ctx, delegator, s.meEarthValidator.GetOperator())
	s.Require().True(found)
	s.Require().Equal(s.meEarthValidator.OperatorAddress, delegation.ValidatorAddress)

	// Case 2: No Experience Region Validator Address
	delegator2 := sdk.MustAccAddressFromBech32(s.mockAddress2)
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: delegator2.String(),
		ValidatorAddress: "",
	})

	v2_0_10.MigrateDelegation(s.Ctx, s.App.StakingKeeper, s.App.KycKeeper)

	delegation2, found := s.App.StakingKeeper.GetDelegation(s.Ctx, delegator2, s.experienceValidator.GetOperator())
	s.Require().True(found)
	s.Require().Equal(expRegion.OperatorAddress, delegation2.ValidatorAddress)

	// Case 3: Validator Address Needs to Be Changed
	delegator3 := sdk.MustAccAddressFromBech32(s.mockAddress3)
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: delegator3.String(),
		ValidatorAddress: s.experienceValidator.OperatorAddress,
	})

	s.App.KycKeeper.SetDID(s.Ctx, delegator3, "testDid3")
	s.App.KycKeeper.SetKYC(s.Ctx, "testDid3", didtypes.Credential{
		Did:  "testDid3",
		Sid:  "kyc",
		Data: []byte(wstakingtypes.MeEarthRegionId),
	})

	v2_0_10.MigrateDelegation(s.Ctx, s.App.StakingKeeper, s.App.KycKeeper)

	delegation3, found := s.App.StakingKeeper.GetDelegation(s.Ctx, delegator3, s.meEarthValidator.GetOperator())
	s.Require().True(found)
	s.Require().Equal(testRegion.OperatorAddress, delegation3.ValidatorAddress)
}

func (s *UpgradeTestSuite) TestMigrateNftUri() {
	// Setup initial state with 3 NFT classes and 6 NFTs
	class1 := "class1"
	class2 := "class2"
	class3 := "class3"

	s.NftReader = MockNftReader{
		Data: map[string]map[string]v2_0_10.NftUri{
			class1: {
				"nft1": {URI: "https://example.com/new_nft1", URIHash: "new_hash1"},
				"nft2": {URI: "https://example.com/new_nft2", URIHash: "new_hash2"},
			},
			class2: {
				"nft1": {URI: "https://example.com/new_nft1", URIHash: "new_hash1"},
				"nft2": {URI: "https://example.com/new_nft2", URIHash: "new_hash2"},
			},
			class3: {
				"nft1": {URI: "https://example.com/new_nft1", URIHash: "new_hash1"},
			},
		},
		Err: nil,
	}

	nftUriData := map[string]map[string]v2_0_10.NftUri{
		class1: {
			"nft1": {URI: "https://example.com/class1/nft1", URIHash: utils.CalculateUriHash("https://example.com/class1/nft1")},
			"nft2": {URI: "https://example.com/class1/nft2", URIHash: utils.CalculateUriHash("https://example.com/class1/nft2")},
		},
		class2: {
			"nft1": {URI: "https://example.com/class2/nft1", URIHash: utils.CalculateUriHash("https://example.com/class2/nft1")},
			"nft2": {URI: "https://example.com/class2/nft2", URIHash: utils.CalculateUriHash("https://example.com/class2/nft2")},
		},
		class3: {
			"nft1": {URI: "https://example.com/class3/nft1", URIHash: utils.CalculateUriHash("https://example.com/class3/nft1")},
			"nft2": {URI: "https://example.com/class3/nft2", URIHash: utils.CalculateUriHash("https://example.com/class3/nft2")},
		},
	}

	// Save the NFT classes and NFTs in the keeper
	for classID, nftData := range nftUriData {
		err := s.App.WNFTKeeper.SaveClass(s.Ctx, nft.Class{Id: classID})
		s.Require().NoError(err)
		for nftID, n := range nftData {
			err = s.App.WNFTKeeper.Mint(s.Ctx, nft.NFT{
				ClassId: classID,
				Id:      nftID,
				Uri:     n.URI,
				UriHash: n.URIHash,
			},
				sdk.MustAccAddressFromBech32(s.mockAddress1))
			s.Require().NoError(err)
		}
	}

	// Call the MigrateNftUri function
	v2_0_10.MigrateNftUri(s.Ctx, s.App.WNFTKeeper, "dummy_path", s.NftReader)

	expectedNft, err := s.NftReader.ReadNft("dummy_path")
	s.Require().NoError(err)

	// Verify the URI changes for each NFT
	for classID, nftData := range nftUriData {
		for nftID, _ := range nftData {
			updatedNFT, found := s.App.WNFTKeeper.GetNFT(s.Ctx, classID, nftID)
			s.Require().True(found, "NFT not found: %s/%s", classID, nftID)

			expectedURI := expectedNft[classID][nftID].URI
			expectedURIHash := expectedNft[classID][nftID].URIHash
			if classID == class3 && nftID == "nft2" {
				expectedURI = nftUriData[classID][nftID].URI
				expectedURIHash = nftUriData[classID][nftID].URIHash
			}

			s.Require().Equal(expectedURI, updatedNFT.Uri, "URI mismatch for NFT: %s/%s", classID, nftID)
			s.Require().Equal(expectedURIHash, updatedNFT.UriHash, "URI hash mismatch for NFT: %s/%s", classID, nftID)
		}
	}
}
