package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

func (m msgServer) CreateSubAccount(goCtx context.Context, msg *types.MsgCreateSubAccount) (*types.MsgCreateSubAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	holderInfo, found := m.GetDidInfo(ctx, msg.Did)
	if !found || holderInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrHolderNotFound
	}

	if !m.HasKYC(ctx, msg.Did) {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrCredentialNotFound
	}

	if holderInfo.SubAccount != "" {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountAlreadyExists
	}

	if holderInfo.SubAccount == msg.SubAccount {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountAlreadyRegistered
	}

	if m.didKeeper.HasDidBySubAccount(ctx, msg.SubAccount) {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountAlreadyRegistered
	}

	// check sub account address is valid and not exist
	subAccount, err := sdk.AccAddressFromBech32(msg.SubAccount)
	if err != nil {
		return &types.MsgCreateSubAccountResponse{}, err
	}
	
	if !m.accountKeeper.HasAccount(ctx, subAccount) {
		m.accountKeeper.SetAccount(ctx, m.accountKeeper.NewAccountWithAddress(ctx, subAccount))
	}


	holderInfo.SubAccount = msg.SubAccount
	m.SetDidInfo(ctx, msg.Did, holderInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateSubAccount,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeySubAccount, msg.SubAccount),
			sdk.NewAttribute(types.AttributeKeyDid, msg.Did),
		),
	})
	return &types.MsgCreateSubAccountResponse{}, nil
}
