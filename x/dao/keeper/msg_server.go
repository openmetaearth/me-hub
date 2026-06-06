package keeper

import (
	"context"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/openmetaearth/me-hub/x/dao/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdateDao(goCtx context.Context, msg *types.MsgUpdateDao) (*types.MsgUpdateDaoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isGlobalDao := k.IsGlobalDao(ctx, msg.Creator)
	if !isGlobalDao {
		return nil, types.ErrCreatorNotDao
	}

	oldDao, found := k.GetDaoAddresses(ctx)
	if !found {
		return nil, types.ErrNotFound
	}

	if err := k.validateDaoAddresses(ctx, msg.DaoAddresses); err != nil {
		return nil, err
	}

	k.SetDaoAddresses(ctx, msg.DaoAddresses)

	err := k.kycHook.SetKycIssers(ctx, []string{oldDao.GlobalDao, oldDao.MeidDao}, []string{msg.DaoAddresses.GlobalDao, msg.DaoAddresses.MeidDao})
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrSetKycIssuer, err.Error())
	}

	oldByte, _ := json.Marshal(oldDao)
	newByte, _ := json.Marshal(msg.DaoAddresses)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDaoUpdated,
			sdk.NewAttribute(types.AttributeKeyLastDaoAddresses, string(oldByte)),
			sdk.NewAttribute(types.AttributeKeyNewDaoAddresses, string(newByte)),
		),
	)
	return &types.MsgUpdateDaoResponse{}, nil
}

type accountGetter interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

func (k msgServer) validateDaoAddresses(ctx sdk.Context, daoAddresses types.DaoAddresses) error {
	if err := validateDaoAddressNotModuleAccount(ctx, k.authKeeper, "GlobalDao", daoAddresses.GlobalDao); err != nil {
		return err
	}
	if err := validateDaoAddressNotModuleAccount(ctx, k.authKeeper, "MeidDao", daoAddresses.MeidDao); err != nil {
		return err
	}
	if err := validateDaoAddressNotModuleAccount(ctx, k.authKeeper, "DevOperator", daoAddresses.DevOperator); err != nil {
		return err
	}
	return validateDaoAddressNotModuleAccount(ctx, k.authKeeper, "AirdropAddress", daoAddresses.AirdropAddress)
}

func validateDaoAddressNotModuleAccount(ctx sdk.Context, accountKeeper accountGetter, fieldName, address string) error {
	if address == "" {
		return nil
	}

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "%s: %s", fieldName, address)
	}

	account := accountKeeper.GetAccount(ctx, addr)
	if _, ok := account.(authtypes.ModuleAccountI); ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "%s cannot be a module account: %s", fieldName, address)
	}
	return nil
}

func (k msgServer) FreeGasAccount(goCtx context.Context, msg *types.MsgFreeGasAccount) (*types.MsgFreeGasAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isGlobalDao := k.IsGlobalDao(ctx, msg.Creator)
	if !isGlobalDao {
		return nil, types.ErrCreatorNotDao
	}

	for _, account := range msg.Accounts {
		isExist := k.CheckFreeGasAccount(ctx, account.Address)
		if isExist {
			if account.IsFree {
				return nil, errorsmod.Wrap(types.ErrFreeGasAccountAlreadyExist, account.Address)
			} else {
				k.RemoveFreeGasAccount(ctx, account.Address)
			}
		}

		if !isExist {
			if account.IsFree {
				k.SetFreeGasAccount(ctx, account.Address)
			} else {
				return nil, errorsmod.Wrap(types.ErrAccountIsNotFree, account.Address)
			}
		}
	}

	return &types.MsgFreeGasAccountResponse{}, nil
}
