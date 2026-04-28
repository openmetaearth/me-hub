package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

/*
todo: remove these functions, for now we only need to support SetSbtClass
*/

func (k *Keeper) HasSBT(ctx sdk.Context, did string) bool {
	return k.nftKeeper.HasNFT(ctx, types.ModuleName, did)
}

func (k *Keeper) GetSBT(ctx sdk.Context, did string) (sbt nft.NFT, found bool) {
	return k.nftKeeper.GetNFT(ctx, types.ModuleName, did)
}

//func (k *Keeper) GetAllSBTs(ctx sdk.Context, regionId string) (SBTs []nft.NFT) {
//	// todo: implement, for export genesis
//	return SBTs
//}

func (k *Keeper) SetSBT(ctx sdk.Context, sbt nft.NFT, receiver sdk.AccAddress) error {
	return k.nftKeeper.Mint(ctx, sbt, receiver)
}

func (k *Keeper) RemoveSBT(ctx sdk.Context, did string) error {
	return k.nftKeeper.Burn(ctx, types.ModuleName, did)
}

func (k *Keeper) SetSbtClass(ctx sdk.Context) error {
	class := nft.Class{
		Id:   types.ModuleName,
		Name: types.ModuleName,
	}

	return k.nftKeeper.SaveClass(ctx, class)
}
