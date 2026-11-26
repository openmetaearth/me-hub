#!/usr/bin/env sh
set -eu

# Prepare on the host: proto-builder's Go version predates this repository's
# go.mod. Generate in the pinned image, without reusing swagger-proto or
# tmp-swagger-gen (both may contain schemas from an older SDK).
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"
for tool in go docker jq curl; do
  command -v "$tool" >/dev/null 2>&1 || {
    printf 'Swagger generation requires %s on the host\n' "$tool" >&2
    exit 1
  }
done

WORK_DIR=$(mktemp -d "${TMPDIR:-/tmp}/me-hub-swagger.XXXXXX")
trap 'chmod -R u+w "$WORK_DIR"; rm -rf "$WORK_DIR"' EXIT
trap 'exit 1' HUP INT TERM
mkdir -p "$WORK_DIR/proto" "$WORK_DIR/tmp-swagger-gen"

module_dir() {
  go list -m -mod=readonly -f '{{.Dir}}' "$1"
}

# Use the selected go.mod replacements/versions, not upstream branches or the
# (potentially older) SDK/IBC versions in proto/buf.lock.
for module in github.com/cosmos/cosmos-sdk github.com/cosmos/ibc-go/v8 \
  github.com/evmos/ethermint github.com/CosmWasm/wasmd \
  github.com/cosmos/cosmos-proto github.com/cometbft/cometbft; do
  module_path=$(module_dir "$module")
  if [ ! -d "$module_path/proto" ]; then
    printf 'Missing proto sources for %s; run go mod download first\n' "$module" >&2
    exit 1
  fi
  printf 'Swagger schemas: %s\n' "$module_path"
  cp -R "$module_path/proto/." "$WORK_DIR/proto/"
  chmod -R u+w "$WORK_DIR/proto"
done

# SDK's proto tree also contains the extracted nft/feegrant/upgrade/evidence
# schemas; those Go module archives do not include .proto files.
gogo_path=$(module_dir github.com/cosmos/gogoproto)
mkdir -p "$WORK_DIR/proto/gogoproto"
cp "$gogo_path/gogoproto/gogo.proto" "$WORK_DIR/proto/gogoproto/"
gateway_path=$(module_dir github.com/grpc-ecosystem/grpc-gateway)
cp -R "$gateway_path/third_party/googleapis/google" "$WORK_DIR/proto/"

# ICS23's Go submodule omits its proto sources. Fetch only the schema at the
# exact VCS revision recorded for the selected module, never a mutable branch.
ics23_module=$(go list -m -mod=readonly -f '{{if .Replace}}{{.Replace.Path}}@{{.Replace.Version}}{{else}}{{.Path}}@{{.Version}}{{end}}' github.com/cosmos/ics23/go)
go mod download -json "$ics23_module" > "$WORK_DIR/ics23-module.json"
ics23_revision=$(jq -er '.Origin.Hash // empty' "$WORK_DIR/ics23-module.json")
case "$ics23_revision" in
  *[!0-9a-f]*|'') printf 'Missing immutable ICS23 source revision\n' >&2; exit 1 ;;
esac
mkdir "$WORK_DIR/ics23-download"
curl --fail --silent --show-error --location \
  "https://raw.githubusercontent.com/cosmos/ics23/$ics23_revision/proto/cosmos/ics23/v1/proofs.proto" \
  -o "$WORK_DIR/ics23-download/proofs.proto"
mkdir -p "$WORK_DIR/proto/cosmos/ics23/v1"
cp "$WORK_DIR/ics23-download/proofs.proto" "$WORK_DIR/proto/cosmos/ics23/v1/"
printf 'Swagger schemas: ICS23 %s\n' "$ics23_revision"

# A single fresh Buf module: no dependency's buf.lock/config may choose an
# external SDK over the exact selected sources staged above.
rm -f "$WORK_DIR/proto/buf.yaml" "$WORK_DIR/proto/buf.lock" "$WORK_DIR/proto/buf.work.yaml"
printf 'version: v1\n' > "$WORK_DIR/proto/buf.yaml"
cp -R "$ROOT/proto/metaearth" "$WORK_DIR/proto/"
cp "$ROOT/proto/buf.gen.swagger.yaml" "$WORK_DIR/buf.gen.swagger.yaml"
cp "$ROOT/docs/config.json" "$WORK_DIR/config.json"
jq -r '.apis[].url | sub("^\\./tmp-swagger-gen/"; "") | sub("\\.swagger\\.json$"; ".proto")' \
  "$WORK_DIR/config.json" > "$WORK_DIR/query-paths"

# PROTO_DOCKER_ARGS intentionally permits Docker CLI options, as in Makefile's
# other proto targets. The generator itself needs no remote Buf plugins.
# shellcheck disable=SC2086
docker run --rm ${PROTO_DOCKER_ARGS:-} -v "$WORK_DIR:/workspace" \
  --user "$(id -u):$(id -g)" --env HOME=/workspace --workdir /workspace \
  "${PROTO_SWAGGER_IMAGE:-ghcr.io/cosmos/proto-builder:0.14.0}" sh -ec '
    cd /workspace/proto
    while IFS= read -r query; do
      buf generate --template /workspace/buf.gen.swagger.yaml --path "$query"
    done < /workspace/query-paths
    cd /workspace
    swagger-combine config.json -o openapi.yml -f yaml \
      --continueOnConflictingPaths true --includeDefinitions true
    test -s openapi.yml
  '

# Publish only after every fresh fragment and the combined document succeed.
cp "$WORK_DIR/openapi.yml" "$ROOT/docs/static/openapi.yml"
printf 'Updated docs/static/openapi.yml\n'
