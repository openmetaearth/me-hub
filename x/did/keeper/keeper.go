package keeper

import (
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/did/types"
)

type Keeper struct {
	cdc       codec.Codec
	storeKey  storetypes.StoreKey
	daoKeeper types.DaoKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	daoKeeper types.DaoKeeper,
) *Keeper {
	return &Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		daoKeeper: daoKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) PubKeyFromString(s string) (pk cryptotypes.PubKey, err error) {
	err = k.cdc.UnmarshalInterfaceJSON([]byte(s), &pk)
	return pk, err
}

func (k Keeper) MustAccAddressFromPubkeyString(s string) sdk.AccAddress {
	pk, err := k.PubKeyFromString(s)
	if err != nil {
		panic(err)
	}

	return sdk.AccAddress(pk.Address())
}
