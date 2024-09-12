package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/rollup/types"
)

func (k Keeper) StakeForDaFraud(ctx sdk.Context, rollappID, challenger string, challengeKey []byte) error {
	challengerAddr, err := sdk.AccAddressFromBech32(challenger)
	if err != nil {
		return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("AccAddressFromBech32 error. err = %s", err.Error()))
	}
	//一次DA欺诈挑战只质押一次，此时的Key为challengeKey，如果之前已经存在，则认为之前已经质押，此时出错
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupAppStakeForChallengeDaFraud(rollappID))
	data := store.Get(challengeKey)
	if data != nil {
		return errorsmod.Wrapf(types.ErrStakeDaFraudRepeated, fmt.Sprintf("challengeKey = %s", hex.EncodeToString(challengeKey)))
	}
	/*
		balanceCoin := k.bk.GetBalance(ctx, challengerAddr, params.BaseDenom)
		if balanceCoin.IsLT(stakeCoin) {
			return errorsmod.Wrapf(types.ErrInsufficientBalance, fmt.Sprintf(",user's balanceCoin = %s,but need stake= %dumec",
				balanceCoin.String(), stakeForChallenge))
		}*/
	stakeForChallenge := uint64(k.GetDaFraudChallengeStake(ctx)) * types.MecPrecision
	stakeCoin := sdk.NewCoin(params.BaseDenom, sdk.NewInt(int64(stakeForChallenge)))
	//如果金额不够的话，SendCoinsFromAccountToModule这里就已经会判断处理了
	if err = k.bk.SendCoinsFromAccountToModule(ctx, challengerAddr, types.MODULE_NAME, sdk.NewCoins(stakeCoin)); err != nil {
		return errorsmod.Wrapf(types.ErrStakeDataErr, fmt.Sprintf("stake coin to module error.err = %s", err.Error()))
	}
	stakeMsg := &types.MsgStakeChallengeDaFraud{
		Challenger: challenger,
		Denom:      params.BaseDenom,
		Amount:     stakeForChallenge,
	}
	store.Set(challengeKey, k.cdc.MustMarshal(stakeMsg))
	return nil

}
func (t *rollupServer) SubmitDaFraudProof(ctx context.Context, req *types.MsgDaFraudProofRequest) (*types.MsgDaFraudProofResponse, error) {

}
