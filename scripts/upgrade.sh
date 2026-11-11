#!/bin/bash
# Usage: ./scripts/upgrade.sh propose v2.0.18.rc1
#        ./scripts/upgrade.sh vote 123
set -euo pipefail

HUB_VAL_KEYS=(val1 val2 val3 val4 val5 val6 val7 val8 val9 val10 val11 val12 val13 val14 val15 val16 val17 val18 val19 val20 val21 val22 val23 val24)

MED_BIN=${MED_BIN:-med}
NODE=${NODE:-http://118.175.0.230:26657}
ME_NODE1_HOME=${ME_NODE1_HOME:-./node1}
KEYRING_BACKEND_NAME=${KEYRING_BACKEND_NAME:-test}
ME_CHAIN_ID=${ME_CHAIN_ID:-}
HALT_OFFSET=${HALT_OFFSET:-10000}
GOV_FEES=${GOV_FEES:-20000umec}
VOTE_FEES=${VOTE_FEES:-20000umec}
GAS=${GAS:-auto}
GAS_ADJUSTMENT=${GAS_ADJUSTMENT:-1.5}
VOTE_OPTION=${VOTE_OPTION:-yes}
TX_WAIT_ATTEMPTS=${TX_WAIT_ATTEMPTS:-60}
TX_WAIT_INTERVAL=${TX_WAIT_INTERVAL:-2}

usage() {
    printf '%s\n' \
        'Usage: upgrade.sh propose <version> | vote <proposal-id>' \
        '  propose: val1 submits an SDK gov software-upgrade proposal.' \
        '  vote:    val1-val24 vote sequentially; default option is yes.' \
        'Defaults: NODE=http://118.175.0.230:26657, ME_NODE1_HOME=./node1 (relative to cwd),' \
        '          KEYRING_BACKEND_NAME=test, fees=20000umec, GAS=auto.' \
        'ME_CHAIN_ID defaults to the RPC chain ID; an explicit value must match it.' \
        'GOV_DEPOSIT defaults to the minimum deposit queried from gov params.' \
        'HALT_HEIGHT sets the upgrade height; otherwise current height + HALT_OFFSET (10000).' \
        'Choose a height AFTER the governance voting period, allowing time to deploy binaries.' \
        'VOTE_OPTION may be yes, no, abstain, or no_with_veto.' \
        'UPGRADE_INFO optionally provides a JSON object containing binary download information.' \
        'TX_WAIT_ATTEMPTS=60 and TX_WAIT_INTERVAL=2 control transaction polling.' \
        'Transactions are never automatically rebroadcast; errors stop the script.'
}

fail() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }

positive_integer() {
    [[ $1 =~ ^[1-9][0-9]*$ && ${#1} -le 15 ]] || fail "$2 must be a positive integer (at most 15 digits)"
}

query() {
    "$MED_BIN" query "$@" --node "$NODE" --home "$ME_NODE1_HOME" --chain-id "$ME_CHAIN_ID" -o json
}

prepare() {
    command -v "$MED_BIN" >/dev/null 2>&1 || fail "command not found: $MED_BIN"
    command -v jq >/dev/null 2>&1 || fail 'jq is required'
    [[ -d $ME_NODE1_HOME ]] || fail "node home does not exist: $ME_NODE1_HOME"
    positive_integer "$TX_WAIT_ATTEMPTS" TX_WAIT_ATTEMPTS
    [[ $TX_WAIT_INTERVAL =~ ^[0-9]+$ ]] || fail 'TX_WAIT_INTERVAL must be a non-negative integer'
    local status remote_chain
    status=$("$MED_BIN" status --node "$NODE" --home "$ME_NODE1_HOME")
    remote_chain=$(jq -er '.NodeInfo.network // .node_info.network | select(type == "string" and length > 0)' <<< "$status")
    [[ -z $ME_CHAIN_ID || $ME_CHAIN_ID == "$remote_chain" ]] || fail "chain ID mismatch: configured=$ME_CHAIN_ID RPC=$remote_chain"
    ME_CHAIN_ID=$remote_chain
    CURRENT_HEIGHT=$(jq -er '.SyncInfo.latest_block_height // .sync_info.latest_block_height | tostring' <<< "$status")
    positive_integer "$CURRENT_HEIGHT" 'RPC block height'
    RESULT_DIR=$(mktemp -d "${TMPDIR:-/tmp}/me-hub-upgrade.XXXXXX")
    printf 'RPC=%s chain=%s home=%s\nTransaction results: %s\n' "$NODE" "$ME_CHAIN_ID" "$ME_NODE1_HOME" "$RESULT_DIR"
}

require_key() {
    "$MED_BIN" keys show "$1" --address --home "$ME_NODE1_HOME" \
        --keyring-backend "$KEYRING_BACKEND_NAME" >/dev/null || fail "key unavailable: $1"
}

# Both CheckTx and committed DeliverTx must explicitly report success.
# A timeout/transport error is ambiguous: preserve the response and never retry a broadcast.
broadcast() {
    local key=$1 fees=$2 label=$3 response hash attempt committed
    shift 3
    response="$RESULT_DIR/$label-broadcast.json"
    printf 'Broadcasting %s with %s\n' "$label" "$key"
    if ! "$MED_BIN" tx gov "$@" --from "$key" \
        --node "$NODE" --home "$ME_NODE1_HOME" --chain-id "$ME_CHAIN_ID" \
        --keyring-backend "$KEYRING_BACKEND_NAME" --fees "$fees" \
        --gas "$GAS" --gas-adjustment "$GAS_ADJUSTMENT" \
        --broadcast-mode sync --yes -o json > "$response"; then
        fail "broadcast failed or outcome unknown; inspect $response before retrying"
    fi
    if ! jq -e '(.tx_response // .) | (.code != null and (.code | tonumber) == 0)' "$response" >/dev/null; then
        jq . "$response" >&2 || true
        fail "CheckTx failed or response invalid: $response"
    fi
    hash=$(jq -er '(.tx_response // .).txhash | select(type == "string" and test("^[A-Fa-f0-9]{64}$"))' "$response")
    printf 'txhash=%s\n' "$hash"
    committed="$RESULT_DIR/$label-committed.json"
    for ((attempt = 1; attempt <= TX_WAIT_ATTEMPTS; attempt++)); do
        if query tx "$hash" > "$committed" 2> "$RESULT_DIR/$label-query.err"; then
            if ! jq -e --arg hash "$hash" '(.tx_response // .) |
                (.txhash | ascii_upcase) == ($hash | ascii_upcase) and
                (.height | tonumber) > 0 and (.code != null and (.code | tonumber) == 0)' "$committed" >/dev/null; then
                jq . "$committed" >&2 || true
                fail "DeliverTx failed or response invalid: $committed (do not blindly rebroadcast)"
            fi
            COMMITTED_TX=$committed
            return
        fi
        if ((attempt < TX_WAIT_ATTEMPTS)); then sleep "$TX_WAIT_INTERVAL"; fi
    done
    fail "tx $hash not confirmed; inspect $RESULT_DIR and query it manually; no rebroadcast performed"
}

propose() {
    local version=$1 halt deposit info proposal_id
    require_key val1
    if [[ -n ${HALT_HEIGHT:-} ]]; then
        positive_integer "$HALT_HEIGHT" HALT_HEIGHT
        halt=$HALT_HEIGHT
    else
        positive_integer "$HALT_OFFSET" HALT_OFFSET
        halt=$((CURRENT_HEIGHT + HALT_OFFSET))
    fi
    ((halt > CURRENT_HEIGHT)) || fail "upgrade height $halt must be above current height $CURRENT_HEIGHT"
    deposit=${GOV_DEPOSIT:-}
    if [[ -z $deposit ]]; then
        query gov params > "$RESULT_DIR/gov-params.json"
        deposit=$(jq -er '(.params.min_deposit // .deposit_params.min_deposit // .min_deposit) |
            select(type == "array" and length > 0) | map(.amount + .denom) | join(",")' "$RESULT_DIR/gov-params.json")
    fi
    info=${UPGRADE_INFO:-'{}'}
    jq -e 'type == "object"' <<< "$info" >/dev/null || fail 'UPGRADE_INFO must be a JSON object'
    printf 'Proposing %s from val1: current height=%s upgrade height=%s deposit=%s\n' "$version" "$CURRENT_HEIGHT" "$halt" "$deposit"
    printf 'Ensure the upgrade height is after the voting period and allows time to deploy the new binary.\n'
    broadcast val1 "$GOV_FEES" propose submit-legacy-proposal software-upgrade "$version" \
        --title "Upgrade to $version" --description "Upgrade ME-Hub to $version" \
        --upgrade-height "$halt" --upgrade-info "$info" --no-validate --deposit "$deposit"
    proposal_id=$(jq -er '
        def decoded: try @base64d catch "";
        (.tx_response // .) |
        [.logs[]?.events[]?, .events[]?] |
        [ .[] | select(.type == "submit_proposal" or .type == "message") |
          .attributes[]? |
          if .key == "proposal_id" then .value
          elif (.key | decoded) == "proposal_id" then (.value | decoded)
          else empty end ] | unique |
        select(length == 1) | .[0] | select(test("^[1-9][0-9]*$"))
    ' "$COMMITTED_TX") || fail "proposal committed but ID could not be extracted; inspect $COMMITTED_TX; do not resubmit"
    printf 'proposal_id=%s\n' "$proposal_id"
    query gov proposal "$proposal_id"
}

vote() {
    local proposal_id=$1 key status
    # Fail before any broadcast if a configured key is absent. No key generation or funding.
    for key in "${HUB_VAL_KEYS[@]}"; do require_key "$key"; done
    query gov proposal "$proposal_id" > "$RESULT_DIR/proposal-before.json"
    status=$(jq -er '(.proposal // .).status | tostring' "$RESULT_DIR/proposal-before.json")
    [[ $status == PROPOSAL_STATUS_VOTING_PERIOD || $status == 2 ]] || fail "proposal $proposal_id is not in voting period: $status"
    for key in "${HUB_VAL_KEYS[@]}"; do
        broadcast "$key" "$VOTE_FEES" "vote-$key" vote "$proposal_id" "$VOTE_OPTION"
    done
    query gov proposal "$proposal_id"
    printf 'All configured keys voted %s on proposal %s. Wait for governance tally; voting does not execute the upgrade immediately.\n' "$VOTE_OPTION" "$proposal_id"
}

main() {
    case ${1:-} in
        -h|--help|help) usage; return ;;
        propose)
            [[ $# -eq 2 ]] || { usage >&2; exit 1; }
            [[ $2 =~ ^[A-Za-z0-9][A-Za-z0-9._-]*$ ]] || fail 'invalid upgrade version'
            ;;
        vote)
            [[ $# -eq 2 ]] || { usage >&2; exit 1; }
            positive_integer "$2" 'proposal ID'
            case $VOTE_OPTION in yes|no|abstain|no_with_veto) ;; *) fail 'invalid VOTE_OPTION' ;; esac
            ;;
        *) usage >&2; exit 1 ;;
    esac
    prepare
    "$1" "$2"
}

main "$@"
