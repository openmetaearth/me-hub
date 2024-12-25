package keeper

import (
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/kyc/handler"
	"github.com/st-chain/me-hub/x/kyc/types"
)

// Keeper is wrapper of did keeper and nft keeper.
type Keeper struct {
	cdc       codec.Codec
	storeKey  storetypes.StoreKey
	stkKeeper types.StakingKeeper
	didKeeper types.DIDKeeper
	nftKeeper types.NFTKeeper

	handlerReg *handler.HandlerRegistry
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	stkKeeper types.StakingKeeper,
	didKeeper types.DIDKeeper,
	nftKeeper types.NFTKeeper,
) *Keeper {
	return &Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		stkKeeper: stkKeeper,
		didKeeper: didKeeper,
		nftKeeper: nftKeeper,

		handlerReg: handler.NewEventRegistry(),
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) PubKeyFromString(s string) (pk cryptotypes.PubKey, err error) {
	err = k.cdc.UnmarshalInterfaceJSON([]byte(s), &pk)
	return pk, err
}

func (k *Keeper) MustAccAddressFromPubkeyString(s string) sdk.AccAddress {
	if len(s) > 0 {
		pk, err := k.PubKeyFromString(s)
		if err != nil {
			panic(err)
		}
		return sdk.AccAddress(pk.Address())
	}
	return sdk.AccAddress{}
}

func (k *Keeper) RegisterEventHandler(eventType string, priority int, module string, handler handler.HandlerFunc) {
	k.handlerReg.RegisterHandler(eventType, priority, module, handler)
}
