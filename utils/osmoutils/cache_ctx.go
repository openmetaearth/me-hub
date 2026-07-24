package osmoutils

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ApplyFuncIfNoError runs f in a cache context and commits its state changes
// only when f returns without an error or non-gas panic.
func ApplyFuncIfNoError(ctx sdk.Context, f func(ctx sdk.Context) error) (err error) {
	defer func() {
		if recoveryError := recover(); recoveryError != nil {
			if isOutOfGas, _ := IsOutOfGasError(recoveryError); isOutOfGas {
				panic(recoveryError)
			}
			PrintPanicRecoveryError(ctx, recoveryError)
			err = errors.New("panic occurred during execution")
		}
	}()

	cacheCtx, write := ctx.CacheContext()
	err = f(cacheCtx)
	if err != nil {
		ctx.Logger().Error(err.Error())
	} else {
		write()
	}
	return err
}

// IsOutOfGasError reports whether err is an SDK out-of-gas panic value.
func IsOutOfGasError(err any) (bool, string) {
	switch e := err.(type) {
	case storetypes.ErrorOutOfGas:
		return true, e.Descriptor
	case storetypes.ErrorGasOverflow:
		return true, e.Descriptor
	default:
		return false, ""
	}
}

// PrintPanicRecoveryError logs a recovered panic and its stack trace.
func PrintPanicRecoveryError(ctx sdk.Context, recoveryError any) {
	errStackTrace := string(debug.Stack())
	switch e := recoveryError.(type) {
	case storetypes.ErrorOutOfGas:
		ctx.Logger().Debug("out of gas error inside panic recovery block: " + e.Descriptor)
		return
	case string:
		ctx.Logger().Error("recovering from string panic: " + e)
	case runtime.Error:
		ctx.Logger().Error("recovered runtime panic: " + e.Error())
	case error:
		ctx.Logger().Error("recovered error panic: " + e.Error())
	default:
		ctx.Logger().Error("recovered panic; could not capture value in context, see stdout")
		fmt.Println("recovering from panic", recoveryError)
		debug.PrintStack()
		return
	}
	ctx.Logger().Error("stack trace: " + errStackTrace)
}
