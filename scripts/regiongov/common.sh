#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
source "$ROOT/config.sh"
command -v "$MED_BIN" >/dev/null || { echo "med binary not found: $MED_BIN" >&2; exit 1; }
command -v jq >/dev/null || { echo 'jq is required' >&2; exit 1; }
NET=()
KEY=()
[[ -z $NODE ]] || NET+=(--node "$NODE")
[[ -z $KEYRING_BACKEND ]] || KEY+=(--keyring-backend "$KEYRING_BACKEND")
fail() { echo "FAIL: $*" >&2; exit 1; }
query() { "$MED_BIN" query regiongov "$@" ${NET[@]+"${NET[@]}"} --output json; }
key_address() { "$MED_BIN" keys show "$1" -a ${KEY[@]+"${KEY[@]}"}; }
assert_json() { jq -e "$@" >/dev/null || fail 'JSON assertion failed'; }

# Never retry broadcasting: a timeout may mean the transaction has already committed.
send_tx() {
 local signer=$1 hash deadline response
 shift
 local options=(--from "$signer" --chain-id "$CHAIN_ID" --gas "$GAS" --gas-adjustment "$GAS_ADJUSTMENT" --broadcast-mode sync --yes --output json)
 [[ -z $FEES ]] || options+=(--fees "$FEES")
 [[ -z $GAS_PRICES ]] || options+=(--gas-prices "$GAS_PRICES")
 echo "TX: $* (from $signer)" >&2
 response=$("$MED_BIN" tx regiongov "$@" "${options[@]}" ${NET[@]+"${NET[@]}"} ${KEY[@]+"${KEY[@]}"}) || fail "broadcast failed: $*"
 printf '%s\n' "$response" > "$RESULT_DIR/broadcast.json"
 printf '%s' "$response" | assert_json '(.code // 0 | tonumber) == 0'
 hash=$(printf '%s' "$response" | jq -er '.txhash | select(length > 0)')
 deadline=$((SECONDS + TX_TIMEOUT))
 while (( SECONDS < deadline )); do
  if "$MED_BIN" query tx "$hash" ${NET[@]+"${NET[@]}"} --output json > "$RESULT_DIR/$hash.json" 2> "$RESULT_DIR/query-tx-error.log"; then
   TX_JSON=$(cat "$RESULT_DIR/$hash.json")
   printf '%s' "$TX_JSON" | assert_json '(.code | tonumber) == 0'
   echo "PASS committed: $hash" >&2
   return
  fi
  sleep "$POLL_INTERVAL"
 done
 fail "timeout waiting for $hash; inspect chain before rerunning; results: $RESULT_DIR"
}
event_id() {
 printf '%s' "$TX_JSON" | jq -er --arg key "$1" '
 [(.logs[]?.events[]?, .events[]?) | .attributes[]? |
 select(.key == $key) | .value] | unique | if length == 1 then .[0] else error("missing/ambiguous event id") end'
}
action_status() {
 query pending-action "$1" | assert_json --arg status "$2" '.action.status == $status'
}
approve() {
 local id=$1
 send_tx "$DAO_KEY" approve-pending-action "$id"
 action_status "$id" PENDING_ACTION_STATUS_APPROVED
}
assert_equity() {
 query equity "$REGION_ID" | assert_json --argjson expected "$1" '(.equity.addresses // [] | sort) == ($expected | sort)'
}
assert_frozen() {
 query equity-freezes "$REGION_ID" | assert_json --argjson addresses "$TARGET_JSON" --argjson frozen "$1" '
 .freezes as $states | all($addresses[]; . as $a | ([$states[]? | select(.address == $a) | .frozen][0] // false) == $frozen)'
}
