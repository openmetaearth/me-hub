package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/openmetaearth/me-hub/app/params"
	minttypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const DayPerYear uint64 = 365

var minDepositAmount = sdk.NewInt(1_000_000) // 0.01mec == 1000000umec

func (k MsgServer) TermToDuration(term int64) (time.Duration, error) {
	// for formal environment
	minutesPerDay := (24 * 60 * time.Minute)
	// for test environment
	//minutesPerDay := time.Minute

	return time.Duration(int64(minutesPerDay) * term), nil
}

func (k MsgServer) GetFixedDepositInterest(cfg *types.FixedDepositCfg, principal sdk.Coin, term int64) (sdk.Coin, error) {
	principalNormed, err := sdk.ParseCoinsNormalized(principal.String())
	if err != nil {
		return sdk.Coin{}, types.ErrPayInterest.Wrap(err.Error())
	}
	principalAmount := principalNormed.AmountOf(params.BaseDenom)
	interest := cfg.Rate.MulInt(principalAmount).MulInt(math.NewInt(term)).QuoInt(sdk.NewIntFromUint64(DayPerYear))

	return sdk.NewCoin(params.BaseDenom, interest.TruncateInt()), nil
}

func (k MsgServer) DoFixedDeposit(goCtx context.Context, msg *types.MsgDoFixedDeposit) (*types.MsgDoFixedDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	regionId, err := k.MustGetKycRegionIdByAccount(ctx, msg.Account)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrDidNotExists, err.Error())
	}

	if !msg.Principal.Amount.IsPositive() {
		return nil, types.ErrDoFixedDeposit.Wrapf("fixed deposit amount error (%s)", types.ErrAmountNotPositive)
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Principal.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Principal.Denom, bondDenom,
		)
	}

	// minimum amount 0.01mec == 1000000umec
	if msg.Principal.Amount.LT(minDepositAmount) {
		return nil, types.ErrDoFixedDeposit.Wrapf("fixed deposit amount error (%s)", types.ErrAmountLessThanMin)
	}

	duration, err := k.TermToDuration(msg.Term)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, "period value error (%s)", err)
	}
	startTime := ctx.BlockTime()
	endTime := startTime.Add(duration)

	config, ok := k.GetFixedDepositCfg(ctx, regionId, msg.Term)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, "get fixed deposit config error (%s)", types.ErrNoFixedDepositCountOfCfgFound)
	}
	if config.Status != types.RegionFixedDepositCfgStatusActive {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit,
			"fixed deposit config status is inactive, fixed deposit not allowed (%s)", types.ErrFixedDepositConfigInactive)
	}

	interest, err := k.GetFixedDepositInterest(&config, msg.Principal, msg.Term)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, "get calculate interests error (%s)", err)
	}

	accAddr, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, "account format error (%s)", err)
	}

	//principal from user account to principal vault, interest from base vault to interest vault
	region, found := k.GetRegion(ctx, regionId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, "region id(%s) no exist", regionId)
	}
	regionBaseAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, "region base account(%s) format error (%s)", region.RegionTreasureAddr, err)
	}
	regionInterestAddr, err := sdk.AccAddressFromBech32(region.DepositInterestAddr)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, "region interest account(%s) format error (%s)", region.DepositInterestAddr, err)
	}

	principalAddr := k.authKeeper.GetModuleAddress(types.FixedDepositPrincipalPool)
	if principalAddr == nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, fmt.Sprintf("%s module account has not been set", types.FixedDepositPrincipalPool))
	}

	if coin := k.bankKeeper.GetBalance(ctx, accAddr, msg.Principal.Denom); coin.IsLT(msg.Principal) {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit, "account %s balance(%s) less deposit coin(%s)", accAddr.String(), coin.String(), msg.Principal.String())
	}

	totalRewardsPerBlockTemp := k.mintKeeper.GetPerBlockMintCoinAmount(ctx)
	totalRewardsPerBlock := sdk.NewIntFromBigInt(&totalRewardsPerBlockTemp)
	totalSupply := sdk.NewInt(types.CaclTotalSupply).MulRaw(100000000)
	initAllocationFunds := sdk.NewInt(minttypes.TotalMintCoinsAmount)

	deAmount := region.DelegateAmount

	interestAmountDec := sdk.NewDecFromInt(deAmount).Mul(sdk.NewDecFromInt(region.RegionShare)).Mul(sdk.NewDecFromInt(totalRewardsPerBlock)).
		Quo(sdk.NewDecFromInt(totalSupply).Mul(sdk.NewDecFromInt(initAllocationFunds)))

	remainingBalance := sdk.NewDecFromInt(k.bankKeeper.GetBalance(ctx, regionBaseAddr, interest.Denom).Amount).Sub(interestAmountDec.Add(region.DelegateInterest))
	if remainingBalance.Sub(sdk.NewDecFromInt(interest.Amount)).LT(sdk.ZeroDec()) {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit,
			"region account base address %s balance(%s) less interest coin(%s)",
			regionBaseAddr.String(),
			k.bankKeeper.GetBalance(ctx, regionBaseAddr, interest.Denom).String(),
			interest.String())
	}

	//1. send principal from user account principal module account
	err = k.bankKeeper.Extend().SendCoinsFromAccountToModuleWithTag(
		ctx,
		accAddr,
		types.FixedDepositPrincipalPool,
		sdk.NewCoins(msg.Principal),
		fmt.Sprintf("SendFixedPrincipal_%d", msg.Term),
	)
	if err != nil {
		return nil, types.ErrDoFixedDeposit.Wrapf("send coin from region base account(%s) to principal module account(%s) error (%s)",
			regionBaseAddr.String(), types.FixedDepositPrincipalPool, err)
	}

	//2. send interest from region base account to region interest account
	err = k.bankKeeper.Extend().SendCoinsWithTag(
		ctx,
		regionBaseAddr,
		regionInterestAddr,
		sdk.NewCoins(interest),
		fmt.Sprintf("SendFixedInterest_%d", msg.Term),
	)
	if err != nil {
		return nil, types.ErrDoFixedDeposit.Wrapf("send coin from account(%s) to interest account(%s) error (%s)",
			accAddr.String(), regionInterestAddr.String(), err)
	}

	fixedDeposit := types.FixedDeposit{
		Account:   msg.Account,
		Principal: msg.Principal,
		Interest:  interest,
		StartTime: startTime,
		EndTime:   endTime,
		Term:      msg.Term,
		Rate:      config.Rate,
	}
	id := k.AppendFixedDeposit(ctx, fixedDeposit)

	err = k.Keeper.IncreaseFixedDepositCountOfCfg(ctx, region.RegionId, msg.Term)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit,
			"increase fixed deposit count under the current config error, region(%s) term (%s) error(%s)",
			region.RegionId, msg.Term, err)
	}

	region.FixedDepositAmount = region.FixedDepositAmount.Add(fixedDeposit.Principal.Amount)
	k.SetRegion(ctx, region)

	amount, found := k.GetFixedDepositTotalAmount(ctx)
	if !found {
		k.SetFixedDepositTotalAmount(ctx, types.FixedDepositTotal{Amount: fixedDeposit.Principal})
	} else {
		k.SetFixedDepositTotalAmount(ctx, types.FixedDepositTotal{Amount: amount.Amount.Add(fixedDeposit.Principal)})
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeDoFixedDeposit,
		sdk.NewAttribute(types.AttributeKeyId, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyRegionId, regionId),
		sdk.NewAttribute(types.AttributeKeyAccount, msg.Account),
		sdk.NewAttribute(types.AttributeKeyPrincipalAddr, principalAddr.String()),
		sdk.NewAttribute(types.AttributeKeyPrincipal, msg.Principal.String()),
		sdk.NewAttribute(types.AttributeKeyTreasureAddr, regionBaseAddr.String()),
		sdk.NewAttribute(types.AttributeKeyInterestAddr, regionInterestAddr.String()),
		sdk.NewAttribute(types.AttributeKeyInterest, fixedDeposit.Interest.String()),
		sdk.NewAttribute(types.AttributeKeyStartTime, fixedDeposit.StartTime.String()),
		sdk.NewAttribute(types.AttributeKeyEndTime, fixedDeposit.EndTime.String()),
		sdk.NewAttribute(types.AttributeKeyRate, fixedDeposit.Rate.String()),
		sdk.NewAttribute(types.AttributeKeyTerm, strconv.FormatInt(fixedDeposit.Term, 10)),
	))
	return &types.MsgDoFixedDepositResponse{Id: id}, nil
}

func (k MsgServer) WithdrawFixedDeposit(goCtx context.Context, msg *types.MsgWithdrawFixedDeposit) (*types.MsgWithdrawFixedDepositResponse, error) {
	var interest sdk.Coin
	ctx := sdk.UnwrapSDKContext(goCtx)
	log := ctx.Logger()

	regionId, err := k.MustGetKycRegionIdByAccount(ctx, msg.Account)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedWithDraw, err.Error())
	}

	fixedDeposit, isFound := k.GetFixedDeposit(ctx, msg.Id)
	if !isFound {
		return nil, types.ErrDoFixedWithDraw.Wrapf("fixed deposit not found (%s)", types.ErrNoFixedDepositFound)
	}

	if fixedDeposit.Account != msg.Account {
		return nil, types.ErrDoFixedWithDraw.Wrapf("only depositor can withdraw, depositor:(%s), current retreiver:(%s), error:(%s)",
			fixedDeposit.Account, msg.Account, types.ErrFixedDepositInvalidPayee)
	}

	accAddr, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedWithDraw, "account format error (%s)", err)
	}

	//expired: principal from principal vault to user account addr; interest from interest vault to user account
	//no expired: principal from principal vault to user account addr; interest from interest vault to base vault
	region, found := k.GetRegion(ctx, regionId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedWithDraw, "region id(%s) no exist", regionId)
	}
	regionInterestAddr, err := sdk.AccAddressFromBech32(region.DepositInterestAddr)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedWithDraw, "region interest account(%s) format error (%s)", region.DepositInterestAddr, err)
	}
	principalAddr := k.authKeeper.GetModuleAddress(types.FixedDepositPrincipalPool)
	if principalAddr == nil {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedWithDraw, fmt.Sprintf("%s module account has not been set", types.FixedDepositPrincipalPool))
	}

	if coin := k.bankKeeper.GetBalance(ctx, principalAddr, fixedDeposit.Principal.Denom); coin.IsLT(fixedDeposit.Principal) {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedWithDraw,
			"principal account %s balance(%s) less principal coin(%s)",
			principalAddr.String(),
			coin.String(),
			fixedDeposit.Principal.String())
	}
	if coin := k.bankKeeper.GetBalance(ctx, regionInterestAddr, fixedDeposit.Interest.Denom); coin.IsLT(fixedDeposit.Interest) {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedDeposit,
			"region interest account %s balance(%s) less interest coin(%s)",
			regionInterestAddr.String(),
			coin.String(),
			fixedDeposit.Interest.String())
	}

	expired := fixedDeposit.EndTime.Unix() <= ctx.BlockTime().Unix()
	if expired {
		//1. deposit period has expired, send the principal from principal module account to user account
		err = k.bankKeeper.Extend().SendCoinsFromModuleToAccountWithTag(ctx,
			types.FixedDepositPrincipalPool,
			accAddr,
			sdk.NewCoins(fixedDeposit.Principal),
			fmt.Sprintf("WithdrawFixedPrincipal_%d", fixedDeposit.Term),
		)
		if err != nil {
			return nil, types.ErrDoFixedWithDraw.Wrapf("send coin from principal vault to account error (%s)", err)
		}

		//2. deposit period has expired, send the interest from interest account to user account
		err = k.bankKeeper.Extend().SendCoinsWithTag(ctx,
			regionInterestAddr,
			accAddr,
			sdk.NewCoins(fixedDeposit.Interest),
			fmt.Sprintf("WithdrawFixedInterest_%d", fixedDeposit.Term),
		)
		if err != nil {
			return nil, types.ErrDoFixedWithDraw.Wrapf("send coin from interest vault to account error (%s)", err)
		}

		interest = fixedDeposit.Interest

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeDoFixedWithDraw,
				sdk.NewAttribute(types.AttributeKeyId, fmt.Sprintf("%d", fixedDeposit.Id)),
				sdk.NewAttribute(types.AttributeKeyExpired, "expired"),
				sdk.NewAttribute(types.AttributeKeyRegionId, regionId),
				sdk.NewAttribute(types.AttributeKeyAccount, msg.Account),
				sdk.NewAttribute(types.AttributeKeyPrincipalAddr, principalAddr.String()),
				sdk.NewAttribute(types.AttributeKeyPrincipal, fixedDeposit.Principal.String()),
				sdk.NewAttribute(types.AttributeKeyInterestAddr, regionInterestAddr.String()),
				sdk.NewAttribute(types.AttributeKeyInterest, fixedDeposit.Interest.String()),
				sdk.NewAttribute(types.AttributeKeyStartTime, fixedDeposit.StartTime.String()),
				sdk.NewAttribute(types.AttributeKeyEndTime, fixedDeposit.EndTime.String()),
				sdk.NewAttribute(types.AttributeKeyRate, fixedDeposit.Rate.String()),
				sdk.NewAttribute(types.AttributeKeyTerm, strconv.FormatInt(fixedDeposit.Term, 10)),
			),
		)

		err = k.Keeper.DecreaseFixedDepositCountOfCfg(ctx, region.RegionId, fixedDeposit.Term)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrDoFixedWithDraw,
				"decrease fixed deposit count under the current config error, region(%s) term (%s) error(%s)",
				region.RegionId, fixedDeposit.Term, err)
		}
	} else {
		return nil, types.ErrDoFixedWithDraw.Wrapf("withdraw fixed deposit error (%s)", types.ErrFixedDepositNotExpired)
	}

	k.RemoveFixedDeposit(ctx, msg.Id)

	amount, found := k.GetFixedDepositTotalAmount(ctx)
	if found {
		k.SetFixedDepositTotalAmount(ctx, types.FixedDepositTotal{Amount: amount.Amount.Sub(fixedDeposit.Principal)})
	}

	region.FixedDepositAmount = region.FixedDepositAmount.Sub(fixedDeposit.Principal.Amount)
	if region.FixedDepositAmount.IsNegative() {
		return nil, sdkerrors.Wrapf(types.ErrDoFixedWithDraw, "region fixed deposit amount(%s) less than zero", region.FixedDepositAmount.String())
	}
	k.SetRegion(ctx, region)

	log.Info("withdraw fixed deposit",
		"period expired:", expired,
		", account:", fixedDeposit.Account,
		", principal:", fixedDeposit.Principal,
		", interest:", interest,
		", start time:", fixedDeposit.StartTime,
		", end time:", fixedDeposit.EndTime,
		", term", strconv.FormatInt(fixedDeposit.Term, 10),
		", rate", fixedDeposit.Rate.String(),
	)

	return &types.MsgWithdrawFixedDepositResponse{
		Principal: fixedDeposit.Principal,
		Interest:  fixedDeposit.Interest,
		Term:      fixedDeposit.Term,
		Rate:      fixedDeposit.Rate,
	}, nil
}
