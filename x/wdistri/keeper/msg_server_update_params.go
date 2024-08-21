package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/wdistri/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(types.ErrPermissionDenied, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}
	if msg.Params == nil {
		return nil, types.ErrInvalidParams
	}
	if (!msg.Params.DistributionParams.BaseProposerReward.IsNil() && !msg.Params.DistributionParams.BaseProposerReward.IsZero()) || //nolint:staticcheck
		(!msg.Params.DistributionParams.BonusProposerReward.IsNil() && !msg.Params.DistributionParams.BonusProposerReward.IsZero()) { //nolint:staticcheck
		return nil, errors.Wrapf(types.ErrInvalidParams, "cannot update base or bonus proposer reward because these are deprecated fields")
	}

	k.WrapDistrKeeper.SetParams(ctx, msg.Params.DistributionParams)
	return &types.MsgUpdateParamsResponse{}, nil
}
