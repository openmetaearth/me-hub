# Release Build Guide

This document explains how to build and use the minimal release Docker image for ME-Chain.

## Overview

The release image is a minimal, production-ready Docker image that contains only the compiled `med` binary and required libraries. Unlike the private network image, it does **not** include any pre-initialized chain data or test accounts.

## Image Comparison

| Feature | Release Image | Private Network Image |
|---------|---------------|----------------------|
| **Image Tag** | `me-hub/release:*` | `me-hub/private-net:*` |
| **Size** | ~200MB | ~423MB |
| **Contains** | Binary + libraries only | Binary + initialized chain |
| **Pre-configured** | No | Yes (test accounts) |
| **Dockerfile** | `Dockerfile.release` | `Dockerfile` |
| **Entrypoint** | `med` | Chain startup script |
| **Default CMD** | `version` | Starts chain automatically |
| **Use Case** | Production/Custom setups | Development/Testing |
| **Build Target** | `make docker-release` | `make docker-private-net` |
| **CI Workflow** | `build-push-release.yml` | `build-push-private-net.yml` |

## Building Release Image

### Local Build

```bash
# Build with default tag (latest)
make docker-release

# Build with specific tag
make docker-release TAG=v1.0.0

# Build with custom version info
VERSION=v1.0.0 make docker-release TAG=v1.0.0
```

### CI/CD Build

The release image is automatically built when you push a version tag:

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

This triggers the GitHub Actions workflow `.github/workflows/build-push-release.yml` which will:
1. Build the release Docker image
2. Tag it with the git tag (e.g., `v1.0.0`)
3. Push to Harbor registry as:
   - `{HARBOR_REGISTRY}/openmetaearth/me_hub:v1.0.0`
   - `{HARBOR_REGISTRY}/openmetaearth/me_hub:latest`
4. Verify the image by running `med version`

## Using Release Image

### 1. Check Version

```bash
# Run version command
docker run --rm me-hub/release:latest version

# Show detailed version information
docker run --rm me-hub/release:latest version --long
```

Expected output:
```
med version, commit <git-commit>, built at <timestamp>
```

### 2. View Help

```bash
# Show all available commands
docker run --rm me-hub/release:latest --help

# Show help for specific command
docker run --rm me-hub/release:latest init --help
docker run --rm me-hub/release:latest keys --help
docker run --rm me-hub/release:latest start --help
```

### 3. Initialize Custom Chain

```bash
# Create a volume for chain data
docker volume create mychain-data

# Initialize a new chain
docker run --rm \
  -v mychain-data:/root/.mechain \
  me-hub/release:latest \
  init mynode --chain-id mychain_100-1

# Add a key
docker run --rm -it \
  -v mychain-data:/root/.mechain \
  me-hub/release:latest \
  keys add mykey --keyring-backend test

# Add genesis account
docker run --rm \
  -v mychain-data:/root/.mechain \
  me-hub/release:latest \
  add-genesis-account mykey 1000000000000000000000umec --keyring-backend test

# Create genesis transaction
docker run --rm \
  -v mychain-data:/root/.mechain \
  me-hub/release:latest \
  gentx mykey 100000000000000000000umec --chain-id mychain_100-1 --keyring-backend test

# Collect genesis transactions
docker run --rm \
  -v mychain-data:/root/.mechain \
  me-hub/release:latest \
  collect-gentxs

# Start the chain
docker run -d \
  -v mychain-data:/root/.mechain \
  -p 36657:36657 -p 1318:1318 -p 9545:9545 \
  --name mychain \
  --entrypoint med \
  me-hub/release:latest \
  start
```

### 4. Use as Base Image

Create your own Dockerfile:

```dockerfile
FROM me-hub/release:latest

# Copy your custom configuration files
COPY genesis.json /root/.mechain/config/genesis.json
COPY config.toml /root/.mechain/config/config.toml
COPY app.toml /root/.mechain/config/app.toml

# Copy initialization script
COPY init-chain.sh /scripts/init-chain.sh
RUN chmod +x /scripts/init-chain.sh

# Set custom entrypoint
ENTRYPOINT ["/scripts/init-chain.sh"]
```

### 5. Interactive Shell

```bash
# Enter container with shell for manual setup
docker run -it --rm \
  -v mychain-data:/root/.mechain \
  --entrypoint sh \
  me-hub/release:latest

# Inside the container, you can run any med commands:
# $ med init mynode --chain-id test
# $ med keys add mykey
# $ med start
```

### 6. Extract Binary

```bash
# Copy binary from container to host
docker create --name temp me-hub/release:latest
docker cp temp:/usr/local/bin/med ./med
docker rm temp

# Verify extracted binary
./med version
```

## Production Deployment Examples

### Example 1: Kubernetes Deployment

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mechain-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mechain-node
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mechain
  template:
    metadata:
      labels:
        app: mechain
    spec:
      containers:
      - name: mechain
        image: harbor.example.com/openmetaearth/me_hub:v1.0.0
        command: ["med"]
        args: ["start"]
        ports:
        - containerPort: 36657
          name: rpc
        - containerPort: 1318
          name: api
        - containerPort: 9545
          name: jsonrpc
        volumeMounts:
        - name: data
          mountPath: /root/.mechain
        resources:
          requests:
            memory: "4Gi"
            cpu: "2"
          limits:
            memory: "8Gi"
            cpu: "4"
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: mechain-data
```

### Example 2: Docker Compose for Production

```yaml
version: '3.8'

services:
  mechain:
    image: harbor.example.com/openmetaearth/me_hub:v1.0.0
    container_name: mechain-prod
    entrypoint: ["med"]
    command: ["start"]
    ports:
      - "36657:36657"
      - "1318:1318"
      - "9545:9545"
      - "8090:8090"
    volumes:
      - mechain-data:/root/.mechain
      - ./config/config.toml:/root/.mechain/config/config.toml:ro
      - ./config/app.toml:/root/.mechain/config/app.toml:ro
    restart: unless-stopped
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "10"
    healthcheck:
      test: ["CMD", "med", "status"]
      interval: 30s
      timeout: 10s
      retries: 5
      start_period: 60s

volumes:
  mechain-data:
    driver: local
```

### Example 3: Systemd Service (Binary Extraction)

```bash
# 1. Extract binary
docker run --rm -v /usr/local/bin:/target \
  me-hub/release:latest \
  sh -c "cp /usr/local/bin/med /target/"

# 2. Extract library
docker run --rm -v /usr/lib:/target \
  me-hub/release:latest \
  sh -c "cp /lib/libwasmvm.*.so /target/"

# 3. Create systemd service
cat > /etc/systemd/system/mechain.service <<EOF
[Unit]
Description=ME-Chain Node
After=network-online.target

[Service]
User=mechain
WorkingDirectory=/home/mechain
ExecStart=/usr/local/bin/med start
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

# 4. Start service
systemctl daemon-reload
systemctl enable mechain
systemctl start mechain
```

## Verification Steps

After deploying, verify your installation:

```bash
# 1. Check version
docker run --rm your-registry/openmetaearth/me_hub:v1.0.0 version

# 2. Check that help works
docker run --rm your-registry/openmetaearth/me_hub:v1.0.0 --help

# 3. Test initialization
docker run --rm \
  -v test-data:/root/.mechain \
  your-registry/openmetaearth/me_hub:v1.0.0 \
  init testnode --chain-id test_100-1

# 4. Verify files were created
docker run --rm \
  -v test-data:/root/.mechain \
  --entrypoint ls \
  your-registry/openmetaearth/me_hub:v1.0.0 \
  -la /root/.mechain/config/

# 5. Cleanup
docker volume rm test-data
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy to Production

on:
  push:
    tags:
      - 'v*'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Pull release image
        run: |
          docker pull ${{ secrets.HARBOR_REGISTRY }}/openmetaearth/me_hub:${{ github.ref_name }}

      - name: Verify image
        run: |
          docker run --rm ${{ secrets.HARBOR_REGISTRY }}/openmetaearth/me_hub:${{ github.ref_name }} version

      - name: Deploy to production
        run: |
          # Your deployment logic here
          kubectl set image deployment/mechain \
            mechain=${{ secrets.HARBOR_REGISTRY }}/openmetaearth/me_hub:${{ github.ref_name }}
```

## Dockerfile Reference

The release image is built from `docker/Dockerfile.release`:

```dockerfile
# Build stage - compiles the binary
FROM golang:1.23-bullseye as go-builder
WORKDIR /me-hub
COPY . .
RUN make build-vendor
RUN cp `ldd ./build/med | grep -P '/.+libwasmvm.*.so' -o` /go/

# Runtime stage - minimal image
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y ca-certificates tzdata libc6
COPY --from=go-builder /me-hub/build/med /usr/local/bin/med
COPY --from=go-builder /go/libwasmvm.*.so /lib/
WORKDIR /root
ENTRYPOINT ["med"]
CMD ["version"]
```

## Troubleshooting

### Issue: "med: command not found"

**Solution:** The binary is at `/usr/local/bin/med`. If using custom entrypoint, use full path or ensure PATH includes `/usr/local/bin`.

### Issue: "error while loading shared libraries: libwasmvm.*.so"

**Solution:** The library should be in `/lib/`. Verify with:
```bash
docker run --rm --entrypoint ls me-hub/release:latest -l /lib/libwasmvm*.so
```

### Issue: "permission denied" when writing to volume

**Solution:** Ensure proper volume permissions:
```bash
docker run --rm -v mydata:/data --entrypoint sh me-hub/release:latest -c "chmod 755 /data"
```

### Issue: Version shows wrong git tag

**Solution:** Ensure you're building with correct build args:
```bash
VERSION=$(git describe --tags) make docker-release
```

## Security Considerations

1. **Use specific version tags** in production (not `latest`)
2. **Scan images** for vulnerabilities before deployment
3. **Run as non-root user** if possible (requires custom Dockerfile)
4. **Use secrets management** for sensitive data (keys, passwords)
5. **Enable read-only filesystem** where possible
6. **Limit container resources** (CPU, memory)

## Best Practices

1. ✅ Always use specific version tags in production
2. ✅ Use volumes for persistent data
3. ✅ Configure proper logging and monitoring
4. ✅ Set up health checks
5. ✅ Use orchestration (K8s, Docker Swarm) for production
6. ✅ Implement backup strategy for chain data
7. ✅ Test deployments in staging environment first
8. ✅ Document your custom configuration

## Related Documentation

- [Private Network Image](README.md#-private-network-image) - Pre-initialized test network
- [Docker README](README.md) - Complete Docker documentation
- [GitHub Workflows](../.github/workflows/README.md) - CI/CD documentation
- [Build Workflow](../.github/workflows/build-push-release.yml) - Release build workflow
- [Dockerfile.release](Dockerfile.release) - Dockerfile source

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review logs: `docker logs <container-name>`
3. Check GitHub Actions logs for build issues
4. Submit an issue on the project repository

---

**Last Updated:** 2025-11-12
**Image:** `me-hub/release:*`
**Dockerfile:** `docker/Dockerfile.release`
**Workflow:** `.github/workflows/build-push-release.yml`
**Makefile Target:** `make docker-release`
