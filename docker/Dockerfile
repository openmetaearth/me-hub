# syntax=docker/dockerfile:1

# Dockerfile for ME-Chain private network
# This creates a single-node private network that's pre-compiled and pre-initialized
# Simply run the container to start a fully configured private network

FROM golang:1.23-bullseye as go-builder

ARG arch=x86_64
ARG GIT_VERSION
ARG GIT_COMMIT
ARG GENESIS_ACCOUNTS=""
ARG GENESIS_ACCOUNTS_JSON=""

WORKDIR /me-hub

# Copy source code
COPY . .

# Build using vendor mode (make will run go mod vendor)
RUN make build-vendor

# Find and copy the dynamic libwasmvm library
RUN cp `ldd ./build/med | grep -P '/.+libwasmvm.*.so' -o` /go/

# --------------------------------------------------------
# Runtime Stage
# --------------------------------------------------------

FROM ubuntu:24.04

# Install required packages
RUN apt-get update && \
    apt-get install -y \
    ca-certificates \
    tzdata \
    libc6 \
    wget \
    git \
    jq \
    curl && \
    rm -rf /var/lib/apt/lists/*

# Set the timezone to Asia/Shanghai
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Copy binary from builder
COPY --from=go-builder /me-hub/build/med /usr/local/bin/med
RUN chmod +x /usr/local/bin/med

# Copy dynamic library
COPY --from=go-builder /go/libwasmvm.*.so /lib/

# Copy initialization scripts
COPY ./docker/scripts/setup_local_docker.sh /scripts/setup_local_docker.sh
COPY ./docker/scripts/start_private_net.sh /scripts/start_private_net.sh
COPY ./scripts/src/genesis_config_commands.sh /scripts/src/genesis_config_commands.sh
RUN chmod +x /scripts/*.sh

# Set working directory
WORKDIR /root

# Environment variables with defaults
ENV CHAIN_ID="mechain_202404-1"
ENV MONIKER_NAME="local"
ENV KEY_NAME="global_dao"
ENV HOME="/root"

# Genesis accounts configuration
# Format: "account1:amount1,account2:amount2,..."
# Example: GENESIS_ACCOUNTS="test1:1000000000000000000000umec,test2:500000000000000000000umec"
ARG GENESIS_ACCOUNTS
ENV GENESIS_ACCOUNTS=${GENESIS_ACCOUNTS}

# Alternative JSON format for more control (supports mnemonic)
# Example: GENESIS_ACCOUNTS_JSON='[{"name":"test1","amount":"1000000000000000000000umec","mnemonic":"..."}]'
ARG GENESIS_ACCOUNTS_JSON
ENV GENESIS_ACCOUNTS_JSON=${GENESIS_ACCOUNTS_JSON}

# Expose ports
# P2P
EXPOSE 36656
# RPC
EXPOSE 36657
# API
EXPOSE 1318
# gRPC
EXPOSE 8090
# gRPC-web
EXPOSE 8091
# JSON-RPC
EXPOSE 9545
# JSON-RPC WebSocket
EXPOSE 9546

# Volume for persistent data (optional)
VOLUME ["/root/.mechain"]

# Use startup script as entrypoint
ENTRYPOINT ["/scripts/start_private_net.sh"]
