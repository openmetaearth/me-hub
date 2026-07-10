# Use Alpine 3.18 (musl 1.2.4) for compatibility with wasmvm v1.4.1's prebuilt muslc static lib,
# which references LFS64 symbols (fstat64, etc.) that musl >= 1.2.5 (Alpine 3.19+) no longer exports.
# Go 1.21 here is just the bootstrap toolchain; go.mod's `toolchain go1.24.7` triggers auto-download
# of Go 1.24.7 for the actual build.
FROM golang:1.21-alpine3.18 AS builder

RUN apk add --no-cache git build-base linux-headers binutils-gold

WORKDIR /app

COPY go.mod go.sum ./

# Cosmwasm - Download correct libwasmvm version
RUN ARCH=$(uname -m) && \
    WASMVM_VERSION=$(awk '$1 == "github.com/CosmWasm/wasmvm/v2" { print $2; exit }' go.mod) && \
    test -n "$WASMVM_VERSION" && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.$ARCH.a \
    -O /lib/libwasmvm_muslc.$ARCH.a && \
    # verify checksum
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm_muslc.$ARCH.a | grep $(cat /tmp/checksums.txt | grep libwasmvm_muslc.$ARCH.a | cut -d ' ' -f 1) && \
    # wasmvm's link_muslc.go uses `-lwasmvm_muslc` (no arch suffix), so expose the lib under that name
    cp /lib/libwasmvm_muslc.$ARCH.a /lib/libwasmvm_muslc.a

COPY . .

ARG GIT_VERSION=dev
ARG GIT_COMMIT=unknown

RUN make build VERSION="$GIT_VERSION" COMMIT="$GIT_COMMIT" LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true

FROM alpine:3.19

WORKDIR /root

COPY --from=builder /app/build/med /usr/bin/med

EXPOSE 26656/tcp 26657/tcp 26660/tcp 9090/tcp 1317/tcp 8545/tcp 8546/tcp

VOLUME ["/root"]

ENTRYPOINT ["med"]