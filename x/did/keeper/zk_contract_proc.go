package keeper

import (
	"math/big"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/openmetaearth/me-hub/x/did/types"
)

const (
	stateIDExistsCallGas = uint64(600_000)
	authV3VerifyCallGas  = uint64(5_000_000)
)

const stateABIJSON = `[
  {"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"idExists","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}
]`

const authV3ValidatorABIJSON = `[
  {"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"bytes","name":"proof","type":"bytes"},{"internalType":"bytes","name":"params","type":"bytes"}],"name":"verify","outputs":[{"internalType":"uint256","name":"userID","type":"uint256"},{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"value","type":"uint256"}],"internalType":"struct IAuthValidator.AuthResponseField[]","name":"authResponseFields","type":"tuple[]"}],"stateMutability":"view","type":"function"}
]`

var (
	stateABI           = mustParseABI(stateABIJSON)
	authV3ValidatorABI = mustParseABI(authV3ValidatorABIJSON)
)

type authResponseField struct {
	Name  string
	Value *big.Int
}

type authV3VerifyResult struct {
	UserID             *big.Int
	AuthResponseFields []authResponseField
}

func mustParseABI(raw string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(raw))
	if err != nil {
		panic(err)
	}
	return parsed
}

func (k Keeper) configuredZkContract(ctx sdk.Context, contractType types.ZkContractAddressType) (common.Address, error) {
	info := k.GetZkContractAddress(ctx, contractType)
	if nil == info {
		return common.Address{}, errors.Wrapf(types.ErrZkContractAddressNotSet,
			"contract type %s", contractType.String())
	}
	if !common.IsHexAddress(info.ContractAddress) {
		return common.Address{}, types.ErrInvalidZkContractAddress
	}
	return common.HexToAddress(info.ContractAddress), nil
}

func (k Keeper) checkZkCoreIDExists(ctx sdk.Context, caller common.Address, coreID *big.Int) error {
	stateAddress, err := k.configuredZkContract(ctx, types.PROXY_STATE)
	if err != nil {
		return err
	}
	input, err := stateABI.Pack("idExists", coreID)
	if err != nil {
		return errors.Wrap(types.ErrZkEVMCall, err.Error())
	}

	ret, err := k.evmKeeper.StaticCall(ctx, caller, stateAddress, input, stateIDExistsCallGas)
	if err != nil {
		return errors.Wrap(types.ErrZkEVMCall, err.Error())
	}
	values, err := stateABI.Unpack("idExists", ret)
	if err != nil || len(values) != 1 {
		return errors.Wrap(types.ErrZkEVMCall, "invalid State.idExists response")
	}
	exists, ok := values[0].(bool)
	if !ok || !exists {
		return errors.Wrap(types.ErrInvalidZkCoreID, "core ID is not registered in State")
	}
	return nil
}

func (k Keeper) verifyZkAuthProof(ctx sdk.Context, caller common.Address, authProof []byte) (*big.Int, *big.Int, error) {
	validatorAddress, err := k.configuredZkContract(ctx, types.PROXY_AUTH_V3_VALIDATOR)
	if err != nil {
		return nil, nil, err
	}
	input, err := authV3ValidatorABI.Pack("verify", caller, authProof, []byte{})
	if err != nil {
		return nil, nil, errors.Wrap(types.ErrInvalidZkProof, err.Error())
	}
	ret, err := k.evmKeeper.StaticCall(ctx, caller, validatorAddress, input, authV3VerifyCallGas)
	if err != nil {
		return nil, nil, errors.Wrap(types.ErrInvalidZkProof, err.Error())
	}

	var result authV3VerifyResult
	if err := authV3ValidatorABI.UnpackIntoInterface(&result, "verify", ret); err != nil {
		return nil, nil, errors.Wrap(types.ErrInvalidZkProof, err.Error())
	}
	if result.UserID == nil {
		return nil, nil, errors.Wrap(types.ErrInvalidZkProof, "missing userID")
	}
	for _, field := range result.AuthResponseFields {
		if field.Name == "challenge" && field.Value != nil {
			return result.UserID, field.Value, nil
		}
	}
	return nil, nil, errors.Wrap(types.ErrInvalidZkProof, "missing challenge response field")
}

