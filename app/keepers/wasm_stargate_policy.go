package keepers

import (
	errorsmod "cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func rejectStargateMsg(_ sdk.AccAddress, _ *wasmvmtypes.StargateMsg) ([]sdk.Msg, error) {
	return nil, errorsmod.Wrap(wasmtypes.ErrUnknownMsg, "stargate messages are disabled")
}
