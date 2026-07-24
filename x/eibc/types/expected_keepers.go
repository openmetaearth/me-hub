package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	commontypes "github.com/openmetaearth/me-hub/x/common/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	BlockedAddr(addr sdk.AccAddress) bool
}

type DelayedAckKeeper interface {
	GetRollappPacket(ctx sdk.Context, rollappPacketKey string) (*commontypes.RollappPacket, error)
	BridgingFeeFromAmt(ctx sdk.Context, amt sdkmath.Int) (res sdkmath.Int)
	BridgingFee(ctx sdk.Context) (res sdkmath.LegacyDec)
}
