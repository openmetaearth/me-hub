package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrValidatorOwnerAddress = sdkerrors.Register(ModuleName, 56, "cannot update validator owner address to module account address")

	ErrSample           = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrDoDeposit        = sdkerrors.Register(ModuleName, 1101, "do deposit error")
	ErrDoWithdraw       = sdkerrors.Register(ModuleName, 1102, "do withdraw error")
	ErrNoDeposit        = sdkerrors.Register(ModuleName, 1103, "deposit not exist error")
	ErrNotEnoughDeposit = sdkerrors.Register(ModuleName, 1104, "deposit not enough error")
	ErrAmountNotInteger = sdkerrors.Register(ModuleName, 1105, "amount is not integer error")

	ErrDoFixedDeposit           = sdkerrors.Register(ModuleName, 1113, "do fixed deposit error")
	ErrDoFixedWithDraw          = sdkerrors.Register(ModuleName, 1114, "do fixed withdraw error")
	ErrNoFixedDepositFound      = sdkerrors.Register(ModuleName, 1115, "fixed deposit not found")
	ErrFixedDepositNotExpired   = sdkerrors.Register(ModuleName, 1116, "fixed deposit not reach deadline")
	ErrFixedDepositInvalidPayee = sdkerrors.Register(ModuleName, 1117, "fixed deposit payee error")
	ErrPayInterest              = sdkerrors.Register(ModuleName, 1118, "pay interest err.")

	ErrAddFixedDepositConfig           = sdkerrors.Register(ModuleName, 1120, "add fixed deposit config error")
	ErrRemoveFixedDepositConfig        = sdkerrors.Register(ModuleName, 1121, "remove fixed deposit config error")
	ErrNoFixedDepositCountOfCfgFound   = sdkerrors.Register(ModuleName, 1122, "fixed deposit count under config not found")
	ErrFixedDepositCountOfCfgIsZero    = sdkerrors.Register(ModuleName, 1123, "fixed deposit count under config is zero")
	ErrFixedDepositConfigAlreadyExists = sdkerrors.Register(ModuleName, 1124, "fixed deposit config already exists")
	ErrFixedDepositExistUnderConfig    = sdkerrors.Register(ModuleName, 1125, "fixed deposit exists under the config")
	ErrSetFixedDepositConfigRate       = sdkerrors.Register(ModuleName, 1126, "set fixed deposit config rate error")
	ErrSetFixedDepositConfigStatus     = sdkerrors.Register(ModuleName, 1127, "set fixed deposit config status error")
	ErrFixedDepositConfigInactive      = sdkerrors.Register(ModuleName, 1128, "fixed deposit config status inactive")
	ErrFixedDepositConfigRateInvalid   = sdkerrors.Register(ModuleName, 1129, "fixed deposit config rate invalid, rate range:[0.0001, 10000]")

	ErrorBonusNew      = sdkerrors.Register(ModuleName, 1130, "new bonus error")
	ErrorBonusRetrieve = sdkerrors.Register(ModuleName, 1131, "bonus retrieve error")
	ErrorBonusStatus   = sdkerrors.Register(ModuleName, 1132, "bonus status error")
	ErrorBonusNotExist = sdkerrors.Register(ModuleName, 1133, "bonus not exist")

	ErrVaultAccountExists    = sdkerrors.Register(ModuleName, 1140, "region vault account exists")
	ErrSendCoinToRegionAdmin = sdkerrors.Register(ModuleName, 1141, "send coins to region admin error")
	ErrSendCoinToRegionVault = sdkerrors.Register(ModuleName, 1142, "send coins to region vault error")
	ErrSetRegionAnnualRate   = sdkerrors.Register(ModuleName, 1143, "set region annual rate error")
	ErrSendCoinToValOwner    = sdkerrors.Register(ModuleName, 1144, "send coins to validator owner error")
	ErrSendCoinToDevOperator = sdkerrors.Register(ModuleName, 1145, "send coins to region vault error")
	ErrSendCoinToGlobalAdmin = sdkerrors.Register(ModuleName, 1146, "send coins to region vault error")
	ErrGetGlobalAdmin        = sdkerrors.Register(ModuleName, 1147, "get global admin error")
	ErrCheckGlobalDao        = sdkerrors.Register(ModuleName, 1148, "invalid global dao")

	ErrMeidNew               = sdkerrors.Register(ModuleName, 1150, "new meid error")
	ErrMeidRemove            = sdkerrors.Register(ModuleName, 1151, "remove meid error")
	ErrMeidExists            = sdkerrors.Register(ModuleName, 1152, "meid already exists")
	ErrMeidNotExists         = sdkerrors.Register(ModuleName, 1153, "meid not exists")
	ErrRegionVaultNotExists  = sdkerrors.Register(ModuleName, 1154, "region vault not exists")
	ErrInvalidMeidRole       = sdkerrors.Register(ModuleName, 1155, "invalid meid role value")
	ErrSetMeidMinStaking     = sdkerrors.Register(ModuleName, 1156, "set meid min staking error")
	ErrSetMeidMaxStaking     = sdkerrors.Register(ModuleName, 1157, "set meid max staking error")
	ErrMeidRegionAdminExists = sdkerrors.Register(ModuleName, 1158, "meid region admin already exists")
	ErrMeidRemoveRegionAdmin = sdkerrors.Register(ModuleName, 1159, "remove meid region admin error")

	ErrSetClearHeight = sdkerrors.Register(ModuleName, 1160, "set clear height failed")
	ErrGetClearHeight = sdkerrors.Register(ModuleName, 1161, "get clear height failed")

	ErrRegionAlreadyExist      = sdkerrors.Register(ModuleName, 1162, "region already exist")
	ErrRegionNotExist          = sdkerrors.Register(ModuleName, 1163, "region not exist")
	ErrRegionValidatorNotExist = sdkerrors.Register(ModuleName, 1164, "region validator not exist")
	ErrTransferRegion          = sdkerrors.Register(ModuleName, 1165, "transfer region err")
	ErrRegion                  = sdkerrors.Register(ModuleName, 1166, "region err")

	ErrFixedDepositAnnualRateNotExists = sdkerrors.Register(ModuleName, 1170, "fixed deposit vault not exists")
	ErrAmountNotPositive               = sdkerrors.Register(ModuleName, 1171, "amount not positive")
	ErrAmountLessThanMin               = sdkerrors.Register(ModuleName, 1172, "amount is less than 0.01ME")
	ErrRegionValidatorDuplicate        = sdkerrors.Register(ModuleName, 1173, "input validator duplicates with already bonded validators")
	ErrRegionNameDuplicate             = sdkerrors.Register(ModuleName, 1174, "input region name duplicates with the name of existing region")
	ErrRegionName                      = sdkerrors.Register(ModuleName, 1175, "invalid region name")
	ErrExpRegionNotExist               = sdkerrors.Register(ModuleName, 1176, "experience region not exist")

	ErrMeidNFTNew       = sdkerrors.Register(ModuleName, 1180, "new meid-nft error")
	ErrMeidNFTRemove    = sdkerrors.Register(ModuleName, 1181, "remove meid-nft error")
	ErrMeidNFTExists    = sdkerrors.Register(ModuleName, 1182, "meid-nft already exists")
	ErrMeidNFTNotExists = sdkerrors.Register(ModuleName, 1183, "meid-nft not exists")

	ErrParameter      = sdkerrors.Register(ModuleName, 1201, "parameter error")
	ErrUnknownAccount = sdkerrors.Register(ModuleName, 1202, "Unknown account")

	ErrNodeLimitExceeded             = sdkerrors.Register(ModuleName, 1203, "Node delegation limit exceeded.")
	ErrAssertionFailed               = sdkerrors.Register(ModuleName, 1204, "type Delegation assertion failed")
	ErrCalculateInterest             = sdkerrors.Register(ModuleName, 1205, "delegator calculate interest err.")
	ErrValidatorDelegationAmount     = sdkerrors.Register(ModuleName, 1206, "Validator DelegationAmount less than requested value.")
	ErrEmptyDelegationDistInfo       = sdkerrors.Register(ModuleName, 1207, "no delegation distribution info")
	ErrNotEnoughDelegationAmount     = sdkerrors.Register(ModuleName, 1208, "not enough delegation amount")
	ErrMaxUnbondingDelegationEntries = sdkerrors.Register(ModuleName, 1209, "too many unbonding delegation entries for (delegator, validator) tuple")
	ErrNoDelegatorForAddress         = sdkerrors.Register(ModuleName, 1210, "delegator does not contain delegation")
)
