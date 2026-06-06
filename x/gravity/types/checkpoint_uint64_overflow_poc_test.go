package types

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/openmetaearth/me-hub/utils"
)

func TestOutgoingBatchCheckpointMustEncodeUint64TimeoutWithoutInt64Truncation(t *testing.T) {
	const gravityID = "me-gravity"
	const firstUint64AboveInt64 = uint64(1) << 63

	batch := &OutgoingTxBatch{
		BatchNonce:    1,
		BatchTimeout:  firstUint64AboveInt64,
		TokenContract: "0x3fC91A3afd03b6dd40d11400240d25ffb4458d99",
		FeeReceive:    "0x3fC91A3afd03b6dd40d11400240d25ffb4458d99",
	}

	actual, err := batch.GetCheckpoint(gravityID)
	require.NoError(t, err)

	// Calculate the correct expected solidity encoding using big.Int SetUint64
	expected := correctOutgoingBatchCheckpointForUint64Timeout(t, batch, gravityID)

	require.Equal(t, hex.EncodeToString(expected), hex.EncodeToString(actual),
		"batch checkpoint must match the Solidity uint256 encoding of the stored uint64 timeout")
}

func correctOutgoingBatchCheckpointForUint64Timeout(t *testing.T, m *OutgoingTxBatch, gravityIDString string) []byte {
	gravityID, err := utils.StrToByte32(gravityIDString)
	require.NoError(t, err)

	batchMethodName, err := utils.StrToByte32("transactionBatch")
	require.NoError(t, err)

	txAmounts := make([]*big.Int, len(m.Transactions))
	txDestinations := make([]gethcommon.Address, len(m.Transactions))
	txFees := make([]*big.Int, len(m.Transactions))
	for i, tx := range m.Transactions {
		txAmounts[i] = tx.Token.Amount.BigInt()
		txDestinations[i] = toHexAddr(gravityIDString, tx.DestAddress)
		txFees[i] = tx.Fee.Amount.BigInt()
	}

	abiEncodedBatch, err := outgoingBatchTxCheckpointABI.Pack("submitBatch",
		gravityID,
		batchMethodName,
		txAmounts,
		txDestinations,
		txFees,
		new(big.Int).SetUint64(m.BatchNonce),
		toHexAddr(gravityIDString, m.TokenContract),
		new(big.Int).SetUint64(m.BatchTimeout),
		toHexAddr(gravityIDString, m.FeeReceive),
	)
	require.NoError(t, err)

	return crypto.Keccak256Hash(abiEncodedBatch[4:]).Bytes()
}
