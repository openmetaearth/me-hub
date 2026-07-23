package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

func (m msgServer) CreateSubAccount(goCtx context.Context, msg *types.MsgCreateSubAccount) (*types.MsgCreateSubAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creatorAccount, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgCreateSubAccountResponse{}, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	did, ok := m.GetDID(ctx, creatorAccount)
	if !ok {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrDidNotFound
	}

	holderInfo, found := m.GetDidInfo(ctx, did)
	if !found {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrHolderNotFound
	}

	if holderInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrDidNotActive
	}

	if !m.HasKYC(ctx, did) {
		return &types.MsgCreateSubAccountResponse{}, didtypes.ErrCredentialNotFound
	}

	account := m.accountKeeper.GetAccount(ctx, creatorAccount)
	if account == nil {
		return &types.MsgCreateSubAccountResponse{}, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "account not found")
	}

	pk := account.GetPubKey()
	if pk == nil {
		return &types.MsgCreateSubAccountResponse{}, types.ErrMainAccountPubkeyNotSet
	}

	if _, ok := pk.(*ethsecp256k1.PubKey); ok {
		return &types.MsgCreateSubAccountResponse{}, types.ErrEthAccountNotAllowed
	}

	if holderInfo.SubAccount != "" {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountAlreadyExists
	}

	if m.didKeeper.HasDidBySubAccount(ctx, msg.SubAccount) {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountAlreadyRegistered
	}

	// validate that sub_account is an ethsecp256k1 address derived from the provided pubkey
	subAccountPubKey, err := m.PubKeyFromString(msg.SubAccountPubkey)
	if err != nil {
		return &types.MsgCreateSubAccountResponse{}, errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
	}

	if _, ok := subAccountPubKey.(*ethsecp256k1.PubKey); !ok {
		return &types.MsgCreateSubAccountResponse{}, errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, "sub_account must be an ethsecp256k1 address")
	}

	derivedAddr := sdk.AccAddress(subAccountPubKey.Address())

	subAccount, err := sdk.AccAddressFromBech32(msg.SubAccount)
	if err != nil {
		return &types.MsgCreateSubAccountResponse{}, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if !derivedAddr.Equals(subAccount) {
		return &types.MsgCreateSubAccountResponse{}, errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey does not match sub_account address")
	}

	_, ok = m.GetDID(ctx, subAccount)
	if ok {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountHasDID
	}

	// check sub_account delegation, if has delegation, return error
	_, err = m.stkKeeper.GetDelegation(ctx, subAccount, sdk.ValAddress{})
	if err == nil {
		return &types.MsgCreateSubAccountResponse{}, types.ErrSubAccountHasDelegation
	}

	if !m.accountKeeper.HasAccount(ctx, subAccount) {
		newAccount := m.accountKeeper.NewAccountWithAddress(ctx, subAccount)
		err = newAccount.SetPubKey(subAccountPubKey)
		if err != nil {
			return &types.MsgCreateSubAccountResponse{}, errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
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
		),
	})
	return &types.MsgCreateSubAccountResponse{}, nil
}
