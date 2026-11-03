#!/usr/bin/env bash
set -Eeuo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

prepare_equity() {
 read -r -a TARGETS <<< "$EQUITY_ADDRESSES"
 (( ${#TARGETS[@]} >= 2 )) || fail 'set EQUITY_ADDRESSES to at least two test addresses'
 TARGET_JSON=$(printf '%s\n' "${TARGETS[@]}" | jq -R . | jq -s .)
 OLD_EQUITY=$(query equity "$REGION_ID" | jq -c '.equity.addresses // []')
 printf '%s\n' "$OLD_EQUITY" > "$RESULT_DIR/original-equity.json"
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
 local kind=$1 DAO PERIOD
 DAO=$(key_address "$DAO_KEY")
 PERIOD=$(query params | jq -er '.params.min_voting_period')
 NEW_OPERATOR=${NEW_OPERATOR:-$("$MED_BIN" keys show "$OPERATOR_KEY" -a --bech val ${KEY[@]+"${KEY[@]}"})}
 jq -n --arg region "$REGION_ID" --arg receiver "$DAO" --arg operator "$NEW_OPERATOR" --arg period "$PERIOD" --arg kind "$kind" '
 {title:("med binary test " + $kind),summary:"med binary proposal test",voting_period:$period} +
 (if $kind == "treasury" then {kind:"PROPOSAL_KIND_TREASURY_WITHDRAW",treasury:{withdrawals:[{region_id:$region,receiver:$receiver,amount:[{denom:"umec",amount:"1"}]}]}}
 else {kind:"PROPOSAL_KIND_REPLACE_REGION_OPERATOR",replace_operator:{region_id:$region,new_operator:$operator}} end)' > "$RESULT_DIR/$kind.json"
 echo "Proposal JSON: $RESULT_DIR/$kind.json"
 cat "$RESULT_DIR/$kind.json"
 send_tx "$DAO_KEY" submit-proposal "$RESULT_DIR/$kind.json"
 PROPOSAL_ID=$(event_id proposal_id)
 query proposal "$PROPOSAL_ID" | assert_json '.proposal.status == "PROPOSAL_STATUS_VOTING"'
 query electors "$PROPOSAL_ID" > "$RESULT_DIR/electors-$PROPOSAL_ID.json"
 assert_json 'all(.electors[]; (.region_id | ascii_downcase) != "me_earth" and (.region_id | ascii_downcase) != "experience_region")' < "$RESULT_DIR/electors-$PROPOSAL_ID.json"
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
 local OPERATOR
 OPERATOR=$(key_address "$OPERATOR_KEY")
 submit_test_proposal treasury
 send_tx "$OPERATOR_KEY" vote "$PROPOSAL_ID" yes
 query votes "$PROPOSAL_ID" | assert_json --arg voter "$OPERATOR" 'any(.votes[]; .voter == $voter and .option == "VOTE_OPTION_YES")'
 veto_test_proposal
}

pass_test_proposal() {
 local kind=$1 voter status deadline period
 local voters=()
 read -r -a voters <<< "$VOTER_KEYS"
 for voter in "${voters[@]}"; do key_address "$voter" >/dev/null; done
 period=$(query params | jq -er '.params.min_voting_period')
 submit_test_proposal "$kind"
 for voter in "${voters[@]}"; do send_tx "$voter" vote "$PROPOSAL_ID" yes; done
 send_tx "$DAO_KEY" vote "$PROPOSAL_ID" yes
 deadline=$((SECONDS + period + TX_TIMEOUT))
 echo "Waiting for proposal $PROPOSAL_ID to pass (voting period: $period seconds)..."
 while (( SECONDS < deadline )); do
  status=$(query proposal "$PROPOSAL_ID" | jq -er '.proposal.status')
  case $status in
   PROPOSAL_STATUS_PASSED) query proposal "$PROPOSAL_ID" | jq .; return ;;
   PROPOSAL_STATUS_VOTING) sleep "$POLL_INTERVAL" ;;
   *) fail "proposal $PROPOSAL_ID ended with $status" ;;
  esac
 done
 fail "proposal $PROPOSAL_ID timed out; inspect chain before rerunning"
}

test_treasury_proposal_pass() { pass_test_proposal treasury; }
test_operator_proposal_pass() { pass_test_proposal replace-operator; }

test_queries() {
 "$ROOT/queries.sh" all > "$RESULT_DIR/queries.jsonl"
 jq -se 'length == 9' "$RESULT_DIR/queries.jsonl" >/dev/null
}

run_case() (
 echo "=== $1 ==="
 RESULT_DIR=$(mktemp -d "${TMPDIR:-/tmp}/regiongov.$1.XXXXXX")
 echo "Results: $RESULT_DIR"
 case $1 in
  test_*equity*|test_approve_pending_action|test_reject_pending_action) prepare_equity ;;
  test_freeze_region|test_unfreeze_region)
   REGIONS=("$REGION_ID")
   [[ -z $SECOND_REGION_ID ]] || REGIONS+=("$SECOND_REGION_ID") ;;
 esac
 "$1"
 echo "PASS: $1"
)

TESTS=(test_submit_equity_update test_approve_pending_action test_reject_pending_action test_submit_freeze_equity test_submit_unfreeze_equity test_freeze_equity test_unfreeze_equity test_freeze_region test_unfreeze_region test_submit_treasury_proposal test_submit_operator_proposal test_vote test_treasury_proposal_pass test_operator_proposal_pass test_queries)

main() {
 if (( $# == 0 )); then
  printf 'Usage: run.sh command_name\n'
  printf '%s\n' "${TESTS[@]}"
  return
 fi
 (( $# == 1 )) || fail 'specify one command name'
 case " ${TESTS[*]} " in
  *" $1 "*) run_case "$1" ;;
  *) fail "unknown command: $1" ;;
 esac
}

if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
 main "$@"
fi
