package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/openmetaearth/me-hub/x/kyc/handler"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

// Keeper is wrapper of did keeper and nft keeper.
type Keeper struct {
	cdc           codec.Codec
	storeKey      storetypes.StoreKey
	stkKeeper     types.StakingKeeper
	accountKeeper authkeeper.AccountKeeper
	didKeeper     types.DIDKeeper
	nftKeeper     types.NFTKeeper

	handlerReg *handler.HandlerRegistry
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	stkKeeper types.StakingKeeper,
	accountKeeper authkeeper.AccountKeeper,
	didKeeper types.DIDKeeper,
	nftKeeper types.NFTKeeper,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		stkKeeper:     stkKeeper,
		accountKeeper: accountKeeper,
		didKeeper:     didKeeper,
		nftKeeper:     nftKeeper,

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

func (k *Keeper) MustAccAddressFromPubkeyString(s string) (sdk.AccAddress, error) {
	if len(s) > 0 {
		pk, err := k.PubKeyFromString(s)
		if err != nil {
			return sdk.AccAddress{}, err
		}
		return sdk.AccAddress(pk.Address()), nil
	}
	return sdk.AccAddress{}, fmt.Errorf("pubkey is empty")
}

func (k *Keeper) RegisterEventHandler(eventType string, priority int, module string, handler handler.HandlerFunc) {
	k.handlerReg.RegisterHandler(eventType, priority, module, handler)
}
