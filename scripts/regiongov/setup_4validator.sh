#!/bin/bash

# A backup suffix works with both BSD and GNU sed.
sed_inplace() {
    local file="${!#}"
    sed -i.bak "$@" || return
    rm -f "$file.bak"
}

#======================================================================================
# Hub-only setup: four Docker validators; requires med, jq, bc and Docker Compose.
# Existing hub data is backed up; other chains are not initialized or moved.
#config start
#======================================================================================
echo "Services must keep running in the background, so start this script in its own session: screen -S install"
echo "Press Enter after you confirm the settings"
read -p "Do you want to continue?"

read -rp "Please enter PROJECT_ENV[alpha]: " PROJECT_ENV
PROJECT_ENV=${PROJECT_ENV:-alpha}
echo "PROJECT_ENV is set to: $PROJECT_ENV"
BASE_DIR="$(pwd)"
echo "BASE_DIR is set to: $BASE_DIR"
read -rp "Please enter init_NODE_NAME[node1]: " NODE_NAME
NODE_NAME=${NODE_NAME:-node1}
echo "NODE_NAME is set to: $NODE_NAME"

read -rp "Please enter ME_CHAIN_ID[mechain_2401-1]: " ME_CHAIN_ID
ME_CHAIN_ID=${ME_CHAIN_ID:-"mechain_2401-1"}
echo "ME_CHAIN_ID is set to: $ME_CHAIN_ID"

read -rp "Please enter KEYRING_BACKEND_NAME[test]: " KEYRING_BACKEND_NAME
KEYRING_BACKEND_NAME=${KEYRING_BACKEND_NAME:-"test"}
echo "KEYRING_BACKEND_NAME is set to: $KEYRING_BACKEND_NAME"


read -rp "Please enter ENABLE_MONITORING[N]: " ENABLE_MONITORING
ENABLE_MONITORING=${ENABLE_MONITORING:-"N"}
echo "ENABLE_MONITORING is set to: $ENABLE_MONITORING"

read -rp "Please enter Docker_ADDR_RANGE[172.30.50.0/24]: " Docker_ADDR_RANGE
Docker_ADDR_RANGE=${Docker_ADDR_RANGE:-"172.30.50.0/24"}
echo "Docker_ADDR_RANGE is set to: $Docker_ADDR_RANGE"
DOCKER_NETWORK_PREFIX=$(echo $Docker_ADDR_RANGE | awk -F'[./]' '{print $1"."$2"."$3}')
echo "DOCKER_NETWORK_PREFIX is set to: $DOCKER_NETWORK_PREFIX"

read -rp "Please enter Docker_Network_Name[me-alpha]: " Docker_Network_Name
Docker_Network_Name=${Docker_Network_Name:-"me-alpha"}
echo "Docker_Network_Name is set to: $Docker_Network_Name"


read -rp "Whether to enable interactive mode?[N]: " Interactive_Mode
Interactive_Mode=${Interactive_Mode:-"N"}
echo "interactive mode is set to: $Interactive_Mode"

read -rp "Randomly generate mnemonic words?[N]: " Randomly_Generate_Mnemonic_Words
Randomly_Generate_Mnemonic_Words=${Randomly_Generate_Mnemonic_Words:-"N"}
echo "randomly generate mnemonic words is set to: $Randomly_Generate_Mnemonic_Words"

#Select test or production block times and governance periods.
read -rp "Quick install for test?[Y]: " INSTALL_TYPE
INSTALL_TYPE=${INSTALL_TYPE:-"Y"}
echo "Quick install for test is set to: $INSTALL_TYPE"



read -p "All Settings are initialized, Do you want to continue? (Y/n) "  REPLY
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Operation cancelled by user."
    exit 1
fi

cd "${BASE_DIR}" || { echo "error: BASE_DIR not found: ${BASE_DIR}"; exit 1; }
command -v med >/dev/null || { echo "error: med not found in PATH" >&2; exit 1; }
#======================================================================================
# mnemonic
if [[ ! $Randomly_Generate_Mnemonic_Words =~ ^[Yy]$ ]]
then
    global_dict="sign wise lyrics find brass grass already urge invite vibrant arrow grow dune theme trust mix napkin ball pause sponsor bacon capable flip people"
else
    global_dict=$(med keys mnemonic)
fi
#======================================================================================
# mnemonic index
USER_SEQ_ME_GLOBAL_DAO=1                                                # ME super-admin global_dao
USER_SEQ_ME_MEID_DAO=2                                                  # meid  address
USER_SEQ_ME_DEV_OPERATOR=3                                              # association address
USER_SEQ_ME_AIRDROP=4                                                   # airdrop address
USER_SEQ_ME_USER=5                                                      # regular user

USER_SEQ_ME_VAL1=11                                                     # validator node1
USER_SEQ_ME_VAL2=12                                                     # validator node2
USER_SEQ_ME_VAL3=13                                                     # validator node3
USER_SEQ_ME_VAL4=14                                                     # validator node4

#User addresses. Empty means derive from mnemonic + index. Fill in an external address for custody mode
USER_ADDR_ME_GLOBAL_DAO=                                                 # 1
USER_ADDR_ME_MEID_DAO=                                                   # 2
USER_ADDR_ME_DEV_OPERATOR=                                               # 3
USER_ADDR_ME_AIRDROP=                                                    # 4
USER_ADDR_ME_USER=                                                       # 5
USER_ADDR_ME_VAL1=                                                       # 11
USER_ADDR_ME_VAL2=                                                       # 12
USER_ADDR_ME_VAL3=                                                       # 13
USER_ADDR_ME_VAL4=                                                       # 14
MED_IMAGE="${MED_IMAGE:-ghcr.io/openmetaearth/med:v2.0.18}"
# Hub containers use this fixed home path.
HUB_CONTAINER_HOME="${HUB_CONTAINER_HOME:-/root/.mechain}"
HUB_MIN_GAS_PRICES="${HUB_MIN_GAS_PRICES:-0.02umec}"

# ME Hub directory
ME_NODES_HOME="${BASE_DIR}/nodes/hub-nodes"
# node directory
ME_NODE1_HOME="${ME_NODES_HOME}/${NODE_NAME}"

ME_NODE1_IP="${DOCKER_NETWORK_PREFIX}.10"
ME_NODE2_IP="${DOCKER_NETWORK_PREFIX}.11"
ME_NODE3_IP="${DOCKER_NETWORK_PREFIX}.12"
ME_NODE4_IP="${DOCKER_NETWORK_PREFIX}.13"

ME_TENDERMINT_CONFIG_FILE="${ME_NODE1_HOME}/config/config.toml"
ME_CLIENT_CONFIG_FILE="${ME_NODE1_HOME}/config/client.toml"
ME_APP_CONFIG_FILE="${ME_NODE1_HOME}/config/app.toml"
ME_GENESIS_FILE="${ME_NODE1_HOME}/config/genesis.json"


#Minimum deposit when creating a governance proposal. The denom is umec so amounts stay uniform (config:app_state.gov.deposit_params.min_deposit[0].denom )"failed to validate genesis state: invalid minimum deposit: 0umec"
ME_MIN_DEPOSIT_AMOUNT="100000000"
#ME genesis stake amount
#Modular supply is 10 billion MEC, i.e. 1000000000000000000umec.
TOTAL_SUPPLY="1000000000000000000umec"  # 100 billion MEC
calculate() {
    echo "$1" | bc
}
# ME-HUB
#======================================================================================

TOTAL_AMOUNT=${TOTAL_SUPPLY%umec}

#experience_region:ME_VAL2 -experience, share of stake_tokens_pool 10%
ME_STAKING_AMOUNT_NODE2="$(calculate "${TOTAL_AMOUNT} * 10 / 100")umec"   # 10% = 100,000,000 MEC = 100000000000000000umec
ME_STAKING_AMOUNT_NODE2_REGION="experience_region"
#usa:ME_VAL3 -ind, share of stake_tokens_pool 17.76%
ME_STAKING_AMOUNT_NODE3="$(calculate "${TOTAL_AMOUNT} * 1776 / 10000")umec"  # 17.76% = 177,600,000 MEC = 177600000000000000umec
ME_STAKING_AMOUNT_NODE3_REGION="ind"
#rus:ME_VAL4 -chn, share of stake_tokens_pool 17.72%
ME_STAKING_AMOUNT_NODE4="$(calculate "${TOTAL_AMOUNT} * 1772 / 10000")umec"   # 17.72% = 177,200,000 MEC = 177200000000000000umec
ME_STAKING_AMOUNT_NODE4_REGION="chn"

#me_earth:ME_VAL1. After the three regions above, the rest of the 10 billion is staked to me_earth  
STAKING_AMOUNT="$(calculate "${TOTAL_AMOUNT} - ${ME_STAKING_AMOUNT_NODE2%umec} - ${ME_STAKING_AMOUNT_NODE3%umec} - ${ME_STAKING_AMOUNT_NODE4%umec}")umec"  # 54.52% = 545,200,000 MEC =545200000000000000umec
ME_STAKING_AMOUNT_NODE1_REGION="me_earth"
# staking module param: tokens must stay locked for 604800s/7 days, 18000s/5 hours, 3024000s/28 days, 21600s/6 hours.
# sequencer module param For Cosmos SDK-based chains, trust_period should usually be about 2/3 of the unbonding time (~2
# weeks) during which they can be financially punished (slashed) for misbehavior.trust_period is configured in the config.toml.
ME_UNBONDING_TIME="3024000s"

if [[ $INSTALL_TYPE =~ ^[Yy]$ ]]
then
    #Dispute period. A block is final only after the dispute window
    ME_DISPUTE_PERIOD_IN_BLOCKS="5"
    # Incentive params. Production should use "day" and test should use "minute" to control release frequency. The code is hard-coded to day, which releases once every 17280 blocks. Creating a region, a region_share change, or a height that is a multiple of 17280 all distribute block rewards.
    DISTR_EPOCH_IDENTIFIER="day"
    # Incentive params. Production should use 2 weeks "14d"; test uses "60s" to control the lock time
    LOCKABLE_DURATIONS="60s"
    # max_deposit_period: test uses 600s, production uses 2 days
    MAX_DEPOSIT_PERIOD="600s"
    # voting_period: test uses 3 minutes; expedited must be strictly shorter than voting_period
    VOTING_PERIOD="180s"
    EXPEDITED_VOTING_PERIOD="90s"
    # Block time. Test uses 500ms; production currently uses 5s. 1330 blocks is about 11 minutes
    ME_TIMEOUT_COMMIT="5s"
else
    ME_DISPUTE_PERIOD_IN_BLOCKS="50"
    DISTR_EPOCH_IDENTIFIER="day"
    LOCKABLE_DURATIONS="14d"
    MAX_DEPOSIT_PERIOD="172800s"
    VOTING_PERIOD="172800s"
    EXPEDITED_VOTING_PERIOD="300s"
    ME_TIMEOUT_COMMIT="5s"
fi
add_key() {
    local action=$1
    local key_name=$2
    local user_seq=$3
    local extra_args=$4
    local mnemonics=${global_dict}

    if [ "${action}" == "med" ]; then
        if ! output=$(echo -e "y\n${mnemonics}\n" | ${action} keys add "$key_name" \
                --keyring-backend ${KEYRING_BACKEND_NAME} \
                --hd-path "m/44'/118'/0'/0/0/${user_seq}" \
                --key-type secp256k1 \
                --recover $extra_args --coin-type 118 2>&1); then

                # If the above fails, fall through to the no-prompt y confirmation
                echo -e "${mnemonics}\n" | ${action} keys add "$key_name" \
                        --keyring-backend ${KEYRING_BACKEND_NAME} \
                        --key-type secp256k1 \
                        --hd-path "m/44'/118'/0'/0/0/${user_seq}" \
                        --recover $extra_args --coin-type 118
        else
                echo "$output"
        fi
    else
        if ! output=$(echo -e "y\n${mnemonics}\n" | ${action} keys add "$key_name" \
                --hd-path "m/44'/118'/0'/0/0/${user_seq}" \
                --recover $extra_args --coin-type 118 2>&1); then

                # If the above fails, fall through to the no-prompt y confirmation
                echo -e "${mnemonics}\n" | ${action} keys add "$key_name" \
                        --keyring-backend ${KEYRING_BACKEND_NAME} \
                        --hd-path "m/44'/118'/0'/0/0/${user_seq}" \
                        --recover $extra_args --coin-type 118
        else
                echo "$output"
        fi
    fi
}


# Check whether the tx hash is on chain. Wait at most 60s
check_tx_status() {
    local cosmosapp=$1
    local txhash=$2
    local extra_args=$3
    sleepTime=120

    # Check that the argument is non-empty
    if [ -z "$txhash" ]; then
        echo "Error: cosmosapp or txhash cannot be empty."
        echo "Usage: check_cosmos_tx_status <txhash>"
        exit 1
        # return 1
    fi

    echo "Check tx status: ${txhash}"
    while true; do
        # Query the transaction and capture the exit code
        if ! $cosmosapp query tx ${txhash} $extra_args > /dev/null 2>&1; then
            sleep 1
            sleepTime=$((sleepTime - 1))
            if [ $sleepTime -eq 0 ]; then  
                echo "shell:$cosmosapp query tx ${txhash} $extra_args -o json"
                echo "Tx failed, not found in chain, txhash: ${txhash}"
                echo "Current result:"
                $cosmosapp query tx ${txhash} $extra_args
                exit 1
            fi
            continue
        fi

        # If the transaction is found, check the response code
        resp_code=$($cosmosapp q tx ${txhash} -o json $extra_args | jq -r '.code' 2>/dev/null)

        if [ "$resp_code" != "0" ]; then 
            echo "shell:$cosmosapp query tx ${txhash} $extra_args -o json"
            echo "Tx failed in chain, code: ${resp_code}"
            exit 1
        fi

        # If the code is 0, the transaction is successful
        echo "Tx succeeded, code: ${resp_code}"
        break
    done
}


echo "# ---------------------------------------------------------------------------- #"
echo "#                       Create the docker compose network                                  #"
echo "# ---------------------------------------------------------------------------- #"

if docker network inspect "$Docker_Network_Name" >/dev/null 2>&1; then
    echo "Error: Docker network '$Docker_Network_Name' already exists." >&2
    echo "Confirm deletion of network $Docker_Network_Name!!!!!!!!!!? (Y/n) "
    read -r del_network
    del_network=${del_network:-N}
    if [ ! "$del_network" != "${del_network#[Yy]}" ] ;then
        echo "You chose not to delete the network. Enter a new network name"
        echo "Existing networks:"
        docker network ls --format '{{.Name}}'
        while true; do
            read -rp "Please enter a new unique Docker_Network_Name: " Docker_Network_Name
            # Check whether the network name exists
            if ! docker network inspect "${Docker_Network_Name}" >/dev/null 2>&1; then
                break
            fi
            echo "Error: Network ${Docker_Network_Name} already exists. Please choose a different name."
        done
        echo "Docker_Network_Name is set to: ${Docker_Network_Name}"
    else
        echo "Deleting network..."
        docker network rm ${Docker_Network_Name}
        if [ $? -ne 0 ]; then
            echo "Failed to delete the network, exiting"
            exit 1
        fi
    fi
fi

verify_return_code() {
    # Check the return value
    if [ $? -eq 0 ]; then
        echo "Network created"
    else
        echo "That subnet may already exist. Enter a different subnet"
        exit 1
    fi
}
docker network create ${Docker_Network_Name} --subnet ${Docker_ADDR_RANGE} --ip-range ${Docker_ADDR_RANGE}
# Exit if creating this network fails
verify_return_code


if [ -d "${ME_NODES_HOME}" ]; then
    mv "${ME_NODES_HOME}" "${ME_NODES_HOME}.old-$(date "+%Y-%m-%d-%H-%M-%S")"
fi
mkdir -p "${ME_NODES_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                            init ME node1                                   #"
echo "# ---------------------------------------------------------------------------- #"
med init ${NODE_NAME} --chain-id="${ME_CHAIN_ID}" --home ${ME_NODE1_HOME}


echo "# ---------------------------------------------------------------------------- #"
echo "#                              Set configurations                               #"
echo "# ---------------------------------------------------------------------------- #"
echo "Edit/config/config.toml"
# Mempool version to use:
#   1) "v0" - (default) FIFO mempool.
#   2) "v1" - prioritized mempool (deprecated; will be removed in the next release).
#sed_inplace -e "/\[mempool\]/,/^\[/ s/v0/v1/" "$TENDERMINT_CONFIG_FILE"
#Allow multiple services on one machine on different ports
sed_inplace '/^allow_duplicate_ip/s/false/true/' "${ME_NODE1_HOME}"/config/config.toml
#Address-book strictness. false allows non-routable addresses (such as private IPs) as peers and non-global unicast addresses.
sed_inplace '/^addr_book_strict/s/true/false/' "${ME_NODE1_HOME}"/config/config.toml
sed_inplace 's/0.0.0.0:26656/0.0.0.0:26656/' "${ME_NODE1_HOME}"/config/config.toml
sed_inplace 's/127.0.0.1:26657/0.0.0.0:26657/' "${ME_NODE1_HOME}"/config/config.toml
# Block time
sed_inplace 's/^timeout_commit = .*/timeout_commit = "'${ME_TIMEOUT_COMMIT}'"/' "${ME_NODE1_HOME}"/config/config.toml
#CORS
sed_inplace 's/cors_allowed_origins.*$/cors_allowed_origins = ["*"]/' "${ME_NODE1_HOME}"/config/config.toml
sed_inplace 's/localhost:6060/0.0.0.0:6060/' "${ME_NODE1_HOME}"/config/config.toml
echo "Edit/config/app.toml"
sed_inplace -e '/\[api\]/,/^\[/ s/enable *= .*/enable = true/' "${ME_NODE1_HOME}/config/app.toml"
sed_inplace 's/swagger = false/swagger = true/'   "${ME_NODE1_HOME}/config/app.toml"
sed_inplace -e "/\[api\]/,/^\[/ s/address *= .*/address = \"tcp:\/\/0.0.0.0:1317\"/" "${ME_NODE1_HOME}/config/app.toml"
sed_inplace 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "${ME_NODE1_HOME}/config/app.toml"

sed_inplace -e "/\[grpc\]/,/^\[/ s/address *= .*/address = \"0.0.0.0:9090\"/"  "${ME_NODE1_HOME}/config/app.toml"

sed_inplace -e "/\[grpc-web\]/,/^\[/ s/address *= .*/address = \"0.0.0.0:9091\"/" "${ME_NODE1_HOME}/config/app.toml"
sed_inplace 's/enable-unsafe-cors = false/enable-unsafe-cors = true/' "${ME_NODE1_HOME}/config/app.toml"

sed_inplace -e "/\[json-rpc\]/,/^\[/ s/address *= .*/address = \"0.0.0.0:8545\"/" "${ME_NODE1_HOME}/config/app.toml"
sed_inplace -e "/\[json-rpc\]/,/^\[/ s/^ws-address *= .*/ws-address = \"0.0.0.0:8546\"/" "${ME_NODE1_HOME}/config/app.toml"
sed_inplace 's/127.0.0.1:6065/0.0.0.0:6065/' "${ME_NODE1_HOME}/config/app.toml"

#Minimum gas price a node accepts. Fee = gas price times gas used. Governance can change it; validators may set different values.
sed_inplace -e 's/^minimum-gas-prices *= .*/minimum-gas-prices = "'"${HUB_MIN_GAS_PRICES}"'"/' "${ME_NODE1_HOME}/config/app.toml"
# default: the last 362880 states are kept, pruning at 10 block intervals
# nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
# everything: 2 latest states will be kept; pruning at 10 block intervals.
# custom: allow pruning options to be manually specified through 'pruning-keep-recent', and 'pruning-interval'
sed_inplace -e 's/^pruning *= .*/pruning = "nothing"/' "${ME_NODE1_HOME}/config/app.toml"

sed_inplace -e 's/^max-recv-msg-size *= .*/max-recv-msg-size = "1048576000"/' "${ME_NODE1_HOME}/config/app.toml"

echo "Edit/config/client.toml"
sed_inplace "s/localhost:26657/127.0.0.1:26657/" "${ME_NODE1_HOME}"/config/client.toml
# The keyring's backend, where the keys are stored (os|file|kwallet|pass|test|memory)
echo "Keyring backend is now:${KEYRING_BACKEND_NAME}"
sed_inplace -e "s/^keyring-backend *= .*/keyring-backend = \"${KEYRING_BACKEND_NAME}\"/" "${ME_NODE1_HOME}"/config/client.toml

med config chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME}
med config keyring-backend ${KEYRING_BACKEND_NAME} --home ${ME_NODE1_HOME}
med config node "tcp://127.0.0.1:26657" --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#                              set Genesis config                               #"
echo "# ---------------------------------------------------------------------------- #"
tmp=$(mktemp)
set_gov_params() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                              setting gov params                               #"
    echo "# ---------------------------------------------------------------------------- #"
    jq '.app_state.gov.deposit_params.min_deposit[0].denom = "umec"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg amount "$ME_MIN_DEPOSIT_AMOUNT" '.app_state.gov.deposit_params.min_deposit[0].amount = $amount' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg period "${VOTING_PERIOD}" '.app_state.gov.voting_params.voting_period = $period' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.deposit_params.max_deposit_period = $period' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq '.app_state.gov.params.min_deposit[0].denom = "umec"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg amount "$ME_MIN_DEPOSIT_AMOUNT" '.app_state.gov.params.min_deposit[0].amount = $amount' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    #flase means do not burn. Controls the veto deposit: true burns veto-vote tokens (removed from circulation forever); false refunds them to the voter
    jq '.app_state.gov.params.burn_vote_veto = false' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg period "${VOTING_PERIOD}" '.app_state.gov.params.voting_period = $period' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.params.max_deposit_period = $period' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    # expedited_voting_period must be shorter than voting_period, or genesis validate fails
    jq --arg period "${EXPEDITED_VOTING_PERIOD}" 'if .app_state.gov.params.expedited_voting_period then .app_state.gov.params.expedited_voting_period = $period else . end' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}

set_hub_params() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                              setting hub params                               #"
    echo "# ---------------------------------------------------------------------------- #"
    sed_inplace -e 's/bond_denom": ".*"/bond_denom": "umec"/' "${ME_GENESIS_FILE}"
    sed_inplace -e 's/mint_denom": ".*"/mint_denom": "umec"/' "${ME_GENESIS_FILE}"

    jq ".app_state.rollapp.params.dispute_period_in_blocks = \"$ME_DISPUTE_PERIOD_IN_BLOCKS\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    #bridging_fee is a bridge param that sets the cross-chain fee, usually paid in the native denom
    jq ".app_state.delayedack.params.bridging_fee = \"0\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    #epoch_identifier is a bridge param that sets the cross-chain epoch, for example week, month, or year
    jq ".app_state.delayedack.params.epoch_identifier = \"week\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq ".app_state.eibc.params.epoch_identifier = \"week\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    #IBC packets removed per epoch, set to 1000000;
    jq ".app_state.delayedack.params.delete_packets_epoch_limit = \"1000000\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"


    #increase the tx size cost per byte from 10 to 100
    jq ".app_state.auth.params.tx_size_cost_per_byte = \"100\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"


    # jail validators faster, and shorten recovery time, no slash for downtime
    jq ".app_state.slashing.params.signed_blocks_window = \"10000\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq ".app_state.slashing.params.min_signed_per_window = \"0.800000000000000000\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq ".app_state.slashing.params.downtime_jail_duration = \"120s\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq ".app_state.slashing.params.slash_fraction_downtime = \"0.0\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}

set_consenus_params() {
    # cometbft's updated values
	  # 	MaxBytes: 4194304,  // four megabytes
	  # 	MaxGas:   10000000, // ten million
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                        setting consensus params                           #"
    echo "# ---------------------------------------------------------------------------- #"
    jq '.consensus_params["block"]["max_bytes"] = "4194304"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq '.consensus_params["block"]["max_gas"] = "-1"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}

set_EVM_params() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                        setting EVM params                                    #"
    echo "# ---------------------------------------------------------------------------- #"
    jq '.app_state["feemarket"]["params"]["no_base_fee"] = false' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq '.app_state.evm.params.evm_denom = "umec"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq '.app_state.evm.params.enable_create = true' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}

#Adding a "minute" epoch
set_epochs_params() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                        setting epochs params                              #"
    echo "# ---------------------------------------------------------------------------- #"
    jq '.app_state.epochs.epochs += [{
    "identifier": "minute",
    "start_time": "0001-01-01T00:00:00Z",
    "duration": "60s",
    "current_epoch": "0",
    "current_epoch_start_time": "0001-01-01T00:00:00Z",
    "epoch_counting_started": false,
    "current_epoch_start_height": "0"
    }]' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}

#should be set to days on live net and lockable duration to 2 weeks
set_incentives_params() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                        setting incentives params(incentives)                        #"
    echo "# ---------------------------------------------------------------------------- #"
    jq --arg time "$DISTR_EPOCH_IDENTIFIER" '.app_state.incentives.params.distr_epoch_identifier = $time' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg time "$LOCKABLE_DURATIONS" '.app_state.incentives.lockable_durations = ["$time"]' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"

}


set_misc_params() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  setting misc params(crisis-tx constant fee denom)         #"
    echo "# ---------------------------------------------------------------------------- #"
    jq '.app_state.crisis.constant_fee.denom = "umec"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq -r '.app_state.gamm.params.pool_creation_fee[0].denom = "umec"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq '.app_state["txfees"]["basedenom"] = "umec"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq '.app_state["txfees"]["params"]["epoch_identifier"] = "minute"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"

    jq -r '.app_state.gamm.params.enable_global_pool_fees = true' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}

set_bank_denom_metadata() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                       setting bank denom params                           #"
    echo "# ---------------------------------------------------------------------------- #"
    jq '.app_state.bank.denom_metadata = [
        {
            "base": "umec",
            "denom_units": [
                {
                    "aliases": [],
                    "denom": "umec",
                    "exponent": 0
                },
                {
                    "aliases": [],
                    "denom": "MEC",
                    "exponent": 8
                }
            ],
            "description": "Denom metadata for MEC (umec)",
            "display": "MEC",
            "name": "MEC",
            "symbol": "MEC"
        }
    ]' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}

set_consenus_params
set_gov_params
set_hub_params
set_misc_params
set_EVM_params
set_bank_denom_metadata
set_epochs_params
set_incentives_params


if [[ "$ENABLE_MONITORING" == [Yy] ]] ;then
  sed_inplace -e "s/prometheus = false/prometheus = true/" "${ME_TENDERMINT_CONFIG_FILE}"
fi

echo "# ---------------------------------------------------------------------------- #"
echo "# Decide whether (global_dao, meid_dao, dev_operator, airdrop, val1, user) need to be recovered and add #"
echo "# -------------------------------------------------------------------------------------#"
if [ -z "${USER_ADDR_ME_GLOBAL_DAO}" ]; then
    echo "add global_dao"
    me_global_dao=$(add_key "med" "global_dao"  ${USER_SEQ_ME_GLOBAL_DAO} "--home $ME_NODE1_HOME")
    echo "${me_global_dao}"
    USER_ADDR_ME_GLOBAL_DAO=$(echo "${me_global_dao}" | awk '/^[[:space:]-]*address:/ {print $NF}')
    USER_ADDR_ME_GLOBAL_DAO_PUBKEY=$(echo "${me_global_dao}" | sed -n 's/^[[:space:]]*pubkey: //p' | sed "s/'//g")
    echo "-------------------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_MEID_DAO}" ]; then
    echo "add meid_dao"
    me_meid_dao=$(add_key "med"  "meid_dao" ${USER_SEQ_ME_MEID_DAO} "--home $ME_NODE1_HOME")
    echo "${me_meid_dao}"
    USER_ADDR_ME_MEID_DAO=$(echo "${me_meid_dao}" | awk '/^[[:space:]-]*address:/ {print $NF}')
    USER_ADDR_ME_MEID_DAO_PUBKEY=$(echo "${me_meid_dao}" | sed -n 's/^[[:space:]]*pubkey: //p' | sed "s/'//g")
    echo "------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_DEV_OPERATOR}" ]; then
    echo "add dev_operator"
    me_dev_operator=$(add_key  "med" "dev_operator" ${USER_SEQ_ME_DEV_OPERATOR} "--home $ME_NODE1_HOME")
    echo "${me_dev_operator}"
    USER_ADDR_ME_DEV_OPERATOR=$(echo "${me_dev_operator}" | awk '/^[[:space:]-]*address:/ {print $NF}')
    USER_ADDR_ME_DEV_OPERATOR_PUBKEY=$(echo "${me_dev_operator}" | sed -n 's/^[[:space:]]*pubkey: //p' | sed "s/'//g")
    echo "------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_AIRDROP}" ]; then
    echo "add airdrop"
    me_airdrop=$(add_key  "med" "airdrop" ${USER_SEQ_ME_AIRDROP} "--home $ME_NODE1_HOME")
    echo "${me_airdrop}"
    USER_ADDR_ME_AIRDROP=$(echo "${me_airdrop}" | awk '/^[[:space:]-]*address:/ {print $NF}')
    USER_ADDR_ME_AIRDROP_PUBKEY=$(echo "${me_airdrop}" | sed -n 's/^[[:space:]]*pubkey: //p' | sed "s/'//g")
    echo "-------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_VAL1}" ]; then
    echo "add val1"
    me_val1=$(add_key  "med"  "val1" ${USER_SEQ_ME_VAL1} "--home $ME_NODE1_HOME")
    echo "${me_val1}"
    USER_ADDR_ME_VAL1=$(echo "${me_val1}" | awk '/^[[:space:]-]*address:/ {print $NF}')
    USER_ADDR_ME_VAL1_PUBKEY=$(echo "${me_val1}" | sed -n 's/^[[:space:]]*pubkey: //p' | sed "s/'//g")
    echo "--------------------------------------------------------------------------------------"
fi
echo "-----------------------------------------------------------------------------------------"
echo "# ---------------------------------------------------------------------------------------#"
echo "# Initialize global_dao with an extra 100,000,000 MEC (10^16 umec) #"
echo "# (meid_dao, dev_operator, airdrop,val1)-"0umec"                                        #"                            
echo "# --------------------------------------------------------------------------------------#"
med add-genesis-account global_dao   "10000000000000000umec" --home ${ME_NODE1_HOME}
med add-genesis-account meid_dao     "0umec" --home ${ME_NODE1_HOME}
med add-genesis-account dev_operator "0umec" --home ${ME_NODE1_HOME}
med add-genesis-account airdrop      "0umec" --home ${ME_NODE1_HOME}
med add-genesis-account val1         "0umec" --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#                              Set the stake pool                                       #"
echo "# ---------------------------------------------------------------------------- #"

med add-genesis-stake-pool  --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#                         Genesis module accounts: internal module addresses                           #"
echo "# ---------------------------------------------------------------------------- #"

med add-genesis-m-accounts --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#   Set dao_addresses(global_dao, meid_dao, dev_operator,airdrop_address)      #"
echo "# ---------------------------------------------------------------------------- #"
jq --arg addr "${USER_ADDR_ME_GLOBAL_DAO}" \
  '.app_state.dao.dao_addresses.global_dao = $addr' \
  "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"

jq --arg addr "${USER_ADDR_ME_MEID_DAO}" \
  '.app_state.dao.dao_addresses.meid_dao = $addr' \
  "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"

jq --arg addr "${USER_ADDR_ME_DEV_OPERATOR}" \
  '.app_state.dao.dao_addresses.dev_operator = $addr' \
  "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"

jq --arg addr "${USER_ADDR_ME_AIRDROP}" \
  '.app_state.dao.dao_addresses.airdrop_address = $addr' \
  "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"

jq --arg unbonding_time "${ME_UNBONDING_TIME}" '.app_state["sequencer"]["params"]["unbonding_time"] = $unbonding_time'  "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
jq --arg unbonding_time "${ME_UNBONDING_TIME}" '.app_state["staking"]["params"]["unbonding_time"] = $unbonding_time'  "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
set_kyc_issuers() {
    echo "# ----------------------------------------------------------------------------------------------- #"
    echo "#  setting issuer:USER_ADDR_ME_GLOBAL_DAO, USER_ADDR_ME_MEID_DAO (addresses allowed to issue or verify KYC) #"
    echo "# ----------------------------------------------------------------------------------------------- #"
    jq --arg addr1 "${USER_ADDR_ME_GLOBAL_DAO}" \
       --arg pubkey1 "${USER_ADDR_ME_GLOBAL_DAO_PUBKEY}" \
       --arg addr2 "${USER_ADDR_ME_MEID_DAO}" \
       --arg pubkey2 "${USER_ADDR_ME_MEID_DAO_PUBKEY}" \
       '.app_state.kyc.issuers = [
          {
            "did": "0000000000001",
            "address": $addr1,
            "pubkey": $pubkey1,
            "kycLevel": "KYC_LEVEL_TWO",
            "status": "DID_STATUS_ACTIVE"
          },
          {
            "did": "0000000000002",
            "address": $addr2,
            "pubkey": $pubkey2,
            "kycLevel": "KYC_LEVEL_TWO",
            "status": "DID_STATUS_ACTIVE"
          }
        ]' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}


echo "# ---------------------------------------------------------------------------- #"
echo "#                Genesis stake (production uses the built-in stake_tokens_pool account for the stake)       #"
echo "# ---------------------------------------------------------------------------- #"

#global_dao # account name that sends the tx
#STAKING_AMOUNT # stake amount
#--chain-id chain ID
#--keyring-backend # keyring backend (test or file)
#--region-id. # validator region id (earth region)
#--validator-address  # validator operator address


# Check module account
jq '.app_state.auth.accounts[] | select(.name=="stake_tokens_pool")' $ME_NODE1_HOME/config/genesis.json

# Check balance
jq '.app_state.bank.balances[] | select(.address | contains("me1gmpxkchcdgfq995zye5efwzfw86zfa4vcxms9l"))' $ME_NODE1_HOME/config/genesis.json
# Check module-account balance
#med q bank balances $(med q auth module-account stake_tokens_pool -o json --home $ME_NODE1_HOME | jq -r '.base_account.address') --home $ME_NODE1_HOME


#stake
echo "med gentx global_dao "${STAKING_AMOUNT}" --chain-id "${ME_CHAIN_ID}" --keyring-backend ${KEYRING_BACKEND_NAME} --region-id ${ME_STAKING_AMOUNT_NODE1_REGION} --validator-address "${USER_ADDR_ME_VAL1}" --home "${ME_NODE1_HOME}""
med gentx global_dao "${STAKING_AMOUNT}" --chain-id "${ME_CHAIN_ID}" --keyring-backend ${KEYRING_BACKEND_NAME} --region-id ${ME_STAKING_AMOUNT_NODE1_REGION}  --validator-address "${USER_ADDR_ME_VAL1}" --home "${ME_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                              Collect genesis txs                                     #"
echo "# ---------------------------------------------------------------------------- #"

med collect-gentxs --home "${ME_NODE1_HOME}"


echo "# ---------------------------------------------------------------------------- #"
echo "#                Configure kyc issuers                                               #"
echo "# ---------------------------------------------------------------------------- #"

set_kyc_issuers

echo "# ---------------------------------------------------------------------------- #"
echo "#                              Validate genesis txs                                     #"
echo "# ---------------------------------------------------------------------------- #"

med validate-genesis --home "${ME_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                              ME sync node config                                   #"
echo "# ---------------------------------------------------------------------------- #"
me_init_sync_node() {
    local NODE_NAME=$1
    local NODE_HOME="${ME_NODES_HOME}/${NODE_NAME}"
    local NODE_IP=$2
    local CONFIG_DIRECTORY="${NODE_HOME}/config"
    local TENDERMINT_CONF="${CONFIG_DIRECTORY}/config.toml"
    local APP_CONFIG_FILE="${CONFIG_DIRECTORY}/app.toml"
    local CLIENT_CONFIG_FILE="${CONFIG_DIRECTORY}/client.toml"

    echo "========================================"
    echo "init node: ${NODE_NAME} - ${NODE_HOME}"

    rm -rf "${NODE_HOME}"

    med init "$NODE_NAME" --chain-id="$ME_CHAIN_ID" --home "${NODE_HOME}"

    /bin/cp "${ME_GENESIS_FILE}"  "${CONFIG_DIRECTORY}"
    /bin/cp -rf ${ME_NODE1_HOME}/config/gentx  ${CONFIG_DIRECTORY}/
    /bin/cp -rf ${ME_NODE1_HOME}/config/app.toml  ${CONFIG_DIRECTORY}/
    /bin/cp -rf ${ME_NODE1_HOME}/config/client.toml  ${CONFIG_DIRECTORY}/
    /bin/cp -rf ${ME_NODE1_HOME}/config/config.toml ${CONFIG_DIRECTORY}/
    sed_inplace "s|^moniker = .*$|moniker = \"${NODE_NAME}\"|" "${TENDERMINT_CONF}"
    sed_inplace "s|^node = .*$|node = \"tcp://${NODE_IP}:26657\"|" "${CLIENT_CONFIG_FILE}"
    
}

me_init_sync_node node2 ${ME_NODE2_IP}
me_init_sync_node node3 ${ME_NODE3_IP}
me_init_sync_node node4 ${ME_NODE4_IP}

echo "# ---------------------------------------------------------------------------- #"
echo "#                              ME Read Node ID                                     #"
echo "# ---------------------------------------------------------------------------- #"
NODEID_NODE1=$(med tendermint show-node-id --home ${ME_NODE1_HOME})
NODEID_NODE2=$(med tendermint show-node-id --home ${ME_NODES_HOME}/node2)
NODEID_NODE3=$(med tendermint show-node-id --home ${ME_NODES_HOME}/node3)
NODEID_NODE4=$(med tendermint show-node-id --home ${ME_NODES_HOME}/node4)
echo Node1_ID:${NODEID_NODE1}
echo Node2_ID:${NODEID_NODE2}
echo Node3_ID:${NODEID_NODE3}
echo Node4_ID:${NODEID_NODE4}

echo "# ---------------------------------------------------------------------------- #"
echo "#                         Edit Node2-Node5 persistent_peers                #"
echo "# ---------------------------------------------------------------------------- #"
NodeID_Node1_P2P="${NODEID_NODE2}@${ME_NODE2_IP}:26656,${NODEID_NODE3}@${ME_NODE3_IP}:26656,${NODEID_NODE4}@${ME_NODE4_IP}:26656"

echo NodeID_Node1_P2P:$NodeID_Node1_P2P
NodeID_Node2_P2P="${NODEID_NODE1}@${ME_NODE1_IP}:26656,${NODEID_NODE3}@${ME_NODE3_IP}:26656,${NODEID_NODE4}@${ME_NODE4_IP}:26656"

echo NodeID_Node2_P2P:$NodeID_Node2_P2P
NodeID_Node3_P2P="${NODEID_NODE1}@${ME_NODE1_IP}:26656,${NODEID_NODE2}@${ME_NODE2_IP}:26656,${NODEID_NODE4}@${ME_NODE4_IP}:26656"

echo NodeID_Node3_P2P:${NodeID_Node3_P2P}
NodeID_Node4_P2P="${NODEID_NODE1}@${ME_NODE1_IP}:26656,${NODEID_NODE2}@${ME_NODE2_IP}:26656,${NODEID_NODE3}@${ME_NODE3_IP}:26656"

echo NodeID_Node4_P2P:${NodeID_Node4_P2P}
sed_inplace 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node2_P2P}"'"/' ${ME_NODES_HOME}/node2/config/config.toml
sed_inplace 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node3_P2P}"'"/' ${ME_NODES_HOME}/node3/config/config.toml
sed_inplace 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node4_P2P}"'"/' ${ME_NODES_HOME}/node4/config/config.toml

echo "# ---------------------------------------------------------------------------- #"
echo "#                          Generate docker-compose file                            #"
echo "# ---------------------------------------------------------------------------- #"
cat >docker-compose.yml<<EOF
x-me-template: &me-template
  restart: unless-stopped
  image: ${MED_IMAGE}
  entrypoint: ["med"]
  networks:
    - ${Docker_Network_Name}
  command:
    - start
    - --home
    - ${HUB_CONTAINER_HOME}
    - --minimum-gas-prices
    - "${HUB_MIN_GAS_PRICES}"

networks:
  ${Docker_Network_Name}:
    name: ${Docker_Network_Name} 
    external: true

services:
  hub-node1:
    <<: *me-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${ME_NODE1_IP}
    ports:
      - 26657:26657
      - 9090:9090
      - 1317:1317
      - 8545:8545
    volumes:
      - ${ME_NODE1_HOME}:${HUB_CONTAINER_HOME}

  hub-node2:
    <<: *me-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${ME_NODE2_IP}
    volumes:
      - ${ME_NODES_HOME}/node2:${HUB_CONTAINER_HOME}

  hub-node3:
    <<: *me-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${ME_NODE3_IP}
    volumes:
      - ${ME_NODES_HOME}/node3:${HUB_CONTAINER_HOME}

  hub-node4:
    <<: *me-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${ME_NODE4_IP}
    volumes:
      - ${ME_NODES_HOME}/node4:${HUB_CONTAINER_HOME}
EOF
docker compose up -d hub-node1
echo "Wait 5s for services to finish starting"
sleep 5
docker compose up -d hub-node2
echo "Wait 5s for services to finish starting"
sleep 5
docker compose up -d hub-node3
echo "Wait 5s for services to finish starting"
sleep 5
docker compose up -d hub-node4
echo "Wait 5s for services to finish starting"
sleep 5

echo "# ---------------------------------------------------------------------------- #"
echo "#                  Edit node1 persistent_peers                              #"
echo "# --------------------- ------------------------------------------------------- #"
docker compose stop hub-node1
sed_inplace 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node1_P2P}"'"/' ${ME_NODE1_HOME}/config/config.toml
echo NodeID_Node1_P2P:${NodeID_Node1_P2P}
docker compose  up -d hub-node1
echo "Wait 5s for services to finish starting"
sleep 5
echo "# ---------------------------------------------------------------------------- #"
echo "#                  Current container status                                              #"
echo "# ---------------------------------------------------------------------------- #"
docker compose  ps

echo "# ---------------------------------------------------------------------------- #"
echo "#               Decide whether val2, val3, val4 need to be add                                  #"
echo "# ---------------------------------------------------------------------------- #"
if [ -z "${USER_ADDR_ME_VAL2}" ]; then
    val2=$(add_key  "med"  "val2" ${USER_SEQ_ME_VAL2}  "--home $ME_NODE1_HOME")
    echo "${val2}"
    USER_ADDR_ME_VAL2=$(echo "${val2}" | awk '/^[[:space:]-]*address:/ {print $NF}')
    USER_ADDR_ME_VAL2_PUBKEY=$(echo "${val2}" | sed -n 's/^[[:space:]]*pubkey: //p' | sed "s/'//g")
    echo "--------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_VAL3}" ]; then
    val3=$(add_key  "med"  "val3" ${USER_SEQ_ME_VAL3} "--home $ME_NODE1_HOME")
    echo "${val3}"
    USER_ADDR_ME_VAL3=$(echo "${val3}" | awk '/^[[:space:]-]*address:/ {print $NF}')
    USER_ADDR_ME_VAL3_PUBKEY=$(echo "${val3}" | sed -n 's/^[[:space:]]*pubkey: //p' | sed "s/'//g")
    echo "--------------------------------------------------------------------------"
fi    
if [ -z "${USER_ADDR_ME_VAL4}" ]; then
    val4=$(add_key  "med"  "val4" ${USER_SEQ_ME_VAL4} "--home $ME_NODE1_HOME")
    echo "${val4}"
    USER_ADDR_ME_VAL4=$(echo "${val4}" | awk '/^[[:space:]-]*address:/ {print $NF}')
    USER_ADDR_ME_VAL4_PUBKEY=$(echo "${val4}" | sed -n 's/^[[:space:]]*pubkey: //p' | sed "s/'//g")
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#                  Wait for the ME genesis block and get the block hash                              #"
echo "# ---------------------------------------------------------------------------- #"
wait_genesis_block_hash() {
    local attempt block_hash
    for ((attempt=1; attempt<=30; attempt++)); do
        if block_hash=$(med q block 1 --home "${ME_NODE1_HOME}" --node tcp://127.0.0.1:26657 |
            jq -er '.block_id.hash | select(type == "string" and length > 0)'); then
            echo "Genesis block hash: $block_hash"
            return
        fi
        echo "Waiting for Hub RPC ($attempt/30), retrying in 2s..." >&2
        sleep 2
    done
    echo "Cannot query block 1; check docker compose logs hub-node1 and port 26657." >&2
    exit 1
}
wait_genesis_block_hash


echo "# ---------------------------------------------------------------------------- #"
echo "#                  Get ME Node2, Node3, Node4 pubkeys                                   #"
echo "# ---------------------------------------------------------------------------- #"
ME_Node1_Pubkey=$(med tendermint show-validator --home ${ME_NODE1_HOME})
ME_Node2_Pubkey=$(med tendermint show-validator --home ${ME_NODES_HOME}/node2)
ME_Node3_Pubkey=$(med tendermint show-validator --home ${ME_NODES_HOME}/node3)
ME_Node4_Pubkey=$(med tendermint show-validator --home ${ME_NODES_HOME}/node4)
echo Node1 pubkey:${ME_Node1_Pubkey}
echo Node2 pubkey:${ME_Node2_Pubkey}
echo Node3 pubkey:${ME_Node3_Pubkey}
echo Node4 pubkey:${ME_Node4_Pubkey}

# offline tx
# A normal online tx can send only one tx per block, so the next tx waits until the previous one is on chain
# Use an offline tx when several txs must be sent back to back
# An offline tx works by tracking and setting the send sequence yourself
#-a $account_number   # on-chain account number of this address id
#--sequence ${seqID}  # on-chain sequence of this address. Increment by 1 on every tx to prevent replay
#--offline offline signing
#--fees tx fee
#--gas set the gas limit

#staking create-validator
##Create a new validator
##Stake tokens as the bond
##The node joins consensus and produces blocks
##Must meet the minimum stake
##Earns block rewards and tx fees
#staking new-region
##Add a new region on an existing validator
##No extra stake
##Main goal is to improve decentralization
##Usually used to run nodes in different places
##Inherits the primary validator's stake and rewards
##In short, create-validator creates a brand-new validator, while new-region extends an existing validator into a new location.

if [[ $Interactive_Mode =~ ^[Yy]$ ]]
then
    echo "# -------------------------------------------------------------------------------- #"
    echo "#      Custody create-validator mode                                                            #"
    echo "# -------------------------------------------------------------------------------- #"
    echo "Send the ME Node1, Node2, Node3, and Node4 pubkeys to custody"
    echo "Press Enter after custody finishes"
    read -p "Do you want to continue? (Y/n)"
else
    echo "# -------------------------------------------------------------------------------- #"
    echo "#      Automatic create-validator mode                                                            #"
    echo "# -------------------------------------------------------------------------------- #"
    seqID=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO} -o json --home $ME_NODE1_HOME | jq -r .base_account.sequence)
    account_number=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO} -o json --home $ME_NODE1_HOME | jq -r .base_account.account_number)
    echo "# ---------------------------------------------------------------------------- #"
    echo "#               create-validator ${ME_STAKING_AMOUNT_NODE2_REGION}                            #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "shell:med tx staking create-validator --pubkey "${ME_Node2_Pubkey}" --moniker node2 --amount "${ME_STAKING_AMOUNT_NODE2}" --validator-address "${USER_ADDR_ME_VAL2}" --region-id ${ME_STAKING_AMOUNT_NODE2_REGION} --from global_dao --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --chain-id "${ME_CHAIN_ID}" --fees=10000umec -a "$account_number" --sequence "${seqID}" --offline -y --home $ME_NODE1_HOME -o json  | jq -r '.txhash'"
    txhash=$(med tx staking create-validator --pubkey "${ME_Node2_Pubkey}" --moniker node2 \
        --amount "${ME_STAKING_AMOUNT_NODE2}" \
        --validator-address "${USER_ADDR_ME_VAL2}" \
        --region-id ${ME_STAKING_AMOUNT_NODE2_REGION} \
        --from global_dao \
        --commission-rate="0.10" \
        --commission-max-rate="0.20" \
        --commission-max-change-rate="0.01" \
        --chain-id "${ME_CHAIN_ID}" \
        --fees=10000umec -a "$account_number" --sequence "${seqID}" --offline -y --home ${ME_NODE1_HOME} -o json  | jq -r '.txhash')
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"


    seqID=$((seqID+1))
    echo "# ---------------------------------------------------------------------------- #"
    echo "#               create-validator ${ME_STAKING_AMOUNT_NODE3_REGION}             #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "shell:med tx staking create-validator --pubkey "${ME_Node3_Pubkey}" --moniker node3 --amount "${ME_STAKING_AMOUNT_NODE3}" --validator-address "${USER_ADDR_ME_VAL3}" --region-id ${ME_STAKING_AMOUNT_NODE3_REGION} --from global_dao --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --chain-id "${ME_CHAIN_ID}" --fees=10000umec -a "$account_number" --sequence "${seqID}" --offline -y --home $ME_NODE1_HOME -o json  | jq -r '.txhash'"
    txhash=$(med tx staking create-validator --pubkey "${ME_Node3_Pubkey}" --moniker node3 \
        --amount "${ME_STAKING_AMOUNT_NODE3}" \
        --validator-address "${USER_ADDR_ME_VAL3}" \
        --region-id ${ME_STAKING_AMOUNT_NODE3_REGION} \
        --from global_dao \
        --commission-rate="0.10" \
        --commission-max-rate="0.20" \
        --commission-max-change-rate="0.01" \
        --chain-id "${ME_CHAIN_ID}" \
        --fees=10000umec -a "$account_number" --sequence "${seqID}" --offline -y --home ${ME_NODE1_HOME} -o json  | jq -r '.txhash')
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"


    seqID=$((seqID+1))
    echo "# ---------------------------------------------------------------------------- #"
    echo "#               create-validator ${ME_STAKING_AMOUNT_NODE4_REGION}                                           #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "shell:med tx staking create-validator --pubkey "${ME_Node4_Pubkey}" --moniker node4 --amount "${ME_STAKING_AMOUNT_NODE4}" --validator-address "${USER_ADDR_ME_VAL4}" --region-id ${ME_STAKING_AMOUNT_NODE4_REGION} --from global_dao --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --chain-id "${ME_CHAIN_ID}"  --fees=10000umec -a "$account_number" --sequence "${seqID}" --offline -y -o json  --home $ME_NODE1_HOME -o json  | jq -r '.txhash'"
    txhash=$(med tx staking create-validator --pubkey "${ME_Node4_Pubkey}" --moniker node4 \
        --amount "${ME_STAKING_AMOUNT_NODE4}" \
        --validator-address "${USER_ADDR_ME_VAL4}" \
        --region-id ${ME_STAKING_AMOUNT_NODE4_REGION} \
        --from global_dao \
        --commission-rate="0.10" \
        --commission-max-rate="0.20" \
        --commission-max-change-rate="0.01" \
        --chain-id "${ME_CHAIN_ID}" \
        --fees=10000umec -a "${account_number}" --sequence "${seqID}" --offline -y -o json  --home ${ME_NODE1_HOME} -o json  | jq -r '.txhash')
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#               Query the validator list regionID                                    #"
echo "# ---------------------------------------------------------------------------- #"
while [ $(med query staking validators --home ${ME_NODE1_HOME} | grep -c "regionID") -ne 4 ]; do
    echo "4 validators' region IDs are not ready. Wait 5s and query again"
    sleep 5
    echo "Waiting for 4 region IDs... validators that already exist:"
    med query staking validators --home ${ME_NODE1_HOME} --output json | jq -r '.validators[] | .description.regionID'
done
echo "4 validators are ready!"
med query staking validators --home ${ME_NODE1_HOME} --output json | jq -r '.validators[] | .description.regionID'

if [[ $Interactive_Mode =~ ^[Nn]$ ]]
then
    echo "# ---------------------------------------------------------------------------- #"
    echo "#               Create region new-region zone                                           #"
    echo "# ---------------------------------------------------------------------------- #"
    seqID=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO}   --home ${ME_NODE1_HOME} --output json | jq -r .base_account.sequence)
    validator_count=$(med query staking validators --home ${ME_NODE1_HOME} --output json | jq '.validators | length')
    echo Current sequence:$seqID
    for ((i=0; i<validator_count; i++)); do
        operator_address=$(med query staking validators  --home ${ME_NODE1_HOME} --output json  | jq -r .validators[$i].operator_address)
        regionID=$(med query staking validators   --home ${ME_NODE1_HOME} --output json | jq -r .validators[$i].description.regionID)
        # Check whether the region already exists
        if med q staking regions --home ${ME_NODE1_HOME} --output json | jq -e ".region[] | select(.regionId == \"$regionID\")" > /dev/null; then
            echo "Region $regionID already exists, skipping..."
            continue
        fi
        # Create region
        echo "==============================================="
        echo "create new-region: ${regionID}"
        echo "shell:med tx staking new-region ${regionID} ${operator_address} --from global_dao  --chain-id ${ME_CHAIN_ID} --fees=10000umec --gas 500000 -a $account_number --sequence ${seqID} --offline -y -o json --home $ME_NODE1_HOME| jq -r .txhash"
        txhash=$(med tx staking new-region ${regionID} ${operator_address} --from global_dao  --chain-id ${ME_CHAIN_ID} --fees=10000umec --gas 500000 -a $account_number --sequence ${seqID} --offline -y -o json --home ${ME_NODE1_HOME}| jq -r .txhash)
        echo "Check whether the tx hash is on chain: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        seqID=$((seqID+1))
    done
fi

echo "Hub setup complete: four validators."
docker compose ps
