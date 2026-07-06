package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

func (m msgServer) CreateSubAccount(goCtx context.Context, msg *types.MsgCreateSubAccount) (*types.MsgCreateSubAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sdkAccount, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgCreateSubAccountResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	
	did, ok := m.GetDID(ctx, sdkAccount)
	if !ok {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrDidNotFound
	}

	holderInfo, found := m.GetDidInfo(ctx, did)
	if !found || holderInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrHolderNotFound
	}

	if !m.HasKYC(ctx, did) {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrCredentialNotFound
	}


	if holderInfo.Address != msg.Creator {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrUnauthorized
	}

	if holderInfo.SubAccount != "" {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountAlreadyExists
	}

	if m.didKeeper.HasDidBySubAccount(ctx, msg.SubAccount) {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountAlreadyRegistered
	}

	// validate that sub_account is an ethsecp256k1 address derived from the provided pubkey
	pk, err := m.PubKeyFromString(msg.Pubkey)
	if err != nil {
		return &types.MsgCreateSubAccountResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
	}

	if _, ok := pk.(*ethsecp256k1.PubKey); !ok {
		return &types.MsgCreateSubAccountResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "sub_account must be an ethsecp256k1 address")
	}

	derivedAddr := sdk.AccAddress(pk.Address())

	subAccount, err := sdk.AccAddressFromBech32(msg.SubAccount)
	if err != nil {
		return &types.MsgCreateSubAccountResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if !derivedAddr.Equals(subAccount) {
		return &types.MsgCreateSubAccountResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey does not match sub_account address")
	}

	if !m.accountKeeper.HasAccount(ctx, subAccount) {
		newAccount := m.accountKeeper.NewAccountWithAddress(ctx, subAccount)
		err = newAccount.SetPubKey(pk)
		if err != nil {
			return &types.MsgCreateSubAccountResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
		}
		m.accountKeeper.SetAccount(ctx, newAccount)
	}

	holderInfo.SubAccount = msg.SubAccount
	m.SetDidInfo(ctx, did, holderInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateSubAccount,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeySubAccount, msg.SubAccount),
			sdk.NewAttribute(types.AttributeKeyDid, did),
		),
	})
	return &types.MsgCreateSubAccountResponse{}, nil
}
