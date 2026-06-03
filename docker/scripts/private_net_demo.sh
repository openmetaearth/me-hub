#!/bin/bash
# ME-Chain Private Network Demo Script
# This script demonstrates common operations on the private network

set -e

CONTAINER_NAME="mechain-private-net"
KEYRING_BACKEND="test"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_step() {
    echo -e "\n${BLUE}==>${NC} ${GREEN}$1${NC}\n"
}

print_info() {
    echo -e "${YELLOW}$1${NC}"
}

print_command() {
    echo -e "${BLUE}$ $1${NC}"
}

# Check if container is running
check_container() {
    if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        echo -e "${RED}Error: Container ${CONTAINER_NAME} is not running${NC}"
        echo "Please start it first with: make docker-private-net-start"
        exit 1
    fi
}

# Demo 1: Show all accounts
demo_accounts() {
    print_step "Demo 1: Listing all pre-configured accounts"

    print_command "med keys list --keyring-backend test"
    docker exec "${CONTAINER_NAME}" med keys list --keyring-backend ${KEYRING_BACKEND}

    echo ""
    print_info "Getting address of global_dao account:"
    print_command "med keys show global_dao -a --keyring-backend test"
    GLOBAL_DAO_ADDR=$(docker exec "${CONTAINER_NAME}" med keys show global_dao -a --keyring-backend ${KEYRING_BACKEND})
    echo "${GLOBAL_DAO_ADDR}"
}

# Demo 2: Check balances
demo_balances() {
    print_step "Demo 2: Checking account balances"

    for account in global_dao pools user; do
        echo ""
        print_info "Balance of ${account}:"
        ADDR=$(docker exec "${CONTAINER_NAME}" med keys show ${account} -a --keyring-backend ${KEYRING_BACKEND})
        print_command "med query bank balances ${ADDR}"
        docker exec "${CONTAINER_NAME}" med query bank balances ${ADDR}
    done
}

# Demo 3: Send tokens
demo_send_tokens() {
    print_step "Demo 3: Sending tokens between accounts"

    FROM_ADDR=$(docker exec "${CONTAINER_NAME}" med keys show global_dao -a --keyring-backend ${KEYRING_BACKEND})
    TO_ADDR=$(docker exec "${CONTAINER_NAME}" med keys show user -a --keyring-backend ${KEYRING_BACKEND})
    AMOUNT="1000000umec"

    print_info "Sending ${AMOUNT} from global_dao to user"
    print_info "From: ${FROM_ADDR}"
    print_info "To: ${TO_ADDR}"

    print_command "med tx bank send ${FROM_ADDR} ${TO_ADDR} ${AMOUNT} --chain-id mechain_202404-1 --keyring-backend test --yes"

    TX_HASH=$(docker exec "${CONTAINER_NAME}" med tx bank send \
        ${FROM_ADDR} \
        ${TO_ADDR} \
        ${AMOUNT} \
        --chain-id mechain_202404-1 \
        --keyring-backend ${KEYRING_BACKEND} \
        --yes \
        --output json | jq -r '.txhash')

    echo ""
    print_info "Transaction Hash: ${TX_HASH}"

    echo ""
    print_info "Waiting for transaction to be included in a block..."
    sleep 6

    echo ""
    print_command "med query tx ${TX_HASH}"
    docker exec "${CONTAINER_NAME}" med query tx ${TX_HASH} --output json | jq .
}

# Demo 4: Query chain status
demo_chain_status() {
    print_step "Demo 4: Querying chain status"

    print_command "med status"
    docker exec "${CONTAINER_NAME}" med status | jq .

    echo ""
    print_info "Chain ID:"
    docker exec "${CONTAINER_NAME}" med status | jq -r '.NodeInfo.network'

    echo ""
    print_info "Latest Block Height:"
    docker exec "${CONTAINER_NAME}" med status | jq -r '.SyncInfo.latest_block_height'

    echo ""
    print_info "Node Moniker:"
    docker exec "${CONTAINER_NAME}" med status | jq -r '.NodeInfo.moniker'
}

# Demo 5: Query validator info
demo_validator() {
    print_step "Demo 5: Querying validator information"

    print_command "med query staking validators"
    docker exec "${CONTAINER_NAME}" med query staking validators --output json | jq .

    echo ""
    print_info "Validator Operator Address:"
    docker exec "${CONTAINER_NAME}" med query staking validators --output json | jq -r '.validators[0].operator_address'

    echo ""
    print_info "Validator Status:"
    docker exec "${CONTAINER_NAME}" med query staking validators --output json | jq -r '.validators[0].status'

    echo ""
    print_info "Validator Tokens:"
    docker exec "${CONTAINER_NAME}" med query staking validators --output json | jq -r '.validators[0].tokens'
}

# Demo 6: REST API queries
demo_rest_api() {
    print_step "Demo 6: Using REST API endpoints"

    print_info "Node Info:"
    print_command "curl http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info"
    curl -s http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info | jq .

    echo ""
    print_info "Latest Block:"
    print_command "curl http://localhost:1318/cosmos/base/tendermint/v1beta1/blocks/latest"
    curl -s http://localhost:1318/cosmos/base/tendermint/v1beta1/blocks/latest | jq '.block.header'
}

# Demo 7: JSON-RPC queries (Ethereum compatible)
demo_jsonrpc() {
    print_step "Demo 7: Using JSON-RPC endpoint (Ethereum compatible)"

    print_info "Getting current block number:"
    print_command "curl -X POST http://localhost:9545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"eth_blockNumber\",\"params\":[],\"id\":1}'"
    BLOCK_HEX=$(curl -s -X POST http://localhost:9545 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result')

    echo "Block number (hex): ${BLOCK_HEX}"
    echo "Block number (decimal): $((${BLOCK_HEX}))"

    echo ""
    print_info "Getting chain ID:"
    print_command "curl -X POST http://localhost:9545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"eth_chainId\",\"params\":[],\"id\":1}'"
    curl -s -X POST http://localhost:9545 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' | jq .
}

# Demo 8: Create a new account
demo_create_account() {
    print_step "Demo 8: Creating a new account"

    ACCOUNT_NAME="demo_account_$$"

    print_info "Creating new account: ${ACCOUNT_NAME}"
    print_command "med keys add ${ACCOUNT_NAME} --keyring-backend test"
    docker exec "${CONTAINER_NAME}" med keys add ${ACCOUNT_NAME} --keyring-backend ${KEYRING_BACKEND} --output json | jq .

    echo ""
    NEW_ADDR=$(docker exec "${CONTAINER_NAME}" med keys show ${ACCOUNT_NAME} -a --keyring-backend ${KEYRING_BACKEND})
    print_info "New account address: ${NEW_ADDR}"

    echo ""
    print_info "Sending some tokens to the new account..."
    FROM_ADDR=$(docker exec "${CONTAINER_NAME}" med keys show global_dao -a --keyring-backend ${KEYRING_BACKEND})
    docker exec "${CONTAINER_NAME}" med tx bank send \
        ${FROM_ADDR} \
        ${NEW_ADDR} \
        5000000umec \
        --chain-id mechain_202404-1 \
        --keyring-backend ${KEYRING_BACKEND} \
        --yes > /dev/null

    sleep 6

    echo ""
    print_info "New account balance:"
    docker exec "${CONTAINER_NAME}" med query bank balances ${NEW_ADDR}
}

# Demo 9: Query parameters
demo_parameters() {
    print_step "Demo 9: Querying chain parameters"

    print_info "Staking parameters:"
    print_command "med query staking params"
    docker exec "${CONTAINER_NAME}" med query staking params --output json | jq .

    echo ""
    print_info "Bank metadata:"
    print_command "med query bank denom-metadata"
    docker exec "${CONTAINER_NAME}" med query bank denom-metadata --output json | jq .

    echo ""
    print_info "Governance parameters:"
    print_command "med query gov params"
    docker exec "${CONTAINER_NAME}" med query gov params --output json | jq .
}

# Demo 10: Monitor blocks
demo_monitor_blocks() {
    print_step "Demo 10: Monitoring block production"

    print_info "Monitoring block production for 30 seconds..."
    print_info "Press Ctrl+C to stop"

    echo ""
    PREV_HEIGHT=0
    END_TIME=$(($(date +%s) + 30))

    while [ $(date +%s) -lt $END_TIME ]; do
        HEIGHT=$(curl -s http://localhost:36657/status | jq -r '.result.sync_info.latest_block_height')
        BLOCK_TIME=$(curl -s http://localhost:36657/status | jq -r '.result.sync_info.latest_block_time')

        if [ "$HEIGHT" != "$PREV_HEIGHT" ]; then
            echo "Block #${HEIGHT} - ${BLOCK_TIME}"
            PREV_HEIGHT=$HEIGHT
        fi

        sleep 2
    done
}

# Main menu
show_menu() {
    echo ""
    echo "╔════════════════════════════════════════════╗"
    echo "║  ME-Chain Private Network Demo Script     ║"
    echo "╚════════════════════════════════════════════╝"
    echo ""
    echo "Available demos:"
    echo "  1. List all accounts"
    echo "  2. Check account balances"
    echo "  3. Send tokens between accounts"
    echo "  4. Query chain status"
    echo "  5. Query validator information"
    echo "  6. Use REST API endpoints"
    echo "  7. Use JSON-RPC endpoint"
    echo "  8. Create a new account"
    echo "  9. Query chain parameters"
    echo " 10. Monitor block production"
    echo "  a. Run all demos"
    echo "  q. Quit"
    echo ""
}

# Run all demos
run_all_demos() {
    demo_accounts
    demo_balances
    demo_chain_status
    demo_validator
    demo_rest_api
    demo_jsonrpc
    demo_parameters
    demo_send_tokens
    demo_create_account

    print_step "All demos completed!"
    print_info "You can now explore the private network on your own."
}

# Main function
main() {
    check_container

    if [ $# -eq 0 ]; then
        # Interactive mode
        while true; do
            show_menu
            read -p "Select a demo (1-10, a, q): " choice

            case $choice in
                1) demo_accounts ;;
                2) demo_balances ;;
                3) demo_send_tokens ;;
                4) demo_chain_status ;;
                5) demo_validator ;;
                6) demo_rest_api ;;
                7) demo_jsonrpc ;;
                8) demo_create_account ;;
                9) demo_parameters ;;
                10) demo_monitor_blocks ;;
                a|A) run_all_demos ;;
                q|Q) echo "Goodbye!"; exit 0 ;;
                *) echo -e "${RED}Invalid option${NC}" ;;
            esac
        done
    else
        # Non-interactive mode - run all demos
        run_all_demos
    fi
}

# Check dependencies
for cmd in docker curl jq; do
    if ! command -v "$cmd" > /dev/null 2>&1; then
        echo -e "${RED}Error: Required command not found: $cmd${NC}"
        exit 1
    fi
done

main "$@"
