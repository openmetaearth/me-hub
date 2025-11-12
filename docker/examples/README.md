# Genesis Accounts Configuration Examples

This directory contains examples and scripts demonstrating how to build ME-Chain private network Docker images with custom genesis accounts.

## Files

### `build_with_accounts.sh`

Interactive script showing various examples of how to configure genesis accounts.

**Usage:**

```bash
# Run interactively
./docker/examples/build_with_accounts.sh

# Show specific example
./docker/examples/build_with_accounts.sh 1

# Show all examples
./docker/examples/build_with_accounts.sh 0
```

## Quick Reference

### Simple Format

Create accounts with comma-separated values:

```bash
# Format: "name1:amount1,name2:amount2,..."
make docker-private-net GENESIS_ACCOUNTS="alice:1000000000000000000000umec,bob:500000000000000000000umec"
```

### JSON Format

Create accounts with JSON for more control (supports mnemonics):

```bash
# Format: [{"name":"...","amount":"...","mnemonic":"..."}]
make docker-private-net GENESIS_ACCOUNTS_JSON='[
  {"name":"alice","amount":"2000000000000000000000umec"},
  {"name":"bob","amount":"1000000000000000000000umec"}
]'
```

### Runtime Configuration

Add accounts when starting the container:

```bash
docker run -d \
  -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
  -e GENESIS_ACCOUNTS="test1:1000000000000000000000umec,test2:500000000000000000000umec" \
  --name mechain-private-net \
  me-hub/private-net:latest
```

## Common Scenarios

### 1. Transfer Testing

```bash
make docker-private-net \
  GENESIS_ACCOUNTS="sender:10000000000000000000000umec,receiver:1000000000000000000000umec"
```

### 2. Smart Contract Deployment

```bash
make docker-private-net \
  GENESIS_ACCOUNTS="deployer:100000000000000000000000umec,user1:5000000000000000000000umec,user2:5000000000000000000000umec"
```

### 3. Multi-Signature Testing

```bash
make docker-private-net \
  GENESIS_ACCOUNTS="signer1:5000000000000000000000umec,signer2:5000000000000000000000umec,signer3:5000000000000000000000umec"
```

### 4. Integration Testing (Multiple Users)

```bash
make docker-private-net \
  GENESIS_ACCOUNTS="user1,user2,user3,user4,user5"
```

Each user gets the default amount of 1,000 MEC.

## Amount Conversion

ME-Chain uses `umec` as the base unit:

| MEC | umec | Note |
|-----|------|------|
| 1 MEC | 1000000000000000000umec | 10^18 |
| 10 MEC | 10000000000000000000umec | 10^19 |
| 100 MEC | 100000000000000000000umec | 10^20 |
| 1,000 MEC | 1000000000000000000000umec | 10^21 |
| 10,000 MEC | 10000000000000000000000umec | 10^22 |
| 100,000 MEC | 100000000000000000000000umec | 10^23 |
| 1,000,000 MEC | 1000000000000000000000000umec | 10^24 |

## Verification

After building and starting the network, verify your accounts:

```bash
# List all accounts
docker exec mechain-private-net med keys list --keyring-backend test

# Check specific account balance
docker exec mechain-private-net med keys show alice --keyring-backend test -a
docker exec mechain-private-net med query bank balances <alice_address>

# Use automated test script
./docker/scripts/test_genesis_accounts.sh
```

## GitHub Actions Example

```yaml
name: Build Test Network

on:
  workflow_dispatch:
    inputs:
      accounts:
        description: 'Genesis accounts'
        default: 'alice:10000000000000000000000umec,bob:5000000000000000000000umec'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      
      - name: Build private network
        run: |
          make docker-private-net \
            GENESIS_ACCOUNTS="${{ github.event.inputs.accounts }}"
      
      - name: Start network
        run: |
          docker run -d \
            -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
            --name mechain-test \
            me-hub/private-net:latest
      
      - name: Run tests
        run: |
          sleep 15
          ./docker/scripts/test_genesis_accounts.sh
```

## Documentation

For complete documentation, see:

- [GENESIS_ACCOUNTS.md](../GENESIS_ACCOUNTS.md) - Full configuration guide
- [README.md](../README.md) - Main Docker documentation
- [QUICKSTART.md](../QUICKSTART.md) - Quick start guide

## Support

If you encounter issues:

1. Check the [troubleshooting section](../GENESIS_ACCOUNTS.md#故障排查) in GENESIS_ACCOUNTS.md
2. Run the test script: `./docker/scripts/test_genesis_accounts.sh`
3. View container logs: `docker logs mechain-private-net`
4. Submit an issue to the repository