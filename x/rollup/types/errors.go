package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/rollapp module sentinel errors
var (
	ErrRollappDisable            = errorsmod.Register(MODULE_NAME, 1100, "rollapp is disable")
	ErrRollappNotExist           = errorsmod.Register(MODULE_NAME, 1101, "rollapp is not exist")
	ErrRollappVersionMismatch    = errorsmod.Register(MODULE_NAME, 1102, "rollapp's version is mismatch")
	ErrInputDataErr              = errorsmod.Register(MODULE_NAME, 1103, "Input data error.")
	ErrInsufficientBalance       = errorsmod.Register(MODULE_NAME, 1105, "Insufficient Balance")
	ErrParserDataErr             = errorsmod.Register(MODULE_NAME, 1106, "Parser data error")
	ErrStakeDataErr              = errorsmod.Register(MODULE_NAME, 1107, "Stake error")
	ErrUnStakeLimit              = errorsmod.Register(MODULE_NAME, 1108, "Unstake limit")
	ErrProcessErr                = errorsmod.Register(MODULE_NAME, 1109, "process error")
	ErrStakeTimeoutLimit         = errorsmod.Register(MODULE_NAME, 1110, "Stake's time exceeds stake's end time")
	ErrUnStakeProc               = errorsmod.Register(MODULE_NAME, 1111, "unStake process error")
	ErrRollappIdRegisterRepeated = errorsmod.Register(MODULE_NAME, 1112, "Rollapp Id has been registered")
	//ErrRollappIdNotRegister      = errorsmod.Register(MODULE_NAME, 1112, "RollappId unregistere")
	ErrRollappIDMismatch    = errorsmod.Register(MODULE_NAME, 1113, "rollappid  mismatch")
	ErrNotFound             = errorsmod.Register(MODULE_NAME, 1115, "data not found")
	ErrStakeDaFraudRepeated = errorsmod.Register(MODULE_NAME, 1116, "Stake challenge DA fraud Repeated ")
	ErrInBlackList          = errorsmod.Register(MODULE_NAME, 1117, "address in black list ")
)
