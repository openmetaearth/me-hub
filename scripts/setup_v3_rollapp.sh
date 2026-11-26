#!/bin/bash
# Register a new rollapp on a Hub that is already running v3.
# init creates a new rollapp home, the sequencer account, and the dymint key.
# The later commands register that rollapp and sequencer on the Hub.
# This script does not start the rollapp process or open an IBC channel.
#
# The rollapp chain id must use revision 1, for example myrollapp_402-1.
# The new home is separate from the existing rollapp-nodes/node1.
#
#   cd /data/docker-beta
#   ROLLAPP_CHAIN_ID=myrollapp_402-1 \
#   ./scripts/setup_v3_rollapp.sh all
#
#   ./scripts/setup_v3_rollapp.sh init|rollapp|fund|sequencer|show
set -euo pipefail

PROJECT_ENV="${PROJECT_ENV:-beta}"
BASE_DIR="${BASE_DIR:-/data/docker-${PROJECT_ENV}}"
BIN_DIR="${BIN_DIR:-${BASE_DIR}/bin}"
KEYRING_BACKEND_NAME="${KEYRING_BACKEND_NAME:-test}"
ME_CHAIN_ID="${ME_CHAIN_ID:-mechain_400-1}"
ME_NODE1_HOME="${ME_NODE1_HOME:-${BASE_DIR}/nodes/hub-nodes/node1}"
DENOM="${DENOM:-umec}"
FEES="${FEES:-2000000${DENOM}}"
ROLLAPP_DENOM="${ROLLAPP_DENOM:-urax}"
ROLLAPP_DISPLAY_DENOM="${ROLLAPP_DISPLAY_DENOM:-RAX}"

ROLLAPP_CHAIN_ID="${ROLLAPP_CHAIN_ID:-}"
ROLLAPP_ALIAS="${ROLLAPP_ALIAS:-}"
ROLLAPP_HOME="${ROLLAPP_HOME:-}"
VM_TYPE="${VM_TYPE:-WASM}"
BECH32_PREFIX="${BECH32_PREFIX:-me}"
# Must be greater than 0. A zero supply makes the Hub store an empty native denom.
INITIAL_SUPPLY="${INITIAL_SUPPLY:-1000000000000000000}"
GENESIS_CHECKSUM="${GENESIS_CHECKSUM:-}"
ROLLAPP_GENESIS="${ROLLAPP_GENESIS:-}"
NATIVE_DENOM_FILE="${NATIVE_DENOM_FILE:-}"

CREATOR="${CREATOR:-global_dao}"
SEQUENCER_KEY="${SEQUENCER_KEY:-sequencer}"
SEQUENCER_MNEMONIC="${SEQUENCER_MNEMONIC:-}"
SEQUENCER_MONIKER="${SEQUENCER_MONIKER:-sequencer}"
SEQUENCER_RPC="${SEQUENCER_RPC:-http://127.0.0.1:26657}"
SEQUENCER_REST="${SEQUENCER_REST:-http://127.0.0.1:1317}"
DA_BASE_URL="${DA_BASE_URL:-http://127.0.0.1:26658}"
SEQUENCER_PUBKEY_FILE="${SEQUENCER_PUBKEY_FILE:-}"
# Integer in base-denom units. Empty means read rollapp params min_sequencer_bond_global.
# create-rollapp --min-sequencer-bond takes this integer and treats it as umec, not whole MEC.
MIN_SEQUENCER_BOND="${MIN_SEQUENCER_BOND:-}"
SEQUENCER_BOND="${SEQUENCER_BOND:-}"

usage() {
    cat <<EOF
Usage: $0 <command>

  all         Init the rollapp, create it on the Hub, fund and register the sequencer, then query
  init        rollappd init, create the sequencer account, and export the dymint pubkey
  rollapp     ${CREATOR} creates ${ROLLAPP_CHAIN_ID:-<ROLLAPP_CHAIN_ID>}
  fund        Send bond and fees from ${CREATOR} when the sequencer balance is short
  sequencer   Register the dymint pubkey from sequencer.info
  show        Query the rollapp and its sequencers

Env:
  BASE_DIR=${BASE_DIR}
  ME_CHAIN_ID=${ME_CHAIN_ID}
  ROLLAPP_CHAIN_ID=${ROLLAPP_CHAIN_ID:-<required>}
  ROLLAPP_ALIAS=${ROLLAPP_ALIAS:-<name part of the chain id>}
  ROLLAPP_HOME=${ROLLAPP_HOME:-<rollapp-nodes/alias>}
  VM_TYPE=${VM_TYPE}
  BECH32_PREFIX=${BECH32_PREFIX}
  ROLLAPP_DENOM=${ROLLAPP_DENOM}
  INITIAL_SUPPLY=${INITIAL_SUPPLY}
  GENESIS_CHECKSUM=${GENESIS_CHECKSUM:-<sha256 of the new genesis>}
  ROLLAPP_GENESIS=${ROLLAPP_GENESIS:-<ROLLAPP_HOME/config/genesis.json>}
  NATIVE_DENOM_FILE=${NATIVE_DENOM_FILE:-<unset>}
  CREATOR=${CREATOR}
  SEQUENCER_KEY=${SEQUENCER_KEY}
  SEQUENCER_MNEMONIC=${SEQUENCER_MNEMONIC:-<generated on init>}
  MIN_SEQUENCER_BOND=${MIN_SEQUENCER_BOND:-<on-chain global minimum>}
  SEQUENCER_BOND=${SEQUENCER_BOND:-<MIN_SEQUENCER_BOND>${DENOM}}
EOF
}

need_base_dir() {
    if [ ! -d "${BASE_DIR}" ]; then
        echo "error: BASE_DIR not found: ${BASE_DIR}" >&2
        exit 1
    fi
    cd "${BASE_DIR}"
    export PATH="${BIN_DIR}:${PATH}"
}

need_rollapp_id() {
    if [ -z "${ROLLAPP_CHAIN_ID}" ]; then
        echo "error: set ROLLAPP_CHAIN_ID, for example myrollapp_402-1" >&2
        exit 1
    fi
    if ! echo "${ROLLAPP_CHAIN_ID}" | grep -Eq '^[a-z]+_[1-9][0-9]*-1$'; then
        echo "error: rollapp id must be <lowercase-name>_<eip155>-1, got ${ROLLAPP_CHAIN_ID}" >&2
        exit 1
    fi
    if [ -z "${ROLLAPP_ALIAS}" ]; then
        ROLLAPP_ALIAS="${ROLLAPP_CHAIN_ID%%_*}"
    fi
    if [ -z "${ROLLAPP_HOME}" ]; then
        ROLLAPP_HOME="${BASE_DIR}/nodes/rollapp-nodes/${ROLLAPP_ALIAS}"
    fi
    if [ -z "${ROLLAPP_GENESIS}" ]; then
        ROLLAPP_GENESIS="${ROLLAPP_HOME}/config/genesis.json"
    fi
    if [ -z "${SEQUENCER_PUBKEY_FILE}" ]; then
        SEQUENCER_PUBKEY_FILE="${ROLLAPP_HOME}/sequencer.info"
    fi
}

check_tx_status() {
    local cosmosapp=$1
    local txhash=$2
    local extra_args=$3
    local sleepTime=120
    if [ -z "${txhash}" ] || [ "${txhash}" = "null" ]; then
        echo "error: empty txhash" >&2
        exit 1
    fi
    echo "Check tx status: ${txhash}"
    while true; do
        if ! ${cosmosapp} query tx "${txhash}" ${extra_args} >/dev/null 2>&1; then
            sleep 1
            sleepTime=$((sleepTime - 1))
            if [ "${sleepTime}" -eq 0 ]; then
                echo "Tx not found: ${txhash}" >&2
                ${cosmosapp} query tx "${txhash}" ${extra_args} || true
                exit 1
            fi
            continue
        fi
        local resp_code
        resp_code=$(${cosmosapp} q tx "${txhash}" -o json ${extra_args} | jq -r '.code')
        if [ "${resp_code}" != "0" ]; then
            echo "Tx failed, code: ${resp_code}" >&2
            ${cosmosapp} q tx "${txhash}" ${extra_args} || true
            exit 1
        fi
        echo "Tx succeeded, code: ${resp_code}"
        break
    done
}

hub_tx() {
    local from=$1
    shift
    local keyring_args=()
    if [ -n "${HUB_KEYRING_DIR:-}" ]; then
        keyring_args+=(--keyring-dir "${HUB_KEYRING_DIR}")
    fi
    med tx "$@" \
        --from "${from}" \
        --keyring-backend "${KEYRING_BACKEND_NAME}" \
        --chain-id "${ME_CHAIN_ID}" \
        --home "${ME_NODE1_HOME}" \
        "${keyring_args[@]}" \
        --fees "${FEES}" \
        --gas auto --gas-adjustment 1.5 \
        -y -o json
}

# dymint ComputeGenesisHash: unmarshal, remarshal with sorted keys, then sha256.
# A raw file hash does not match the checksum the rollapp sends at startup.
file_sha256() {
    python3 - "$1" <<'PY'
import hashlib, json, sys
with open(sys.argv[1], encoding="utf-8") as fh:
    obj = json.load(fh)
canonical = json.dumps(obj, sort_keys=True, separators=(",", ":"), ensure_ascii=False)
canonical = canonical.replace("&", r"\u0026").replace("<", r"\u003c").replace(">", r"\u003e")
print(hashlib.sha256(canonical.encode("utf-8")).hexdigest())
PY
}

resolve_genesis_checksum() {
    if [ -n "${GENESIS_CHECKSUM}" ]; then
        return 0
    fi
    if [ ! -f "${ROLLAPP_GENESIS}" ]; then
        echo "error: set GENESIS_CHECKSUM, or point ROLLAPP_GENESIS at the rollapp genesis file" >&2
        exit 1
    fi
    GENESIS_CHECKSUM="$(file_sha256 "${ROLLAPP_GENESIS}")"
    echo "genesis checksum ${GENESIS_CHECKSUM}"
}

resolve_min_bond() {
    if [ -z "${MIN_SEQUENCER_BOND}" ]; then
        local params
        params=$(med q rollapp params --home "${ME_NODE1_HOME}" -o json)
        MIN_SEQUENCER_BOND=$(echo "${params}" | jq -r '
            def amt:
                if . == null then empty
                elif type == "object" then (.amount // empty)
                elif type == "string" then (capture("^(?<n>[0-9]+)") | .n)
                else empty end;
            (.params.minSequencerBondGlobal
                // .params.min_sequencer_bond_global
                // .minSequencerBondGlobal
                // .min_sequencer_bond_global) | amt
        ')
    fi
    if [ -z "${MIN_SEQUENCER_BOND}" ] || [ "${MIN_SEQUENCER_BOND}" = "null" ]; then
        echo "error: could not read min_sequencer_bond_global; set MIN_SEQUENCER_BOND to the integer umec amount" >&2
        echo "med q rollapp params --home ${ME_NODE1_HOME} -o json" >&2
        exit 1
    fi
    if [ -z "${SEQUENCER_BOND}" ]; then
        SEQUENCER_BOND="${MIN_SEQUENCER_BOND}${DENOM}"
    fi
}

key_addr() {
    local bin=$1
    local name=$2
    local home=$3
    "${bin}" keys show "${name}" -a \
        --keyring-backend "${KEYRING_BACKEND_NAME}" --home "${home}" 2>/dev/null || true
}

sequencer_home() {
    echo "${ROLLAPP_HOME}/sequencer_keys"
}

sequencer_addr() {
    key_addr rollappd "${SEQUENCER_KEY}" "$(sequencer_home)"
}

set_toml_value() {
    local file=$1
    local key=$2
    local value=$3
    if grep -q "^${key} *=" "${file}"; then
        sed -i "s|^${key} *=.*|${key} = ${value}|" "${file}"
    else
        echo "${key} = ${value}" >> "${file}"
    fi
}

# NewNode indexes DAConfig[i] for every DALayer entry and panics when da_config is empty.
ensure_da_config() {
    local dymint=$1
    python3 - "${dymint}" "${DA_BASE_URL}" <<'PY'
import json, sys
path, base = sys.argv[1], sys.argv[2]
lines = open(path).read().splitlines()
def value_of(prefix):
    for line in lines:
        stripped = line.strip()
        if stripped.startswith(prefix):
            return stripped.split("=", 1)[1].strip()
    return None
current = value_of("da_config")
if current not in (None, '""', "''", "[]", ""):
    sys.exit(0)
cfg = {
    "base_url": base,
    "timeout": 50000000000,
    "gas_prices": 0.1,
    "auth_token": "TOKEN",
    "backoff": {"initial_delay": 6000000000, "max_delay": 6000000000, "growth_factor": 2},
    "retry_attempts": 4,
    "retry_delay": 3000000000,
}
layer = 'da_layer = ["me-da"]'
config = "da_config = ['%s']" % json.dumps(cfg, separators=(",", ":"))
out, seen_layer, seen_config = [], False, False
for line in lines:
    stripped = line.strip()
    if stripped.startswith("da_layer"):
        out.append(layer)
        seen_layer = True
    elif stripped.startswith("da_config"):
        out.append(config)
        seen_config = True
    else:
        out.append(line)
if not seen_layer:
    out.append(layer)
if not seen_config:
    out.append(config)
open(path, "w").write("\n".join(out) + "\n")
print("set da_layer=[\"me-da\"] and a matching da_config")
PY
}

# Rollapp gas, staking, and gov use urax. hub-genesis looks up metadata for the bond denom during InitGenesis.
prepare_genesis() {
    local genesis=$1
    local tmp
    tmp=$(mktemp)
    jq --arg denom "${ROLLAPP_DENOM}" --arg display "${ROLLAPP_DISPLAY_DENOM}" '
        .app_state.staking.params.bond_denom = $denom
        | (if .app_state.mint.minter then .app_state.mint.minter.mint_denom = $denom else . end)
        | (if .app_state.mint.params.mint_denom then .app_state.mint.params.mint_denom = $denom else . end)
        | (if .app_state.gov.deposit_params.min_deposit[0] then .app_state.gov.deposit_params.min_deposit[0].denom = $denom else . end)
        | (if .app_state.crisis.constant_fee then .app_state.crisis.constant_fee.denom = $denom else . end)
        | (if .app_state.rollappparams.params then .app_state.rollappparams.params.da = "me-da" else . end)
        | .app_state.bank.denom_metadata = [
            {
                "description": ("Denom metadata for " + $display + " (" + $denom + ")"),
                "denom_units": [
                    {"denom": $denom, "exponent": 0, "aliases": []},
                    {"denom": $display, "exponent": 18, "aliases": []}
                ],
                "base": $denom,
                "display": $display,
                "name": $display,
                "symbol": $display
            }
        ]
    ' "${genesis}" > "${tmp}" && mv "${tmp}" "${genesis}"
    local app_toml="${ROLLAPP_HOME}/config/app.toml"
    if [ -f "${app_toml}" ]; then
        set_toml_value "${app_toml}" minimum-gas-prices "\"0${ROLLAPP_DENOM}\""
    fi
    echo "rollapp denom set to ${ROLLAPP_DENOM}"
    if [ "${INITIAL_SUPPLY}" = "0" ]; then
        echo "error: INITIAL_SUPPLY must be greater than 0 so the Hub keeps native denom ${ROLLAPP_DENOM}" >&2
        exit 1
    fi
    if [ -d "${ROLLAPP_HOME}/data" ]; then
        echo "genesis changed; reset chain data before the next start:"
        echo "  rollappd dymint unsafe-reset-all --home ${ROLLAPP_HOME}"
    fi
}

# hub-genesis reports this balance as the native-denom initial supply.
ensure_genesis_supply() {
    local addr=$1
    local genesis="${ROLLAPP_HOME}/config/genesis.json"
    local have
    have=$(jq -r --arg addr "${addr}" --arg denom "${ROLLAPP_DENOM}" '
        [.app_state.bank.balances[]? | select(.address == $addr) | .coins[]? | select(.denom == $denom) | .amount] | first // "0"
    ' "${genesis}")
    if [ "${have}" = "${INITIAL_SUPPLY}" ]; then
        echo "genesis supply ${INITIAL_SUPPLY}${ROLLAPP_DENOM} already belongs to ${addr}"
        return 0
    fi
    if [ "${have}" != "0" ]; then
        echo "error: ${addr} already has ${have}${ROLLAPP_DENOM} in genesis, expected ${INITIAL_SUPPLY}" >&2
        exit 1
    fi
    echo "# add-genesis-account ${addr} ${INITIAL_SUPPLY}${ROLLAPP_DENOM}"
    rollappd add-genesis-account "${addr}" "${INITIAL_SUPPLY}${ROLLAPP_DENOM}" \
        --home "${ROLLAPP_HOME}" \
        --keyring-backend "${KEYRING_BACKEND_NAME}"
}

write_native_denom() {
    NATIVE_DENOM_FILE="${ROLLAPP_HOME}/native_denom.json"
    jq -n --arg base "${ROLLAPP_DENOM}" --arg display "${ROLLAPP_DISPLAY_DENOM}" \
        '{base: $base, display: $display, exponent: 18}' > "${NATIVE_DENOM_FILE}"
}

cmd_init() {
    need_base_dir
    need_rollapp_id
    if ! command -v rollappd >/dev/null 2>&1; then
        echo "error: rollappd is not on PATH (${BIN_DIR})" >&2
        exit 1
    fi

    local genesis="${ROLLAPP_HOME}/config/genesis.json"
    if [ -f "${genesis}" ]; then
        local existing
        existing=$(jq -r '.chain_id' "${genesis}")
        if [ "${existing}" != "${ROLLAPP_CHAIN_ID}" ]; then
            echo "error: ${ROLLAPP_HOME} is already chain ${existing}" >&2
            exit 1
        fi
        echo "rollapp home already initialized: ${ROLLAPP_HOME}"
    else
        echo "# rollappd init ${ROLLAPP_ALIAS} --chain-id ${ROLLAPP_CHAIN_ID}"
        mkdir -p "${ROLLAPP_HOME}"
        rollappd init "${ROLLAPP_ALIAS}" --chain-id "${ROLLAPP_CHAIN_ID}" --home "${ROLLAPP_HOME}"
    fi

    local seq_home roll_addr
    seq_home="$(sequencer_home)"
    roll_addr="$(sequencer_addr)"
    if [ -z "${roll_addr}" ]; then
        mkdir -p "${seq_home}"
        echo "# create rollapp sequencer key ${SEQUENCER_KEY} in ${seq_home}"
        if [ -n "${SEQUENCER_MNEMONIC}" ]; then
            echo "${SEQUENCER_MNEMONIC}" | rollappd keys add "${SEQUENCER_KEY}" \
                --keyring-backend "${KEYRING_BACKEND_NAME}" \
                --home "${seq_home}" \
                --recover
        else
            rollappd keys add "${SEQUENCER_KEY}" \
                --keyring-backend "${KEYRING_BACKEND_NAME}" \
                --home "${seq_home}"
        fi
        roll_addr="$(sequencer_addr)"
    fi
    if [ -z "${roll_addr}" ]; then
        echo "error: rollapp sequencer key ${SEQUENCER_KEY} was not created" >&2
        exit 1
    fi
    echo "sequencer address: ${roll_addr}"

    local dymint="${ROLLAPP_HOME}/config/dymint.toml"
    set_toml_value "${dymint}" keyring_home_dir "\"${seq_home}\""
    set_toml_value "${dymint}" dym_account_name "\"sequencer\""
    ensure_da_config "${dymint}"
    prepare_genesis "${ROLLAPP_HOME}/config/genesis.json"
    ensure_genesis_supply "${roll_addr}"

    echo "# export dymint pubkey to ${SEQUENCER_PUBKEY_FILE}"
    rollappd dymint show-sequencer --home "${ROLLAPP_HOME}" > "${SEQUENCER_PUBKEY_FILE}"
    cat "${SEQUENCER_PUBKEY_FILE}"
}

rollapp_exists() {
    med q rollapp show "${ROLLAPP_CHAIN_ID}" --home "${ME_NODE1_HOME}" >/dev/null 2>&1
}

sequencer_exists() {
    local addr=$1
    med q sequencer show-sequencer "${addr}" --home "${ME_NODE1_HOME}" >/dev/null 2>&1
}

# create-rollapp stores the genesis checksum. A later genesis edit has to be written back before the rollapp is sealed.
sync_registered_checksum() {
    local show onchain sealed launched txhash
    show=$(med q rollapp show "${ROLLAPP_CHAIN_ID}" --home "${ME_NODE1_HOME}" -o json)
    onchain=$(echo "${show}" | jq -r '
        .rollapp.genesisInfo.genesisChecksum
        // .rollapp.genesis_info.genesis_checksum
        // .rollapp.genesisInfo.genesis_checksum
        // .rollapp.genesis_info.genesisChecksum
        // empty')
    local onchain_supply onchain_base
    onchain_supply=$(echo "${show}" | jq -r '
        .rollapp.genesisInfo.initialSupply
        // .rollapp.genesis_info.initial_supply
        // .rollapp.genesisInfo.initial_supply
        // "0"')
    onchain_base=$(echo "${show}" | jq -r '
        .rollapp.genesisInfo.nativeDenom.base
        // .rollapp.genesis_info.native_denom.base
        // .rollapp.genesisInfo.native_denom.base
        // empty')
    sealed=$(echo "${show}" | jq -r '.rollapp.genesisInfo.sealed // .rollapp.genesis_info.sealed // false')
    launched=$(echo "${show}" | jq -r '.rollapp.launched // false')
    if [ "${onchain}" = "${GENESIS_CHECKSUM}" ] && [ "${onchain_supply}" = "${INITIAL_SUPPLY}" ] && [ "${onchain_base}" = "${ROLLAPP_DENOM}" ]; then
        echo "hub genesis info matches checksum ${GENESIS_CHECKSUM}"
        return 0
    fi
    echo "hub genesis info checksum=${onchain:-<empty>} supply=${onchain_supply} denom=${onchain_base:-<empty>}"
    echo "local genesis info checksum=${GENESIS_CHECKSUM} supply=${INITIAL_SUPPLY} denom=${ROLLAPP_DENOM}"
    if [ "${sealed}" = "true" ] || [ "${launched}" = "true" ]; then
        echo "error: genesis info is sealed, so the hub record cannot be updated" >&2
        exit 1
    fi
    write_native_denom
    echo "# update-rollapp ${ROLLAPP_CHAIN_ID} genesis info"
    txhash=$(hub_tx "${CREATOR}" rollapp update-rollapp "${ROLLAPP_CHAIN_ID}" \
        --genesis-checksum "${GENESIS_CHECKSUM}" \
        --bech32-prefix "${BECH32_PREFIX}" \
        --initial-supply "${INITIAL_SUPPLY}" \
        --native-denom "${NATIVE_DENOM_FILE}" | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
}

cmd_rollapp() {
    need_base_dir
    need_rollapp_id
    resolve_genesis_checksum
    resolve_min_bond
    local seq_addr
    seq_addr="$(sequencer_addr)"
    if rollapp_exists; then
        echo "rollapp already registered: ${ROLLAPP_CHAIN_ID}"
        sync_registered_checksum
        med q rollapp show "${ROLLAPP_CHAIN_ID}" --home "${ME_NODE1_HOME}"
        return 0
    fi
    write_native_denom
    echo "# create-rollapp ${ROLLAPP_CHAIN_ID} alias=${ROLLAPP_ALIAS} vm=${VM_TYPE} sequencer=${seq_addr}"
    local txhash
    txhash=$(hub_tx "${CREATOR}" rollapp create-rollapp \
        "${ROLLAPP_CHAIN_ID}" "${ROLLAPP_ALIAS}" "${VM_TYPE}" \
        --init-sequencer "${seq_addr}" \
        --min-sequencer-bond "${MIN_SEQUENCER_BOND}" \
        --genesis-checksum "${GENESIS_CHECKSUM}" \
        --bech32-prefix "${BECH32_PREFIX}" \
        --initial-supply "${INITIAL_SUPPLY}" \
        --native-denom "${NATIVE_DENOM_FILE}" | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    med q rollapp show "${ROLLAPP_CHAIN_ID}" --home "${ME_NODE1_HOME}"
}

coin_amount() {
    echo "$1" | sed "s/${DENOM}$//"
}

cmd_fund() {
    need_base_dir
    need_rollapp_id
    resolve_min_bond
    local seq_addr bal bond fee need
    seq_addr="$(sequencer_addr)"
    bal=$(med q bank balances "${seq_addr}" --home "${ME_NODE1_HOME}" -o json \
        | jq -r --arg denom "${DENOM}" '.balances[]? | select(.denom == $denom) | .amount')
    bal="${bal:-0}"
    bond="$(coin_amount "${SEQUENCER_BOND}")"
    fee="$(coin_amount "${FEES}")"
    need=$(python3 - "${bal}" "${bond}" "${fee}" <<'PY'
import sys
bal, bond, fee = (int(x) for x in sys.argv[1:])
gap = bond + fee - bal
print(gap if gap > 0 else 0)
PY
)
    if [ "${need}" = "0" ]; then
        echo "sequencer ${seq_addr} already has ${bal}${DENOM}"
        return 0
    fi
    echo "# bank send ${need}${DENOM} ${CREATOR} -> ${SEQUENCER_KEY}"
    local txhash
    txhash=$(hub_tx "${CREATOR}" bank send "${CREATOR}" "${seq_addr}" "${need}${DENOM}" | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    med q bank balances "${seq_addr}" --home "${ME_NODE1_HOME}"
}

cmd_sequencer() {
    need_base_dir
    need_rollapp_id
    resolve_min_bond
    if ! rollapp_exists; then
        echo "error: rollapp ${ROLLAPP_CHAIN_ID} is not registered" >&2
        exit 1
    fi
    local seq_addr pubkey meta txhash
    seq_addr="$(sequencer_addr)"
    if sequencer_exists "${seq_addr}"; then
        echo "sequencer already registered: ${seq_addr}"
        med q sequencer show-sequencer "${seq_addr}" --home "${ME_NODE1_HOME}"
        return 0
    fi
    if [ ! -f "${SEQUENCER_PUBKEY_FILE}" ]; then
        echo "error: sequencer pubkey file not found: ${SEQUENCER_PUBKEY_FILE}" >&2
        exit 1
    fi
    pubkey="$(tr -d '\n' < "${SEQUENCER_PUBKEY_FILE}")"
    meta="$(mktemp)"
    jq -n --arg moniker "${SEQUENCER_MONIKER}" --arg rpc "${SEQUENCER_RPC}" --arg rest "${SEQUENCER_REST}" \
        '{moniker: $moniker, details: "", rpcs: [$rpc], rest_api_urls: [$rest]}' > "${meta}"
    echo "# create-sequencer ${seq_addr} bond=${SEQUENCER_BOND}"
    txhash=$(HUB_KEYRING_DIR="$(sequencer_home)" hub_tx "${SEQUENCER_KEY}" sequencer create-sequencer \
        "${pubkey}" "${ROLLAPP_CHAIN_ID}" "${SEQUENCER_BOND}" "${meta}" | jq -r '.txhash')
    rm -f "${meta}"
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    med q sequencer show-sequencer "${seq_addr}" --home "${ME_NODE1_HOME}"
}

cmd_show() {
    need_base_dir
    need_rollapp_id
    echo "# rollapp ${ROLLAPP_CHAIN_ID}"
    med q rollapp show "${ROLLAPP_CHAIN_ID}" --home "${ME_NODE1_HOME}"
    echo "# sequencers"
    med q sequencer show-sequencers-by-rollapp "${ROLLAPP_CHAIN_ID}" --home "${ME_NODE1_HOME}"
}

cmd_all() {
    cmd_init
    cmd_rollapp
    cmd_fund
    cmd_sequencer
    cmd_show
}

main() {
    local cmd=${1:-}
    case "${cmd}" in
        all) cmd_all ;;
        init) cmd_init ;;
        rollapp) cmd_rollapp ;;
        fund) cmd_fund ;;
        sequencer) cmd_sequencer ;;
        show) cmd_show ;;
        *) usage; exit 1 ;;
    esac
}

main "${1:-}"
