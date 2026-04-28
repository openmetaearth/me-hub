package ante

import (
	"fmt"
	"math"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	megrouptypes "github.com/openmetaearth/me-hub/x/megroup/types"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
)

const (
	gasEstimationDeductFeeDecorator = 100_000
	priorityScalingFactor           = 100_000_000
	msgLimits                       = 1000
)

var minimumFee = sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(10000)))

// DeductFeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeDecorator struct {
	ak             ante.AccountKeeper
	BankKeeper     BankKeeper
	feegrantKeeper ante.FeegrantKeeper
	daoKeeper      DaoKeeper
	stakingKeeper  StakingKeeper
	kycKeeper      KycKeeper
	txFeeChecker   ante.TxFeeChecker
	wasmKeeper     WasmKeeper
}

func NewDeductFeeDecorator(
	ak ante.AccountKeeper,
	bk BankKeeper,
	fk ante.FeegrantKeeper,
	dk DaoKeeper,
	sk StakingKeeper,
	kycKeeper KycKeeper,
	tfc ante.TxFeeChecker,
	wk WasmKeeper,
) DeductFeeDecorator {
	if tfc == nil {
		tfc = checkTxFeeWithValidatorMinGasPrices
	}

	if ak == nil || fk == nil || dk == nil || sk == nil || wk == nil {
		panic("invalid parameter")
	}

	return DeductFeeDecorator{
		ak:             ak,
		BankKeeper:     bk,
		feegrantKeeper: fk,
		daoKeeper:      dk,
		stakingKeeper:  sk,
		kycKeeper:      kycKeeper,
		txFeeChecker:   tfc,
		wasmKeeper:     wk,
	}
}

func (dfd DeductFeeDecorator) ParseWasmMsgContractCreator(ctx sdk.Context, tx sdk.Tx) (string, bool) {
	// wasm exec message should be the only message in tx
	// to be considered as a wasm transaction
	// this criterion is coarse, refine it later!

	allwasm := true
	var contract string
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *wasmtypes.MsgExecuteContract:
			contract = msg.Contract
		case *wasmtypes.MsgMigrateContract:
			contract = msg.Contract
		case *wasmtypes.MsgUpdateAdmin:
			contract = msg.Contract
		case *wasmtypes.MsgClearAdmin:
			contract = msg.Contract
		case *wasmtypes.MsgSudoContract:
			contract = msg.Contract
		default:
			allwasm = false
		}
	}

	if allwasm {
		addr, err := sdk.AccAddressFromBech32(contract)
		if err != nil {
			return "", false
		}
		contractInfo := dfd.wasmKeeper.GetContractInfo(ctx, addr)
		if contractInfo == nil {
			return "", false
		}
		admin := contractInfo.Creator
		return admin, true
	}
	return "", false
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if len(feeTx.GetMsgs()) > msgLimits {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "messages should not exceed %d", msgLimits)
	}

	if simulate {
		ctx.GasMeter().ConsumeGas(gasEstimationDeductFeeDecorator, "deduct fee decorator")
	}

	var (
		priority int64
		err      error
	)

	feePending := feeTx.GetFee()
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()

	isDao := dfd.daoKeeper.IsDao(ctx, feePayer.String())
	isFreeGasAccount := dfd.daoKeeper.CheckFreeGasAccount(ctx, feePayer.String())
	freeGas := isFreeGasAccount || isDao

	// freeGas for MsgJoinGroup only when ALL messages in the tx are MsgJoinGroup.
	// Mixing MsgJoinGroup with other message types is not allowed to get free gas,
	// preventing attackers from bundling arbitrary messages with MsgJoinGroup to bypass fees.
	if !freeGas && len(feeTx.GetMsgs()) > 0 {
		allJoinGroup := true
		for _, msg := range feeTx.GetMsgs() {
			if _, ok := msg.(*megrouptypes.MsgJoinGroup); !ok {
				allJoinGroup = false
				break
			}
		}
		freeGas = allJoinGroup
	}

	if !freeGas && !simulate {
		_, priority, err = dfd.txFeeChecker(ctx, tx)
		if err != nil {
			return ctx, err
		}
		fee, err := sdk.ParseCoinsNormalized(feePending.String())
		if err != nil {
			return ctx, sdkerrors.Wrap(err, "")
		}

		deductFeesFrom := feePayer

		// if fee granter set deduct fee from fee granter account.
		// this works with only when fee grant enabled.
		if feeGranter != nil {
			if dfd.daoKeeper == nil {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee grants are not enabled")
			} else if !feeGranter.Equals(feePayer) {
				err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, tx.GetMsgs())
				if err != nil {
					return ctx, sdkerrors.Wrapf(err, "%s not allowed to pay fees from %s", feeGranter, feePayer)
				}
			}
			deductFeesFrom = feeGranter
		}

		err = dfd.CheckFunds(ctx, tx, deductFeesFrom.String(), fee)
		if err != nil {
			return ctx, err
		}
		// deduct the fees
		if !fee.IsZero() {
			// DeductFees deducts fees from the given account.
			if !fee.IsValid() {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fee)
			}

			fee10 := make(sdk.Coins, len(fee))
			fee20 := make(sdk.Coins, len(fee))
			fee30 := make(sdk.Coins, len(fee))
			fee40 := make(sdk.Coins, len(fee))

			rate10 := sdk.MustNewDecFromStr("0.1")
			rate20 := sdk.MustNewDecFromStr("0.2")
			rate30 := sdk.MustNewDecFromStr("0.3")

			for i, f := range fee {
				if f.Amount.LT(sdk.NewInt(10)) {
					return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "fee must greater than 10: %s", fee)
				}
				fee10[i] = sdk.NewCoin(f.Denom, rate10.MulInt(f.Amount).TruncateInt())
				fee20[i] = sdk.NewCoin(f.Denom, rate20.MulInt(f.Amount).TruncateInt())
				fee30[i] = sdk.NewCoin(f.Denom, rate30.MulInt(f.Amount).TruncateInt())
				fee40[i] = sdk.NewCoin(f.Denom, f.Amount.Sub(fee10[i].Amount).Sub(fee20[i].Amount).Sub(fee30[i].Amount))
			}
			inputs := []banktypes.Input{
				{
					Address: deductFeesFrom.String(),
					Coins:   fee,
				},
			}

			outputs := []banktypes.Output{}
			feeReceiverTypes := []wbanktypes.FeeReceiverType{}
			outputs = append(outputs, banktypes.Output{
				Address: dfd.daoKeeper.GetDevOperator(ctx),
				Coins:   fee10,
			})
			feeReceiverTypes = append(feeReceiverTypes, wbanktypes.FeeReceiverDevOperator)

			fee20Address := ""
			fee20ReceiverType := wbanktypes.FeeReceiverKycRegionOwner

			kyc, isKyc := didtypes.Credential{}, false
			did, hasDid := dfd.kycKeeper.GetDID(ctx, deductFeesFrom)
			if hasDid {
				kyc, isKyc = dfd.kycKeeper.GetKYC(ctx, did)
			}
			if isKyc {
				fee20Address, err = dfd.stakingKeeper.GetValOwnerAddress(ctx, string(kyc.Data))
				if err != nil {
					return ctx, fmt.Errorf("couldn't get validator from kyc address: %s", deductFeesFrom.String())
				}
			} else {
				fee20Address, err = dfd.stakingKeeper.GetProposerOwnerAddress(ctx)
				if err != nil {
					return ctx, err
				}
				fee20ReceiverType = wbanktypes.FeeReceiverProposerOwner
			}

			outputs = append(outputs, banktypes.Output{Address: fee20Address, Coins: fee20})
			feeReceiverTypes = append(feeReceiverTypes, fee20ReceiverType)

			fee40Address := ""
			globalFee := fee30
			contractOwner, ok := dfd.ParseWasmMsgContractCreator(ctx, tx)
			if ok {
				fee40Address = contractOwner
				fee40ReceiverTypes := wbanktypes.FeeReceiverContractCreator
				feeReceiverTypes = append(feeReceiverTypes, fee40ReceiverTypes)
				outputs = append(outputs, banktypes.Output{
					Address: fee40Address,
					Coins:   fee40,
				})
			} else {
				globalFee = fee30.Add(fee40...)
			}

			outputs = append(outputs, banktypes.Output{
				Address: dfd.daoKeeper.GetGlobalDaoFeePoolAddr(ctx).String(),
				Coins:   globalFee})
			feeReceiverTypes = append(feeReceiverTypes, wbanktypes.FeeReceiverGlobalDaoFeePool)

			err = dfd.BankKeeper.FeeToReceivers(ctx, inputs, outputs, feeReceiverTypes)
			if err != nil {
				return ctx, err
			}
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(sdk.EventTypeTx,
					sdk.NewAttribute(sdk.AttributeKeyFee, feeTx.GetFee().String()),
				),
			)
		}
	}
	newCtx := ctx.WithPriority(priority)
	return next(newCtx, tx, simulate)
}

func (dfd DeductFeeDecorator) CheckFunds(ctx sdk.Context, tx sdk.Tx, feePayer string, fees sdk.Coins) error {
	if len(fees.Denoms()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "denom is empty")
	}

	fromAddress := ""
	userSendAmount := make(map[string]sdk.Coins)
	for _, msg := range tx.GetMsgs() {
		switch txMsg := msg.(type) {
		case *banktypes.MsgSend:
			fromAddress = txMsg.FromAddress
			sendAmount := userSendAmount[txMsg.FromAddress]
			sendAmount = sendAmount.Add(txMsg.Amount...)
			userSendAmount[txMsg.FromAddress] = sendAmount
		case *banktypes.MsgMultiSend:
			if len(txMsg.Inputs) == 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "no input coins provided")
			}
			fromAddress = txMsg.Inputs[0].Address
			sendAmount := userSendAmount[fromAddress]
			for _, output := range txMsg.Outputs {
				sendAmount = sendAmount.Add(output.Coins...)
			}
			userSendAmount[fromAddress] = sendAmount
		case *stakingtypes.MsgDelegate:
			fromAddress = txMsg.DelegatorAddress
			sendAmount := userSendAmount[txMsg.DelegatorAddress]
			sendAmount = sendAmount.Add(txMsg.Amount)
			userSendAmount[txMsg.DelegatorAddress] = sendAmount
		case *wstakingtypes.MsgDoFixedDeposit:
			fromAddress = txMsg.Account
			sendAmount := userSendAmount[txMsg.Account]
			sendAmount = sendAmount.Add(txMsg.Principal)
			userSendAmount[txMsg.Account] = sendAmount
		}
	}

	if _, exists := userSendAmount[feePayer]; !exists {
		userSendAmount[feePayer] = fees
	} else {
		if fromAddress == feePayer {
			userSendAmount[feePayer] = userSendAmount[feePayer].Add(fees...)
		}
	}

	for address, sendAmount := range userSendAmount {
		balance := dfd.BankKeeper.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(address))
		if !balance.IsAllGTE(sendAmount) {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "check funds for %s; got: %s required: %s",
				address, balance, sendAmount)
		}
	}
	return nil
}

// checkTxFeeWithValidatorMinGasPrices implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, can the tx priority is computed from the gas price.
func checkTxFeeWithValidatorMinGasPrices(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() {
		if !feeCoins.IsAllGTE(minimumFee) {
			return sdk.Coins{}, 0, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "fee must greater than or equal %s: got %s", minimumFee.String(), feeCoins.String())
		}
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdk.NewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAllGTE(requiredFees) {
				return nil, 0, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	priority := getTxPriorityByFee(feeCoins)
	return feeCoins, priority, nil
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritize as expected.
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.MulRaw(int64(priorityScalingFactor)).QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}
	return priority
}

func getTxPriorityByFee(fee sdk.Coins) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}
	return priority
}
