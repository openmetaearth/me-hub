package keeper

import (
	"context"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

type Querier struct {
	*Keeper
}

var _ types.QueryServer = Querier{}

func (k Keeper) Region(goCtx context.Context, req *types.QueryRegionRequest) (*types.QueryRegionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	region, found := k.GetRegion(ctx, req.RegionId)
	if !found {
		return nil, types.ErrRegionNotExist
	}
	return &types.QueryRegionResponse{Region: region}, nil
}

func (k Keeper) AllRegion(goCtx context.Context, req *types.QueryAllRegionRequest) (*types.QueryAllRegionResponse, error) {
	var regions []types.Region

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	regionStore := prefix.NewStore(store, types.KeyPrefix(types.RegionKeyPrefix))

	pageRes, err := query.Paginate(regionStore, req.Pagination, func(key []byte, value []byte) error {
		var region types.Region
		if err := k.cdc.Unmarshal(value, &region); err != nil {
			return err
		}
		regions = append(regions, region)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRegionResponse{Region: regions, Pagination: pageRes}, nil
}
func (k Keeper) DelegationRewards(c context.Context, req *types.QueryDelegationRewardsRequest) (*types.QueryDelegationRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	del := k.Delegation(ctx, delAdr, sdk.ValAddress{})
	if del == nil {
		return nil, status.Error(codes.NotFound, "delegator not found, address="+delAdr.String())
	}
	delegation, ok := del.(stakingtypes.Delegation)
	if !ok {
		return nil, types.ErrAssertDelegation
	}
	interest, err := k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
	if err != nil {
		return nil, err
	}
	//endingPeriod := k.IncrementValidatorPeriod(ctx, val)
	//rewards := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)
	rewards := sdk.NewDecCoins(sdk.NewDecCoinFromDec(params.BaseDenom, interest))
	return &types.QueryDelegationRewardsResponse{Rewards: rewards}, nil
	//return &types.QueryDelegationRewardsResponse{Rewards: sdk.NewDecCoinsFromCoins(sdk.NewCoin(sdk.BaseMEDenom, interest.TruncateInt()))}, nil

}
