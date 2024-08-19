package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

// BaseKeeperWrapper is a wrapper of the cosmos-sdk bank module.
type BaseKeeperWrapper struct {
	bankkeeper.BaseKeeper
	ak banktypes.AccountKeeper
	dk types.DaoKeeper
}

// NewKeeper returns a new BaseKeeperWrapper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak banktypes.AccountKeeper,
	dk types.DaoKeeper,
	blockedAddrs map[string]bool,
	authority string,
) BaseKeeperWrapper {
	return BaseKeeperWrapper{
		BaseKeeper: bankkeeper.NewBaseKeeper(cdc, storeKey, ak, blockedAddrs, authority),
		ak:         ak,
		dk:         dk,
	}
}

// StakeCoinsFromModuleToModule stakes coins and transfers them from stake pool
// module account to a module account. It will panic if the module account
// does not exist or is unauthorized.
func (k BaseKeeperWrapper) StakeCoinsFromModuleToModule(
	ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins,
) error {
	senderAcc := k.ak.GetModuleAccount(ctx, senderModule)
	if senderAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}
	if !senderAcc.HasPermission(authtypes.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to send stake coins", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	if !recipientAcc.HasPermission(authtypes.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to receive stake coins", recipientModule))
	}

	return k.SendCoins(ctx, senderAcc.GetAddress(), recipientAcc.GetAddress(), amt)
}

// UnstakeCoinsFromModuleToModule unstakes the unbonding coins and transfers
// them from a module account to the stake_tokens_pool module's account. It will panic if the
// module account does not exist or is unauthorized.
func (k BaseKeeperWrapper) UnstakeCoinsFromModuleToModule(
	ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins,
) error {
	senderAcc := k.ak.GetModuleAccount(ctx, senderModule)
	if senderAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if !senderAcc.HasPermission(authtypes.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to send unstake coins", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientAcc))
	}

	if !recipientAcc.HasPermission(authtypes.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to receive unstake coins", recipientModule))
	}

	return k.SendCoins(ctx, senderAcc.GetAddress(), recipientAcc.GetAddress(), amt)
}
