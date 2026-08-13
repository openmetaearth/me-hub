#!/usr/bin/env bash
#
# Generate TypeScript proto types into ts-client/metaearth.<module>/types
# using buf + stephenh-ts-proto (see proto/buf.gen.ts.yaml).
#
# Why not Ignite?
#   Official/local ignite often fails on tool/version mismatches. This script
#   only regenerates types/ (protobuf codecs). It does NOT rewrite Ignite
#   wrappers (module.ts / rest.ts / registry.ts / index.ts). Update those
#   manually when Msg/Query RPCs change.
#
# Why an isolated workdir?
#   Repo root may have both buf.yaml (v2) and buf.work.yaml (v1), which buf
#   rejects. Generating from a copy of proto/ avoids that conflict and uses
#   proto/buf.yaml deps reliably.
#
# Usage:
#   ./scripts/protocgen-ts.sh              # all metaearth modules that have ts-client dirs
#   ./scripts/protocgen-ts.sh dao rollapp  # only selected modules
#   make proto-gen-ts
#
# Prerequisites:
#   - buf CLI on PATH (brew install buf / go install ...)
#   - Network access to buf.build (remote plugin + deps)
#   - Optional: buf registry login  (avoids BSR rate limits)
#
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if ! command -v buf >/dev/null 2>&1; then
  echo "error: buf not found on PATH" >&2
  exit 1
fi

if [[ ! -f proto/buf.gen.ts.yaml ]]; then
  echo "error: proto/buf.gen.ts.yaml missing" >&2
  exit 1
fi

# Modules requested on CLI, or discover from proto/metaearth + existing ts-client.
if [[ $# -gt 0 ]]; then
  MODULES=("$@")
else
  MODULES=()
  for d in proto/metaearth/*/; do
    [[ -d "$d" ]] || continue
    mod="$(basename "$d")"
    # skip non-module dirs if any
    case "$mod" in
      common|migratetest) continue ;;
    esac
    if [[ -d "ts-client/metaearth.${mod}" ]]; then
      MODULES+=("$mod")
    fi
  done
fi

if [[ ${#MODULES[@]} -eq 0 ]]; then
  echo "error: no modules to generate (need proto/metaearth/<mod> and ts-client/metaearth.<mod>)" >&2
  exit 1
fi

WORKDIR="$(mktemp -d "${TMPDIR:-/tmp}/me-hub-ts-client.XXXXXX")"
cleanup() { rm -rf "$WORKDIR"; }
trap cleanup EXIT

echo ">> preparing isolated proto workdir at $WORKDIR"
cp -R proto "$WORKDIR/proto"
# Drop parent-style configs if somehow copied; keep only proto module config.
rm -f "$WORKDIR/buf.yaml" "$WORKDIR/buf.work.yaml" 2>/dev/null || true

echo ">> buf dep update"
(cd "$WORKDIR/proto" && buf dep update)

echo ">> generating TypeScript types for: ${MODULES[*]}"
for mod in "${MODULES[@]}"; do
  out="$ROOT/ts-client/metaearth.${mod}/types"
  if [[ ! -d "proto/metaearth/${mod}" ]]; then
    echo "skip ${mod}: proto/metaearth/${mod} not found" >&2
    continue
  fi
  if [[ ! -d "ts-client/metaearth.${mod}" ]]; then
    echo "skip ${mod}: ts-client/metaearth.${mod} not found (create module wrappers first)" >&2
    continue
  fi

  mkdir -p "$out"
  echo "   - metaearth.${mod} -> ts-client/metaearth.${mod}/types"
  (cd "$WORKDIR/proto" && buf generate \
    --template buf.gen.ts.yaml \
    -o "$out" \
    --path "metaearth/${mod}")
done

echo ">> done"
echo "note: only types/ were regenerated; review module.ts/registry.ts/rest.ts/index.ts if APIs changed."
