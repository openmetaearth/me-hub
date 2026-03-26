package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	errorsmod "cosmossdk.io/errors"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/st-chain/me-hub/x/wbank/types"
)

// BaseKeeperWrapper is a wrapper of the cosmos-sdk bank module.
type BaseKeeperWrapper struct {
	bankkeeper.BaseKeeper
	ak banktypes.AccountKeeper
	dk types.DaoKeeper
}

// NewKeeper returns a new BaseKeeperWrapper instance.
func NewKeeper(
	appCodec codec.BinaryCodec,
	storeService store.KVStoreService,
	ak banktypes.AccountKeeper,
	dk types.DaoKeeper,
	moduleAccountAddrs map[string]bool,
	logger log.Logger,
) BaseKeeperWrapper {
	return BaseKeeperWrapper{
		BaseKeeper: bankkeeper.NewBaseKeeper(
			appCodec,
			storeService,
			ak,
			moduleAccountAddrs,
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			logger,
		),
		ak: ak,
		dk: dk,
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
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}
	if !senderAcc.HasPermission(authtypes.Staking) {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to send stake coins", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	if !recipientAcc.HasPermission(authtypes.Staking) {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to receive stake coins", recipientModule))
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
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if !senderAcc.HasPermission(authtypes.Staking) {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to send unstake coins", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientAcc))
	}

	if !recipientAcc.HasPermission(authtypes.Staking) {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to receive unstake coins", recipientModule))
	}

	return k.SendCoins(ctx, senderAcc.GetAddress(), recipientAcc.GetAddress(), amt)
}

func (k BaseKeeperWrapper) FeeToReceivers(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output, receiverTypes []types.FeeReceiverType) error {
	if len(inputs) == 0 {
		return errors.New("inputs error")
	}
	err := k.InputOutputCoins(ctx, inputs[0], outputs)
	if err != nil {
		return errorsmod.Wrap(err, "failed to process input-output coins")
	}

	if len(receiverTypes) != len(outputs) {
		return errorsmod.Wrap(err, "fee receiver types and outputs are not equal")
	}

	attributes := []sdk.Attribute{}
	attributes = append(attributes, sdk.NewAttribute(sdk.AttributeKeySender, inputs[0].Address))
	for index, output := range outputs {
		attributes = append(attributes, sdk.NewAttribute(fmt.Sprintf("%s", receiverTypes[index]), output.Address))
		attributes = append(attributes, sdk.NewAttribute(fmt.Sprintf("%s_amount", receiverTypes[index]), output.Coins.String()))
	}
	event := sdk.NewEvent(types.EventTypeFeeToReceivers, attributes...)
	ctx.EventManager().EmitEvent(event)
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
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if k.BlockedAddr(recipientAddr) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr)
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
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
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
		return errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule)
	}

	return k.SendCoinsWithTag(ctx, senderAddr, recipientAcc.GetAddress(), amt, tag...)
}

func (k BankKeeperExtend) SendCoinsWithTag(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins, tag ...string) error {
	err := k.SendCoins(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err
	}
	l := len(ctx.EventManager().Events()) - 2

	for _, t := range tag {
		ctx.EventManager().Events()[l].Attributes = append(ctx.EventManager().Events()[l].Attributes, abci.EventAttribute{
			Key:   "tag",
			Value: t,
		})
	}

	return nil
}
