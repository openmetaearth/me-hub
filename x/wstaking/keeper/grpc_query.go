package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//regionID := strings.ToLower(types.ExperienceRegionName)
	//meid, found := k.GetMeid(ctx, req.DelegatorAddress)
	//if found {
	//	regionID = meid.RegionId
	//}
	//region, isFound := k.GetRegion(ctx, regionID)
	//if !isFound {
	//	return nil, types.ErrRegionNotExist.Wrapf("region not found=%s", regionID)
	//}
	//valAddr, valErr := sdk.ValAddressFromBech32(region.OperatorAddress)
	//if valErr != nil {
	//	return nil, valErr
	//}
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

// Delegation queries delegate info for given validator delegator pair
func (k Querier) Delegation(c context.Context, req *stakingtypes.QueryDelegationRequest) (*stakingtypes.QueryDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	regionId := k.GetRegionIdByAccount(ctx, delAddr)
	region, isFound := k.GetRegion(ctx, regionId)
	if !isFound {
		return nil, types.ErrRegionNotExist.Wrapf("region not found=%s", regionId)
	}
	valAddr, valErr := sdk.ValAddressFromBech32(region.OperatorAddress)
	if valErr != nil {
		return nil, valErr
	}
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"delegation with delegator %s not found for validator",
			req.DelegatorAddr)
	}
	delResponse, err := DelegationToDelegationResponse(ctx, k.Keeper, delegation)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &stakingtypes.QueryDelegationResponse{DelegationResponse: &delResponse}, nil
}

func DelegationToDelegationResponse(ctx sdk.Context, k *Keeper, del stakingtypes.Delegation) (stakingtypes.DelegationResponse, error) {
	if del.Unmovable.GT(sdk.ZeroInt()) {
		_, found := k.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			return stakingtypes.DelegationResponse{}, stakingtypes.ErrNoValidatorFound
		}
	}

	_, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
	if err != nil {
		return stakingtypes.DelegationResponse{}, err
	}
	amount := del.Amount.Add(del.UnMeidAmount).Add(del.Unmovable)
	return NewDelegationResp(del, sdk.NewCoin(k.BondDenom(ctx), amount)), nil
}

func (k Keeper) Stakes(goCtx context.Context, req *types.QueryStakesRequest) (*types.QueryStakesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	stakes := k.GetAllStakes(ctx)
	return &types.QueryStakesResponse{Stakes: stakes}, nil
}

func (k Keeper) QueryAllRecord(goCtx context.Context, req *types.QueryAllRecords) (*types.QueryAllRecordsResponse, error) {
	var records []types.Record
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	meidStore := prefix.NewStore(store, types.NewRecordKey)

	pageRes, err := query.Paginate(meidStore, req.Pagination, func(key []byte, value []byte) error {
		var record types.Record
		if err := k.cdc.Unmarshal(value, &record); err != nil {
			return err
		}

		records = append(records, record)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("query all records err=%s", err.Error()))
	}

	return &types.QueryAllRecordsResponse{Records: records, Pagination: pageRes}, nil
}

func (k Querier) QueryRecordByAddress(goCtx context.Context, req *types.QueryRecordsByAddress) (*types.QueryRecordsByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	from, err := sdk.AccAddressFromBech32(req.Account)
	if err != nil {
		return nil, err
	}
	records := k.GetRecordsByAddress(ctx, from)
	return &types.QueryRecordsByAddressResponse{Records: records}, nil
}

func (k Querier) QueryReviewRecordByID(goCtx context.Context, req *types.QueryReviewRecordByNumber) (*types.QueryReviewRecordByNumberResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if req.ActionNumber == "" {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("ActionNumber is empty"))
	}
	rr := k.GetReviewRecordByID(ctx, req.ActionNumber)
	return &types.QueryReviewRecordByNumberResponse{ReviewRecord: rr}, nil
}

// Delegation queries delegate info for given validator delegator pair
func (k Querier) AllDelegations(c context.Context, req *types.QueryAllDelegationsRequest) (*types.QueryAllDelegationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	queryStore := prefix.NewStore(store, stakingtypes.DelegationKey)

	delegations := []stakingtypes.Delegation{}
	pageRes, err := query.Paginate(queryStore, req.Pagination, func(key []byte, value []byte) error {
		delegation := types.MustUnmarshalDelegation(k.cdc, value)
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllDelegationsResponse{Delegations: delegations, Pagination: pageRes}, nil
}

// RegionWithdrawer returns the address that is granted withdraw
// for the given region, or an empty address if not set.
func (k Querier) RegionWithdrawer(goCtx context.Context, req *types.QueryRegionWithdrawerRequest) (*types.QueryRegionWithdrawerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.RegionId == "" {
		return nil, status.Error(codes.InvalidArgument, "region_id cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, _ := k.GetRegionWithdraw(ctx, req.RegionId)
	return &types.QueryRegionWithdrawerResponse{
		RegionId: req.RegionId,
		Address:  address,
	}, nil
}
