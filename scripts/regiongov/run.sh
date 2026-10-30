#!/usr/bin/env bash
set -Eeuo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"
TESTS=(test_submit_equity_update test_approve_pending_action test_reject_pending_action test_submit_freeze_equity test_submit_unfreeze_equity test_freeze_equity test_unfreeze_equity test_freeze_region test_unfreeze_region test_submit_treasury_proposal test_submit_operator_proposal test_vote test_queries)

preflight() {
 local test=$1
STATUS=$("$MED_BIN" status ${NET[@]+"${NET[@]}"})
NETWORK=$(printf '%s' "$STATUS" | jq -er '.NodeInfo.network // .node_info.network')
CHAIN_ID=${CHAIN_ID:-$NETWORK}
[[ $CHAIN_ID == "$NETWORK" ]] || fail "chain ID mismatch: $CHAIN_ID != $NETWORK"
printf '%s' "$STATUS" | assert_json '(.SyncInfo.catching_up // .sync_info.catching_up // false) == false'
for cmd in submit-proposal vote submit-equity-update submit-freeze-equity submit-unfreeze-equity approve-pending-action reject-pending-action freeze-region unfreeze-region freeze-equity unfreeze-equity; do
 "$MED_BIN" tx regiongov "$cmd" --help >/dev/null
done
query params >/dev/null
[[ $test != test_queries ]] || return 0
# These checks run before any transaction or state mutation.
REGION=$("$MED_BIN" query staking region "$REGION_ID" ${NET[@]+"${NET[@]}"} --output json) || fail "region $REGION_ID does not exist"
query region-freeze "$REGION_ID" | assert_json '.frozen == false'
REGIONS=("$REGION_ID")
if [[ $test == test_freeze_region || $test == test_unfreeze_region ]]; then
 [[ -z $SECOND_REGION_ID ]] || REGIONS+=("$SECOND_REGION_ID")
 for region in "${REGIONS[@]}"; do
  "$MED_BIN" query staking region "$region" ${NET[@]+"${NET[@]}"} --output json >/dev/null
  query region-freeze "$region" | assert_json '.frozen == false'
 done
 DAO=$(key_address "$DAO_KEY")
 "$MED_BIN" query dao addresses ${NET[@]+"${NET[@]}"} --output json | assert_json --arg dao "$DAO" '.dao_addresses.global_dao == $dao'
 return 0
fi
DAO=$(key_address "$DAO_KEY")
OPERATOR=$(key_address "$OPERATOR_KEY")
OPERATOR_VAL=$("$MED_BIN" keys show "$OPERATOR_KEY" -a --bech val ${KEY[@]+"${KEY[@]}"})
printf '%s' "$REGION" | assert_json --arg operator "$OPERATOR_VAL" '(.region.operator_address // .operator_address) == $operator'
"$MED_BIN" query dao addresses ${NET[@]+"${NET[@]}"} --output json | assert_json --arg dao "$DAO" '.dao_addresses.global_dao == $dao'
[[ $DAO != "$OPERATOR" ]] || fail 'operator must not be DAO'
PERIOD=$(query params | jq -er '.params.min_voting_period')
NEW_OPERATOR=${NEW_OPERATOR:-$OPERATOR_VAL}
if [[ $test == test_submit_treasury_proposal || $test == test_submit_operator_proposal || $test == test_vote ]]; then
 return 0
fi
read -r -a TARGETS <<< "$EQUITY_ADDRESSES"
(( ${#TARGETS[@]} >= 2 )) || fail 'set EQUITY_ADDRESSES to at least two test addresses in config.sh or the environment'
TARGET_JSON=$(printf '%s\n' "${TARGETS[@]}" | jq -R . | jq -s .)
printf '%s' "$TARGET_JSON" | assert_json 'length == (unique | length)'
assert_frozen false
OLD_EQUITY=$(query equity "$REGION_ID" | jq -c '.equity.addresses // []')
# Refuse concurrent pending operations for this region, including later pages.
PAGE=1
while :; do
 ACTIONS=$(query pending-actions --status pending --limit 100 --page "$PAGE")
 printf '%s' "$ACTIONS" | assert_json --arg region "$REGION_ID" 'all(.actions[]?; .region_id != $region)'
 [[ $(printf '%s' "$ACTIONS" | jq '.actions | length') == 100 ]] || break
 PAGE=$((PAGE + 1))
done
PERIOD=$(query params | jq -er '.params.min_voting_period')
NEW_OPERATOR=${NEW_OPERATOR:-$OPERATOR_VAL}
echo "Preflight passed: chain=$CHAIN_ID region=$REGION_ID operator=$OPERATOR_KEY dao=$DAO_KEY"
}

setup_equity() {
 send_tx "$OPERATOR_KEY" submit-equity-update "$REGION_ID" "${TARGETS[@]}"
 approve "$(event_id action_id)"
 assert_equity "$TARGET_JSON"
}

restore_equity() {
 local restore=()
 while IFS= read -r address; do restore+=("$address"); done < <(printf '%s' "$OLD_EQUITY" | jq -r '.[]')
 send_tx "$OPERATOR_KEY" submit-equity-update "$REGION_ID" ${restore[@]+"${restore[@]}"}
 approve "$(event_id action_id)"
 assert_equity "$OLD_EQUITY"
}

reject_action() {
 local id=$1
 send_tx "$DAO_KEY" reject-pending-action "$id"
 action_status "$id" PENDING_ACTION_STATUS_REJECTED
}

test_submit_equity_update() {
 local id
 send_tx "$OPERATOR_KEY" submit-equity-update "$REGION_ID" "${TARGETS[@]}"
 id=$(event_id action_id)
 action_status "$id" PENDING_ACTION_STATUS_PENDING
 query pending-action "$id" | assert_json --argjson expected "$TARGET_JSON" '.action.addresses == $expected'
 assert_equity "$OLD_EQUITY"
 reject_action "$id"
 assert_equity "$OLD_EQUITY"
}

test_approve_pending_action() {
 setup_equity
 restore_equity
}

test_reject_pending_action() {
 local id
 send_tx "$OPERATOR_KEY" submit-equity-update "$REGION_ID" "${TARGETS[@]}"
 id=$(event_id action_id)
 reject_action "$id"
 assert_equity "$OLD_EQUITY"
}

test_submit_freeze_equity() {
 local id address
 setup_equity
 send_tx "$OPERATOR_KEY" submit-freeze-equity "$REGION_ID" "${TARGETS[@]}"
 id=$(event_id action_id)
 action_status "$id" PENDING_ACTION_STATUS_PENDING
 assert_frozen false
 approve "$id"
 assert_frozen true
 for address in "${TARGETS[@]}"; do send_tx "$DAO_KEY" unfreeze-equity "$REGION_ID" "$address"; done
 assert_frozen false
 restore_equity
}

test_submit_unfreeze_equity() {
 local id address
 setup_equity
 for address in "${TARGETS[@]}"; do send_tx "$DAO_KEY" freeze-equity "$REGION_ID" "$address"; done
 assert_frozen true
 send_tx "$OPERATOR_KEY" submit-unfreeze-equity "$REGION_ID" "${TARGETS[@]}"
 id=$(event_id action_id)
 action_status "$id" PENDING_ACTION_STATUS_PENDING
 assert_frozen true
 approve "$id"
 assert_frozen false
 restore_equity
}

assert_first_frozen() {
 query equity-freezes "$REGION_ID" | assert_json --arg a "${TARGETS[0]}" 'any(.freezes[]; .address == $a and .frozen)'
}

test_freeze_equity() {
 setup_equity
 send_tx "$DAO_KEY" freeze-equity "$REGION_ID" "${TARGETS[0]}"
 assert_first_frozen
 send_tx "$DAO_KEY" unfreeze-equity "$REGION_ID" "${TARGETS[0]}"
 assert_frozen false
 restore_equity
}

test_unfreeze_equity() {
 setup_equity
 send_tx "$DAO_KEY" freeze-equity "$REGION_ID" "${TARGETS[0]}"
 assert_first_frozen
 send_tx "$DAO_KEY" unfreeze-equity "$REGION_ID" "${TARGETS[0]}"
 assert_frozen false
 restore_equity
}

assert_regions_frozen() {
 local region
 for region in "${REGIONS[@]}"; do
  query region-freeze "$region" | assert_json --argjson expected "$1" '.frozen == $expected'
 done
}

test_freeze_region() {
 send_tx "$DAO_KEY" freeze-region "${REGIONS[@]}"
 assert_regions_frozen true
 send_tx "$DAO_KEY" unfreeze-region "${REGIONS[@]}"
 assert_regions_frozen false
}

test_unfreeze_region() {
 send_tx "$DAO_KEY" freeze-region "${REGIONS[@]}"
 assert_regions_frozen true
 send_tx "$DAO_KEY" unfreeze-region "${REGIONS[@]}"
 assert_regions_frozen false
}

submit_test_proposal() {
 local kind=$1
 jq -n --arg region "$REGION_ID" --arg receiver "$DAO" --arg operator "$NEW_OPERATOR" --arg period "$PERIOD" --arg kind "$kind" '
 {title:("med binary test " + $kind),summary:"Integration test: DAO veto; no execution",voting_period:$period} +
 (if $kind == "treasury" then {kind:"PROPOSAL_KIND_TREASURY_WITHDRAW",treasury:{withdrawals:[{region_id:$region,receiver:$receiver,amount:[{denom:"umec",amount:"1"}]}]}}
 else {kind:"PROPOSAL_KIND_REPLACE_REGION_OPERATOR",replace_operator:{region_id:$region,new_operator:$operator}} end)' > "$RESULT_DIR/$kind.json"
 send_tx "$DAO_KEY" submit-proposal "$RESULT_DIR/$kind.json"
 PROPOSAL_ID=$(event_id proposal_id)
 query proposal "$PROPOSAL_ID" | assert_json '.proposal.status == "PROPOSAL_STATUS_VOTING"'
 query electors "$PROPOSAL_ID" > "$RESULT_DIR/electors-$PROPOSAL_ID.json"
 query votes "$PROPOSAL_ID" | assert_json '(.votes // [] | length) == 0'
}

veto_test_proposal() {
 send_tx "$DAO_KEY" vote "$PROPOSAL_ID" no
 query proposal "$PROPOSAL_ID" | assert_json '.proposal.status == "PROPOSAL_STATUS_REJECTED"'
}

test_submit_treasury_proposal() {
 submit_test_proposal treasury
 veto_test_proposal
}

test_submit_operator_proposal() {
 submit_test_proposal replace-operator
 veto_test_proposal
}

test_vote() {
 submit_test_proposal treasury
 send_tx "$OPERATOR_KEY" vote "$PROPOSAL_ID" yes
 query votes "$PROPOSAL_ID" | assert_json --arg voter "$OPERATOR" 'any(.votes[]; .voter == $voter and .option == "VOTE_OPTION_YES")'
 veto_test_proposal
}

test_queries() {
 "$ROOT/queries.sh" all > "$RESULT_DIR/queries.jsonl"
 jq -se 'length == 9' "$RESULT_DIR/queries.jsonl" >/dev/null
}

run_case() (
 # A subshell isolates variables and traps between cases, not the chain itself.
 local name=$1 mode=$2
 echo "=== $name ==="
 preflight "$name"
 [[ $mode == --run ]] || { echo "PASS preflight: $name"; exit 0; }
 RESULT_DIR=$(mktemp -d "${TMPDIR:-/tmp}/regiongov.$name.XXXXXX")
 echo "Results: $RESULT_DIR"
 printf '%s\n' "${OLD_EQUITY:-[]}" > "$RESULT_DIR/original-equity.json"
 printf '%s\n' "${TARGET_JSON:-[]}" > "$RESULT_DIR/test-addresses.json"
 trap 'code=$?; if (( code != 0 )); then echo "FAIL: $name (exit $code). Inspect $RESULT_DIR and chain state; no blind rollback." >&2; fi' EXIT
 "$name"
 echo "PASS: $name"
)

main() {
 local mode=${1:---list} name selected
 case $mode in
  --list|--help)
   printf 'Usage: run.sh --list | --check [test_name ...|all] | --run test_name ...|all\n'
   printf '%s\n' "${TESTS[@]}"
   return ;;
  --check|--run) shift ;;
  *) fail "unknown mode: $mode" ;;
 esac
 if (( $# == 0 )); then
  [[ $mode == --check ]] || fail 'choose a test name or all explicitly'
  set -- all
 fi
 if [[ $1 == all ]]; then
  (( $# == 1 )) || fail 'all cannot be combined with test names'
  set -- "${TESTS[@]}"
 fi
 # Validate every requested name before any transaction is sent.
 for name in "$@"; do
  selected=false
  for candidate in "${TESTS[@]}"; do [[ $candidate != "$name" ]] || selected=true; done
  [[ $selected == true ]] || fail "unknown test: $name (use --list)"
 done
 for name in "$@"; do run_case "$name" "$mode"; done
}

if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
 main "$@"
fi

