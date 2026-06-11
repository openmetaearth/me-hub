#!/bin/bash

# Set up contract address and your Ethereum address
CONTRACT_ADDRESS="0x...ContractAddress..."
YOUR_ADDRESS="0x...YourAddress..."

# Set up amount to claim
AMOUNT=100

# Call claimReward function using med
med claim-reward --contract-address $CONTRACT_ADDRESS --amount $AMOUNT --address $YOUR_ADDRESS