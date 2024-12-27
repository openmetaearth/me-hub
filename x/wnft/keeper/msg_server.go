package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/st-chain/me-hub/utils"
	"github.com/st-chain/me-hub/x/wnft/types"
	"strconv"
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

	class := nft.Class{
		Id:          msg.ClassId,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		TotalSupply: msg.TotalSupply,
	}

	k.SaveClass(ctx, class)
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
	tokenId, err := strconv.ParseUint(msg.TokenId, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token id")
	}

	if tokenId < 1 || tokenId > class.TotalSupply {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token id")
	}
	receiver, err := sdk.AccAddressFromBech32(msg.Sender)

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

	ctx.EventManager().EmitTypedEvent(&nft.EventMint{
		ClassId: msg.ClassId,
		Id:      msg.TokenId,
		Owner:   msg.Sender,
	})

	return &types.MsgMintNFTResponse{}, nil
}
