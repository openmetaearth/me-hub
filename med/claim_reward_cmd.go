package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

func claimRewardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-reward",
		Short: "Claim reward from contract",
		Run: func(cmd *cobra.Command, args []string) {
			contractAddress, err := cmd.Flags().GetString("contract-address")
			if err != nil {
				fmt.Println(err)
				return
			}

			amount, err := cmd.Flags().GetUint64("amount")
			if err != nil {
				fmt.Println(err)
				return
			}

			address, err := cmd.Flags().GetString("address")
			if err != nil {
				fmt.Println(err)
				return
			}

			// Connect to Ethereum client
			client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
			if err != nil {
				fmt.Println(err)
				return
			}

			// Set up contract instance
			contractInstance, err := NewRewardClaim(common.HexToAddress(contractAddress), client)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Set up transaction options
			txOpts, err := bind.NewKeyedTransactorWithChainID(context.Background(), common.HexToAddress(address), big.NewInt(1))
			if err != nil {
				fmt.Println(err)
				return
			}

			// Call claimReward function
			tx, err := contractInstance.ClaimReward(txOpts, big.NewInt(int64(amount)))
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
		},
	}

	cmd.Flags().String("contract-address", "", "Contract address")
	cmd.Flags().Uint64("amount", 0, "Amount to claim")
	cmd.Flags().String("address", "", "Your Ethereum address")

	return cmd
}