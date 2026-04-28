package keeper

import (
	"errors"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/openmetaearth/me-hub/x/wbank/types"
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

func (k BaseKeeperWrapper) FeeToReceivers(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output, receiverTypes []types.FeeReceiverType) error {
	err := k.InputOutputCoins(ctx, inputs, outputs)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to process input-output coins")
	}

	if len(receiverTypes) != len(outputs) {
		return sdkerrors.Wrap(err, "fee receiver types and outputs are not equal")
	}

	attributes := []sdk.Attribute{}
	if len(inputs) > 0 {
		attributes = append(attributes, sdk.NewAttribute(sdk.AttributeKeySender, inputs[0].Address))
		for index, output := range outputs {
			attributes = append(attributes, sdk.NewAttribute(fmt.Sprintf("%s", receiverTypes[index]), output.Address))
			attributes = append(attributes, sdk.NewAttribute(fmt.Sprintf("%s_amount", receiverTypes[index]), output.Coins.String()))
		}
		event := sdk.NewEvent(types.EventTypeFeeToReceivers, attributes...)
		ctx.EventManager().EmitEvent(event)
	} else {
		return errors.New("inputs error")
	}
	return nil
}

func (k BaseKeeperWrapper) Extend() BankKeeperExtend {
	return BankKeeperExtend{
		BaseKeeperWrapper: k,
	}
}

type BankKeeperExtend struct {
	BaseKeeperWrapper
}

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist. An error is returned if
// the recipient address is black-listed or if sending the tokens fails.
func (k BankKeeperExtend) SendCoinsFromModuleToAccountWithTag(
	ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins, tag ...string,
) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if k.BlockedAddr(recipientAddr) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr)
	}

	return k.SendCoinsWithTag(ctx, senderAddr, recipientAddr, amt, tag...)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another.
// It will panic if either module account does not exist.
func (k BankKeeperExtend) SendCoinsFromModuleToModuleWithTag(
	ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins, tag ...string,
) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.SendCoinsWithTag(ctx, senderAddr, recipientAcc.GetAddress(), amt, tag...)
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (k BankKeeperExtend) SendCoinsFromAccountToModuleWithTag(
	ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins, tag ...string,
) error {
	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule)
	}

	return k.SendCoinsWithTag(ctx, senderAddr, recipientAcc.GetAddress(), amt, tag...)
}

func (k BankKeeperExtend) SendCoinsWithTag(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins, tags ...string) error {
	if len(tags) == 0 {
		return k.SendCoins(ctx, fromAddr, toAddr, amt)
	}

	beforeCount := len(ctx.EventManager().Events())
	err := k.SendCoins(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err
	}

	// Find the transfer event emitted by SendCoins (search from the new events only)
	events := ctx.EventManager().Events()
	transferIdx := -1
	for i := len(events) - 1; i >= beforeCount; i-- {
		if events[i].Type == banktypes.EventTypeTransfer {
			transferIdx = i
			break
		}
	}
	if transferIdx == -1 {
		return nil
	}

	for _, t := range tags {
		events[transferIdx].Attributes = append(events[transferIdx].Attributes, abci.EventAttribute{
			Key:   "tag",
			Value: t,
		})
	}

	return nil
}
