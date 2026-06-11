package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func claimReward(amount *big.Int) {
	// Connect to Ethereum client
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set up contract instance
	contractAddress := common.HexToAddress("0x...ContractAddress...")
	contractInstance, err := NewRewardClaim(contractAddress, client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set up transaction options
	txOpts, err := bind.NewKeyedTransactorWithChainID(context.Background(), common.HexToAddress("0x...YourAddress..."), big.NewInt(1))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Call claimReward function
	tx, err := contractInstance.ClaimReward(txOpts, amount)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(receipt)
}