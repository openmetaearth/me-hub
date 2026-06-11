package main

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func TestClaimReward(t *testing.T) {
	// Connect to Ethereum client
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
	if err != nil {
		t.Fatal(err)
	}

	// Set up contract instance
	contractAddress := common.HexToAddress("0x...ContractAddress...")
	contractInstance, err := NewRewardClaim(contractAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	// Set up transaction options
	txOpts, err := bind.NewKeyedTransactorWithChainID(context.Background(), common.HexToAddress("0x...YourAddress..."), big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}

	// Call claimReward function
	amount := big.NewInt(100)
	tx, err := contractInstance.ClaimReward(txOpts, amount)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, receipt)
}