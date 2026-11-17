#!/usr/bin/env bash
# Offline regression tests: public test keys only; no queries or transactions.
set -Eeuo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/run.sh"

TEST_DIR=$(mktemp -d "${TMPDIR:-/tmp}/regiongov-run-test.XXXXXX")
printf 'Offline test results: %s\n' "$TEST_DIR"
# Ensure an accidental query/broadcast fails instead of touching a live chain.
MED_UNDER_TEST=$(command -v "$MED_BIN")
med_offline() {
 case "$1" in
  keys) "$MED_UNDER_TEST" "$@" ;;
  *) fail "offline test forbids med command: $*" ;;
 esac
}
MED_BIN=med_offline
CURRENT_OPERATOR=mevaloper1different
current_region_operator() { printf '%s\n' "$CURRENT_OPERATOR"; }
# A DAO-submitted replacement must not look up this unrelated key.
OPERATOR_KEY=nonexistent-operator-key
NEW_OPERATOR=
NEW_OPERATOR_INDEX=23
RESULT_DIR="$TEST_DIR/default"
mkdir "$RESULT_DIR"

require_new_operator
operator_23=$NEW_OPERATOR
account_23=$(jq -er '.address' "$RESULT_DIR/new-operator-key.json")
[[ $operator_23 == mevaloper1* ]] || fail 'expected a validator operator address'
[[ $account_23 == me1* ]] || fail 'expected an account address'
[[ $("$MED_BIN" keys show regiongov-operator-23 -a --bech val \
 --home "$RESULT_DIR/derived-keys" --keyring-backend test) == "$operator_23" ]] \
 || fail 'the replacement signing key was not saved'
"$MED_BIN" keys show regiongov-operator-23 --output json \
 --home "$RESULT_DIR/derived-keys" --keyring-backend test \
 | assert_json --arg address "$account_23" '.type == "local" and .address == $address'
# A second preflight reuses the existing key/address without an overwrite prompt.
require_new_operator
"$MED_BIN" keys list --output json --home "$RESULT_DIR/derived-keys" --keyring-backend test \
 | assert_json 'length == 1'
printf 'PASS: recovered signing key, reusable preflight, no OPERATOR_KEY dependency\n'

# A fresh result directory with the same index derives the same identity.
NEW_OPERATOR=
RESULT_DIR="$TEST_DIR/repeated"
mkdir "$RESULT_DIR"
require_new_operator
[[ $NEW_OPERATOR == "$operator_23" ]] || fail 'HD derivation must be deterministic'
printf 'PASS: same index preserves the operator identity\n'

NEW_OPERATOR=
NEW_OPERATOR_INDEX=26
RESULT_DIR="$TEST_DIR/next-index"
mkdir "$RESULT_DIR"
require_new_operator
operator_26=$NEW_OPERATOR
[[ $operator_26 != "$operator_23" ]] || fail 'different indices must derive different operators'
"$MED_BIN" keys show regiongov-operator-26 --output json \
 --home "$RESULT_DIR/derived-keys" --keyring-backend test \
 | assert_json '.type == "local"'
printf 'PASS: configurable replacement index saves another signing key\n'

# Explicit addresses must not trigger recovery or write to any keyring.
NEW_OPERATOR=$operator_26
RESULT_DIR="$TEST_DIR/explicit"
mkdir "$RESULT_DIR"
require_new_operator
[[ ! -e $RESULT_DIR/derived-keys && ! -e $RESULT_DIR/new-operator-key.json ]] \
 || fail 'explicit NEW_OPERATOR unexpectedly recovered a key'
printf 'PASS: explicit address does not recover a key\n'

CURRENT_OPERATOR=$operator_26
if (require_new_operator > "$TEST_DIR/same-operator.log" 2>&1); then
 fail 'accepted replacement with the current operator'
fi
grep -q 'NEW_OPERATOR equals current operator' "$TEST_DIR/same-operator.log" \
 || fail 'missing actionable same-operator error'
printf 'PASS: current operator cannot be reused as the replacement\n'

NEW_OPERATOR=
NEW_OPERATOR_INDEX=invalid
RESULT_DIR="$TEST_DIR/invalid-index"
mkdir "$RESULT_DIR"
if (require_new_operator > "$TEST_DIR/invalid-index.log" 2>&1); then
 fail 'accepted an invalid HD index'
fi
grep -q 'NEW_OPERATOR_INDEX must be a nonnegative HD index' "$TEST_DIR/invalid-index.log" \
 || fail 'missing invalid-index error'
[[ ! -e $RESULT_DIR/derived-keys ]] || fail 'invalid index wrote to a keyring'
printf 'PASS: invalid index is rejected before recovery\n'
printf 'All offline operator-key tests passed; no chain requests were sent.\n'
