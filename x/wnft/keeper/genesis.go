package keeper

import (
	"cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis loads nft genesis. Classes or tokens already created by other
// modules (kyc SBT class, region classes) are skipped so export/import works
// when those modules InitGenesis first.
func (k Keeper) InitGenesis(ctx sdk.Context, data *nft.GenesisState) {
	for _, class := range data.Classes {
		if class == nil {
			continue
		}
		if k.HasClass(ctx, class.Id) {
			continue
		}
		if err := k.SaveClass(ctx, *class); err != nil {
			panic(err)
		}
	}

	for _, entry := range data.Entries {
		if entry == nil {
			continue
		}
		for _, token := range entry.Nfts {
			if token == nil {
				continue
			}
			if k.HasNFT(ctx, token.ClassId, token.Id) {
				continue
			}
			owner, err := sdk.AccAddressFromBech32(entry.Owner)
			if err != nil {
				panic(err)
			}
			if err := k.Mint(ctx, *token, owner); err != nil {
				panic(err)
			}
		}
	}
}
