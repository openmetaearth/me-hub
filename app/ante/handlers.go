package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	ethante "github.com/evmos/ethermint/app/ante"
	"github.com/st-chain/me-hub/x/rollapp/transfergenesis"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	proofheightante "github.com/st-chain/me-hub/x/delayedack/ante"
)

func newEthAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.NewEthSetUpContextDecorator(options.EvmKeeper),

		// TODO: need to allow universal fees for Eth as well
		ethante.NewEthMempoolFeeDecorator(options.EvmKeeper),                           // Check eth effective gas price against minimal-gas-prices
		ethante.NewEthMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper), // Check eth effective gas price against the global MinGasPrice

		ethante.NewEthValidateBasicDecorator(options.EvmKeeper),
		ethante.NewEthSigVerificationDecorator(options.EvmKeeper),
		ethante.NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper),
		ethante.NewCanTransferDecorator(options.EvmKeeper),
		ethante.NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted),
		ethante.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper), // innermost AnteDecorator.
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
		ethante.NewEthEmitEventDecorator(options.EvmKeeper), // emit eth tx hash and index at the very last ante handler.
	)
}

func newCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	deductFeeDecorator := NewDeductFeeDecorator(
		options.AccountKeeper,
		options.BankKeeper,
		options.FeegrantKeeper,
		options.DaoKeeper,
		options.StakingKeeper,
		options.KycKeeper,
		options.TxFeeChecker,
		options.WasmViewKeeper,
	)
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		// reject MsgEthereumTxs and disable the Msg types that cannot be included on an authz.MsgExec msgs field
		NewRejectMessagesDecorator().
			WithPredicate(BlockTypeUrls(
				0,
				sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
				sdk.MsgTypeURL(&ibcclienttypes.MsgSubmitMisbehaviour{}))), // blocked to avoid skipping our validation logic in lightclient ante handler

		deductFeeDecorator,
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, ethante.DefaultSigVerificationGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),

		// decorator that runs our custom logic for all IBC messages, even wrapped msgs
		NewInnerDecorator(
			proofheightante.NewIBCProofHeightDecorator().InnerCallback,
		),

		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),

		transfergenesis.NewTransferEnabledDecorator(options.RollappKeeper.GetRollapp, options.IBCKeeper.ChannelKeeper),
	)
}
