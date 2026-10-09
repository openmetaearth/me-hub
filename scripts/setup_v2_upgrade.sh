#!/bin/bash
# Server-side Hub v2 → v3.0.0 upgrade drill, driven by release images:
#   V2: ghcr.io/openmetaearth/med:v2.0.17
#   V3: ghcr.io/openmetaearth/med:v3.0.0
#
# Host needs: docker, jq. Chain data lives on the host and is mounted at
# /root/.mechain inside the container (same path the release image uses).
#
#   ./scripts/setup_v2_upgrade.sh pull
#   ./scripts/setup_v2_upgrade.sh init          # wipe + v2 genesis
#   ./scripts/setup_v2_upgrade.sh start-v2
#   ./scripts/setup_v2_upgrade.sh propose       # software-upgrade v3.0.0 + vote
#   ./scripts/setup_v2_upgrade.sh wait-halt
#   ./scripts/setup_v2_upgrade.sh start-v3      # same data dir, v3 image
#
# Optional:
#   HOME_DIR=/data/mechain-upgrade CONTAINER_NAME=mechain-upgrade \
#     MED_V2_IMAGE=ghcr.io/openmetaearth/med:v2.0.17 \
#     MED_V3_IMAGE=ghcr.io/openmetaearth/med:v3.0.0 \
#     ./scripts/setup_v2_upgrade.sh init
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

MED_V2_IMAGE="${MED_V2_IMAGE:-ghcr.io/openmetaearth/med:v2.0.17}"
MED_V3_IMAGE="${MED_V3_IMAGE:-ghcr.io/openmetaearth/med:v3.0.0}"
HOME_DIR="${HOME_DIR:-$HOME/.mechain-upgrade}"
CONTAINER_HOME="/root/.mechain"
CONTAINER_NAME="${CONTAINER_NAME:-mechain-upgrade}"
CHAIN_ID="${CHAIN_ID:-mechain_2401-1}"
MONIKER_NAME="${MONIKER_NAME:-local}"
KEY_NAME="${KEY_NAME:-global_dao}"
KEY_NAME_SEQUENCER="${KEY_NAME_SEQUENCER:-sequencer}"
MNEMONIC="${MNEMONIC:-curtain hat remain song receive tower stereo hope frog cheap brown plate raccoon post reflect wool sail salmon game salon group glimpse adult shift}"
STAKING_AMOUNT="${STAKING_AMOUNT:-10000000000000000umec}"
DAO_GENESIS_AMOUNT="${DAO_GENESIS_AMOUNT:-1000000000000000000umec}"
VOTING_PERIOD="${VOTING_PERIOD:-15s}"
UNBONDING_TIME="${UNBONDING_TIME:-10s}"
UPGRADE_NAME="${UPGRADE_NAME:-v3.0.0}"
HALT_OFFSET="${HALT_OFFSET:-50}"
DEPOSIT="${DEPOSIT:-100000000umec}"
FEES="${FEES:-200000000umec}"

# Listen addresses inside the container. Host publish ports are mapped 1:1.
SETTLEMENT_ADDR="${SETTLEMENT_ADDR:-0.0.0.0:36657}"
P2P_ADDRESS="${P2P_ADDRESS:-0.0.0.0:26656}"
GRPC_ADDRESS="${GRPC_ADDRESS:-0.0.0.0:9090}"
GRPC_WEB_ADDRESS="${GRPC_WEB_ADDRESS:-0.0.0.0:9091}"
API_ADDRESS="${API_ADDRESS:-0.0.0.0:1317}"
JSONRPC_ADDRESS="${JSONRPC_ADDRESS:-0.0.0.0:8545}"
JSONRPC_WS_ADDRESS="${JSONRPC_WS_ADDRESS:-0.0.0.0:8546}"

GENESIS_FILE="$HOME_DIR/config/genesis.json"
TENDERMINT_CONFIG_FILE="$HOME_DIR/config/config.toml"
CLIENT_CONFIG_FILE="$HOME_DIR/config/client.toml"
APP_CONFIG_FILE="$HOME_DIR/config/app.toml"

usage() {
  cat <<EOF
Usage: $0 <command>

Commands:
  pull        Pull $MED_V2_IMAGE and $MED_V3_IMAGE
  init        Wipe HOME_DIR and write a v2 genesis (uses v2 image)
  start-v2    Start the v2 container on HOME_DIR
  start-v3    Start the v3 container on the same HOME_DIR (after halt)
  stop        Stop and remove the container
  logs        Follow container logs
  status      med status inside the running container
  exec ...    med ... inside the running container
  propose     Submit + vote software-upgrade $UPGRADE_NAME
  wait-halt   Wait until the v2 container hits the upgrade height
  help        Show this help

Images:
  MED_V2_IMAGE=$MED_V2_IMAGE
  MED_V3_IMAGE=$MED_V3_IMAGE
  HOME_DIR=$HOME_DIR
  CONTAINER_NAME=$CONTAINER_NAME

If GHCR is private:
  echo "\$GHCR_TOKEN" | docker login ghcr.io -u USER --password-stdin
EOF
}

need_bin() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command not found: $1" >&2
    exit 1
  fi
}

host_port() {
  printf '%s\n' "${1##*:}"
}

docker_publish_args() {
  printf -- '-p %s:%s -p %s:%s -p %s:%s -p %s:%s -p %s:%s -p %s:%s -p %s:%s' \
    "$(host_port "$SETTLEMENT_ADDR")" "$(host_port "$SETTLEMENT_ADDR")" \
    "$(host_port "$P2P_ADDRESS")" "$(host_port "$P2P_ADDRESS")" \
    "$(host_port "$GRPC_ADDRESS")" "$(host_port "$GRPC_ADDRESS")" \
    "$(host_port "$GRPC_WEB_ADDRESS")" "$(host_port "$GRPC_WEB_ADDRESS")" \
    "$(host_port "$API_ADDRESS")" "$(host_port "$API_ADDRESS")" \
    "$(host_port "$JSONRPC_ADDRESS")" "$(host_port "$JSONRPC_ADDRESS")" \
    "$(host_port "$JSONRPC_WS_ADDRESS")" "$(host_port "$JSONRPC_WS_ADDRESS")"
}

# Run a one-shot med command in IMAGE against the host data dir.
med_image() {
  local image="$1"
  shift
  # -i so keys --recover can read the mnemonic from stdin. No -t (piped input).
  # shellcheck disable=SC2086
  docker run --rm -i \
    --entrypoint med \
    -v "$HOME_DIR:$CONTAINER_HOME" \
    "$image" \
    "$@" --home "$CONTAINER_HOME"
}

med_v2() {
  med_image "$MED_V2_IMAGE" "$@"
}

container_running() {
  docker inspect -f '{{.State.Running}}' "$CONTAINER_NAME" 2>/dev/null | grep -q true
}

med_exec() {
  if ! container_running; then
    echo "error: container $CONTAINER_NAME is not running" >&2
    exit 1
  fi
  docker exec -i "$CONTAINER_NAME" med --home "$CONTAINER_HOME" "$@"
}

cmd_pull() {
  need_bin docker
  echo "pulling $MED_V2_IMAGE"
  docker pull "$MED_V2_IMAGE"
  echo "pulling $MED_V3_IMAGE"
  docker pull "$MED_V3_IMAGE"
}

cmd_stop() {
  need_bin docker
  if docker inspect "$CONTAINER_NAME" >/dev/null 2>&1; then
    docker rm -f "$CONTAINER_NAME" >/dev/null
    echo "removed $CONTAINER_NAME"
  else
    echo "container $CONTAINER_NAME not present"
  fi
}

start_container() {
  local image="$1"
  local label="$2"
  need_bin docker
  if docker inspect "$CONTAINER_NAME" >/dev/null 2>&1; then
    echo "error: container $CONTAINER_NAME already exists; run: $0 stop" >&2
    exit 1
  fi
  if [ ! -f "$GENESIS_FILE" ]; then
    echo "error: genesis not found at $GENESIS_FILE (run: $0 init)" >&2
    exit 1
  fi
  # Do not use --restart: v2 is supposed to stay dead after halt.
  # shellcheck disable=SC2046
  docker run -d \
    --name "$CONTAINER_NAME" \
    --entrypoint med \
    -v "$HOME_DIR:$CONTAINER_HOME" \
    $(docker_publish_args) \
    "$image" \
    start --home "$CONTAINER_HOME" >/dev/null
  echo "started $label ($image) as $CONTAINER_NAME"
  echo "logs: $0 logs"
}

cmd_start_v2() {
  start_container "$MED_V2_IMAGE" "v2"
}

cmd_start_v3() {
  start_container "$MED_V3_IMAGE" "v3"
}

cmd_logs() {
  need_bin docker
  docker logs -f --tail 200 "$CONTAINER_NAME"
}

cmd_status() {
  med_exec status
}

cmd_exec() {
  med_exec "$@"
}

block_height() {
  med_exec status --output json | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height'
}

cmd_propose() {
  need_bin jq
  local height halt
  height="$(block_height)"
  halt="${HALT_HEIGHT:-$((height + HALT_OFFSET))}"
  echo "current height=$height halt_height=$halt plan=$UPGRADE_NAME"

  med_exec tx gov submit-legacy-proposal software-upgrade "$UPGRADE_NAME" \
    --title "$UPGRADE_NAME" \
    --description "ME-Hub SDK 0.50 / IBC v8 / settlement v3" \
    --upgrade-height "$halt" \
    --upgrade-info '{"binaries":{}}' \
    --no-validate \
    --deposit "$DEPOSIT" \
    --from "$KEY_NAME" \
    --keyring-backend test \
    --chain-id "$CHAIN_ID" \
    --gas auto --gas-adjustment 1.5 \
    --fees "$FEES" \
    --yes

  # Voting window is 15s in this genesis; vote immediately.
  sleep 2
  local proposal_id
  proposal_id="$(med_exec q gov proposals --output json | jq -r '
    .proposals[-1].proposal_id // .proposals[-1].id // .proposals[-1].proposalId
  ')"
  if [ -z "$proposal_id" ] || [ "$proposal_id" = "null" ]; then
    echo "error: could not read proposal id" >&2
    exit 1
  fi
  echo "voting yes on proposal $proposal_id"
  med_exec tx gov vote "$proposal_id" yes \
    --from "$KEY_NAME" \
    --keyring-backend test \
    --chain-id "$CHAIN_ID" \
    --fees "$FEES" \
    --yes
  echo "plan:"
  med_exec q upgrade plan || true
}

cmd_wait_halt() {
  need_bin docker
  echo "waiting for $CONTAINER_NAME to halt on $UPGRADE_NAME ..."
  local i=0
  while true; do
    if docker logs "$CONTAINER_NAME" 2>&1 | grep -q "UPGRADE \"${UPGRADE_NAME}\" NEEDED"; then
      echo "upgrade halt detected"
      # Give the process a moment to exit.
      sleep 2
      cmd_stop
      return 0
    fi
    if ! container_running; then
      echo "container exited before halt log was seen; last logs:" >&2
      docker logs --tail 80 "$CONTAINER_NAME" >&2 || true
      exit 1
    fi
    i=$((i + 1))
    if [ "$i" -ge 360 ]; then
      echo "error: timed out waiting for halt (~30m)" >&2
      exit 1
    fi
    sleep 5
  done
}

cmd_init() {
  need_bin docker
  need_bin jq

  genesis_config_cmds="${GENESIS_CONFIG_CMDS:-$SCRIPT_DIR/src/genesis_config_commands_v2.sh}"
  if [ ! -f "$genesis_config_cmds" ]; then
    echo "error: v2 genesis helpers not found: $genesis_config_cmds" >&2
    echo "this file is SDK 0.47-only; do not source scripts/src/genesis_config_commands.sh" >&2
    exit 1
  fi
  # shellcheck source=/dev/null
  . "$genesis_config_cmds"

  if docker inspect "$CONTAINER_NAME" >/dev/null 2>&1; then
    echo "removing existing container $CONTAINER_NAME"
    docker rm -f "$CONTAINER_NAME" >/dev/null
  fi

  echo "wiping $HOME_DIR and initializing v2 genesis with $MED_V2_IMAGE"
  rm -rf "$HOME_DIR"
  mkdir -p "$HOME_DIR"

  if ! docker image inspect "$MED_V2_IMAGE" >/dev/null 2>&1; then
    docker pull "$MED_V2_IMAGE"
  fi

  med_v2 init "$MONIKER_NAME" --chain-id="$CHAIN_ID"

  sed -i'' -e "/\[rpc\]/,+3 s/laddr *= .*/laddr = \"tcp:\/\/$SETTLEMENT_ADDR\"/" "$TENDERMINT_CONFIG_FILE"
  sed -i'' -e "/\[p2p\]/,+3 s/laddr *= .*/laddr = \"tcp:\/\/$P2P_ADDRESS\"/" "$TENDERMINT_CONFIG_FILE"
  sed -i'' -e 's/^timeout_propose = .*/timeout_propose = "1s"/' "$TENDERMINT_CONFIG_FILE"
  sed -i'' -e 's/^timeout_commit = .*/timeout_commit = "1s"/' "$TENDERMINT_CONFIG_FILE"
  sed -i'' -e "/\[grpc\]/,+6 s/address *= .*/address = \"$GRPC_ADDRESS\"/" "$APP_CONFIG_FILE"
  sed -i'' -e "/\[grpc-web\]/,+7 s/address *= .*/address = \"$GRPC_WEB_ADDRESS\"/" "$APP_CONFIG_FILE"
  sed -i'' -e "/\[json-rpc\]/,+6 s/address *= .*/address = \"$JSONRPC_ADDRESS\"/" "$APP_CONFIG_FILE"
  sed -i'' -e "/\[json-rpc\]/,+9 s/^ws-address *= .*/ws-address = \"$JSONRPC_WS_ADDRESS\"/" "$APP_CONFIG_FILE"
  sed -i'' -e '/\[api\]/,+3 s/enable *= .*/enable = true/' "$APP_CONFIG_FILE"
  sed -i'' -e "/\[api\]/,+9 s/address *= .*/address = \"tcp:\/\/$API_ADDRESS\"/" "$APP_CONFIG_FILE"
  sed -i'' -e 's/^minimum-gas-prices *= .*/minimum-gas-prices = "0.02umec"/' "$APP_CONFIG_FILE"
  sed -i'' -e "s/^chain-id *= .*/chain-id = \"$CHAIN_ID\"/" "$CLIENT_CONFIG_FILE"
  sed -i'' -e "s/^keyring-backend *= .*/keyring-backend = \"test\"/" "$CLIENT_CONFIG_FILE"
  sed -i'' -e "s/^node *= .*/node = \"tcp:\/\/127.0.0.1:$(host_port "$SETTLEMENT_ADDR")\"/" "$CLIENT_CONFIG_FILE"

  set_consenus_params
  set_gov_params
  set_hub_params
  set_misc_params
  set_EVM_params
  set_bank_denom_metadata
  set_epochs_params
  set_incentives_params

  jq --arg p "$VOTING_PERIOD" '.app_state.gov.voting_params.voting_period = $p' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq --arg p "$VOTING_PERIOD" '.app_state.gov.params.voting_period = $p' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq --arg p "$VOTING_PERIOD" '.app_state.gov.params.max_deposit_period = $p' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq --arg p "$VOTING_PERIOD" '.app_state.gov.deposit_params.max_deposit_period = $p' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"

  echo "$MNEMONIC" | med_v2 keys add "$KEY_NAME" --recover --keyring-backend test
  local dao_addr
  dao_addr="$(med_v2 keys show "$KEY_NAME" -a --keyring-backend test)"
  med_v2 add-genesis-account "$dao_addr" "$DAO_GENESIS_AMOUNT"
  med_v2 add-genesis-stake-pool
  med_v2 add-genesis-m-accounts
  med_v2 gen-relayers "me1frjhlw9slyy7mrhmk0r4vytkyldxqtkf326amv,me1c5zp26c0gq2klk87nrpff3y52u34zn4ydug2yd,me1hrxxjeqae2y5wx3kxcljzns9f2lguygu9qngxh,me14jazxhme3ptv00k52fza5rravx4xn27qs0slz2,me1qdhu5h5g0qwhdpl4q553v7gcmltdr4w3lnqnjg" "10000000000umec"
  med_v2 gentx_DAO --pubkey "$(med_v2 keys show "$KEY_NAME" -p --keyring-backend test)"

  med_v2 keys add "$KEY_NAME_SEQUENCER" --key-type secp256k1 --keyring-backend test

  jq '.app_state["dao"]["dao_addresses"]["global_dao"] = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq '.app_state["dao"]["dao_addresses"]["meid_dao"] = "me1p7s6k4ecrm2kl0rs6399k99pyuk322dc78dcxq"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq '.app_state["dao"]["dao_addresses"]["dev_operator"] = "me16qle3emp70kr08wt5508t7gk7trst0zwclnscj"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq '.app_state["dao"]["dao_addresses"]["airdrop_address"] = "me1uzt6kk6ra9x0ap3au3xuqwp94l2rnw4zqscn2s"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq '.app_state["rollapp"]["params"]["dispute_period_in_blocks"] = "1"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  jq --arg t "$UNBONDING_TIME" '.app_state["sequencer"]["params"]["unbonding_time"] = $t' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
  set_kyc_issuers

  med_v2 gentx "$KEY_NAME" "$STAKING_AMOUNT" --chain-id "$CHAIN_ID" --keyring-backend test --region-id me_earth --validator-address "$dao_addr"
  med_v2 collect-gentxs

  set_authorised_deployer_account "$dao_addr"
  set_authorised_deployer_account "$(med_v2 keys show "$KEY_NAME_SEQUENCER" -a --keyring-backend test)"

  med_v2 validate-genesis
  echo "v2 genesis ready at $HOME_DIR (mounted as $CONTAINER_HOME)"
  echo "start v2 with: $0 start-v2"
}

main() {
  local cmd="${1:-help}"
  if [ "$#" -gt 0 ]; then
    shift
  fi
  case "$cmd" in
    pull) cmd_pull ;;
    init) cmd_init ;;
    start-v2 | start) cmd_start_v2 ;;
    start-v3 | upgrade) cmd_start_v3 ;;
    stop) cmd_stop ;;
    logs) cmd_logs ;;
    status) cmd_status ;;
    exec) cmd_exec "$@" ;;
    propose) cmd_propose ;;
    wait-halt) cmd_wait_halt ;;
    help | -h | --help) usage ;;
    *)
      echo "error: unknown command: $cmd" >&2
      usage >&2
      exit 1
      ;;
  esac
}

main "$@"
