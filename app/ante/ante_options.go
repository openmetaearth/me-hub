package ante

import (
	ante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ethante "github.com/evmos/ethermint/app/ante"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	wbankkeeper "github.com/openmetaearth/me-hub/x/wbank/keeper"

	errorsmod "cosmossdk.io/errors"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

type HandlerOptions struct {
	AccountKeeper          *authkeeper.AccountKeeper
	BankKeeper             wbankkeeper.BaseKeeperWrapper
	IBCKeeper              *ibckeeper.Keeper
	FeeMarketKeeper        ethante.FeeMarketKeeper
	EvmKeeper              ethante.EVMKeeper
	FeegrantKeeper         ante.FeegrantKeeper
	SignModeHandler        authsigning.SignModeHandler
	MaxTxGasWanted         uint64
	ExtensionOptionChecker ante.ExtensionOptionChecker
	RollappKeeper          rollappkeeper.Keeper

	DaoKeeper      DaoKeeper
	StakingKeeper  StakingKeeper
	KycKeeper      KycKeeper
	WasmViewKeeper wasmTypes.ViewKeeper
	TxFeeChecker   ante.TxFeeChecker
}

func (options HandlerOptions) validate() error {
	if options.AccountKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "account keeper is required for AnteHandler")
	}
	//if options.BankKeeper == nil {
	//	return errorsmod.Wrap(errortypes.ErrLogic, "bank keeper is required for AnteHandler")
	//}
	if options.SignModeHandler == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.FeeMarketKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "fee market keeper is required for AnteHandler")
	}
	if options.EvmKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "evm keeper is required for AnteHandler")
	}
	if options.DaoKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "dao keeper is required for AnteHandler")
	}
	if options.StakingKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "staking keeper is required for AnteHandler")
	}
	if options.WasmViewKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "wasm view keeper is required for AnteHandler")
	}
	return nil
}
