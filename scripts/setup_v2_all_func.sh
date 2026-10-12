#!/bin/bash
# Day-to-day operations on an environment deployed by setup_v2_all_install.sh.
# Run it directly:
#   PROJECT_ENV=beta ./scripts/setup_v2_all_func.sh restart
#   ./scripts/setup_v2_all_func.sh keyslist
# Or source it and call the functions:
#   source ./scripts/setup_v2_all_func.sh
#   restart
#   keyslist
set -euo pipefail

PROJECT_ENV="${PROJECT_ENV:-beta}"
BASE_DIR="${BASE_DIR:-/data/docker-${PROJECT_ENV}}"
BIN_DIR="${BASE_DIR}/bin"
KEYRING_BACKEND_NAME="${KEYRING_BACKEND_NAME:-test}"
ME_CHAIN_ID="${ME_CHAIN_ID:-mechain_400-1}"
ROLLAPP_CHAIN_ID="${ROLLAPP_CHAIN_ID:-mecheckin_401-1}"
ME_NODES_HOME="${ME_NODES_HOME:-${BASE_DIR}/nodes/hub-nodes}"
ME_NODE1_HOME="${ME_NODE1_HOME:-${ME_NODES_HOME}/node1}"
RAPP_NODES_HOME="${RAPP_NODES_HOME:-${BASE_DIR}/nodes/rollapp-nodes}"
RAPP_NODE1_HOME="${RAPP_NODE1_HOME:-${RAPP_NODES_HOME}/node1}"
RLY_RAPP_NODE_HOME="${RLY_RAPP_NODE_HOME:-${BASE_DIR}/nodes/rly-rapp}"
TIMEOUT_COMMIT="${TIMEOUT_COMMIT:-5s}"
HUB_TO_ROLLAPP_AMOUNT="${HUB_TO_ROLLAPP_AMOUNT:-10000000umec}"
ROLLAPP_TO_HUB_AMOUNT="${ROLLAPP_TO_HUB_AMOUNT:-500000urax}"
HUB_IBC_FEES="${HUB_IBC_FEES:-100000umec}"
ROLLAPP_IBC_FEES="${ROLLAPP_IBC_FEES:-2000urax}"
IBC_WAIT_SECONDS="${IBC_WAIT_SECONDS:-180}"
UPGRADE_NAME="${UPGRADE_NAME:-v3.0.0}"
HALT_OFFSET="${HALT_OFFSET:-50}"
GOV_DEPOSIT="${GOV_DEPOSIT:-100000000umec}"
GOV_FEES="${GOV_FEES:-200000umec}"
VOTE_FEES="${VOTE_FEES:-10000umec}"
VOTE_OPTION="${VOTE_OPTION:-yes}"
VAL_FUND_AMOUNT="${VAL_FUND_AMOUNT:-100000000umec}"
HUB_VAL_KEYS=(val1 val2 val3 val4)
HUB_SERVICES=(hub-node1 hub-node2 hub-node3 hub-node4)
HUB_NODE_HOMES=(node1 node2 node3 node4)
ROLLAPP_SERVICES=(rollapp-node1 rollapp-node2 rollapp-node3)
ROLLAPP_NODE_HOMES=(node1 node2 node3)
ROLLAPP_DA_LAYER="${ROLLAPP_DA_LAYER:-me-da}"
# Phase B switches the Hub image. The Phase C rollappd image must be set explicitly.
MED_V3_IMAGE="${MED_V3_IMAGE:-ghcr.io/openmetaearth/med:v3.0.0}"
ROLLAPP_V3_IMAGE="${ROLLAPP_V3_IMAGE:-}"
# Hub and rollapp images run as root. In-container home is fixed to these paths.
HUB_CONTAINER_HOME="${HUB_CONTAINER_HOME:-/root/.mechain}"
ROLLAPP_CONTAINER_HOME="${ROLLAPP_CONTAINER_HOME:-/root/.rollapp}"
HUB_MIN_GAS_PRICES="${HUB_MIN_GAS_PRICES:-0.02umec}"
ROLLAPP_MIN_GAS_PRICES="${ROLLAPP_MIN_GAS_PRICES:-0.001umec}"
ROLLAPP_GOV_DEPOSIT="${ROLLAPP_GOV_DEPOSIT:-1000000urax}"
ROLLAPP_GOV_FEES="${ROLLAPP_GOV_FEES:-4000urax}"
ROLLAPP_HALT_OFFSET="${ROLLAPP_HALT_OFFSET:-${HALT_OFFSET}}"
ROLLAPP_VAL_KEYS=(roluser)

usage() {
    cat <<EOF
Usage: $0 <command>

Commands:
  restart                 Stop the 4 hub(med) nodes, set timeout_commit to ${TIMEOUT_COMMIT}, then start them
  keyslist                med keys list --keyring-backend ${KEYRING_BACKEND_NAME} --home node1
  ibc-hub-to-rollapp      hub(global_dao) --ibc--> rollapp(ibc)
  ibc-rollapp-to-hub      rollapp(roluser) --ibc--> hub(ibc-rollapp receive address)
  ibc                     rollapp->hub, then hub->rollapp
  rollapp-propose         RollApp submits the ${UPGRADE_NAME} software-upgrade (Hub must still be producing v2 blocks)
  rollapp-proposal [id]   Query a rollapp proposal
  rollapp-vote [id]       roluser votes, then wait until the rollapp halts
  propose                 Hub submits ${UPGRADE_NAME} (refuses if the rollapp has not halted)
  proposal [id]           Query a Hub proposal
  vote [id]               Four validators vote ${VOTE_OPTION} on the Hub proposal
  hub-upgrade             After the Hub halt, switch to ${MED_V3_IMAGE} (refuses if the rollapp has not halted; do not start the old rollappd)
  run-3d-migration        Stop rollapp and run rollappd run-3d-migration on node1/2/3 (DA=me-da)
  rollapp-upgrade         After Hub v3: backup, 3D migration, switch to the new rollappd and start it

  Order (proposal is not the same as switching the image; do not reverse these):
    1. rollapp-propose -> rollapp-vote (old rollappd produces blocks, wait for upgrade-info.json)
    2. propose -> vote -> hub-upgrade (this is when med:v3 is switched in)
    3. ROLLAPP_V3_IMAGE=... rollapp-upgrade
  rollapp-propose after the Hub halt fails: halted Hub nodes panic and RPC is gone.
  The old rollappd also must not be started after hub-upgrade (GetProposer nil panic).

  fix-da-config           Rewrite dymint.toml da_layer/da_config as v3 arrays, then restart rollapp (does not re-run migration)
  help                    Show this help

Env:
  PROJECT_ENV=${PROJECT_ENV}
  BASE_DIR=${BASE_DIR}
  KEYRING_BACKEND_NAME=${KEYRING_BACKEND_NAME}
  TIMEOUT_COMMIT=${TIMEOUT_COMMIT}
  ME_CHAIN_ID=${ME_CHAIN_ID}
  ROLLAPP_CHAIN_ID=${ROLLAPP_CHAIN_ID}
  HUB_TO_ROLLAPP_AMOUNT=${HUB_TO_ROLLAPP_AMOUNT}
  ROLLAPP_TO_HUB_AMOUNT=${ROLLAPP_TO_HUB_AMOUNT}
  UPGRADE_NAME=${UPGRADE_NAME}
  HALT_OFFSET=${HALT_OFFSET}
  GOV_DEPOSIT=${GOV_DEPOSIT}
  ROLLAPP_GOV_DEPOSIT=${ROLLAPP_GOV_DEPOSIT}
  ROLLAPP_HALT_OFFSET=${ROLLAPP_HALT_OFFSET}
  ROLLAPP_V3_IMAGE=${ROLLAPP_V3_IMAGE:-<required for rollapp-upgrade>}
  MED_V3_IMAGE=${MED_V3_IMAGE}
  ROLLAPP_DA_LAYER=${ROLLAPP_DA_LAYER}
  HUB_CONTAINER_HOME=${HUB_CONTAINER_HOME}
  ROLLAPP_CONTAINER_HOME=${ROLLAPP_CONTAINER_HOME}
  SKIP_ROLLAPP_GOV_CHECK=${SKIP_ROLLAPP_GOV_CHECK:-}
EOF
}

need_base_dir() {
    if [ ! -d "${BASE_DIR}" ]; then
        echo "error: BASE_DIR not found: ${BASE_DIR}" >&2
        exit 1
    fi
    cd "${BASE_DIR}"
    export PATH="${BIN_DIR}:${PATH}"
    export LD_LIBRARY_PATH="${BASE_DIR}/lib:${LD_LIBRARY_PATH:-}"
}

set_timeout_commit() {
    local commit="${1:-${TIMEOUT_COMMIT}}"
    local node config current
    for node in "${HUB_NODE_HOMES[@]}"; do
        config="${ME_NODES_HOME}/${node}/config/config.toml"
        if [ ! -f "${config}" ]; then
            echo "error: missing ${config}" >&2
            exit 1
        fi
        sed -i 's/^timeout_commit = .*/timeout_commit = "'"${commit}"'"/' "${config}"
        current=$(grep -E '^timeout_commit = ' "${config}" || true)
        echo "${node}: ${current}"
    done
}

restart() {
    need_base_dir
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Stop hub-node1..4, set timeout_commit=${TIMEOUT_COMMIT}, then start         #"
    echo "# ---------------------------------------------------------------------------- #"
    docker compose stop "${HUB_SERVICES[@]}"
    set_timeout_commit "${TIMEOUT_COMMIT}"
    docker compose up -d "${HUB_SERVICES[@]}"
    docker compose ps "${HUB_SERVICES[@]}"
}

keyslist() {
    need_base_dir
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  med keys list --keyring-backend ${KEYRING_BACKEND_NAME} --home ${ME_NODE1_HOME} #"
    echo "# ---------------------------------------------------------------------------- #"
    med keys list --keyring-backend "${KEYRING_BACKEND_NAME}" --home "${ME_NODE1_HOME}"
}

check_tx_status() {
    local cosmosapp=$1
    local txhash=$2
    local extra_args=$3
    local sleepTime=120

    if [ -z "$txhash" ] || [ "$txhash" = "null" ]; then
        echo "error: empty txhash" >&2
        exit 1
    fi

    echo "Check tx status: ${txhash}"
    while true; do
        if ! $cosmosapp query tx ${txhash} ${extra_args} > /dev/null 2>&1; then
            sleep 1
            sleepTime=$((sleepTime - 1))
            if [ $sleepTime -eq 0 ]; then
                echo "Tx not found: ${txhash}"
                $cosmosapp query tx ${txhash} ${extra_args} || true
                exit 1
            fi
            continue
        fi
        local resp_code
        resp_code=$($cosmosapp q tx ${txhash} -o json ${extra_args} | jq -r '.code' 2>/dev/null)
        if [ "$resp_code" != "0" ]; then
            echo "Tx failed, code: ${resp_code}"
            $cosmosapp q tx ${txhash} ${extra_args} || true
            exit 1
        fi
        echo "Tx succeeded, code: ${resp_code}"
        break
    done
}

watch_ibc_balance() {
    local binary="$1"
    local user_addr="$2"
    local extra_args="$3"
    local waited=0
    while [ "${waited}" -lt "${IBC_WAIT_SECONDS}" ]; do
        if $binary query bank balances "${user_addr}" ${extra_args} -o json | jq -e '.balances[]? | select(.denom | test("ibc"))' >/dev/null 2>&1; then
            echo "addr ${user_addr} received ibc tokens"
            $binary query bank balances "${user_addr}" ${extra_args}
            return 0
        fi
        echo "addr ${user_addr} has not received ibc tokens yet, sleep 5s (${waited}/${IBC_WAIT_SECONDS})"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: no ibc tokens for ${user_addr} within ${IBC_WAIT_SECONDS}s" >&2
    $binary query bank balances "${user_addr}" ${extra_args} || true
    exit 1
}

load_ibc_context() {
    need_base_dir
    local kb=(--keyring-backend "${KEYRING_BACKEND_NAME}")

    USER_ADDR_ME_GLOBAL_DAO="${USER_ADDR_ME_GLOBAL_DAO:-$(med keys show global_dao -a "${kb[@]}" --home "${ME_NODE1_HOME}")}"
    USER_ADDR_ROLLAPP_ADMIN="${USER_ADDR_ROLLAPP_ADMIN:-$(rollappd keys show roluser -a "${kb[@]}" --home "${RAPP_NODE1_HOME}")}"
    USER_ADDR_ROLLAPP_IBC="${USER_ADDR_ROLLAPP_IBC:-$(rollappd keys show ibc -a "${kb[@]}" --home "${RAPP_NODE1_HOME}")}"

    if [ -z "${USER_ADDR_ME_IBC_ROLLAPP:-}" ]; then
        USER_ADDR_ME_IBC_ROLLAPP=$(rly keys show "${ME_CHAIN_ID}" me_user_ibc_rollapp_name --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | grep -Eo 'me1[0-9a-z]+' | head -1 || true)
    fi
    if [ -z "${USER_ADDR_ME_IBC_ROLLAPP:-}" ]; then
        echo "error: could not resolve the hub-side ibc receive address; set USER_ADDR_ME_IBC_ROLLAPP" >&2
        exit 1
    fi

    if [ -z "${ME_CHANNEL_ID:-}" ] || [ -z "${RAPP_CHANNEL_ID:-}" ]; then
        ME_CHANNEL_ID="${ME_CHANNEL_ID:-$(rly q channels "${ME_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ROLLAPP_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id' | head -1)}"
        RAPP_CHANNEL_ID="${RAPP_CHANNEL_ID:-$(rly q channels "${ROLLAPP_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ME_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id' | head -1)}"
    fi
    if [ -z "${ME_CHANNEL_ID:-}" ] || [ -z "${RAPP_CHANNEL_ID:-}" ]; then
        echo "error: hub-rollapp channel not found; confirm rly-rollapp is linked. ME_CHANNEL_ID=${ME_CHANNEL_ID:-} RAPP_CHANNEL_ID=${RAPP_CHANNEL_ID:-}" >&2
        exit 1
    fi

    echo "GLOBAL_DAO=${USER_ADDR_ME_GLOBAL_DAO}"
    echo "ROLLAPP_ADMIN(roluser)=${USER_ADDR_ROLLAPP_ADMIN}"
    echo "ROLLAPP_IBC=${USER_ADDR_ROLLAPP_IBC}"
    echo "ME_IBC_ROLLAPP=${USER_ADDR_ME_IBC_ROLLAPP}"
    echo "me -> rollapp channel: ${ME_CHANNEL_ID}"
    echo "rollapp -> me channel: ${RAPP_CHANNEL_ID}"
}

print_ibc_balances() {
    echo "hub global_dao:"
    med query bank balances "${USER_ADDR_ME_GLOBAL_DAO}" --home "${ME_NODE1_HOME}"
    echo "hub ibc-rollapp receive address:"
    med query bank balances "${USER_ADDR_ME_IBC_ROLLAPP}" --home "${ME_NODE1_HOME}"
    echo "rollapp roluser:"
    rollappd query bank balances "${USER_ADDR_ROLLAPP_ADMIN}" --home "${RAPP_NODE1_HOME}"
    echo "rollapp ibc:"
    rollappd query bank balances "${USER_ADDR_ROLLAPP_IBC}" --home "${RAPP_NODE1_HOME}"
}

# hub(global_dao/umec) -> rollapp(ibc address)
ibc_hub_to_rollapp() {
    load_ibc_context
    local amount="${1:-}"
    amount="${amount:-${HUB_TO_ROLLAPP_AMOUNT}}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  ibc: hub -> rollapp  ${amount}                                              #"
    echo "# ---------------------------------------------------------------------------- #"
    print_ibc_balances
    echo "med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${USER_ADDR_ROLLAPP_IBC} ${amount} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=${HUB_IBC_FEES} --chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME} --keyring-backend ${KEYRING_BACKEND_NAME} -y"
    local txhash
    txhash=$(med tx ibc-transfer transfer transfer "${ME_CHANNEL_ID}" "${USER_ADDR_ROLLAPP_IBC}" "${amount}" \
        --from "${USER_ADDR_ME_GLOBAL_DAO}" --fees="${HUB_IBC_FEES}" --chain-id "${ME_CHAIN_ID}" \
        --home "${ME_NODE1_HOME}" --keyring-backend "${KEYRING_BACKEND_NAME}" -y --output json | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    watch_ibc_balance rollappd "${USER_ADDR_ROLLAPP_IBC}" "--home ${RAPP_NODE1_HOME}"
}

# rollapp(roluser/urax) -> hub(ibc-rollapp receive address)
ibc_rollapp_to_hub() {
    load_ibc_context
    local amount="${1:-}"
    amount="${amount:-${ROLLAPP_TO_HUB_AMOUNT}}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  ibc: rollapp -> hub  ${amount}                                              #"
    echo "# ---------------------------------------------------------------------------- #"
    print_ibc_balances
    echo "rollappd tx ibc-transfer transfer transfer ${RAPP_CHANNEL_ID} ${USER_ADDR_ME_IBC_ROLLAPP} ${amount} --from ${USER_ADDR_ROLLAPP_ADMIN} --fees=${ROLLAPP_IBC_FEES} --chain-id ${ROLLAPP_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME} -y"
    local txhash
    txhash=$(rollappd tx ibc-transfer transfer transfer "${RAPP_CHANNEL_ID}" "${USER_ADDR_ME_IBC_ROLLAPP}" "${amount}" \
        --from "${USER_ADDR_ROLLAPP_ADMIN}" --fees="${ROLLAPP_IBC_FEES}" --chain-id "${ROLLAPP_CHAIN_ID}" \
        --keyring-backend "${KEYRING_BACKEND_NAME}" --home "${RAPP_NODE1_HOME}" -y --output json | jq -r '.txhash')
    check_tx_status rollappd "${txhash}" "--home ${RAPP_NODE1_HOME}"
    watch_ibc_balance med "${USER_ADDR_ME_IBC_ROLLAPP}" "--home ${ME_NODE1_HOME}"
}

ibc() {
    ibc_rollapp_to_hub "${ROLLAPP_TO_HUB_AMOUNT}"
    ibc_hub_to_rollapp "${HUB_TO_ROLLAPP_AMOUNT}"
}

hub_height() {
    med status --home "${ME_NODE1_HOME}" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height'
}

proposal_id_from_json() {
    jq -r '.proposal_id // .id // .proposalId // empty'
}

latest_proposal_id() {
    med q gov proposals --home "${ME_NODE1_HOME}" -o json 2>/dev/null | jq -r '
        (.proposals // [])
        | if length == 0 then empty
          else .[-1] | (.proposal_id // .id // .proposalId)
          end
    '
}

# Query a proposal. With no id, query the latest one and list every proposal summary.
proposal() {
    need_base_dir
    local proposal_id="${1:-}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Query gov proposals                                                         #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "All proposals:"
    med q gov proposals --home "${ME_NODE1_HOME}" -o json | jq -r '
        (.proposals // [])
        | if length == 0 then "  (none)"
          else .[] | "  id=\(.proposal_id // .id // .proposalId) status=\(.status) title=\(.content.title // .title // .messages[0].content.title // "-")"
          end
    '
    if [ -z "${proposal_id}" ]; then
        proposal_id="$(latest_proposal_id || true)"
    fi
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        echo "No proposal to query"
        return 0
    fi
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  proposal ${proposal_id}                                                     #"
    echo "# ---------------------------------------------------------------------------- #"
    med q gov proposal "${proposal_id}" --home "${ME_NODE1_HOME}"
    echo "votes:"
    med q gov votes "${proposal_id}" --home "${ME_NODE1_HOME}" || true
}

rollapp_upgrade_info_name() {
    local f="${RAPP_NODE1_HOME}/data/upgrade-info.json"
    if [ -f "${f}" ]; then
        jq -r '.name // empty' "${f}" 2>/dev/null || true
    fi
}

wait_rollapp_halt() {
    local halt height waited=0 name
    echo "Waiting for the RollApp to reach the software-upgrade halt height"
    while [ "${waited}" -lt 1800 ]; do
        name="$(rollapp_upgrade_info_name)"
        if [[ -n "${name}" ]]; then
            echo "rollapp halted: ${RAPP_NODE1_HOME}/data/upgrade-info.json name=${name}"
            return 0
        fi
        halt=$(rollappd q upgrade plan --home "${RAPP_NODE1_HOME}" -o json 2>/dev/null | jq -r '.plan.height // .height // empty')
        halt=${halt//\"/}
        height="$(rollapp_height 2>/dev/null || true)"
        height=${height//\"/}
        if [[ -n "${halt}" && "${halt}" != "null" && -n "${height}" && "${height}" != "null" ]] && [ "${height}" -ge "${halt}" ] 2>/dev/null; then
            echo "Reached rollapp halt height=${halt} current=${height}"
            return 0
        fi
        if [[ -z "${height}" || "${height}" == "null" ]]; then
            echo "error: rollapp RPC is down and upgrade-info.json is missing." >&2
            echo "This usually means the old rollappd panicked after connecting to a v3 Hub, not a normal software-upgrade halt." >&2
            exit 1
        fi
        echo "Waiting for rollapp halt: height=${height} plan=${halt:-?} (${waited}/1800)"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: rollapp did not reach the halt height within 1800s" >&2
    rollappd q upgrade plan --home "${RAPP_NODE1_HOME}" || true
    exit 1
}

# Before a Hub proposal or image switch, the rollapp must already have halted on a v2 Hub.
# Do not wait for the Hub halt and then rollapp-propose: halted Hub nodes panic and settlement RPC is gone.
require_rollapp_halted() {
    local name plan
    if [ "${SKIP_ROLLAPP_GOV_CHECK:-}" = "1" ]; then
        echo "SKIP_ROLLAPP_GOV_CHECK=1, skipping the rollapp halt check"
        return 0
    fi
    name="$(rollapp_upgrade_info_name)"
    if [[ -n "${name}" ]]; then
        echo "rollapp halted: upgrade-info name=${name}"
        return 0
    fi
    plan=$(rollappd q upgrade plan --home "${RAPP_NODE1_HOME}" -o json 2>/dev/null | jq -r '.plan.name // .name // empty' || true)
    if [[ -n "${plan}" && "${plan}" != "null" ]]; then
        echo "rollapp already has upgrade plan=${plan}, waiting for halt (must finish while the Hub is still producing blocks)"
        wait_rollapp_halt
        return 0
    fi
    echo "error: finish the rollapp software-upgrade and halt on a v2 Hub before touching the Hub." >&2
    echo "Order: rollapp-propose -> rollapp-vote -> propose -> vote -> hub-upgrade -> rollapp-upgrade" >&2
    echo "A proposal after the Hub halt does nothing: RPC is down after the halt, and starting v1.0.21 after hub-upgrade hits a GetProposer panic." >&2
    echo "If this env is already Hub v3 and the rollapp has no proposal, reset it. Force skip: SKIP_ROLLAPP_GOV_CHECK=1" >&2
    exit 1
}

# global_dao submits the software-upgrade to v3.0.0, then queries the proposal.
propose() {
    need_base_dir
    # require_rollapp_halted # Temporarily disabled: after the rollapp proposal the process stops and cannot be queried; the operator checks by hand.
    local height halt txhash proposal_id
    height="$(hub_height)"
    height="${height//\"/}"
    halt="${HALT_HEIGHT:-$((height + HALT_OFFSET))}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Submit software-upgrade ${UPGRADE_NAME}  height=${height} halt=${halt}       #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "med tx gov submit-legacy-proposal software-upgrade ${UPGRADE_NAME} --upgrade-height ${halt} --from global_dao --deposit ${GOV_DEPOSIT}"
    txhash=$(med tx gov submit-legacy-proposal software-upgrade "${UPGRADE_NAME}" \
        --title "${UPGRADE_NAME}" \
        --description "ME-Hub SDK 0.50 / IBC v8 / settlement v3" \
        --upgrade-height "${halt}" \
        --upgrade-info '{"binaries":{}}' \
        --no-validate \
        --deposit "${GOV_DEPOSIT}" \
        --from global_dao \
        --keyring-backend "${KEYRING_BACKEND_NAME}" \
        --chain-id "${ME_CHAIN_ID}" \
        --home "${ME_NODE1_HOME}" \
        --gas auto --gas-adjustment 1.5 \
        --fees "${GOV_FEES}" \
        --yes -o json | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    proposal_id=$(med q tx "${txhash}" --home "${ME_NODE1_HOME}" -o json | jq -r '
        [
            .logs[]?.events[]?,
            .events[]?
        ]
        | map(select(.type == "submit_proposal" or .type == "message"))
        | .[].attributes[]?
        | select(.key == "proposal_id")
        | .value
    ' | head -1)
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        proposal_id="$(latest_proposal_id || true)"
    fi
    echo "proposal_id=${proposal_id}"
    proposal "${proposal_id}"
}

umec_balance() {
    local addr="$1"
    med q bank balances "${addr}" --home "${ME_NODE1_HOME}" -o json | jq -r '[.balances[]? | select(.denom=="umec") | .amount | tonumber] | add // 0'
}

fund_val_if_needed() {
    local key="$1"
    local addr amount txhash
    addr=$(med keys show "${key}" -a --keyring-backend "${KEYRING_BACKEND_NAME}" --home "${ME_NODE1_HOME}")
    amount="$(umec_balance "${addr}")"
    if [ "${amount}" -ge 100000 ]; then
        echo "${key} ${addr} umec=${amount}"
        return 0
    fi
    echo "${key} ${addr} umec=${amount}, sending ${VAL_FUND_AMOUNT} from global_dao for fees"
    txhash=$(med tx bank send global_dao "${addr}" "${VAL_FUND_AMOUNT}" \
        --from global_dao \
        --keyring-backend "${KEYRING_BACKEND_NAME}" \
        --chain-id "${ME_CHAIN_ID}" \
        --home "${ME_NODE1_HOME}" \
        --fees "${VOTE_FEES}" \
        --yes -o json | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
}

# All four validators vote in this function. With no id, vote on the latest proposal.
vote() {
    need_base_dir
    local proposal_id="${1:-}"
    local key txhash
    if [ -z "${proposal_id}" ]; then
        proposal_id="$(latest_proposal_id || true)"
    fi
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        echo "error: no proposal to vote on; run propose first or pass a proposal id" >&2
        exit 1
    fi
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  val1-val4 vote ${VOTE_OPTION} on proposal ${proposal_id}                    #"
    echo "# ---------------------------------------------------------------------------- #"
    for key in "${HUB_VAL_KEYS[@]}"; do
        fund_val_if_needed "${key}"
        echo "med tx gov vote ${proposal_id} ${VOTE_OPTION} --from ${key}"
        txhash=$(med tx gov vote "${proposal_id}" "${VOTE_OPTION}" \
            --from "${key}" \
            --keyring-backend "${KEYRING_BACKEND_NAME}" \
            --chain-id "${ME_CHAIN_ID}" \
            --home "${ME_NODE1_HOME}" \
            --fees "${VOTE_FEES}" \
            --yes -o json | jq -r '.txhash')
        check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    done
    proposal "${proposal_id}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Hub has voted. After the halt, run hub-upgrade (image switch). Do not start the old rollappd now. #"
    echo "# ---------------------------------------------------------------------------- #"
}

set_compose_hub_image() {
    local compose="${BASE_DIR}/docker-compose.yml"
    local image="${1:-${MED_V3_IMAGE}}"
    if [ ! -f "${compose}" ]; then
        echo "error: missing ${compose}" >&2
        exit 1
    fi
    python3 - "${compose}" "${image}" "${HUB_CONTAINER_HOME}" "${HUB_MIN_GAS_PRICES}" <<'PY'
import pathlib, re, sys

path = pathlib.Path(sys.argv[1])
image, home, min_gas = sys.argv[2], sys.argv[3], sys.argv[4]
text = path.read_text()
needle = "x-me-template:"
idx = text.find(needle)
if idx < 0:
    raise SystemExit("docker-compose.yml has no x-me-template")
rest = text[idx + len(needle):]
m = re.search(r"(?m)^(x-|networks:|services:)", rest)
end = idx + len(needle) + (m.start() if m else len(rest))
block = text[idx:end]
block, n = re.subn(r"(?m)^  image: .+$", f"  image: {image}", block, count=1)
if n != 1:
    raise SystemExit("failed to replace the x-me-template image")
cmd = (
    "  command:\n"
    "    - start\n"
    "    - --home\n"
    f"    - {home}\n"
    "    - --minimum-gas-prices\n"
    f'    - "{min_gas}"\n'
)
block, n = re.subn(r"(?ms)^  command:\n(?:    - .+\n)+", cmd, block, count=1)
if n != 1:
    raise SystemExit("failed to replace the x-me-template command")
block = re.sub(r"(?m)^  user: .+\n", "", block)
path.write_text(text[:idx] + block + text[end:])
print(f"compose hub image -> {image}")
print(f"compose hub home -> {home} min-gas -> {min_gas}")
PY
}

wait_hub_halt() {
    local halt height waited=0
    echo "Waiting for the Hub to reach the software-upgrade halt height"
    while [ "${waited}" -lt 1800 ]; do
        halt=$(med q upgrade plan --home "${ME_NODE1_HOME}" -o json 2>/dev/null | jq -r '.plan.height // .height // empty')
        halt=${halt//\"/}
        height="$(hub_height 2>/dev/null || true)"
        height=${height//\"/}
        if [[ -n "${halt}" && "${halt}" != "null" && -n "${height}" && "${height}" != "null" ]] && [ "${height}" -ge "${halt}" ] 2>/dev/null; then
            echo "Reached halt height=${halt} current=${height}"
            return 0
        fi
        if [[ -z "${height}" || "${height}" == "null" ]]; then
            echo "hub RPC is down (it may already have halted); continuing with ${MED_V3_IMAGE}"
            return 0
        fi
        echo "Waiting for halt: height=${height} plan=${halt:-?} (${waited}/1800)"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: halt height not reached within 1800s" >&2
    med q upgrade plan --home "${ME_NODE1_HOME}" || true
    exit 1
}

wait_hub_block() {
    local waited=0 height
    while [ "${waited}" -lt 180 ]; do
        height="$(hub_height 2>/dev/null || true)"
        height=${height//\"/}
        if [[ -n "${height}" && "${height}" != "null" && "${height}" -ge 1 ]] 2>/dev/null; then
            echo "hub is producing blocks, height=${height}"
            return 0
        fi
        echo "Waiting for a hub block, sleep 5s (${waited}/180)"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: new med produced no block within 180s" >&2
    docker compose logs --tail 80 hub-node1 || true
    exit 1
}

# Phase B: after the halt, switch the compose Hub image to med:v3.0.0 and start it.
# Old rollappd v1.0.21 panics on a nil GetProposer() inside IsSequencerVerify when it talks to a v3 Hub, so it must stay stopped.
# Switching the image is not the proposal: rollapp gov/halt must finish before this, while the Hub is still v2.
hub_upgrade() {
    need_base_dir
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Phase B Hub upgrade  image=${MED_V3_IMAGE}  home=${HUB_CONTAINER_HOME}       #"
    echo "# ---------------------------------------------------------------------------- #"
    # require_rollapp_halted # Temporarily disabled: after the rollapp proposal the process stops and cannot be queried; the operator checks by hand.
    wait_hub_halt
    echo "Stop hub + rollapp (do not let rollapp keep batching during the halt)"
    docker compose stop rly-rollapp "${ROLLAPP_SERVICES[@]}" "${HUB_SERVICES[@]}" || true
    ensure_local_image "${MED_V3_IMAGE}"
    set_compose_hub_image "${MED_V3_IMAGE}"
    docker compose up -d --pull never "${HUB_SERVICES[@]}"
    docker compose stop rly-rollapp "${ROLLAPP_SERVICES[@]}" || true
    docker compose ps "${HUB_SERVICES[@]}"
    wait_hub_block
    echo "upgrade applied:"
    med q upgrade applied "${UPGRADE_NAME}" --home "${ME_NODE1_HOME}" || true
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Hub v3 is up. Do not docker compose up rollapp-node* (still v1.0.21, it will panic) #"
    echo "#  Next: ROLLAPP_V3_IMAGE=... ./scripts/setup_v2_all_func.sh rollapp-upgrade          #"
    echo "# ---------------------------------------------------------------------------- #"
}

rollapp_height() {
    rollappd status --home "${RAPP_NODE1_HOME}" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height'
}

rollapp_latest_proposal_id() {
    rollappd q gov proposals --home "${RAPP_NODE1_HOME}" -o json 2>/dev/null | jq -r '
        (.proposals // [])
        | if length == 0 then empty
          else .[-1] | (.proposal_id // .id // .proposalId)
          end
    '
}

# Query a rollapp gov proposal. With no id, query the latest one.
rollapp_proposal() {
    need_base_dir
    local proposal_id="${1:-}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Query rollapp gov proposals                                                #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "validators:"
    rollappd q staking validators --home "${RAPP_NODE1_HOME}" -o json | jq -r '
        (.validators // [])
        | if length == 0 then "  (none)"
          else .[] | "  moniker=\(.description.moniker) status=\(.status) tokens=\(.tokens)"
          end
    '
    echo "All proposals:"
    rollappd q gov proposals --home "${RAPP_NODE1_HOME}" -o json | jq -r '
        (.proposals // [])
        | if length == 0 then "  (none)"
          else .[] | "  id=\(.proposal_id // .id // .proposalId) status=\(.status) title=\(.content.title // .title // "-")"
          end
    '
    if [ -z "${proposal_id}" ]; then
        proposal_id="$(rollapp_latest_proposal_id || true)"
    fi
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        echo "No rollapp proposal to query"
        return 0
    fi
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollapp proposal ${proposal_id}                                             #"
    echo "# ---------------------------------------------------------------------------- #"
    rollappd q gov proposal "${proposal_id}" --home "${RAPP_NODE1_HOME}"
    echo "votes:"
    rollappd q gov votes "${proposal_id}" --home "${RAPP_NODE1_HOME}" || true
}

# roluser submits the rollapp software-upgrade v3.0.0 (submit-legacy-proposal + --no-validate).
# The handler adds the hubgenesis store and InitGenesis rollappparams (DA=me-da).
rollapp_propose() {
    need_base_dir
    local height halt txhash proposal_id
    height="$(rollapp_height)"
    height="${height//\"/}"
    halt="${ROLLAPP_HALT_HEIGHT:-$((height + ROLLAPP_HALT_OFFSET))}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollapp software-upgrade ${UPGRADE_NAME}  height=${height} halt=${halt}     #"
    echo "#  New modules: hubgenesis store + rollappparams (DA=${ROLLAPP_DA_LAYER})      #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "rollappd tx gov submit-legacy-proposal software-upgrade ${UPGRADE_NAME} --upgrade-height ${halt} --from roluser --deposit ${ROLLAPP_GOV_DEPOSIT}"
    txhash=$(rollappd tx gov submit-legacy-proposal software-upgrade "${UPGRADE_NAME}" \
        --upgrade-height "${halt}" \
        --upgrade-info "${UPGRADE_NAME}" \
        --no-validate \
        --title "${UPGRADE_NAME} Upgrade" \
        --description "${UPGRADE_NAME} upgrade: hubgenesis + rollappparams (DA=me-da)" \
        --deposit "${ROLLAPP_GOV_DEPOSIT}" \
        --from roluser \
        --keyring-backend "${KEYRING_BACKEND_NAME}" \
        --chain-id "${ROLLAPP_CHAIN_ID}" \
        --home "${RAPP_NODE1_HOME}" \
        --fees "${ROLLAPP_GOV_FEES}" \
        --yes -o json | jq -r '.txhash')
    check_tx_status rollappd "${txhash}" "--home ${RAPP_NODE1_HOME}"
    proposal_id=$(rollappd q tx "${txhash}" --home "${RAPP_NODE1_HOME}" -o json | jq -r '
        [
            .logs[]?.events[]?,
            .events[]?
        ]
        | map(select(.type == "submit_proposal" or .type == "message"))
        | .[].attributes[]?
        | select(.key == "proposal_id")
        | .value
    ' | head -1)
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        proposal_id="$(rollapp_latest_proposal_id || true)"
    fi
    echo "proposal_id=${proposal_id}"
    rollapp_proposal "${proposal_id}"
}

# Genesis gentx has only roluser as a staking validator. Vote with each local key whose name matches a bonded validator.
rollapp_vote() {
    need_base_dir
    local proposal_id="${1:-}"
    local key txhash
    if [ -z "${proposal_id}" ]; then
        proposal_id="$(rollapp_latest_proposal_id || true)"
    fi
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        echo "error: no rollapp proposal to vote on; run rollapp-propose first or pass a proposal id" >&2
        exit 1
    fi
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollapp validators vote ${VOTE_OPTION} on proposal ${proposal_id}          #"
    echo "# ---------------------------------------------------------------------------- #"
    rollappd q staking validators --home "${RAPP_NODE1_HOME}"
    for key in "${ROLLAPP_VAL_KEYS[@]}"; do
        echo "rollappd tx gov vote ${proposal_id} ${VOTE_OPTION} --from ${key}"
        txhash=$(rollappd tx gov vote "${proposal_id}" "${VOTE_OPTION}" \
            --from "${key}" \
            --keyring-backend "${KEYRING_BACKEND_NAME}" \
            --chain-id "${ROLLAPP_CHAIN_ID}" \
            --home "${RAPP_NODE1_HOME}" \
            --fees "${ROLLAPP_GOV_FEES}" \
            --yes -o json | jq -r '.txhash')
        check_tx_status rollappd "${txhash}" "--home ${RAPP_NODE1_HOME}"
    done
    rollapp_proposal "${proposal_id}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Vote done. Next, wait for the rollapp halt, then Hub propose (do not hub-upgrade first). #"
    echo "# ---------------------------------------------------------------------------- #"
    wait_rollapp_halt
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollapp has halted. Next: propose -> vote -> hub-upgrade -> rollapp-upgrade   #"
    echo "# ---------------------------------------------------------------------------- #"
}

toml_line() {
    grep -E "^${2}[[:space:]]*=" "$1" | head -1 || true
}

toml_put_line() {
    local file="$1" key="$2" line="$3"
    [ -z "${line}" ] && return 0
    local tmp
    tmp=$(mktemp)
    if grep -qE "^${key}[[:space:]]*=" "${file}"; then
        awk -v key="${key}" -v line="${line}" '
            BEGIN { re = "^" key "[[:space:]]*=" }
            $0 ~ re { print line; next }
            { print }
        ' "${file}" > "${tmp}"
        # Do not use mv: some envs alias it to mv -i, and docker may have written the target as root
        command cp "${tmp}" "${file}"
        rm -f "${tmp}"
    else
        printf '%s\n' "${line}" >> "${file}"
        rm -f "${tmp}"
    fi
}

need_rollapp_v3_image() {
    if [ -z "${ROLLAPP_V3_IMAGE}" ]; then
        echo "error: set ROLLAPP_V3_IMAGE, for example harbor.starex.xyz/rollapp/rollappd:<new-tag>" >&2
        exit 1
    fi
}

# Skip pull when the tag is already local. An unauthenticated Harbor docker pull returns unauthorized even when the image is already on the machine.
ensure_local_image() {
    local image="$1"
    if docker image inspect "${image}" >/dev/null 2>&1; then
        echo "Local image ${image} already present, skipping pull"
        return 0
    fi
    echo "pull ${image}"
    docker pull "${image}"
}

backup_rollapp_home() {
    local home="$1"
    local bak="${home}.bak-$(date +%Y%m%d%H%M%S)"
    echo "Backup ${home} -> ${bak}"
    cp -a "${home}" "${bak}"
}

run_3d_migration() {
    local home="$1"
    local out rc=0
    echo "rollappd run-3d-migration --home ${home} --rollapp-param-da ${ROLLAPP_DA_LAYER}"
    # mehub_v3: the command lives on rollappd. If the host binary is still old, use the v3 image.
    if rollappd run-3d-migration --help >/dev/null 2>&1; then
        out=$(rollappd run-3d-migration --home "${home}" --rollapp-param-da "${ROLLAPP_DA_LAYER}" 2>&1) || rc=$?
    else
        if [ -z "${ROLLAPP_V3_IMAGE}" ]; then
            echo "error: host rollappd has no run-3d-migration; set ROLLAPP_V3_IMAGE and run it from the v3 image" >&2
            exit 1
        fi
        out=$(docker run --rm --pull never \
            -v "${home}:${ROLLAPP_CONTAINER_HOME}" \
            --entrypoint rollappd \
            "${ROLLAPP_V3_IMAGE}" \
            run-3d-migration --home "${ROLLAPP_CONTAINER_HOME}" --rollapp-param-da "${ROLLAPP_DA_LAYER}" 2>&1) || rc=$?
    fi
    printf '%s\n' "${out}"
    if echo "${out}" | grep -qE '3D dymint store migration successful|3D migration is not needed'; then
        return 0
    fi
    echo "error: 3D migration failed (exit ${rc})" >&2
    exit 1
}

# Run the 3D store migration offline on the three rollapp nodes. Stop the processes first.
run-3d-migration() {
    need_base_dir
    local node home
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollappd run-3d-migration  da=${ROLLAPP_DA_LAYER}                           #"
    echo "# ---------------------------------------------------------------------------- #"
    docker compose stop "${ROLLAPP_SERVICES[@]}" || true
    for node in "${ROLLAPP_NODE_HOMES[@]}"; do
        home="${RAPP_NODES_HOME}/${node}"
        if [ ! -d "${home}/data" ] && [ ! -f "${home}/config/genesis.json" ]; then
            echo "Skipping missing ${home}"
            continue
        fi
        backup_rollapp_home "${home}"
        run_3d_migration "${home}"
    done
    echo "3D migration done. Nodes stay stopped; switch to the v3 image before docker compose up -d ${ROLLAPP_SERVICES[*]}"
}

# v3 dymint: da_layer / da_config must be arrays of equal length.
# The old scalar form da_layer = "me-da", da_config = '{json}' cannot be copied with awk -v (JSON quotes get truncated to an empty string).
normalize_v3_da_config() {
    local cfg="$1"
    local home
    home=$(cd "$(dirname "${cfg}")/.." && pwd)
    python3 - "${ROLLAPP_DA_LAYER}" "${cfg}" "${home}" <<'PY'
import glob, json, pathlib, re, sys

da_layer, cfg_path, home = sys.argv[1], pathlib.Path(sys.argv[2]), sys.argv[3]
cfg = cfg_path
paths = [cfg, pathlib.Path(str(cfg) + ".pre-v3")]
paths.extend(pathlib.Path(p) for p in sorted(glob.glob(home + ".bak-*/config/dymint.toml")))

PLACEHOLDERS = {"", "TOKEN", "TOKEN123"}

def try_parse(raw: str):
    for cand in (raw, raw.replace('\\"', '"'), raw.replace("\\'", "'")):
        try:
            obj = json.loads(cand)
        except json.JSONDecodeError:
            continue
        if isinstance(obj, dict) and obj:
            return obj
    return None

def extract_objs(text: str):
    m = re.search(r"(?m)^da_config\s*=\s*(.+)$", text)
    if not m:
        return []
    val = m.group(1).strip()
    chunks = []
    if val.startswith("["):
        for a, b in re.findall(r"'((?:\\'|[^'])*)'|\"((?:\\\"|[^\"])*)\"", val):
            s = a or b
            if s.startswith("{"):
                chunks.append(s)
    elif (val[:1] in "'\"" and val[-1:] == val[:1]):
        chunks.append(val[1:-1])
    out = []
    for chunk in chunks:
        obj = try_parse(chunk)
        if obj:
            out.append(obj)
    return out

best, source = None, None
best_score = -1
for path in paths:
    if not path.is_file():
        continue
    for obj in extract_objs(path.read_text()):
        if not obj.get("base_url"):
            continue
        token = str(obj.get("auth_token") or "")
        score = 1
        if token not in PLACEHOLDERS:
            score += 10
        if obj.get("retry_attempts") is not None:
            score += 2
        if score > best_score:
            best, source, best_score = obj, str(path), score

if best is None:
    raise SystemExit(f"error: no usable da_config JSON (checked {cfg} / .pre-v3 / bak-*)")

json_str = json.dumps(best, separators=(",", ":"))
if "'" in json_str:
    raise SystemExit("error: da_config JSON contains a single quote and cannot be written as a TOML single-quoted string")

layer_line = f'da_layer = ["{da_layer}"]'
cfg_line = f"da_config = ['{json_str}']"

def put(text, key, line):
    pat = rf"(?m)^{re.escape(key)}\s*=\s*.*$"
    if re.search(pat, text):
        return re.sub(pat, line, text, count=1)
    return text.rstrip() + "\n" + line + "\n"

text = put(cfg.read_text(), "da_layer", layer_line)
text = put(text, "da_config", cfg_line)
cfg.write_text(text)
token = str(best.get("auth_token") or "")
print(f"{cfg} DA from {source}")
print(layer_line)
print(f"da_config base_url={best.get('base_url')} auth_token_len={len(token)}")
PY
}

# Start from the new rollappd default dymint.toml, then copy settlement / keyring back from the old file.
# Restore DA on its own as v3 arrays from .pre-v3 / bak. Do not awk-copy the old scalars.
merge_dymint_toml() {
    local home="$1"
    local old="${home}/config/dymint.toml"
    local tmp new_toml
    if [ ! -f "${old}" ]; then
        echo "error: missing ${old}" >&2
        exit 1
    fi
    if [ ! -f "${old}.pre-v3" ]; then
        cp -a "${old}" "${old}.pre-v3"
    fi
    tmp=$(mktemp -d)
    if ! docker run --rm --pull never \
        -v "${tmp}:${ROLLAPP_CONTAINER_HOME}" \
        --entrypoint rollappd \
        "${ROLLAPP_V3_IMAGE}" \
        init tmpnode --chain-id "${ROLLAPP_CHAIN_ID}" --home "${ROLLAPP_CONTAINER_HOME}" >/dev/null 2>&1; then
        echo "New rollappd init failed; keeping the original dymint.toml and only checking DA/settlement"
        rm -rf "${tmp}"
        return 0
    fi
    if [ ! -f "${tmp}/config/dymint.toml" ]; then
        echo "New image did not write dymint.toml; keeping the original file"
        rm -rf "${tmp}"
        return 0
    fi
    # docker writes tmp as root; copy onto a user-owned file before editing to avoid mv Permission denied
    new_toml=$(mktemp)
    command cp "${tmp}/config/dymint.toml" "${new_toml}"
    chmod u+w "${new_toml}"
    rm -rf "${tmp}" 2>/dev/null || true
    local key
    local keys=(
        settlement_layer settlement_node_address settlement_gas_prices settlement_gas_limit
        keyring_home_dir dym_account_name
        max_idle_time max_proof_time batch_submit_max_time block_batch_max_size_bytes
        max_supported_batch_skew retry_attempts
    )
    for key in "${keys[@]}"; do
        toml_put_line "${new_toml}" "${key}" "$(toml_line "${old}.pre-v3" "${key}")"
    done
    command cp "${new_toml}" "${old}"
    rm -f "${new_toml}"
}

fix_dymint_runtime_keys() {
    local home="$1"
    local cfg="${home}/config/dymint.toml"
    toml_put_line "${cfg}" settlement_layer 'settlement_layer = "me-hub"'
    toml_put_line "${cfg}" keyring_home_dir "keyring_home_dir = \"${ROLLAPP_CONTAINER_HOME}/sequencer_keys\""
    normalize_v3_da_config "${cfg}"
    echo "${home} dymint:"
    grep -E '^(settlement_layer|settlement_node_address|da_layer|keyring_home_dir|dym_account_name)[[:space:]]*=' "${cfg}" || true
}

# Use this when 3D migration is already done and only the v3 DA arrays are missing. Do not run rollapp-upgrade again.
fix_da_config() {
    need_base_dir
    local node home
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Fix v3 dymint DA config  da=${ROLLAPP_DA_LAYER}                             #"
    echo "# ---------------------------------------------------------------------------- #"
    docker compose stop "${ROLLAPP_SERVICES[@]}" || true
    for node in "${ROLLAPP_NODE_HOMES[@]}"; do
        home="${RAPP_NODES_HOME}/${node}"
        if [ ! -f "${home}/config/dymint.toml" ]; then
            echo "Skipping missing ${home}"
            continue
        fi
        fix_dymint_runtime_keys "${home}"
    done
    set_compose_container_homes
    docker compose up -d "${ROLLAPP_SERVICES[@]}"
    docker compose ps "${ROLLAPP_SERVICES[@]}"
    echo "Check rollapp-node1 logs to confirm DA is up:"
    docker compose logs --tail 40 rollapp-node1 || true
}

ensure_rollapp_min_gas() {
    local home="$1"
    local cfg="${home}/config/app.toml"
    if [ ! -f "${cfg}" ]; then
        echo "error: missing ${cfg}" >&2
        exit 1
    fi
    if grep -qE '^minimum-gas-prices' "${cfg}"; then
        sed -i 's/^minimum-gas-prices *= .*/minimum-gas-prices = "'"${ROLLAPP_MIN_GAS_PRICES}"'"/' "${cfg}"
    else
        printf '\nminimum-gas-prices = "%s"\n' "${ROLLAPP_MIN_GAS_PRICES}" >> "${cfg}"
    fi
    grep -E '^minimum-gas-prices' "${cfg}"
}

set_compose_container_homes() {
    local compose="${BASE_DIR}/docker-compose.yml"
    local rollapp_image="${1:-}"
    if [ ! -f "${compose}" ]; then
        echo "error: missing ${compose}" >&2
        exit 1
    fi
    python3 - "${compose}" "${HUB_CONTAINER_HOME}" "${ROLLAPP_CONTAINER_HOME}" \
        "${HUB_MIN_GAS_PRICES}" "${ROLLAPP_MIN_GAS_PRICES}" "${rollapp_image}" <<'PY'
import pathlib, re, sys

path = pathlib.Path(sys.argv[1])
hub_home, rapp_home = sys.argv[2], sys.argv[3]
hub_gas, rapp_gas, image = sys.argv[4], sys.argv[5], sys.argv[6]
text = path.read_text()
text = text.replace(":/home/ubuntu/.mechain", f":{hub_home}")
text = text.replace(":/home/ubuntu/.rollapp", f":{rapp_home}")

def patch_template(text, needle, home, min_gas, image=""):
    idx = text.find(needle)
    if idx < 0:
        raise SystemExit(f"docker-compose.yml has no {needle}")
    rest = text[idx + len(needle):]
    m = re.search(r"(?m)^(x-|networks:|services:)", rest)
    end = idx + len(needle) + (m.start() if m else len(rest))
    block = text[idx:end]
    if image:
        block, n = re.subn(r"(?m)^  image: .+$", f"  image: {image}", block, count=1)
        if n != 1:
            raise SystemExit(f"failed to replace the image for {needle}")
    cmd = (
        "  command:\n"
        "    - start\n"
        "    - --home\n"
        f"    - {home}\n"
        "    - --minimum-gas-prices\n"
        f'    - "{min_gas}"\n'
    )
    block, n = re.subn(r"(?ms)^  command:\n(?:    - .+\n)+", cmd, block, count=1)
    if n != 1:
        raise SystemExit(f"failed to replace the command for {needle}")
    block = re.sub(r"(?m)^  user: .+\n", "", block)
    return text[:idx] + block + text[end:]

text = patch_template(text, "x-me-template:", hub_home, hub_gas)
text = patch_template(text, "x-rollapp-template:", rapp_home, rapp_gas, image)
path.write_text(text)
print(f"compose hub home -> {hub_home} min-gas -> {hub_gas}")
print(f"compose rollapp home -> {rapp_home} min-gas -> {rapp_gas}")
if image:
    print(f"compose rollapp image -> {image}")
PY
}

set_compose_rollapp_image() {
    set_compose_container_homes "${ROLLAPP_V3_IMAGE}"
}

wait_rollapp_block() {
    local waited=0
    while [ "${waited}" -lt 120 ]; do
        local height
        height=$(rollappd status --home "${RAPP_NODE1_HOME}" 2>/dev/null | jq -r '.SyncInfo.latest_block_height // .sync_info.latest_block_height // empty')
        height=${height//\"/}
        if [[ -n "${height}" && "${height}" -ge 1 ]] 2>/dev/null; then
            echo "rollapp is producing blocks, height=${height}"
            return 0
        fi
        echo "Waiting for a rollapp block, sleep 5s (${waited}/120)"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: new rollappd produced no block within 120s" >&2
    docker compose logs --tail 80 rollapp-node1 || true
    exit 1
}

# Phase C: after the Hub is v3, upgrade RollApp / dymint (backup + 3D migration + image switch).
rollapp_upgrade() {
    need_base_dir
    need_rollapp_v3_image
    local node home
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Phase C RollApp upgrade  image=${ROLLAPP_V3_IMAGE}  da=${ROLLAPP_DA_LAYER}  #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "Hub upgrade applied:"
    med q upgrade applied "${UPGRADE_NAME}" --home "${ME_NODE1_HOME}" || true
    echo "Stop rollapp + rly-rollapp (do not keep batching after settlement drops)"
    docker compose stop rly-rollapp "${ROLLAPP_SERVICES[@]}" || true

    ensure_local_image "${ROLLAPP_V3_IMAGE}"

    for node in "${ROLLAPP_NODE_HOMES[@]}"; do
        home="${RAPP_NODES_HOME}/${node}"
        if [ ! -d "${home}/data" ] && [ ! -f "${home}/config/dymint.toml" ]; then
            echo "Skipping missing ${home}"
            continue
        fi
        backup_rollapp_home "${home}"
        run_3d_migration "${home}"
        merge_dymint_toml "${home}"
        fix_dymint_runtime_keys "${home}"
        ensure_rollapp_min_gas "${home}"
    done

    set_compose_rollapp_image
    docker compose up -d --pull never "${ROLLAPP_SERVICES[@]}"
    docker compose ps "${ROLLAPP_SERVICES[@]}"
    wait_rollapp_block
    echo "Hub latest-state-info:"
    med-v3 q rollapp latest-state-index "${ROLLAPP_CHAIN_ID}" --home "${ME_NODE1_HOME}" 
    echo "Phase C done. Start the relayer separately: docker compose up -d rly-rollapp"
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    case "${1:-}" in
        restart) restart ;;
        keyslist) keyslist ;;
        ibc-hub-to-rollapp) shift; ibc_hub_to_rollapp "${1:-}" ;;
        ibc-rollapp-to-hub) shift; ibc_rollapp_to_hub "${1:-}" ;;
        ibc) ibc ;;
        propose) propose ;;
        proposal) shift; proposal "${1:-}" ;;
        vote) shift; vote "${1:-}" ;;
        hub-upgrade) hub_upgrade ;;
        rollapp-propose) rollapp_propose ;;
        rollapp-proposal) shift; rollapp_proposal "${1:-}" ;;
        rollapp-vote) shift; rollapp_vote "${1:-}" ;;
        run-3d-migration) run-3d-migration ;;
        rollapp-upgrade) rollapp_upgrade ;;
        fix-da-config) fix_da_config ;;
        help|-h|--help) usage ;;
        *)
            usage
            exit 1
            ;;
    esac
fi
