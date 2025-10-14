package blacklist

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/st-chain/me-hub/x/blacklist/keeper"
	"github.com/st-chain/me-hub/x/blacklist/types"
)

// NewHandler returns a handler for "blacklist" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		sdkCtx := sdk.WrapSDKContext(ctx)

		switch msg := msg.(type) {
		case *types.MsgUpdateBlacklist:
			res, err := msgServer.UpdateBlacklist(sdkCtx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
