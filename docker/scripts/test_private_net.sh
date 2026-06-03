#!/bin/bash
# Test script for ME-Chain Docker private network
# This script verifies that the private network is running correctly

set -e

CONTAINER_NAME=${CONTAINER_NAME:-"mechain-private-net"}
RPC_ENDPOINT=${RPC_ENDPOINT:-"http://localhost:36657"}
API_ENDPOINT=${API_ENDPOINT:-"http://localhost:1318"}
JSONRPC_ENDPOINT=${JSONRPC_ENDPOINT:-"http://localhost:9545"}
GRPC_ENDPOINT=${GRPC_ENDPOINT:-"localhost:8090"}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

print_header() {
    echo ""
    echo "=========================================="
    echo "$1"
    echo "=========================================="
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}✗${NC} $1"
    ((TESTS_FAILED++))
}

print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Check if container is running
test_container_running() {
    print_header "Testing Container Status"

    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        print_success "Container is running"
        return 0
    else
        print_error "Container is not running"
        return 1
    fi
}

# Test RPC endpoint
test_rpc_endpoint() {
    print_header "Testing RPC Endpoint"

    if curl -s --max-time 5 "${RPC_ENDPOINT}/status" > /dev/null 2>&1; then
        local chain_id=$(curl -s "${RPC_ENDPOINT}/status" | jq -r '.result.node_info.network')
        local latest_height=$(curl -s "${RPC_ENDPOINT}/status" | jq -r '.result.sync_info.latest_block_height')

        print_success "RPC endpoint is accessible"
        print_info "Chain ID: ${chain_id}"
        print_info "Latest block height: ${latest_height}"

        if [ "$latest_height" -gt 0 ]; then
            print_success "Chain is producing blocks"
        else
            print_error "Chain is not producing blocks"
        fi
    else
        print_error "RPC endpoint is not accessible at ${RPC_ENDPOINT}"
    fi
}

# Test REST API endpoint
test_api_endpoint() {
    print_header "Testing REST API Endpoint"

    if curl -s --max-time 5 "${API_ENDPOINT}/cosmos/base/tendermint/v1beta1/node_info" > /dev/null 2>&1; then
        local moniker=$(curl -s "${API_ENDPOINT}/cosmos/base/tendermint/v1beta1/node_info" | jq -r '.default_node_info.moniker')
        print_success "REST API endpoint is accessible"
        print_info "Node moniker: ${moniker}"
    else
        print_error "REST API endpoint is not accessible at ${API_ENDPOINT}"
    fi
}

# Test JSON-RPC endpoint
test_jsonrpc_endpoint() {
    print_header "Testing JSON-RPC Endpoint"

    local response=$(curl -s --max-time 5 -X POST "${JSONRPC_ENDPOINT}" \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' 2>/dev/null)

    if [ ! -z "$response" ]; then
        local block_number=$(echo "$response" | jq -r '.result')
        if [ "$block_number" != "null" ]; then
            print_success "JSON-RPC endpoint is accessible"
            print_info "Block number (hex): ${block_number}"
            print_info "Block number (decimal): $((${block_number}))"
        else
            print_error "JSON-RPC endpoint returned null"
        fi
    else
        print_error "JSON-RPC endpoint is not accessible at ${JSONRPC_ENDPOINT}"
    fi
}

# Test accounts
test_accounts() {
    print_header "Testing Pre-configured Accounts"

    # Check if global_dao account exists
    if docker exec "${CONTAINER_NAME}" med keys show global_dao --keyring-backend test > /dev/null 2>&1; then
        local address=$(docker exec "${CONTAINER_NAME}" med keys show global_dao -a --keyring-backend test)
        print_success "global_dao account exists"
        print_info "Address: ${address}"

        # Check balance
        local balance=$(docker exec "${CONTAINER_NAME}" med query bank balances "${address}" --output json 2>/dev/null | jq -r '.balances[0].amount')
        if [ ! -z "$balance" ] && [ "$balance" != "null" ]; then
            print_success "Account has balance"
            print_info "Balance: ${balance} umec"
        else
            print_error "Account has no balance"
        fi
    else
        print_error "global_dao account not found"
    fi

    # Check other accounts
    for account in pools user sequencer; do
        if docker exec "${CONTAINER_NAME}" med keys show "${account}" --keyring-backend test > /dev/null 2>&1; then
            print_success "${account} account exists"
        else
            print_error "${account} account not found"
        fi
    done
}

# Test validator
test_validator() {
    print_header "Testing Validator Status"

    local validator_info=$(docker exec "${CONTAINER_NAME}" med query staking validators --output json 2>/dev/null)
    local validator_count=$(echo "$validator_info" | jq '.validators | length')

    if [ "$validator_count" -gt 0 ]; then
        print_success "Validator is active"
        local validator_address=$(echo "$validator_info" | jq -r '.validators[0].operator_address')
        local validator_status=$(echo "$validator_info" | jq -r '.validators[0].status')
        print_info "Validator address: ${validator_address}"
        print_info "Validator status: ${validator_status}"
    else
        print_error "No validator found"
    fi
}

# Test genesis configuration
test_genesis() {
    print_header "Testing Genesis Configuration"

    # Check chain-id
    local chain_id=$(docker exec "${CONTAINER_NAME}" med status 2>/dev/null | jq -r '.NodeInfo.network')
    if [ "$chain_id" == "mechain_202404-1" ] || [ ! -z "$chain_id" ]; then
        print_success "Chain ID is configured"
        print_info "Chain ID: ${chain_id}"
    else
        print_error "Chain ID is not properly configured"
    fi

    # Check staking denom
    local staking_params=$(curl -s "${API_ENDPOINT}/cosmos/staking/v1beta1/params" 2>/dev/null)
    local bond_denom=$(echo "$staking_params" | jq -r '.params.bond_denom')
    if [ "$bond_denom" == "umec" ]; then
        print_success "Staking denom is correctly set to umec"
    else
        print_error "Staking denom is not umec: ${bond_denom}"
    fi
}

# Test block production
test_block_production() {
    print_header "Testing Block Production"

    local height1=$(curl -s "${RPC_ENDPOINT}/status" | jq -r '.result.sync_info.latest_block_height')
    print_info "Current height: ${height1}"

    print_info "Waiting 10 seconds for new blocks..."
    sleep 10

    local height2=$(curl -s "${RPC_ENDPOINT}/status" | jq -r '.result.sync_info.latest_block_height')
    print_info "New height: ${height2}"

    if [ "$height2" -gt "$height1" ]; then
        local blocks_produced=$((height2 - height1))
        print_success "Chain is producing blocks"
        print_info "Blocks produced: ${blocks_produced}"
    else
        print_error "Chain is not producing blocks"
    fi
}

# Main test execution
main() {
    echo "╔════════════════════════════════════════╗"
    echo "║  ME-Chain Private Network Test Suite  ║"
    echo "╚════════════════════════════════════════╝"
    echo ""
    print_info "Container: ${CONTAINER_NAME}"
    print_info "RPC: ${RPC_ENDPOINT}"
    print_info "API: ${API_ENDPOINT}"
    print_info "JSON-RPC: ${JSONRPC_ENDPOINT}"

    # Run all tests
    test_container_running || exit 1
    test_rpc_endpoint
    test_api_endpoint
    test_jsonrpc_endpoint
    test_accounts
    test_validator
    test_genesis
    test_block_production

    # Summary
    print_header "Test Summary"
    echo "Tests Passed: ${TESTS_PASSED}"
    echo "Tests Failed: ${TESTS_FAILED}"

    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        exit 0
    else
        echo -e "${RED}Some tests failed!${NC}"
        exit 1
    fi
}

# Check dependencies
check_dependencies() {
    local missing_deps=0

    for cmd in docker curl jq; do
        if ! command -v "$cmd" > /dev/null 2>&1; then
            print_error "Required command not found: $cmd"
            ((missing_deps++))
        fi
    done

    if [ $missing_deps -gt 0 ]; then
        echo -e "${RED}Please install missing dependencies${NC}"
        exit 1
    fi
}

# Run checks and tests
check_dependencies
main
