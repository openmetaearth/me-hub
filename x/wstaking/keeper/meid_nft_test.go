package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestRemoveMeidNFTDeletesRegionIndex() {
	s.SetupTest()

	account := s.TestAccs[0].String()
	regionId := "alpha"
	meidNFT := types.MeidNFT{
		Creator:    s.Dao.GlobalDao,
		Account:    account,
		RegionId:   regionId,
		RegionName: "Alpha",
		Umeid:      "umeid-alpha-1",
		NftId:      "nft-alpha-1",
	}

	s.Keeper().SetMeidNFT(s.Ctx, meidNFT)

	regionStore := prefix.NewStore(s.Ctx.KVStore(s.Keeper().GetStoreKey()), types.KeyPrefix(types.MeidNFTAccountKeyPrefix+regionId))
	s.Require().True(regionStore.Has(types.MeidNFTKey(account)))

	s.Keeper().RemoveMeidNFT(s.Ctx, account, regionId)

	_, found := s.Keeper().GetMeidNFT(s.Ctx, account)
	s.Require().False(found)
	s.Require().False(regionStore.Has(types.MeidNFTKey(account)))
}
