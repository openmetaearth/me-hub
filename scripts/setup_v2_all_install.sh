#!/bin/bash
#======================================================================================
#config start
#======================================================================================
echo "Services must keep running in the background, so start this script in its own session: screen -S install"
echo "Press Enter after you confirm the settings"
read -p "Do you want to continue?"

read -rp "Please enter PROJECT_ENV[beta]: " PROJECT_ENV
PROJECT_ENV=${PROJECT_ENV:-beta}
echo "PROJECT_ENV is set to: $PROJECT_ENV"
BASE_DIR="/data/docker-${PROJECT_ENV}"
BIN_DIR="${BASE_DIR}/bin"
echo "BASE_DIR is set to: $BASE_DIR"
echo "BIN_DIR is set to: $BIN_DIR"
read -rp "Please enter init_NODE_NAME[node1]: " NODE_NAME
NODE_NAME=${NODE_NAME:-node1}
echo "NODE_NAME is set to: $NODE_NAME"

read -rp "Please enter ME_CHAIN_ID[mechain_400-1]: " ME_CHAIN_ID
ME_CHAIN_ID=${ME_CHAIN_ID:-"mechain_400-1"}
echo "ME_CHAIN_ID is set to: $ME_CHAIN_ID"

read -rp "Please enter DA_CHAIN_ID[me-da]: " DA_CHAIN_ID
DA_CHAIN_ID=${DA_CHAIN_ID:-"me-da"}
echo "DA_CHAIN_ID is set to: $DA_CHAIN_ID"

read -rp "Please enter ROLLAPP_CHAIN_ID[mecheckin_401-1]: " ROLLAPP_CHAIN_ID
ROLLAPP_CHAIN_ID=${ROLLAPP_CHAIN_ID:-"mecheckin_401-1"}
echo "ROLLAPP_CHAIN_ID is set to: $ROLLAPP_CHAIN_ID"

read -rp "Please enter KEYRING_BACKEND_NAME[test]: " KEYRING_BACKEND_NAME
KEYRING_BACKEND_NAME=${KEYRING_BACKEND_NAME:-"test"}
echo "KEYRING_BACKEND_NAME is set to: $KEYRING_BACKEND_NAME"

MANAGER_IP=$(hostname -I | awk '{print $1}')
read -rp "Please enter LOCAL_ADDR[${MANAGER_IP}]: " LOCAL_ADDR
LOCAL_ADDR=${LOCAL_ADDR:-"${MANAGER_IP}"}
echo "LOCAL_ADDR is set to: $LOCAL_ADDR"

read -rp "Please enter ENABLE_MONITORING[N]: " ENABLE_MONITORING
ENABLE_MONITORING=${ENABLE_MONITORING:-"N"}
echo "ENABLE_MONITORING is set to: $ENABLE_MONITORING"

read -rp "Please enter Docker_ADDR_RANGE[172.28.20.0/24]: " Docker_ADDR_RANGE
Docker_ADDR_RANGE=${Docker_ADDR_RANGE:-"172.28.20.0/24"}
echo "Docker_ADDR_RANGE is set to: $Docker_ADDR_RANGE"
DOCKER_NETWORK_PREFIX=$(echo $Docker_ADDR_RANGE | awk -F'[./]' '{print $1"."$2"."$3}')
echo "DOCKER_NETWORK_PREFIX is set to: $DOCKER_NETWORK_PREFIX"

read -rp "Please enter Docker_Network_Name[me-beta]: " Docker_Network_Name
Docker_Network_Name=${Docker_Network_Name:-"me-beta"}
echo "Docker_Network_Name is set to: $Docker_Network_Name"

read -rp "manual sync keys?[Y]: " Sync_Key
Sync_Key=${Sync_Key:-"Y"}
echo "manual sync keys is set to: $Sync_Key"

read -rp "Whether to enable interactive mode?[N]: " Interactive_Mode
Interactive_Mode=${Interactive_Mode:-"N"}
echo "interactive mode is set to: $Interactive_Mode"

read -rp "Randomly generate mnemonic words?[N]: " Randomly_Generate_Mnemonic_Words
Randomly_Generate_Mnemonic_Words=${Randomly_Generate_Mnemonic_Words:-"N"}
echo "randomly generate mnemonic words is set to: $Randomly_Generate_Mnemonic_Words"

#Quick-install mode withdraws stake from me_earth so a test env can be stood up quickly. Turn it off to withdraw from each production region
read -rp "Quick install for test?[Y]: " INSTALL_TYPE
INSTALL_TYPE=${INSTALL_TYPE:-"Y"}
echo "Quick install for test is set to: $INSTALL_TYPE"

#read -rp "Enable Wasm Contract?[N]: " ENABLE_WASM_CONTRACT
#ENABLE_WASM_CONTRACT=${ENABLE_WASM_CONTRACT:-"N"}
#echo "Enable Wasm Contract is set to: $ENABLE_WASM_CONTRACT"

read -rp "link bin to system?[N]: " LINK_BIN_ACTION
LINK_BIN_ACTION=${LINK_BIN_ACTION:-"N"}
echo "link bin to system is set to: $LINK_BIN_ACTION"

# all=all = full install; rollapp_after = resume from rollapp IBC (do not rebuild hub/da/rollapp-node1)
# Or export RUN_STAGE=rollapp_after beforehand; pressing Enter keeps that value
read -rp "Please enter RUN_STAGE[all/rollapp_after] (default ${RUN_STAGE:-all}): " RUN_STAGE_INPUT
RUN_STAGE=${RUN_STAGE_INPUT:-${RUN_STAGE:-all}}
echo "RUN_STAGE is set to: $RUN_STAGE"



read -p "All Settings are initialized, Do you want to continue? (Y/n) "  REPLY
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Operation cancelled by user."
    exit 1
fi

#======================================================================================
echo "# ---------------------------------------------------------------------------- #"
echo "#                              set bin env                                     #"
echo "# ---------------------------------------------------------------------------- #"

cd "${BASE_DIR}" || { echo "error: BASE_DIR not found: ${BASE_DIR}"; exit 1; }
export PATH="${BIN_DIR}:${PATH}"
export LD_LIBRARY_PATH="${BASE_DIR}/lib:${LD_LIBRARY_PATH}"
#======================================================================================
# mnemonic
if [[ ! $Randomly_Generate_Mnemonic_Words =~ ^[Yy]$ ]]
then
    global_dict="sign wise lyrics find brass grass already urge invite vibrant arrow grow dune theme trust mix napkin ball pause sponsor bacon capable flip people"
else
    global_dict=$(meda-appd keys mnemonic)
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

USER_SEQ_ME_SEQUENCER=20                                                # rollapp account bound on the hub 
USER_SEQ_ME_IBC_DA=21                                                   # ibc -> da
USER_SEQ_ME_IBC_ROLLAPP=22                                              # ibc -> rollapp

USER_SEQ_ME_SEQUENCER2=30                                               # rollapp account bound on the hub 2
USER_SEQ_ME_SEQUENCER3=31                                               # rollapp account bound on the hub 3
da_user_index_val0=50
da_user_index_val1=51
da_user_index_val2=52
da_user_index_val3=53
da_user_index_val4=54

da_user_index_dao=56
USER_SEQ_DA_IBC_ME=57

da_user_index_bridge=60
da_user_index_full=61
da_user_index_light=62

rollapp_user_index_roluser=80                                            # rollapp rollapp user
rollapp_user_index_dao=81                                                # rollapp dao
rollapp_user_index_dev_operator=82                                       # devOperator devOperator association address
rollapp_user_index_ibc=83                                                # rollapp ibc
rollapp_user_index_sync_user=84                                          # rollapp rollapp sync user
#======================================================================================
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
USER_ADDR_ME_SEQUENCER=                                                  # 20
USER_ADDR_ME_IBC_DA=                                                     # 21
USER_ADDR_ME_IBC_ROLLAPP=                                                # 22
USER_ADDR_ME_SEQUENCER2=                                                 # 30
USER_ADDR_ME_SEQUENCER3=                                                 # 31
USER_ADDR_DA_DAO=                                                        # 56
USER_ADDR_DA_IBC=                                                        # 57
USER_ADDR_DA_VAL0=                                                       # 50
USER_ADDR_DA_VAL1=                                                       # 51
USER_ADDR_DA_VAL2=                                                       # 52
USER_ADDR_DA_VAL3=                                                       # 53
USER_ADDR_DA_VAL4=                                                       # 54
USER_ADDR_DA_BRIDGE=                                                     # 60
USER_ADDR_DA_FULL=                                                       # 61
USER_ADDR_DA_LIGHT=                                                      # 62
USER_ADDR_ROLLAPP_ADMIN=                                                 # 80
USER_ADDR_ROLLAPP_DAO=                                                   # 81
USER_ADDR_ROLLAPP_OPERATOR=                                              # 82
USER_ADDR_ROLLAPP_IBC=                                                   # 83
USER_ADDR_ROLLAPP_SYNC=                                                  # 84
#======================================================================================
#With Randomly_Generate_Mnemonic_Words=N the mnemonic produces the addresses below. Record only; do not edit
#USER_ADDR_ME_GLOBAL_DAO=me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t       # 1
#USER_ADDR_ME_MEID_DAO=me1mhhv44t0yyrm6rfgck2kl43ep7g2c4tj0hnqcf         # 2
#USER_ADDR_ME_DEV_OPERATOR=me1ega32ph9s9ahq5redwzaag05pd58m6xj37lfkw     # 3
#USER_ADDR_ME_AIRDROP=me1a5x2v9apf573suhdsxytf9rj0348m46xs86zh5          # 4
#USER_ADDR_ME_USER=me16dknggrc6jafynx90pjmcsk2zxz9j9v6pjv4qr             # 5
#USER_ADDR_ME_VAL1=me1rhqcuw38nytl29e5mlx0fyy6drp7wcucaasc9w             # 11
#USER_ADDR_ME_VAL2=me1sz05asn4svpr05d7y7r8d0vt2sj62u8rm39det             # 12
#USER_ADDR_ME_VAL3=me1hajj0g6tak48gs634jus8m8hkg0p3avcu6nw20             # 13
#USER_ADDR_ME_VAL4=me1uf459xuas5yvjdgugf34w62lkrx7a9zlqlf0ck             # 14
#USER_ADDR_ME_SEQUENCER=me1tes4hhrew6uurw5e4dnnlydenre9y7g4x52c39        # 20
#USER_ADDR_ME_IBC_DA=me15xdtgn8zex65v620whvz984xsefzzzyr2rckh0           # 21
#USER_ADDR_ME_IBC_ROLLAPP=me1zd2rgpqavfy8g9067gewvsrfa5d6wxfkavcklm      # 22
#USER_ADDR_ME_SEQUENCER2=me1p4r40t6cyehfhsyvsxdlct8a3m6gxjwlvxg05a       # 30
#USER_ADDR_ME_SEQUENCER3=me1psrnsjepss4dt7v6sr9rekh9z8j3qxp96s52xw       # 31
#USER_ADDR_DA_DAO=me1hvfxr8aywaqcxnk8vzjthgzpvrkgkzqm7gzpw3              # 56
#USER_ADDR_DA_IBC=me1zps5zt4tvz8uu4ajuq59sdjj2cf48w0nrv0gvd              # 57
#USER_ADDR_DA_VAL0=me1x3tj5zk2eczv5637a535a44d2wndm0efatx9rt             # 50
#USER_ADDR_DA_VAL1=me1jr6vw85ssgs872zf62weq43gltcwurq8gluzj5             # 51
#USER_ADDR_DA_VAL2=me1plg2v6ypykffzwk6jdkdcr8fznzgfqa3fxa459             # 52
#USER_ADDR_DA_VAL3=me165eudfny56tm59qkm62u28xx3wm8h57rj5kaya             # 53
#USER_ADDR_DA_VAL4=me1yf5a70fm2wkatgdycnxmzdgw76xs0y8jnnzuzq             # 54
#USER_ADDR_DA_BRIDGE=me1hhu97ug8dy2ffmrvsx3wrlxv79ha4n06m5eygp           # 60
#USER_ADDR_DA_FULL=me19fc63cm5gxq2zzep9f9srq7ccxyj6pqu9tchxj             # 61
#USER_ADDR_DA_LIGHT=me1c0lq9trgcxg9mqj8weg0vvqmeh8vlzc6s9594s            # 62
#USER_ADDR_ROLLAPP_ADMIN=me1fdhthjk30zvj44scz3239964ft4927mrupn7g6       # 80
#USER_ADDR_ROLLAPP_DAO=me1zpeg6t2rsgf356dc9zlpjw369zpcxlm55u7jm8         # 81
#USER_ADDR_ROLLAPP_OPERATOR=me1dtnx49zywysyxj0mt2ru335tpskllaks9cnn6l    # 82
#USER_ADDR_ROLLAPP_IBC=me1qk7jwgdmq5lfraq5l4p3qa5vq9vwa84rzhmhhc         # 83
#USER_ADDR_ROLLAPP_SYNC=me1hq4erh9g2guvxmexfd7c9la0qhamsx2r2ek3ga        # 84

#======================================================================================
# Public Setting
#======================================================================================
# Docker runtime images (genesis/keys still use the host ${BIN_DIR} binary)
MED_IMAGE="${MED_IMAGE:-ghcr.io/openmetaearth/med:v2.0.17}"
ROLLAPP_IMAGE="${ROLLAPP_IMAGE:-harbor.starex.xyz/rollapp/rollappd:v1.0.21}"
# hub / rollapp images run as root. In-container home is fixed to these paths (DA/rly still use ubuntu:24.04)
HUB_CONTAINER_HOME="${HUB_CONTAINER_HOME:-/root/.mechain}"
ROLLAPP_CONTAINER_HOME="${ROLLAPP_CONTAINER_HOME:-/root/.rollapp}"
HUB_MIN_GAS_PRICES="${HUB_MIN_GAS_PRICES:-0.02umec}"
ROLLAPP_MIN_GAS_PRICES="${ROLLAPP_MIN_GAS_PRICES:-0.001umec}"

# ME Hub directory
ME_NODES_HOME="${BASE_DIR}/nodes/hub-nodes"
# node directory
ME_NODE1_HOME="${ME_NODES_HOME}/${NODE_NAME}"

ME_NODE1_IP="${DOCKER_NETWORK_PREFIX}.10"
ME_NODE2_IP="${DOCKER_NETWORK_PREFIX}.11"
ME_NODE3_IP="${DOCKER_NETWORK_PREFIX}.12"
ME_NODE4_IP="${DOCKER_NETWORK_PREFIX}.13"
HUB_RPC_ADDR="${ME_NODE1_IP}:26657"
CHAIN_RPC_ADDR_ME="http://${ME_NODE1_IP}:26657"

ME_TENDERMINT_CONFIG_FILE="${ME_NODE1_HOME}/config/config.toml"
ME_CLIENT_CONFIG_FILE="${ME_NODE1_HOME}/config/client.toml"
ME_APP_CONFIG_FILE="${ME_NODE1_HOME}/config/app.toml"
ME_GENESIS_FILE="${ME_NODE1_HOME}/config/genesis.json"


#Minimum deposit when creating a governance proposal. The denom is umec so amounts stay uniform (config:app_state.gov.deposit_params.min_deposit[0].denom )"failed to validate genesis state: invalid minimum deposit: 0umec"
ME_MIN_DEPOSIT_AMOUNT="100000000"
#ME genesis stake amount
#Modular supply is 10 billion MEC, i.e. 1000000000000000000umec.
TOTAL_SUPPLY="1000000000000000000umec"  # 100 billion MEC
#minimum stake, unit udmec
DA_MIN_DEPOSIT_AMOUNT="1000000000000"
MIN_TREASURY_BALANCE="100000000"  # 1E umec minimum to keep
dao_TIA_AMOUNT="0udmec" # da dao own balance
# ME DA node directory
DA_NODES_HOME="${BASE_DIR}/nodes/da-nodes"

# node1 name
DA_NODE1_NAME=${NODE_NAME}

# node1 node directory
DA_NODE1_HOME="${DA_NODES_HOME}/${DA_NODE1_NAME}"
DA_BRIDGE_HOME="${DA_NODES_HOME}/bridge"
DA_FULL_HOME="${DA_NODES_HOME}/full"
DA_LIGHT_HOME="${DA_NODES_HOME}/light"
DA_NODE1_IP="${DOCKER_NETWORK_PREFIX}.101"
DA_NODE2_IP="${DOCKER_NETWORK_PREFIX}.102"
DA_NODE3_IP="${DOCKER_NETWORK_PREFIX}.103"
DA_NODE4_IP="${DOCKER_NETWORK_PREFIX}.104"
DA_BRIDGE_IP="${DOCKER_NETWORK_PREFIX}.190"
DA_LIGHT_IP="${DOCKER_NETWORK_PREFIX}.191"
DA_FULL_IP="${DOCKER_NETWORK_PREFIX}.192"
CHAIN_RPC_ADDR_DA="http://${DA_NODE1_IP}:26657"
CHAIN_RPC_ADDR_DA_BRIDGE="http://${DA_BRIDGE_IP}:26658"
# light node
DA_RPC_ADDR="${DA_LIGHT_IP}:26658"
DA_TENDERMINT_CONFIG_FILE="${DA_NODE1_HOME}/config/config.toml"
DA_APP_CONFIG_FILE="${DA_NODE1_HOME}/config/app.toml"
DA_GENESIS_FILE="${DA_NODE1_HOME}/config/genesis.json"


# ROLLAPP node directory
RAPP_NODE1_NAME=${NODE_NAME}
RAPP_NODES_HOME="${BASE_DIR}/nodes/rollapp-nodes"
RAPP_NODE1_HOME="${RAPP_NODES_HOME}/${RAPP_NODE1_NAME}"
RAPP_NODE1_IP="${DOCKER_NETWORK_PREFIX}.130"
RAPP_NODE2_IP="${DOCKER_NETWORK_PREFIX}.131"
RAPP_NODE3_IP="${DOCKER_NETWORK_PREFIX}.132"
CHAIN_RPC_ADDR_RAPP="http://${RAPP_NODE1_IP}:26657"

# denom name, default is stake
DENOM="urax"
#minimum stake, unit urax
ROLLAPP_MIN_DEPOSIT_AMOUNT="1000000"

# rly node directory
IBC_PORT=transfer
IBC_VERSION=ics20-1
RLY_CONFIG_FILE="${RLY_PATH}/config/config.yaml"
HUB_IBC_CONF_FILE="${RLY_PATH}/hub.json"
DA_IBC_CONF_FILE="${RLY_PATH}/da.json"
#======================================================================================
# rly-da node directory
RLY_DA_NODE_HOME="${BASE_DIR}/nodes/rly-da"
# rly-rapp node directory
RLY_RAPP_NODE_HOME="${BASE_DIR}/nodes/rly-rapp"

calculate() {
    echo "$1" | bc
}
# Convert timeout_commit (1s / 0.5s / 500ms) to seconds for sleep. bash $(( )) cannot handle decimals
timeout_commit_seconds() {
    local d="$1"
    case "$d" in
        *ms) echo "scale=3; ${d%ms}/1000" | bc ;;
        *s)  echo "${d%s}" ;;
        *)   echo "$d" ;;
    esac
}
umec_to_udmec() {
    local umec_amount=${1%umec}
    local result=$(echo "$umec_amount" | bc)
    echo "${result}udmec"   #1mec=100000000udmec
}


#======================================================================================
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
    ME_TIMEOUT_COMMIT="500ms"
else
    ME_DISPUTE_PERIOD_IN_BLOCKS="50"
    DISTR_EPOCH_IDENTIFIER="day"
    LOCKABLE_DURATIONS="14d"
    MAX_DEPOSIT_PERIOD="172800s"
    VOTING_PERIOD="172800s"
    EXPEDITED_VOTING_PERIOD="300s"
    ME_TIMEOUT_COMMIT="5s"
fi
REGION_HIGHT_WAIT="17280"
#======================================================================================
# ME-DA
#======================================================================================

val0_STAKING_AMOUNT="1100000000000udmec"   # da node1 genesis stake for this node
val0_TIA_AMOUNT="2000000000000udmec"     # da val0 own balance; DA coins are burned at expiry
#DA Node genesis stake for new validators
#Production mode
#- Create 4 nodes. Custody must sign and send MsgCreateValidator. @Lee.
#- Source of funds: withdraw from the top four ME-Hub mainnet region treasuries (IND, CHN, ME_EARTH, EXPERIENCE). Production should have more than four national treasuries, so exclude the EXPERIENCE region for stability.
#- Withdrawal ratio: 0.01% (1/10000) of each region's staked amount, computed dynamically as ME_SEND_DA_AMOUNT and IBC-transferred to the ME-DA network.
#Currently funded by the global_dao account, so the four amounts below must not exceed TARGET_BALANCE.
ME_SEND_DA_AMOUNT_VAL1="$(calculate "${STAKING_AMOUNT%umec} / 10000")umec"
ME_SEND_DA_AMOUNT_VAL2="$(calculate "${ME_STAKING_AMOUNT_NODE2%umec} / 10000")umec"
ME_SEND_DA_AMOUNT_VAL3="$(calculate "${ME_STAKING_AMOUNT_NODE3%umec} / 10000")umec"
ME_SEND_DA_AMOUNT_VAL4="$(calculate "${ME_STAKING_AMOUNT_NODE4%umec} / 10000")umec"
ME_SEND_IBC_AMOUNT="10000000umec"
ME_SEND_SEQUENCER1_AMOUNT="1000000000000umec"
ME_SEND_SEQUENCER2_AMOUNT="1000000000000umec"
ME_SEND_SEQUENCER3_AMOUNT="1000000000000umec"
#MEC reserved for other uses
ME_RESERVE_AMOUNT="200000000umec"
TARGET_BALANCE="$(calculate "${ME_SEND_DA_AMOUNT_VAL1%umec} + ${ME_SEND_DA_AMOUNT_VAL2%umec} + ${ME_SEND_DA_AMOUNT_VAL3%umec} + ${ME_SEND_DA_AMOUNT_VAL4%umec} + ${ME_SEND_IBC_AMOUNT%umec} + ${ME_RESERVE_AMOUNT%umec} + ${ME_SEND_SEQUENCER1_AMOUNT%umec} + ${ME_SEND_SEQUENCER2_AMOUNT%umec} + ${ME_SEND_SEQUENCER3_AMOUNT%umec}")"
# ME hub Wait until this hub height before creating a region. Each block rewards 792 MEC, i.e. 79200000000umec, which is how the height is estimated. The spare 30 is only so the treasury is not drained in the extreme case. Imprecise sleeps usually overshoot
ME_CREAT_REGION_HEIGHT="$(calculate "(${TARGET_BALANCE} / 79200000000) + 30")"
#DA_STAKING_AMOUNT must be greater than 1000000000000udmec and the sum must be less than ME_SEND_DA_AMOUNT * 100000000 (equality can still fail even with zero fees)
DA_STAKING_AMOUNT_VAL1=$(umec_to_udmec "${ME_SEND_DA_AMOUNT_VAL1}")
DA_STAKING_AMOUNT_NODE1_REGION="me_earth"
DA_STAKING_AMOUNT_VAL2=$(umec_to_udmec "${ME_SEND_DA_AMOUNT_VAL2}")
DA_STAKING_AMOUNT_NODE2_REGION="usa"
DA_STAKING_AMOUNT_VAL3=$(umec_to_udmec "${ME_SEND_DA_AMOUNT_VAL3}")
DA_STAKING_AMOUNT_NODE3_REGION="ind"
DA_STAKING_AMOUNT_VAL4=$(umec_to_udmec "${ME_SEND_DA_AMOUNT_VAL4}")
DA_STAKING_AMOUNT_NODE4_REGION="chn"

# staking module param: tokens must stay locked for 604800s/7 days, 18000s/5 hours.For Cosmos SDK-based chains, trust_period should usually be about 2/3 of the unbonding time (~2
# weeks) during which they can be financially punished (slashed) for misbehavior.The value is divided by 100 and multiplied by 85. That is the IBC trusting period (med q ibc client state 07-tendermint-0). 3024000s estimates trust_period at 28 days
# trusting_period=UnbondingTime / 100 * 85 // TODO: replace with percentage
# Watch the relationship with RLY_TIME_THRESHOLD
DA_UNBONDING_TIME=${ME_UNBONDING_TIME}
# staking module param: max validators
MAX_VALIDATORS="4"
#Dispute period. A block is final only after the dispute window. Recommended 50

if [[ $INSTALL_TYPE =~ ^[Yy]$  ]]
then
    DA_DISPUTE_PERIOD_IN_BLOCKS="5"
else
    DA_DISPUTE_PERIOD_IN_BLOCKS="50"
fi
#======================================================================================
# ME-Rollapp
#======================================================================================


# Initial balance rollapp_user_roluser_TOKEN_AMOUNT must be >= rollapp_user_roluser_NODE_STAKING_AMOUNT + rollapp_transfer_ibc_amount
if [[ $INSTALL_TYPE =~ ^[Yy]$  ]]
then
    rollapp_user_roluser_TOKEN_AMOUNT="600000000000$DENOM"
else
    rollapp_user_roluser_TOKEN_AMOUNT="6000000$DENOM"
fi
rollapp_user_roluser_NODE_STAKING_AMOUNT="5000000$DENOM"
rollapp_transfer_ibc_amount="500000$DENOM"

rollapp_user_dao_TOKEN_AMOUNT="0$DENOM"
rollapp_user_dev_operator_TOKEN_AMOUNT="0$DENOM"
#rollapp Time to reach the hub depends on state submission, so dymint batch_submit_max_time is 60s in test and also 60s in production;
#From rollapp to hub, a submit every 10 minutes means funds arrive in at least 15 minutes;
BATCH_SUBMIT_MAX_TIME="60s"
#max sequencers
MAX_SEQUENCERS=10
# Hub Register only node1. node2/node3 are full nodes and do not create-sequencer.
# Equal bond on three nodes can make v3 upgrade backfill point the proposer at node2.
SEQUENCER_MONIKER_NAME1="sequencer1"
SEQUENCER_MONIKER_NAME2="sequencer2"
SEQUENCER_MONIKER_NAME3="sequencer3"
# On-chain min bond is 100000000umec (create-sequencer with 100umec returns code 1011 insufficient bond)
# Must be less than the amount already sent to the sequencer account, ME_SEND_SEQUENCER*_AMOUNT
SEQUENCER_AMOUNT="100000000umec"
#======================================================================================
# ME-Rly
#======================================================================================


# Relayer update interval must be less than (UnbondingTime / 100 * 85). When more than 1/3 of validators change in a short window (adding one validator while bootstrapping 4 nodes), the interval must be larger than that so the trusting period still validates and the client can be refreshed
DA_UNBONDING_SECONDS=${DA_UNBONDING_TIME%s}
#docker Default update-period formula when docker does not set one
#RLY_TIME_THRESHOLD="$(( ((DA_UNBONDING_SECONDS /100 ) *85 )/ 3))s"
RLY_TIME_THRESHOLD="3600s"
SLEEP_RLY_TIME=${RLY_TIME_THRESHOLD%s}

#custom light client trusting period percentage ex. 66 (default: 85); this flag overrides the client-tp flag (default 85) --client-tp-percentage
#CLIENT_TP_PERCENTAGE="85"
#custom light client trusting period ex. 24h (default: 85% of chains reported unbonding time).  --client-tp
CLIENT_TP="$(( ( (DA_UNBONDING_SECONDS / 100) *85 )/ 3600 ))h"
#======================================================================================
# ME-DA
#======================================================================================
# da ibc token activation height. Set ibc_deadline_block_height in genesis. At 1000 (about 1.5 hours) the hub must IBC-stake on DA within 1000 blocks or the DA chain errors. Before that height only the whitelist and the DAO address can send txs. Default 100.
IBC_DEADLINE_BLOCK_HEIGHT=300


#======================================================================================
# Edit the settings above
#======================================================================================

#======================================================================================
#config end
#======================================================================================
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

create_sequencer(){

    local SEQUENCER_MONIKER_NAME=$1
    local NODE_NAME=$2
    local SEQ_NAME=$3
    SEQUENCER_PUBLIC_KEY=$(cat ${RAPP_NODES_HOME}/${NODE_NAME}/sequencer.info)
    echo "Create sequencer with the stake amount"
    local sequencer_metadata=$(jq -n --arg moniker "${SEQUENCER_MONIKER_NAME}" \
        '{Moniker: $moniker, Identity: "", Website: "", SecurityContact: "", Details: ""}')
    
    echo "Create sequencer, metadata: ${sequencer_metadata}"
    
    txhash=$(med tx sequencer create-sequencer ${SEQUENCER_PUBLIC_KEY} ${ROLLAPP_CHAIN_ID} "${sequencer_metadata}" ${SEQUENCER_AMOUNT} \
        --from ${SEQ_NAME} --chain-id ${ME_CHAIN_ID} \
        --fees 40000umec --gas 2000000 --keyring-backend ${KEYRING_BACKEND_NAME} --home ${ME_NODE1_HOME} --yes -o json | jq -r .txhash)
    
    echo "Sequencer create tx hash: ${txhash}"

    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
}


load_existing_addrs() {
    local kb=(--keyring-backend "${KEYRING_BACKEND_NAME}")
    echo "# Read addresses from the existing keyring for a rollapp_after resume"
    USER_ADDR_ME_GLOBAL_DAO=$(med keys show global_dao -a "${kb[@]}" --home "${ME_NODE1_HOME}")
    USER_ADDR_ME_SEQUENCER=$(med keys show sequencer -a "${kb[@]}" --home "${ME_NODE1_HOME}" 2>/dev/null || true)
    USER_ADDR_ME_SEQUENCER2=$(med keys show sequencer2 -a "${kb[@]}" --home "${ME_NODE1_HOME}" 2>/dev/null || true)
    USER_ADDR_ME_SEQUENCER3=$(med keys show sequencer3 -a "${kb[@]}" --home "${ME_NODE1_HOME}" 2>/dev/null || true)
    USER_ADDR_ROLLAPP_ADMIN=$(rollappd keys show roluser -a "${kb[@]}" --home "${RAPP_NODE1_HOME}")
    USER_ADDR_ROLLAPP_IBC=$(rollappd keys show ibc -a "${kb[@]}" --home "${RAPP_NODE1_HOME}")
    USER_ADDR_DA_VAL0=$(meda-appd keys show val0 -a "${kb[@]}" --home "${DA_NODE1_HOME}" 2>/dev/null || true)
    USER_ADDR_DA_VAL1=$(meda-appd keys show val1 -a "${kb[@]}" --home "${DA_NODE1_HOME}" 2>/dev/null || true)
    USER_ADDR_DA_VAL2=$(meda-appd keys show val2 -a "${kb[@]}" --home "${DA_NODE1_HOME}" 2>/dev/null || true)
    USER_ADDR_DA_VAL3=$(meda-appd keys show val3 -a "${kb[@]}" --home "${DA_NODE1_HOME}" 2>/dev/null || true)
    USER_ADDR_DA_VAL4=$(meda-appd keys show val4 -a "${kb[@]}" --home "${DA_NODE1_HOME}" 2>/dev/null || true)
    # hub->rollapp ibc receive address matches mnemonic index 22 and may not be in the hub keyring
    if [ -z "${USER_ADDR_ME_IBC_ROLLAPP}" ]; then
        local tmp_home
        tmp_home=$(mktemp -d)
        local out
        out=$(add_key "meda-appd" "me_user_ibc_rollapp_name" ${USER_SEQ_ME_IBC_ROLLAPP} "--home ${tmp_home}")
        USER_ADDR_ME_IBC_ROLLAPP=$(echo "${out}" | grep -oP 'address: \K[^\s]+')
        rm -rf "${tmp_home}"
    fi
    echo "GLOBAL_DAO=${USER_ADDR_ME_GLOBAL_DAO}"
    echo "ROLLAPP_ADMIN=${USER_ADDR_ROLLAPP_ADMIN}"
    echo "ROLLAPP_IBC=${USER_ADDR_ROLLAPP_IBC}"
    echo "ME_IBC_ROLLAPP=${USER_ADDR_ME_IBC_ROLLAPP}"
}

setup_after_rollapp_node1() {
echo "# ---------------------------------------------------------------------------- #"
echo "#       Wait until rollapp RPC is up and producing blocks (the IBC client is created later by rly tx link)      #"
echo "# ---------------------------------------------------------------------------- #"
while true; do
    height=$(rollappd status --home ${RAPP_NODE1_HOME} 2>/dev/null | jq -r '.SyncInfo.latest_block_height // .sync_info.latest_block_height // empty')
    height=${height//\"/}
    if [[ -n "${height}" && "${height}" -ge 1 ]] 2>/dev/null; then
        echo "rollapp is producing blocks, height=${height}"
        break
    fi
    echo "rollapp is not producing blocks yet, or RPC is not ready, sleep 5s"
    sleep 5
done

echo "# ---------------------------------------------------------------------------- #"
echo "#      Deploy rly-rollapp                                                          #"
echo "# ---------------------------------------------------------------------------- #"
echo "Set channel info"
echo "Initialize the rly config file:rly config init --home ${RLY_RAPP_NODE_HOME}"
rly config init --home ${RLY_RAPP_NODE_HOME}
echo Import config
cat > ${RLY_RAPP_NODE_HOME}/config/config.yaml <<EOF
global:
    api-listen-addr: :5183
    timeout: 10s
    memo: ""
    light-cache-size: 20
    log-level: info
    ics20-memo-limit: 0
    max-receiver-size: 150
chains:
    ${ME_CHAIN_ID}:
        type: cosmos
        value:
            key-directory: ${RLY_RAPP_NODE_HOME}/keys/${ME_CHAIN_ID}
            key: me_user_ibc_rollapp_name
            chain-id: ${ME_CHAIN_ID}
            rpc-addr: ${CHAIN_RPC_ADDR_ME}
            account-prefix: me
            keyring-backend: ${KEYRING_BACKEND_NAME}
            gas-adjustment: 1.5
            gas-prices: 0.02umec
            min-gas-amount: 0
            max-gas-amount: 0
            debug: true
            timeout: 10s
            block-timeout: ""
            output-format: json
            sign-mode: direct
            extra-codecs: []
            coin-type: null
            signing-algorithm: ""
            broadcast-mode: batch
            min-loop-duration: 0s
            extension-options: []
            feegrants: null
    ${ROLLAPP_CHAIN_ID}:
        type: cosmos
        value:
            key-directory: ${RLY_RAPP_NODE_HOME}/keys/${ROLLAPP_CHAIN_ID}
            key: rollapp_user_ibc_me_name
            chain-id: ${ROLLAPP_CHAIN_ID}
            rpc-addr: ${CHAIN_RPC_ADDR_RAPP}
            account-prefix: me
            keyring-backend: ${KEYRING_BACKEND_NAME}
            gas-adjustment: 1.2
            gas-prices: 0.02urax
            min-gas-amount: 0
            max-gas-amount: 0
            debug: true
            timeout: 10s
            block-timeout: ""
            output-format: json
            sign-mode: direct
            extra-codecs: []
            coin-type: null
            signing-algorithm: ""
            broadcast-mode: batch
            min-loop-duration: 0s
            extension-options: []
            feegrants: null
paths:
    hub-rollapp:
        src:
            chain-id: ${ME_CHAIN_ID}
        dst:
            chain-id: ${ROLLAPP_CHAIN_ID}
        src-channel-filter:
            rule: ""
            channel-list: []
EOF

echo "# ---------------------------------------------------------------------------- #"
echo "#      Whether to import me_user_ibc_rollapp_name and rollapp_user_ibc_me_name keys            #"
echo "# ---------------------------------------------------------------------------- #"
if [[ "${Sync_Key}" == [Yy] ]]; then
    mkdir -p ${RLY_RAPP_NODE_HOME}/keys/${ME_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    mkdir -p ${RLY_RAPP_NODE_HOME}/keys/${ROLLAPP_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    echo "Handle me_user_ibc_rollapp_name, same as USER_ADDR_ME_IBC_ROLLAPP"
    me_user_ibc_rollapp_name=$(add_key "meda-appd" "me_user_ibc_rollapp_name" ${USER_SEQ_ME_IBC_ROLLAPP} "--home ./tmpxx10")
    echo "${me_user_ibc_rollapp_name}"
    USER_ME_IBC_DA_NAME=$(echo "${me_user_ibc_rollapp_name}" | grep -oP 'address: \K[^\s]+')
    USER_ME_IBC_DA_NAME_PUBKEY=$(echo "${me_user_ibc_rollapp_name}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    cp -rf ./tmpxx10/keyring-${KEYRING_BACKEND_NAME}/*  ${RLY_RAPP_NODE_HOME}/keys/${ME_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    rm -rf ./tmpxx10

    echo "Handle rollapp_user_ibc_me_name, same as the DA-side USER_ADDR_ROLLAPP_IBC-ibc"
    rollapp_user_ibc_me_name=$(add_key "meda-appd" "rollapp_user_ibc_me_name" ${rollapp_user_index_ibc} "--home ./tmpxx10")
    echo "$rollapp_user_ibc_me_name"
    USER_ROLLAPP_IBC_ME_NAME=$(echo "$rollapp_user_ibc_me_name" | grep -oP 'address: \K[^\s]+')
    USER_ROLLAPP_IBC_ME_NAME_PUBKEY=$(echo "$rollapp_user_ibc_me_name" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    cp -rf ./tmpxx10/keyring-${KEYRING_BACKEND_NAME}/*  ${RLY_RAPP_NODE_HOME}/keys/${ROLLAPP_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    rm -rf ./tmpxx10
fi


echo "# Avoid the 07-tendermint-0: light client not found error. Wait 10s"
sleep 10

echo "create link hub-rollapp:rly tx link hub-rollapp --src-port transfer --dst-port transfer --version ics20-1 --home ${RLY_RAPP_NODE_HOME}"

rly tx link hub-rollapp --src-port transfer --dst-port transfer --version ics20-1 --home ${RLY_RAPP_NODE_HOME}

# Check the previous command's exit status
if [ $? -ne 0 ]; then
    echo "creat command failed, exiting"
    exit 1
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#      rly-rollapp check                                                          #"
echo "# ---------------------------------------------------------------------------- #"
echo Check ${ME_CHAIN_ID} client status
rly q clients ${ME_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}
echo Check ${ROLLAPP_CHAIN_ID} client status
rly q clients ${ROLLAPP_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}

echo Check connection status
rly q connections ${ME_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}



echo "Query on-chain channel info"
rly q channels ${ME_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}
echo "-------------------------------------------------------------------------------"
rly q channels ${ROLLAPP_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}

echo "Get channel ID"
echo "Query ME_CHANNEL_ID:rly q channels "${ME_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ROLLAPP_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id'"
ME_CHANNEL_ID=$(rly q channels "${ME_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ROLLAPP_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id')

echo "Query RAPP_CHANNEL_ID:rly q channels "${ROLLAPP_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ME_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id'"
RAPP_CHANNEL_ID=$(rly q channels "${ROLLAPP_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ME_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id')

echo "me -> rapp, channel_id: ${ME_CHANNEL_ID}"
echo "rapp -> me, channel_id: ${RAPP_CHANNEL_ID}"


echo "Set channel timeout"
rly transact client ${ME_CHAIN_ID} ${ROLLAPP_CHAIN_ID} hub-rollapp --client-tp ${CLIENT_TP} --home ${RLY_RAPP_NODE_HOME}


echo "Add relayer to docker-compose.yml"
cat >>docker-compose.yml<<EOF
  rly-rollapp:
    restart: unless-stopped
    image: ubuntu:24.04
    user: "1000:1000"
    networks:
      - ${Docker_Network_Name}
    volumes:
      - ${BIN_DIR}/rly:/bin/rly
      - ${RLY_RAPP_NODE_HOME}:/home/ubuntu/.relayer
    command:
      - "rly"
      - "start"
      - "hub-rollapp"
      - "--time-threshold"
      - "${RLY_TIME_THRESHOLD}"
      # - "--no-flush"

EOF

docker compose up -d rly-rollapp

echo "# ---------------------------------------------------------------------------- #"
echo "#      rollapp to me-hub  ibc                                                     #"
echo "# ---------------------------------------------------------------------------- #"

if [[ $Interactive_Mode =~ ^[Nn]$ ]]
then
    echo "ibc: rollapp to me-hub"
    echo "Query me-hub ibc-rollapp address ${USER_ADDR_ME_IBC_ROLLAPP} balance:"
    med query bank balances ${USER_ADDR_ME_IBC_ROLLAPP} --home ${ME_NODE1_HOME}

    echo "rollapp native-token holder sends an IBC transfer of amount ${rollapp_transfer_ibc_amount}"
    echo "rollappd tx ibc-transfer transfer transfer ${RAPP_CHANNEL_ID} ${USER_ADDR_ME_IBC_ROLLAPP} ${rollapp_transfer_ibc_amount} --from ${USER_ADDR_ROLLAPP_ADMIN} --fees=2000urax --chain-id ${ROLLAPP_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME} -y"

    txhash=$(rollappd tx ibc-transfer transfer transfer ${RAPP_CHANNEL_ID} ${USER_ADDR_ME_IBC_ROLLAPP} ${rollapp_transfer_ibc_amount} --from ${USER_ADDR_ROLLAPP_ADMIN} --fees=2000urax --chain-id ${ROLLAPP_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME} -y --output json | jq -r '.txhash')

    check_tx_status rollappd ${txhash} "--home $RAPP_NODE1_HOME"
else
    echo "# -------------------------------------------------------------------------------- #"
    echo "#      Custody transfer mode                                                                  #"
    echo "# -------------------------------------------------------------------------------- #"
    echo "ibc: rollapp to me-hub"
    echo "Query me-hub ibc-rollapp address ${USER_ADDR_ME_IBC_ROLLAPP} balance:"
    echo  "med query bank balances ${USER_ADDR_ME_IBC_ROLLAPP} --home ${ME_NODE1_HOME}"
    echo "Press Enter if this looks correct"
    read -p "Do you want to continue? (Y/n)"
fi

echo  "Check whether the rollapp is active. Fast finality starts after activation"
echo "Run shell: med q rollapp show ${ROLLAPP_CHAIN_ID}  --home \"${ME_NODE1_HOME}\" -o json | jq .rollapp.genesis_state.transfers_enabled"
while true; do
  rollapp_status=$(med q rollapp show ${ROLLAPP_CHAIN_ID}  --home "${ME_NODE1_HOME}" -o json | jq .rollapp.genesis_state.transfers_enabled)
  if [[ ${rollapp_status} == "true" ]];then
      echo "rollapp activated"
      break
  else
      echo "rollapp not active, sleep 1s"
      sleep 1
      continue
  fi
done


if [[ $Interactive_Mode =~ ^[Nn]$ ]]
then
    echo "ibc: rollapp to me-hub"
    echo "Run shell: med tx rollapp skip-delay-rollapp ${ROLLAPP_CHAIN_ID} true --from ${USER_ADDR_ME_GLOBAL_DAO} --gas 500000 --chain-id ${ME_CHAIN_ID} --keyring-backend=${KEYRING_BACKEND_NAME} --home  $ME_NODE1_HOME  -y --output json | jq -r '.txhash'"
    txhash=$(med tx rollapp skip-delay-rollapp ${ROLLAPP_CHAIN_ID} true --from ${USER_ADDR_ME_GLOBAL_DAO} --gas 500000 --chain-id ${ME_CHAIN_ID} --keyring-backend=${KEYRING_BACKEND_NAME} --home  $ME_NODE1_HOME  -y --output json | jq -r '.txhash')
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    echo "set skip-delay-rollapp success"
    sleep 2

    #===================================================================
    # me-hub to rollapp

    echo ""
    echo "ibc: me-hub to rollapp"
    echo ""

    echo "Query IBC balance shell: rollappd query bank balances ${USER_ADDR_ROLLAPP_IBC} --home ${RAPP_NODE1_HOME}"
    rollappd query bank balances ${USER_ADDR_ROLLAPP_IBC} --home ${RAPP_NODE1_HOME}

    echo "ME user transfers 10mec to the rollapp ibc address shell: med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${USER_ADDR_ROLLAPP_IBC}  ${ME_SEND_IBC_AMOUNT} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=100000umec  --chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME} --keyring-backend ${KEYRING_BACKEND_NAME}  -y"
    txhash=$(med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${USER_ADDR_ROLLAPP_IBC} ${ME_SEND_IBC_AMOUNT} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=100000umec  --chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME} --keyring-backend ${KEYRING_BACKEND_NAME}  -y --output json | jq -r '.txhash')
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"

    #===================================================================
else
    echo "# -------------------------------------------------------------------------------- #"
    echo "#      Custody transfer mode                                                                  #"
    echo "# -------------------------------------------------------------------------------- #"
    echo "ibc: rollapp to me-hub"
    echo "Query me-hub ibc-rollapp address ${USER_ADDR_ROLLAPP_IBC} balance:"
    echo  "rollappd query bank balances ${USER_ADDR_ROLLAPP_IBC} --home ${RAPP_NODE1_HOME}"
    echo "Press Enter if this looks correct"
    read -p "Do you want to continue? (Y/n)"
fi
echo "========================="
# Keep checking whether the address holds the given denom
watch_addr_balance_type() {
  local binary="$1"
  local token_prefix="$2"   # denom prefix
  local user_addr="$3"      # user address
  local extra_args="$4"     # extra args
  local check_status=0      # check status

  while [ $check_status -eq 0 ]; do
    token_prefix='ibc'
    COIN_TYPES=$($binary query bank balances ${user_addr} ${extra_args}  -o json | jq -r .balances[].denom)
    for COIN_TYPE in ${COIN_TYPES}; do
      if [[ $COIN_TYPE == *"$token_prefix"* ]];then
        check_status=1
        break
      fi
    done
    if [ ${check_status} -eq 0 ]; then
      if [ "${binary}" == "med" ]; then
        echo "addr ${user_addr} not found coin type: ${token_prefix} sleep 30s,at last need about 200s"
        sleep 30
      else
        echo "addr ${user_addr} not found coin type: ${token_prefix} sleep 2s,at last need about 8s"
        sleep 2
      fi
    fi
  done

  echo ""
  echo "addr ${user_addr} coin type: ${token_prefix}"
  $binary query bank balances ${user_addr} ${extra_args}
}


watch_addr_balance_type rollappd ibc ${USER_ADDR_ROLLAPP_IBC} "--home ${RAPP_NODE1_HOME}"

echo "========================="

echo ""
echo "check me chain ibc_rollapp user urax"
echo ""

watch_addr_balance_type med ibc ${USER_ADDR_ME_IBC_ROLLAPP} "--home ${ME_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#              Initialize RollApp Node2/Node3 as full nodes (they do not produce blocks)                  #"
echo "# ---------------------------------------------------------------------------- #"
echo "Check static config:"
echo "dym_account_name: $(grep '^dym_account_name' ${RAPP_NODES_HOME}/node1/config/dymint.toml)"
echo "keyring_home_dir: $(grep '^keyring_home_dir' ${RAPP_NODES_HOME}/node1/config/dymint.toml)"
echo "da_config: $(grep '^da_config' ${RAPP_NODES_HOME}/node1/config/dymint.toml)"
echo "settlement_node_address: $(grep '^settlement_node_address' ${RAPP_NODES_HOME}/node1/config/dymint.toml)"
# full node: copy node1 genesis/dymint and keep this machine's priv_validator_key (not registered on the Hub).
# Do not change dym_account_name. Keep using sequencer only for settlement queries; it will not produce blocks.
rollapp_init_full_node() {
    local NODE_NAME=$1
    local NODE_HOME="${RAPP_NODES_HOME}/${NODE_NAME}"
    local NODE_IP=$2

    local CONFIG_DIRECTORY="$NODE_HOME/config"
    local TENDERMINT_CONF="$CONFIG_DIRECTORY/config.toml"
    local CLIENT_CONFIG_FILE="$CONFIG_DIRECTORY/client.toml"
    local DYMINT_CONFIG_FILE="$CONFIG_DIRECTORY/dymint.toml"

    echo "========================================"
    echo "init full node: ${NODE_NAME} - ${NODE_HOME}"

    rollappd init ${NODE_NAME} --chain-id="${ROLLAPP_CHAIN_ID}" --home ${NODE_HOME}
    cp -rf ${RAPP_NODE1_HOME}/config/app.toml ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/client.toml ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/config.toml ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/dymint.toml ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/genesis.json ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/gentx ${CONFIG_DIRECTORY}/
    mkdir -p "${NODE_HOME}/sequencer_keys"
    cp -rf ${RAPP_NODE1_HOME}/sequencer_keys/. "${NODE_HOME}/sequencer_keys/"
    sed -i "s|^moniker = .*$|moniker = \"${NODE_NAME}\"|" "${TENDERMINT_CONF}"
    sed -i "s|^node = .*$|node = \"tcp://${NODE_IP}:26657\"|" "${CLIENT_CONFIG_FILE}"
    echo "Check static config:"
    echo "dym_account_name: $(grep '^dym_account_name' ${DYMINT_CONFIG_FILE})"
    echo "keyring_home_dir: $(grep '^keyring_home_dir' ${DYMINT_CONFIG_FILE})"
    echo "da_config: $(grep '^da_config' ${DYMINT_CONFIG_FILE})"
    echo "settlement_node_address: $(grep '^settlement_node_address' ${DYMINT_CONFIG_FILE})"
}
rollapp_init_full_node node2 ${RAPP_NODE2_IP}
rollapp_init_full_node node3 ${RAPP_NODE3_IP}

echo "# ---------------------------------------------------------------------------- #"
echo "#      Node2/Node3 are full nodes. They do not create-sequencer on the Hub and do not produce blocks.        #"
echo "# ---------------------------------------------------------------------------- #"
echo "Check sequencer status (expect only node1/sequencer): med q sequencer list-sequencer --home ${ME_NODE1_HOME}"
med q sequencer list-sequencer --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#      Create rollapp Node2-Node3 Docker Compose.                                 #"
echo "# ---------------------------------------------------------------------------- #"
cat >>docker-compose.yml<<EOF
  rollapp-node2:
    <<: *rollapp-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${RAPP_NODE2_IP}
    volumes:
      - ${RAPP_NODES_HOME}/node2:${ROLLAPP_CONTAINER_HOME}
  rollapp-node3:
    <<: *rollapp-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${RAPP_NODE3_IP}
    volumes:
      - ${RAPP_NODES_HOME}/node3:${ROLLAPP_CONTAINER_HOME}
EOF
echo "# ---------------------------------------------------------------------------- #"
echo "#                              RollApp Read Node ID                               #"
echo "# ---------------------------------------------------------------------------- #"
RAPP_NODEID_NODE1=$(rollappd dymint show-node-id --home ${RAPP_NODES_HOME}/node1  2>&1 | tail -n 1)
RAPP_NODEID_NODE2=$(rollappd dymint show-node-id --home ${RAPP_NODES_HOME}/node2 2>/dev/null | tail -n 1)
RAPP_NODEID_NODE3=$(rollappd dymint show-node-id --home ${RAPP_NODES_HOME}/node3 2>/dev/null | tail -n 1)
echo Rapp_Node1_ID:${RAPP_NODEID_NODE1}
echo Rapp_Node2_ID:${RAPP_NODEID_NODE2}
echo Rapp_Node3_ID:${RAPP_NODEID_NODE3}
echo "# ---------------------------------------------------------------------------- #"
echo "#                         Edit Node2-Node5 persistent_peers                #"
echo "# ---------------------------------------------------------------------------- #"
RAPP_NodeID_Node1_DYMINT_P2P="/ip4/${RAPP_NODE2_IP}/tcp/26656/p2p/${RAPP_NODEID_NODE2},/ip4/${RAPP_NODE3_IP}/tcp/28656/p2p/${RAPP_NODEID_NODE3}"
RAPP_NodeID_Node1_CONFIG_P2P="${RAPP_NODEID_NODE2}@${RAPP_NODE2_IP}:26656,${RAPP_NODEID_NODE3}@${RAPP_NODE3_IP}:28656"
echo RAPP_NodeID_Node1_DYMINT_P2P:$RAPP_NodeID_Node1_DYMINT_P2P

RAPP_NodeID_Node2_DYMINT_P2P="/ip4/${RAPP_NODE1_IP}/tcp/26656/p2p/${RAPP_NODEID_NODE1},/ip4/${RAPP_NODE3_IP}/tcp/28656/p2p/${RAPP_NODEID_NODE3}"
RAPP_NodeID_Node2_CONFIG_P2P="${RAPP_NODEID_NODE1}@${RAPP_NODE1_IP}:26656,${RAPP_NODEID_NODE3}@${RAPP_NODE3_IP}:28656"
echo RAPP_NodeID_Node2_DYMINT_P2P:$RAPP_NodeID_Node2_DYMINT_P2P

RAPP_NodeID_Node3_DYMINT_P2P="/ip4/${RAPP_NODE1_IP}/tcp/26656/p2p/${RAPP_NODEID_NODE1},/ip4/${RAPP_NODE2_IP}/tcp/28656/p2p/${RAPP_NODEID_NODE2}"
RAPP_NodeID_Node3_CONFIG_P2P="${RAPP_NODEID_NODE1}@${RAPP_NODE1_IP}:26656,${RAPP_NODEID_NODE2}@${RAPP_NODE2_IP}:28656"
echo RAPP_NodeID_Node3_DYMINT_P2P:$RAPP_NodeID_Node3_DYMINT_P2P


sed -i 's|^p2p_persistent_nodes = .*$|p2p_persistent_nodes = "'"${RAPP_NodeID_Node2_DYMINT_P2P}"'"|' ${RAPP_NODES_HOME}/node2/config/dymint.toml
sed -i 's|^persistent_peers = .*$|persistent_peers = "'"${RAPP_NodeID_Node2_CONFIG_P2P}"'"|' ${RAPP_NODES_HOME}/node2/config/config.toml


sed -i 's|^p2p_persistent_nodes = .*$|p2p_persistent_nodes = "'"${RAPP_NodeID_Node3_DYMINT_P2P}"'"|' ${RAPP_NODES_HOME}/node3/config/dymint.toml
sed -i 's|^persistent_peers = .*$|persistent_peers = "'"${RAPP_NodeID_Node3_CONFIG_P2P}"'"|' ${RAPP_NODES_HOME}/node3/config/config.toml

echo "# ---------------------------------------------------------------------------- #"
echo "#                  Start rollapp  Node2 Node3                                     #"
echo "# --------------------- -------------------------------------------------------#"

docker compose up -d rollapp-node2
echo "Wait 5s for services to finish starting"
sleep 5
rly  tx update-clients hub-rollapp --home ${RLY_RAPP_NODE_HOME}
docker compose up -d rollapp-node3
echo "Wait 5s for services to finish starting"
sleep 5
rly  tx update-clients hub-rollapp --home ${RLY_RAPP_NODE_HOME}
echo "# ---------------------------------------------------------------------------- #"
echo "#                  Edit node1 persistent_peers                              #"
echo "# --------------------- ------------------------------------------------------- #"
docker compose  down rollapp-node1
sed -i 's|^p2p_persistent_nodes = .*$|p2p_persistent_nodes = "'"${RAPP_NodeID_Node1_DYMINT_P2P}"'"|' ${RAPP_NODES_HOME}/node1/config/dymint.toml
sed -i 's|^persistent_peers = .*$|persistent_peers = "'"${RAPP_NodeID_Node1_CONFIG_P2P}"'"|' ${RAPP_NODES_HOME}/node1/config/config.toml
docker compose up -d rollapp-node1


if [[ $LINK_BIN_ACTION =~ ^[Yy]$ ]]
then
    ln -sfn ${ME_NODE1_HOME}  $HOME/.mechain
    ln -sfn ${RAPP_NODE1_HOME}  $HOME/.rollapp
    ln -sfn ${DA_NODE1_HOME}  $HOME/.meda-app
    ln -sfn ${DA_LIGHT_HOME}  $HOME/.meda-light-me-da
fi

if [[ $ENABLE_WASM_CONTRACT =~ ^[Yy]$ ]]
then
    echo "# ---------------------------------------------------------------------------- #"
    echo "#      Deploy contract                                                                 #"
    echo "# ---------------------------------------------------------------------------- #"
    #upload contract
    #med tx wasm store contract/matching_contract.wasm --from global_dao --gas auto --gas-adjustment 1.3 --chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME} -y
    #txhash=$(med tx wasm store contract/matching_contract.wasm --from global_dao --gas auto --gas-adjustment 1.3 --chain-id mechain_400-1 --home ${ME_NODE1_HOME} -y --output json | jq -r '.txhash')
    #check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    #instantiate contract
    #med tx wasm instantiate <CODE_ID> '{"global_dao":"<admin_address>","fee_collector":"<fee_collector_address>"}' --from global_dao --label "matching-v1" --global_dao <admin_address> --chain-id mechain_400-1 --home ${ME_NODE1_HOME} -y
    #txhash=$(med tx wasm instantiate <CODE_ID> '{"admin":"<admin_address>","fee_collector":"<fee_collector_address>"}' --from global_dao --label "matching-v1" --global_dao <admin_address> --chain-id mechain_400-1 --home ${ME_NODE1_HOME} -y --output json | jq -r '.txhash')
    #check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    #record contract address CONTRACT_ADDR
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#      Deploy status                                                                 #"
echo "# ---------------------------------------------------------------------------- #"
echo "Current HUB status"
med status --node tcp://${ME_NODE1_IP}:26657
echo "Current DA status"
meda-appd status --node tcp://${DA_NODE1_IP}:26657
echo "Current rollapp status"
rollappd status --node tcp://${RAPP_NODE1_IP}:26657
echo "Current IBC status"
med  q ibc client state 07-tendermint-0 --home ${ME_NODE1_HOME} -o json
echo "Current ME global_dao balance"
med query bank balances ${USER_ADDR_ME_GLOBAL_DAO} --home ${ME_NODE1_HOME}  -o json
echo "Current DA val0 balance (a balance before DA reaches height ${IBC_DEADLINE_BLOCK_HEIGHT} is normal)"
meda-appd query bank balances ${USER_ADDR_DA_VAL0} --home ${DA_NODE1_HOME}  -o json
echo "Current DA val1 balance"
meda-appd query bank balances ${USER_ADDR_DA_VAL1} --home ${DA_NODE1_HOME}  -o json
echo "Current DA val2 balance"
meda-appd query bank balances ${USER_ADDR_DA_VAL2} --home ${DA_NODE1_HOME}  -o json
echo "Current DA val3 balance"
meda-appd query bank balances ${USER_ADDR_DA_VAL3} --home ${DA_NODE1_HOME}  -o json
echo "Current DA val4 balance"
meda-appd query bank balances ${USER_ADDR_DA_VAL4} --home ${DA_NODE1_HOME}  -o json
echo "Current rollapp roluser balance"
rollappd query bank balances ${USER_ADDR_ROLLAPP_ADMIN} --home ${RAPP_NODE1_HOME}  -o json
echo "Current rollapp ibc balance"
rollappd query bank balances ${USER_ADDR_ROLLAPP_IBC} --home ${RAPP_NODE1_HOME}  -o json
echo "hub REGION_ID"
med query staking validators --home ${ME_NODE1_HOME} -o json | jq -r '.validators[] | [.description.regionID, .tokens, .status] | @tsv' | column -t -s $'\t' -N "REGION_ID,TOKENS,STATUS"
echo "hub MONIKER"
meda-appd query staking validators --home ${DA_NODE1_HOME} -o json | jq -r '.validators[] | [.description.moniker, .tokens, .status] | @tsv' | column -t -s $'\t' -N "MONIKER,TOKENS,STATUS"
echo "IBC trusting period"
rly query clients-expiration hub-rollapp --home ${RLY_RAPP_NODE_HOME}
rly query clients-expiration hub-meda   --home  ${RLY_DA_NODE_HOME}

end_time=$(date +%s)
duration=$((end_time - start_time))
echo "-----------------------------------"
echo "Script finished"
echo "Elapsed: $duration seconds"
}

start_time=$(date +%s)
echo "# ---------------------------------------------------------------------------- #"
echo "#                              check bin env                                  #"
echo "# ---------------------------------------------------------------------------- #"
pre_check_file() {
  local files=("$@")
  for file in "${files[@]}"; do

    if [ ! -f "$file" ]; then
      echo "error: $file file not found"
      exit 1
    else
      echo "check: $file"
    fi

  done
}
pre_check_file ${BIN_DIR}/med ${BIN_DIR}/meda-appd ${BIN_DIR}/meda ${BIN_DIR}/rly ${BIN_DIR}/rollappd ${BASE_DIR}/lib/libwasmvm.x86_64.so /lib/libwasmvm.x86_64.so
sudo ldconfig

if [[ $ENABLE_WASM_CONTRACT =~ ^[Yy]$ ]]
then
    pre_check_file ${BASE_DIR}/contract/matching_contract.wasm
fi

if [[ $LINK_BIN_ACTION =~ ^[Yy]$ ]]
then
    ln -sfn ${BIN_DIR}/med /bin/med
    ln -sfn ${BIN_DIR}/meda /bin/meda
    ln -sfn ${BIN_DIR}/meda-appd /bin/meda-appd
    ln -sfn ${BIN_DIR}/rly /bin/rly
    ln -sfn ${BIN_DIR}/rollappd /bin/rollappd 
fi

if [[ "${RUN_STAGE}" == "rollapp_after" ]]; then
    echo "# ---------------------------------------------------------------------------- #"
    echo "#     Resume mode RUN_STAGE=rollapp_after: skip genesis and continue from rollapp IBC              #"
    echo "# ---------------------------------------------------------------------------- #"
    docker compose up -d hub-node1 hub-node2 hub-node3 hub-node4 \
        da-node1 da-node2 da-node3 da-node4 da-bridge da-full da-light \
        rly-da rollapp-node1
    docker compose ps
    load_existing_addrs
    setup_after_rollapp_node1
    exit 0
fi

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


mv   ${BASE_DIR}/nodes ${BASE_DIR}/old.nodes-$(date "+%Y-%m-%d-%H-%M-%S")
mkdir -p ${BASE_DIR}/nodes

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
#sed -i'' -e "/\[mempool\]/,+8 s/v0/v1/" "$TENDERMINT_CONFIG_FILE"
#Allow multiple services on one machine on different ports
sed -i '/^allow_duplicate_ip/s/false/true/' "${ME_NODE1_HOME}"/config/config.toml
#Address-book strictness. false allows non-routable addresses (such as private IPs) as peers and non-global unicast addresses.
sed -i '/^addr_book_strict/s/true/false/' "${ME_NODE1_HOME}"/config/config.toml
sed -i 's/0.0.0.0:26656/0.0.0.0:26656/' "${ME_NODE1_HOME}"/config/config.toml
sed -i 's/127.0.0.1:26657/0.0.0.0:26657/' "${ME_NODE1_HOME}"/config/config.toml
# Block time
sed -i 's/^timeout_commit = .*/timeout_commit = "'${ME_TIMEOUT_COMMIT}'"/' "${ME_NODE1_HOME}"/config/config.toml
#CORS
sed -i 's/cors_allowed_origins.*$/cors_allowed_origins = ["*"]/' "${ME_NODE1_HOME}"/config/config.toml
sed -i 's/localhost:6060/0.0.0.0:6060/' "${ME_NODE1_HOME}"/config/config.toml
echo "Edit/config/app.toml"
sed -i'' -e '/\[api\]/,+3 s/enable *= .*/enable = true/' "${ME_NODE1_HOME}/config/app.toml"
sed -i 's/swagger = false/swagger = true/'   "${ME_NODE1_HOME}/config/app.toml"
sed -i'' -e "/\[api\]/,+9 s/address *= .*/address = \"tcp:\/\/0.0.0.0:1317\"/" "${ME_NODE1_HOME}/config/app.toml"
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "${ME_NODE1_HOME}/config/app.toml"

sed -i'' -e "/\[grpc\]/,+6 s/address *= .*/address = \"0.0.0.0:9090\"/"  "${ME_NODE1_HOME}/config/app.toml"

sed -i'' -e "/\[grpc-web\]/,+7 s/address *= .*/address = \"0.0.0.0:9091\"/" "${ME_NODE1_HOME}/config/app.toml"
sed -i 's/enable-unsafe-cors = false/enable-unsafe-cors = true/' "${ME_NODE1_HOME}/config/app.toml"

sed -i'' -e "/\[json-rpc\]/,+6 s/address *= .*/address = \"0.0.0.0:8545\"/" "${ME_NODE1_HOME}/config/app.toml"
sed -i'' -e "/\[json-rpc\]/,+9 s/^ws-address *= .*/ws-address = \"0.0.0.0:8546\"/" "${ME_NODE1_HOME}/config/app.toml"
sed -i 's/127.0.0.1:6065/0.0.0.0:6065/' "${ME_NODE1_HOME}/config/app.toml"

#Minimum gas price a node accepts. Fee = gas price times gas used. Governance can change it; validators may set different values.
sed -i'' -e 's/^minimum-gas-prices *= .*/minimum-gas-prices = "'"${HUB_MIN_GAS_PRICES}"'"/' "${ME_NODE1_HOME}/config/app.toml"
# default: the last 362880 states are kept, pruning at 10 block intervals
# nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
# everything: 2 latest states will be kept; pruning at 10 block intervals.
# custom: allow pruning options to be manually specified through 'pruning-keep-recent', and 'pruning-interval'
sed -i'' -e 's/^pruning *= .*/pruning = "nothing"/' "${ME_NODE1_HOME}/config/app.toml"

sed -i'' -e 's/^max-recv-msg-size *= .*/max-recv-msg-size = "1048576000"/' "${ME_NODE1_HOME}/config/app.toml"

echo "Edit/config/client.toml"
sed -i "s/localhost:26657/${ME_NODE1_IP}:26657/" "${ME_NODE1_HOME}"/config/client.toml
# The keyring's backend, where the keys are stored (os|file|kwallet|pass|test|memory)
echo "Keyring backend is now:${KEYRING_BACKEND_NAME}"
sed -i'' -e "s/^keyring-backend *= .*/keyring-backend = \"${KEYRING_BACKEND_NAME}\"/" "${ME_NODE1_HOME}"/config/client.toml

med config chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME}
med config keyring-backend ${KEYRING_BACKEND_NAME} --home ${ME_NODE1_HOME}
med config node "tcp://${ME_NODE1_IP}:26657" --home ${ME_NODE1_HOME}

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
    sed -i'' -e 's/bond_denom": ".*"/bond_denom": "umec"/' "${ME_GENESIS_FILE}"
    sed -i'' -e 's/mint_denom": ".*"/mint_denom": "umec"/' "${ME_GENESIS_FILE}"

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

set_authorised_deployer_account() {
  #rollapp creation authority address
  jq --arg address $1 '.app_state.rollapp.params.deployer_whitelist += [{ "address": $address }]' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}



set_consenus_params
set_gov_params
set_hub_params
set_misc_params
set_EVM_params
set_bank_denom_metadata
set_epochs_params
set_incentives_params


if [[ "$ENABLE_MONITPORING" == [Yy] ]] ;then
  sed  -i'' -e "s/prometheus = false/prometheus = true/" "${ME_TENDERMINT_CONFIG_FILE}"
fi

echo "# ---------------------------------------------------------------------------- #"
echo "# Decide whether (global_dao, meid_dao, dev_operator, airdrop, val1, user, sequencer) need to be recovered and add #"
echo "# -------------------------------------------------------------------------------------#"
if [ -z "${USER_ADDR_ME_GLOBAL_DAO}" ]; then
    echo "add global_dao"
    me_global_dao=$(add_key "med" "global_dao"  ${USER_SEQ_ME_GLOBAL_DAO} "--home $ME_NODE1_HOME")
    echo "${me_global_dao}"
    USER_ADDR_ME_GLOBAL_DAO=$(echo "${me_global_dao}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_GLOBAL_DAO_PUBKEY=$(echo "${me_global_dao}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "-------------------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_MEID_DAO}" ]; then
    echo "add meid_dao"
    me_meid_dao=$(add_key "med"  "meid_dao" ${USER_SEQ_ME_MEID_DAO} "--home $ME_NODE1_HOME")
    echo "${me_meid_dao}"
    USER_ADDR_ME_MEID_DAO=$(echo "${me_meid_dao}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_MEID_DAO_PUBKEY=$(echo "${me_meid_dao}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_DEV_OPERATOR}" ]; then
    echo "add dev_operator"
    me_dev_operator=$(add_key  "med" "dev_operator" ${USER_SEQ_ME_DEV_OPERATOR} "--home $ME_NODE1_HOME")
    echo "${me_dev_operator}"
    USER_ADDR_ME_DEV_OPERATOR=$(echo "${me_dev_operator}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_DEV_OPERATOR_PUBKEY=$(echo "${me_dev_operator}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_AIRDROP}" ]; then
    echo "add airdrop"
    me_airdrop=$(add_key  "med" "airdrop" ${USER_SEQ_ME_AIRDROP} "--home $ME_NODE1_HOME")
    echo "${me_airdrop}"
    USER_ADDR_ME_AIRDROP=$(echo "${me_airdrop}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_AIRDROP_PUBKEY=$(echo "${me_airdrop}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "-------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_VAL1}" ]; then
    echo "add val1"
    me_val1=$(add_key  "med"  "val1" ${USER_SEQ_ME_VAL1} "--home $ME_NODE1_HOME")
    echo "${me_val1}"
    USER_ADDR_ME_VAL1=$(echo "${me_val1}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_VAL1_PUBKEY=$(echo "${me_val1}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "--------------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_SEQUENCER}" ]; then
    echo "add sequencer"
    me_sequencer=$(add_key  "med"  "sequencer" ${USER_SEQ_ME_SEQUENCER} "--home $ME_NODE1_HOME")
    echo "${me_sequencer}"
    USER_ADDR_ME_SEQUENCER=$(echo "${me_sequencer}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_SEQUENCER_PUBKEY=$(echo "${me_sequencer}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
if [ -z "${USER_ADDR_ME_SEQUENCER2}" ]; then
    echo "add sequencer2"
    me_sequencer2=$(add_key  "med"  "sequencer2" ${USER_SEQ_ME_SEQUENCER2} "--home $ME_NODE1_HOME")
    echo "${me_sequencer2}"
    USER_ADDR_ME_SEQUENCER2=$(echo "${me_sequencer2}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_SEQUENCER2_PUBKEY=$(echo "${me_sequencer2}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
if [ -z "${USER_ADDR_ME_SEQUENCER3}" ]; then
    echo "add sequencer3"
    me_sequencer3=$(add_key  "med"  "sequencer3" ${USER_SEQ_ME_SEQUENCER3} "--home $ME_NODE1_HOME")
    echo "${me_sequencer3}"
    USER_ADDR_ME_SEQUENCER3=$(echo "${me_sequencer3}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_SEQUENCER3_PUBKEY=$(echo "${me_sequencer3}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi

echo "-----------------------------------------------------------------------------------------"
echo "# ---------------------------------------------------------------------------------------#"
echo "# Initialize balances add-genesis-account global_dao-"0umec"                                       #"
echo "# (meid_dao, dev_operator, airdrop,val1)-"0umec"                                        #"                            
echo "# --------------------------------------------------------------------------------------#"
med add-genesis-account global_dao   "0umec" --home ${ME_NODE1_HOME}
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
echo "#   Users allowed to edit the whitelist, super-admin (global_dao, sequencer, sequencer2,sequencer3)    #"
echo "# ---------------------------------------------------------------------------- #"

set_authorised_deployer_account "${USER_ADDR_ME_GLOBAL_DAO}"

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
    sed -i "s|^moniker = .*$|moniker = \"${NODE_NAME}\"|" "${TENDERMINT_CONF}"
    sed -i "s|^node = .*$|node = \"tcp://${NODE_IP}:26657\"|" "${CLIENT_CONFIG_FILE}"
    
}

me_init_sync_node node2 ${ME_NODE2_IP}
me_init_sync_node node3 ${ME_NODE3_IP}
me_init_sync_node node4 ${ME_NODE4_IP}

echo "# ---------------------------------------------------------------------------- #"
echo "#                              ME Read Node ID                                     #"
echo "# ---------------------------------------------------------------------------- #"
NODEID_NODE1=$(med tendermint show-node-id --home ${ME_NODES_HOME}/node1)
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
sed -i 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node2_P2P}"'"/' ${ME_NODES_HOME}/node2/config/config.toml
sed -i 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node3_P2P}"'"/' ${ME_NODES_HOME}/node3/config/config.toml
sed -i 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node4_P2P}"'"/' ${ME_NODES_HOME}/node4/config/config.toml

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

x-rollapp-template: &rollapp-template
  restart: unless-stopped
  image: ${ROLLAPP_IMAGE}
  entrypoint: ["rollappd"]
  networks:
    - ${Docker_Network_Name}
  environment:
    HUB_RPC_ADDR: ${ME_NODE1_IP}:26657
    DA_RPC_ADDR: ${DA_LIGHT_IP}:26658
    RE_INDEX_START: 1
    OTEL_SDK_DISABLED: true
  command:
    - start
    - --home
    - ${ROLLAPP_CONTAINER_HOME}
    - --minimum-gas-prices
    - "${ROLLAPP_MIN_GAS_PRICES}"

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
      - ${ME_NODES_HOME}/node1:${HUB_CONTAINER_HOME}

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
docker compose  down hub-node1
sed -i 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node1_P2P}"'"/' ${ME_NODES_HOME}/node1/config/config.toml
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
    USER_ADDR_ME_VAL2=$(echo "${val2}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_VAL2_PUBKEY=$(echo "${val2}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "--------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ME_VAL3}" ]; then
    val3=$(add_key  "med"  "val3" ${USER_SEQ_ME_VAL3} "--home $ME_NODE1_HOME")
    echo "${val3}"
    USER_ADDR_ME_VAL3=$(echo "${val3}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_VAL3_PUBKEY=$(echo "${val3}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "--------------------------------------------------------------------------"
fi    
if [ -z "${USER_ADDR_ME_VAL4}" ]; then
    val4=$(add_key  "med"  "val4" ${USER_SEQ_ME_VAL4} "--home $ME_NODE1_HOME")
    echo "${val4}"
    USER_ADDR_ME_VAL4=$(echo "${val4}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_VAL4_PUBKEY=$(echo "${val4}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#                  Wait for the ME genesis block and get the block hash                              #"
echo "# ---------------------------------------------------------------------------- #"
wait_genesis_block_hash() {
    app_name=$1
    extra_args=$2
    while true; do
        sleepTime=1

        BLOCK_HASH=$($app_name q block 1  $extra_args 2> /dev/null | jq -r .block_id.hash)

        if [ -z "${BLOCK_HASH}" ]; then
            echo "no get genesis block, sleep 2s"
            sleep 2
            continue
        fi

        echo "get genesis block hash: ${BLOCK_HASH}"
        break
    done
}
wait_genesis_block_hash med "--home ${ME_NODE1_HOME}"


echo "# ---------------------------------------------------------------------------- #"
echo "#                  Get ME Node2, Node3, Node4 pubkeys                                   #"
echo "# ---------------------------------------------------------------------------- #"
ME_Node1_Pubkey=$(med tendermint show-validator --home ${ME_NODES_HOME}/node1)
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
    echo "#            Self-hosted mode waits for height ${ME_CREAT_REGION_HEIGHT} so region creation receives enough tokens     #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "Need about target height ${ME_CREAT_REGION_HEIGHT} times block time ${ME_TIMEOUT_COMMIT} seconds"
    while true; do
        CURRENT_HEIGHT=$(med status --home ${ME_NODE1_HOME} 2>/dev/null | jq -r '.SyncInfo.latest_block_height' 2>/dev/null)
        CURRENT_HEIGHT=${CURRENT_HEIGHT//\"/}
        
        echo "Current block height: $CURRENT_HEIGHT"

        if [ "$CURRENT_HEIGHT" -ge ${ME_CREAT_REGION_HEIGHT} ] 2>/dev/null; then
            echo -e "\nBlock height has reached ${ME_CREAT_REGION_HEIGHT}. Proceeding with next steps..."
            break
        fi
        timeout_seconds=$(timeout_commit_seconds "${ME_TIMEOUT_COMMIT}")
        sleepTime=$(echo "(${ME_CREAT_REGION_HEIGHT} * ${timeout_seconds}) / 3" | bc)
        echo "Waiting for block height to reach ${ME_CREAT_REGION_HEIGHT} (next check in ${sleepTime} seconds)..."
        sleep "$sleepTime"
    done

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
echo "# ---------------------------------------------------------------------- ------ #"
echo "#               Query the region's current balance                                                 #"
echo "# ---------------------------------------------------------------------------- #"
region_data=$(med query staking regions --home ${ME_NODE1_HOME} -o json)
echo "${region_data}" | jq -r '.region[] | [.regionId, .region_treasure_addr] | @tsv' | while IFS=$'\t' read -r region_id treasure_addr; do
    balance_json=$(med query bank balances "${treasure_addr}" --home ${ME_NODE1_HOME} -o json 2>/dev/null)
    balance=$(echo "${balance_json}" | jq -r '.balances[] | select(.denom=="umec") | .amount // "0"')
    echo "$region_id current balance is ${balance:-0} umec"

done


if [[ $Interactive_Mode =~ ^[Nn]$ ]]
then
    if [[ $INSTALL_TYPE =~ ^[Yy]$ ]]
    then
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${TARGET_BALANCE}umec --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        echo "Quick-withdraw mode takes the Send DA amount from ${ME_STAKING_AMOUNT_NODE1_REGION}${TARGET_BALANCE}"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${TARGET_BALANCE}umec --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
    else
        echo "# ---------------------------------------------------------------------------- #"
        echo "# Self-hosted production params wait for height ${REGION_HIGHT_WAIT} so every region has enough to withdraw (or create a 5th validator at a chosen height) #"
        echo "# ---------------------------------------------------------------------------- #"
        while true; do
            CURRENT_HEIGHT=$(med status --home ${ME_NODE1_HOME} 2>/dev/null | jq -r '.SyncInfo.latest_block_height' 2>/dev/null)
            CURRENT_HEIGHT=${CURRENT_HEIGHT//\"/}
            
            echo "Current block height: $CURRENT_HEIGHT"
            
            if [ "$CURRENT_HEIGHT" -ge "${REGION_HIGHT_WAIT}" ] 2>/dev/null; then
                echo -e "\nBlock height has reached ${REGION_HIGHT_WAIT}. Proceeding with next steps..."
                break
            fi
            timeout_seconds=$(timeout_commit_seconds "${ME_TIMEOUT_COMMIT}")
            sleepTime=$(echo "(${REGION_HIGHT_WAIT} * ${timeout_seconds}) / 100" | bc)
            echo "Waiting for block height to reach ${REGION_HIGHT_WAIT} (next check in ${sleepTime} seconds)..."
            sleep "$sleepTime"
        done
        echo "Withdraw the Send DA val1 amount from ${ME_STAKING_AMOUNT_NODE1_REGION}${ME_SEND_DA_AMOUNT_VAL1}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_SEND_DA_AMOUNT_VAL1} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_SEND_DA_AMOUNT_VAL1} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "Check whether the tx hash is on chain: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        
        echo "Withdraw the Send IBC amount from ${ME_STAKING_AMOUNT_NODE1_REGION}${ME_SEND_IBC_AMOUNT}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_SEND_IBC_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_SEND_IBC_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "Check whether the tx hash is on chain: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"

        echo "Withdraw the reserved amount from ${ME_STAKING_AMOUNT_NODE1_REGION}${ME_RESERVE_AMOUNT}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_RESERVE_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_RESERVE_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "Check whether the tx hash is on chain: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"

        echo "Withdraw the Send DA val2 amount from ${ME_STAKING_AMOUNT_NODE2_REGION}${ME_SEND_DA_AMOUNT_VAL2}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE2_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL2} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE2_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL2} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "Check whether the tx hash is on chain: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        
        echo "Withdraw the Send DA val3 amount from ${ME_STAKING_AMOUNT_NODE3_REGION}${ME_SEND_DA_AMOUNT_VAL3}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE3_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL3} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE3_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL3} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "Check whether the tx hash is on chain: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"

        echo "Withdraw the Send DA val4 amount from ${ME_STAKING_AMOUNT_NODE4_REGION}${ME_SEND_DA_AMOUNT_VAL4}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE4_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL4} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE4_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL4} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "Check whether the tx hash is on chain: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        echo "Current ${USER_ADDR_ME_GLOBAL_DAO} account:"
        med query bank balances "${USER_ADDR_ME_GLOBAL_DAO}" --home ${ME_NODE1_HOME} -o json 2>/dev/null
    fi
fi

echo "# ---------------------------------------------------------------------- ------ #"
echo "#               Deploy ME-DA-Node1                                                 #"
echo "# ---------------------------------------------------------------------------- #"

meda-appd init "${DA_NODE1_NAME}" --chain-id="${DA_CHAIN_ID}" --home "${DA_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#   Decide whether to add users (dao, ibc, val0, val1, val2, val3, val4). No output means do not create      #"
echo "# ---------------------------------------------------------------------------- #"
if [ -z "${USER_ADDR_DA_DAO}" ]; then
    da_dao=$(add_key "meda-appd" "dao" ${da_user_index_dao} "--home ${DA_NODE1_HOME}")
    echo "${da_dao}"
    USER_ADDR_DA_DAO=$(echo "${da_dao}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_DAO_PUBKEY=$(echo "${da_dao}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
if [ -z "${USER_ADDR_DA_IBC}" ]; then
    da_ibc=$(add_key "meda-appd" "ibc" ${USER_SEQ_DA_IBC_ME} "--home ${DA_NODE1_HOME}")
    echo "${da_ibc}"
    USER_ADDR_DA_IBC=$(echo "${da_ibc}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_IBC_PUBKEY=$(echo "${da_ibc}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g") 
fi
if [ -z "${USER_ADDR_DA_VAL0}" ]; then
    da_val0=$(add_key "meda-appd"  "val0" ${da_user_index_val0} "--home ${DA_NODE1_HOME}")
    echo "${da_val0}"
    USER_ADDR_DA_VAL0=$(echo "${da_val0}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_VAL0_PUBKEY=$(echo "${da_val0}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g") 
fi
if [ -z "${USER_ADDR_DA_VAL1}" ]; then
    da_val1=$(add_key "meda-appd" "val1" ${da_user_index_val1} "--home ${DA_NODE1_HOME}")
    echo "${da_val1}"
    USER_ADDR_DA_VAL1=$(echo "${da_val1}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_VAL1_PUBKEY=$(echo "${da_val1}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
if [ -z "${USER_ADDR_DA_VAL2}" ]; then
    da_val2=$(add_key "meda-appd" "val2" ${da_user_index_val2} "--home ${DA_NODE1_HOME}")
    echo "${da_val2}"
    USER_ADDR_DA_VAL2=$(echo "${da_val2}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_VAL2_PUBKEY=$(echo "${da_val2}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
if [ -z "${USER_ADDR_DA_VAL3}" ]; then
    da_val3=$(add_key "meda-appd"  "val3" ${da_user_index_val3} "--home ${DA_NODE1_HOME}")
    echo "${da_val3}"
    USER_ADDR_DA_VAL3=$(echo "${da_val3}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_VAL3_PUBKEY=$(echo "${da_val3}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
if [ -z "${USER_ADDR_DA_VAL4}" ]; then
    da_val4=$(add_key "meda-appd"  "val4" ${da_user_index_val4} "--home ${DA_NODE1_HOME}")
    echo "${da_val4}"
    USER_ADDR_DA_VAL4=$(echo "${da_val4}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_VAL4_PUBKEY=$(echo "${da_val4}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi

#echo "create channel hub-meda, the ME ibc address is required"
if [ -z "${USER_ADDR_ME_IBC_DA}" ]; then
    da_me_ibc_da=$(add_key "meda-appd"  "me_user_ibc_da" ${USER_SEQ_ME_IBC_DA} "--home ${DA_NODE1_HOME}")
    echo "${da_me_ibc_da}"
    USER_ADDR_ME_IBC_DA=$(echo "${da_me_ibc_da}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_IBC_DA_PUBKEY=$(echo "${da_me_ibc_da}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
#echo "create channel hub-rollapp, the ME rollapp address is required"
if [ -z "${USER_ADDR_ME_IBC_ROLLAPP}" ]; then
    da_ibc_me=$(add_key "meda-appd"  "da_user_ibc_rollapp" ${USER_SEQ_ME_IBC_ROLLAPP} "--home ${DA_NODE1_HOME}")
    echo "${da_ibc_me}"
    USER_ADDR_ME_IBC_ROLLAPP=$(echo "${da_ibc_me}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_IBC_ROLLAPP_PUBKEY=$(echo "${da_ibc_me}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi


echo "# ----------------------------------------------------------------------------------------------------------- #"
echo "#  ME Add 6 free-gas-accounts: ${USER_ADDR_ME_IBC_DA}, ${USER_ADDR_ME_IBC_ROLLAPP}, ${USER_ADDR_ME_SEQUENCER}, ${USER_ADDR_ME_SEQUENCER2}, ${USER_ADDR_ME_SEQUENCER3},${USER_ADDR_ME_GLOBAL_DAO}#"
echo "# ----------------------------------------------------------------------------------------------------------- #"
seqID=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO} -o json   --home ${ME_NODE1_HOME} | jq -r .base_account.sequence)
me_add_free_gas_account_list() {
    local addrs=("$@") 
    local current_free_gas_address=$(med q dao free-gas-accounts --home ${ME_NODE1_HOME} -o json | jq -c '.addresses')
    echo  current free gas address:${current_free_gas_address}
    for addr in "${addrs[@]}"; do
        echo "Current addr:${addr}"
        if  echo "$current_free_gas_address" | jq --arg addr "${addr}" '. as $arr | $arr | index($addr) != null | not' | grep -q "true"; then
            echo "#Current sequence: $seqID, query free-gas-account:${addr} #"
            echo "Run shell: med tx dao free-gas-account "[{\"address\":\"${addr}\", \"is_free\":true}]" --from ${USER_ADDR_ME_GLOBAL_DAO} --home ${ME_NODE1_HOME} --gas 500000 --fees 200000umec --account-number ${account_number} --sequence ${seqID} --offline -y -o json | jq -r '.txhash'"
            txhash=$(med tx dao free-gas-account "[{\"address\":\"${addr}\", \"is_free\":true}]" \
            --from ${USER_ADDR_ME_GLOBAL_DAO} \
            --home ${ME_NODE1_HOME} \
            --gas 500000 \
            --fees 200000umec \
            --account-number ${account_number} \
            --sequence ${seqID} \
            --offline \
            -y \
            -o json | jq -r '.txhash')
            # Check whether the tx succeeded
            if [ -z "${txhash}" ] || [ "${txhash}" = "null" ]; then
                echo "error: failed to add free-gas account ${addr}"
                med q tx ${txhash} --home ${ME_NODE1_HOME}
                exit 1
            fi
            
            echo "tx hash: ${txhash}"
            check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        else
             echo "This addr already exists, skip"
        fi 
        seqID=$((seqID+1))
    done
    # Check that every account was added to the free-gas list
    echo "Current free-gas account list:"
    med q dao free-gas-accounts --home ${ME_NODE1_HOME}  -o json | jq -c '.addresses'
    echo "----------------------------"
}


me_add_free_gas_account_list ${USER_ADDR_ME_IBC_DA} ${USER_ADDR_ME_IBC_ROLLAPP} ${USER_ADDR_ME_SEQUENCER}  ${USER_ADDR_ME_SEQUENCER2} ${USER_ADDR_ME_SEQUENCER3} ${USER_ADDR_ME_GLOBAL_DAO}




echo "# ---------------------------------------------------------------------------- #"
echo "#      Initialize Rollapp Node1 early, step 1, so rollapp user addresses exist                    #"
echo "#      Deploy Rollapp Node1 step 1                                                  #"
echo "# ---------------------------------------------------------------------------- #"


echo "Initialize the node..."
rollappd init "${RAPP_NODE1_NAME}" --chain-id "${ROLLAPP_CHAIN_ID}" --home ${RAPP_NODE1_HOME}

echo "# ------------------------------------------------------------------------------------------- #"
echo "#      Rollapp add roluser, ibc, sync_user, dao, dev_operator, sequencer. No output means do not create         #"
echo "# ------------------------------------------------------------------------------------------- #"
if [ -z "${USER_ADDR_ROLLAPP_ADMIN}" ]; then
    echo "add roluser"
    rollapp_roluser=$(add_key  "rollappd"  "roluser" ${rollapp_user_index_roluser} "--home ${RAPP_NODE1_HOME}")
    echo "${rollapp_roluser}"
    USER_ADDR_ROLLAPP_ADMIN=$(echo "${rollapp_roluser}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ROLLAPP_ADMIN_PUBKEY=$(echo "${rollapp_roluser}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "-----------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ROLLAPP_IBC}" ]; then
    echo "add ibc "
    rollapp_ibc=$(add_key  "rollappd"  "ibc" ${rollapp_user_index_ibc} "--home ${RAPP_NODE1_HOME}")
    echo "${rollapp_ibc}"
    USER_ADDR_ROLLAPP_IBC=$(echo "${rollapp_ibc}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ROLLAPP_IBC_PUBKEY=$(echo "${rollapp_ibc}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "-----------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ROLLAPP_SYNC}" ]; then
    echo "add sync_user"
    rollapp_sync_user=$(add_key  "rollappd"  "sync_user" ${rollapp_user_index_sync_user} "--home ${RAPP_NODE1_HOME}")
    echo "${rollapp_sync_user}"
    USER_ADDR_ROLLAPP_SYNC=$(echo "${rollapp_sync_user}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ROLLAPP_SYNC_PUBKEY=$(echo "${rollapp_sync_user}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "-----------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ROLLAPP_DAO}" ]; then
    echo "add dao"
    rollapp_dao=$(add_key  "rollappd"  "dao" ${rollapp_user_index_dao} "--home ${RAPP_NODE1_HOME}")
    echo "${rollapp_dao}"
    USER_ADDR_ROLLAPP_DAO=$(echo "${rollapp_dao}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ROLLAPP_DAO_PUBKEY=$(echo "${rollapp_dao}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "-----------------------------------------------------------------------------------"
fi
if [ -z "${USER_ADDR_ROLLAPP_OPERATOR}" ]; then
    echo "add dev_operator"
    rollapp_dev_operator=$(add_key  "rollappd"  "dev_operator" ${rollapp_user_index_dev_operator} "--home ${RAPP_NODE1_HOME}")
    echo "${rollapp_dev_operator}"
    USER_ADDR_ROLLAPP_OPERATOR=$(echo "${rollapp_dev_operator}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ROLLAPP_OPERATOR_PUBKEY=$(echo "${rollapp_dev_operator}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "-----------------------------------------------------------------------------------"
fi

# rollapp Recover seq user
if [ -z "${USER_ADDR_ROLLAPP_SEQUENCER}" ]; then
    echo "Recover the sequencer account..."
    mkdir -p "${RAPP_NODE1_HOME}"/sequencer_keys/keyring-${KEYRING_BACKEND_NAME}
    rollapp_seq=$(add_key  "rollappd"  "sequencer" ${USER_SEQ_ME_SEQUENCER} "--home ${RAPP_NODE1_HOME}/sequencer_keys")
    echo "${rollapp_seq}"
    USER_ADDR_ROLLAPP_SEQUENCER=$(echo "${rollapp_seq}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ROLLAPP_SEQUENCER_PUBKEY=$(echo "${rollapp_seq}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "sequencer account recovered"
fi

# Keep querying whether the address has a balance
check_user_is_balances() {
    local binary=$1
    local user_addr=$2
    local extra_args=$3
    local sleepTime=30000

    while true; do

        if [[ $sleepTime -eq 0 ]]; then
            echo "check_user_is_balances ERROR 3000s not found balances: ${user_addr}"
            exit 1
        fi

        resp_status=$(${binary} query bank balances ${user_addr} -o json ${extra_args} | jq -r .balances[0].denom)

        if [[ ${resp_status} == "null" ]] ; then
            echo "INFO: user not balances, sleep 5s"
            sleep 5
            sleepTime=$((sleepTime - 5))
            continue
        fi

        echo "check user balances success"
        echo "${binary} query bank balances ${user_addr} ${extra_args}"
        account_status=$(${binary} query bank balances ${user_addr} ${extra_args} -o json)
        balance=$(echo "$account_status" | jq -r .balances[0].amount)
        denom=$(echo "$account_status" | jq -r .balances[0].denom)
        echo User ${user_addr} holds ${balance} ${denom}
        break

    done
}

echo "# ---------------------------------------------------------------------------- #"
echo "#                           DA genesis                                           #"
echo "# ---------------------------------------------------------------------------- #"

echo "# ---------------------------------------------------------------------------- #"
echo "# Initial funds ${USER_ADDR_DA_VAL0}:${val0_TIA_AMOUNT},${USER_ADDR_DA_DAO}:${dao_TIA_AMOUNT}  #"
echo "# ---------------------------------------------------------------------------- #"
meda-appd add-genesis-account ${USER_ADDR_DA_VAL0} $val0_TIA_AMOUNT --keyring-backend ${KEYRING_BACKEND_NAME} --home "$DA_NODE1_HOME"
meda-appd add-genesis-account ${USER_ADDR_DA_DAO} $dao_TIA_AMOUNT --keyring-backend ${KEYRING_BACKEND_NAME} --home "$DA_NODE1_HOME"

echo "# ---------------------------------------------------------------------------- #"
echo "#               genesis stake val0:${val0_STAKING_AMOUNT}                            #"
echo "# ---------------------------------------------------------------------------- #"

meda-appd gentx val0 ${val0_STAKING_AMOUNT} --chain-id ${DA_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} --gas-prices 0.0001udmec --moniker ${DA_STAKING_AMOUNT_NODE1_REGION}  --home "${DA_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               Collect genesis txs                                                   #"
echo "# ---------------------------------------------------------------------------- #"

# collect-gentxs writes node1 persistent_peers as this machine's physical IP, so fix it afterwards
meda-appd collect-gentxs --home "${DA_NODE1_HOME}"
sed -i 's|^persistent_peers = .*$|persistent_peers = ""|' "${DA_NODES_HOME}/node1/config/config.toml"

echo "# ---------------------------------------------------------------------------- #"
echo "#               Validate genesis txs                                                   #"
echo "# ---------------------------------------------------------------------------- #"
meda-appd validate-genesis --home "${DA_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               Edit the app config                                                #"
echo "# ---------------------------------------------------------------------------- #"

sed -i'' -e "/\[api\]/,+3 s/enable = false/enable = true/"  "${DA_NODE1_HOME}"/config/app.toml
sed -i'' -e "/\[api\]/,+10 s/swagger = false/swagger = true/" "${DA_NODE1_HOME}"/config/app.toml
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "${DA_NODE1_HOME}"/config/app.toml
sed -i'' -e "/\[grpc\]/,+3 s/enable = false/enable = true/" "${DA_NODE1_HOME}"/config/app.toml
sed -i'' -e "/\[grpc-web\]/,+5 s/enable = false/enable = true/" "${DA_NODE1_HOME}"/config/app.toml

echo "# ---------------------------------------------------------------------------- #"
echo "#               Edit the config file                                             #"
echo "# ---------------------------------------------------------------------------- #"
sed -i '/^laddr/s/127.0.0.1/0.0.0.0/' "${DA_NODE1_HOME}"/config/config.toml
#null: index no txs (default). kv: index txs in a key-value store (LevelDB or RocksDB). psql: index txs in PostgreSQL (extra config required). Changing null to kv enables key-value indexing so a tx can be queried by hash
sed -i '/^indexer/s/null/kv/' "${DA_NODE1_HOME}"/config/config.toml
#Allow one IP to use many ports, so different nodes can run on one machine
sed -i '/^allow_duplicate_ip/s/false/true/' "${DA_NODE1_HOME}"/config/config.toml
#Address-book strictness. false allows non-routable addresses (such as private IPs) as peers and non-global unicast addresses.
sed -i '/^addr_book_strict/s/true/false/' "${DA_NODE1_HOME}"/config/config.toml
#On a validator or full node, keep this false, because state sync and historical queries may need the data. On an archive node with little disk, true saves space
sed -i '/^discard_abci_responses/s/true/false/' "${DA_NODE1_HOME}"/config/config.toml
#is a Tendermint config that controls tx trace detail. noop: default, record nothing; local: record traces locally.
sed -i '/^trace_type/s/noop/local/'  "${DA_NODE1_HOME}"/config/config.toml
sed -i '/^trace_pull_address/s/.*/trace_pull_address = \":26661\"/' "${DA_NODE1_HOME}"/config/config.toml
#CORS
sed -ie 's/cors_allowed_origins.*$/cors_allowed_origins = ["*"]/' "${DA_NODE1_HOME}"/config/config.toml
if [[ "$ENABLE_MONITPORING" == [Yy] ]] ;then
    sed  -i'' -e "s/prometheus = false/prometheus = true/" "${DA_NODE1_HOME}"/config/config.toml
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#               Edit the client config                                                #"
echo "# ---------------------------------------------------------------------------- #"
meda-appd config chain-id ${DA_CHAIN_ID} --home "${DA_NODE1_HOME}"
meda-appd config keyring-backend ${KEYRING_BACKEND_NAME} --home "${DA_NODE1_HOME}"
meda-appd config node "tcp://${DA_NODE1_IP}:26657" --home "${DA_NODE1_HOME}"
echo "Edit/config/client.toml"
sed -i "s/localhost:26657/${DA_NODE1_IP}:26657/" "${DA_NODE1_HOME}"/config/client.toml
# The keyring's backend, where the keys are stored (os|file|kwallet|pass|test|memory)
echo "Replace the wallet name"
sed -i'' -e "s/^keyring-backend *= .*/keyring-backend = \"${KEYRING_BACKEND_NAME}\"/" "${DA_NODE1_HOME}"/config/client.toml

echo "# ---------------------------------------------------------------------------- #"
echo "#               Edit the genesis config                                            #"
echo "# ---------------------------------------------------------------------------- #"

# Override the VotingPeriod from 1 week to 1 minute
sed -i'' 's#"604800s"#"60s"#g' "${DA_GENESIS_FILE}"

# Override the genesis to use app version 1 and then upgrade to app version 2 later.
sed -i'' 's#"app_version": "2"#"app_version": "1"#g' "${DA_GENESIS_FILE}"

# Raise the IBC token activation height
sed -i "/ibc_deadline_block_height/s/100/${IBC_DEADLINE_BLOCK_HEIGHT}/" "${DA_GENESIS_FILE}"

# Write the dao address into the genesis file
tmp=$(mktemp)


sed -i "/max_validators/s/100/${MAX_VALIDATORS}/" "${DA_GENESIS_FILE}"

jq --arg addr ${USER_ADDR_DA_DAO} '.app_state.dao.dao_address = $addr' "${DA_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${DA_GENESIS_FILE}"
#minimum stake
jq --arg amount ${DA_MIN_DEPOSIT_AMOUNT} '.app_state.gov.deposit_params.min_deposit[0].amount = $amount' "${DA_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${DA_GENESIS_FILE}"
#Lower bound on validator commission, so a validator cannot undercut others unfairly
jq --arg var2 "0.000000000000000000" '.app_state.staking.params.min_commission_rate = $var2' "${DA_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${DA_GENESIS_FILE}"
#unbonding time
jq --arg var2 "${DA_UNBONDING_TIME}" '.app_state.staking.params.unbonding_time = $var2' "${DA_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${DA_GENESIS_FILE}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               Configure the other sync nodes                                             #"
echo "# ---------------------------------------------------------------------------- #"
da_init_sync_node() {
    local NODE_NAME=$1
    local NODE_HOME="${DA_NODES_HOME}/${NODE_NAME}"
    local NODE_IP=$2

    local CONFIG_DIRECTORY="$NODE_HOME/config"
    local TENDERMINT_CONF="$CONFIG_DIRECTORY/config.toml"
    local DA_APP_CONFIG_FILE="$CONFIG_DIRECTORY/app.toml"

    echo "========================================"
    echo "init node: ${NODE_NAME} - ${NODE_HOME}"


    meda-appd init ${NODE_NAME} --chain-id="${DA_CHAIN_ID}" --home ${NODE_HOME}

    cp ${DA_GENESIS_FILE} ${CONFIG_DIRECTORY}
    /bin/cp -rf ${DA_NODE1_HOME}/config/gentx ${CONFIG_DIRECTORY}/

    /bin/cp "${DA_GENESIS_FILE}"  "${CONFIG_DIRECTORY}"
    /bin/cp -rf ${DA_NODE1_HOME}/config/gentx  ${CONFIG_DIRECTORY}/
    /bin/cp -rf ${DA_NODE1_HOME}/config/app.toml  ${CONFIG_DIRECTORY}/
    /bin/cp -rf ${DA_NODE1_HOME}/config/client.toml  ${CONFIG_DIRECTORY}/
    /bin/cp -rf ${DA_NODE1_HOME}/config/config.toml ${CONFIG_DIRECTORY}/
    sed -i "s|^moniker = .*$|moniker = \"${NODE_NAME}\"|" "${TENDERMINT_CONF}"
    sed -i "s|^node = .*$|node = \"tcp://${NODE_IP}:26657\"|" "${NODE_HOME}/config/client.toml"

}
da_init_sync_node node2 ${DA_NODE2_IP}
da_init_sync_node node3 ${DA_NODE3_IP}
da_init_sync_node node4 ${DA_NODE4_IP}

echo "# ---------------------------------------------------------------------------- #"
echo "#               Get DA Node ID                                               #"
echo "# ---------------------------------------------------------------------------- #"
DA_NODEID_NODE1=$(meda-appd  tendermint show-node-id --home "${DA_NODE1_HOME}")
DA_NODEID_NODE2=$(meda-appd  tendermint show-node-id --home "${DA_NODES_HOME}"/node2)
DA_NODEID_NODE3=$(meda-appd  tendermint show-node-id --home "${DA_NODES_HOME}"/node3)
DA_NODEID_NODE4=$(meda-appd  tendermint show-node-id --home "${DA_NODES_HOME}"/node4)

echo Node1_ID:"${DA_NODEID_NODE1}"
echo Node2_ID:"${DA_NODEID_NODE2}"
echo Node3_ID:"${DA_NODEID_NODE3}"
echo Node4_ID:"${DA_NODEID_NODE4}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                         Edit Node2-Node5 persistent_peers                #"
echo "# ---------------------------------------------------------------------------- #"
echo NodeID_Node1_P2P:"${NodeID_Node2_P2P}"
NodeID_Node1_P2P="${DA_NODEID_NODE2}@${DA_NODE2_IP}:26656,${DA_NODEID_NODE3}@${DA_NODE3_IP}:26656,${DA_NODEID_NODE4}@${DA_NODE4_IP}:26656"
echo NodeID_Node2_P2P:"${NodeID_Node2_P2P}"
NodeID_Node2_P2P="${DA_NODEID_NODE1}@${DA_NODE1_IP}:26656,${DA_NODEID_NODE3}@${DA_NODE3_IP}:26656,${DA_NODEID_NODE4}@${DA_NODE4_IP}:26656"

echo NodeID_Node3_P2P:"${NodeID_Node2_P2P}"
NodeID_Node3_P2P="${DA_NODEID_NODE1}@${DA_NODE1_IP}:26656,${DA_NODEID_NODE2}@${DA_NODE2_IP}:26656,${DA_NODEID_NODE4}@${DA_NODE4_IP}:26656"

echo NodeID_Node4_P2P:"${NodeID_Node3_P2P}"
NodeID_Node4_P2P="${DA_NODEID_NODE1}@${DA_NODE1_IP}:26656,${DA_NODEID_NODE2}@${DA_NODE2_IP}:26656,${DA_NODEID_NODE3}@${DA_NODE3_IP}:26656"


echo NodeID_Node4_P2P:"${NodeID_Node4_P2P}"
sed -i 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node2_P2P}"'"/' "${DA_NODES_HOME}"/node2/config/config.toml
sed -i 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node3_P2P}"'"/' "${DA_NODES_HOME}"/node3/config/config.toml
sed -i 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node4_P2P}"'"/' "${DA_NODES_HOME}"/node4/config/config.toml


echo "# ---------------------------------------------------------------------------- #"
echo "#                          Generate docker-compose file                            #"
echo "# ---------------------------------------------------------------------------- #"
{
cat <<EOF
x-da-template: &da-template
  restart: unless-stopped
  image: ubuntu:24.04
  user: "1000:1000"
  networks:
    - ${Docker_Network_Name}
  command:
    - "meda-appd"
    - "start"
    - "--v2-upgrade-height"
    - "3"

EOF
  cat docker-compose.yml
} > temp238rdy2 && mv temp238rdy2 docker-compose.yml
cat >>docker-compose.yml<<EOF

  da-node1:
    restart: unless-stopped
    image: ubuntu:24.04
    user: "1000:1000"
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${DA_NODE1_IP}
    ports:
      - 36657:26657
      - 9190:9090
      - 2317:1317
    volumes:
      - ${BIN_DIR}/meda-appd:/bin/meda-appd
      - ${DA_NODES_HOME}/node1:/home/ubuntu/.meda-app
    command:
      - "meda-appd"
      - "start"
      - "--v2-upgrade-height"
      - "3"

  da-node2:
    <<: *da-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${DA_NODE2_IP}
    volumes:
      - ${BIN_DIR}/meda-appd:/bin/meda-appd
      - ${DA_NODES_HOME}/node2:/home/ubuntu/.meda-app

  da-node3:
    <<: *da-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${DA_NODE3_IP}
    volumes:
      - ${BIN_DIR}/meda-appd:/bin/meda-appd
      - ${DA_NODES_HOME}/node3:/home/ubuntu/.meda-app

  da-node4:
    <<: *da-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${DA_NODE4_IP}
    volumes:
      - ${BIN_DIR}/meda-appd:/bin/meda-appd
      - ${DA_NODES_HOME}/node4:/home/ubuntu/.meda-app
EOF
docker compose  up -d da-node1
echo "Wait 5s for services to finish starting"
sleep 5
docker compose up -d da-node2
echo "Wait 5s for services to finish starting"
sleep 5
docker compose up -d da-node3
echo "Wait 5s for services to finish starting"
sleep 5
docker compose  up -d da-node4
echo "Wait 5s for services to finish starting"
sleep 5
echo "# ---------------------------------------------------------------------------- #"
echo "#                  Edit node1 persistent_peers                               #"
echo "# ---------------------------------------------------------------------------- #"
docker compose down da-node1
sed -i 's|^persistent_peers = .*$|persistent_peers = "'"${NodeID_Node1_P2P}"'"|' "${DA_NODES_HOME}/node1/config/config.toml"
echo "NodeID_Node1_P2P:${NodeID_Node1_P2P}"
docker compose up -d da-node1
echo "Wait 5s for services to finish starting"
sleep 5
echo "# ---------------------------------------------------------------------------- #"
echo "#                  Current container status                                              #"
echo "# ---------------------------------------------------------------------------- #"
docker compose ps


echo "# ---------------------------------------------------------------------------- #"
echo "#                  Wait for the DA genesis block and get the block hash                          #"
echo "# ---------------------------------------------------------------------------- #"
wait_genesis_block_hash meda-appd "--home ${DA_NODE1_HOME}"
#Current wait_genesis_block_hash value $BLOCK_HASH
TrustedHash=${BLOCK_HASH}
echo "TrustedHash:${TrustedHash}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                  Set up the bridge node                                             #"
echo "# ---------------------------------------------------------------------------- #"
meda bridge init --p2p.network ${DA_CHAIN_ID} --node.store ${DA_BRIDGE_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#               Decide whether to add the bridge user                                     #"
echo "# ---------------------------------------------------------------------------- #"

if [ -z "${USER_ADDR_DA_BRIDGE}" ]; then
    meda_bridge=$(add_key  "meda-appd" "bridge" ${da_user_index_bridge}  "--home ${DA_BRIDGE_HOME}/keys")
    echo "${meda_bridge}"
    USER_ADDR_DA_BRIDGE=$(echo "${meda_bridge}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_BRIDGE_PUBKEY=$(echo "${meda_bridge}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
echo "# ---------------------------------------------------------------------------- #"
echo "#               config file                                                       #"
echo "# ---------------------------------------------------------------------------- #"


sed -i "s/DefaultKeyName .*/DefaultKeyName = \"bridge\"/" "${DA_BRIDGE_HOME}"/config.toml

sed -i "/TrustedHash/s/TrustedHash = \"\"/TrustedHash = \"${TrustedHash}\"/" "${DA_BRIDGE_HOME}"/config.toml

sed -i "s/SkipAuth = false/SkipAuth = true/" "${DA_BRIDGE_HOME}"/config.toml

sed -i'' -e "/\[Core\]/,+3 s/IP = .*/IP = \"${DA_NODE1_IP}\"/" "${DA_BRIDGE_HOME}"/config.toml

sed -i'' -e "/\[RPC\]/,+3 s/Address = \"localhost\"/Address = \"0.0.0.0\"/" "${DA_BRIDGE_HOME}"/config.toml

sed -i'' -e "/\[Gateway\]/,+3 s/Address = \"localhost\"/Address = \"0.0.0.0\"/" "${DA_BRIDGE_HOME}"/config.toml

sed -i "/ListenAddresses/s/.*/ListenAddresses = \[\"\/ip4\/0.0.0.0\/udp\/2121\/quic-v1\/webtransport\", \"\/ip4\/0.0.0.0\/udp\/2121\/quic-v1\", \"\/ip4\/0.0.0.0\/tcp\/2121\"\]/" "${DA_BRIDGE_HOME}"/config.toml


cat >>docker-compose.yml<<EOF

  da-bridge:
    restart: unless-stopped
    image: ubuntu:24.04
    user: "1000:1000"
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${DA_BRIDGE_IP}
    ports:
      - 26658:26658
    volumes:
      - ${BIN_DIR}/meda:/bin/meda
      - ${DA_BRIDGE_HOME}:/home/ubuntu/.meda-bridge-me-da
    command:
      - "meda"
      - "bridge"
      - "start"
      - "--p2p.network"
      - "me-da"

EOF

docker compose up -d da-bridge
#===================================================================
# Get the bridge node's TRUSTED_PEERS_LIST

echo ""
while true; do
  br_peers_id=$(meda p2p info  --url "${CHAIN_RPC_ADDR_DA_BRIDGE}" --token "aa" |jq -r '.result.id' 2>/dev/null)
  if [ $? -ne 0 ]; then
    echo "no get ${CHAIN_RPC_ADDR_DA_BRIDGE} peers_id, sleep 1s"
    sleep 1
    continue
  fi

  TRUSTED_PEERS_LIST="/ip4/${DA_BRIDGE_IP}/tcp/2121/p2p/${br_peers_id}"
  echo "TRUSTED_PEERS_LIST: ${TRUSTED_PEERS_LIST}"
  echo ""
  break

done

echo "# ---------------------------------------------------------------------------- #"
echo "#                  Set up the full node                                             #"
echo "# ---------------------------------------------------------------------------- #"

MEDA_FULL_INIT=$(meda full init --p2p.network ${DA_CHAIN_ID} --node.store ${DA_FULL_HOME})
echo ${MEDA_FULL_INIT}
FULL_CELES_KEY_ADDR=$(echo "${MEDA_FULL_INIT}" | grep -A 1 "ADDRESS:" | tail -n 1 | tr -d '[:space:]')
FULL_CELES_DICT=$(echo "${MEDA_FULL_INIT}" | grep "MNEMONIC" -A 1 | tail -n 1 | sed 's/^[ \t]*//;s/[ \t]*$//')

echo "FULL_CELES_KEY_ADDR: ${FULL_CELES_KEY_ADDR}"
echo "FULL_CELES_DICT: ${FULL_CELES_DICT}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               Decide whether to add the full user                                           #"
echo "# ---------------------------------------------------------------------------- #"
if [ -z "${USER_ADDR_DA_FULL}" ]; then
    meda_full=$(add_key "meda-appd"  "full" ${da_user_index_full} "--home ${DA_FULL_HOME}/keys")
    echo "${meda_full}"
    USER_ADDR_DA_FULL=$(echo "${meda_full}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_FULL_PUBKEY=$(echo "${meda_full}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi


echo "# ---------------------------------------------------------------------------- #"
echo "#               Full config file                                                    #"
echo "# ---------------------------------------------------------------------------- #"

sed -i "s/DefaultKeyName .*/DefaultKeyName = \"full\"/"   "${DA_FULL_HOME}"/config.toml

sed -i "/TrustedHash/s/TrustedHash = \"\"/TrustedHash = \"${TrustedHash}\"/"  "${DA_FULL_HOME}"/config.toml
# sed -i "s/SkipAuth = false/SkipAuth = true/"  "${DA_FULL_HOME}"/config.toml

sed -i'' -e "/\[Core\]/,+3 s/IP = .*/IP = \"${DA_NODE1_IP}\"/"  "${DA_FULL_HOME}"/config.toml
# sed -i'' -e "/\[RPC\]/,+3 s/Port = \"26658\"/Port = \"20003\"/"  "${DA_FULL_HOME}"/config.toml

sed -i "/ListenAddresses/s/.*/ListenAddresses =\[\"\/ip4\/0.0.0.0\/udp\/2121\/quic-v1\/webtransport\", \"\/ip4\/0.0.0.0\/udp\/2121\/quic-v1\", \"\/ip4\/0.0.0.0\/udp\/2121\/webrtc-direct\", \"\/ip4\/0.0.0.0\/tcp\/2121\"\]/"  "${DA_FULL_HOME}"/config.toml
# sed -i '/2121/s/2121/2123/g'  "${DA_FULL_HOME}"/config.toml

sed -i "s#TrustedPeers = \[\]#TrustedPeers = \[\"${TRUSTED_PEERS_LIST}\"\]#"  "${DA_FULL_HOME}"/config.toml


echo "# ---------------------------------------------------------------------------- #"
echo "#                  Set up the light node                                             #"
echo "# ---------------------------------------------------------------------------- #"

MEDA_LIGHT_INIT=$(meda light init --p2p.network ${DA_CHAIN_ID} --node.store ${DA_LIGHT_HOME})
echo ${MEDA_LIGHT_INIT}
LIGHT_CELES_KEY_ADDR=$(echo "${MEDA_LIGHT_INIT}" | grep -A 1 "ADDRESS:" | tail -n 1 | tr -d '[:space:]')
LIGHT_CELES_DICT=$(echo "${MEDA_LIGHT_INIT}" | grep "MNEMONIC" -A 1 | tail -n 1 | sed 's/^[ \t]*//;s/[ \t]*$//')

echo "LIGHT_CELES_KEY_ADDR: ${LIGHT_CELES_KEY_ADDR}"
echo "LIGHT_CELES_DICT: ${LIGHT_CELES_DICT}"


echo "# ---------------------------------------------------------------------------- #"
echo "#               Decide whether to add the light user                                           #"
echo "# ---------------------------------------------------------------------------- #"
if [ -z "${USER_ADDR_DA_LIGHT}" ]; then
    meda_light=$(add_key "meda-appd" "light" ${da_user_index_light} "--home ${DA_LIGHT_HOME}/keys")
    echo "${meda_light}"
    USER_ADDR_DA_LIGHT=$(echo "${meda_light}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_LIGHT_PUBKEY=$(echo "${meda_light}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
echo "# ---------------------------------------------------------------------------- #"
echo "#               Light config file                                                    #"
echo "# ---------------------------------------------------------------------------- #"

sed -i "s/DefaultKeyName .*/DefaultKeyName = \"light\"/" "${DA_LIGHT_HOME}"/config.toml

sed -i "/TrustedHash/s/TrustedHash = \"\"/TrustedHash = \"${TrustedHash}\"/" "${DA_LIGHT_HOME}"/config.toml
# Enable auth. Token command:light_token=$(meda light auth --node.store ${DA_LIGHT_HOME} write)
sed -i "s/SkipAuth = false/SkipAuth = true/" "${DA_LIGHT_HOME}"/config.toml

sed -i'' -e "/\[Core\]/,+3 s/IP = .*/IP = \"${DA_NODE1_IP}\"/" "${DA_LIGHT_HOME}"/config.toml
# sed -i'' -e "/\[RPC\]/,+3 s/Port = \"26658\"/Port = \"20002\"/" "${DA_LIGHT_HOME}"/config.toml
sed -i'' -e "/\[RPC\]/,+3 s/Address = \"localhost\"/Address = \"0.0.0.0\"/" "${DA_LIGHT_HOME}"/config.toml

sed -i "/ListenAddresses/s/.*/ListenAddresses =\[\"\/ip4\/0.0.0.0\/udp\/2121\/quic-v1\/webtransport\", \"\/ip4\/0.0.0.0\/udp\/2121\/quic-v1\", \"\/ip4\/0.0.0.0\/udp\/2121\/webrtc-direct\", \"\/ip4\/0.0.0.0\/tcp\/2121\"\]/" "${DA_LIGHT_HOME}"/config.toml

# sed -i '/2121/s/2121/2122/g' "${DA_LIGHT_HOME}"/config.toml

sed -i "s#TrustedPeers = \[\]#TrustedPeers = \[\"${TRUSTED_PEERS_LIST}\"\]#" "${DA_LIGHT_HOME}"/config.toml


#===================================================================
echo "# ---------------------------------------------------------------------------- #"
echo "#               Configure da-full and da-light docker-compose                            #"
echo "# ---------------------------------------------------------------------------- #"
cat >>docker-compose.yml<<EOF

  da-full:
    restart: unless-stopped
    image: ubuntu:24.04
    user: "1000:1000"
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${DA_FULL_IP}
    volumes:
      - ${BIN_DIR}/meda:/bin/meda
      - ${DA_FULL_HOME}:/home/ubuntu/.meda-full-me-da
    command:
      - "meda"
      - "full"
      - "start"
      - "--p2p.network"
      - "me-da"

  da-light:
    restart: unless-stopped
    image: ubuntu:24.04
    user: "1000:1000"
    command:
      - "meda"
      - "light"
      - "start"
      - "--p2p.network"
      - "me-da"
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${DA_LIGHT_IP}
    ports:
      - 36658:26658
    volumes:
      - ${BIN_DIR}/meda:/bin/meda
      - ${DA_LIGHT_HOME}:/home/ubuntu/.meda-light-me-da

EOF

#===================================================================
docker compose up -d da-full
docker compose up -d da-light


echo "Wait 60s, otherwise the later free-whitelist step can fail"
sleep 60

echo "# ---------------------------------------------------------------------------- #"
echo "#               Set DA 10 fee-whitelist                                        #"
echo "# ---------------------------------------------------------------------------- #"
echo "# DA_IBC:${USER_ADDR_DA_IBC}                                                   #"
echo "# ME_SEQUENCER:${USER_ADDR_ME_SEQUENCER}                                       #"
echo "# ME_SEQUENCER2:${USER_ADDR_ME_SEQUENCER2}                                     #"
echo "# ME_SEQUENCER3:${USER_ADDR_ME_SEQUENCER3}                                     #"
echo "# DA_LIGHT:${USER_ADDR_DA_LIGHT}                                               #"
echo "# DA_VAL1:${USER_ADDR_DA_VAL1}                                                 #"
echo "# DA_VAL2:${USER_ADDR_DA_VAL2}                                                 #"
echo "# DA_VAL3:${USER_ADDR_DA_VAL3}                                                 #"
echo "# DA_VAL4:${USER_ADDR_DA_VAL4}                                                 #"
echo "# ---------------------------------------------------------------------------- #"



add_addr_to_da_whitelist() {
    addrs=("$@") 
    seqID=$(meda-appd q auth account ${USER_ADDR_DA_DAO} -o json --home ${DA_NODE1_HOME} | jq -r .sequence)
    account_number=$(meda-appd q auth account ${USER_ADDR_DA_DAO} -o json --home ${DA_NODE1_HOME} | jq -r .account_number)

    for addr in "${addrs[@]}"; do

        echo "Run shell: meda-appd tx dao add-to-fee-whitelist ${addr} --from dao --chain-id ${DA_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} -a $account_number --sequence ${seqID} --offline -y -o json --home ${DA_NODE1_HOME} | jq -r '.txhash'"
        txhash=$(meda-appd tx dao add-to-fee-whitelist ${addr} --from dao --chain-id ${DA_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} -a $account_number --sequence ${seqID} --offline -y -o json --home ${DA_NODE1_HOME} | jq -r '.txhash')
        seqID=$((seqID+1))
        check_tx_status meda-appd $txhash "--home ${DA_NODE1_HOME}"
    done
}

add_addr_to_da_whitelist ${USER_ADDR_DA_IBC} ${USER_ADDR_ME_SEQUENCER} ${USER_ADDR_ME_SEQUENCER2} ${USER_ADDR_ME_SEQUENCER3} ${USER_ADDR_DA_LIGHT} ${USER_ADDR_DA_VAL1} ${USER_ADDR_DA_VAL2} ${USER_ADDR_DA_VAL3} ${USER_ADDR_DA_VAL4}


echo "# ---------------------------------------------------------------------------- #"
echo "#               Query the current whitelist                                                #"
echo "# ---------------------------------------------------------------------------- #"
meda-appd q dao fee-white-list --home ${DA_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#               Initialize rly-da                                                    #"
echo "# ---------------------------------------------------------------------------- #"

rly config init --home ${RLY_DA_NODE_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#               Set channel config me-meda                                          #"
echo "# ---------------------------------------------------------------------------- #"
cat > ${RLY_DA_NODE_HOME}/config/config.yaml <<EOF
global:
    api-listen-addr: :5183
    timeout: 10s
    memo: ""
    light-cache-size: 20
    log-level: info
    ics20-memo-limit: 0
    max-receiver-size: 150
chains:
    ${ME_CHAIN_ID}:
        type: cosmos
        value:
            key-directory: ${RLY_DA_NODE_HOME}/keys/${ME_CHAIN_ID}
            key: me_user_ibc_da_name
            chain-id: ${ME_CHAIN_ID}
            rpc-addr: ${CHAIN_RPC_ADDR_ME}
            account-prefix: me
            keyring-backend: ${KEYRING_BACKEND_NAME}
            gas-adjustment: 1.5
            gas-prices: 0.02umec
            min-gas-amount: 0
            max-gas-amount: 0
            debug: true
            timeout: 10s
            block-timeout: ""
            output-format: json
            sign-mode: direct
            extra-codecs: []
            coin-type: null
            signing-algorithm: ""
            broadcast-mode: batch
            min-loop-duration: 0s
            extension-options: []
            feegrants: null
    ${DA_CHAIN_ID}:
        type: cosmos
        value:
            key-directory: ${RLY_DA_NODE_HOME}/keys/${DA_CHAIN_ID}
            key: da_user_ibc_me_name
            chain-id: ${DA_CHAIN_ID}
            rpc-addr: ${CHAIN_RPC_ADDR_DA}
            account-prefix: me
            keyring-backend: ${KEYRING_BACKEND_NAME}
            gas-adjustment: 1.2
            gas-prices: 0.1udmec
            min-gas-amount: 0
            max-gas-amount: 0
            debug: true
            timeout: 10s
            block-timeout: ""
            output-format: json
            sign-mode: direct
            extra-codecs: []
            coin-type: null
            signing-algorithm: ""
            broadcast-mode: batch
            min-loop-duration: 0s
            extension-options: []
            feegrants: null
paths:
    hub-meda:
        src:
            chain-id: ${ME_CHAIN_ID}
        dst:
            chain-id: ${DA_CHAIN_ID}
        src-channel-filter:
            rule: ""
            channel-list: []
EOF

echo "# ---------------------------------------------------------------------------- #"
echo "#               Whether to sync me_user_ibc_da_name and da_user_ibc_me_name keys          #"
echo "# ---------------------------------------------------------------------------- #"
if [[ "$Sync_Key" == [Yy] ]]; then
    mkdir -p ${RLY_DA_NODE_HOME}/keys/${ME_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    mkdir -p ${RLY_DA_NODE_HOME}/keys/${DA_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    echo "Handle me_user_ibc_da_name, same as USER_ADDR_ME_IBC_DA-me_user_ibc_da"
    me_user_ibc_da_name=$(add_key "meda-appd" "me_user_ibc_da_name" ${USER_SEQ_ME_IBC_DA} "--home ${RLY_DA_NODE_HOME}/tmpxx20")
    echo "${me_user_ibc_da_name}"
    USER_ME_IBC_DA_NAME=$(echo "${me_user_ibc_da_name}" | grep -oP 'address: \K[^\s]+')
    USER_ME_IBC_DA_NAME_PUBKEY=$(echo "${me_user_ibc_da_name}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    cp -rf ${RLY_DA_NODE_HOME}/tmpxx20/keyring-${KEYRING_BACKEND_NAME}/*  ${RLY_DA_NODE_HOME}/keys/${ME_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    rm -rf ${RLY_DA_NODE_HOME}/tmpxx20

    echo "Handle da_user_ibc_me_name, same as USER_ADDR_DA_IBC-ibc"
    da_user_ibc_me_name=$(add_key  "meda-appd"  "da_user_ibc_me_name" ${USER_SEQ_DA_IBC_ME} "--home ${RLY_DA_NODE_HOME}/tmpxx20")
    echo "${da_user_ibc_me_name}"
    USER_DA_IBC_ME_NAME=$(echo "${da_user_ibc_me_name}" | grep -oP 'address: \K[^\s]+')
    USER_DA_IBC_ME_NAME_PUBKEY=$(echo "${da_user_ibc_me_name}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    mv  ${RLY_DA_NODE_HOME}/tmpxx20/keyring-${KEYRING_BACKEND_NAME}/*   ${RLY_DA_NODE_HOME}/keys/${DA_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    rm -rf ${RLY_DA_NODE_HOME}/tmpxx20
fi
#===================================================================


echo "# ---------------------------------------------------------------------------- #"
echo "#              create channel hub-meda                                       #"
echo "# ---------------------------------------------------------------------------- #"
#error   Error sending messages  {"path_name": "hub-meda", "src_chain_id": "mechain_400-1", "dst_chain_id": "me-da", "src_client_id": "07-tendermint-0", "dst_client_id": "07-tendermint-0", "error": "rpc error: code = Unknown desc = rpc error: code = Unknown desc = failed to execute message; message index: 1: connection handshake open try failed: consensus height is greater than or equal to the current block height (0-17 >= 0-17): invalid height [cosmos/ibc-go/v6@v6.2.2/modules/core/03-connection/keeper/handshake.go:79] With gas wanted: '18446744073709551615' and gas used: '128736' : unknown request"}
#rly tx link The rly tx link sync period must be above height 17, otherwise IBC logs a similar error. It does not matter: sync waits until the required height, so the error can be ignored
echo "Wait 60 seconds"
sleep 60
echo rly tx link hub-meda --home ${RLY_DA_NODE_HOME} --override 
rly tx link hub-meda --home ${RLY_DA_NODE_HOME} --override

echo "# ---------------------------------------------------------------------------- #"
echo "#               Set channel timeout                                                     #"
echo "# ---------------------------------------------------------------------------- #"
echo rly transact client ${ME_CHAIN_ID} ${DA_CHAIN_ID} hub-meda --client-tp ${CLIENT_TP} --home ${RLY_DA_NODE_HOME}
rly transact client ${ME_CHAIN_ID} ${DA_CHAIN_ID} hub-meda --client-tp ${CLIENT_TP} --home ${RLY_DA_NODE_HOME}
echo "# ---------------------------------------------------------------------------- #"
echo "#               Deploy the relayer                                                    #"
echo "# ---------------------------------------------------------------------------- #"
cat >> docker-compose.yml<<EOF

  rly-da:
    restart: unless-stopped
    image: ubuntu:24.04
    user: "1000:1000"
    networks:
      - ${Docker_Network_Name}
    volumes:
      - ${BIN_DIR}/rly:/bin/rly
      - ${RLY_DA_NODE_HOME}:/home/ubuntu/.relayer
    command:
      - "rly"
      - "start"
      - "hub-meda"
      - "--time-threshold"
      - "${RLY_TIME_THRESHOLD}"
EOF

docker compose up -d rly-da

echo "# ---------------------------------------------------------------------------- #"
echo "#      rly-da check                                                                  #"
echo "# ---------------------------------------------------------------------------- #"
echo "sleep 5s Wait for the relayer to finish starting"
sleep 5
echo Check ${ME_CHAIN_ID} client status
rly q clients ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME}
echo Check ${DA_CHAIN_ID} client status
rly q clients ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME}

echo Check connection status
rly q connections ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME}


echo "Check ${ME_CHAIN_ID} channel status shell:rly q channels ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME}"
rly q channels ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME}
echo "Check ${DA_CHAIN_ID} channel status shell:rly q channels ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME}"
rly q channels ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME}
echo "Get ${ME_CHAIN_ID} channel ID shell:rly q channels ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME} 2>/dev/null | jq -r 'select(.state == "STATE_OPEN") | .channel_id'"
ME_CHANNEL_ID=$(rly q channels ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME} 2>/dev/null | jq -r 'select(.state == "STATE_OPEN") | .channel_id')
echo "channel_id: ${ME_CHANNEL_ID}"
echo "Get ${DA_CHAIN_ID} channel ID shell:rly q channels ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME} 2>/dev/null | jq -r 'select(.state == "STATE_OPEN") | .channel_id'"
DA_CHANNEL_ID=$(rly q channels ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME} 2>/dev/null | jq -r 'select(.state == "STATE_OPEN") | .channel_id')
echo "channel_id: ${DA_CHANNEL_ID}"


# Query the original stake
echo "Query the original stake shell:meda-appd query staking validators --home ${DA_NODE1_HOME}"
meda-appd query staking validators --home ${DA_NODE1_HOME}
echo "-------------------------------------------"

echo "# ---------------------------------------------------------------------------- #"
echo "#                Create DA validators                                                   #"
echo "# ---------------------------------------------------------------------------- #"
DaNode1Pubkey=$(meda-appd tendermint show-validator --home ${DA_NODES_HOME}/node1)
DaNode2Pubkey=$(meda-appd tendermint show-validator --home ${DA_NODES_HOME}/node2)
DaNode3Pubkey=$(meda-appd tendermint show-validator --home ${DA_NODES_HOME}/node3)
DaNode4Pubkey=$(meda-appd tendermint show-validator --home ${DA_NODES_HOME}/node4)


if [[ ! $Interactive_Mode =~ ^[Yy]$ ]]
then
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                         Non-custody mode                                            #"
    echo "#     ME gold_dao user IBC-transfers to DA val1, val2, val3, and val4 for later validator stake              #"
    echo "# ---------------------------------------------------------------------------- #"
    #If there is no initial balance, drop .base_account before reading the value
    account_number=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO} -o json --home ${ME_NODE1_HOME} | jq -r .base_account.account_number)
    seqID=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO} -o json --home ${ME_NODE1_HOME} | jq -r .base_account.sequence)
    me_ibc_transfer_list() {
        local addr=$1
        local send_amount=$2
        echo "Handle ibc: me gold_dao -> da: ${addr} ${send_amount}"
        echo "Running: med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${addr} ${send_amount} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=100000umec  --chain-id ${ME_CHAIN_ID}  --keyring-backend ${KEYRING_BACKEND_NAME} --home $ME_NODE1_HOME -a $account_number --sequence ${seqID} --offline -y -o json | jq -r '.txhash' "
        txhash=$(med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${addr} ${send_amount} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=100000umec  --chain-id ${ME_CHAIN_ID}  --keyring-backend ${KEYRING_BACKEND_NAME} --home $ME_NODE1_HOME -a $account_number --sequence ${seqID} --offline -y -o json | jq -r '.txhash')
        echo "tx hash: ${txhash}"
        echo "----------------------------"
        seqID=$((seqID+1))
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    }
    me_region_transfer_list() {
        local addr=$1
        local region=$2
        local send_amount=$3
        echo "Handle ibc: me region -> da: ${addr} ${send_amount}"
        echo "Running: med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}   $addr $send_amount --from global_dao --home $ME_NODE1_HOME -y -o json | jq -r '.txhash' "
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  $addr $send_amount --from global_dao --home $ME_NODE1_HOME -y -o json | jq -r '.txhash')
        echo "tx hash: ${txhash}"
        echo "----------------------------"
        seqID=$((seqID+1))
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    }
    da_create_validator(){
        local Node_Region_Pubkey=$1
        local DA_STAKING_AMOUNT=$2
        local USER_ADDR=$3
        local REGION=$4
        local Account_User_Number=$(meda-appd query account ${USER_ADDR} --home ${DA_NODES_HOME}/node1 -o json | jq -r '.account_number')
        local SEQUENCE=$(meda-appd query account ${USER_ADDR} --home ${DA_NODES_HOME}/node1 -o json | jq -r '.sequence')
        echo "shell:meda-appd tx staking create-validator --pubkey "${Node_Region_Pubkey}" --moniker ${REGION} --amount ${DA_STAKING_AMOUNT} --from ${USER_ADDR} --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --min-self-delegation 10 --account-number $Account_User_Number  --sequence $SEQUENCE  --home ${DA_NODES_HOME}/node1 --keyring-backend=${KEYRING_BACKEND_NAME}  -y -o json | jq -r '.txhash'"
        txhash=$(meda-appd tx staking create-validator --pubkey "${Node_Region_Pubkey}" --moniker ${REGION} \
            --amount ${DA_STAKING_AMOUNT} \
            --from ${USER_ADDR} \
            --commission-rate="0.10" \
            --commission-max-rate="0.20" \
            --commission-max-change-rate="0.01" \
            --min-self-delegation 10 \
            --account-number $Account_User_Number \
            --sequence $SEQUENCE \
            --home ${DA_NODES_HOME}/node1 \
            --keyring-backend=${KEYRING_BACKEND_NAME}  -y -o json | jq -r '.txhash')
        check_tx_status meda-appd ${txhash} "--home ${DA_NODES_HOME}/node1"
    }
    da_delegate_validator(){
        local REGION=$1
        local DA_STAKING_AMOUNT=$2
        local USER_ADDR=$3
        local Account_User_Number=$(meda-appd query account ${USER_ADDR} --home ${DA_NODES_HOME}/node1 -o json | jq -r '.account_number')
        local SEQUENCE=$(meda-appd query account ${USER_ADDR} --home ${DA_NODES_HOME}/node1 -o json | jq -r '.sequence')
        operator_address=$(meda-appd query staking validators --home "${DA_NODES_HOME}/node1" -o json | jq -r --arg region "$REGION" '.validators[] | select(.description.moniker==$region) | .operator_address')
        echo "shell:meda-appd tx staking delegate ${operator_address} ${DA_STAKING_AMOUNT}  --from ${USER_ADDR}  --account-number ${Account_User_Number} --sequence ${SEQUENCE} --home ${DA_NODES_HOME}/node1 --keyring-backend=${KEYRING_BACKEND_NAME}  -y -o json"
        txhash=$(meda-appd tx staking delegate ${operator_address} ${DA_STAKING_AMOUNT}  --from ${USER_ADDR}  --account-number ${Account_User_Number} --sequence ${SEQUENCE} --home ${DA_NODES_HOME}/node1 --keyring-backend=${KEYRING_BACKEND_NAME}  -y -o json| jq -r '.txhash')
        check_tx_status meda-appd ${txhash} "--home ${DA_NODES_HOME}/node1"
    }
    echo "--------------------------------------------------------------------------------------"
    echo "Send to ${USER_ADDR_DA_VAL1}"
    me_ibc_transfer_list ${USER_ADDR_DA_VAL1}   ${ME_SEND_DA_AMOUNT_VAL1}
    echo "Wait 10s"
    sleep 10
    echo "da-val1 Add IBC-token stake. The initial native genesis stake ${val0_STAKING_AMOUNT} is burned when DA height reaches ${IBC_DEADLINE_BLOCK_HEIGHT} and stake becomes 0. Without this IBC-token stake the consensus node unbonds"
    echo "During the unbonding period, da_delegate_validator can still restore the consensus node by adding stake. After unbonding period ${DA_UNBONDING_SECONDS}, stake can no longer be added and the node must be recreated with da_create_validator"
    da_delegate_validator ${DA_STAKING_AMOUNT_NODE1_REGION} ${DA_STAKING_AMOUNT_VAL1} ${USER_ADDR_DA_VAL1}
    echo "After overriding relayer rly, wait 5s and confirm hub-meda sync"
    echo "shell:rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME}"
    rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME} 
    echo "Wait 5s, about one block time"
    sleep 5
    echo "--------------------------------------------------------------------------------------"
    echo "Send to ${USER_ADDR_DA_VAL2}"
    me_ibc_transfer_list ${USER_ADDR_DA_VAL2}  ${ME_SEND_DA_AMOUNT_VAL2}
    echo "Wait 10s"
    sleep 10
    echo "Create validator da-val2"
    da_create_validator   ${DaNode2Pubkey}  ${DA_STAKING_AMOUNT_VAL2}  ${USER_ADDR_DA_VAL2} ${DA_STAKING_AMOUNT_NODE2_REGION}

    echo "After overriding relayer rly, wait 5s and confirm hub-meda sync"
    echo "shell:rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME}"
    rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME} 
    echo "Wait 5s, about one block time" 
    sleep 5
    echo "--------------------------------------------------------------------------------------"
    echo "Send to ${USER_ADDR_DA_VAL3}"
    me_ibc_transfer_list ${USER_ADDR_DA_VAL3}  ${ME_SEND_DA_AMOUNT_VAL3}
    echo "Wait 10s"
    sleep 10
    echo "Create validator da-val3"
    da_create_validator   ${DaNode3Pubkey}  ${DA_STAKING_AMOUNT_VAL3}  ${USER_ADDR_DA_VAL3} ${DA_STAKING_AMOUNT_NODE3_REGION}

    echo "After overriding relayer rly, wait 5s and confirm hub-meda sync"
    echo "shell:rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME}"
    rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME} 
    echo "Wait 5s, about one block time" 
    sleep 5
    echo "--------------------------------------------------------------------------------------"
    echo "Send to ${USER_ADDR_DA_VAL4}"
    me_ibc_transfer_list ${USER_ADDR_DA_VAL4}  ${ME_SEND_DA_AMOUNT_VAL4}
    echo "Wait 10s"
    sleep 10
    echo "Create validator da-val4"
    da_create_validator   ${DaNode4Pubkey}  ${DA_STAKING_AMOUNT_VAL4}  ${USER_ADDR_DA_VAL4} ${DA_STAKING_AMOUNT_NODE4_REGION}

    echo "After overriding relayer rly, wait 5s and confirm hub-meda sync"
    echo "shell:rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME}"
    rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME} 
    echo "Wait 5s, about one block time" 
    sleep 5
else
     echo "Press Enter after custody finishes the transfer"
     read -p "Do you want to continue? (Y/n) " 
fi
echo "Current consensus nodes" 
meda-appd query staking validators --home ${DA_NODE1_HOME} -o json | jq '.validators[].description.moniker'
# Query stake
echo "-------------------------------------------"
echo "meda new staking validators data"
meda-appd query staking validators -o json  --home ${DA_NODES_HOME}/node1 | jq .validators[].tokens
echo "-------------------------------------------"

if [[  ${Interactive_Mode} =~ ^[Yy]$ ]]
then
     echo "Custody mode: check the output above"
     echo "Press Enter if this looks correct"
     read -p "Do you want to continue? (Y/n)"
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#        DA deploy finished                                                             #"
echo "# ---------------------------------------------------------------------------- #"

echo "# ---------------------------------------------------------------------------- #"
echo "#      Deploy Rollapp Node1 step 2                                               #"
echo "# ---------------------------------------------------------------------------- #"

ROLLAPP_TENDERMINT_CONFIG_FILE="${RAPP_NODE1_HOME}/config/config.toml"
ROLLAPP_CONFIG_FILE="${RAPP_NODE1_HOME}/config/app.toml"
ROLLAPP_DYMINT_FILE="${RAPP_NODE1_HOME}/config/dymint.toml"
ROLLAPP_GENESIS_FILE="${RAPP_NODE1_HOME}/config/genesis.json"
ROLLAPP_CLIENT_FILE="${RAPP_NODE1_HOME}/config/client.toml"


echo "Configure node params..."
rollappd config chain-id ${ROLLAPP_CHAIN_ID} --home ${RAPP_NODE1_HOME}
rollappd config keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}
rollappd config node "tcp://${RAPP_NODE1_IP}:26657" --home ${RAPP_NODE1_HOME}

# sync me-hub key name
echo 'sync-hub-key-name = "sync_user"' >> ${ROLLAPP_CLIENT_FILE}

echo "# ---------------------------------------------------------------------------- #"
echo "#               Edit the app config                                                #"
echo "# ---------------------------------------------------------------------------- #"

sed -i'' -e "/\[api\]/,+3 s/enable = false/enable = true/"  "${ROLLAPP_CONFIG_FILE}"
sed -i'' -e "/\[api\]/,+10 s/swagger = false/swagger = true/" "${ROLLAPP_CONFIG_FILE}"
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "${ROLLAPP_CONFIG_FILE}"
sed -i'' -e "/\[grpc\]/,+3 s/enable = false/enable = true/"  "${ROLLAPP_CONFIG_FILE}"
sed -i'' -e "/\[grpc-web\]/,+5 s/enable = false/enable = true/"  "${ROLLAPP_CONFIG_FILE}"
sed -i'' -e 's/^pruning *= .*/pruning = "nothing"/' "${ROLLAPP_CONFIG_FILE}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               Edit the config file                                             #"
echo "# ---------------------------------------------------------------------------- #"
sed -i '/^laddr/s/127.0.0.1/0.0.0.0/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
sed -i '/^addr_book_strict/s/true/false/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
sed -i 's/cors_allowed_origins.*$/cors_allowed_origins = ["*"]/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
if [[ "$ENABLE_MONITPORING" == [Yy] ]] ;then
    sed  -i'' -e "s/prometheus = false/prometheus = true/" "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
fi


echo "Configure the genesis file genesis..."
tmp=$(mktemp)

jq --arg denom ${DENOM} '.app_state.mint.params.mint_denom = $denom' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
echo "Updated mint params"
jq --arg denom ${DENOM} '.app_state.staking.params.bond_denom = $denom' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
echo "Updated staking params"
jq --arg denom ${DENOM} '.app_state.gov.deposit_params.min_deposit[0].denom = $denom' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
jq --arg amount ${ROLLAPP_MIN_DEPOSIT_AMOUNT} '.app_state.gov.deposit_params.min_deposit[0].amount = $amount' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
# rollapp uses cosmos-sdk v0.46. gov GenesisState has only deposit_params/voting_params, not params
jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.deposit_params.max_deposit_period = $period' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
jq --arg period "${VOTING_PERIOD}" '.app_state.gov.voting_params.voting_period = $period' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
jq 'del(.app_state.gov.params)' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
echo "Updated gov params voting_period=${VOTING_PERIOD} max_deposit_period=${MAX_DEPOSIT_PERIOD} (sdk v0.46 fields only)"

echo "Genesis file configured"

#echo "Validate the genesis file..."
#jq '.app_state.mint.params.mint_denom' "${ROLLAPP_GENESIS_FILE}"
#jq '.app_state.staking.params.bond_denom' "${ROLLAPP_GENESIS_FILE}"
#jq '.app_state.gov.deposit_params.min_deposit[0].denom' "${ROLLAPP_GENESIS_FILE}"

#===================================================================

# Export a module pubkey - node pubkey
rollappd dymint show-sequencer --home ${RAPP_NODE1_HOME} > ${RAPP_NODE1_HOME}/sequencer.info

echo "Validate the node pubkey..."
cat ${RAPP_NODE1_HOME}/sequencer.info
echo "Node pubkey configured"
echo  sync me-hub key address
echo "sync-hub-key-address = \"${USER_ADDR_ROLLAPP_SYNC}\"" >> ${ROLLAPP_CLIENT_FILE}

echo "# ---------------------------------------------------------------------------- #"
echo "#  Initialize roluser, dao, and dev_operator balances                                         #"
echo "# ---------------------------------------------------------------------------- #"
rollappd add-genesis-account roluser "$rollapp_user_roluser_TOKEN_AMOUNT" --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}
rollappd add-genesis-account dao "$rollapp_user_dao_TOKEN_AMOUNT" --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}
rollappd add-genesis-account dev_operator "$rollapp_user_dev_operator_TOKEN_AMOUNT" --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}


rollappd gentx_seq --pubkey "$(rollappd dymint show-sequencer --home ${RAPP_NODE1_HOME})" --from roluser --home ${RAPP_NODE1_HOME}

rollappd gentx roluser "$rollapp_user_roluser_NODE_STAKING_AMOUNT" --chain-id "$ROLLAPP_CHAIN_ID" --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "# Add 5 fee whitelist entries                                                             #"
echo "# rollapp_dao:${USER_ADDR_ROLLAPP_DAO}                                         #"
echo "# rollapp_operator:${USER_ADDR_ROLLAPP_OPERATOR}                               #"
echo "# rollapp_ibc:${USER_ADDR_ROLLAPP_IBC}                                         #"
echo "# rollapp_sync:${USER_ADDR_ROLLAPP_SYNC}                                       #"
echo "# rollapp_admin:${USER_ADDR_ROLLAPP_ADMIN}                                     #"
echo "# ---------------------------------------------------------------------------- #"
jq --arg user_addr ${USER_ADDR_ROLLAPP_DAO} '.app_state.region.relayerList[0].address = $user_addr' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"

jq --arg user_addr ${USER_ADDR_ROLLAPP_OPERATOR} '.app_state.region.relayerList[1].address = $user_addr' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"

jq --arg user_addr ${USER_ADDR_ROLLAPP_IBC} '.app_state.region.relayerList[2].address = $user_addr' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"

jq --arg user_addr ${USER_ADDR_ROLLAPP_SYNC} '.app_state.region.relayerList[3].address = $user_addr' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"

jq --arg user_addr ${USER_ADDR_ROLLAPP_ADMIN} '.app_state.region.relayerList[4].address = $user_addr' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"

echo "# 5 fee whitelist entries added, shown below #"
jq '.app_state.region.relayerList' "${ROLLAPP_GENESIS_FILE}"
echo "# Collect the genesis tx file #"
rollappd collect-gentxs --home ${RAPP_NODE1_HOME}

echo "# Validate the genesis tx #"
rollappd validate-genesis --home ${RAPP_NODE1_HOME} 


echo "# ---------------------------------------------------------------------------- #"
echo "#      Edit the DYMINT_FILE                                                    #"
echo "# ---------------------------------------------------------------------------- #"
#settlement
sed -i "s/settlement_layer = \"mock\"/settlement_layer = \"me-hub\"/" ${ROLLAPP_DYMINT_FILE}

#sed -i "/settlement_gas_limit/s/= .*/= 10000000/" ${DYMINT_FILE}

# hub
sed -i "/settlement_node_address/s/127.0.0.1:36657/${HUB_RPC_ADDR}/" ${ROLLAPP_DYMINT_FILE}

sed -i '/^settlement_gas_prices/s/1000000000/0.02umec/' ${ROLLAPP_DYMINT_FILE}

sed -i '/^settlement_gas_prices/s/1000000000/0.02umec/' ${ROLLAPP_DYMINT_FILE}

sed -i '/max_supported_batch_skew/s/20/2000/' ${ROLLAPP_DYMINT_FILE}
sed -i '/retry_attempts/s/10/20/' ${ROLLAPP_DYMINT_FILE}

echo "# DYMINT_FILE file updated #"
echo "# ---------------------------------------------------------------------------- #"
# Read JSON data
# base_url = da light addr
# timeout param in the DA (data availability) client sets the request timeout in nanoseconds. 50 seconds is 50,000,000,000 nanoseconds
json_data=$(cat <<EOF
{
    "base_url": "http://${DA_RPC_ADDR}",
    "timeout": 50000000000,
    "gas_prices": 0.1,
    "auth_token": "TOKEN123",
    "backoff": {
        "initial_delay": 6000000000,
        "max_delay": 6000000000,
        "growth_factor": 2 
    },
    "retry_attempts": 4,
    "retry_delay": 3000000000
}
EOF
)

escaped_json=$(echo "${json_data}" | jq -c . | sed 's/"/\\"/g')
echo "escaped_json: ${escaped_json}"
#===================================================================
# Dymint

sed -i '/max_idle_time/s/=.*/= "5s"/' ${ROLLAPP_DYMINT_FILE}
#Max time allowed to build a block proof, so proof generation cannot run forever and stall block production. In production, raise it if proofs often time out
sed -i '/max_proof_time/s/=.*/= "5s"/' ${ROLLAPP_DYMINT_FILE}
#Timeout for submitting a tx batch to the data-availability layer, so a network problem cannot wait forever
sed -i "/batch_submit_max_time/s/=.*/= \"${BATCH_SUBMIT_MAX_TIME}\"/" ${ROLLAPP_DYMINT_FILE}
#max batch size
sed -i '/^block_batch_max_size_bytes/s/500000/1874272/' ${ROLLAPP_DYMINT_FILE}
#Which data-availability implementation the rollup uses, so tx data is available to every participant. "mock": simulation, for tests only. "me-da": the custom DA layer named "me-da", paired with da_config.
sed -i "s/da_layer = \"mock\"/da_layer = \"me-da\"/" ${ROLLAPP_DYMINT_FILE}
# Use sed to replace the da_config field in the TOML file
sed -i "s|da_config = \"\"|da_config = \'$escaped_json\'|" ${ROLLAPP_DYMINT_FILE}

# Replace the old da light address
sed -i "/^da_config/s/http:\/\/.*26658/http\:\/\/${DA_RPC_ADDR}/"  ${ROLLAPP_DYMINT_FILE}

echo "--------------------------------------"
echo "create light token"
light_token=$(meda light auth --node.store ${DA_LIGHT_HOME} write)
echo "light1 token: ${light_token}"

sed -i "/^da_config/s/TOKEN123/${light_token}/" ${ROLLAPP_DYMINT_FILE}

echo "------------------------"
sed -n '/^keyring_home_dir/p' ${ROLLAPP_DYMINT_FILE}
echo "Bind the in-container address, not the host address"
sed -i '/^keyring_home_dir/s#=.*#= "'"${ROLLAPP_CONTAINER_HOME}"'/sequencer_keys"#' ${ROLLAPP_DYMINT_FILE}
sed -i 's|^dym_account_name = .*$|dym_account_name = "sequencer"|' ${ROLLAPP_DYMINT_FILE}
sed -n '/^keyring_home_dir/p' ${ROLLAPP_DYMINT_FILE}
sed -n '/^dym_account_name/p' ${ROLLAPP_DYMINT_FILE}
#===================================================================

echo "Set the DAO and devOperator addresses"
echo "DAO address: ${USER_ADDR_ROLLAPP_DAO}"
echo "devOperator address: ${USER_ADDR_ROLLAPP_OPERATOR}"

jq --arg dao_address ${USER_ADDR_ROLLAPP_DAO} '.app_state.region.dao.address = $dao_address' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"

jq --arg devOperator ${USER_ADDR_ROLLAPP_OPERATOR} '.app_state.region.devOperator.address = $devOperator' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
sed -i '/timeout_broadcast_tx_commit/s/10/360/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
echo "Open ports 26657 and 26656"
sed -i '/^laddr/s/127.0.0.1/0.0.0.0/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"

echo "gas fee 0.001 umec"
# sed -i'' -e "s/^minimum-gas-prices *= .*/minimum-gas-prices = \"0$DENOM\"/" "$APP_CONFIG_FILE"
# sed -i '/^minimum-gas-prices/s/=.*/= "0.001urax,0.001umec"/' "$APP_CONFIG_FILE"
sed -i '/^minimum-gas-prices/s/=.*/= "'"${ROLLAPP_MIN_GAS_PRICES}"'"/' "${ROLLAPP_CONFIG_FILE}"


echo "Set the API port"
sed -i -e '/\[api\]/,+3 s/enable *= .*/enable = true/' "${ROLLAPP_CONFIG_FILE}"
sed -i -e '/\[api\]/,+10 s/swagger *= .*/swagger = true/' "${ROLLAPP_CONFIG_FILE}"
sed -i -e '/\[api\]/,+10 s/localhost/0.0.0.0/' "${ROLLAPP_CONFIG_FILE}"


echo "# ---------------------------------------------------------------------------- #"
echo "#      Create rollapp Node1 Docker Compose.                                             #"
echo "# ---------------------------------------------------------------------------- #"
cat >>docker-compose.yml<<EOF
  rollapp-node1:
    <<: *rollapp-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${RAPP_NODE1_IP}
    environment:
      HUB_RPC_ADDR: ${ME_NODE1_IP}:26657
      DA_RPC_ADDR: ${DA_LIGHT_IP}:26658
      RE_INDEX_START: 1
      ALLOC_THRESHOLD: 6
      SYS_THRESHOLD: 10
      TENDERMINT_CLIENT: 07-tendermint-1
    ports:
      - 46657:26657
      - 9290:9090
      - 3317:1317
    volumes:
      - ${RAPP_NODE1_HOME}:${ROLLAPP_CONTAINER_HOME}
EOF

echo "# ---------------------------------------------------------------------------- #"
echo "#                        Create Rollapp  Node1                                 #"
echo "# ---------------------------------------------------------------------------- #"
if [[ $Interactive_Mode =~ ^[Nn]$ ]]
then
    txhash=$(med tx rollapp create-rollapp ${ROLLAPP_CHAIN_ID}  ${MAX_SEQUENCERS} '{"Addresses":[]}' --from global_dao --broadcast-mode sync --fees 2000000umec --yes --keyring-backend ${KEYRING_BACKEND_NAME} --home ${ME_NODE1_HOME}  -o json | jq -r .txhash)

    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    echo "med q rollapp list --home ${ME_NODE1_HOME}"
    med q rollapp list --home ${ME_NODE1_HOME}
    echo "Withdraw the node1 sequencer amount from ${ME_STAKING_AMOUNT_NODE1_REGION}${ME_SEND_SEQUENCER1_AMOUNT}"
    echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_SEQUENCER}   ${ME_SEND_SEQUENCER1_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
    txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_SEQUENCER}   ${ME_SEND_SEQUENCER1_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
    echo "Check whether the tx hash is on chain: ${txhash}"
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
else
    echo "Give custody the following params (bond only the node1 sequencer)"
    echo "Belongs to Rollapp:${ROLLAPP_CHAIN_ID}"
    echo "max sequencers:${MAX_SEQUENCERS}"
    echo "sequence address:${USER_ADDR_ME_SEQUENCER}"
    read -rp "Press Enter after custody runs rollapp create-rollapp" 
    med q rollapp list --home ${ME_NODE1_HOME}
    read -rp "Press Enter if the sequencer list matches. Otherwise ask custody to run rollapp create-rollapp again rollapp create-rollapp"
    echo "${USER_ADDR_ME_SEQUENCER}account info"
    med query bank balances ${USER_ADDR_ME_SEQUENCER} --home ${ME_NODE1_HOME} -o json
    read -rp "Confirm the balance above is enough, otherwise custody must transfer to this account. The amount should be greater than ${SEQUENCER_AMOUNT}. Press Enter to continue. If the balance is short, press N" REPLY
    if [[ $REPLY =~ ^[Nn]$ ]]
    then
        echo "${USER_ADDR_ME_SEQUENCER}account info"
        med query bank balances ${USER_ADDR_ME_SEQUENCER} --home ${ME_NODE1_HOME} -o json
    fi
    read -rp "Confirm the balance above is enough" 
fi
echo "me bind sequencer (node1 only)"
echo "Get the Node1 sequencer pubkey"
echo "sequencer.info contents:"
cat ${RAPP_NODE1_HOME}/sequencer.info
#Production - stake is 0
#Register the sequencer
#The sequencer address needs umec on the HUB. Custody should fund it regularly. After genesis, fund 1 MEC first. In theory that covers 10000 txs

echo "Prepare to bind the sequencer..."
create_sequencer ${SEQUENCER_MONIKER_NAME1} node1 sequencer

echo "Sequencer status shell:med q sequencer list-sequencer --home ${ME_NODE1_HOME}"
med q sequencer list-sequencer --home ${ME_NODE1_HOME}
echo "-------------------------------------------------------------------------------"

docker compose  up -d rollapp-node1
echo "Wait 5s for services to finish starting"
sleep 5
echo "# ---------------------------------------------------------------------------- #"
echo "#                  Current container status                                                  #"
echo "# ---------------------------------------------------------------------------- #"
docker compose  ps


setup_after_rollapp_node1

#test commands
#Withdraw treasury funds to global_dao
#./bin/med  tx staking withdraw-from-region me_earth  me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t  10mec --from global_dao --home nodes/hub-nodes/node1 -y -o json
#Transfer to rollapp ibc and query the rollapp ibc balance
#./bin/med query bank balances me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t --home nodes/hub-nodes/node1 -o json
#rollappd query bank balances me1qk7jwgdmq5lfraq5l4p3qa5vq9vwa84rzhmhhc  --home nodes/rollapp -o json
#./bin/med tx ibc-transfer transfer transfer channel-1 me1qk7jwgdmq5lfraq5l4p3qa5vq9vwa84rzhmhhc 1mec --from me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t --fees=100000umec  --chain-id mechain_400-1 --home /data/docker-dev/nodes/hub-nodes/node1 --keyring-backend test  -y -o json
#./bin/rollappd query bank balances me1qk7jwgdmq5lfraq5l4p3qa5vq9vwa84rzhmhhc --home nodes/rollapp -o json

#Transfer to ibc and query the ibc balance
#./bin/meda-appd query bank balances me1zps5zt4tvz8uu4ajuq59sdjj2cf48w0nrv0gvd --home nodes/da-nodes/node1
#./bin/med tx ibc-transfer transfer transfer channel-0 me1zps5zt4tvz8uu4ajuq59sdjj2cf48w0nrv0gvd 2mec --from me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t --fees=100000umec  --chain-id mechain_400-1  --keyring-backend test  -y -o json --home nodes/hub-nodes/node1
#./bin/med query tx 9EF49DA3EA7EAE9D1E2CF25515A5934292CF4289CA6B401E7F56065ED7A5BBD4  --home nodes/hub-nodes/node1 -o json
#./bin/meda-appd query bank balances me1zps5zt4tvz8uu4ajuq59sdjj2cf48w0nrv0gvd --home nodes/da-nodes/node1


#Query ibc status
#trusting_period is 511200s, while unbonding_period is 604800s
#./bin/med q ibc client state 07-tendermint-0 --home nodes/hub-nodes/node1 -o json
#./bin/med tx ibc client update 07-tendermint-0 --from <your-key> --chain-id mechain_400-1 --keyring-backend test -y
#Check the ibc client expiry
#How to keep it from expiring:
#1. Send an ibc tx to refresh it before it expires;
#2. IBC The IBC client's built-in refresh
#./bin/rly query clients-expiration hub-rollapp --home ./nodes/rly-rapp  -o json
#./bin/rly query clients-expiration hub-meda   --home ./nodes/rly-da -o json
#update
#./bin/rly  tx update-clients hub-meda --home  ./nodes/rly-da
#./bin/rly  tx update-clients hub-rollapp --home ./nodes/rly-rapp

#Region Query the balance and withdraw to an account
#ME_NODE1_HOME="nodes/hub-nodes/node1" 
#region_data=$(./bin/med  query staking regions --home ${ME_NODE1_HOME} -o json)
#echo "$region_data" | jq -r '.region[] | [.regionId, .region_treasure_addr] | @tsv' | while IFS=$'\t' read -r region_id treasure_addr; do
#    balance_json=$(./bin/med  query bank balances "$treasure_addr" --home ${ME_NODE1_HOME} -o json 2>/dev/null)
#    balance=$(echo "$balance_json" | jq -r '.balances[] | select(.denom=="umec") | .amount // "0"')
#    echo "$region_id current balance is ${balance:-0} umec"
#
#done
#echo "me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t balance is:"
#./bin/med query bank balances me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t --home nodes/hub-nodes/node1 -o json
#./bin/med tx staking withdraw-from-region me_earth me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t  200umec --from global_dao --home nodes/hub-nodes/node1 -y -o json | jq -r '.txhash'
#./bin/med query tx 65BF86C4AFD1E6C5D281A13D21B089E6620956F6EA6F6BD19A65419A1611F4C3 --home nodes/hub-nodes/node1 -o json


# Wipe all node data and return to genesis block 0
# Rarely useful. Only when the whole chain must restart from 1 without a new genesis
# ./bin/med tendermint unsafe-reset-all --home ./nodes/node1


# Query the validator list
# ./bin/med query staking validators --home ./nodes/node1

# Export a user private key
#./bin/med keys export user_name --unsafe --unarmored-hex --keyring-backend test --home nodes/hub-nodes/node1/ | xxd -r -p | base64
# Import a user private key
#./bin/med keys import <name> <keyfile> --home nodes/hub-nodes/node1/ --keyring-backend test


## Edge proxy node settings
# Turn off the p2p address exchange, so this node stops sharing its address book and stays network-isolated
# Run sed -i 's/pex = true/pex = false/' "${ME_NODE1_HOME}"/config/config.toml
# An external node must use a persistent_peers long connection to reach this sentry

#contract upgrade
# 1.  Upload the new code
#med tx wasm store artifacts/matching_contract_v2.wasm --from admin --gas auto --gas-adjustment 1.3 --chain-id your-chain-id -y
# Record the new CODE_ID
# 2.  Run the migration
#med tx wasm migrate $CONTRACT_ADDR <NEW_CODE_ID> '{"new_version":"1.1.0"}' --from admin --chain-id your-chain-id -y

