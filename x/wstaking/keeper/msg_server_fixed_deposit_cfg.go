package keeper

import (
	"context"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k MsgServer) NewFixedDepositCfg(goCtx context.Context, msg *types.MsgNewFixedDepositCfg) (*types.MsgNewFixedDepositCfgResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Dao) {
		return nil, types.ErrCheckGlobalDao
	}

	_, found := k.GetRegionCache(msg.RegionId)
	if !found {
		return nil, types.ErrRegionName.Wrapf("add fixed deposit config error, region not exist (%s)", msg.RegionId)
	}

	if msg.Term <= 0 {
		return nil, types.ErrAddFixedDepositConfig.Wrapf("add fixed deposit config error, term is not positive 0 (%d)", msg.Term)
	}

	if !msg.Rate.IsPositive() {
		return nil, types.ErrAddFixedDepositConfig.Wrapf("add fixed deposit config error, rate is not positive 0 (%s)", msg.Rate.String())
	}

	minRate := sdkmath.LegacyMustNewDecFromStr("0.0001")
	maxRate := sdkmath.LegacyMustNewDecFromStr("10000")
	if msg.Rate.LT(minRate) || msg.Rate.GT(maxRate) {
		return nil, types.ErrAddFixedDepositConfig.Wrapf("add fixed deposit config rate(%s) error (%s)",
			msg.Rate.String(), types.ErrFixedDepositConfigRateInvalid)
	}

	_, ok := k.GetFixedDepositCfg(ctx, msg.RegionId, msg.Term)
	if ok {
		return nil, types.ErrAddFixedDepositConfig.Wrapf("add fixed deposit config error (%s)", types.ErrFixedDepositConfigAlreadyExists)
	}

	cfg := types.FixedDepositCfg{
		RegionId: msg.RegionId,
		Term:     msg.Term,
		Rate:     msg.Rate,
		Status:   types.RegionFixedDepositCfgStatusActive,
	}
	k.Keeper.SetFixedDepositCfg(ctx, cfg)
	k.InitFixedDepositCountOfCfg(ctx, msg.RegionId, msg.Term)

	return &types.MsgNewFixedDepositCfgResp{}, nil
}

func (k MsgServer) RemoveFixedDepositCfg(goCtx context.Context, msg *types.MsgRemoveFixedDepositCfg) (*types.MsgRemoveFixedDepositCfgResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Admin) {
		return nil, types.ErrCheckGlobalDao
	}

	_, found := k.GetFixedDepositCfg(ctx, msg.RegionId, msg.Term)
	if !found {
		return nil, types.ErrRemoveFixedDepositConfig.Wrapf("fixed deposit config not found  for region(%s) and term(%d)", msg.RegionId, msg.Term)
	}

	count := k.GetFixedDepositCountOfCfg(ctx, msg.RegionId, msg.Term)
	if count > 0 {
		return nil, types.ErrRemoveFixedDepositConfig.Wrapf("remove fixed deposit config error:(%s)", types.ErrFixedDepositExistUnderConfig)
	}

	k.Keeper.RemoveFixedDepositCfg(ctx, msg.RegionId, msg.Term)

	return &types.MsgRemoveFixedDepositCfgResp{}, nil
}

func (k MsgServer) SetFixedDepositCfgStatus(goCtx context.Context, msg *types.MsgSetFixedDepositCfgStatus) (*types.MsgSetFixedDepositCfgStatusResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Admin) {
		return nil, types.ErrCheckGlobalDao
	}

	config, ok := k.GetFixedDepositCfg(ctx, msg.RegionId, msg.Term)
	if !ok {
		return nil, types.ErrSetFixedDepositConfigStatus.Wrapf("set fixed deposit config status error (%s)", types.ErrNoFixedDepositCountOfCfgFound)
	}
	config.Status = msg.Status
	k.Keeper.SetFixedDepositCfg(ctx, config)

	return &types.MsgSetFixedDepositCfgStatusResp{}, nil
}

func (k MsgServer) SetFixedDepositCfgRate(goCtx context.Context, msg *types.MsgSetFixedDepositCfgRate) (*types.MsgSetFixedDepositCfgRateResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Admin) {
		return nil, types.ErrCheckGlobalDao
	}

	config, ok := k.GetFixedDepositCfg(ctx, msg.RegionId, msg.Term)
	if !ok {
		return nil, types.ErrSetFixedDepositConfigRate.Wrapf("set fixed deposit config rate error (%s)", types.ErrNoFixedDepositCountOfCfgFound)
	}

	config.Rate = msg.Rate
	k.Keeper.SetFixedDepositCfg(ctx, config)

	return &types.MsgSetFixedDepositCfgRateResp{}, nil
}
