# ME-Chain Docker Quick Start Guide

This guide covers two types of Docker images for ME-Chain:

1. **Private Network Image** - Pre-initialized single-node test network
2. **Release Image** - Minimal production-ready binary image

## 📦 Image Types

### Private Network Image

A pre-compiled and pre-initialized single-node private network with test accounts. Perfect for:
- Local development
- Testing and CI/CD
- Quick prototyping

**Build**: `make docker-private-net`  
**Tag**: `me-hub/private-net:latest`

### Release Image

A minimal runtime image containing only the compiled `med` binary. Perfect for:
- Production deployments
- Custom chain configurations
- Base image for other containers

**Build**: `make docker-release`  
**Tag**: `me-hub/release:latest`

---

# 🚀 Private Network Image

## 🚀 Quick Start

### Option A: Using Makefile (Recommended)

```bash
# 1. Build Docker image
make docker-private-net

# 2. Start private network
make docker-private-net-start

# 3. Test network functionality
make docker-private-net-test

# 4. Stop network
make docker-private-net-stop
```

### Option B: Using Docker Compose

```bash
# Start
docker compose -f docker/docker-compose.yml up -d

# View logs
docker compose -f docker/docker-compose.yml logs -f

# Stop
docker compose -f docker/docker-compose.yml down
```

### Option C: Using Docker Commands

```bash
# Start with persistent data
docker run -d \
  -p 36657:36657 \
  -p 1318:1318 \
  -p 9545:9545 \
  -p 8090:8090 \
  -v mechain-data:/root/.mechain \
  --name mechain-private-net \
  me-hub/private-net:latest

# View logs
docker logs -f mechain-private-net

# Stop and clean up
docker stop mechain-private-net
docker rm mechain-private-net
```

## 📡 Access Endpoints

After successful startup, you can access the network through the following endpoints:

| Service                | Port  | URL                    |
| ---------------------- | ----- | ---------------------- |
| **RPC**                | 36657 | http://localhost:36657 |
| **REST API**           | 1318  | http://localhost:1318  |
| **JSON-RPC**           | 9545  | http://localhost:9545  |
| **JSON-RPC WebSocket** | 9546  | ws://localhost:9546    |
| **gRPC**               | 8090  | localhost:8090         |
| **gRPC-web**           | 8091  | localhost:8091         |
| **P2P**                | 36656 | -                      |

## 🔑 Pre-configured Accounts

The private network includes the following pre-configured test accounts (using test keyring):

### Main Account (global_dao)

```bash
# Mnemonic
curtain hat remain song receive tower stereo hope frog cheap brown plate 
raccoon post reflect wool sail salmon game salon group glimpse adult shift

# View address
docker exec mechain-private-net med keys show global_dao -a --keyring-backend test

# Expected address
me139mq752delxv78jvtmwxhasyrycufsvr0mue6u

# Check balance
docker exec mechain-private-net med query bank balances \
  me139mq752delxv78jvtmwxhasyrycufsvr0mue6u
```

**Initial Balance**: 1,000,000,000,000,000,000,000,000 umec (1,000,000 MEC)

### Other Pre-configured Accounts

| Account Name   | Purpose           | Initial Balance |
| -------------- | ----------------- | --------------- |
| **global_dao** | Main account/DAO  | 1,000,000 MEC   |
| **pools**      | AMM pool account  | 1,000,000 MEC   |
| **user**       | Test user account | 1,000,000 MEC   |
| **sequencer**  | Sequencer account | 0.0001 MEC      |

## 🎯 Creating Custom Genesis Accounts

You can create additional genesis accounts at **build time** or **runtime** using environment variables. This is especially useful for CI/CD pipelines and automated testing.

📖 **See [GENESIS_ACCOUNTS.md](GENESIS_ACCOUNTS.md) for detailed documentation.**

### Quick Examples

#### Build Time Configuration

```bash
# Simple format: Create multiple accounts with specified amounts
make docker-private-net \
  GENESIS_ACCOUNTS="alice:1000000000000000000000umec,bob:1000000000000000000000umec"

# JSON format: More control (can specify mnemonics)
make docker-private-net \
  GENESIS_ACCOUNTS_JSON='[{"name":"alice","amount":"2000000000000000000000umec"}]'
```

#### Runtime Configuration

```bash
# Add accounts when starting the container
docker run -d \
  -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
  -e GENESIS_ACCOUNTS="test1:1000000000000000000000umec,test2:500000000000000000000umec" \
  --name mechain-private-net \
  me-hub/private-net:latest
```

#### Docker Compose Configuration

Edit `docker/docker-compose.yml`:

```yaml
services:
  mechain-private-net:
    environment:
      - GENESIS_ACCOUNTS=alice:1000000000000000000000umec,bob:1000000000000000000000umec
```

### Verify Custom Accounts

```bash
# List all accounts
docker exec mechain-private-net med keys list --keyring-backend test

# Check balance of custom account
docker exec mechain-private-net med keys show alice --keyring-backend test -a
docker exec mechain-private-net med query bank balances <alice_address>

# Or use the test script
chmod +x docker/scripts/test_genesis_accounts.sh
./docker/scripts/test_genesis_accounts.sh
```

### GitHub Actions Integration

Create custom test networks automatically:

```yaml
- name: Build test network with custom accounts
  run: |
    make docker-private-net \
      GENESIS_ACCOUNTS="alice:10000000000000000000000umec,bob:5000000000000000000000umec"
```

For complete documentation including:
- JSON format with mnemonic support
- Amount calculation reference
- GitHub Actions workflow examples
- Troubleshooting guide

**👉 See [GENESIS_ACCOUNTS.md](GENESIS_ACCOUNTS.md)**

## 🧪 Test Connection

```bash
# 1. Check node status
curl http://localhost:36657/status

# 2. View Chain ID and block height
curl -s http://localhost:36657/status | jq '.result.node_info.network, .result.sync_info.latest_block_height'

# 3. Check REST API
curl http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info

# 4. Check JSON-RPC
curl -X POST http://localhost:9545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# 5. Run complete test suite
make docker-private-net-test
```

## 🛠️ Common Commands

### Account Management

```bash
# List all accounts
docker exec mechain-private-net med keys list --keyring-backend test

# Check specific account balance
docker exec mechain-private-net med query bank balances \
  me139mq752delxv78jvtmwxhasyrycufsvr0mue6u

# View account details
docker exec mechain-private-net med keys show global_dao --keyring-backend test
```

### Transaction Operations

```bash
# Send tokens
docker exec mechain-private-net med tx bank send \
  me139mq752delxv78jvtmwxhasyrycufsvr0mue6u \
  <target_address> \
  1000000umec \
  --chain-id mechain_202404-1 \
  --keyring-backend test \
  --yes

# Query transaction
docker exec mechain-private-net med query tx <TX_HASH>
```

### Validator Operations

```bash
# View validator list
docker exec mechain-private-net med query staking validators

# View validator details
docker exec mechain-private-net med query staking validator \
  mevaloper139mq752delxv78jvtmwxhasyrycufsvr707ate
```

### Chain State Queries

```bash
# Query block height
docker exec mechain-private-net med status | jq .SyncInfo.latest_block_height

# Query specific block
docker exec mechain-private-net med query block <HEIGHT>

# View chain parameters
docker exec mechain-private-net med query staking params
```

### Container Management

```bash
# Enter container shell
docker exec -it mechain-private-net sh

# View real-time logs
docker logs -f mechain-private-net

# View recent logs
docker logs --tail 100 mechain-private-net

# Check container status
docker ps | grep mechain-private-net

# Check container resource usage
docker stats mechain-private-net --no-stream
```

## ⚙️ Environment Variables

Customize chain configuration through environment variables:

| Variable                | Default Value    | Description                             |
| ----------------------- | ---------------- | --------------------------------------- |
| `CHAIN_ID`              | mechain_202404-1 | Chain ID                                |
| `MONIKER_NAME`          | local            | Node name                               |
| `KEY_NAME`              | global_dao       | Main account name                       |
| `TZ`                    | Asia/Shanghai    | Timezone setting                        |
| `GENESIS_ACCOUNTS`      | ""               | Custom genesis accounts (simple format) |
| `GENESIS_ACCOUNTS_JSON` | ""               | Custom genesis accounts (JSON format)   |

**Example**:
```bash
docker run -d \
  -e CHAIN_ID="mechain_202404-1" \
  -e MONIKER_NAME="my-node" \
  -e KEY_NAME="global_dao" \
  -e TZ="UTC" \
  -p 36657:36657 \
  -p 1318:1318 \
  -p 9545:9545 \
  -p 8090:8090 \
  --name mechain-private-net \
  me-hub/private-net:latest
```

---

# 🎯 Release Image

The release image is a minimal runtime image without any chain initialization. It only contains the `med` binary and required libraries.

## Building Release Image

```bash
# Build with default tag (latest)
make docker-release

# Build with specific tag
make docker-release TAG=v1.0.0
```

## Using Release Image

### Check Version

```bash
# Run from Docker image
docker run --rm me-hub/release:latest version

# Check detailed version
docker run --rm me-hub/release:latest version --long
```

### View Help

```bash
# Show all available commands
docker run --rm me-hub/release:latest --help

# Show help for specific command
docker run --rm me-hub/release:latest init --help
```

### Initialize Custom Chain

```bash
# Create a volume for chain data
docker volume create mychain-data

# Initialize a new chain
docker run --rm \
  -v mychain-data:/root/.mechain \
  me-hub/release:latest \
  init mynode --chain-id mychain_100-1

# Add keys
docker run --rm \
  -v mychain-data:/root/.mechain \
  me-hub/release:latest \
  keys add mykey --keyring-backend test

# Run your custom chain
docker run -d \
  -v mychain-data:/root/.mechain \
  -p 36657:36657 -p 1318:1318 -p 9545:9545 \
  --name mychain \
  --entrypoint med \
  me-hub/release:latest \
  start
```

### Use as Base Image

Create your own Dockerfile:

```dockerfile
FROM me-hub/release:latest

# Copy your custom genesis file
COPY genesis.json /root/.mechain/config/genesis.json

# Copy your custom configuration
COPY config.toml /root/.mechain/config/config.toml
COPY app.toml /root/.mechain/config/app.toml

# Set custom entrypoint
COPY start-chain.sh /scripts/start-chain.sh
RUN chmod +x /scripts/start-chain.sh

ENTRYPOINT ["/scripts/start-chain.sh"]
```

### Interactive Shell

```bash
# Enter container with shell
docker run -it --rm \
  -v mychain-data:/root/.mechain \
  --entrypoint sh \
  me-hub/release:latest

# Now you can run any med commands
# med init mynode --chain-id test
# med keys add mykey
# med start
```

## GitHub Actions Workflow

The release image is automatically built and pushed when a tag is pushed:

```yaml
# .github/workflows/build-push-release.yml
name: Build and Push Release Docker Image

on:
  push:
    tags:
      - "v*"
```

This will:
1. Build the release Docker image
2. Tag it with the git tag (e.g., `v1.0.0`)
3. Push to Harbor registry
4. Verify by running `med version`

### Pull from Registry

```bash
# Pull specific version
docker pull your-registry/openmetaearth/me_hub:v1.0.0

# Pull latest
docker pull your-registry/openmetaearth/me_hub:latest

# Check version
docker run --rm your-registry/openmetaearth/me_hub:v1.0.0 version
```

## Comparison: Release vs Private Network

| Feature            | Release Image     | Private Network Image      |
| ------------------ | ----------------- | -------------------------- |
| **Size**           | ~200MB            | ~423MB                     |
| **Includes**       | Binary only       | Binary + initialized chain |
| **Pre-configured** | No                | Yes (test accounts)        |
| **Use Case**       | Production/Custom | Development/Testing        |
| **Startup**        | Requires init     | Ready to run               |
| **Entrypoint**     | `med`             | Startup script             |
| **Default CMD**    | `version`         | Starts chain               |

## 🔄 Data Persistence

### Using Volume for Data Persistence

```bash
# Start with specified volume
docker run -d \
  -v mechain-data:/root/.mechain \
  -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
  --name mechain-private-net \
  me-hub/private-net:latest

# View volume
docker volume ls | grep mechain

# Inspect volume contents
docker volume inspect mechain-data
```

### Reset Network (Clear All Data)

```bash
# Using Makefile
make docker-private-net-stop
docker volume rm docker_mechain-data

# Or using Docker Compose
docker compose -f docker/docker-compose.yml down -v

# Or using Docker commands
docker stop mechain-private-net
docker rm mechain-private-net
docker volume rm mechain-data

# Restart
make docker-private-net-start
```

## 📊 Performance & Resources

### System Requirements

- **CPU**: 2+ cores
- **Memory**: Minimum 2GB, Recommended 4GB
- **Disk**: At least 1GB available space
- **Network**: Ports 36656-36657, 1318, 8090-8091, 9545-9546 available

### Resource Usage

- **Image Size**: ~423MB
- **Runtime Memory**: ~200MB
- **CPU Usage**: < 5% (idle)
- **Startup Time**: < 10 seconds
- **Block Time**: ~5 seconds/block

## 🔍 Health Check

Docker container has built-in health check functionality:

```bash
# View health status
docker inspect mechain-private-net --format '{{.State.Health.Status}}'

# View health check logs
docker inspect mechain-private-net --format '{{json .State.Health}}' | jq
```

Health check configuration:
- **Command**: `med status`
- **Interval**: 30 seconds
- **Timeout**: 10 seconds
- **Retries**: 5 times
- **Start Period**: 60 seconds

## 🐛 Troubleshooting

### Container Won't Start

```bash
# Check logs
docker logs mechain-private-net

# Check port usage
netstat -tuln | grep -E "36657|1318|9545|8090"

# Check disk space
df -h
docker system df
```

### Port Conflicts

Use custom port mapping:

```bash
docker run -d \
  -p 26657:36657 \
  -p 2318:1318 \
  -p 8545:9545 \
  -p 9090:8090 \
  --name mechain-private-net \
  me-hub/private-net:latest
```

### Chain Not Producing Blocks

```bash
# Check validator status
docker exec mechain-private-net med query staking validators

# Check node status
docker exec mechain-private-net med status

# View detailed logs
docker logs -f mechain-private-net | grep -E "ERROR|WARN"
```

### API Not Accessible

```bash
# Check configuration file
docker exec mechain-private-net cat /root/.mechain/config/app.toml | grep -A 5 "\[api\]"

# Check port mapping
docker port mechain-private-net

# Test internal container connection
docker exec mechain-private-net curl http://localhost:1318/cosmos/base/tendermint/v1beta1/node_info
```

## 📚 Related Documentation

- [Genesis Accounts Configuration](GENESIS_ACCOUNTS.md) - **Detailed guide for creating custom accounts**
- [Quick Start Guide](QUICKSTART.md)
- [ME-Chain Main Documentation](../README.md)
- [Build Release Workflow](../.github/workflows/build-push-release.yml)
- [Build Private Network Workflow](../.github/workflows/build-push-private-net.yml)

## 📁 Related Files

```
docker/
├── Dockerfile                        # Private network image definition
├── Dockerfile.release                # Release image definition (minimal)
├── docker-compose.yml                # Docker Compose configuration
├── README.md                         # This document
├── QUICKSTART.md                     # Quick start guide
├── GENESIS_ACCOUNTS.md               # Genesis accounts configuration guide
└── scripts/
    ├── setup_local_docker.sh         # Chain initialization script
    ├── start_private_net.sh          # Container startup script
    ├── test_private_net.sh           # Automated test script
    └── test_genesis_accounts.sh      # Genesis accounts verification script

.github/workflows/
├── build-push-private-net.yml        # CI/CD for private network image
└── build-push-release.yml            # CI/CD for release image
```

## 💡 Best Practices

### Development & Testing

```bash
# 1. Start private network
make docker-private-net-start

# 2. Wait for initialization (about 10 seconds)
sleep 10

# 3. Run tests
make docker-private-net-test

# 4. Clean up after development
make docker-private-net-stop
```

### CI/CD Integration

```bash
#!/bin/bash
# Use in CI pipeline

# Start test network
make docker-private-net-start

# Wait for chain to start
sleep 15

# Run integration tests
./run_integration_tests.sh

# Clean up
make docker-private-net-stop
docker volume rm docker_mechain-data
```

### Local Development

- Use data persistence to avoid repeated initialization
- Regularly clean logs to avoid disk usage
- Use `docker stats` to monitor resource usage
- Keep test account private keys for development

## ⚠️ Security Notice

- ⚠️ **Production Use Prohibited**: This configuration is for development and testing only
- ⚠️ **Private Key Security**: Pre-configured mnemonic is public, never use in production
- ⚠️ **Network Exposure**: Do not expose ports to the public internet
- ⚠️ **Keyring Backend**: Use `test` backend for development only, use `file` or `os` in production

## 🆘 Getting Help

Having issues?

1. Check the [Troubleshooting](#-troubleshooting) section
2. Run `make docker-private-net-test` to diagnose problems
3. View container logs: `docker logs mechain-private-net`
4. Submit an issue to the project repository

## 📝 Changelog

### v1.0.0 (2024-11-12)
- ✅ Initial release
- ✅ One-click private network startup
- ✅ Built-in complete test suite
- ✅ Data persistence support
- ✅ Health check functionality
- ✅ Release image for production deployments
- ✅ GitHub Actions CI/CD workflows

---

**Last Updated**: 2024-11-12  
**Maintainer**: ME-Chain Team