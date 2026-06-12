# ME-Chain Docker Private Network Quick Start

> One-click launch of a pre-compiled and pre-initialized single-node private test network

## 📋 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Testing & Validation](#testing--validation)
- [FAQ](#faq)

## Overview

This Docker solution provides a fully configured ME-Chain single-node private network, eliminating the need to manually execute `make build` and `setup_local.sh`.

**Comparison**:

```
Traditional:  make build (5 min) + setup_local.sh (2 min) + interactive config = ~10 min
Docker:       make docker-private-net (5 min) + docker run (30 sec) = ~6 min (no interaction)
```

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+ (recommended)
- 4GB+ RAM
- 10GB+ disk space
- Command-line tools: `curl`, `jq` (for testing)

## Quick Start

### Method 1: Using Makefile (Recommended)

```bash
# 1. Build the image
make docker-private-net

# 2. Start the network
make docker-private-net-start

# 3. View logs
docker compose -f docker/docker-compose.yml logs -f

# 4. Run tests
make docker-private-net-test

# 5. Stop the network
make docker-private-net-stop
```

### Method 2: Using Docker Commands

```bash
# 1. Build the image
make docker-private-net

# 2. Start network (with persistent data)
docker run -d \
  -p 36657:36657 \
  -p 1318:1318 \
  -p 9545:9545 \
  -p 8090:8090 \
  -v mechain-data:/root/.mechain \
  --name mechain-private-net \
  me-hub/private-net:latest

# 3. View logs
docker logs -f mechain-private-net

# 4. Stop and clean up
docker stop mechain-private-net
docker rm mechain-private-net
docker volume rm mechain-data
```

### Method 3: Temporary Test Network

```bash
# Start temporary network (data lost on restart)
docker run -d \
  -p 36657:36657 \
  -p 1318:1318 \
  -p 9545:9545 \
  -p 8090:8090 \
  --name mechain-private-net \
  me-hub/private-net:latest
```

## Usage

### Access Endpoints

| Service         | Port  | URL                    |
| --------------- | ----- | ---------------------- |
| **RPC**         | 36657 | http://localhost:36657 |
| **REST API**    | 1318  | http://localhost:1318  |
| **JSON-RPC**    | 9545  | http://localhost:9545  |
| **JSON-RPC WS** | 9546  | ws://localhost:9546    |
| **gRPC**        | 8090  | localhost:8090         |
| **gRPC-web**    | 8091  | localhost:8091         |

### Quick Tests

```bash
# Check node status
curl http://localhost:36657/status | jq

# Check API
curl http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info | jq

# Check JSON-RPC
curl -X POST http://localhost:9545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq
```

### Pre-configured Accounts

The private network includes the following pre-configured accounts:

#### 1. Main Account (global_dao)

```bash
# Mnemonic
curtain hat remain song receive tower stereo hope frog cheap brown plate 
raccoon post reflect wool sail salmon game salon group glimpse adult shift

# View address
docker exec mechain-private-net med keys show global_dao -a --keyring-backend test

# Check balance
docker exec mechain-private-net med query bank balances \
  $(docker exec mechain-private-net med keys show global_dao -a --keyring-backend test)
```

**Initial Balance**: 1,000,000 MEC (1e24 umec)

#### 2. Other Accounts

- **pools**: AMM pool account (1,000,000 MEC)
- **user**: Test user account (1,000,000 MEC)
- **sequencer**: Sequencer account (0.0001 MEC)

### Common Operations

```bash
# List all accounts
docker exec mechain-private-net med keys list --keyring-backend test

# Query account balance
ADDR=$(docker exec mechain-private-net med keys show global_dao -a --keyring-backend test)
docker exec mechain-private-net med query bank balances $ADDR

# Send tokens
docker exec mechain-private-net med tx bank send \
  <from_address> \
  <to_address> \
  1000000umec \
  --chain-id mechain_202404-1 \
  --keyring-backend test \
  --yes

# View block height
docker exec mechain-private-net med status | jq .SyncInfo.latest_block_height

# View validator information
docker exec mechain-private-net med query staking validators

# Enter container interactive shell
docker exec -it mechain-private-net sh
```

## Testing & Validation

### Automated Tests

```bash
# Run complete test suite
make docker-private-net-test

# Or run test script directly
./docker/scripts/test_private_net.sh
```

Test coverage includes:
- ✅ Container running status
- ✅ RPC endpoint accessibility
- ✅ REST API endpoint
- ✅ JSON-RPC endpoint
- ✅ Pre-configured accounts
- ✅ Validator status
- ✅ Genesis configuration
- ✅ Block production

### Interactive Demo

```bash
# Run interactive demo script
./docker/scripts/private_net_demo.sh

# Or run all demos
./docker/scripts/private_net_demo.sh --all
```

Demo includes:
1. List all accounts
2. Check account balances
3. Send token transactions
4. Query chain state
5. Query validator information
6. REST API usage
7. JSON-RPC usage
8. Create new account
9. Query chain parameters
10. Monitor block production

## FAQ

### Q1: How to reset chain state?

```bash
# Using Docker Compose
docker compose -f docker/docker-compose.yml down -v
docker compose -f docker/docker-compose.yml up -d

# Using Docker commands
docker stop mechain-private-net
docker rm mechain-private-net
docker volume rm mechain-data
docker run -d -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
  -v mechain-data:/root/.mechain --name mechain-private-net me-hub/private-net:latest
```

### Q2: What about port conflicts?

```bash
# Use different port mappings
docker run -d \
  -p 26657:36657 \
  -p 2318:1318 \
  -p 8545:9545 \
  -p 9090:8090 \
  --name mechain-private-net \
  me-hub/private-net:latest
```

### Q3: How to view detailed logs?

```bash
# Real-time logs
docker logs -f mechain-private-net

# Last 100 lines
docker logs --tail 100 mechain-private-net

# Specific time range
docker logs --since 10m mechain-private-net
docker logs --until 2024-01-01T00:00:00 mechain-private-net
```

### Q4: How to backup data?

```bash
# Backup
docker run --rm \
  -v mechain-data:/data \
  -v $(pwd):/backup \
  ubuntu tar czf /backup/mechain-backup-$(date +%Y%m%d).tar.gz /data

# Restore
docker run --rm \
  -v mechain-data:/data \
  -v $(pwd):/backup \
  ubuntu tar xzf /backup/mechain-backup-20240101.tar.gz -C /
```

### Q5: How to customize configuration?

```bash
# Through environment variables
docker run -d \
  -e CHAIN_ID="custom_chain_100-1" \
  -e MONIKER_NAME="my-node" \
  -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
  --name mechain-private-net \
  me-hub/private-net:latest

# Or modify environment variables in docker/docker-compose.yml
```

### Q6: Container fails to start?

```bash
# Check logs
docker logs mechain-private-net

# Check port usage
netstat -tuln | grep -E '36657|1318|9545|8090'

# Check resources
docker stats mechain-private-net

# Remove and recreate
docker rm -f mechain-private-net
docker volume rm mechain-data
make docker-private-net-start
```

### Q7: How to connect to the private network?

**JavaScript/TypeScript:**
```javascript
const { SigningStargateClient } = require("@cosmjs/stargate");

const rpcEndpoint = "http://localhost:36657";
const client = await SigningStargateClient.connect(rpcEndpoint);
```

**Python:**
```python
from cosmpy.aerial.client import LedgerClient, NetworkConfig

network = NetworkConfig(
    chain_id="mechain_202404-1",
    url="grpc+http://localhost:8090",
    fee_minimum_gas_price=1,
    fee_denomination="umec",
    staking_denomination="umec",
)
client = LedgerClient(network)
```

**Go:**
```go
import "github.com/cosmos/cosmos-sdk/client"

clientCtx := client.Context{}.
    WithNodeURI("tcp://localhost:36657").
    WithChainID("mechain_202404-1")
```

### Q8: How to use in CI/CD?

**GitHub Actions Example:**
```yaml
name: Integration Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Start ME-Chain private network
        run: |
          make docker-private-net
          make docker-private-net-start
          
      - name: Wait for chain to start
        run: |
          timeout 60 bash -c 'until curl -s http://localhost:36657/status > /dev/null; do sleep 2; done'
          
      - name: Run tests
        run: make docker-private-net-test
        
      - name: Cleanup
        if: always()
        run: make docker-private-net-stop
```

## Performance Tuning

```bash
# Limit resource usage
docker run -d \
  --memory="4g" \
  --cpus="2" \
  --memory-swap="4g" \
  -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
  --name mechain-private-net \
  me-hub/private-net:latest
```

## Project Files

```
me-hub/
├── docker/
│   ├── Dockerfile                      # Docker image definition
│   ├── docker-compose.yml              # Docker Compose configuration
│   ├── README.md                       # Quick guide
│   ├── QUICKSTART.md                   # This file
│   └── scripts/
│       ├── setup_local_docker.sh       # Non-interactive init script
│       ├── start_private_net.sh        # Container startup script
│       ├── test_private_net.sh         # Automated test script
│       └── private_net_demo.sh         # Interactive demo script
└── scripts/
    └── src/
        └── genesis_config_commands.sh  # Genesis config commands
```

## Advanced Usage

### Docker Network Integration

```bash
# Create custom network
docker network create mechain-network

# Start private network
docker run -d \
  --network mechain-network \
  --name mechain-private-net \
  me-hub/private-net:latest

# Use in other containers
# RPC: http://mechain-private-net:36657
# API: http://mechain-private-net:1318
```

### Multi-Instance Deployment

```bash
# Instance 1
docker run -d \
  -p 36657:36657 -p 1318:1318 \
  --name mechain-node1 \
  me-hub/private-net:latest

# Instance 2 (using different ports)
docker run -d \
  -p 36667:36657 -p 1328:1318 \
  --name mechain-node2 \
  me-hub/private-net:latest
```

## Resources

- [Full Documentation](README.md)
- [ME-Chain Main Documentation](../README.md)
- [Docker Official Documentation](https://docs.docker.com/)
- [Cosmos SDK Documentation](https://docs.cosmos.network/)

## Support

Having issues?
1. Check [FAQ](#faq)
2. Run `make docker-private-net-test` for diagnostics
3. View container logs: `docker logs mechain-private-net`
4. Submit an issue

---

**Happy Building!** 🚀