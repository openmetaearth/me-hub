#!/usr/bin/env bash
set -Eeuo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

recover_test_key() {
 # Public test mnemonic only; never use these accounts for real funds.
 local global_dict="sign wise lyrics find brass grass already urge invite vibrant arrow grow dune theme trust mix napkin ball pause sponsor bacon capable flip people"
 local index=$1 name=$2
 shift 2
 printf '%s\n' "$global_dict" | "$MED_BIN" keys add "$name" \
  --recover --no-backup --key-type secp256k1 --coin-type 118 \
  --hd-path "m/44'/118'/0'/0/0/$index" --keyring-backend test \
  --home "$RESULT_DIR/derived-keys" --output json "$@"
}

derive_test_address() {
 recover_test_key "$1" "$2" --dry-run | jq -er '.address | select(type == "string" and startswith("me1"))'
}

generate_equity_addresses() {
 local index address
 TARGETS=()
 # setup_4validator.sh uses 11-14 for val1-val4; reserve 21-22 for equity tests.
 for index in 21 22; do
  address=$(derive_test_address "$index" "regiongov-equity-$index") \
   || fail "could not derive equity address at index $index"
  TARGETS+=("$address")
 done
}

generate_operator_address() {
 local address
 # Reserve index 23 for the replacement operator, separate from validators/equity.
 address=$(derive_test_address 23 regiongov-operator-23) || fail 'could not derive operator address at index 23'
 # Re-encode with the validator prefix AND checksum; replacing "me1" is invalid.
 NEW_OPERATOR=$("$MED_BIN" debug addr "$address" --home "$RESULT_DIR/derived-keys" \
  | jq -Rer 'select(startswith("Bech32 Val: ")) | ltrimstr("Bech32 Val: ") | select(startswith("mevaloper1"))') \
  || fail 'could not convert the derived account to a validator operator address'
 printf 'Derived operator account: %s (HD index 23)\n' "$address"
}

prepare_equity() {
 local key address
 if [[ -n $EQUITY_ADDRESSES ]]; then
  read -r -a TARGETS <<< "$EQUITY_ADDRESSES"
 else
  generate_equity_addresses
 fi
 (( ${#TARGETS[@]} >= 2 )) || fail 'set EQUITY_ADDRESSES to at least two test addresses'
 TARGET_JSON=$(printf '%s\n' "${TARGETS[@]}" | jq -R . | jq -s .)
 printf '%s\n' "$TARGET_JSON" | assert_json 'map(ascii_downcase) | length == (unique | length)'
 for key in val1 val2 val3 val4 "$OPERATOR_KEY"; do
  address=$(key_address "$key") || fail "could not read operator key $key"
  if printf '%s\n' "$TARGET_JSON" | jq -e --arg address "$address" 'map(ascii_downcase) | index($address | ascii_downcase) != null' >/dev/null; then
   fail "equity addresses must not include operator key $key ($address)"
  fi
 done
 printf '%s\n' "$TARGET_JSON" > "$RESULT_DIR/equity-addresses.json"
 printf 'Equity addresses: %s\n' "$TARGET_JSON"
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

prepare_equity_treasury() {
 [[ $REGION_ID == chn ]] || fail 'this test requires REGION_ID=chn'
 local index address dao
 EQUITY_KEY_OPTIONS=("${KEY[@]}")
 if [[ -z $EQUITY_ADDRESSES ]]; then
  local addresses=()
  # Three equities: one proposes and the other two vote (2/3 exceeds 51%).
  # Index 23 remains reserved for operator replacement.
  for index in 21 22 24; do
   address=$(recover_test_key "$index" "regiongov-equity-$index" | jq -er '.address | select(type == "string" and startswith("me1"))') \
    || fail "could not recover equity test key at index $index"
   addresses+=("$address")
  done
  EQUITY_ADDRESSES="${addresses[*]}"
  # Keep the node's client configuration, but sign with an isolated test keyring.
  EQUITY_KEY_OPTIONS=(--home ./nodes/hub-nodes/node1 --keyring-dir "$RESULT_DIR/derived-keys" --keyring-backend test)
 fi
 prepare_equity
 (( ${#TARGETS[@]} >= 3 )) || fail 'set EQUITY_ADDRESSES to at least three test addresses with keys in the configured keyring'
 dao=$(key_address "$DAO_KEY")
 for address in "${TARGETS[@]}"; do
  [[ $address != "$dao" ]] || fail 'equity test accounts must differ from the DAO'
  "$MED_BIN" keys show "$address" -a "${EQUITY_KEY_OPTIONS[@]}" >/dev/null \
   || fail "no signing key for equity address $address"
 done
}

send_equity_tx() {
 local KEY=("${EQUITY_KEY_OPTIONS[@]}")
 send_tx "$@"
}

expect_equity_unauthorized() {
 local signer=$1 reason=$2
 shift 2
 # Auto gas would fail during simulation, without a committed permission check.
 [[ $REJECT_GAS =~ ^[1-9][0-9]*$ ]] || fail 'REJECT_GAS must be a positive fixed gas limit'
 local GAS=$REJECT_GAS
 local KEY=("${EQUITY_KEY_OPTIONS[@]}")
 send_module_tx_expect_result "$signer" regiongov 2 regiongov "$@"
 printf '%s' "$TX_JSON" | assert_json --arg reason "$reason" '(.raw_log // "" | contains($reason))'
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
}

test_unfreeze_region() {
 send_tx "$DAO_KEY" unfreeze-region "${REGIONS[@]}"
 assert_regions_frozen false
}

submit_test_proposal() {
 local kind=$1 DAO PERIOD
 DAO=$(key_address "$DAO_KEY")
 PERIOD=$(query params | jq -er '.params.min_voting_period')
 if [[ $kind == replace-operator ]]; then require_new_operator; fi
 jq -n --arg region "$REGION_ID" --arg receiver "$DAO" --arg operator "$NEW_OPERATOR" --arg period "$PERIOD" --arg kind "$kind" '
 {title:("med binary test " + $kind),summary:"med binary proposal test",voting_period:$period} +
 (if $kind == "treasury" then {kind:"PROPOSAL_KIND_TREASURY_WITHDRAW",treasury:{withdrawals:[{region_id:$region,receiver:$receiver,amount:[{denom:"umec",amount:"10"}]}]}}
 elif $kind == "replace-operator" then {kind:"PROPOSAL_KIND_REPLACE_REGION_OPERATOR",replace_operator:{region_id:$region,new_operator:$operator}}
 elif $kind == "update-min-voting-period" then {kind:"PROPOSAL_KIND_UPDATE_MIN_VOTING_PERIOD",update_min_voting_period:{min_voting_period:"121"}}
 else error("unknown proposal kind: " + $kind) end)' > "$RESULT_DIR/$kind.json"
 echo "Proposal JSON: $RESULT_DIR/$kind.json"
 cat "$RESULT_DIR/$kind.json"
 send_tx "$DAO_KEY" submit-proposal "$RESULT_DIR/$kind.json"
 PROPOSAL_ID=$(event_id proposal_id)
 query proposal "$PROPOSAL_ID" | assert_json '.proposal.status == "PROPOSAL_STATUS_VOTING"'
 if [[ $kind == update-min-voting-period ]]; then
  query proposal "$PROPOSAL_ID" | assert_json '
   .proposal.kind == "PROPOSAL_KIND_UPDATE_MIN_VOTING_PERIOD" and
   (.proposal.voting_period | tonumber) == 120 and
   (.proposal.update_min_voting_period.min_voting_period | tonumber) == 121'
  # Submission alone must not change the minimum before voting/execution.
  query params | assert_json '(.params.min_voting_period | tonumber) == 120'
 fi
 query electors "$PROPOSAL_ID" > "$RESULT_DIR/electors-$PROPOSAL_ID.json"
 assert_json '(.electors | length) > 0 and all(.electors[]; (.operator | startswith("me1")) and (.region_id | ascii_downcase) != "me_earth" and (.region_id | ascii_downcase) != "experience_region")' < "$RESULT_DIR/electors-$PROPOSAL_ID.json"
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
 assert_json --arg voter "$OPERATOR" 'any(.electors[]; .operator == $voter)' < "$RESULT_DIR/electors-$PROPOSAL_ID.json"
 send_tx "$OPERATOR_KEY" vote "$PROPOSAL_ID" yes
 query votes "$PROPOSAL_ID" | assert_json --arg voter "$OPERATOR" 'any(.votes[]; .voter == $voter and .option == "VOTE_OPTION_YES")'
 veto_test_proposal
}

pass_test_proposal() {
 local kind=$1 voter period
 local voters=()
 read -r -a voters <<< "$VOTER_KEYS"
 for voter in "${voters[@]}"; do key_address "$voter" >/dev/null; done
 period=$(query params | jq -er '.params.min_voting_period')
 submit_test_proposal "$kind"
 for voter in "${voters[@]}"; do send_tx "$voter" vote "$PROPOSAL_ID" yes; done
 send_tx "$DAO_KEY" vote "$PROPOSAL_ID" yes
 wait_for_proposal_pass "$period"
}

wait_for_proposal_pass() { wait_for_proposal_status "$1" PROPOSAL_STATUS_PASSED; }

wait_for_proposal_status() {
 local period=$1 expected=$2 status deadline
 deadline=$((SECONDS + period + TX_TIMEOUT))
 echo "Waiting for proposal $PROPOSAL_ID to reach $expected (voting period: $period seconds)..."
 while (( SECONDS < deadline )); do
  status=$(query proposal "$PROPOSAL_ID" | jq -er '.proposal.status')
  if [[ $status == "$expected" ]]; then
   query proposal "$PROPOSAL_ID" | jq .
   return
  fi
  [[ $status == PROPOSAL_STATUS_VOTING ]] || fail "proposal $PROPOSAL_ID ended with $status, expected $expected"
  sleep "$POLL_INTERVAL"
 done
 fail "proposal $PROPOSAL_ID timed out; inspect chain before rerunning"
}

test_treasury_proposal_pass() { pass_test_proposal treasury; }

test_equity_treasury_proposal_pass() {
 local params period percent threshold minimum treasury dao operator proposer voter
 local voters=("${TARGETS[@]:1}")
 proposer=${TARGETS[0]}
 params=$(query params)
 period=$(printf '%s' "$params" | jq -er '.params.min_voting_period')
 percent=$(printf '%s' "$params" | jq -er '.params.equity_withdraw_percent | select(. >= 1 and . <= 100)')
 threshold=$(printf '%s' "$params" | jq -er '.params.pass_threshold_percent')
 (( ${#voters[@]} * 100 >= ${#TARGETS[@]} * threshold )) \
  || fail "the other equity accounts cannot reach the $threshold% threshold without the proposer"
 dao=$(key_address "$DAO_KEY")

 # A 10umec withdrawal must also fit the equity percentage cap. Do not silently
 # fund the treasury: fail before changing equity state if test funds are missing.
 "$MED_BIN" query staking region "$REGION_ID" "${NET[@]}" --output json > "$RESULT_DIR/region.json"
 treasury=$(jq -er '.region.region_treasure_addr | select(length > 0)' "$RESULT_DIR/region.json")
 operator=$(key_address "$OPERATOR_KEY")
 "$MED_BIN" keys show "$OPERATOR_KEY" --bech val -a "${KEY[@]}" \
  | jq -Re --slurpfile region "$RESULT_DIR/region.json" '. == $region[0].region.operator_address' >/dev/null \
  || fail 'OPERATOR_KEY must be the current chn operator'
 "$MED_BIN" query bank spendable-balances "$treasury" --denom umec "${NET[@]}" --output json > "$RESULT_DIR/treasury-before.json"
 minimum=$(((10 * 100 + percent - 1) / percent))
 assert_json --argjson minimum "$minimum" '(.balance.amount // "0" | tonumber) >= $minimum' < "$RESULT_DIR/treasury-before.json"

 # Fund only the isolated public test accounts so they can pay transaction fees.
 for voter in "${TARGETS[@]}"; do
  send_module_tx "$DAO_KEY" bank send "$dao" "$voter" "$FUND_AMOUNT"
 done
 setup_equity
 jq -n --arg region "$REGION_ID" --arg receiver "$dao" --arg period "$period" '{
  title:"chn equity withdrawal approved by other equity accounts",
  summary:"An equity account proposes; other equities vote, with separate region operator and DAO approvals",
  kind:"PROPOSAL_KIND_TREASURY_WITHDRAW",voting_period:$period,
  treasury:{withdrawals:[{region_id:$region,receiver:$receiver,amount:[{denom:"umec",amount:"10"}]}]}
 }' > "$RESULT_DIR/equity-treasury.json"
 send_equity_tx "$proposer" submit-proposal "$RESULT_DIR/equity-treasury.json"
 PROPOSAL_ID=$(event_id proposal_id)
 query proposal "$PROPOSAL_ID" > "$RESULT_DIR/equity-proposal.json"
 assert_json --arg proposer "$proposer" --arg operator "$operator" '
  .proposal.proposer == $proposer and .proposal.status == "PROPOSAL_STATUS_VOTING" and
  .proposal.electorate == "ELECTORATE_KIND_REGION_EQUITY" and
  .proposal.required_operator == $operator
 ' < "$RESULT_DIR/equity-proposal.json"
 query electors "$PROPOSAL_ID" > "$RESULT_DIR/electors-$PROPOSAL_ID.json"
 assert_json --argjson targets "$TARGET_JSON" '
  ([.electors[].operator] | sort) == ($targets | sort) and all(.electors[]; .region_id == "chn")
 ' < "$RESULT_DIR/electors-$PROPOSAL_ID.json"

 # Equity ballots form the percentage; operator/DAO approvals are separate.
 for voter in "${voters[@]}"; do send_equity_tx "$voter" vote "$PROPOSAL_ID" yes; done
 send_tx "$OPERATOR_KEY" vote "$PROPOSAL_ID" yes
 send_tx "$DAO_KEY" vote "$PROPOSAL_ID" yes
 query votes "$PROPOSAL_ID" --limit "$((${#TARGETS[@]} + 1))" > "$RESULT_DIR/equity-votes.json"
 assert_json --arg proposer "$proposer" --arg operator "$operator" --argjson voters "$(printf '%s\n' "${voters[@]}" | jq -R . | jq -s .)" '
  (.votes // []) as $votes |
  all($voters[]; . as $voter | any($votes[]; .voter == $voter and .option == "VOTE_OPTION_YES")) and
  any($votes[]; .voter == $operator and .option == "VOTE_OPTION_YES") and
  all($votes[]; .voter != $proposer)
 ' < "$RESULT_DIR/equity-votes.json"
 wait_for_proposal_pass "$period"
 restore_equity
}

test_equity_treasury_other_region_rejected() {
 local period treasury receiver dao voter before_count receiver_before
 [[ $OTHER_REGION_ID != "$REGION_ID" && $OTHER_REGION_ID != me_earth && $OTHER_REGION_ID != experience_region ]] \
  || fail 'OTHER_REGION_ID must name a different country region'
 [[ $REJECT_GAS =~ ^[1-9][0-9]*$ ]] || fail 'REJECT_GAS must be a positive fixed gas limit'
 period=$(query params | jq -er '.params.min_voting_period')
 dao=$(key_address "$DAO_KEY")
 receiver=$(derive_test_address 25 regiongov-cross-region-receiver) || fail 'could not derive withdrawal receiver'
 for voter in "${TARGETS[@]}"; do
  [[ $voter != "$receiver" ]] || fail 'test receiver must not be an equity signer'
 done
 [[ $receiver != "$dao" ]] || fail 'test receiver must not be the fee-paying DAO'

 # Require an existing, funded foreign treasury and chn-only equity memberships.
 # Do not modify the other region's equity list or silently fund its treasury.
 "$MED_BIN" query staking region "$OTHER_REGION_ID" "${NET[@]}" --output json > "$RESULT_DIR/other-region.json"
 treasury=$(jq -er '.region.region_treasure_addr | select(length > 0)' "$RESULT_DIR/other-region.json")
 "$MED_BIN" query bank spendable-balances "$treasury" --denom umec "${NET[@]}" --output json \
  | assert_json '(.balance.amount // "0" | tonumber) >= 10'
 query equity "$OTHER_REGION_ID" | assert_json --argjson targets "$TARGET_JSON" '
  (.equity.addresses // [] | map(ascii_downcase)) as $foreign |
  all($targets[]; ascii_downcase as $a | $foreign | index($a) == null)'
 for voter in "$REGION_ID" "$OTHER_REGION_ID"; do
  query region-freeze "$voter" | assert_json '.frozen == false'
 done

 for voter in "${TARGETS[@]}"; do
  send_module_tx "$DAO_KEY" bank send "$dao" "$voter" "$FUND_AMOUNT"
 done
 setup_equity
 assert_frozen false
 "$MED_BIN" query bank balances "$receiver" --denom umec "${NET[@]}" --output json > "$RESULT_DIR/receiver-before.json"
 receiver_before=$(jq -er '.amount' "$RESULT_DIR/receiver-before.json")
 jq -n --arg region "$OTHER_REGION_ID" --arg receiver "$receiver" --arg period "$period" '{
  title:"chn equities cannot withdraw from another region",
  summary:"Foreign equity submission and ballots must be rejected without a treasury transfer",
  kind:"PROPOSAL_KIND_TREASURY_WITHDRAW",voting_period:$period,
  treasury:{withdrawals:[{region_id:$region,receiver:$receiver,amount:[{denom:"umec",amount:"10"}]}]}
 }' > "$RESULT_DIR/other-region-treasury.json"

 before_count=$(query proposals --count-total --limit 1 | jq -er '.pagination.total')
 expect_equity_unauthorized "${TARGETS[0]}" "an equity address can withdraw only from its own region, got $OTHER_REGION_ID" \
  submit-proposal "$RESULT_DIR/other-region-treasury.json"
 query proposals --count-total --limit 1 | assert_json --arg count "$before_count" '.pagination.total == $count'

 # A valid DAO proposal lets us test ballots, rather than voting on a nonexistent ID.
 # It retains the existing country-operator electorate, never foreign equities.
 send_tx "$DAO_KEY" submit-proposal "$RESULT_DIR/other-region-treasury.json"
 PROPOSAL_ID=$(event_id proposal_id)
 query proposal "$PROPOSAL_ID" | assert_json --arg dao "$dao" --arg region "$OTHER_REGION_ID" '
  .proposal.proposer == $dao and .proposal.status == "PROPOSAL_STATUS_VOTING" and
  .proposal.electorate == "ELECTORATE_KIND_COUNTRY_REGION_OPERATORS" and
  .proposal.treasury.withdrawals[0].region_id == $region'
 query electors "$PROPOSAL_ID" > "$RESULT_DIR/other-region-electors.json"
 assert_json '(.electors | length) > 0' < "$RESULT_DIR/other-region-electors.json"
 query votes "$PROPOSAL_ID" --count-total --limit 1 | assert_json '(.pagination.total | tonumber) == 0'
 for voter in "${TARGETS[@]}"; do
  expect_equity_unauthorized "$voter" "$voter is not an elector of proposal $PROPOSAL_ID" vote "$PROPOSAL_ID" yes
  query votes "$PROPOSAL_ID" --count-total --limit 1 | assert_json '(.pagination.total | tonumber) == 0'
  query proposal "$PROPOSAL_ID" | assert_json '.proposal.status == "PROPOSAL_STATUS_VOTING"'
 done
 send_tx "$DAO_KEY" vote "$PROPOSAL_ID" yes
 query proposal "$PROPOSAL_ID" | assert_json '.proposal.status == "PROPOSAL_STATUS_VOTING" and .proposal.dao_vote == "DAO_VOTE_YES"'
 wait_for_proposal_status "$period" PROPOSAL_STATUS_REJECTED
 query proposal "$PROPOSAL_ID" | assert_json '.proposal.dao_vote == "DAO_VOTE_YES"'
 query votes "$PROPOSAL_ID" --count-total --limit 1 | assert_json '(.pagination.total | tonumber) == 0'
 # Treasury rewards may accrue during the wait; use a non-signing receiver to
 # prove no withdrawal happened, instead of requiring a static treasury balance.
 "$MED_BIN" query bank balances "$receiver" --denom umec "${NET[@]}" --output json \
  | assert_json --arg amount "$receiver_before" '.amount == $amount'
 restore_equity
}

test_min_voting_period_proposal_pass() {
 local before current
 before=$(query params)
 printf '%s\n' "$before" > "$RESULT_DIR/params-before.json"
 current=$(printf '%s' "$before" | jq -er '.params.min_voting_period')
 [[ $current == 120 ]] || fail "this test requires min_voting_period=120, got $current; it leaves the minimum at 121 after success"
 pass_test_proposal update-min-voting-period
 query params > "$RESULT_DIR/params-after.json"
 assert_json --argjson before "$before" '
  (.params.min_voting_period | tonumber) == 121 and
  (.params | del(.min_voting_period)) == ($before.params | del(.min_voting_period))
 ' < "$RESULT_DIR/params-after.json"
 jq . "$RESULT_DIR/params-after.json"
 printf 'min_voting_period updated: 120 -> 121 seconds (not restored).\n'
}

require_new_operator() {
 [[ -n $NEW_OPERATOR ]] || generate_operator_address
 local current
 current=$("$MED_BIN" keys show "$OPERATOR_KEY" -a --bech val ${KEY[@]+"${KEY[@]}"})
 [[ $(printf '%s' "$NEW_OPERATOR" | tr '[:upper:]' '[:lower:]') != $(printf '%s' "$current" | tr '[:upper:]' '[:lower:]') ]] \
  || fail 'NEW_OPERATOR must differ from the current operator; if index 23 was already used, set NEW_OPERATOR to another unused address'
 printf '%s\n' "$NEW_OPERATOR" > "$RESULT_DIR/new-operator.txt"
 printf 'New operator: %s\n' "$NEW_OPERATOR"
}

test_operator_proposal_pass() { pass_test_proposal replace-operator; }

test_treasury_proposal_51pass_daoveto() {
 local voter params threshold period eligible required yes end_time deadline block_time
 local voters=()
 read -r -a voters <<< "$VOTER_KEYS"
 for voter in "${voters[@]}"; do key_address "$voter" >/dev/null; done
 params=$(query params)
 threshold=$(printf '%s' "$params" | jq -er '.params.pass_threshold_percent')
 [[ $threshold == 51 ]] || fail "this test requires pass_threshold_percent=51, got $threshold"
 period=$(printf '%s' "$params" | jq -er '.params.min_voting_period')

 submit_test_proposal treasury
 eligible=$(jq -er '.electors | length | select(. > 0)' "$RESULT_DIR/electors-$PROPOSAL_ID.json")
 # Reach the 51% threshold with the smallest whole number of operator votes.
 required=$(((eligible * threshold + 99) / 100))
 (( ${#voters[@]} >= required )) || fail "set VOTER_KEYS to at least $required eligible operator keys ($eligible electors)"
 echo "Voting YES with $required/$eligible operators (threshold: $threshold%)..."
 yes=0
 for voter in "${voters[@]}"; do
  (( yes < required )) || break
  send_tx "$voter" vote "$PROPOSAL_ID" yes
  yes=$((yes + 1))
 done
 query votes "$PROPOSAL_ID" --limit "$eligible" > "$RESULT_DIR/operator-yes-votes.json"
 assert_json --argjson eligible "$eligible" --argjson threshold "$threshold" '
  ([.votes[]? | select(.option == "VOTE_OPTION_YES")] | length) * 100 >= $eligible * $threshold
 ' < "$RESULT_DIR/operator-yes-votes.json"

 send_tx "$DAO_KEY" vote "$PROPOSAL_ID" yes
 query proposal "$PROPOSAL_ID" > "$RESULT_DIR/dao-yes.json"
 assert_json '.proposal.status == "PROPOSAL_STATUS_VOTING" and .proposal.dao_vote == "DAO_VOTE_YES"' < "$RESULT_DIR/dao-yes.json"
 end_time=$(jq -er '.proposal.voting_end_time | sub("\\.[0-9]+Z$"; "Z") | fromdateiso8601' "$RESULT_DIR/dao-yes.json")

 # DAO may overwrite YES with NO before voting closes, rejecting immediately.
 veto_test_proposal
 query proposal "$PROPOSAL_ID" > "$RESULT_DIR/dao-veto.json"
 assert_json '.proposal.status == "PROPOSAL_STATUS_REJECTED" and .proposal.dao_vote == "DAO_VOTE_NO"' < "$RESULT_DIR/dao-veto.json"

 deadline=$((SECONDS + period + TX_TIMEOUT))
 echo "DAO vetoed proposal $PROPOSAL_ID; waiting past its voting end time to verify it stays rejected..."
 while (( SECONDS < deadline )); do
  block_time=$("$MED_BIN" query block ${NET[@]+"${NET[@]}"} | jq -er '.block.header.time | sub("\\.[0-9]+Z$"; "Z") | fromdateiso8601')
  # Use committed chain time, not the local clock; wait past fractional seconds too.
  if (( block_time > end_time )); then
   query proposal "$PROPOSAL_ID" > "$RESULT_DIR/after-voting-end.json"
   assert_json '.proposal.status == "PROPOSAL_STATUS_REJECTED" and .proposal.dao_vote == "DAO_VOTE_NO"' < "$RESULT_DIR/after-voting-end.json"
   jq . "$RESULT_DIR/after-voting-end.json"
   return
  fi
  sleep "$POLL_INTERVAL"
 done
 fail "timeout waiting past proposal $PROPOSAL_ID voting end time; inspect chain before rerunning"
}

test_queries() {
 "$ROOT/queries.sh" all > "$RESULT_DIR/queries.jsonl"
 jq -se 'length == 9' "$RESULT_DIR/queries.jsonl" >/dev/null
}

run_case() (
 echo "=== $1 ==="
 RESULT_DIR=$(mktemp -d "${TMPDIR:-/tmp}/regiongov.$1.XXXXXX")
 echo "Results: $RESULT_DIR"
 case $1 in
  test_equity_treasury_proposal_pass|test_equity_treasury_other_region_rejected) prepare_equity_treasury ;;
  test_*equity*|test_approve_pending_action|test_reject_pending_action) prepare_equity ;;
  test_freeze_region|test_unfreeze_region)
   REGIONS=("$REGION_ID")
   [[ -z $SECOND_REGION_ID ]] || REGIONS+=("$SECOND_REGION_ID") ;;
 esac
 "$1"
 echo "PASS: $1"
)

TESTS=(test_submit_equity_update test_approve_pending_action test_reject_pending_action test_submit_freeze_equity test_submit_unfreeze_equity test_freeze_equity test_unfreeze_equity test_freeze_region test_unfreeze_region test_submit_treasury_proposal test_submit_operator_proposal test_vote test_treasury_proposal_pass test_treasury_proposal_51pass_daoveto test_equity_treasury_proposal_pass test_equity_treasury_other_region_rejected test_min_voting_period_proposal_pass test_operator_proposal_pass test_queries)

run_all() {
 local RESULT_DIR
 RESULT_DIR=$(mktemp -d "${TMPDIR:-/tmp}/regiongov.all.XXXXXX")
 printf 'Preflight results: %s\n' "$RESULT_DIR"
 require_new_operator
 # Keep the same generated/explicit address in every fresh test process.
 export NEW_OPERATOR
 local test result=0 passed=0 failed=0 skipped=0 index=0 started=$SECONDS
 local results=()
 for test in "${TESTS[@]}"; do
  index=$((index + 1))
  if (( failed > 0 )); then
   results+=("SKIP  $test")
   skipped=$((skipped + 1))
   continue
  fi
  printf '\n[%d/%d] %s\n' "$index" "${#TESTS[@]}" "$test"
  # Use a fresh Bash process: testing run_case in an if would disable its errexit.
  if "$BASH" "$ROOT/run.sh" "$test"; then
   results+=("PASS  $test")
   passed=$((passed + 1))
  else
   result=$?
   results+=("FAIL  $test (exit $result)")
   failed=$((failed + 1))
  fi
 done
 printf '\n=== Test summary ===\n'
 printf '%s\n' "${results[@]}"
 printf 'Total: %d, passed: %d, failed: %d, skipped: %d, elapsed: %ds\n' \
  "${#TESTS[@]}" "$passed" "$failed" "$skipped" "$((SECONDS - started))"
 if (( failed > 0 )); then
  printf 'Stopped after the first failure; inspect the chain and test results before rerunning.\n' >&2
 fi
 return "$result"
}

main() {
 if (( $# == 0 )); then
  printf 'Usage: run.sh <all|command_name>\n'
  printf '  all: run every test in order, stopping on the first failure\n'
  printf '%s\n' "${TESTS[@]}"
  return
 fi
 (( $# == 1 )) || fail 'specify one command name'
 if [[ $1 == all ]]; then
  run_all
  return
 fi
 case " ${TESTS[*]} " in
  *" $1 "*) run_case "$1" ;;
  *) fail "unknown command: $1" ;;
 esac
}

if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
 main "$@"
fi
