FROM --platform=linux/amd64 golang:1.22-bullseye as go-builder
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    git \
    libc6-dev \
    jq bash vim file

RUN go mod download
RUN make build -j$(nproc)
RUN cp `ldd ./build/med | grep -oP '(/.*libwasmvm.x86_64.so)' -o` /go/

FROM ubuntu:22.04

COPY --from=go-builder /app/build/med /usr/local/bin/
COPY --from=go-builder /go/libwasmvm.x86_64.so /lib/x86_64-linux-gnu/libwasmvm.x86_64.so
WORKDIR /app

COPY scripts/ ./scripts/

ENV KEY_NAME=local-user
ENV MONIKER_NAME=local

RUN chmod +x ./scripts/*.sh

EXPOSE 26656 26657 1317 9090
