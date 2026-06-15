package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/x/wstaking/types"

	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

func (k MsgServer) IbcTransferFromRegionTreasure(goCtx context.Context, msg *types.MsgIbcTransferFromRegionTreasure) (*types.MsgIbcTransferFromRegionTreasureResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dao := k.daoKeeper.GetGlobalDao(ctx)
	if msg.Creator != dao {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "sender is not the dao")
	}
	region, found := k.GetRegion(ctx, msg.RegionId)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "region not found")
	}

	treasureAddress := region.RegionTreasureAddr

	_, err := k.IbcTransferKeeper.Transfer(ctx, newRegionTreasureIBCTransferMsg(
		msg,
		treasureAddress,
	))
	if err != nil {
		return nil, err
	}

	return &types.MsgIbcTransferFromRegionTreasureResponse{}, nil
}

func newRegionTreasureIBCTransferMsg(msg *types.MsgIbcTransferFromRegionTreasure, treasureAddress string) *ibctransfertypes.MsgTransfer {
	return ibctransfertypes.NewMsgTransfer(
		msg.SourcePort,
		msg.SourceChannel,
		msg.Token,
		treasureAddress,
		msg.Receiver,
		ibcclienttypes.Height{RevisionNumber: msg.TimeoutHeight.RevisionNumber, RevisionHeight: msg.TimeoutHeight.RevisionHeight},
		msg.TimeoutTimestamp,
		msg.Memo,
	)
}
