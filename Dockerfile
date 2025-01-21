FROM golang:1.23-bullseye AS go-builder
ARG arch=x86_64
ARG LINK_STATICALLY

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    git \
    libc6-dev

RUN go mod download
RUN make build
RUN ldd ./build/med && \
    LIB_PATH=$(ldd ./build/med | grep -o '/go/pkg/mod/github.com/!cosm!wasm/wasmvm@v1.3.0/internal/api/libwasmvm\.aarch64\.so') && \
    echo "Library path found: $LIB_PATH" && \
    cp "$LIB_PATH" /go/libwasmvm.aarch64.so

FROM ubuntu:22.04
WORKDIR root
RUN apt-get update && apt-get install -y curl jq bash vim

COPY --from=go-builder /app/build/med /usr/local/bin/
COPY --from=go-builder /go/libwasmvm.aarch64.so /lib/x86_64-linux-gnu/libwasmvm.aarch64.so

ENV LD_LIBRARY_PATH=/lib/x86_64-linux-gnu

EXPOSE 36656/tcp 36657/tcp 36660/tcp 8090/tcp 1318/tcp 8545/tcp 8546/tcp
VOLUME ["/root"]
ENTRYPOINT ["med"]
