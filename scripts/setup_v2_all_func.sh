#!/bin/bash
# 对 setup_v2_all_install.sh 部署好的环境做日常操作。
# 可直接执行：
#   PROJECT_ENV=beta ./scripts/setup_v2_all_func.sh restart
#   ./scripts/setup_v2_all_func.sh keyslist
# 也可 source 后调用函数：
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
HALT_OFFSET="${HALT_OFFSET:-100}"
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
# Phase B 换 Hub 镜像；Phase C 新 rollappd 镜像必须显式指定
MED_V3_IMAGE="${MED_V3_IMAGE:-ghcr.io/openmetaearth/med:v3.0.0}"
ROLLAPP_V3_IMAGE="${ROLLAPP_V3_IMAGE:-}"
# hub / rollapp 镜像都以 root 跑，容器内 home 固定这两个路径
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
  restart                 停掉 4 个 hub(med) 节点，把 timeout_commit 改成 ${TIMEOUT_COMMIT}，再启动
  keyslist                med keys list --keyring-backend ${KEYRING_BACKEND_NAME} --home node1
  ibc-hub-to-rollapp      hub(global_dao) --ibc--> rollapp(ibc)
  ibc-rollapp-to-hub      rollapp(roluser) --ibc--> hub(ibc-rollapp 收款地址)
  ibc                     先 rollapp->hub，再 hub->rollapp
  rollapp-propose         RollApp 提交 ${UPGRADE_NAME} software-upgrade（必须在 Hub 仍是 v2 出块时）
  rollapp-proposal [id]   查询 rollapp 提案
  rollapp-vote [id]       roluser 投票，并等到 rollapp halt
  propose                 Hub 提交 ${UPGRADE_NAME}（拒绝：若 rollapp 还没 halt）
  proposal [id]           查询 Hub 提案
  vote [id]               四个 validator 对 Hub 提案投 ${VOTE_OPTION}
  hub-upgrade             Hub halt 后换 ${MED_V3_IMAGE}（拒绝：若 rollapp 还没 halt；不要启动旧 rollappd）
  run-3d-migration        停 rollapp，对 node1/2/3 跑 rollappd run-3d-migration（DA=me-da）
  rollapp-upgrade         Hub v3 之后：备份、3D migration、换新 rollappd 启动

  顺序（提案 ≠ 换镜像，不能颠倒）：
    1. rollapp-propose → rollapp-vote（旧 rollappd 出块，等 upgrade-info.json）
    2. propose → vote → hub-upgrade（这时才换 med:v3）
    3. ROLLAPP_V3_IMAGE=... rollapp-upgrade
  Hub halt 后再 rollapp-propose 会失败：halt 后 Hub 节点 panic，RPC 没了。
  旧 rollappd 也不能在 hub-upgrade 之后启动（GetProposer nil panic）。

  fix-da-config           把 dymint.toml 的 da_layer/da_config 改成 v3 数组格式后重启 rollapp（不重跑 migration）
  help                    显示本说明

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
    echo "#  停掉 hub-node1..4，timeout_commit=${TIMEOUT_COMMIT}，再启动                   #"
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
            echo "addr ${user_addr} 已收到 ibc 代币"
            $binary query bank balances "${user_addr}" ${extra_args}
            return 0
        fi
        echo "addr ${user_addr} 尚未收到 ibc 代币, sleep 5s (${waited}/${IBC_WAIT_SECONDS})"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: ${IBC_WAIT_SECONDS}s 内未收到 ibc 代币: ${user_addr}" >&2
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
        echo "error: 无法得到 hub 侧 ibc 收款地址，请设置 USER_ADDR_ME_IBC_ROLLAPP" >&2
        exit 1
    fi

    if [ -z "${ME_CHANNEL_ID:-}" ] || [ -z "${RAPP_CHANNEL_ID:-}" ]; then
        ME_CHANNEL_ID="${ME_CHANNEL_ID:-$(rly q channels "${ME_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ROLLAPP_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id' | head -1)}"
        RAPP_CHANNEL_ID="${RAPP_CHANNEL_ID:-$(rly q channels "${ROLLAPP_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ME_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id' | head -1)}"
    fi
    if [ -z "${ME_CHANNEL_ID:-}" ] || [ -z "${RAPP_CHANNEL_ID:-}" ]; then
        echo "error: 找不到 hub-rollapp 通道，请确认 rly-rollapp 已 link。ME_CHANNEL_ID=${ME_CHANNEL_ID:-} RAPP_CHANNEL_ID=${RAPP_CHANNEL_ID:-}" >&2
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
    echo "hub ibc-rollapp 收款地址:"
    med query bank balances "${USER_ADDR_ME_IBC_ROLLAPP}" --home "${ME_NODE1_HOME}"
    echo "rollapp roluser:"
    rollappd query bank balances "${USER_ADDR_ROLLAPP_ADMIN}" --home "${RAPP_NODE1_HOME}"
    echo "rollapp ibc:"
    rollappd query bank balances "${USER_ADDR_ROLLAPP_IBC}" --home "${RAPP_NODE1_HOME}"
}

# hub(global_dao/umec) -> rollapp(ibc 地址)
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

# rollapp(roluser/urax) -> hub(ibc-rollapp 收款地址)
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

# 查询提案。不传 id 则查最新一条，并列出全部提案摘要。
proposal() {
    need_base_dir
    local proposal_id="${1:-}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  查询 gov 提案                                                                #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "全部提案:"
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
        echo "没有可查询的提案"
        return 0
    fi
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  proposal ${proposal_id}                                                     #"
    echo "# ---------------------------------------------------------------------------- #"
    med q gov proposal "${proposal_id}" --home "${ME_NODE1_HOME}"
    echo "votes:"
    med q gov votes "${proposal_id}" --home "${ME_NODE1_HOME}" || true
    echo "upgrade plan:"
    med q upgrade plan --home "${ME_NODE1_HOME}" || true
}

rollapp_upgrade_info_name() {
    local f="${RAPP_NODE1_HOME}/data/upgrade-info.json"
    if [ -f "${f}" ]; then
        jq -r '.name // empty' "${f}" 2>/dev/null || true
    fi
}

wait_rollapp_halt() {
    local halt height waited=0 name
    echo "等待 RollApp 到达 software-upgrade halt 高度"
    while [ "${waited}" -lt 1800 ]; do
        name="$(rollapp_upgrade_info_name)"
        if [[ -n "${name}" ]]; then
            echo "rollapp 已 halt: ${RAPP_NODE1_HOME}/data/upgrade-info.json name=${name}"
            return 0
        fi
        halt=$(rollappd q upgrade plan --home "${RAPP_NODE1_HOME}" -o json 2>/dev/null | jq -r '.plan.height // .height // empty')
        halt=${halt//\"/}
        height="$(rollapp_height 2>/dev/null || true)"
        height=${height//\"/}
        if [[ -n "${halt}" && "${halt}" != "null" && -n "${height}" && "${height}" != "null" ]] && [ "${height}" -ge "${halt}" ] 2>/dev/null; then
            echo "已到 rollapp halt height=${halt} current=${height}"
            return 0
        fi
        if [[ -z "${height}" || "${height}" == "null" ]]; then
            echo "error: rollapp RPC 不可用且没有 upgrade-info.json。" >&2
            echo "这通常是旧 rollappd 连了 v3 Hub 后 panic，不是正常 software-upgrade halt。" >&2
            exit 1
        fi
        echo "等待 rollapp halt: height=${height} plan=${halt:-?} (${waited}/1800)"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: 1800s 内 rollapp 未到达 halt 高度" >&2
    rollappd q upgrade plan --home "${RAPP_NODE1_HOME}" || true
    exit 1
}

# Hub 提案/换镜像之前，rollapp 必须已经在 v2 Hub 上 halt。
# 不能等到 Hub halt 再 rollapp-propose：Hub halt 后节点 panic，settlement RPC 没了。
require_rollapp_halted() {
    local name plan
    if [ "${SKIP_ROLLAPP_GOV_CHECK:-}" = "1" ]; then
        echo "SKIP_ROLLAPP_GOV_CHECK=1，跳过 rollapp halt 检查"
        return 0
    fi
    name="$(rollapp_upgrade_info_name)"
    if [[ -n "${name}" ]]; then
        echo "rollapp 已 halt: upgrade-info name=${name}"
        return 0
    fi
    plan=$(rollappd q upgrade plan --home "${RAPP_NODE1_HOME}" -o json 2>/dev/null | jq -r '.plan.name // .name // empty' || true)
    if [[ -n "${plan}" && "${plan}" != "null" ]]; then
        echo "rollapp 已有 upgrade plan=${plan}，等待 halt（必须在 Hub 仍出块时完成）"
        wait_rollapp_halt
        return 0
    fi
    echo "error: 先让 rollapp 在 v2 Hub 上完成 software-upgrade 并 halt，再动 Hub。" >&2
    echo "顺序: rollapp-propose → rollapp-vote → propose → vote → hub-upgrade → rollapp-upgrade" >&2
    echo "Hub halt 后再提案没用：halt 后 RPC 挂了；hub-upgrade 后再起 v1.0.21 会 GetProposer panic。" >&2
    echo "这条环境若已经 Hub v3 且 rollapp 没提案，只能重置。强制跳过: SKIP_ROLLAPP_GOV_CHECK=1" >&2
    exit 1
}

# 由 global_dao 提交 software-upgrade 到 v3.0.0，然后查询提案。
propose() {
    need_base_dir
    require_rollapp_halted
    local height halt txhash proposal_id
    height="$(hub_height)"
    height="${height//\"/}"
    halt="${HALT_HEIGHT:-$((height + HALT_OFFSET + 20))}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  提交 software-upgrade ${UPGRADE_NAME}  height=${height} halt=${halt}        #"
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
    echo "${key} ${addr} umec=${amount}, 从 global_dao 转入 ${VAL_FUND_AMOUNT} 作手续费"
    txhash=$(med tx bank send global_dao "${addr}" "${VAL_FUND_AMOUNT}" \
        --from global_dao \
        --keyring-backend "${KEYRING_BACKEND_NAME}" \
        --chain-id "${ME_CHAIN_ID}" \
        --home "${ME_NODE1_HOME}" \
        --fees "${VOTE_FEES}" \
        --yes -o json | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
}

# 四个 validator 在同一个函数里投票。不传 id 则投最新提案。
vote() {
    need_base_dir
    local proposal_id="${1:-}"
    local key txhash
    if [ -z "${proposal_id}" ]; then
        proposal_id="$(latest_proposal_id || true)"
    fi
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        echo "error: 没有可投票的提案，请先 propose 或传入 proposal id" >&2
        exit 1
    fi
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  val1-val4 对提案 ${proposal_id} 投 ${VOTE_OPTION}                           #"
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
    echo "#  Hub 已投票。等 halt 后执行 hub-upgrade（换镜像）。不要现在起旧 rollappd。     #"
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
    raise SystemExit("docker-compose.yml 里没有 x-me-template")
rest = text[idx + len(needle):]
m = re.search(r"(?m)^(x-|networks:|services:)", rest)
end = idx + len(needle) + (m.start() if m else len(rest))
block = text[idx:end]
block, n = re.subn(r"(?m)^  image: .+$", f"  image: {image}", block, count=1)
if n != 1:
    raise SystemExit("未能替换 x-me-template 的 image")
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
    raise SystemExit("未能替换 x-me-template 的 command")
block = re.sub(r"(?m)^  user: .+\n", "", block)
path.write_text(text[:idx] + block + text[end:])
print(f"compose hub image -> {image}")
print(f"compose hub home -> {home} min-gas -> {min_gas}")
PY
}

wait_hub_halt() {
    local halt height waited=0
    echo "等待 Hub 到达 software-upgrade halt 高度"
    while [ "${waited}" -lt 1800 ]; do
        halt=$(med q upgrade plan --home "${ME_NODE1_HOME}" -o json 2>/dev/null | jq -r '.plan.height // .height // empty')
        halt=${halt//\"/}
        height="$(hub_height 2>/dev/null || true)"
        height=${height//\"/}
        if [[ -n "${halt}" && "${halt}" != "null" && -n "${height}" && "${height}" != "null" ]] && [ "${height}" -ge "${halt}" ] 2>/dev/null; then
            echo "已到 halt height=${halt} current=${height}"
            return 0
        fi
        if [[ -z "${height}" || "${height}" == "null" ]]; then
            echo "hub RPC 不可用（可能已 halt），继续换 ${MED_V3_IMAGE}"
            return 0
        fi
        echo "等待 halt: height=${height} plan=${halt:-?} (${waited}/1800)"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: 1800s 内未到达 halt 高度" >&2
    med q upgrade plan --home "${ME_NODE1_HOME}" || true
    exit 1
}

wait_hub_block() {
    local waited=0 height
    while [ "${waited}" -lt 180 ]; do
        height="$(hub_height 2>/dev/null || true)"
        height=${height//\"/}
        if [[ -n "${height}" && "${height}" != "null" && "${height}" -ge 1 ]] 2>/dev/null; then
            echo "hub 已出块, height=${height}"
            return 0
        fi
        echo "等待 hub 出块, sleep 5s (${waited}/180)"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: 新 med 180s 内未出块" >&2
    docker compose logs --tail 80 hub-node1 || true
    exit 1
}

# Phase B：halt 后把 compose 的 Hub 镜像换成 med:v3.0.0 并启动。
# 旧 rollappd v1.0.21 连 v3 Hub 会在 IsSequencerVerify 对 GetProposer() 空指针 panic，必须保持停止。
# 换镜像 ≠ 提案：rollapp 的 gov/halt 必须在这之前、Hub 还是 v2 时完成。
hub_upgrade() {
    need_base_dir
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Phase B Hub 升级  image=${MED_V3_IMAGE}  home=${HUB_CONTAINER_HOME}          #"
    echo "# ---------------------------------------------------------------------------- #"
    require_rollapp_halted
    wait_hub_halt
    echo "停 hub + rollapp（halt 期间不要让 rollapp 继续堆 batch）"
    docker compose stop rly-rollapp "${ROLLAPP_SERVICES[@]}" "${HUB_SERVICES[@]}" || true
    echo "pull ${MED_V3_IMAGE}"
    docker pull "${MED_V3_IMAGE}"
    set_compose_hub_image "${MED_V3_IMAGE}"
    docker compose up -d "${HUB_SERVICES[@]}"
    docker compose stop rly-rollapp "${ROLLAPP_SERVICES[@]}" || true
    docker compose ps "${HUB_SERVICES[@]}"
    wait_hub_block
    echo "upgrade applied:"
    med q upgrade applied "${UPGRADE_NAME}" --home "${ME_NODE1_HOME}" || true
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Hub v3 已起来。不要 docker compose up rollapp-node*（仍是 v1.0.21，会 panic） #"
    echo "#  下一步：ROLLAPP_V3_IMAGE=... ./scripts/setup_v2_all_func.sh rollapp-upgrade     #"
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

# 查询 rollapp gov 提案。不传 id 则查最新一条。
rollapp_proposal() {
    need_base_dir
    local proposal_id="${1:-}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  查询 rollapp gov 提案                                                       #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "validators:"
    rollappd q staking validators --home "${RAPP_NODE1_HOME}" -o json | jq -r '
        (.validators // [])
        | if length == 0 then "  (none)"
          else .[] | "  moniker=\(.description.moniker) status=\(.status) tokens=\(.tokens)"
          end
    '
    echo "全部提案:"
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
        echo "没有可查询的 rollapp 提案"
        return 0
    fi
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollapp proposal ${proposal_id}                                             #"
    echo "# ---------------------------------------------------------------------------- #"
    rollappd q gov proposal "${proposal_id}" --home "${RAPP_NODE1_HOME}"
    echo "votes:"
    rollappd q gov votes "${proposal_id}" --home "${RAPP_NODE1_HOME}" || true
    echo "upgrade plan:"
    rollappd q upgrade plan --home "${RAPP_NODE1_HOME}" || true
}

# 由 roluser 提交 rollapp software-upgrade v3.0.0（submit-legacy-proposal + --no-validate）。
# handler 会新增 hubgenesis store，并 InitGenesis rollappparams（DA=me-da）。
rollapp_propose() {
    need_base_dir
    local height halt txhash proposal_id
    height="$(rollapp_height)"
    height="${height//\"/}"
    halt="${ROLLAPP_HALT_HEIGHT:-$((height + ROLLAPP_HALT_OFFSET + 50))}"
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollapp software-upgrade ${UPGRADE_NAME}  height=${height} halt=${halt}     #"
    echo "#  新增模块: hubgenesis store + rollappparams (DA=${ROLLAPP_DA_LAYER})         #"
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

# 创世 gentx 只有 roluser 一个 staking 验证人；对链上每个 bonded validator 若本地有同名 key 则投票。
rollapp_vote() {
    need_base_dir
    local proposal_id="${1:-}"
    local key txhash
    if [ -z "${proposal_id}" ]; then
        proposal_id="$(rollapp_latest_proposal_id || true)"
    fi
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        echo "error: 没有可投票的 rollapp 提案，请先 rollapp-propose 或传入 proposal id" >&2
        exit 1
    fi
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollapp 验证人对提案 ${proposal_id} 投 ${VOTE_OPTION}                       #"
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
    echo "#  投票完成。下一步等 rollapp halt，再 Hub propose（不要先 hub-upgrade）。       #"
    echo "# ---------------------------------------------------------------------------- #"
    wait_rollapp_halt
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  rollapp 已 halt。现在可以: propose → vote → hub-upgrade → rollapp-upgrade     #"
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
        # 不用 mv：有的环境 alias 成 mv -i，且 docker 写出的目标文件可能是 root
        command cp "${tmp}" "${file}"
        rm -f "${tmp}"
    else
        printf '%s\n' "${line}" >> "${file}"
        rm -f "${tmp}"
    fi
}

need_rollapp_v3_image() {
    if [ -z "${ROLLAPP_V3_IMAGE}" ]; then
        echo "error: 请设置 ROLLAPP_V3_IMAGE，例如 harbor.starex.xyz/rollapp/rollappd:<new-tag>" >&2
        exit 1
    fi
}

backup_rollapp_home() {
    local home="$1"
    local bak="${home}.bak-$(date +%Y%m%d%H%M%S)"
    echo "备份 ${home} -> ${bak}"
    cp -a "${home}" "${bak}"
}

run_3d_migration() {
    local home="$1"
    local out rc=0
    echo "rollappd run-3d-migration --home ${home} --rollapp-param-da ${ROLLAPP_DA_LAYER}"
    # mehub_v3：命令挂在 rollappd 上。宿主机若仍是旧二进制，改用 v3 镜像。
    if rollappd run-3d-migration --help >/dev/null 2>&1; then
        out=$(rollappd run-3d-migration --home "${home}" --rollapp-param-da "${ROLLAPP_DA_LAYER}" 2>&1) || rc=$?
    else
        if [ -z "${ROLLAPP_V3_IMAGE}" ]; then
            echo "error: 宿主机 rollappd 没有 run-3d-migration，请设置 ROLLAPP_V3_IMAGE 用 v3 镜像跑" >&2
            exit 1
        fi
        out=$(docker run --rm \
            -v "${home}:${ROLLAPP_CONTAINER_HOME}" \
            --entrypoint rollappd \
            "${ROLLAPP_V3_IMAGE}" \
            run-3d-migration --home "${ROLLAPP_CONTAINER_HOME}" --rollapp-param-da "${ROLLAPP_DA_LAYER}" 2>&1) || rc=$?
    fi
    printf '%s\n' "${out}"
    if echo "${out}" | grep -qE '3D dymint store migration successful|3D migration is not needed'; then
        return 0
    fi
    echo "error: 3D migration 失败 (exit ${rc})" >&2
    exit 1
}

# 离线对三个 rollapp 节点执行 3D store 迁移。必须先停进程。
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
            echo "跳过不存在的 ${home}"
            continue
        fi
        backup_rollapp_home "${home}"
        run_3d_migration "${home}"
    done
    echo "3D migration 完成。节点保持停止，换 v3 镜像后再 docker compose up -d ${ROLLAPP_SERVICES[*]}"
}

# v3 dymint：da_layer / da_config 必须是等长数组。
# 旧格式 da_layer = "me-da"、da_config = '{json}' 不能再用 awk -v 拷贝（JSON 引号会被截断成空串）。
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
    raise SystemExit(f"error: 找不到可用的 da_config JSON（查过 {cfg} / .pre-v3 / bak-*）")

json_str = json.dumps(best, separators=(",", ":"))
if "'" in json_str:
    raise SystemExit("error: da_config JSON 含单引号，无法写入 TOML 单引号字符串")

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

# 用新 rollappd 的默认 dymint.toml，再把旧文件里的 settlement / keyring 抄回去。
# DA 单独按 v3 数组格式从 .pre-v3 / bak 还原，不要用 awk 拷贝旧标量。
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
    if ! docker run --rm \
        -v "${tmp}:${ROLLAPP_CONTAINER_HOME}" \
        --entrypoint rollappd \
        "${ROLLAPP_V3_IMAGE}" \
        init tmpnode --chain-id "${ROLLAPP_CHAIN_ID}" --home "${ROLLAPP_CONTAINER_HOME}" >/dev/null 2>&1; then
        echo "新 rollappd init 失败，保留原 dymint.toml 并只核对 DA/settlement"
        rm -rf "${tmp}"
        return 0
    fi
    if [ ! -f "${tmp}/config/dymint.toml" ]; then
        echo "新镜像未生成 dymint.toml，保留原文件"
        rm -rf "${tmp}"
        return 0
    fi
    # docker 以 root 写 tmp，拷到当前用户文件再改，避免 mv Permission denied
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

# 现网已经迁完 3D、只差 v3 DA 数组时用这个，不要再跑 rollapp-upgrade。
fix_da_config() {
    need_base_dir
    local node home
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  修复 v3 dymint DA 配置  da=${ROLLAPP_DA_LAYER}                              #"
    echo "# ---------------------------------------------------------------------------- #"
    docker compose stop "${ROLLAPP_SERVICES[@]}" || true
    for node in "${ROLLAPP_NODE_HOMES[@]}"; do
        home="${RAPP_NODES_HOME}/${node}"
        if [ ! -f "${home}/config/dymint.toml" ]; then
            echo "跳过不存在的 ${home}"
            continue
        fi
        fix_dymint_runtime_keys "${home}"
    done
    set_compose_container_homes
    docker compose up -d "${ROLLAPP_SERVICES[@]}"
    docker compose ps "${ROLLAPP_SERVICES[@]}"
    echo "看 rollapp-node1 日志确认 DA 已起来："
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
        raise SystemExit(f"docker-compose.yml 里没有 {needle}")
    rest = text[idx + len(needle):]
    m = re.search(r"(?m)^(x-|networks:|services:)", rest)
    end = idx + len(needle) + (m.start() if m else len(rest))
    block = text[idx:end]
    if image:
        block, n = re.subn(r"(?m)^  image: .+$", f"  image: {image}", block, count=1)
        if n != 1:
            raise SystemExit(f"未能替换 {needle} 的 image")
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
        raise SystemExit(f"未能替换 {needle} 的 command")
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
            echo "rollapp 已出块, height=${height}"
            return 0
        fi
        echo "等待 rollapp 出块, sleep 5s (${waited}/120)"
        sleep 5
        waited=$((waited + 5))
    done
    echo "error: 新 rollappd 120s 内未出块" >&2
    docker compose logs --tail 80 rollapp-node1 || true
    exit 1
}

# Phase C：Hub 已是 v3 之后，升级 RollApp / dymint（备份 + 3D migration + 换镜像）。
rollapp_upgrade() {
    need_base_dir
    need_rollapp_v3_image
    local node home
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  Phase C RollApp 升级  image=${ROLLAPP_V3_IMAGE}  da=${ROLLAPP_DA_LAYER}     #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "Hub upgrade applied:"
    med q upgrade applied "${UPGRADE_NAME}" --home "${ME_NODE1_HOME}" || true
    echo "停 rollapp + rly-rollapp（避免 settlement 断开后继续堆 batch）"
    docker compose stop rly-rollapp "${ROLLAPP_SERVICES[@]}" || true

    echo "pull ${ROLLAPP_V3_IMAGE}"
    docker pull "${ROLLAPP_V3_IMAGE}"

    for node in "${ROLLAPP_NODE_HOMES[@]}"; do
        home="${RAPP_NODES_HOME}/${node}"
        if [ ! -d "${home}/data" ] && [ ! -f "${home}/config/dymint.toml" ]; then
            echo "跳过不存在的 ${home}"
            continue
        fi
        backup_rollapp_home "${home}"
        run_3d_migration "${home}"
        merge_dymint_toml "${home}"
        fix_dymint_runtime_keys "${home}"
        ensure_rollapp_min_gas "${home}"
    done

    set_compose_rollapp_image
    docker compose up -d "${ROLLAPP_SERVICES[@]}"
    docker compose ps "${ROLLAPP_SERVICES[@]}"
    wait_rollapp_block
    echo "Hub latest-state-info:"
    med q rollapp latest-state-info "${ROLLAPP_CHAIN_ID}" --home "${ME_NODE1_HOME}" || true
    echo "Phase C 完成。relayer 请另启：docker compose up -d rly-rollapp"
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
