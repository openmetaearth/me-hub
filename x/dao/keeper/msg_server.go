package keeper

import (
	"context"
	sdkerrors "cosmossdk.io/errors"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	k.SetDaoAddresses(ctx, msg.DaoAddresses)

	err := k.kycHook.SetKycIssers(ctx, []string{oldDao.GlobalDao, oldDao.MeidDao}, []string{msg.DaoAddresses.GlobalDao, msg.DaoAddresses.MeidDao})
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrSetKycIssuer, err.Error())
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
				return nil, sdkerrors.Wrap(types.ErrFreeGasAccountAlreadyExist, account.Address)
			} else {
				k.RemoveFreeGasAccount(ctx, account.Address)
			}
		}

		if !isExist {
			if account.IsFree {
				k.SetFreeGasAccount(ctx, account.Address)
			} else {
				return nil, sdkerrors.Wrap(types.ErrAccountIsNotFree, account.Address)
			}
		}
	}

	return &types.MsgFreeGasAccountResponse{}, nil
}
