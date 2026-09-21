package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

const staticCallGasDescriptor = "EVM static call"

func consumeStaticCallResult(ctx sdk.Context, res *evmtypes.MsgEthereumTxResponse) ([]byte, error) {
	if res == nil {
		return nil, fmt.Errorf("EVM execution returned no result")
	}

	ctx.GasMeter().ConsumeGas(res.GasUsed, staticCallGasDescriptor)
	if res.Failed() {
		return nil, fmt.Errorf("EVM execution failed: %s", res.VmError)
	}
	return res.Ret, nil
}

// StaticCall executes an EVM contract call without committing EVM state.
// Native modules use it for deterministic contract-backed validation.
func (k *Keeper) StaticCall(
	ctx sdk.Context,
	caller common.Address,
	contract common.Address,
	input []byte,
	gas uint64,
) ([]byte, error) {
	// StateDB accesses consume SDK gas too. Isolate that temporary accounting and
	// charge the canonical EVM GasUsed to the original transaction below.
	evmCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	msg := ethtypes.NewMessage(
		caller,
		&contract,
		k.GetNonce(evmCtx, caller),
		big.NewInt(0),
		gas,
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		input,
		ethtypes.AccessList{},
		true,
	)

	res, err := k.ApplyMessage(evmCtx, msg, nil, false)
	if err != nil {
		return nil, err
	}
	return consumeStaticCallResult(ctx, res)
}
