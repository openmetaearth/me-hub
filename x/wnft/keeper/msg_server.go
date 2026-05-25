package keeper

import (
	"context"
	"strconv"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/openmetaearth/me-hub/utils"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/wnft/types"
)

// MsgServer is wrapper staking customParamsKeeper message server.
type MsgServer struct {
	*Keeper
	nft.MsgServer
}

var _ types.MsgServer = MsgServer{}

// NewMsgServerImpl returns an implementation of the staking wrapped MsgServer.
func NewMsgServerImpl(
	keeper *Keeper,
	nftMsgSrv nft.MsgServer,
) MsgServer {
	return MsgServer{
		Keeper:    keeper,
		MsgServer: nftMsgSrv,
	}
}

// NewClass implements the NewClass method of types.MsgServer.
func (k Keeper) NewClass(goCtx context.Context, msg *types.MsgNewClass) (*types.MsgNewClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, ok := k.GetClass(ctx, msg.ClassId)
	if ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "class %s already exists", msg.ClassId)
	}

	//Check if the name occupies the zone name todo
	if utils.CheckIsRegionName(msg.ClassId) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid class name %s", msg.ClassId)
	}

	classMetadata := &types.ClassMetadata{
		Creator: msg.Sender,
	}

	metadata, err := codectypes.NewAnyWithValue(classMetadata)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrLogic, "%v", err)
	}

	class := nft.Class{
		Id:          msg.ClassId,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		TotalSupply: msg.TotalSupply,
		Data:        metadata,
	}

	err = k.SaveClass(ctx, class)
	if err != nil {
		return &types.MsgNewClassResponse{}, err
	}
	ctx.EventManager().EmitTypedEvent(&class)
	return &types.MsgNewClassResponse{}, nil
}

func (k Keeper) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//check token id An integer between 1 and the total supply of the NFT type, non-repeating
	if !k.HasClass(ctx, msg.ClassId) {
		return nil, sdkerrors.Wrap(nft.ErrClassNotExists, msg.ClassId)
	}

	class, ok := k.GetClass(ctx, msg.ClassId)
	if !ok {
		return nil, sdkerrors.Wrap(nft.ErrClassNotExists, msg.ClassId)
	}

	var classMetadata types.ClassMetadata
	if err := k.cdc.Unmarshal(class.Data.GetValue(), &classMetadata); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrLogic, "%v", err)
	}

	if classMetadata.Creator != msg.Creator {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not the creator of class %s", msg.Creator, msg.ClassId)
	}

	tokenId, err := strconv.ParseUint(msg.TokenId, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token id")
	}

	if tokenId < 1 || tokenId > class.TotalSupply {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token id")
	}
	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Receiver)
	}

	if err = k.Mint(ctx,
		nft.NFT{
			ClassId: msg.ClassId,
			Id:      msg.TokenId,
			Uri:     msg.Uri,
			UriHash: msg.UriHash,
			Data:    nil,
		},
		receiver,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Receiver),
			sdk.NewAttribute(types.AttributeKeyClassName, class.Name),
		),
	})

	return &types.MsgMintNFTResponse{}, nil
}

// Send implements Send method of the types.MsgServer.
func (k MsgServer) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if msg.ClassId == kyctypes.ModuleName {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "SBT is not allowed to be transferred to others")
	}

	owner := k.GetOwner(ctx, msg.ClassId, msg.Id)
	if !owner.Equals(sender) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not the owner of nft %s", sender, msg.Id)
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	class, found := k.GetClass(ctx, msg.ClassId)
	if !found {
		return nil, sdkerrors.Wrap(nft.ErrClassNotExists, msg.ClassId)
	}

	myNFT, found := k.GetNFT(ctx, msg.ClassId, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(nft.ErrNFTNotExists, msg.Id)
	}

	if err := k.Transfer(ctx, msg.ClassId, msg.Id, receiver); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSendNFT,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.Id),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
			sdk.NewAttribute(types.AttributeKeyUri, myNFT.Uri),
			sdk.NewAttribute(types.AttributeKeyClassName, class.Name),
		),
	})
	return &types.MsgSendResponse{}, nil
}
