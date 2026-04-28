package keeper

import (
	"unicode"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"golang.org/x/net/context"
)

func (k MsgServer) NewRecord(goCtx context.Context, msg *types.MsgNewRecord) (*types.MsgNewRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	if msg.ActionNumber == "" {
		return nil, types.ErrInvalidRecordParams.Wrap("invalid record number,is empty")
	} else {
		for _, r := range msg.ActionNumber {
			// unicode.IsLetter check if it is a letter
			// unicode.IsDigit check if it is a number
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return nil, types.ErrInvalidRecordParams.Wrap("invalid record number,only letters and numbers are allowed")
			}
		}
	}
	if msg.ActionUrl == "" {
		return nil, types.ErrInvalidRecordParams.Wrap("url is empty")
	}
	r := types.Record{
		RecordNumber: msg.ActionNumber,
		Url:          msg.ActionUrl,
		From:         msg.From,
	}
	k.SetRecord(ctx, r, from)
	return &types.MsgNewRecordResponse{}, nil
}

func (k MsgServer) ReviewRecord(goCtx context.Context, msg *types.MsgReviewRecord) (*types.MsgReviewRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	globalAdmin := k.daoKeeper.GetGlobalDao(ctx)
	meidAdmin := k.daoKeeper.GetMeidDao(ctx)
	if globalAdmin != msg.From && meidAdmin != msg.From {
		// use a constant format string in Wrapf to avoid non-constant format string vet issue
		return nil, sdkerrors.Wrapf(types.ErrParameter, "review record account (%s) should  be global admin", msg.From)
	}
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	if msg.ReviewResult == "" {
		return nil, types.ErrInvalidRecordParams.Wrap("review result is empty")
	}
	if msg.RecordHash == "" {
		return nil, types.ErrInvalidRecordParams.Wrap("invalid record hash,is empty")
	}
	if msg.ActionNumber == "" {
		return nil, types.ErrInvalidRecordParams.Wrap("invalid record number,is empty")
	}
	rr := types.ReviewRecord{
		RecordHash:      msg.RecordHash,
		ActionNumber:    msg.ActionNumber,
		RecordResult:    msg.ReviewResult,
		ReviewedAddress: msg.ReviewedAddress,
	}
	k.SetReviewRecord(ctx, rr)
	return &types.MsgReviewRecordResponse{}, nil
}
