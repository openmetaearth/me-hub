package types

import (
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/abi"
)

// GetCheckpointRelayerSet returns the checkpoint
func GetCheckpointRelayerSet(relayerSet *types.RelayerSet, gravityIDStr string) ([]byte, error) {
	addresses := make([]string, len(relayerSet.Members))
	powers := make([]*big.Int, len(relayerSet.Members))
	for i, member := range relayerSet.Members {
		addresses[i] = member.ExternalAddress
		powers[i] = big.NewInt(int64(member.Power))
	}

	gravityID, err := utils.StrToByte32(gravityIDStr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}
	checkpoint, err := utils.StrToByte32("checkpoint")
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse checkpoint")
	}

	params := []abi.Param{
		{"bytes32": gravityID},
		{"bytes32": checkpoint},
		{"uint256": big.NewInt(int64(relayerSet.Nonce))},
		{"address[]": addresses},
		{"uint256[]": powers},
	}
	encode, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encode), nil
}

func GetCheckpointConfirmBatch(txBatch *types.OutgoingTxBatch, gravityIDStr string) ([]byte, error) {
	txCount := len(txBatch.Transactions)
	amounts := make([]*big.Int, txCount)
	destinations := make([]string, txCount)
	fees := make([]*big.Int, txCount)
	for i, transferTx := range txBatch.Transactions {
		amounts[i] = transferTx.Token.Amount.BigInt()
		destinations[i] = transferTx.DestAddress
		fees[i] = transferTx.Fee.Amount.BigInt()
	}

	gravityID, err := utils.StrToByte32(gravityIDStr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}
	transactionBatch, err := utils.StrToByte32("transactionBatch")
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse checkpoint")
	}

	params := []abi.Param{
		{"bytes32": gravityID},
		{"bytes32": transactionBatch},
		{"uint256[]": amounts},
		{"address[]": destinations},
		{"uint256[]": fees},
		{"uint256": big.NewInt(int64(txBatch.BatchNonce))},
		{"address": txBatch.TokenContract},
		{"uint256": big.NewInt(int64(txBatch.BatchTimeout))},
		{"address": txBatch.FeeReceive},
	}

	encode, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encode), nil
}
