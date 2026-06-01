package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

var (
	fixedDepositMinRate = sdk.MustNewDecFromStr("0.0001")
	fixedDepositMaxRate = sdk.MustNewDecFromStr("10000")
)

func validateFixedDepositCfgRate(rate sdk.Dec) error {
	if !rate.IsPositive() {
		return types.ErrFixedDepositConfigRateInvalid.Wrapf("rate must be > 0 (%s)", rate.String())
	}

	if rate.LT(fixedDepositMinRate) || rate.GT(fixedDepositMaxRate) {
		return types.ErrFixedDepositConfigRateInvalid.Wrapf("rate(%s) out of range [%s, %s]", rate.String(), fixedDepositMinRate.String(), fixedDepositMaxRate.String())
	}
	return nil
}

func (k MsgServer) NewFixedDepositCfg(goCtx context.Context, msg *types.MsgNewFixedDepositCfg) (*types.MsgNewFixedDepositCfgResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Dao) {
		return nil, types.ErrCheckGlobalDao
	}

	_, found := k.GetRegion(ctx, msg.RegionId)
	if !found {
		return nil, types.ErrRegionName.Wrapf("add fixed deposit config error, region not exist (%s)", msg.RegionId)
	}

	if msg.Term <= 0 {
		return nil, types.ErrAddFixedDepositConfig.Wrapf("add fixed deposit config error, term is not positive 0 (%d)", msg.Term)
	}

	if err := validateFixedDepositCfgRate(msg.Rate); err != nil {
		return nil, sdkerrors.Wrap(err, "add fixed deposit config error")
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

	if err := validateFixedDepositCfgRate(msg.Rate); err != nil {
		return nil, sdkerrors.Wrap(err, "set fixed deposit config rate error")
	}

	config.Rate = msg.Rate
	k.Keeper.SetFixedDepositCfg(ctx, config)

	return &types.MsgSetFixedDepositCfgRateResp{}, nil
}
