package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strconv"
)

func (k MsgServer) NewFixedDepositCfg(goCtx context.Context, msg *types.MsgNewFixedDepositCfg) (*types.MsgNewFixedDepositCfgResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//globalAdminAddress := k.GetGlobalAdminAddress(ctx)
	//if msg.Admin != globalAdminAddress {
	//	return nil, types.ErrAddFixedDepositConfig.Wrapf("only global admin can add fixed deposit config")
	//}

	_, found := k.GetRegion(ctx, msg.RegionId)
	if !found {
		return nil, types.ErrAddFixedDepositConfig.Wrapf("add fixed deposit config error, region not exist (%s)", msg.RegionId)
	}

	if msg.Term <= 0 {
		return nil, types.ErrAddFixedDepositConfig.Wrapf("add fixed deposit config error, term is not positive 0 (%d)", msg.Term)
	}

	if !msg.Rate.IsPositive() {
		return nil, types.ErrAddFixedDepositConfig.Wrapf("add fixed deposit config error, rate is not positive 0 (%s)", msg.Rate.String())
	}

	minRate := sdk.MustNewDecFromStr("0.0001")
	maxRate := sdk.MustNewDecFromStr("10000")
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

	event := sdk.NewEvent(types.EventTypeAddFixedDepositCfg,
		sdk.NewAttribute(types.AttributeKeyAccount, msg.Admin),
		sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
		sdk.NewAttribute(types.AttributeKeyTerm, strconv.FormatInt(msg.Term, 10)),
		sdk.NewAttribute(types.AttributeKeyRate, msg.Rate.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgNewFixedDepositCfgResp{}, nil
}

func (k MsgServer) RemoveFixedDepositCfg(goCtx context.Context, msg *types.MsgRemoveFixedDepositCfg) (*types.MsgRemoveFixedDepositCfgResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//globalAdminAddress := k.GetGlobalAdminAddress(ctx)
	//if msg.Admin != globalAdminAddress {
	//	return nil, types.ErrRemoveFixedDepositConfig.Wrapf("only global admin can remove fixed deposit config")
	//}

	count := k.GetFixedDepositCountOfCfg(ctx, msg.RegionId, msg.Term)
	if count != 0 {
		return nil, types.ErrRemoveFixedDepositConfig.Wrapf("remove fixed deposit config error:(%s)", types.ErrFixedDepositExistUnderConfig)
	}

	k.Keeper.RemoveFixedDepositCfg(ctx, msg.RegionId, msg.Term)

	event := sdk.NewEvent(types.EventTypeRemoveFixedDepositCfg,
		sdk.NewAttribute(types.AttributeKeyAccount, msg.Admin),
		sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
		sdk.NewAttribute(types.AttributeKeyTerm, strconv.FormatInt(msg.Term, 10)),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgRemoveFixedDepositCfgResp{}, nil
}

func (k MsgServer) SetFixedDepositCfgStatus(goCtx context.Context, msg *types.MsgSetFixedDepositCfgStatus) (*types.MsgSetFixedDepositCfgStatusResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//globalAdminAddress := k.GetGlobalAdminAddress(ctx)
	//if msg.Admin != globalAdminAddress {
	//	return nil, types.ErrSetFixedDepositConfigStatus.Wrapf("only global admin can set fixed deposit config status")
	//}

	config, ok := k.GetFixedDepositCfg(ctx, msg.RegionId, msg.Term)
	if !ok {
		return nil, types.ErrSetFixedDepositConfigStatus.Wrapf("set fixed deposit config status error (%s)", types.ErrNoFixedDepositCountOfCfgFound)
	}
	config.Status = msg.Status
	k.Keeper.SetFixedDepositCfg(ctx, config)

	event := sdk.NewEvent(types.EventTypeSetFixedDepositCfgStatus,
		sdk.NewAttribute(types.AttributeKeyAccount, msg.Admin),
		sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
		sdk.NewAttribute(types.AttributeKeyTerm, strconv.FormatInt(msg.Term, 10)),
		sdk.NewAttribute(types.AttributeKeyStatus, msg.Status.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgSetFixedDepositCfgStatusResp{}, nil
}

func (k MsgServer) SetFixedDepositCfgRate(goCtx context.Context, msg *types.MsgSetFixedDepositCfgRate) (*types.MsgSetFixedDepositCfgRateResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//globalAdminAddress := k.GetGlobalAdminAddress(ctx)
	//if msg.Admin != globalAdminAddress {
	//	return nil, types.ErrSetFixedDepositConfigRate.Wrapf("only global admin can set fixed deposit config rate")
	//}

	config, ok := k.GetFixedDepositCfg(ctx, msg.RegionId, msg.Term)
	if !ok {
		return nil, types.ErrSetFixedDepositConfigRate.Wrapf("set fixed deposit config rate error (%s)", types.ErrNoFixedDepositCountOfCfgFound)
	}

	config.Rate = msg.Rate
	k.Keeper.SetFixedDepositCfg(ctx, config)
	event := sdk.NewEvent(types.EventTypeSetFixedDepositCfgRate,
		sdk.NewAttribute(types.AttributeKeyAccount, msg.Admin),
		sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
		sdk.NewAttribute(types.AttributeKeyTerm, strconv.FormatInt(msg.Term, 10)),
		sdk.NewAttribute(types.AttributeKeyRate, msg.Rate.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgSetFixedDepositCfgRateResp{}, nil
}
