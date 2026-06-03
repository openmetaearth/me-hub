FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git make gcc libc-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:3.18

RUN apk add --no-cache ca-certificates curl jq

COPY --from=builder /app/build/med /usr/local/bin/med

EXPOSE 26657 1317 9090

# Add healthcheck script
COPY scripts/healthcheck.sh /usr/local/bin/healthcheck.sh
RUN chmod +x /usr/local/bin/healthcheck.sh

HEALTHCHECK --interval=10s --timeout=3s --start-period=30s --retries=3 \
  CMD /usr/local/bin/healthcheck.sh

ENTRYPOINT ["med"]
CMD ["start"]
