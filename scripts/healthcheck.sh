#!/bin/sh

# Healthcheck script for ME Network node
# Checks if the node is running and responding to RPC calls

RPC_ENDPOINT=${RPC_ENDPOINT:-http://localhost:26657}

# Check if the process is running
if ! pgrep -x "med" > /dev/null 2>&1; then
    echo "med process not running"
    exit 1
fi

# Check if the RPC endpoint is responding
if ! curl -sf "${RPC_ENDPOINT}/status" > /dev/null 2>&1; then
    echo "RPC endpoint not responding"
    exit 1
fi

# Get the latest block height
LATEST_HEIGHT=$(curl -sf "${RPC_ENDPOINT}/status" | jq -r '.result.sync_info.latest_block_height' 2>/dev/null)

if [ -z "$LATEST_HEIGHT" ] || [ "$LATEST_HEIGHT" = "null" ]; then
    echo "Unable to get latest block height"
    exit 1
fi

# Check if the block height is increasing (compare with previous check)
HEIGHT_FILE="/tmp/last_block_height"
if [ -f "$HEIGHT_FILE" ]; then
    LAST_HEIGHT=$(cat "$HEIGHT_FILE")
    if [ "$LATEST_HEIGHT" -le "$LAST_HEIGHT" ] 2>/dev/null; then
        echo "Block height not increasing: last=$LAST_HEIGHT, current=$LATEST_HEIGHT"
        exit 1
    fi
fi

echo "$LATEST_HEIGHT" > "$HEIGHT_FILE"
echo "Node healthy at block height: $LATEST_HEIGHT"
exit 0
