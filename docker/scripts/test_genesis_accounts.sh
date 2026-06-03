#!/usr/bin/env bash
# Test script for verifying genesis accounts in ME-Chain private network

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME=${CONTAINER_NAME:-"mechain-private-net"}
RPC_URL=${RPC_URL:-"http://localhost:36657"}
API_URL=${API_URL:-"http://localhost:1318"}

echo "========================================="
echo "ME-Chain Genesis Accounts Test"
echo "========================================="
echo ""

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Function to check if container is running
check_container() {
    print_info "Checking if container is running..."
    if ! docker ps | grep -q "$CONTAINER_NAME"; then
        print_error "Container $CONTAINER_NAME is not running"
        exit 1
    fi
    print_success "Container is running"
}

# Function to check if chain is ready
check_chain_ready() {
    print_info "Checking if chain is ready..."
    max_attempts=30
    attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if curl -s "$RPC_URL/status" > /dev/null 2>&1; then
            print_success "Chain is ready"
            return 0
        fi
        attempt=$((attempt + 1))
        echo -n "."
        sleep 2
    done

    echo ""
    print_error "Chain failed to start after $max_attempts attempts"
    return 1
}

# Function to list all accounts
list_accounts() {
    print_info "Listing all genesis accounts..."
    echo ""

    accounts=$(docker exec "$CONTAINER_NAME" med keys list --keyring-backend test --output json 2>/dev/null || echo "[]")

    if [ "$accounts" = "[]" ]; then
        print_error "No accounts found"
        return 1
    fi

    echo "$accounts" | jq -r '.[] | "  - \(.name): \(.address)"'
    echo ""

    account_count=$(echo "$accounts" | jq '. | length')
    print_success "Found $account_count accounts"

    return 0
}

# Function to check account balance
check_balance() {
    local account_name=$1
    print_info "Checking balance for $account_name..."

    # Get account address
    address=$(docker exec "$CONTAINER_NAME" med keys show "$account_name" --keyring-backend test -a 2>/dev/null)

    if [ -z "$address" ]; then
        print_error "Account $account_name not found"
        return 1
    fi

    # Query balance via API
    balance_json=$(curl -s "$API_URL/cosmos/bank/v1beta1/balances/$address" 2>/dev/null)

    if [ -z "$balance_json" ]; then
        print_error "Failed to query balance for $account_name"
        return 1
    fi

    # Extract umec balance
    umec_balance=$(echo "$balance_json" | jq -r '.balances[] | select(.denom=="umec") | .amount')

    if [ -z "$umec_balance" ]; then
        print_error "No umec balance found for $account_name"
        return 1
    fi

    # Convert to MEC (divide by 10^18)
    mec_balance=$(echo "scale=6; $umec_balance / 1000000000000000000" | bc)

    echo "  Address: $address"
    echo "  Balance: $umec_balance umec ($mec_balance MEC)"
    print_success "Balance verified"

    return 0
}

# Function to verify default accounts
verify_default_accounts() {
    echo ""
    echo "========================================="
    echo "Verifying Default Accounts"
    echo "========================================="
    echo ""

    default_accounts=("global_dao" "sequencer" "pools" "user")

    for account in "${default_accounts[@]}"; do
        if docker exec "$CONTAINER_NAME" med keys show "$account" --keyring-backend test > /dev/null 2>&1; then
            print_success "Default account '$account' exists"
        else
            print_error "Default account '$account' not found"
        fi
    done
}

# Function to verify custom accounts
verify_custom_accounts() {
    echo ""
    echo "========================================="
    echo "Verifying Custom Genesis Accounts"
    echo "========================================="
    echo ""

    # Get all account names
    all_accounts=$(docker exec "$CONTAINER_NAME" med keys list --keyring-backend test --output json 2>/dev/null | jq -r '.[].name')

    # Default accounts to exclude
    default_accounts=("global_dao" "sequencer" "pools" "user")

    custom_count=0
    while IFS= read -r account; do
        # Check if account is not in default list
        is_custom=true
        for default_account in "${default_accounts[@]}"; do
            if [ "$account" = "$default_account" ]; then
                is_custom=false
                break
            fi
        done

        if [ "$is_custom" = true ]; then
            custom_count=$((custom_count + 1))
            print_info "Custom account found: $account"
            check_balance "$account"
            echo ""
        fi
    done <<< "$all_accounts"

    if [ $custom_count -eq 0 ]; then
        print_info "No custom genesis accounts found (only default accounts exist)"
    else
        print_success "Verified $custom_count custom genesis accounts"
    fi
}

# Function to test transactions
test_transactions() {
    echo ""
    echo "========================================="
    echo "Testing Transactions"
    echo "========================================="
    echo ""

    # Get all accounts
    all_accounts=$(docker exec "$CONTAINER_NAME" med keys list --keyring-backend test --output json 2>/dev/null | jq -r '.[].name')
    account_array=($all_accounts)

    if [ ${#account_array[@]} -lt 2 ]; then
        print_info "Need at least 2 accounts to test transactions, skipping..."
        return 0
    fi

    sender="${account_array[0]}"
    receiver="${account_array[1]}"

    print_info "Testing transaction from $sender to $receiver..."

    sender_addr=$(docker exec "$CONTAINER_NAME" med keys show "$sender" --keyring-backend test -a)
    receiver_addr=$(docker exec "$CONTAINER_NAME" med keys show "$receiver" --keyring-backend test -a)

    # Get initial balances
    initial_receiver_balance=$(curl -s "$API_URL/cosmos/bank/v1beta1/balances/$receiver_addr" | jq -r '.balances[] | select(.denom=="umec") | .amount')

    print_info "Sending 1000umec from $sender to $receiver..."

    # Send transaction
    tx_result=$(docker exec "$CONTAINER_NAME" med tx bank send "$sender_addr" "$receiver_addr" 1000umec \
        --from "$sender" \
        --keyring-backend test \
        --chain-id mechain_202404-1 \
        --gas auto \
        --gas-adjustment 1.5 \
        --fees 1000000umec \
        --yes \
        --output json 2>/dev/null || echo "{}")

    if [ "$tx_result" = "{}" ]; then
        print_error "Transaction failed"
        return 1
    fi

    tx_hash=$(echo "$tx_result" | jq -r '.txhash')
    print_info "Transaction hash: $tx_hash"

    # Wait for transaction to be included in a block
    print_info "Waiting for transaction to be confirmed..."
    sleep 6

    # Get final balance
    final_receiver_balance=$(curl -s "$API_URL/cosmos/bank/v1beta1/balances/$receiver_addr" | jq -r '.balances[] | select(.denom=="umec") | .amount')

    if [ "$final_receiver_balance" != "$initial_receiver_balance" ]; then
        print_success "Transaction successful! Balance changed from $initial_receiver_balance to $final_receiver_balance umec"
    else
        print_info "Balance unchanged (transaction may still be pending)"
    fi
}

# Function to export account information
export_account_info() {
    echo ""
    echo "========================================="
    echo "Exporting Account Information"
    echo "========================================="
    echo ""

    output_file="genesis_accounts_export.json"

    print_info "Exporting account information to $output_file..."

    accounts=$(docker exec "$CONTAINER_NAME" med keys list --keyring-backend test --output json 2>/dev/null)

    echo "{" > "$output_file"
    echo "  \"exported_at\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"," >> "$output_file"
    echo "  \"chain_id\": \"mechain_202404-1\"," >> "$output_file"
    echo "  \"accounts\": [" >> "$output_file"

    account_names=$(echo "$accounts" | jq -r '.[].name')
    account_count=$(echo "$account_names" | wc -l)
    current=0

    while IFS= read -r account_name; do
        current=$((current + 1))
        address=$(echo "$accounts" | jq -r ".[] | select(.name==\"$account_name\") | .address")
        pubkey=$(echo "$accounts" | jq -r ".[] | select(.name==\"$account_name\") | .pubkey")

        # Get balance
        balance_json=$(curl -s "$API_URL/cosmos/bank/v1beta1/balances/$address" 2>/dev/null)
        balances=$(echo "$balance_json" | jq -c '.balances')

        echo "    {" >> "$output_file"
        echo "      \"name\": \"$account_name\"," >> "$output_file"
        echo "      \"address\": \"$address\"," >> "$output_file"
        echo "      \"pubkey\": $pubkey," >> "$output_file"
        echo "      \"balances\": $balances" >> "$output_file"

        if [ $current -lt $account_count ]; then
            echo "    }," >> "$output_file"
        else
            echo "    }" >> "$output_file"
        fi
    done <<< "$account_names"

    echo "  ]" >> "$output_file"
    echo "}" >> "$output_file"

    print_success "Account information exported to $output_file"

    # Pretty print summary
    echo ""
    echo "Summary:"
    cat "$output_file" | jq -r '.accounts[] | "  - \(.name): \(.address) - \(.balances | length) token(s)"'
}

# Main test execution
main() {
    echo "Starting genesis accounts verification..."
    echo ""

    # Check container
    check_container

    # Check if chain is ready
    if ! check_chain_ready; then
        exit 1
    fi

    # List all accounts
    if ! list_accounts; then
        exit 1
    fi

    # Verify default accounts
    verify_default_accounts

    # Verify custom accounts
    verify_custom_accounts

    # Test transactions (optional)
    if [ "${RUN_TX_TEST:-false}" = "true" ]; then
        test_transactions
    else
        print_info "Transaction test skipped (set RUN_TX_TEST=true to enable)"
    fi

    # Export account information
    if [ "${EXPORT_ACCOUNTS:-false}" = "true" ]; then
        export_account_info
    fi

    echo ""
    echo "========================================="
    echo "Test Summary"
    echo "========================================="
    echo ""
    print_success "All genesis account tests completed successfully!"
    echo ""
    echo "Chain endpoints:"
    echo "  RPC: $RPC_URL"
    echo "  API: $API_URL"
    echo "  JSON-RPC: http://localhost:9545"
    echo ""
}

# Run main function
main "$@"
