#!/usr/bin/env bash

set -eo pipefail

# proto-builder:0.14.0 ships buf that only accepts config version v1.
# Repo-root buf.yaml is v2 (local ts-client). Generate from proto/ so
# proto/buf.yaml is used instead.
echo "Generating gogo proto code"
cd proto
buf format -w

# get protoc executions
# go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null

proto_dirs=$(find ../proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
    for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
      if grep go_package $file &>/dev/null; then
        buf generate --template buf.gen.gogo.yaml $file
      fi
    done
done

cd ..

# TypeScript client types (ts-client/metaearth.*/types) — use dedicated script:
#   ./scripts/protocgen-ts.sh
#   make proto-gen-ts

# move proto files to the right places
# Note: Proto files are suffixed with the current binary version.
cp -r github.com/openmetaearth/me-hub/* ./
rm -rf github.com