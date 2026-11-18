#!/usr/bin/env bash
# No broadcast: exercise the installed binary's argument and message validation.
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/common.sh"
SIGNER=$(key_address "$DAO_KEY")
CHAIN_ID=${CHAIN_ID:-$("$MED_BIN" status ${NET[@]+"${NET[@]}"} | jq -er '.NodeInfo.network // .node_info.network')}
TMP=$(mktemp -d "${TMPDIR:-/tmp}/regiongov-negative.XXXXXX")
trap 'rm -rf "$TMP"' EXIT
# Positive control: ensure errors below are not simply missing keys or a broken binary.
"$MED_BIN" tx regiongov vote 1 yes --from "$SIGNER" --chain-id "$CHAIN_ID" --generate-only ${KEY[@]+"${KEY[@]}"} > "$TMP/valid.json"
jq -e '.body.messages[0].proposal_id == "1"' "$TMP/valid.json" >/dev/null
reject() {
 local pattern=$1
 shift
 if "$MED_BIN" tx regiongov "$@" --from "$SIGNER" --chain-id "$CHAIN_ID" --generate-only ${KEY[@]+"${KEY[@]}"} > "$TMP/out" 2> "$TMP/error"; then
  fail "invalid command succeeded: $*"
 fi
 grep -Eiq "$pattern" "$TMP/error" || { cat "$TMP/error" >&2; fail "unexpected failure: $*"; }
 echo "PASS rejected: $*"
}
reject 'positive uint64' vote 0 yes
reject 'yes or no' vote 1 abstain
reject 'arg' freeze-region
reject 'duplicate region' freeze-region chn CHN
reject 'arg' submit-freeze-equity chn
reject 'arg' submit-unfreeze-equity chn
reject 'duplicate address' submit-freeze-equity chn "$SIGNER" "$SIGNER"
reject 'duplicate address' submit-unfreeze-equity chn "$SIGNER" "$SIGNER"
printf '{' > "$TMP/invalid.json"
reject 'decode proposal' submit-proposal "$TMP/invalid.json"
echo 'PASS: offline validation tests; no transactions broadcast'
