#!/usr/bin/make -f

COMMIT := $(shell git log -1 --format='%H')
VERSION ?= $(shell git describe --tags --always)
TAG ?= latest

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
TM_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
DOCKER := $(shell which docker)
BUILDDIR ?= $(CURDIR)/build

# Dependencies version
DEPS_COSMOS_SDK_VERSION := $(shell cat go.sum | grep 'github.com/openmetaearth/cosmos-sdk' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_ETHERMINT_VERSION := $(shell cat go.sum | grep 'github.com/openmetaearth/ethermint' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_OSMOSIS_VERSION := $(shell cat go.sum | grep 'github.com/openmetaearth/osmosis' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_IBC_GO_VERSION := $(shell cat go.sum | grep 'github.com/cosmos/ibc-go' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_COSMOS_PROTO_VERSION := $(shell cat go.sum | grep 'github.com/cosmos/cosmos-proto' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_COSMOS_GOGOPROTO_VERSION := $(shell cat go.sum | grep 'github.com/cosmos/gogoproto' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_CONFIO_ICS23_VERSION := go/$(shell cat go.sum | grep 'github.com/confio/ics23/go' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_WASM_VERSION := $(shell cat go.sum | grep 'github.com/CosmWasm/wasmd' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(ME_BUILD_OPTIONS)))
  build_tags += gcc cleveldb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace := $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=me-hub \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=med \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
	      -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TM_VERSION)

ifeq (cleveldb,$(findstring cleveldb,$(ME_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ifeq (,$(findstring nostrip,$(ME_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(ME_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

all: install

.PHONY: install
install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/med

.PHONY: build build-debug build-linux build-test

build: go.sum
	go build $(BUILD_FLAGS) -o $(BUILDDIR)/med ./cmd/med

build-vendor: go.sum
	@echo "Building with vendor mode (without ledger support for Docker)..."
	$(eval temp_build_tags := $(filter-out ledger,$(build_tags)))
	@go build -mod=vendor -tags "$(temp_build_tags)" -ldflags '$(ldflags)' -trimpath -o $(BUILDDIR)/med ./cmd/med
	@echo "Build completed successfully"

TRIGGER_BLOCKS ?= 100
build-test: go.sum
	$(eval temp_ldflags := $(filter-out -w -s,$(ldflags)) -X github.com/openmetaearth/me-hub/x/wmint/types.OneDayTotalBlocks=$(TRIGGER_BLOCKS))
	go build -tags "$(build_tags)" -ldflags '$(temp_ldflags)' -o $(BUILDDIR)/med ./cmd/med

build-debug: go.sum
	$(eval temp_ldflags := $(filter-out -w -s,$(ldflags)))
	go build -tags "$(build_tags)" -ldflags '$(temp_ldflags)' -gcflags "all=-N -l" -o $(BUILDDIR)/med ./cmd/med

build-linux: go.sum
	CC=x86_64-unknown-linux-gnu-gcc CGO_ENABLED=1 TARGET_CC=clang LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/med ./cmd/med

build-linux-debug: go.sum
	$(eval temp_ldflags := $(filter-out -w -s,$(ldflags)))
	CC=x86_64-unknown-linux-gnu-gcc CGO_ENABLED=1 TARGET_CC=clang LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 go build -tags "$(build_tags)" -ldflags '$(temp_ldflags)' -gcflags "all=-N -l" -o $(BUILDDIR)/med-debug ./cmd/med

###############################################################################
###                                Docker                                ###
###############################################################################
.PHONY: docker-github docker-local docker-run-debug docker-private-net docker-private-net-start docker-private-net-stop docker-private-net-test docker-release

docker-github:
	DOCKER_BUILDKIT=1 docker build -t ghcr.io/me-hub/med:latest -f Dockerfile .

# docker pull --platform=linux/amd64 ubuntu:24.04
docker-local: build-linux
	@DOCKER_BUILDKIT=1 docker build -t 192.168.0.79/me-hub/med:$(TAG) -f Dockerfile_local .
	@docker push 192.168.0.79/me-hub/med:$(TAG)

docker-run-debug:
	@DOCKER_BUILDKIT=1 docker-compose -f docker-compose.debug.yml up

# Build and optionally run a pre-initialized single-node private network
# Usage:
#   make docker-private-net
#   make docker-private-net GENESIS_ACCOUNTS="test1:1000000000000000000000umec,test2:500000000000000000000umec"
#   make docker-private-net GENESIS_ACCOUNTS_JSON='[{"name":"alice","amount":"2000000000000000000000umec"}]'
docker-private-net:
	@echo "Preparing vendor directory for Docker build..."
	@go mod vendor
	@echo "Building ME-Chain private network Docker image..."
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		$(if $(GENESIS_ACCOUNTS),--build-arg GENESIS_ACCOUNTS="$(GENESIS_ACCOUNTS)",) \
		$(if $(GENESIS_ACCOUNTS_JSON),--build-arg GENESIS_ACCOUNTS_JSON='$(GENESIS_ACCOUNTS_JSON)',) \
		-t me-hub/private-net:$(TAG) \
		-f docker/Dockerfile .
	@rm -rf vendor
	@echo "Docker image built successfully: me-hub/private-net:$(TAG)"
	@echo ""
	@echo "To run the private network:"
	@echo "  docker run -d -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 --name mechain-private-net me-hub/private-net:$(TAG)"
	@echo ""
	@echo "To run with persistent data:"
	@echo "  docker run -d -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 -v mechain-data:/root/.mechain --name mechain-private-net me-hub/private-net:$(TAG)"
	@echo ""
	@echo "To run with additional genesis accounts at runtime:"
	@echo "  docker run -d -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \\"
	@echo "    -e GENESIS_ACCOUNTS=\"test1:1000000000000000000000umec,test2:500000000000000000000umec\" \\"
	@echo "    --name mechain-private-net me-hub/private-net:$(TAG)"
	@echo ""
	@echo "Or use: make docker-private-net-start"

# Start the private network using docker-compose
docker-private-net-start:
	@echo "Starting ME-Chain private network..."
	@docker compose -f docker/docker-compose.yml up -d
	@echo ""
	@echo "Private network started successfully!"
	@echo "RPC: http://localhost:36657"
	@echo "API: http://localhost:1318"
	@echo "JSON-RPC: http://localhost:9545"
	@echo ""
	@echo "View logs: docker compose -f docker/docker-compose.yml logs -f"
	@echo "Run tests: make docker-private-net-test"

# Stop the private network
docker-private-net-stop:
	@echo "Stopping ME-Chain private network..."
	@docker compose -f docker/docker-compose.yml down
	@echo "Private network stopped."

# Test the private network
docker-private-net-test:
	@echo "Running tests on private network..."
	@chmod +x docker/scripts/test_private_net.sh
	@./docker/scripts/test_private_net.sh

# Build release Docker image (no chain initialization)
# Usage:
#   make docker-release
#   make docker-release TAG=v1.0.0
docker-release:
	@echo "Preparing vendor directory for Docker build..."
	@go mod vendor
	@echo "Building ME-Chain release Docker image..."
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		-t me-hub/release:$(TAG) \
		-f docker/Dockerfile.release .
	@rm -rf vendor
	@echo "Docker image built successfully: me-hub/release:$(TAG)"
	@echo ""
	@echo "To verify the image:"
	@echo "  docker run --rm me-hub/release:$(TAG) version"
	@echo ""
	@echo "To run med commands:"
	@echo "  docker run --rm me-hub/release:$(TAG) --help"
	@echo ""
	@echo "To use as a base for custom chain setup:"
	@echo "  docker run -it --rm me-hub/release:$(TAG) init mynode --chain-id mychain"

###############################################################################
###                                Releasing                                ###
###############################################################################

PACKAGE_NAME := $(shell go list -m)
GOLANG_CROSS_VERSION  = v1.23
GOPATH ?= '$(HOME)/go'
COSMWASM_VERSION := $(shell go list -m github.com/CosmWasm/wasmvm | sed 's/.* //')
release-dry-run:
	docker run --privileged -e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip=validate --skip=publish --snapshot

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	@echo "Running release process"
	@echo "COSMWASM version: $(COSMWASM_VERSION)"
	docker run --rm --privileged \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip=validate --release-notes ./release-note.md

.PHONY: release-dry-run release

###############################################################################
###                                Proto                                    ###
###############################################################################
protoCosmosVer=0.14.0
protoCosmosName=ghcr.io/cosmos/proto-builder:$(protoCosmosVer)
protoCosmosImage=docker run --rm -v $(CURDIR):/workspace --user root --workdir /workspace $(protoCosmosName)

proto-all: proto-format proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoCosmosImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Downloading Protobuf dependencies"
	@#make proto-download-deps
	@echo "Generating Protobuf Swagger"
	@$(protoCosmosImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@$(protoCosmosImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoCosmosImage) buf lint --error-format=json

SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party

proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git clone -b me-hub/v0.47.13 --single-branch --depth 1 https://github.com/openmetaearth/cosmos-sdk.git && \
	rm -f ./cosmos-sdk/proto/buf.* && \
	mv ./cosmos-sdk/proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ethermint_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ethermint_tmp" && \
	git clone -b dev --single-branch --depth 1 https://github.com/openmetaearth/ethermint.git && \
	rm -f ./ethermint/proto/buf.* && \
	mv ./ethermint/proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ethermint_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/wasm_tmp" && \
	cd "$(THIRD_PARTY_DIR)/wasm_tmp" && \
	git clone --branch v0.43.0 --single-branch --depth 1 https://github.com/CosmWasm/wasmd.git && \
	rm -f ./wasmd/proto/buf.* && \
	mv ./wasmd/proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/wasm_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/ibc-go.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git fetch --depth=1 origin "$(DEPS_IBC_GO_VERSION)" && \
	git checkout FETCH_HEAD && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-proto.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git fetch --depth=1 origin "$(DEPS_COSMOS_PROTO_VERSION)" && \
	git checkout FETCH_HEAD && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_proto_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL https://raw.githubusercontent.com/cosmos/gogoproto/$(DEPS_COSMOS_GOGOPROTO_VERSION)/gogoproto/gogo.proto > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/confio/ics23" && \
	curl -sSL https://raw.githubusercontent.com/confio/ics23/$(DEPS_CONFIO_ICS23_VERSION)/proofs.proto > "$(THIRD_PARTY_DIR)/proofs.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos/ics23/v1" && \
	curl -sSL "https://raw.githubusercontent.com/cosmos/ics23/refs/heads/master/proto/cosmos/ics23/v1/proofs.proto" > "$(THIRD_PARTY_DIR)/cosmos/ics23/v1/proofs.proto"


.PHONY: proto-gen proto-swagger-gen proto-format proto-lint proto-download-deps

###############################################################################
###                                Linting                                  ###
###############################################################################

golangci_version=v1.60.3

lint-install:
	@echo "--> Installing golangci-lint $(golangci_version)"
	@if golangci-lint version --format json | jq .version | grep -q $(golangci_version); then \
		echo "golangci-lint $(golangci_version) is already installed"; \
	else \
		echo "Installing golangci-lint $(golangci_version)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version); \
	fi

lint: lint-install
	@echo "--> Running linter"
	@golangci-lint run --build-tags=$(GO_BUILD) --out-format=tab

format: lint-install
	@golangci-lint run --build-tags=$(GO_BUILD) --out-format=tab --fix

shell-lint:
	# install shellcheck > https://github.com/koalaman/shellcheck
	grep -r '^#!/usr/bin/env bash' --exclude-dir={node_modules,build} . | cut -d: -f1 | xargs shellcheck

shell-format:
	# install shfmt > https://github.com/mvdan/sh
	#go install mvdan.cc/sh/v3/cmd/shfmt@v3.8.0
	grep -r '^#!/usr/bin/env bash' --exclude-dir={node_modules,build} . | cut -d: -f1 | xargs shfmt -l -w -i 2

.PHONY: format lint shell-lint shell-format

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@echo "--> Running tests"
	go test -mod=readonly ./...

test-count:
	go test -mod=readonly -cpu 1 -count 1 -cover ./... | grep -v 'types\|cli\|no test files'

test-nightly:
	@TEST_INTEGRATION=true go test -mod=readonly -timeout 20m -cpu 4 -v -run TestIntegrationTest ./tests
	@TEST_CROSSCHAIN=true go test -mod=readonly -cpu 4 -v -run TestCrosschainKeeperTestSuite ./x/crosschain/...

mocks:
mocks:
	@go install github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	#mockgen -source=x/wdistri/types/expected_keepers.go -package mock -destination x/wdistri/types/mock/expected_keepers_mock.go
	mockgen -source=app/ante/expected_keepers.go -package mock -destination app/ante/mock/expected_keepers_mocks.go

.PHONY: test test-count test-nightly mocks
