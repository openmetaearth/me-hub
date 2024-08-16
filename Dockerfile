FROM golang:1.22-alpine3.19 as go-builder

ARG LINK_STATICALLY

WORKDIR /app

COPY . .


RUN go mod download

ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

RUN apk add --no-cache $PACKAGES

RUN make build

FROM alpine:3.16.1

RUN apk add --no-cache curl jq bash vim 

COPY --from=go-builder /app/build/med /usr/local/bin/

WORKDIR /app

COPY scripts/* ./scripts/

ENV KEY_NAME=local-user
ENV MONIKER_NAME=local

RUN chmod +x ./scripts/*.sh

EXPOSE 26656 26657 1317 9090
