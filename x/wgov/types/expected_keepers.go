package types

import (
	"context"
	addresscodec "cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// StakingKeeper expected staking keeper (Validator and Delegator sets)
type StakingKeeper interface {
	ValidatorAddressCodec() addresscodec.Codec
	// iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(
		context.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool),
	) error

	TotalBondedTokens(context.Context) (math.Int, error) // total bonded tokens within the validator set
	IterateDelegations(
		ctx context.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool),
	) error

	TotalBondedStakePool(ctx sdk.Context) math.Int
	IterateStakes(ctx sdk.Context, delAddr sdk.AccAddress,
		fn func(x int64, del types.Stake) (stop bool),
	)
}
