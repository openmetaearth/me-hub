#!/bin/bash
# Post-v3 functional regression. Default runs all of: Hub transfer, IBC, KYC, flexible delegate, fixed deposit, governance.
# The governance proposal is voted no so it does not schedule another software-upgrade halt.
#
#   cd /data/docker-beta
#   ./scripts/setup_v3_test.sh all
#   ./scripts/setup_v3_test.sh transfer|ibc|kyc|delegate|fixed|proposal
set -euo pipefail

PROJECT_ENV="${PROJECT_ENV:-beta}"
BASE_DIR="${BASE_DIR:-/data/docker-${PROJECT_ENV}}"
BIN_DIR="${BASE_DIR}/bin"
KEYRING_BACKEND_NAME="${KEYRING_BACKEND_NAME:-test}"
ME_CHAIN_ID="${ME_CHAIN_ID:-mechain_400-1}"
ROLLAPP_CHAIN_ID="${ROLLAPP_CHAIN_ID:-mecheckin_401-1}"
ME_NODE1_HOME="${ME_NODE1_HOME:-${BASE_DIR}/nodes/hub-nodes/node1}"
RAPP_NODE1_HOME="${RAPP_NODE1_HOME:-${BASE_DIR}/nodes/rollapp-nodes/node1}"
FEES="${FEES:-200000umec}"
TRANSFER_AMOUNT="${TRANSFER_AMOUNT:-20000000umec}"
DELEGATE_AMOUNT="${DELEGATE_AMOUNT:-1000000umec}"
FIXED_AMOUNT="${FIXED_AMOUNT:-10000000umec}"
FIXED_TERM="${FIXED_TERM:-1}"
FIXED_REGION="${FIXED_REGION:-me_earth}"
FIXED_RATE="${FIXED_RATE:-0.1}"
GOV_DEPOSIT="${GOV_DEPOSIT:-100000000umec}"
VOTE_FEES="${VOTE_FEES:-10000umec}"
HUB_VAL_KEYS=(val1 val2 val3 val4)
TEST_KEY="${TEST_KEY:-v3test}"
# Test account used only by this script. Unrelated to the install-script mnemonic.
TEST_MNEMONIC="${TEST_MNEMONIC:-abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about}"
if [ -z "${TEST_HD_PATH:-}" ]; then
    TEST_HD_PATH="m/44'/118'/0'/0/99"
fi
KYC_DID="${KYC_DID:-v3test0000001}"
KYC_LEVEL="${KYC_LEVEL:-2}"

usage() {
    cat <<EOF
Usage: $0 <command>

  all         Run transfer, kyc, delegate, fixed, proposal, then ibc
  transfer    global_dao sends ${TRANSFER_AMOUNT} to ${TEST_KEY}
  kyc         global_dao issues KYC for ${TEST_KEY} (region=${FIXED_REGION})
  delegate    ${TEST_KEY} delegates ${DELEGATE_AMOUNT} (region follows KYC, else experience_region)
  fixed       Ensure a fixed-deposit config exists, then ${TEST_KEY} deposits ${FIXED_AMOUNT} for term ${FIXED_TERM}
  proposal    Submit a software-upgrade proposal; four validators vote no
  ibc         Call setup_v2_all_func.sh ibc

Env:
  BASE_DIR=${BASE_DIR}
  ME_CHAIN_ID=${ME_CHAIN_ID}
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
    med tx "$@" \
        --from "${from}" \
        --keyring-backend "${KEYRING_BACKEND_NAME}" \
        --chain-id "${ME_CHAIN_ID}" \
        --home "${ME_NODE1_HOME}" \
        --fees "${FEES}" \
        --gas auto --gas-adjustment 1.5 \
        -y -o json
}

ensure_test_key() {
    local kb=(--keyring-backend "${KEYRING_BACKEND_NAME}" --home "${ME_NODE1_HOME}")
    if med keys show "${TEST_KEY}" -a "${kb[@]}" >/dev/null 2>&1; then
        return 0
    fi
    echo "Create test account ${TEST_KEY}"
    echo -e "${TEST_MNEMONIC}" | med keys add "${TEST_KEY}" \
        --keyring-backend "${KEYRING_BACKEND_NAME}" \
        --home "${ME_NODE1_HOME}" \
        --key-type secp256k1 \
        --hd-path "${TEST_HD_PATH}" \
        --recover --coin-type 118 >/dev/null
}

test_addr() {
    med keys show "${TEST_KEY}" -a \
        --keyring-backend "${KEYRING_BACKEND_NAME}" --home "${ME_NODE1_HOME}"
}

test_pubkey() {
    med keys show "${TEST_KEY}" -p \
        --keyring-backend "${KEYRING_BACKEND_NAME}" --home "${ME_NODE1_HOME}"
}

cmd_transfer() {
    need_base_dir
    ensure_test_key
    local to
    to="$(test_addr)"
    echo "# hub bank send global_dao -> ${TEST_KEY} ${TRANSFER_AMOUNT}"
    local txhash
    txhash=$(hub_tx global_dao bank send global_dao "${to}" "${TRANSFER_AMOUNT}" | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    med q bank balances "${to}" --home "${ME_NODE1_HOME}"
}

cmd_kyc() {
    need_base_dir
    ensure_test_key
    local addr pubkey dao
    addr="$(test_addr)"
    pubkey="$(test_pubkey)"
    dao=$(med keys show global_dao -a --keyring-backend "${KEYRING_BACKEND_NAME}" --home "${ME_NODE1_HOME}")
    if med q kyc did "${addr}" --home "${ME_NODE1_HOME}" >/dev/null 2>&1; then
        echo "KYC already exists: ${addr}"
        med q kyc did "${addr}" --home "${ME_NODE1_HOME}"
        return 0
    fi
    echo "# kyc approve ${KYC_DID} ${addr}"
    local txhash
    txhash=$(hub_tx global_dao kyc approve \
        "${KYC_DID}" "${FIXED_REGION}" "${addr}" "${pubkey}" \
        "${KYC_LEVEL}" "ipfs://v3-test" "v3-test-hash" "${dao}" | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    med q kyc did "${addr}" --home "${ME_NODE1_HOME}"
    med q kyc kyc "${KYC_DID}" --home "${ME_NODE1_HOME}"
}

cmd_delegate() {
    need_base_dir
    ensure_test_key
    local addr
    addr="$(test_addr)"
    echo "# delegate ${DELEGATE_AMOUNT} from ${TEST_KEY}"
    local txhash
    txhash=$(hub_tx "${TEST_KEY}" staking delegate "${DELEGATE_AMOUNT}" | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    med q staking delegation "${addr}" --home "${ME_NODE1_HOME}"
}

cmd_fixed() {
    need_base_dir
    ensure_test_key
    if ! med q staking show-fixed-deposit-cfg "${FIXED_REGION}" --home "${ME_NODE1_HOME}" -o json \
        | jq -e --argjson term "${FIXED_TERM}" '
            [.RegionFixedDepositCfgs[]?.FixedDepositCfgs[]? | select(.term == $term and .status == "FIXED_DEPOSIT_CFG_ACTIVE")]
            | length > 0' >/dev/null 2>&1; then
        echo "# create fixed-deposit config ${FIXED_REGION} term=${FIXED_TERM} rate=${FIXED_RATE}"
        local txhash
        txhash=$(hub_tx global_dao staking new-fixed-deposit-cfg \
            "${FIXED_REGION}" "${FIXED_TERM}" "${FIXED_RATE}" | jq -r '.txhash')
        check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    fi
    local addr
    addr="$(test_addr)"
    echo "# deposit-fixed ${FIXED_AMOUNT} term=${FIXED_TERM} from ${TEST_KEY}"
    local txhash
    txhash=$(hub_tx "${TEST_KEY}" staking deposit-fixed "${FIXED_AMOUNT}" "${FIXED_TERM}" | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    med q staking show-fixed-deposit-by-acct "${addr}" all --home "${ME_NODE1_HOME}"
}

cmd_proposal() {
    need_base_dir
    local height halt gov_addr proposal_file txhash proposal_id key
    height=$(med status --home "${ME_NODE1_HOME}" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height')
    height=${height//\"/}
    halt=$((height + 1000000))
    gov_addr=$(med q auth module-account gov --home "${ME_NODE1_HOME}" -o json | jq -r '.account.base_account.address // .account.address')
    proposal_file="${BASE_DIR}/v3-test-proposal.json"
    cat > "${proposal_file}" <<EOF
{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "authority": "${gov_addr}",
      "plan": {
        "name": "v3-test",
        "time": "0001-01-01T00:00:00Z",
        "height": "${halt}",
        "info": "v3 regression, vote no",
        "upgraded_client_state": null
      }
    }
  ],
  "metadata": "v3-test",
  "deposit": "${GOV_DEPOSIT}",
  "title": "v3-test software upgrade",
  "summary": "regression only; validators vote no so this plan is not adopted",
  "expedited": false
}
EOF
    echo "# submit software-upgrade proposal halt=${halt}, then vote no"
    txhash=$(hub_tx global_dao gov submit-proposal "${proposal_file}" | jq -r '.txhash')
    check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    proposal_id=$(med q tx "${txhash}" --home "${ME_NODE1_HOME}" -o json | jq -r '
        [ .logs[]?.events[]?, .events[]? ]
        | map(select(.type == "submit_proposal"))
        | .[].attributes[]?
        | select(.key == "proposal_id")
        | .value' | head -1)
    if [ -z "${proposal_id}" ] || [ "${proposal_id}" = "null" ]; then
        echo "error: failed to parse proposal_id" >&2
        exit 1
    fi
    echo "proposal_id=${proposal_id}"
    for key in "${HUB_VAL_KEYS[@]}"; do
        echo "vote no --from ${key}"
        txhash=$(med tx gov vote "${proposal_id}" no \
            --from "${key}" \
            --keyring-backend "${KEYRING_BACKEND_NAME}" \
            --chain-id "${ME_CHAIN_ID}" \
            --home "${ME_NODE1_HOME}" \
            --fees "${VOTE_FEES}" \
            -y -o json | jq -r '.txhash')
        check_tx_status med "${txhash}" "--home ${ME_NODE1_HOME}"
    done
    med q gov proposal "${proposal_id}" --home "${ME_NODE1_HOME}"
    echo "Four validators voted no. After the voting period the proposal is rejected and no halt is scheduled."
}

cmd_ibc() {
    need_base_dir
    local func_sh
    func_sh="$(cd "$(dirname "$0")" && pwd)/setup_v2_all_func.sh"
    if [ ! -f "${func_sh}" ]; then
        func_sh="${BASE_DIR}/setup_v2_all_func.sh"
    fi
    PROJECT_ENV="${PROJECT_ENV}" BASE_DIR="${BASE_DIR}" "${func_sh}" ibc
}

cmd_all() {
    cmd_transfer
    cmd_kyc
    cmd_delegate
    cmd_fixed
    cmd_proposal
    cmd_ibc
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    case "${1:-}" in
        all) cmd_all ;;
        transfer) cmd_transfer ;;
        kyc) cmd_kyc ;;
        delegate) cmd_delegate ;;
        fixed) cmd_fixed ;;
        proposal) cmd_proposal ;;
        ibc) cmd_ibc ;;
        help|-h|--help|"") usage ;;
        *)
            usage
            exit 1
            ;;
    esac
fi
