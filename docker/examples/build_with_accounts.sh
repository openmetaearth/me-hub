#!/usr/bin/env bash
# Examples for building ME-Chain private network with custom genesis accounts

set -e

echo "========================================="
echo "ME-Chain Private Network Build Examples"
echo "========================================="
echo ""

# Example 1: Simple format - Basic accounts
example_1() {
    echo "Example 1: Creating basic test accounts with simple format"
    echo "-----------------------------------------------------------"
    echo "Command:"
    echo "  make docker-private-net GENESIS_ACCOUNTS=\"alice:1000000000000000000000umec,bob:1000000000000000000000umec\""
    echo ""
    echo "This creates:"
    echo "  - alice: 1,000 MEC"
    echo "  - bob: 1,000 MEC"
    echo ""
}

# Example 2: Simple format - Different amounts
example_2() {
    echo "Example 2: Creating accounts with different amounts"
    echo "----------------------------------------------------"
    echo "Command:"
    echo "  make docker-private-net GENESIS_ACCOUNTS=\"whale:100000000000000000000000umec,regular:10000000000000000000000umec,small:1000000000000000000000umec\""
    echo ""
    echo "This creates:"
    echo "  - whale: 100,000 MEC"
    echo "  - regular: 10,000 MEC"
    echo "  - small: 1,000 MEC"
    echo ""
}

# Example 3: Simple format - Using default amounts
example_3() {
    echo "Example 3: Creating accounts with default amounts"
    echo "--------------------------------------------------"
    echo "Command:"
    echo "  make docker-private-net GENESIS_ACCOUNTS=\"user1,user2,user3,user4\""
    echo ""
    echo "This creates 4 accounts, each with default 1,000 MEC"
    echo ""
}

# Example 4: JSON format - Basic usage
example_4() {
    echo "Example 4: Using JSON format for more control"
    echo "----------------------------------------------"
    echo "Command:"
    cat << 'EOF'
  make docker-private-net GENESIS_ACCOUNTS_JSON='[
    {"name":"alice","amount":"2000000000000000000000umec"},
    {"name":"bob","amount":"1000000000000000000000umec"}
  ]'
EOF
    echo ""
    echo "This creates:"
    echo "  - alice: 2,000 MEC"
    echo "  - bob: 1,000 MEC"
    echo ""
}

# Example 5: JSON format - With mnemonic
example_5() {
    echo "Example 5: Using JSON format with custom mnemonic (reproducible accounts)"
    echo "-------------------------------------------------------------------------"
    echo "Command:"
    cat << 'EOF'
  make docker-private-net GENESIS_ACCOUNTS_JSON='[
    {
      "name":"test_account",
      "amount":"5000000000000000000000umec",
      "mnemonic":"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
    }
  ]'
EOF
    echo ""
    echo "This creates a reproducible account using the specified mnemonic"
    echo "  - test_account: 5,000 MEC"
    echo ""
}

# Example 6: Docker run with environment variables
example_6() {
    echo "Example 6: Starting container with custom accounts at runtime"
    echo "--------------------------------------------------------------"
    echo "Command:"
    cat << 'EOF'
  docker run -d \
    -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
    -e GENESIS_ACCOUNTS="test1:1000000000000000000000umec,test2:500000000000000000000umec" \
    --name mechain-private-net \
    me-hub/private-net:latest
EOF
    echo ""
    echo "This starts a container with:"
    echo "  - test1: 1,000 MEC"
    echo "  - test2: 500 MEC"
    echo ""
}

# Example 7: Docker Compose
example_7() {
    echo "Example 7: Using Docker Compose with custom accounts"
    echo "-----------------------------------------------------"
    echo "Add to docker-compose.yml:"
    cat << 'EOF'
  services:
    mechain-private-net:
      environment:
        - GENESIS_ACCOUNTS=deployer:50000000000000000000000umec,user1:5000000000000000000000umec,user2:5000000000000000000000umec
EOF
    echo ""
    echo "This creates:"
    echo "  - deployer: 50,000 MEC (for contract deployment)"
    echo "  - user1: 5,000 MEC"
    echo "  - user2: 5,000 MEC"
    echo ""
}

# Example 8: CI/CD - GitHub Actions
example_8() {
    echo "Example 8: GitHub Actions workflow"
    echo "-----------------------------------"
    echo "Add to .github/workflows/test.yml:"
    cat << 'EOF'
  - name: Build test network
    run: |
      make docker-private-net \
        GENESIS_ACCOUNTS="tester1:10000000000000000000000umec,tester2:10000000000000000000000umec"

  - name: Start network
    run: |
      docker run -d \
        -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
        --name mechain-test \
        me-hub/private-net:latest

  - name: Run tests
    run: |
      sleep 15
      ./run_integration_tests.sh
EOF
    echo ""
}

# Example 9: Testing scenarios
example_9() {
    echo "Example 9: Common testing scenarios"
    echo "------------------------------------"
    echo ""
    echo "A) Transfer testing:"
    echo "  GENESIS_ACCOUNTS=\"sender:10000000000000000000000umec,receiver:1000000000000000000000umec\""
    echo ""
    echo "B) Multi-signature testing:"
    echo "  GENESIS_ACCOUNTS=\"signer1:5000000000000000000000umec,signer2:5000000000000000000000umec,signer3:5000000000000000000000umec\""
    echo ""
    echo "C) Contract deployment testing:"
    echo "  GENESIS_ACCOUNTS=\"deployer:100000000000000000000000umec,user1:1000000000000000000000umec,user2:1000000000000000000000umec\""
    echo ""
    echo "D) Stress testing (many accounts):"
    echo "  GENESIS_ACCOUNTS=\"user1,user2,user3,user4,user5,user6,user7,user8,user9,user10\""
    echo ""
}

# Example 10: Amount reference
example_10() {
    echo "Example 10: Amount calculation reference"
    echo "-----------------------------------------"
    echo ""
    echo "ME-Chain uses umec as base unit (1 MEC = 10^18 umec)"
    echo ""
    printf "%-15s | %-30s | %s\n" "MEC Amount" "umec Amount" "Scientific Notation"
    echo "----------------|--------------------------------|--------------------"
    printf "%-15s | %-30s | %s\n" "1 MEC" "1,000,000,000,000,000,000" "1×10^18"
    printf "%-15s | %-30s | %s\n" "10 MEC" "10,000,000,000,000,000,000" "1×10^19"
    printf "%-15s | %-30s | %s\n" "100 MEC" "100,000,000,000,000,000,000" "1×10^20"
    printf "%-15s | %-30s | %s\n" "1,000 MEC" "1,000,000,000,000,000,000,000" "1×10^21"
    printf "%-15s | %-30s | %s\n" "10,000 MEC" "10,000,000,000,000,000,000,000" "1×10^22"
    printf "%-15s | %-30s | %s\n" "100,000 MEC" "100,000,000,000,000,000,000,000" "1×10^23"
    printf "%-15s | %-30s | %s\n" "1,000,000 MEC" "1,000,000,000,000,000,000,000,000" "1×10^24"
    echo ""
    echo "Common amounts for copy-paste:"
    echo "  1,000 MEC:     1000000000000000000000umec"
    echo "  10,000 MEC:    10000000000000000000000umec"
    echo "  100,000 MEC:   100000000000000000000000umec"
    echo "  1,000,000 MEC: 1000000000000000000000000umec"
    echo ""
}

# Main menu
main() {
    echo "Select an example to view:"
    echo ""
    echo "  1) Simple format - Basic accounts"
    echo "  2) Simple format - Different amounts"
    echo "  3) Simple format - Default amounts"
    echo "  4) JSON format - Basic usage"
    echo "  5) JSON format - With mnemonic"
    echo "  6) Docker run with environment variables"
    echo "  7) Docker Compose configuration"
    echo "  8) GitHub Actions workflow"
    echo "  9) Common testing scenarios"
    echo " 10) Amount calculation reference"
    echo "  0) Show all examples"
    echo ""

    if [ $# -eq 0 ]; then
        read -p "Enter number (0-10): " choice
    else
        choice=$1
    fi

    echo ""

    case $choice in
        1) example_1 ;;
        2) example_2 ;;
        3) example_3 ;;
        4) example_4 ;;
        5) example_5 ;;
        6) example_6 ;;
        7) example_7 ;;
        8) example_8 ;;
        9) example_9 ;;
        10) example_10 ;;
        0)
            example_1
            example_2
            example_3
            example_4
            example_5
            example_6
            example_7
            example_8
            example_9
            example_10
            ;;
        *)
            echo "Invalid choice. Use 0-10."
            exit 1
            ;;
    esac

    echo ""
    echo "For complete documentation, see: docker/GENESIS_ACCOUNTS.md"
    echo ""
}

main "$@"
