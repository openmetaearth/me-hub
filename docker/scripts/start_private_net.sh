#!/bin/sh
# Startup script for ME-Chain private network in Docker

echo "Starting ME-Chain private network..."
echo "Chain ID: ${CHAIN_ID:-mechain_202404-1}"
echo "Moniker: ${MONIKER_NAME:-local}"

# Check if chain is already initialized
if [ ! -f "$HOME/.mechain/config/genesis.json" ]; then
    echo "Chain not initialized. Running initialization..."
    /scripts/setup_local_docker.sh

    if [ $? -ne 0 ]; then
        echo "Failed to initialize chain"
        exit 1
    fi
fi

echo "Starting chain node..."
echo "RPC endpoint: http://localhost:36657"
echo "API endpoint: http://localhost:1318"
echo "JSON-RPC endpoint: http://localhost:9545"
echo "GRPC endpoint: localhost:8090"

# Start the chain
exec med start --home "$HOME/.mechain"
