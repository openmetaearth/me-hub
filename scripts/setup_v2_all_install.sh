#!/bin/bash
#======================================================================================
#config start
#======================================================================================
echo "因为存在服务需要后台运行,所以请使用screen -S install创建独立后台运行本脚本"
echo "请确认无误后再回车"
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

#快速安装模式提取质押额直接为从me_earth上提取，以便快速搭建测试环境；关闭快速安装，则按照生产实际各区提取
read -rp "Quick install for test?[Y]: " INSTALL_TYPE
INSTALL_TYPE=${INSTALL_TYPE:-"Y"}
echo "Quick install for test is set to: $INSTALL_TYPE"

#read -rp "Enable Wasm Contract?[N]: " ENABLE_WASM_CONTRACT
#ENABLE_WASM_CONTRACT=${ENABLE_WASM_CONTRACT:-"N"}
#echo "Enable Wasm Contract is set to: $ENABLE_WASM_CONTRACT"

read -rp "link bin to system?[N]: " LINK_BIN_ACTION
LINK_BIN_ACTION=${LINK_BIN_ACTION:-"N"}
echo "link bin to system is set to: $LINK_BIN_ACTION"

# all=全量安装; rollapp_after=从 rollapp IBC 起续跑（不重建 hub/da/rollapp-node1）
# 也可事先 export RUN_STAGE=rollapp_after，回车即沿用
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
# 助记词
if [[ ! $Randomly_Generate_Mnemonic_Words =~ ^[Yy]$ ]]
then
    global_dict="sign wise lyrics find brass grass already urge invite vibrant arrow grow dune theme trust mix napkin ball pause sponsor bacon capable flip people"
else
    global_dict=$(meda-appd keys mnemonic)
fi
#======================================================================================
# 助记词序号
USER_SEQ_ME_GLOBAL_DAO=1                                                # me超管 global_dao
USER_SEQ_ME_MEID_DAO=2                                                  # meid 地址
USER_SEQ_ME_DEV_OPERATOR=3                                              # 协会地址
USER_SEQ_ME_AIRDROP=4                                                   # 空投地址
USER_SEQ_ME_USER=5                                                      # 普通用户

USER_SEQ_ME_VAL1=11                                                     # 节点验证人 node1
USER_SEQ_ME_VAL2=12                                                     # 节点验证人 node2
USER_SEQ_ME_VAL3=13                                                     # 节点验证人 node3
USER_SEQ_ME_VAL4=14                                                     # 节点验证人 node4

USER_SEQ_ME_SEQUENCER=20                                                # rollapp 绑定在 hub 上的账户
USER_SEQ_ME_IBC_DA=21                                                   # ibc -> da
USER_SEQ_ME_IBC_ROLLAPP=22                                              # ibc -> rollapp

USER_SEQ_ME_SEQUENCER2=30                                               # rollapp 绑定在 hub 上的账户2
USER_SEQ_ME_SEQUENCER3=31                                               # rollapp 绑定在 hub 上的账户3
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

rollapp_user_index_roluser=80                                            # rollapp 用户
rollapp_user_index_dao=81                                                # rollapp dao
rollapp_user_index_dev_operator=82                                       # devOperator 协会地址
rollapp_user_index_ibc=83                                                # rollapp ibc
rollapp_user_index_sync_user=84                                          # rollapp 同步用户
#======================================================================================
#此处为各用户地址，为空时将由助记词+序号推导产生，如需要外部地址如资管模式，则需要填入对应地址
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
#非随机下助记词（Randomly_Generate_Mnemonic_Words=N）将实际产生如下地址,此处仅为记录非修改项目：
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
# Docker 运行时镜像（创世/keys 仍用宿主机 ${BIN_DIR} 二进制）
MED_IMAGE="${MED_IMAGE:-ghcr.io/openmetaearth/med:v2.0.17}"
ROLLAPP_IMAGE="${ROLLAPP_IMAGE:-harbor.starex.xyz/rollapp/rollappd:v1.0.21}"

# ME Hub 目录
ME_NODES_HOME="${BASE_DIR}/nodes/hub-nodes"
# 节点目录
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


#创建治理提案时质押的最小代币数量,单位为方便处理统一已经配置为umec(配置项：app_state.gov.deposit_params.min_deposit[0].denom )"failed to validate genesis state: invalid minimum deposit: 0umec"
ME_MIN_DEPOSIT_AMOUNT="100000000"
#ME创世质押金额
#模块化额度有100亿MEC也就是 1000000000000000000umec 。
TOTAL_SUPPLY="1000000000000000000umec"  # 100亿 MEC
#最小质押额度，单位udmec
DA_MIN_DEPOSIT_AMOUNT="1000000000000"
MIN_TREASURY_BALANCE="100000000"  # 1E umec minimum to keep
dao_TIA_AMOUNT="0udmec" # da dao自身的钱
# ME DA 节点目录
DA_NODES_HOME="${BASE_DIR}/nodes/da-nodes"

# 节点1 名称
DA_NODE1_NAME=${NODE_NAME}

# node1 节点目录
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
# 轻节点
DA_RPC_ADDR="${DA_LIGHT_IP}:26658"
DA_TENDERMINT_CONFIG_FILE="${DA_NODE1_HOME}/config/config.toml"
DA_APP_CONFIG_FILE="${DA_NODE1_HOME}/config/app.toml"
DA_GENESIS_FILE="${DA_NODE1_HOME}/config/genesis.json"


# ROLLAPP 节点目录
RAPP_NODE1_NAME=${NODE_NAME}
RAPP_NODES_HOME="${BASE_DIR}/nodes/rollapp-nodes"
RAPP_NODE1_HOME="${RAPP_NODES_HOME}/${RAPP_NODE1_NAME}"
RAPP_NODE1_IP="${DOCKER_NETWORK_PREFIX}.130"
RAPP_NODE2_IP="${DOCKER_NETWORK_PREFIX}.131"
RAPP_NODE3_IP="${DOCKER_NETWORK_PREFIX}.132"
CHAIN_RPC_ADDR_RAPP="http://${RAPP_NODE1_IP}:26657"

# 货币名字, 默认是 stake
DENOM="urax"
#最小质押额度，单位urax
ROLLAPP_MIN_DEPOSIT_AMOUNT="1000000"

# rly 节点目录
IBC_PORT=transfer
IBC_VERSION=ics20-1
RLY_CONFIG_FILE="${RLY_PATH}/config/config.yaml"
HUB_IBC_CONF_FILE="${RLY_PATH}/hub.json"
DA_IBC_CONF_FILE="${RLY_PATH}/da.json"
#======================================================================================
# rly-da 节点目录
RLY_DA_NODE_HOME="${BASE_DIR}/nodes/rly-da"
# rly-rapp 节点目录
RLY_RAPP_NODE_HOME="${BASE_DIR}/nodes/rly-rapp"

calculate() {
    echo "$1" | bc
}
# 把 timeout_commit (1s / 0.5s / 500ms) 转成秒，供 sleep 计算。bash $(( )) 不能处理小数
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

#experience_region:ME_VAL2 -experience，占比stake_tokens_pool的10%
ME_STAKING_AMOUNT_NODE2="$(calculate "${TOTAL_AMOUNT} * 10 / 100")umec"   # 10% = 100,000,000 MEC = 100000000000000000umec
ME_STAKING_AMOUNT_NODE2_REGION="experience_region"
#usa:ME_VAL3 -ind，占比stake_tokens_pool的17.76%
ME_STAKING_AMOUNT_NODE3="$(calculate "${TOTAL_AMOUNT} * 1776 / 10000")umec"  # 17.76% = 177,600,000 MEC = 177600000000000000umec
ME_STAKING_AMOUNT_NODE3_REGION="ind"
#rus:ME_VAL4 -chn，占比stake_tokens_pool的17.72%
ME_STAKING_AMOUNT_NODE4="$(calculate "${TOTAL_AMOUNT} * 1772 / 10000")umec"   # 17.72% = 177,200,000 MEC = 177200000000000000umec
ME_STAKING_AMOUNT_NODE4_REGION="chn"

#me_earth:ME_VAL1，100亿去掉上面3个国家区后全部质押给me_earth  
STAKING_AMOUNT="$(calculate "${TOTAL_AMOUNT} - ${ME_STAKING_AMOUNT_NODE2%umec} - ${ME_STAKING_AMOUNT_NODE3%umec} - ${ME_STAKING_AMOUNT_NODE4%umec}")umec"  # 54.52% = 545,200,000 MEC =545200000000000000umec
ME_STAKING_AMOUNT_NODE1_REGION="me_earth"
# staking 模块参数 质押的代币需要被锁定的时间， 604800s/7天，18000s/5小时,3024000s/28天，21600s/6小时.
# sequencer 模块参数 For Cosmos SDK-based chains, trust_period should usually be about 2/3 of the unbonding time (~2
# weeks) during which they can be financially punished (slashed) for misbehavior.trust_period is configured in the config.toml.
ME_UNBONDING_TIME="3024000s"

if [[ $INSTALL_TYPE =~ ^[Yy]$ ]]
then
    #区块争议期间,只有超过争议期的才被视为最终确定
    ME_DISPUTE_PERIOD_IN_BLOCKS="5"
    # 激励模块参数，生产环境应该使用“day”，测试环境使用“minute”，控制代币激励的释放频率。实际写死为day，准确是每17280高度释放一次。创建新区，region_share变动，区块高度17280的整数倍，三个条件都会分发区块奖励。
    DISTR_EPOCH_IDENTIFIER="day"
    # 激励模块参数，生产环境应该使用2周"14d"，测试环境使用“60s”，控制代币激励的锁定时间
    LOCKABLE_DURATIONS="60s"
    # voting_period和max_deposit_period：测试环境 600s，正式环境 2 天
    MAX_DEPOSIT_PERIOD="600s"
    # 出块速度。测试环境 500ms，正式环境目前生产为 5s。1330 块大约 11 分钟
    ME_TIMEOUT_COMMIT="500ms"
else
    ME_DISPUTE_PERIOD_IN_BLOCKS="50"
    DISTR_EPOCH_IDENTIFIER="day"
    LOCKABLE_DURATIONS="14d"
    MAX_DEPOSIT_PERIOD="172800s"
    ME_TIMEOUT_COMMIT="5s"
fi
REGION_HIGHT_WAIT="17280"
#======================================================================================
# ME-DA
#======================================================================================

val0_STAKING_AMOUNT="1100000000000udmec"   # da node1进行创世质押的钱
val0_TIA_AMOUNT="2000000000000udmec"     # da val0自身的钱,da的钱到期后会自动销毁
#DA Node新增验证者 创世质押金额
#生产模式下
#‒ 创建4个节点，需要资管签名发送MsgCreateValdator交易 @Lee。
#‒ 资金来源：从 ME-Hub 主网 前四大区金库（IND, CHN, ME_EARTH, EXPERIENCE）提取资金，实际生产region应该有4个以上国家金库故可排除体验区EXPERIENCE以确保稳定性。
#‒ 提取比例：按各区金库质押额的 0.01%（万分之一）动态计算ME_SEND_DA_AMOUNT，实时跨链至 ME-DA 网络。
#当前由global_dao账户转账，故下面四总和不能超过TARGET_BALANCE。
ME_SEND_DA_AMOUNT_VAL1="$(calculate "${STAKING_AMOUNT%umec} / 10000")umec"
ME_SEND_DA_AMOUNT_VAL2="$(calculate "${ME_STAKING_AMOUNT_NODE2%umec} / 10000")umec"
ME_SEND_DA_AMOUNT_VAL3="$(calculate "${ME_STAKING_AMOUNT_NODE3%umec} / 10000")umec"
ME_SEND_DA_AMOUNT_VAL4="$(calculate "${ME_STAKING_AMOUNT_NODE4%umec} / 10000")umec"
ME_SEND_IBC_AMOUNT="10000000umec"
ME_SEND_SEQUENCER1_AMOUNT="1000000000000umec"
ME_SEND_SEQUENCER2_AMOUNT="1000000000000umec"
ME_SEND_SEQUENCER3_AMOUNT="1000000000000umec"
#预留MEC，以便其他用途
ME_RESERVE_AMOUNT="200000000umec"
TARGET_BALANCE="$(calculate "${ME_SEND_DA_AMOUNT_VAL1%umec} + ${ME_SEND_DA_AMOUNT_VAL2%umec} + ${ME_SEND_DA_AMOUNT_VAL3%umec} + ${ME_SEND_DA_AMOUNT_VAL4%umec} + ${ME_SEND_IBC_AMOUNT%umec} + ${ME_RESERVE_AMOUNT%umec} + ${ME_SEND_SEQUENCER1_AMOUNT%umec} + ${ME_SEND_SEQUENCER2_AMOUNT%umec} + ${ME_SEND_SEQUENCER3_AMOUNT%umec}")"
# ME hub 等待块高高度才创建region,每块奖励792mec，也就是79200000000umec，那么预估需要高度，30为随意添加的预留额度数字，仅只是为了保证极端数值下不把金库抽空，一般等待高度在非精确sleep时间中都会超额
ME_CREAT_REGION_HEIGHT="$(calculate "(${TARGET_BALANCE} / 79200000000) + 30")"
#DA_STAKING_AMOUNT应大于 1000000000000udmec 且总和小于ME_SEND_DA_AMOUNT数值*100000000(相等的情况下虽然不收手续费但已经会存在失败)
DA_STAKING_AMOUNT_VAL1=$(umec_to_udmec "${ME_SEND_DA_AMOUNT_VAL1}")
DA_STAKING_AMOUNT_NODE1_REGION="me_earth"
DA_STAKING_AMOUNT_VAL2=$(umec_to_udmec "${ME_SEND_DA_AMOUNT_VAL2}")
DA_STAKING_AMOUNT_NODE2_REGION="usa"
DA_STAKING_AMOUNT_VAL3=$(umec_to_udmec "${ME_SEND_DA_AMOUNT_VAL3}")
DA_STAKING_AMOUNT_NODE3_REGION="ind"
DA_STAKING_AMOUNT_VAL4=$(umec_to_udmec "${ME_SEND_DA_AMOUNT_VAL4}")
DA_STAKING_AMOUNT_NODE4_REGION="chn"

# staking 模块参数 质押的代币需要被锁定的时间， 604800s/7天，18000s/5小时.For Cosmos SDK-based chains, trust_period should usually be about 2/3 of the unbonding time (~2
# weeks) during which they can be financially punished (slashed) for misbehavior.实际为数值减少100倍再乘与85，属于跨链信任时间。（med q ibc client state 07-tendermint-0）。3024000s预估trust_period为28天
# trusting_period=UnbondingTime / 100 * 85 // TODO: replace with percentage
# 需要注意与RLY_TIME_THRESHOLD的关系
DA_UNBONDING_TIME=${ME_UNBONDING_TIME}
# staking 模块参数 最大验证者数量
MAX_VALIDATORS="4"
#区块争议期间,只有超过争议期的才被视为最终确定,推荐50

if [[ $INSTALL_TYPE =~ ^[Yy]$  ]]
then
    DA_DISPUTE_PERIOD_IN_BLOCKS="5"
else
    DA_DISPUTE_PERIOD_IN_BLOCKS="50"
fi
#======================================================================================
# ME-Rollapp
#======================================================================================


# 初始化余额rollapp_user_roluser_TOKEN_AMOUNT应该大于等于rollapp_user_roluser_NODE_STAKING_AMOUNT与rollapp_transfer_ibc_amount的总和
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
#rollapp跨链到hub的时间取决于状态提交速度，所以dymint配置提交状态的参数batch_submit_max_time测试时使用60秒，正式环境也配置batch_submit_max_time为60秒；
#从rollapp跨链到hub，如果10分钟提交一次，那么至少15分钟才能到账；
BATCH_SUBMIT_MAX_TIME="60s"
#最大排序器数量
MAX_SEQUENCERS=10
SEQUENCER_MONIKER_NAME1="sequencer1"
SEQUENCER_MONIKER_NAME2="sequencer2"
SEQUENCER_MONIKER_NAME3="sequencer3"
# 链上 min bond 为 100000000umec（create-sequencer 用 100umec 会报 code 1011 insufficient bond）
# 须小于已打到 sequencer 账户的 ME_SEND_SEQUENCER*_AMOUNT
SEQUENCER_AMOUNT="100000000umec"
#======================================================================================
# ME-Rly
#======================================================================================


# 中继服务更新间隔,需要小于 (UnbondingTime / 100 * 85),链认证者短期出现大于1/3变动时（初始化4节点时，每新增1个认证者时间），间隔时间间隔需要大于该值，以避免出现信任周期内数值校验不通过而无法续期
DA_UNBONDING_SECONDS=${DA_UNBONDING_TIME%s}
#docker不指定更新周期时，代码默认更新周期的计算方式
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
# da ibc代币生效块高,注意事项：创世时需要在创世文件中设置 ibc_deadline_block_height ，如果设置为1000(一个半小时)，那么hub需要在1000个区块内跨链到DA做质押，否则da链会报错。在1000个区块之前，只能白名单和DAO地址可以发交易，默认值100。
IBC_DEADLINE_BLOCK_HEIGHT=300


#======================================================================================
# 自定义修改区：请修改以上配置
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

                # 如果上面失败则走无须确认的y的确认项目
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

                # 如果上面失败则走无须确认的y的确认项目
                echo -e "${mnemonics}\n" | ${action} keys add "$key_name" \
                        --keyring-backend ${KEYRING_BACKEND_NAME} \
                        --hd-path "m/44'/118'/0'/0/0/${user_seq}" \
                        --recover $extra_args --coin-type 118
        else
                echo "$output"
        fi
    fi
}


# 检查hash是否上链, 最多等待60s
check_tx_status() {
    local cosmosapp=$1
    local txhash=$2
    local extra_args=$3
    sleepTime=120

    # 检查参数是否为空
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
                echo "当前结果:"
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
    echo "使用质押金额创建排序器"
    local sequencer_metadata=$(jq -n --arg moniker "${SEQUENCER_MONIKER_NAME}" \
        '{Moniker: $moniker, Identity: "", Website: "", SecurityContact: "", Details: ""}')
    
    echo "创建排序器，元数据: ${sequencer_metadata}"
    
    txhash=$(med tx sequencer create-sequencer ${SEQUENCER_PUBLIC_KEY} ${ROLLAPP_CHAIN_ID} "${sequencer_metadata}" ${SEQUENCER_AMOUNT} \
        --from ${SEQ_NAME} --chain-id ${ME_CHAIN_ID} \
        --fees 40000umec --gas 2000000 --keyring-backend ${KEYRING_BACKEND_NAME} --home ${ME_NODE1_HOME} --yes -o json | jq -r .txhash)
    
    echo "排序器创建交易哈希: ${txhash}"

    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
}


load_existing_addrs() {
    local kb=(--keyring-backend "${KEYRING_BACKEND_NAME}")
    echo "# 从现有 keyring 读取地址，供 rollapp_after 续跑"
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
    # hub->rollapp ibc 收款地址与助记词序号 22 相同，不一定在 hub keyring 里
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
echo "#       等待 rollapp RPC 可用并开始出块（IBC client 由后面 rly tx link 创建）      #"
echo "# ---------------------------------------------------------------------------- #"
while true; do
    height=$(rollappd status --home ${RAPP_NODE1_HOME} 2>/dev/null | jq -r '.SyncInfo.latest_block_height // .sync_info.latest_block_height // empty')
    height=${height//\"/}
    if [[ -n "${height}" && "${height}" -ge 1 ]] 2>/dev/null; then
        echo "rollapp 已出块, height=${height}"
        break
    fi
    echo "rollapp 尚未出块或 RPC 未就绪, sleep 5s"
    sleep 5
done

echo "# ---------------------------------------------------------------------------- #"
echo "#      部署rly-rollapp                                                          #"
echo "# ---------------------------------------------------------------------------- #"
echo "设置通道信息"
echo "初始化 rly 配置文件:rly config init --home ${RLY_RAPP_NODE_HOME}"
rly config init --home ${RLY_RAPP_NODE_HOME}
echo 导入配置
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
echo "#      是否导入me_user_ibc_rollapp_name、rollapp_user_ibc_me_name秘钥            #"
echo "# ---------------------------------------------------------------------------- #"
if [[ "${Sync_Key}" == [Yy] ]]; then
    mkdir -p ${RLY_RAPP_NODE_HOME}/keys/${ME_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    mkdir -p ${RLY_RAPP_NODE_HOME}/keys/${ROLLAPP_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    echo "处理me_user_ibc_rollapp_name等价于USER_ADDR_ME_IBC_ROLLAPP"
    me_user_ibc_rollapp_name=$(add_key "meda-appd" "me_user_ibc_rollapp_name" ${USER_SEQ_ME_IBC_ROLLAPP} "--home ./tmpxx10")
    echo "${me_user_ibc_rollapp_name}"
    USER_ME_IBC_DA_NAME=$(echo "${me_user_ibc_rollapp_name}" | grep -oP 'address: \K[^\s]+')
    USER_ME_IBC_DA_NAME_PUBKEY=$(echo "${me_user_ibc_rollapp_name}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    cp -rf ./tmpxx10/keyring-${KEYRING_BACKEND_NAME}/*  ${RLY_RAPP_NODE_HOME}/keys/${ME_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    rm -rf ./tmpxx10

    echo "处理rollapp_user_ibc_me_name等价DA上的USER_ADDR_ROLLAPP_IBC-ibc"
    rollapp_user_ibc_me_name=$(add_key "meda-appd" "rollapp_user_ibc_me_name" ${rollapp_user_index_ibc} "--home ./tmpxx10")
    echo "$rollapp_user_ibc_me_name"
    USER_ROLLAPP_IBC_ME_NAME=$(echo "$rollapp_user_ibc_me_name" | grep -oP 'address: \K[^\s]+')
    USER_ROLLAPP_IBC_ME_NAME_PUBKEY=$(echo "$rollapp_user_ibc_me_name" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    cp -rf ./tmpxx10/keyring-${KEYRING_BACKEND_NAME}/*  ${RLY_RAPP_NODE_HOME}/keys/${ROLLAPP_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    rm -rf ./tmpxx10
fi


echo "# 避免 07-tendermint-0: light client not found 错误, 等待10s"
sleep 10

echo "create link hub-rollapp:rly tx link hub-rollapp --src-port transfer --dst-port transfer --version ics20-1 --home ${RLY_RAPP_NODE_HOME}"

rly tx link hub-rollapp --src-port transfer --dst-port transfer --version ics20-1 --home ${RLY_RAPP_NODE_HOME}

# 检查上一个命令的退出状态
if [ $? -ne 0 ]; then
    echo "creat命令执行失败, 退出脚本"
    exit 1
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#      rly-rollapp验证                                                          #"
echo "# ---------------------------------------------------------------------------- #"
echo 检查${ME_CHAIN_ID} 客户端状态
rly q clients ${ME_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}
echo 检查${ROLLAPP_CHAIN_ID}客户端状态
rly q clients ${ROLLAPP_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}

echo 检查连接状态
rly q connections ${ME_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}



echo "查询链上的上通道信息"
rly q channels ${ME_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}
echo "-------------------------------------------------------------------------------"
rly q channels ${ROLLAPP_CHAIN_ID} --home ${RLY_RAPP_NODE_HOME}

echo "获取通道ID"
echo "查询ME_CHANNEL_ID:rly q channels "${ME_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ROLLAPP_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id'"
ME_CHANNEL_ID=$(rly q channels "${ME_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ROLLAPP_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id')

echo "查询RAPP_CHANNEL_ID:rly q channels "${ROLLAPP_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ME_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id'"
RAPP_CHANNEL_ID=$(rly q channels "${ROLLAPP_CHAIN_ID}" --home "${RLY_RAPP_NODE_HOME}" 2>/dev/null | jq --arg chain_id "${ME_CHAIN_ID}" -r 'select(.counterparty.chain_id == $chain_id and .state == "STATE_OPEN") | .channel_id')

echo "me -> rapp, channel_id: ${ME_CHANNEL_ID}"
echo "rapp -> me, channel_id: ${RAPP_CHANNEL_ID}"


echo "设定通道时间"
rly transact client ${ME_CHAIN_ID} ${ROLLAPP_CHAIN_ID} hub-rollapp --client-tp ${CLIENT_TP} --home ${RLY_RAPP_NODE_HOME}


echo "添加relayer到docker-compose.yml"
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
    echo "查询 me-hub 链上 ibc-rollapp 地址${USER_ADDR_ME_IBC_ROLLAPP}余额:"
    med query bank balances ${USER_ADDR_ME_IBC_ROLLAPP} --home ${ME_NODE1_HOME}

    echo "rollapp 原生代币持有者发起一笔 ibc 转账交易,金额为${rollapp_transfer_ibc_amount}"
    echo "rollappd tx ibc-transfer transfer transfer ${RAPP_CHANNEL_ID} ${USER_ADDR_ME_IBC_ROLLAPP} ${rollapp_transfer_ibc_amount} --from ${USER_ADDR_ROLLAPP_ADMIN} --fees=2000urax --chain-id ${ROLLAPP_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME} -y"

    txhash=$(rollappd tx ibc-transfer transfer transfer ${RAPP_CHANNEL_ID} ${USER_ADDR_ME_IBC_ROLLAPP} ${rollapp_transfer_ibc_amount} --from ${USER_ADDR_ROLLAPP_ADMIN} --fees=2000urax --chain-id ${ROLLAPP_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME} -y --output json | jq -r '.txhash')

    check_tx_status rollappd ${txhash} "--home $RAPP_NODE1_HOME"
else
    echo "# -------------------------------------------------------------------------------- #"
    echo "#      资管转账模式                                                                  #"
    echo "# -------------------------------------------------------------------------------- #"
    echo "ibc: rollapp to me-hub"
    echo "查询 me-hub 链上 ibc-rollapp 地址${USER_ADDR_ME_IBC_ROLLAPP}余额:"
    echo  "med query bank balances ${USER_ADDR_ME_IBC_ROLLAPP} --home ${ME_NODE1_HOME}"
    echo "无误请确认回车"
    read -p "Do you want to continue? (Y/n)"
fi

echo  "查看是否激活 rollapp, 激活后开启快速到账"
echo "执行shell: med q rollapp show ${ROLLAPP_CHAIN_ID}  --home \"${ME_NODE1_HOME}\" -o json | jq .rollapp.genesis_state.transfers_enabled"
while true; do
  rollapp_status=$(med q rollapp show ${ROLLAPP_CHAIN_ID}  --home "${ME_NODE1_HOME}" -o json | jq .rollapp.genesis_state.transfers_enabled)
  if [[ ${rollapp_status} == "true" ]];then
      echo "rollapp 激活成功"
      break
  else
      echo "rollapp 未激活, sleep 1s"
      sleep 1
      continue
  fi
done


if [[ $Interactive_Mode =~ ^[Nn]$ ]]
then
    echo "ibc: rollapp to me-hub"
    echo "执行shell: med tx rollapp skip-delay-rollapp ${ROLLAPP_CHAIN_ID} true --from ${USER_ADDR_ME_GLOBAL_DAO} --gas 500000 --chain-id ${ME_CHAIN_ID} --keyring-backend=${KEYRING_BACKEND_NAME} --home  $ME_NODE1_HOME  -y --output json | jq -r '.txhash'"
    txhash=$(med tx rollapp skip-delay-rollapp ${ROLLAPP_CHAIN_ID} true --from ${USER_ADDR_ME_GLOBAL_DAO} --gas 500000 --chain-id ${ME_CHAIN_ID} --keyring-backend=${KEYRING_BACKEND_NAME} --home  $ME_NODE1_HOME  -y --output json | jq -r '.txhash')
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    echo "set skip-delay-rollapp success"
    sleep 2

    #===================================================================
    # me-hub to rollapp

    echo ""
    echo "ibc: me-hub to rollapp"
    echo ""

    echo "查询ibc余额shell: rollappd query bank balances ${USER_ADDR_ROLLAPP_IBC} --home ${RAPP_NODE1_HOME}"
    rollappd query bank balances ${USER_ADDR_ROLLAPP_IBC} --home ${RAPP_NODE1_HOME}

    echo "ME的user转账10mec到rollapp的ibc地址shell: med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${USER_ADDR_ROLLAPP_IBC}  ${ME_SEND_IBC_AMOUNT} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=100000umec  --chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME} --keyring-backend ${KEYRING_BACKEND_NAME}  -y"
    txhash=$(med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${USER_ADDR_ROLLAPP_IBC} ${ME_SEND_IBC_AMOUNT} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=100000umec  --chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME} --keyring-backend ${KEYRING_BACKEND_NAME}  -y --output json | jq -r '.txhash')
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"

    #===================================================================
else
    echo "# -------------------------------------------------------------------------------- #"
    echo "#      资管转账模式                                                                  #"
    echo "# -------------------------------------------------------------------------------- #"
    echo "ibc: rollapp to me-hub"
    echo "查询 me-hub 链上 ibc-rollapp 地址${USER_ADDR_ROLLAPP_IBC}余额:"
    echo  "rollappd query bank balances ${USER_ADDR_ROLLAPP_IBC} --home ${RAPP_NODE1_HOME}"
    echo "无误请确认回车"
    read -p "Do you want to continue? (Y/n)"
fi
echo "========================="
# 持续检查地址是否拥有指定代币
watch_addr_balance_type() {
  local binary="$1"
  local token_prefix="$2"   # 指定的代币前缀
  local user_addr="$3"      # 用户地址
  local extra_args="$4"     # 额外参数
  local check_status=0      # 检测状态

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
echo "#              Creat RollApp Node2-Node3                                       #"
echo "# ---------------------------------------------------------------------------- #"
echo "检查静态配置："
echo "dym_account_name: $(grep '^dym_account_name' ${RAPP_NODES_HOME}/node1/config/dymint.toml)"
echo "keyring_home_dir: $(grep '^keyring_home_dir' ${RAPP_NODES_HOME}/node1/config/dymint.toml)"
echo "da_config: $(grep '^da_config' ${RAPP_NODES_HOME}/node1/config/dymint.toml)"
echo "settlement_node_address: $(grep '^settlement_node_address' ${RAPP_NODES_HOME}/node1/config/dymint.toml)"
rollapp_init_sync_node() {
    local NODE_NAME=$1
    local NODE_HOME="${RAPP_NODES_HOME}/${NODE_NAME}"
    local NODE_IP=$2
    local SEQ_NAME=$3

    local CONFIG_DIRECTORY="$NODE_HOME/config"
    local TENDERMINT_CONF="$CONFIG_DIRECTORY/config.toml"
    local ROLLAPP_APP_CONFIG_FILE="$CONFIG_DIRECTORY/app.toml"
    local CLIENT_CONFIG_FILE="$CONFIG_DIRECTORY/client.toml"
    local DYMINT_CONFIG_FILE="$CONFIG_DIRECTORY/dymint.toml"

    echo "========================================"
    echo "init node: ${NODE_NAME} - ${NODE_HOME}"


    rollappd init ${NODE_NAME} --chain-id="${ROLLAPP_CHAIN_ID}" --home ${NODE_HOME}
    # 显示sequencer信息
    rollappd dymint show-sequencer --home ${NODE_HOME} > ${NODE_HOME}/sequencer.info
    # 复制配置文件
    cp -rf ${RAPP_NODE1_HOME}/config/app.toml ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/client.toml ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/config.toml ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/dymint.toml ${CONFIG_DIRECTORY}/
    # 复制genesis和gentx
    cp -rf ${RAPP_NODE1_HOME}/config/genesis.json ${CONFIG_DIRECTORY}/
    cp -rf ${RAPP_NODE1_HOME}/config/gentx ${CONFIG_DIRECTORY}/
    # 修改配置文件
    sed -i "s|^moniker = .*$|moniker = \"${NODE_NAME}\"|" "${TENDERMINT_CONF}"
    sed -i "s|^node = .*$|node = \"tcp://${NODE_IP}:26657\"|" "${CLIENT_CONFIG_FILE}"
    sed -i "s|^dym_account_name = .*$|dym_account_name = \"${SEQ_NAME}\"|" "${DYMINT_CONFIG_FILE}"
    # 验证配置
    echo "检查静态配置："
    echo "dym_account_name: $(grep '^dym_account_name' ${DYMINT_CONFIG_FILE})"
    echo "keyring_home_dir: $(grep '^keyring_home_dir' ${DYMINT_CONFIG_FILE})"
    echo "da_config: $(grep '^da_config' ${DYMINT_CONFIG_FILE})"
    echo "settlement_node_address: $(grep '^settlement_node_address' ${DYMINT_CONFIG_FILE})"
}
rollapp_init_sync_node node2 ${RAPP_NODE2_IP} sequencer2
rollapp_init_sync_node node3 ${RAPP_NODE3_IP} sequencer3


# rollapp node2 恢复 seq2 用户
if [ -z "${USER_ADDR_ROLLAPP_SEQUENCER2}" ]; then
    echo "恢复 sequencer2 账户..."
    mkdir -p "${RAPP_NODES_HOME}"/node2/sequencer_keys/keyring-${KEYRING_BACKEND_NAME}
    rollapp_seq2=$(add_key  "rollappd"  "sequencer2" ${USER_SEQ_ME_SEQUENCER2} "--home "${RAPP_NODES_HOME}"/node2/sequencer_keys")
    echo "sequencer2 账户恢复完成"
fi
# rollapp node3 恢复 seq3 用户
if [ -z "${USER_ADDR_ROLLAPP_SEQUENCER3}" ]; then
    echo "恢复 sequencer3 账户..."
    mkdir -p "${RAPP_NODES_HOME}"/node3/sequencer_keys/keyring-${KEYRING_BACKEND_NAME}
    rollapp_seq3=$(add_key  "rollappd"  "sequencer3" ${USER_SEQ_ME_SEQUENCER3} "--home "${RAPP_NODES_HOME}"/node3/sequencer_keys")
    echo "sequencer3 账户恢复完成"
fi
echo "# ---------------------------------------------------------------------------- #"
echo "#      创建rollapp Node2-Node3 排序器。                                          #"
echo "# ---------------------------------------------------------------------------- #"

create_sequencer ${SEQUENCER_MONIKER_NAME2} node2 sequencer2
create_sequencer ${SEQUENCER_MONIKER_NAME3} node3 sequencer3

echo "检查排序器状态shell:med q sequencer list-sequencer --home ${ME_NODE1_HOME}"
med q sequencer list-sequencer --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#      创建rollapp Node2-Node3 Docker Compose。                                 #"
echo "# ---------------------------------------------------------------------------- #"
cat >>docker-compose.yml<<EOF
  rollapp-node2:
    <<: *rollapp-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${RAPP_NODE2_IP}
    volumes:
      - ${RAPP_NODES_HOME}/node2:/root/.rollapp
  rollapp-node3:
    <<: *rollapp-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${RAPP_NODE3_IP}
    volumes:
      - ${RAPP_NODES_HOME}/node3:/root/.rollapp
EOF
echo "# ---------------------------------------------------------------------------- #"
echo "#                              RollApp读取Node ID                               #"
echo "# ---------------------------------------------------------------------------- #"
RAPP_NODEID_NODE1=$(rollappd dymint show-node-id --home ${RAPP_NODES_HOME}/node1  2>&1 | tail -n 1)
RAPP_NODEID_NODE2=$(rollappd dymint show-node-id --home ${RAPP_NODES_HOME}/node2 2>/dev/null | tail -n 1)
RAPP_NODEID_NODE3=$(rollappd dymint show-node-id --home ${RAPP_NODES_HOME}/node3 2>/dev/null | tail -n 1)
echo Rapp_Node1_ID:${RAPP_NODEID_NODE1}
echo Rapp_Node2_ID:${RAPP_NODEID_NODE2}
echo Rapp_Node3_ID:${RAPP_NODEID_NODE3}
echo "# ---------------------------------------------------------------------------- #"
echo "#                         修改Node2-Node5的persistent_peers地址                #"
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
echo "#                  启动rollapp  Node2 Node3                                     #"
echo "# --------------------- -------------------------------------------------------#"

docker compose up -d rollapp-node2
echo "等待5s服务完全启动"
sleep 5
rly  tx update-clients hub-rollapp --home ${RLY_RAPP_NODE_HOME}
docker compose up -d rollapp-node3
echo "等待5s服务完全启动"
sleep 5
rly  tx update-clients hub-rollapp --home ${RLY_RAPP_NODE_HOME}
echo "# ---------------------------------------------------------------------------- #"
echo "#                  修改node1 persistent_peers地址                              #"
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
    echo "#      部署合约                                                                 #"
    echo "# ---------------------------------------------------------------------------- #"
    #上传合约
    #med tx wasm store contract/matching_contract.wasm --from global_dao --gas auto --gas-adjustment 1.3 --chain-id ${ME_CHAIN_ID} --home ${ME_NODE1_HOME} -y
    #txhash=$(med tx wasm store contract/matching_contract.wasm --from global_dao --gas auto --gas-adjustment 1.3 --chain-id mechain_400-1 --home ${ME_NODE1_HOME} -y --output json | jq -r '.txhash')
    #check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    #实例化合约
    #med tx wasm instantiate <CODE_ID> '{"global_dao":"<admin_address>","fee_collector":"<fee_collector_address>"}' --from global_dao --label "matching-v1" --global_dao <admin_address> --chain-id mechain_400-1 --home ${ME_NODE1_HOME} -y
    #txhash=$(med tx wasm instantiate <CODE_ID> '{"admin":"<admin_address>","fee_collector":"<fee_collector_address>"}' --from global_dao --label "matching-v1" --global_dao <admin_address> --chain-id mechain_400-1 --home ${ME_NODE1_HOME} -y --output json | jq -r '.txhash')
    #check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    #记录合约地址 CONTRACT_ADDR
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#      部署状态                                                                 #"
echo "# ---------------------------------------------------------------------------- #"
echo "当前HUB状态"
med status --node tcp://${ME_NODE1_IP}:26657
echo "当前DA状态"
meda-appd status --node tcp://${DA_NODE1_IP}:26657
echo "当前rollapp状态"
rollappd status --node tcp://${RAPP_NODE1_IP}:26657
echo "当前跨链状态"
med  q ibc client state 07-tendermint-0 --home ${ME_NODE1_HOME} -o json
echo "当前ME global_dao 余额"
med query bank balances ${USER_ADDR_ME_GLOBAL_DAO} --home ${ME_NODE1_HOME}  -o json
echo "当前DA val0 余额(若DA尚未到达块hight ${IBC_DEADLINE_BLOCK_HEIGHT}存在余额属于正常)"
meda-appd query bank balances ${USER_ADDR_DA_VAL0} --home ${DA_NODE1_HOME}  -o json
echo "当前DA val1 余额"
meda-appd query bank balances ${USER_ADDR_DA_VAL1} --home ${DA_NODE1_HOME}  -o json
echo "当前DA val2 余额"
meda-appd query bank balances ${USER_ADDR_DA_VAL2} --home ${DA_NODE1_HOME}  -o json
echo "当前DA val3 余额"
meda-appd query bank balances ${USER_ADDR_DA_VAL3} --home ${DA_NODE1_HOME}  -o json
echo "当前DA val4 余额"
meda-appd query bank balances ${USER_ADDR_DA_VAL4} --home ${DA_NODE1_HOME}  -o json
echo "当前rollapp roluser余额"
rollappd query bank balances ${USER_ADDR_ROLLAPP_ADMIN} --home ${RAPP_NODE1_HOME}  -o json
echo "当前rollapp ibc 余额"
rollappd query bank balances ${USER_ADDR_ROLLAPP_IBC} --home ${RAPP_NODE1_HOME}  -o json
echo "hub REGION_ID"
med query staking validators --home ${ME_NODE1_HOME} -o json | jq -r '.validators[] | [.description.regionID, .tokens, .status] | @tsv' | column -t -s $'\t' -N "REGION_ID,TOKENS,STATUS"
echo "hub MONIKER"
meda-appd query staking validators --home ${DA_NODE1_HOME} -o json | jq -r '.validators[] | [.description.moniker, .tokens, .status] | @tsv' | column -t -s $'\t' -N "MONIKER,TOKENS,STATUS"
echo "跨链信任周期"
rly query clients-expiration hub-rollapp --home ${RLY_RAPP_NODE_HOME}
rly query clients-expiration hub-meda   --home  ${RLY_DA_NODE_HOME}

end_time=$(date +%s)
duration=$((end_time - start_time))
echo "-----------------------------------"
echo "脚本执行完成！"
echo "总耗时: $duration 秒"
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
    echo "#     续跑模式 RUN_STAGE=rollapp_after：跳过创世，从 rollapp IBC 继续              #"
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
echo "#                       创建docker compose网络                                  #"
echo "# ---------------------------------------------------------------------------- #"

if docker network inspect "$Docker_Network_Name" >/dev/null 2>&1; then
    echo "Error: Docker network '$Docker_Network_Name' already exists." >&2
    echo "请确认是否删除$Docker_Network_Name网络!!!!!!!!!!? (Y/n) "
    read -r del_network
    del_network=${del_network:-N}
    if [ ! "$del_network" != "${del_network#[Yy]}" ] ;then
        echo "你选择不删除网络，那请重新提供新的网络名"
        echo "当前已经存在如下网络:"
        docker network ls --format '{{.Name}}'
        while true; do
            read -rp "Please enter a new unique Docker_Network_Name: " Docker_Network_Name
            # 检查网络名是否存在
            if ! docker network inspect "${Docker_Network_Name}" >/dev/null 2>&1; then
                break
            fi
            echo "Error: Network ${Docker_Network_Name} already exists. Please choose a different name."
        done
        echo "Docker_Network_Name is set to: ${Docker_Network_Name}"
    else
        echo "正在删除网络..."
        docker network rm ${Docker_Network_Name}
        if [ $? -ne 0 ]; then
            echo "删除网络失败，退出脚本"
            exit 1
        fi
    fi
fi

verify_return_code() {
    # 验证返回值
    if [ $? -eq 0 ]; then
        echo "网络创建成功"
    else
        echo "可能存在相同的网络段，请重新提供网络段"
        exit 1
    fi
}
docker network create ${Docker_Network_Name} --subnet ${Docker_ADDR_RANGE} --ip-range ${Docker_ADDR_RANGE}
# 如果创建这个网络失败就退出
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
echo "修改/config/config.toml"
# Mempool version to use:
#   1) "v0" - (default) FIFO mempool.
#   2) "v1" - prioritized mempool (deprecated; will be removed in the next release).
#sed -i'' -e "/\[mempool\]/,+8 s/v0/v1/" "$TENDERMINT_CONFIG_FILE"
#允许单机器多服务不同端口
sed -i '/^allow_duplicate_ip/s/false/true/' "${ME_NODE1_HOME}"/config/config.toml
#是否启用地址簿的严格性检查。设置为 false 的影响：允许使用非路由地址（如私有 IP 地址）作为对等节点；允许使用非全局单播地址；
sed -i '/^addr_book_strict/s/true/false/' "${ME_NODE1_HOME}"/config/config.toml
sed -i 's/0.0.0.0:26656/0.0.0.0:26656/' "${ME_NODE1_HOME}"/config/config.toml
sed -i 's/127.0.0.1:26657/0.0.0.0:26657/' "${ME_NODE1_HOME}"/config/config.toml
# 修改出块时间
sed -i 's/^timeout_commit = .*/timeout_commit = "'${ME_TIMEOUT_COMMIT}'"/' "${ME_NODE1_HOME}"/config/config.toml
#允许跨域
sed -i 's/cors_allowed_origins.*$/cors_allowed_origins = ["*"]/' "${ME_NODE1_HOME}"/config/config.toml
sed -i 's/localhost:6060/0.0.0.0:6060/' "${ME_NODE1_HOME}"/config/config.toml
echo "修改/config/app.toml"
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

#节点接受交易的最低 gas 价格。交易费用 = gas 价格 × gas 使用量。可以通过治理提案调整，不同验证人可以设置不同的值。
sed -i'' -e 's/^minimum-gas-prices *= .*/minimum-gas-prices = "0.02umec"/' "${ME_NODE1_HOME}/config/app.toml"
# default: the last 362880 states are kept, pruning at 10 block intervals
# nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
# everything: 2 latest states will be kept; pruning at 10 block intervals.
# custom: allow pruning options to be manually specified through 'pruning-keep-recent', and 'pruning-interval'
sed -i'' -e 's/^pruning *= .*/pruning = "nothing"/' "${ME_NODE1_HOME}/config/app.toml"

sed -i'' -e 's/^max-recv-msg-size *= .*/max-recv-msg-size = "1048576000"/' "${ME_NODE1_HOME}/config/app.toml"

echo "修改/config/client.toml"
sed -i "s/localhost:26657/${ME_NODE1_IP}:26657/" "${ME_NODE1_HOME}"/config/client.toml
# The keyring's backend, where the keys are stored (os|file|kwallet|pass|test|memory)
echo "替换钱包存储方式，当前为：${KEYRING_BACKEND_NAME}"
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
    jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.voting_params.voting_period = $period' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.deposit_params.max_deposit_period = $period' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq '.app_state.gov.params.min_deposit[0].denom = "umec"' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg amount "$ME_MIN_DEPOSIT_AMOUNT" '.app_state.gov.params.min_deposit[0].amount = $amount' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    #flase为不销毁，用于控制否决（veto）投票的押金处理方式：当设置为 true 时：用于否决投票的代币将被销毁（永久从流通中移除），当设置为 false 时：用于否决投票的代币将退还给投票者
    jq '.app_state.gov.params.burn_vote_veto = false' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.params.voting_period = $period' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.params.max_deposit_period = $period' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    # expedited_voting_period 必须小于 voting_period，否则 genesis validate 失败
    jq --arg period "300s" 'if .app_state.gov.params.expedited_voting_period then .app_state.gov.params.expedited_voting_period = $period else . end' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
}

set_hub_params() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                              setting hub params                               #"
    echo "# ---------------------------------------------------------------------------- #"
    sed -i'' -e 's/bond_denom": ".*"/bond_denom": "umec"/' "${ME_GENESIS_FILE}"
    sed -i'' -e 's/mint_denom": ".*"/mint_denom": "umec"/' "${ME_GENESIS_FILE}"

    jq ".app_state.rollapp.params.dispute_period_in_blocks = \"$ME_DISPUTE_PERIOD_IN_BLOCKS\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    #bridging_fee 是跨链桥接功能中的一个参数，用于设置跨链交易的手续费，费用通常以原生代币支付
    jq ".app_state.delayedack.params.bridging_fee = \"0\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    #epoch_identifier 是跨链桥接功能中的一个参数，用于设置跨链交易的周期，周期可以是 week、month、year 等
    jq ".app_state.delayedack.params.epoch_identifier = \"week\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq ".app_state.eibc.params.epoch_identifier = \"week\"" "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    #每个周期删除的IBC包数量，设置为1000000；
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
    echo "#                        setting incentives params(激励)                        #"
    echo "# ---------------------------------------------------------------------------- #"
    jq --arg time "$DISTR_EPOCH_IDENTIFIER" '.app_state.incentives.params.distr_epoch_identifier = $time' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"
    jq --arg time "$LOCKABLE_DURATIONS" '.app_state.incentives.lockable_durations = ["$time"]' "${ME_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ME_GENESIS_FILE}"

}


set_misc_params() {
    echo "# ---------------------------------------------------------------------------- #"
    echo "#  setting misc params(设置危机交易（crisis transaction）的固定费用代币单位)         #"
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
  #rollapp的创建权限地址
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
echo "# 判断(global_dao、meid_dao、dev_operator、airdrop、val1、user、sequencer)是否需要恢复add #"
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
echo "# 初始化余额 add-genesis-account global_dao-"0umec"                                       #"
echo "# (meid_dao、dev_operator、airdrop、val1)-"0umec"                                        #"                            
echo "# --------------------------------------------------------------------------------------#"
med add-genesis-account global_dao   "0umec" --home ${ME_NODE1_HOME}
med add-genesis-account meid_dao     "0umec" --home ${ME_NODE1_HOME}
med add-genesis-account dev_operator "0umec" --home ${ME_NODE1_HOME}
med add-genesis-account airdrop      "0umec" --home ${ME_NODE1_HOME}
med add-genesis-account val1         "0umec" --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#                              设定质押池                                       #"
echo "# ---------------------------------------------------------------------------- #"

med add-genesis-stake-pool  --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#                         创世初始化地址: 内部模块地址                           #"
echo "# ---------------------------------------------------------------------------- #"

med add-genesis-m-accounts --home ${ME_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#   Set dao_addresses(global_dao、meid_dao、dev_operator、airdrop_address)      #"
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
    echo "#  setting issuer:USER_ADDR_ME_GLOBAL_DAO、USER_ADDR_ME_MEID_DAO(定义哪些地址有权限颁发或验证 KYC 认证) #"
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
echo "#                创世质押(生产模式使用的是内置账户stake_tokens_pool提供质押额度)       #"
echo "# ---------------------------------------------------------------------------- #"

#global_dao # 发送交易的账户名称
#STAKING_AMOUNT # 质押金额
#--chain-id 区块链网络ID
#--keyring-backend # 密钥存储后端（如test或file）
#--region-id. # 验证人所在地区ID（地球区）
#--validator-address  # 验证人操作地址


# 检查模块账户
jq '.app_state.auth.accounts[] | select(.name=="stake_tokens_pool")' $ME_NODE1_HOME/config/genesis.json

# 检查余额
jq '.app_state.bank.balances[] | select(.address | contains("me1gmpxkchcdgfq995zye5efwzfw86zfa4vcxms9l"))' $ME_NODE1_HOME/config/genesis.json
# 检查模块账户余额
#med q bank balances $(med q auth module-account stake_tokens_pool -o json --home $ME_NODE1_HOME | jq -r '.base_account.address') --home $ME_NODE1_HOME


#质押
echo "med gentx global_dao "${STAKING_AMOUNT}" --chain-id "${ME_CHAIN_ID}" --keyring-backend ${KEYRING_BACKEND_NAME} --region-id ${ME_STAKING_AMOUNT_NODE1_REGION} --validator-address "${USER_ADDR_ME_VAL1}" --home "${ME_NODE1_HOME}""
med gentx global_dao "${STAKING_AMOUNT}" --chain-id "${ME_CHAIN_ID}" --keyring-backend ${KEYRING_BACKEND_NAME} --region-id ${ME_STAKING_AMOUNT_NODE1_REGION}  --validator-address "${USER_ADDR_ME_VAL1}" --home "${ME_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                              收集创世交易                                     #"
echo "# ---------------------------------------------------------------------------- #"

med collect-gentxs --home "${ME_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#   设定可添加白名单的用户-超管(global_dao、sequencer、sequencer2、sequencer3)    #"
echo "# ---------------------------------------------------------------------------- #"

set_authorised_deployer_account "${USER_ADDR_ME_GLOBAL_DAO}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                配置 kyc issuers                                               #"
echo "# ---------------------------------------------------------------------------- #"

set_kyc_issuers

echo "# ---------------------------------------------------------------------------- #"
echo "#                              校验创世交易                                     #"
echo "# ---------------------------------------------------------------------------- #"

med validate-genesis --home "${ME_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                              ME同步节点配置                                   #"
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
echo "#                              ME读取Node ID                                     #"
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
echo "#                         修改Node2-Node5的persistent_peers地址                #"
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
echo "#                          生成 docker-compose 文件                            #"
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
    - /root/.mechain
    - --minimum-gas-prices
    - "0.02umec"

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
  command:
    - start
    - --home
    - /root/.rollapp

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
      - ${ME_NODES_HOME}/node1:/root/.mechain

  hub-node2:
    <<: *me-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${ME_NODE2_IP}
    volumes:
      - ${ME_NODES_HOME}/node2:/root/.mechain

  hub-node3:
    <<: *me-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${ME_NODE3_IP}
    volumes:
      - ${ME_NODES_HOME}/node3:/root/.mechain

  hub-node4:
    <<: *me-template
    networks:
      ${Docker_Network_Name}:
        ipv4_address: ${ME_NODE4_IP}
    volumes:
      - ${ME_NODES_HOME}/node4:/root/.mechain
EOF
docker compose up -d hub-node1
echo "等待5s服务完全启动"
sleep 5
docker compose up -d hub-node2
echo "等待5s服务完全启动"
sleep 5
docker compose up -d hub-node3
echo "等待5s服务完全启动"
sleep 5
docker compose up -d hub-node4
echo "等待5s服务完全启动"
sleep 5

echo "# ---------------------------------------------------------------------------- #"
echo "#                  修改node1 persistent_peers地址                              #"
echo "# --------------------- ------------------------------------------------------- #"
docker compose  down hub-node1
sed -i 's/persistent_peers = ""/persistent_peers = "'"${NodeID_Node1_P2P}"'"/' ${ME_NODES_HOME}/node1/config/config.toml
echo NodeID_Node1_P2P:${NodeID_Node1_P2P}
docker compose  up -d hub-node1
echo "等待5s服务完全启动"
sleep 5
echo "# ---------------------------------------------------------------------------- #"
echo "#                  当前容器状态                                              #"
echo "# ---------------------------------------------------------------------------- #"
docker compose  ps

echo "# ---------------------------------------------------------------------------- #"
echo "#               判断val2、val3、val4是否需要add                                  #"
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
echo "#                  等待ME出创世区块块 - 获取区块 hash                              #"
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
echo "#                  获取ME的Node2、Node3、Node4公钥                                   #"
echo "# ---------------------------------------------------------------------------- #"
ME_Node1_Pubkey=$(med tendermint show-validator --home ${ME_NODES_HOME}/node1)
ME_Node2_Pubkey=$(med tendermint show-validator --home ${ME_NODES_HOME}/node2)
ME_Node3_Pubkey=$(med tendermint show-validator --home ${ME_NODES_HOME}/node3)
ME_Node4_Pubkey=$(med tendermint show-validator --home ${ME_NODES_HOME}/node4)
echo Node1公钥:${ME_Node1_Pubkey}
echo Node2公钥:${ME_Node2_Pubkey}
echo Node3公钥:${ME_Node3_Pubkey}
echo Node4公钥:${ME_Node4_Pubkey}

# 离线交易
# 因为正常的在线交易1个区块只能发送一笔交易, 必须等交易上链后才能发送下一笔交易
# 在需要连续发送交易的场景时，可以使用离线交易
# 离线交易的关键在于自己维护和设定 发送交易的序号
#-a $account_number   # 当前地址在链上的地址账户编号id
#--sequence ${seqID}  # 当前地址在链上的交易序号, 每次交易需要+1，防止重放攻击
#--offline 离线模式签名
#--fees交易费用
#--gas设置gas限制

#staking create-validator
##用于创建新的验证人节点
##需要质押代币作为保证金
##节点会参与网络共识和区块生产
##需要满足最低质押要求
##可以获得出块奖励和交易费用
#staking new-region
##用于在现有验证人节点上添加新的区域/地区
##不需要额外的质押
##主要目的是提高网络的去中心化程度
##通常用于在不同地理位置部署节点
##继承主验证人的质押和奖励
##简单来说,create-validator是创建全新的验证人,而new-region是在现有验证人基础上扩展节点到新的地理位置。

if [[ $Interactive_Mode =~ ^[Yy]$ ]]
then
    echo "# -------------------------------------------------------------------------------- #"
    echo "#      资管创建验证者模式                                                            #"
    echo "# -------------------------------------------------------------------------------- #"
    echo "请提交ME的Node1、Node2、Node3、Node4公钥给资管"
    echo "资管操作后请确认回车"
    read -p "Do you want to continue? (Y/n)"
else
    echo "# -------------------------------------------------------------------------------- #"
    echo "#      自动创建验证者模式                                                            #"
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
echo "#               查询验证者节点列表 regionID                                    #"
echo "# ---------------------------------------------------------------------------- #"
while [ $(med query staking validators --home ${ME_NODE1_HOME} | grep -c "regionID") -ne 4 ]; do
    echo "4个验证节点regionID未就绪,等待5秒再循环查询"
    sleep 5
    echo "等待创建4个regionID完成... 目前已经存在的验证者节点:"
    med query staking validators --home ${ME_NODE1_HOME} --output json | jq -r '.validators[] | .description.regionID'
done
echo "4验证节点已经就绪!"
med query staking validators --home ${ME_NODE1_HOME} --output json | jq -r '.validators[] | .description.regionID'

if [[ $Interactive_Mode =~ ^[Nn]$ ]]
then
    echo "# ---------------------------------------------------------------------------- #"
    echo "#            自建模式等待块高到${ME_CREAT_REGION_HEIGHT},以便创建区时分到足够代币     #"
    echo "# ---------------------------------------------------------------------------- #"
    echo "Need about 指定高度${ME_CREAT_REGION_HEIGHT}*出块速度${ME_TIMEOUT_COMMIT} seconds"
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
    echo "#               创建区new-region zone                                           #"
    echo "# ---------------------------------------------------------------------------- #"
    seqID=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO}   --home ${ME_NODE1_HOME} --output json | jq -r .base_account.sequence)
    validator_count=$(med query staking validators --home ${ME_NODE1_HOME} --output json | jq '.validators | length')
    echo 当前sequence:$seqID
    for ((i=0; i<validator_count; i++)); do
        operator_address=$(med query staking validators  --home ${ME_NODE1_HOME} --output json  | jq -r .validators[$i].operator_address)
        regionID=$(med query staking validators   --home ${ME_NODE1_HOME} --output json | jq -r .validators[$i].description.regionID)
        # 检查区域是否已存在
        if med q staking regions --home ${ME_NODE1_HOME} --output json | jq -e ".region[] | select(.regionId == \"$regionID\")" > /dev/null; then
            echo "Region $regionID already exists, skipping..."
            continue
        fi
        # 创建区
        echo "==============================================="
        echo "create new-region: ${regionID}"
        echo "shell:med tx staking new-region ${regionID} ${operator_address} --from global_dao  --chain-id ${ME_CHAIN_ID} --fees=10000umec --gas 500000 -a $account_number --sequence ${seqID} --offline -y -o json --home $ME_NODE1_HOME| jq -r .txhash"
        txhash=$(med tx staking new-region ${regionID} ${operator_address} --from global_dao  --chain-id ${ME_CHAIN_ID} --fees=10000umec --gas 500000 -a $account_number --sequence ${seqID} --offline -y -o json --home ${ME_NODE1_HOME}| jq -r .txhash)
        echo "查询hash是否上链: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        seqID=$((seqID+1))
    done
fi
echo "# ---------------------------------------------------------------------- ------ #"
echo "#               查询region当前金额                                                 #"
echo "# ---------------------------------------------------------------------------- #"
region_data=$(med query staking regions --home ${ME_NODE1_HOME} -o json)
echo "${region_data}" | jq -r '.region[] | [.regionId, .region_treasure_addr] | @tsv' | while IFS=$'\t' read -r region_id treasure_addr; do
    balance_json=$(med query bank balances "${treasure_addr}" --home ${ME_NODE1_HOME} -o json 2>/dev/null)
    balance=$(echo "${balance_json}" | jq -r '.balances[] | select(.denom=="umec") | .amount // "0"')
    echo "$region_id 当前余额为 ${balance:-0} umec"

done


if [[ $Interactive_Mode =~ ^[Nn]$ ]]
then
    if [[ $INSTALL_TYPE =~ ^[Yy]$ ]]
    then
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${TARGET_BALANCE}umec --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        echo "快速提取模式直接从${ME_STAKING_AMOUNT_NODE1_REGION}提取Send DA所需额度${TARGET_BALANCE}"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${TARGET_BALANCE}umec --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
    else
        echo "# ---------------------------------------------------------------------------- #"
        echo "# 自建模式生产参数需等待块高到${REGION_HIGHT_WAIT},以便所有区有足够代币各自提取(或特定高度创建第5个验证者) #"
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
        echo "从${ME_STAKING_AMOUNT_NODE1_REGION}提取Send DA val1所需额度${ME_SEND_DA_AMOUNT_VAL1}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_SEND_DA_AMOUNT_VAL1} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_SEND_DA_AMOUNT_VAL1} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "查询hash是否上链: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        
        echo "从${ME_STAKING_AMOUNT_NODE1_REGION}提取Send IBC所需额度${ME_SEND_IBC_AMOUNT}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_SEND_IBC_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_SEND_IBC_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "查询hash是否上链: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"

        echo "从${ME_STAKING_AMOUNT_NODE1_REGION}提取预留额度${ME_RESERVE_AMOUNT}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_RESERVE_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_GLOBAL_DAO}   ${ME_RESERVE_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "查询hash是否上链: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"

        echo "从${ME_STAKING_AMOUNT_NODE2_REGION}提取Send DA val2所需额度${ME_SEND_DA_AMOUNT_VAL2}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE2_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL2} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE2_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL2} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "查询hash是否上链: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        
        echo "从${ME_STAKING_AMOUNT_NODE3_REGION}提取Send DA val3所需额度${ME_SEND_DA_AMOUNT_VAL3}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE3_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL3} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE3_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL3} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "查询hash是否上链: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"

        echo "从${ME_STAKING_AMOUNT_NODE4_REGION}提取Send DA val4所需额度${ME_SEND_DA_AMOUNT_VAL4}"
        echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE4_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL4} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE4_REGION}  ${USER_ADDR_ME_GLOBAL_DAO} ${ME_SEND_DA_AMOUNT_VAL4} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
        echo "查询hash是否上链: ${txhash}"
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        echo “当前${USER_ADDR_ME_GLOBAL_DAO}账户信息为：”
        med query bank balances "${USER_ADDR_ME_GLOBAL_DAO}" --home ${ME_NODE1_HOME} -o json 2>/dev/null
    fi
fi

echo "# ---------------------------------------------------------------------- ------ #"
echo "#               部署ME-DA-Node1                                                 #"
echo "# ---------------------------------------------------------------------------- #"

meda-appd init "${DA_NODE1_NAME}" --chain-id="${DA_CHAIN_ID}" --home "${DA_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#   判断用户是否需要add(dao、ibc、val0、val1、val2、val3、val4)，无输出则为不创建      #"
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

#echo "create channel hub-meda时,需要me的ibc地址"
if [ -z "${USER_ADDR_ME_IBC_DA}" ]; then
    da_me_ibc_da=$(add_key "meda-appd"  "me_user_ibc_da" ${USER_SEQ_ME_IBC_DA} "--home ${DA_NODE1_HOME}")
    echo "${da_me_ibc_da}"
    USER_ADDR_ME_IBC_DA=$(echo "${da_me_ibc_da}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_IBC_DA_PUBKEY=$(echo "${da_me_ibc_da}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
#echo "create channel hub-rollapp时,需要me的rollapp地址"
if [ -z "${USER_ADDR_ME_IBC_ROLLAPP}" ]; then
    da_ibc_me=$(add_key "meda-appd"  "da_user_ibc_rollapp" ${USER_SEQ_ME_IBC_ROLLAPP} "--home ${DA_NODE1_HOME}")
    echo "${da_ibc_me}"
    USER_ADDR_ME_IBC_ROLLAPP=$(echo "${da_ibc_me}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ME_IBC_ROLLAPP_PUBKEY=$(echo "${da_ibc_me}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi


echo "# ----------------------------------------------------------------------------------------------------------- #"
echo "#  ME添加6个free-gas-account: ${USER_ADDR_ME_IBC_DA} 、 ${USER_ADDR_ME_IBC_ROLLAPP} 、${USER_ADDR_ME_SEQUENCER} 、${USER_ADDR_ME_SEQUENCER2} 、${USER_ADDR_ME_SEQUENCER3} 、${USER_ADDR_ME_GLOBAL_DAO}#"
echo "# ----------------------------------------------------------------------------------------------------------- #"
seqID=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO} -o json   --home ${ME_NODE1_HOME} | jq -r .base_account.sequence)
me_add_free_gas_account_list() {
    local addrs=("$@") 
    local current_free_gas_address=$(med q dao free-gas-accounts --home ${ME_NODE1_HOME} -o json | jq -c '.addresses')
    echo  current free gas address:${current_free_gas_address}
    for addr in "${addrs[@]}"; do
        echo "当前addr:${addr}"
        if  echo "$current_free_gas_address" | jq --arg addr "${addr}" '. as $arr | $arr | index($addr) != null | not' | grep -q "true"; then
            echo "#当前sequence:$seqID, 查询free-gas-account:${addr} #"
            echo "执行shell:med tx dao free-gas-account "[{\"address\":\"${addr}\", \"is_free\":true}]" --from ${USER_ADDR_ME_GLOBAL_DAO} --home ${ME_NODE1_HOME} --gas 500000 --fees 200000umec --account-number ${account_number} --sequence ${seqID} --offline -y -o json | jq -r '.txhash'"
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
            # 检查交易是否成功
            if [ -z "${txhash}" ] || [ "${txhash}" = "null" ]; then
                echo "错误: 添加免费gas账户 ${addr} 失败"
                med q tx ${txhash} --home ${ME_NODE1_HOME}
                exit 1
            fi
            
            echo "交易哈希: ${txhash}"
            check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
        else
             echo "当前addr已经存在无需添加"
        fi 
        seqID=$((seqID+1))
    done
    # 验证所有账户是否已添加到免费gas列表
    echo "当前免费gas账户列表:"
    med q dao free-gas-accounts --home ${ME_NODE1_HOME}  -o json | jq -c '.addresses'
    echo "----------------------------"
}


me_add_free_gas_account_list ${USER_ADDR_ME_IBC_DA} ${USER_ADDR_ME_IBC_ROLLAPP} ${USER_ADDR_ME_SEQUENCER}  ${USER_ADDR_ME_SEQUENCER2} ${USER_ADDR_ME_SEQUENCER3} ${USER_ADDR_ME_GLOBAL_DAO}




echo "# ---------------------------------------------------------------------------- #"
echo "#      为获取rollapp用户地址提前初始化Rollapp Node1节点-步骤1开始                    #"
echo "#      部署Rollapp Node1 步骤1                                                  #"
echo "# ---------------------------------------------------------------------------- #"


echo "初始化节点..."
rollappd init "${RAPP_NODE1_NAME}" --chain-id "${ROLLAPP_CHAIN_ID}" --home ${RAPP_NODE1_HOME}

echo "# ------------------------------------------------------------------------------------------- #"
echo "#      Rollapp add roluser、ibc、sync_user、dao、dev_operator、sequencer,无输出则为不创建         #"
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

# rollapp 恢复 seq 用户
if [ -z "${USER_ADDR_ROLLAPP_SEQUENCER}" ]; then
    echo "恢复 sequencer 账户..."
    mkdir -p "${RAPP_NODE1_HOME}"/sequencer_keys/keyring-${KEYRING_BACKEND_NAME}
    rollapp_seq=$(add_key  "rollappd"  "sequencer" ${USER_SEQ_ME_SEQUENCER} "--home ${RAPP_NODE1_HOME}/sequencer_keys")
    echo "${rollapp_seq}"
    USER_ADDR_ROLLAPP_SEQUENCER=$(echo "${rollapp_seq}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_ROLLAPP_SEQUENCER_PUBKEY=$(echo "${rollapp_seq}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    echo "sequencer 账户恢复完成"
fi

# 持续查询地址是否有钱
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
        echo 当前用户${user_addr} 拥有 ${balance} ${denom}
        break

    done
}

echo "# ---------------------------------------------------------------------------- #"
echo "#                           DA创世                                           #"
echo "# ---------------------------------------------------------------------------- #"

echo "# ---------------------------------------------------------------------------- #"
echo "# 初始资金${USER_ADDR_DA_VAL0}:${val0_TIA_AMOUNT}、${USER_ADDR_DA_DAO}:${dao_TIA_AMOUNT}  #"
echo "# ---------------------------------------------------------------------------- #"
meda-appd add-genesis-account ${USER_ADDR_DA_VAL0} $val0_TIA_AMOUNT --keyring-backend ${KEYRING_BACKEND_NAME} --home "$DA_NODE1_HOME"
meda-appd add-genesis-account ${USER_ADDR_DA_DAO} $dao_TIA_AMOUNT --keyring-backend ${KEYRING_BACKEND_NAME} --home "$DA_NODE1_HOME"

echo "# ---------------------------------------------------------------------------- #"
echo "#               创世质押val0:${val0_STAKING_AMOUNT}                            #"
echo "# ---------------------------------------------------------------------------- #"

meda-appd gentx val0 ${val0_STAKING_AMOUNT} --chain-id ${DA_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} --gas-prices 0.0001udmec --moniker ${DA_STAKING_AMOUNT_NODE1_REGION}  --home "${DA_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               收集创世交易                                                   #"
echo "# ---------------------------------------------------------------------------- #"

#collect-gentxs会导致错误配置node1上的persistent_peers信息为本机物理ip,故此修改
meda-appd collect-gentxs --home "${DA_NODE1_HOME}"
sed -i 's|^persistent_peers = .*$|persistent_peers = ""|' "${DA_NODES_HOME}/node1/config/config.toml"

echo "# ---------------------------------------------------------------------------- #"
echo "#               校验创世交易                                                   #"
echo "# ---------------------------------------------------------------------------- #"
meda-appd validate-genesis --home "${DA_NODE1_HOME}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               处理app配置文件                                                #"
echo "# ---------------------------------------------------------------------------- #"

sed -i'' -e "/\[api\]/,+3 s/enable = false/enable = true/"  "${DA_NODE1_HOME}"/config/app.toml
sed -i'' -e "/\[api\]/,+10 s/swagger = false/swagger = true/" "${DA_NODE1_HOME}"/config/app.toml
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "${DA_NODE1_HOME}"/config/app.toml
sed -i'' -e "/\[grpc\]/,+3 s/enable = false/enable = true/" "${DA_NODE1_HOME}"/config/app.toml
sed -i'' -e "/\[grpc-web\]/,+5 s/enable = false/enable = true/" "${DA_NODE1_HOME}"/config/app.toml

echo "# ---------------------------------------------------------------------------- #"
echo "#               处理config配置文件                                             #"
echo "# ---------------------------------------------------------------------------- #"
sed -i '/^laddr/s/127.0.0.1/0.0.0.0/' "${DA_NODE1_HOME}"/config/config.toml
#null：不索引任何交易（默认值）kv：使用键值存储索引交易（LevelDB 或 RocksDB）psql：使用 PostgreSQL 数据库索引交易（需要额外配置）。从 null 修改为 kv，表示启用键值存储索引这使得可以通过交易哈希查询交易详情
sed -i '/^indexer/s/null/kv/' "${DA_NODE1_HOME}"/config/config.toml
#允许单个IP多端口也就是不同节点同一机器不同端口运行
sed -i '/^allow_duplicate_ip/s/false/true/' "${DA_NODE1_HOME}"/config/config.toml
#控制地址簿的严格性检查。设置为 false 的影响：允许使用非路由地址（如私有 IP 地址）作为对等节点；允许使用非全局单播地址；
sed -i '/^addr_book_strict/s/true/false/' "${DA_NODE1_HOME}"/config/config.toml
#在验证节点或全节点上，通常建议保持 false，因为某些功能（如状态同步、查询历史状态等）可能需要这些数据，在归档节点上，如果磁盘空间有限，可以考虑设置为 true 来节省空间
sed -i '/^discard_abci_responses/s/true/false/' "${DA_NODE1_HOME}"/config/config.toml
#是 Tendermint 节点配置中的一个参数，用于控制交易追踪的详细程度。noop：默认值，不记录任何追踪信息；local：在本地记录追踪信息。
sed -i '/^trace_type/s/noop/local/'  "${DA_NODE1_HOME}"/config/config.toml
sed -i '/^trace_pull_address/s/.*/trace_pull_address = \":26661\"/' "${DA_NODE1_HOME}"/config/config.toml
#允许跨域
sed -ie 's/cors_allowed_origins.*$/cors_allowed_origins = ["*"]/' "${DA_NODE1_HOME}"/config/config.toml
if [[ "$ENABLE_MONITPORING" == [Yy] ]] ;then
    sed  -i'' -e "s/prometheus = false/prometheus = true/" "${DA_NODE1_HOME}"/config/config.toml
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#               处理client配置文件                                                #"
echo "# ---------------------------------------------------------------------------- #"
meda-appd config chain-id ${DA_CHAIN_ID} --home "${DA_NODE1_HOME}"
meda-appd config keyring-backend ${KEYRING_BACKEND_NAME} --home "${DA_NODE1_HOME}"
meda-appd config node "tcp://${DA_NODE1_IP}:26657" --home "${DA_NODE1_HOME}"
echo "修改/config/client.toml"
sed -i "s/localhost:26657/${DA_NODE1_IP}:26657/" "${DA_NODE1_HOME}"/config/client.toml
# The keyring's backend, where the keys are stored (os|file|kwallet|pass|test|memory)
echo "替换钱包名称"
sed -i'' -e "s/^keyring-backend *= .*/keyring-backend = \"${KEYRING_BACKEND_NAME}\"/" "${DA_NODE1_HOME}"/config/client.toml

echo "# ---------------------------------------------------------------------------- #"
echo "#               处理genesis配置文件                                            #"
echo "# ---------------------------------------------------------------------------- #"

# Override the VotingPeriod from 1 week to 1 minute
sed -i'' 's#"604800s"#"60s"#g' "${DA_GENESIS_FILE}"

# Override the genesis to use app version 1 and then upgrade to app version 2 later.
sed -i'' 's#"app_version": "2"#"app_version": "1"#g' "${DA_GENESIS_FILE}"

# 增加ibc代币生效块高
sed -i "/ibc_deadline_block_height/s/100/${IBC_DEADLINE_BLOCK_HEIGHT}/" "${DA_GENESIS_FILE}"

# 将 dao 地址写入 genesis 文件
tmp=$(mktemp)


sed -i "/max_validators/s/100/${MAX_VALIDATORS}/" "${DA_GENESIS_FILE}"

jq --arg addr ${USER_ADDR_DA_DAO} '.app_state.dao.dao_address = $addr' "${DA_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${DA_GENESIS_FILE}"
#最小质押额
jq --arg amount ${DA_MIN_DEPOSIT_AMOUNT} '.app_state.gov.deposit_params.min_deposit[0].amount = $amount' "${DA_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${DA_GENESIS_FILE}"
#设置验证人佣金率的下限，防止验证人设置过低的佣金率进行不公平竞争
jq --arg var2 "0.000000000000000000" '.app_state.staking.params.min_commission_rate = $var2' "${DA_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${DA_GENESIS_FILE}"
#解绑时间
jq --arg var2 "${DA_UNBONDING_TIME}" '.app_state.staking.params.unbonding_time = $var2' "${DA_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${DA_GENESIS_FILE}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               配置同步其他节点                                             #"
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
echo "#               获取DA Node ID                                               #"
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
echo "#                         修改Node2-Node5的persistent_peers地址                #"
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
echo "#                          生成 docker-compose 文件                            #"
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
echo "等待5s服务完全启动"
sleep 5
docker compose up -d da-node2
echo "等待5s服务完全启动"
sleep 5
docker compose up -d da-node3
echo "等待5s服务完全启动"
sleep 5
docker compose  up -d da-node4
echo "等待5s服务完全启动"
sleep 5
echo "# ---------------------------------------------------------------------------- #"
echo "#                  修改node1 persistent_peers地址                               #"
echo "# ---------------------------------------------------------------------------- #"
docker compose down da-node1
sed -i 's|^persistent_peers = .*$|persistent_peers = "'"${NodeID_Node1_P2P}"'"|' "${DA_NODES_HOME}/node1/config/config.toml"
echo "NodeID_Node1_P2P:${NodeID_Node1_P2P}"
docker compose up -d da-node1
echo "等待5s服务完全启动"
sleep 5
echo "# ---------------------------------------------------------------------------- #"
echo "#                  当前容器状态                                              #"
echo "# ---------------------------------------------------------------------------- #"
docker compose ps


echo "# ---------------------------------------------------------------------------- #"
echo "#                  等待DA出创世区块块 - 获取区块 hash                          #"
echo "# ---------------------------------------------------------------------------- #"
wait_genesis_block_hash meda-appd "--home ${DA_NODE1_HOME}"
#获取当前wait_genesis_block_hash计算值$BLOCK_HASH
TrustedHash=${BLOCK_HASH}
echo "TrustedHash:${TrustedHash}"

echo "# ---------------------------------------------------------------------------- #"
echo "#                  搭建bridge 节点                                             #"
echo "# ---------------------------------------------------------------------------- #"
meda bridge init --p2p.network ${DA_CHAIN_ID} --node.store ${DA_BRIDGE_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#               判断是否需要add  bridge用户                                     #"
echo "# ---------------------------------------------------------------------------- #"

if [ -z "${USER_ADDR_DA_BRIDGE}" ]; then
    meda_bridge=$(add_key  "meda-appd" "bridge" ${da_user_index_bridge}  "--home ${DA_BRIDGE_HOME}/keys")
    echo "${meda_bridge}"
    USER_ADDR_DA_BRIDGE=$(echo "${meda_bridge}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_BRIDGE_PUBKEY=$(echo "${meda_bridge}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
echo "# ---------------------------------------------------------------------------- #"
echo "#               配置文件                                                       #"
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
# 获取桥节点的 TRUSTED_PEERS_LIST

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
echo "#                  搭建full 节点                                             #"
echo "# ---------------------------------------------------------------------------- #"

MEDA_FULL_INIT=$(meda full init --p2p.network ${DA_CHAIN_ID} --node.store ${DA_FULL_HOME})
echo ${MEDA_FULL_INIT}
FULL_CELES_KEY_ADDR=$(echo "${MEDA_FULL_INIT}" | grep -A 1 "ADDRESS:" | tail -n 1 | tr -d '[:space:]')
FULL_CELES_DICT=$(echo "${MEDA_FULL_INIT}" | grep "MNEMONIC" -A 1 | tail -n 1 | sed 's/^[ \t]*//;s/[ \t]*$//')

echo "FULL_CELES_KEY_ADDR: ${FULL_CELES_KEY_ADDR}"
echo "FULL_CELES_DICT: ${FULL_CELES_DICT}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               判断是否add  full用户                                           #"
echo "# ---------------------------------------------------------------------------- #"
if [ -z "${USER_ADDR_DA_FULL}" ]; then
    meda_full=$(add_key "meda-appd"  "full" ${da_user_index_full} "--home ${DA_FULL_HOME}/keys")
    echo "${meda_full}"
    USER_ADDR_DA_FULL=$(echo "${meda_full}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_FULL_PUBKEY=$(echo "${meda_full}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi


echo "# ---------------------------------------------------------------------------- #"
echo "#               Full配置文件                                                    #"
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
echo "#                  搭建light 节点                                             #"
echo "# ---------------------------------------------------------------------------- #"

MEDA_LIGHT_INIT=$(meda light init --p2p.network ${DA_CHAIN_ID} --node.store ${DA_LIGHT_HOME})
echo ${MEDA_LIGHT_INIT}
LIGHT_CELES_KEY_ADDR=$(echo "${MEDA_LIGHT_INIT}" | grep -A 1 "ADDRESS:" | tail -n 1 | tr -d '[:space:]')
LIGHT_CELES_DICT=$(echo "${MEDA_LIGHT_INIT}" | grep "MNEMONIC" -A 1 | tail -n 1 | sed 's/^[ \t]*//;s/[ \t]*$//')

echo "LIGHT_CELES_KEY_ADDR: ${LIGHT_CELES_KEY_ADDR}"
echo "LIGHT_CELES_DICT: ${LIGHT_CELES_DICT}"


echo "# ---------------------------------------------------------------------------- #"
echo "#               判断是否add light用户                                           #"
echo "# ---------------------------------------------------------------------------- #"
if [ -z "${USER_ADDR_DA_LIGHT}" ]; then
    meda_light=$(add_key "meda-appd" "light" ${da_user_index_light} "--home ${DA_LIGHT_HOME}/keys")
    echo "${meda_light}"
    USER_ADDR_DA_LIGHT=$(echo "${meda_light}" | grep -oP 'address: \K[^\s]+')
    USER_ADDR_DA_LIGHT_PUBKEY=$(echo "${meda_light}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
fi
echo "# ---------------------------------------------------------------------------- #"
echo "#               Light配置文件                                                    #"
echo "# ---------------------------------------------------------------------------- #"

sed -i "s/DefaultKeyName .*/DefaultKeyName = \"light\"/" "${DA_LIGHT_HOME}"/config.toml

sed -i "/TrustedHash/s/TrustedHash = \"\"/TrustedHash = \"${TrustedHash}\"/" "${DA_LIGHT_HOME}"/config.toml
# 开启鉴权：Token获取方式：light_token=$(meda light auth --node.store ${DA_LIGHT_HOME} write)
sed -i "s/SkipAuth = false/SkipAuth = true/" "${DA_LIGHT_HOME}"/config.toml

sed -i'' -e "/\[Core\]/,+3 s/IP = .*/IP = \"${DA_NODE1_IP}\"/" "${DA_LIGHT_HOME}"/config.toml
# sed -i'' -e "/\[RPC\]/,+3 s/Port = \"26658\"/Port = \"20002\"/" "${DA_LIGHT_HOME}"/config.toml
sed -i'' -e "/\[RPC\]/,+3 s/Address = \"localhost\"/Address = \"0.0.0.0\"/" "${DA_LIGHT_HOME}"/config.toml

sed -i "/ListenAddresses/s/.*/ListenAddresses =\[\"\/ip4\/0.0.0.0\/udp\/2121\/quic-v1\/webtransport\", \"\/ip4\/0.0.0.0\/udp\/2121\/quic-v1\", \"\/ip4\/0.0.0.0\/udp\/2121\/webrtc-direct\", \"\/ip4\/0.0.0.0\/tcp\/2121\"\]/" "${DA_LIGHT_HOME}"/config.toml

# sed -i '/2121/s/2121/2122/g' "${DA_LIGHT_HOME}"/config.toml

sed -i "s#TrustedPeers = \[\]#TrustedPeers = \[\"${TRUSTED_PEERS_LIST}\"\]#" "${DA_LIGHT_HOME}"/config.toml


#===================================================================
echo "# ---------------------------------------------------------------------------- #"
echo "#               配置da-ful、da-light的docker-compose                            #"
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


echo "延时60秒,否则后续的free-whitelist存在失败情况"
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

        echo "执行shell:meda-appd tx dao add-to-fee-whitelist ${addr} --from dao --chain-id ${DA_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} -a $account_number --sequence ${seqID} --offline -y -o json --home ${DA_NODE1_HOME} | jq -r '.txhash'"
        txhash=$(meda-appd tx dao add-to-fee-whitelist ${addr} --from dao --chain-id ${DA_CHAIN_ID} --keyring-backend ${KEYRING_BACKEND_NAME} -a $account_number --sequence ${seqID} --offline -y -o json --home ${DA_NODE1_HOME} | jq -r '.txhash')
        seqID=$((seqID+1))
        check_tx_status meda-appd $txhash "--home ${DA_NODE1_HOME}"
    done
}

add_addr_to_da_whitelist ${USER_ADDR_DA_IBC} ${USER_ADDR_ME_SEQUENCER} ${USER_ADDR_ME_SEQUENCER2} ${USER_ADDR_ME_SEQUENCER3} ${USER_ADDR_DA_LIGHT} ${USER_ADDR_DA_VAL1} ${USER_ADDR_DA_VAL2} ${USER_ADDR_DA_VAL3} ${USER_ADDR_DA_VAL4}


echo "# ---------------------------------------------------------------------------- #"
echo "#               查询当前白名单地址                                                #"
echo "# ---------------------------------------------------------------------------- #"
meda-appd q dao fee-white-list --home ${DA_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#               初始化rly-da                                                    #"
echo "# ---------------------------------------------------------------------------- #"

rly config init --home ${RLY_DA_NODE_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "#               设置通道配置信息 me-meda                                          #"
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
echo "#               是否同步me_user_ibc_da_name、 da_user_ibc_me_name 密钥          #"
echo "# ---------------------------------------------------------------------------- #"
if [[ "$Sync_Key" == [Yy] ]]; then
    mkdir -p ${RLY_DA_NODE_HOME}/keys/${ME_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    mkdir -p ${RLY_DA_NODE_HOME}/keys/${DA_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    echo "处理me_user_ibc_da_name等价于USER_ADDR_ME_IBC_DA-me_user_ibc_da"
    me_user_ibc_da_name=$(add_key "meda-appd" "me_user_ibc_da_name" ${USER_SEQ_ME_IBC_DA} "--home ${RLY_DA_NODE_HOME}/tmpxx20")
    echo "${me_user_ibc_da_name}"
    USER_ME_IBC_DA_NAME=$(echo "${me_user_ibc_da_name}" | grep -oP 'address: \K[^\s]+')
    USER_ME_IBC_DA_NAME_PUBKEY=$(echo "${me_user_ibc_da_name}" | grep -oP 'pubkey: \K[^\n]+' | sed "s/'//g")
    cp -rf ${RLY_DA_NODE_HOME}/tmpxx20/keyring-${KEYRING_BACKEND_NAME}/*  ${RLY_DA_NODE_HOME}/keys/${ME_CHAIN_ID}/keyring-${KEYRING_BACKEND_NAME}/
    rm -rf ${RLY_DA_NODE_HOME}/tmpxx20

    echo "处理da_user_ibc_me_name等价于USER_ADDR_DA_IBC-ibc"
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
#rly tx link信息同步周期需要高度大于17，否则跨链会出现上述类似错误，但实际无影响，同步信息会延迟到需求高度后再进行同步，故该类错误可忽略
echo "等待60秒"
sleep 60
echo rly tx link hub-meda --home ${RLY_DA_NODE_HOME} --override 
rly tx link hub-meda --home ${RLY_DA_NODE_HOME} --override

echo "# ---------------------------------------------------------------------------- #"
echo "#               设定通道时间                                                     #"
echo "# ---------------------------------------------------------------------------- #"
echo rly transact client ${ME_CHAIN_ID} ${DA_CHAIN_ID} hub-meda --client-tp ${CLIENT_TP} --home ${RLY_DA_NODE_HOME}
rly transact client ${ME_CHAIN_ID} ${DA_CHAIN_ID} hub-meda --client-tp ${CLIENT_TP} --home ${RLY_DA_NODE_HOME}
echo "# ---------------------------------------------------------------------------- #"
echo "#               部署中继器服务                                                    #"
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
echo "#      rly-da验证                                                                  #"
echo "# ---------------------------------------------------------------------------- #"
echo "sleep 5s 等待中继器完全启动"
sleep 5
echo 检查${ME_CHAIN_ID} 客户端状态
rly q clients ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME}
echo 检查${DA_CHAIN_ID}客户端状态
rly q clients ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME}

echo 检查连接状态
rly q connections ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME}


echo "检查${ME_CHAIN_ID} 通道状态 shell:rly q channels ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME}"
rly q channels ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME}
echo "检查${DA_CHAIN_ID} 通道状态 shell:rly q channels ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME}"
rly q channels ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME}
echo "获取${ME_CHAIN_ID} 通道ID shell:rly q channels ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME} 2>/dev/null | jq -r 'select(.state == "STATE_OPEN") | .channel_id'"
ME_CHANNEL_ID=$(rly q channels ${ME_CHAIN_ID} --home ${RLY_DA_NODE_HOME} 2>/dev/null | jq -r 'select(.state == "STATE_OPEN") | .channel_id')
echo "channel_id: ${ME_CHANNEL_ID}"
echo "获取${DA_CHAIN_ID} 通道ID shell:rly q channels ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME} 2>/dev/null | jq -r 'select(.state == "STATE_OPEN") | .channel_id'"
DA_CHANNEL_ID=$(rly q channels ${DA_CHAIN_ID} --home ${RLY_DA_NODE_HOME} 2>/dev/null | jq -r 'select(.state == "STATE_OPEN") | .channel_id')
echo "channel_id: ${DA_CHANNEL_ID}"


# 查询原质押额
echo "查询原质押额shell:meda-appd query staking validators --home ${DA_NODE1_HOME}"
meda-appd query staking validators --home ${DA_NODE1_HOME}
echo "-------------------------------------------"

echo "# ---------------------------------------------------------------------------- #"
echo "#                创建DA验证者                                                   #"
echo "# ---------------------------------------------------------------------------- #"
DaNode1Pubkey=$(meda-appd tendermint show-validator --home ${DA_NODES_HOME}/node1)
DaNode2Pubkey=$(meda-appd tendermint show-validator --home ${DA_NODES_HOME}/node2)
DaNode3Pubkey=$(meda-appd tendermint show-validator --home ${DA_NODES_HOME}/node3)
DaNode4Pubkey=$(meda-appd tendermint show-validator --home ${DA_NODES_HOME}/node4)


if [[ ! $Interactive_Mode =~ ^[Yy]$ ]]
then
    echo "# ---------------------------------------------------------------------------- #"
    echo "#                         非资管模式                                            #"
    echo "#     ME的gold_dao用户 ibc 转账给DA的val1、val2、val3、val4地址用于后续验证者质押              #"
    echo "# ---------------------------------------------------------------------------- #"
    #如果没有初始化金额的话可能需要去除.base_account方能获取值
    account_number=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO} -o json --home ${ME_NODE1_HOME} | jq -r .base_account.account_number)
    seqID=$(med q auth account ${USER_ADDR_ME_GLOBAL_DAO} -o json --home ${ME_NODE1_HOME} | jq -r .base_account.sequence)
    me_ibc_transfer_list() {
        local addr=$1
        local send_amount=$2
        echo "处理ibc: me gold_dao -> da: ${addr} ${send_amount}"
        echo "现在执行:med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${addr} ${send_amount} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=100000umec  --chain-id ${ME_CHAIN_ID}  --keyring-backend ${KEYRING_BACKEND_NAME} --home $ME_NODE1_HOME -a $account_number --sequence ${seqID} --offline -y -o json | jq -r '.txhash' "
        txhash=$(med tx ibc-transfer transfer transfer ${ME_CHANNEL_ID} ${addr} ${send_amount} --from ${USER_ADDR_ME_GLOBAL_DAO} --fees=100000umec  --chain-id ${ME_CHAIN_ID}  --keyring-backend ${KEYRING_BACKEND_NAME} --home $ME_NODE1_HOME -a $account_number --sequence ${seqID} --offline -y -o json | jq -r '.txhash')
        echo "返回txhash: ${txhash}"
        echo "----------------------------"
        seqID=$((seqID+1))
        check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    }
    me_region_transfer_list() {
        local addr=$1
        local region=$2
        local send_amount=$3
        echo "处理ibc: me region -> da: ${addr} ${send_amount}"
        echo "现在执行:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}   $addr $send_amount --from global_dao --home $ME_NODE1_HOME -y -o json | jq -r '.txhash' "
        txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  $addr $send_amount --from global_dao --home $ME_NODE1_HOME -y -o json | jq -r '.txhash')
        echo "返回txhash: ${txhash}"
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
    echo "转账给${USER_ADDR_DA_VAL1}"
    me_ibc_transfer_list ${USER_ADDR_DA_VAL1}   ${ME_SEND_DA_AMOUNT_VAL1}
    echo "等待10s"
    sleep 10
    echo "da-val1增加跨链代币质押额,初始原生代币创始质押${val0_STAKING_AMOUNT}将在DA块hight到达${IBC_DEADLINE_BLOCK_HEIGHT}时销毁导致质押为0,如若不增加该跨链代币额将会自动进入共识节点解绑"
    echo "解绑冷静期期间依旧可以通过da_delegate_validator增加质押额度恢复共识节点。解绑冷静期${DA_UNBONDING_SECONDS}到期后将无法增加质押只能通过da_create_validator重新创建共识节点"
    da_delegate_validator ${DA_STAKING_AMOUNT_NODE1_REGION} ${DA_STAKING_AMOUNT_VAL1} ${USER_ADDR_DA_VAL1}
    echo "override中继器rly后等待5s确认同步hup-meda"
    echo "shell:rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME}"
    rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME} 
    echo "等待5s,也就是通用1个块hight时间"
    sleep 5
    echo "--------------------------------------------------------------------------------------"
    echo "转账给${USER_ADDR_DA_VAL2}"
    me_ibc_transfer_list ${USER_ADDR_DA_VAL2}  ${ME_SEND_DA_AMOUNT_VAL2}
    echo "等待10s"
    sleep 10
    echo "创建验证者da-val2"
    da_create_validator   ${DaNode2Pubkey}  ${DA_STAKING_AMOUNT_VAL2}  ${USER_ADDR_DA_VAL2} ${DA_STAKING_AMOUNT_NODE2_REGION}

    echo "override中继器rly后等待5s确认同步hup-meda"
    echo "shell:rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME}"
    rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME} 
    echo "等待5s,也就是通用1个块hight时间" 
    sleep 5
    echo "--------------------------------------------------------------------------------------"
    echo "转账给${USER_ADDR_DA_VAL3}"
    me_ibc_transfer_list ${USER_ADDR_DA_VAL3}  ${ME_SEND_DA_AMOUNT_VAL3}
    echo "等待10s"
    sleep 10
    echo "创建验证者da-val3"
    da_create_validator   ${DaNode3Pubkey}  ${DA_STAKING_AMOUNT_VAL3}  ${USER_ADDR_DA_VAL3} ${DA_STAKING_AMOUNT_NODE3_REGION}

    echo "override中继器rly后等待5s确认同步hup-meda"
    echo "shell:rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME}"
    rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME} 
    echo "等待5s,也就是通用1个块hight时间" 
    sleep 5
    echo "--------------------------------------------------------------------------------------"
    echo "转账给${USER_ADDR_DA_VAL4}"
    me_ibc_transfer_list ${USER_ADDR_DA_VAL4}  ${ME_SEND_DA_AMOUNT_VAL4}
    echo "等待10s"
    sleep 10
    echo "创建验证者da-val4"
    da_create_validator   ${DaNode4Pubkey}  ${DA_STAKING_AMOUNT_VAL4}  ${USER_ADDR_DA_VAL4} ${DA_STAKING_AMOUNT_NODE4_REGION}

    echo "override中继器rly后等待5s确认同步hup-meda"
    echo "shell:rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME}"
    rly  tx update-clients hub-meda --home ${RLY_DA_NODE_HOME} 
    echo "等待5s,也就是通用1个块hight时间" 
    sleep 5
else
     echo "等待资管操作转账完成后确认回车"
     read -p "Do you want to continue? (Y/n) " 
fi
echo "当前共识节点" 
meda-appd query staking validators --home ${DA_NODE1_HOME} -o json | jq '.validators[].description.moniker'
# 查询质押额
echo "-------------------------------------------"
echo "meda new staking validators data"
meda-appd query staking validators -o json  --home ${DA_NODES_HOME}/node1 | jq .validators[].tokens
echo "-------------------------------------------"

if [[  ${Interactive_Mode} =~ ^[Yy]$ ]]
then
     echo "当前资管模式，请确认上面输出是否有误"
     echo "无误请确认回车"
     read -p "Do you want to continue? (Y/n)"
fi

echo "# ---------------------------------------------------------------------------- #"
echo "#        DA部署完成                                                             #"
echo "# ---------------------------------------------------------------------------- #"

echo "# ---------------------------------------------------------------------------- #"
echo "#      部署Rollapp Node1 步骤2                                               #"
echo "# ---------------------------------------------------------------------------- #"

ROLLAPP_TENDERMINT_CONFIG_FILE="${RAPP_NODE1_HOME}/config/config.toml"
ROLLAPP_CONFIG_FILE="${RAPP_NODE1_HOME}/config/app.toml"
ROLLAPP_DYMINT_FILE="${RAPP_NODE1_HOME}/config/dymint.toml"
ROLLAPP_GENESIS_FILE="${RAPP_NODE1_HOME}/config/genesis.json"
ROLLAPP_CLIENT_FILE="${RAPP_NODE1_HOME}/config/client.toml"


echo "配置节点参数..."
rollappd config chain-id ${ROLLAPP_CHAIN_ID} --home ${RAPP_NODE1_HOME}
rollappd config keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}
rollappd config node "tcp://${RAPP_NODE1_IP}:26657" --home ${RAPP_NODE1_HOME}

# sync me-hub key name
echo 'sync-hub-key-name = "sync_user"' >> ${ROLLAPP_CLIENT_FILE}

echo "# ---------------------------------------------------------------------------- #"
echo "#               处理app配置文件                                                #"
echo "# ---------------------------------------------------------------------------- #"

sed -i'' -e "/\[api\]/,+3 s/enable = false/enable = true/"  "${ROLLAPP_CONFIG_FILE}"
sed -i'' -e "/\[api\]/,+10 s/swagger = false/swagger = true/" "${ROLLAPP_CONFIG_FILE}"
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "${ROLLAPP_CONFIG_FILE}"
sed -i'' -e "/\[grpc\]/,+3 s/enable = false/enable = true/"  "${ROLLAPP_CONFIG_FILE}"
sed -i'' -e "/\[grpc-web\]/,+5 s/enable = false/enable = true/"  "${ROLLAPP_CONFIG_FILE}"
sed -i'' -e 's/^pruning *= .*/pruning = "nothing"/' "${ROLLAPP_CONFIG_FILE}"

echo "# ---------------------------------------------------------------------------- #"
echo "#               处理config配置文件                                             #"
echo "# ---------------------------------------------------------------------------- #"
sed -i '/^laddr/s/127.0.0.1/0.0.0.0/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
sed -i '/^addr_book_strict/s/true/false/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
sed -i 's/cors_allowed_origins.*$/cors_allowed_origins = ["*"]/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
if [[ "$ENABLE_MONITPORING" == [Yy] ]] ;then
    sed  -i'' -e "s/prometheus = false/prometheus = true/" "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
fi


echo "配置创世文件genesis..."
tmp=$(mktemp)

jq --arg denom ${DENOM} '.app_state.mint.params.mint_denom = $denom' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
echo "已更新 mint 参数"
jq --arg denom ${DENOM} '.app_state.staking.params.bond_denom = $denom' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
echo "已更新 staking 参数"
jq --arg denom ${DENOM} '.app_state.gov.deposit_params.min_deposit[0].denom = $denom' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
jq --arg amount ${ROLLAPP_MIN_DEPOSIT_AMOUNT} '.app_state.gov.deposit_params.min_deposit[0].amount = $amount' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
# rollapp 使用 cosmos-sdk v0.46，gov GenesisState 只有 deposit_params/voting_params，没有 params
jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.deposit_params.max_deposit_period = $period' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
jq --arg period "${MAX_DEPOSIT_PERIOD}" '.app_state.gov.voting_params.voting_period = $period' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
jq 'del(.app_state.gov.params)' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
echo "已更新 gov 参数 voting_period/max_deposit_period=${MAX_DEPOSIT_PERIOD} (sdk v0.46 fields only)"

echo "创世文件genesis配置完成"

#echo "验证创世文件配置..."
#jq '.app_state.mint.params.mint_denom' "${ROLLAPP_GENESIS_FILE}"
#jq '.app_state.staking.params.bond_denom' "${ROLLAPP_GENESIS_FILE}"
#jq '.app_state.gov.deposit_params.min_deposit[0].denom' "${ROLLAPP_GENESIS_FILE}"

#===================================================================

# 导出一个模块的公钥 - 节点公钥
rollappd dymint show-sequencer --home ${RAPP_NODE1_HOME} > ${RAPP_NODE1_HOME}/sequencer.info

echo "验证节点公钥..."
cat ${RAPP_NODE1_HOME}/sequencer.info
echo "节点公钥配置完成"
echo  sync me-hub key address
echo "sync-hub-key-address = \"${USER_ADDR_ROLLAPP_SYNC}\"" >> ${ROLLAPP_CLIENT_FILE}

echo "# ---------------------------------------------------------------------------- #"
echo "#  初始化roluser、dao、dev_operator余额                                         #"
echo "# ---------------------------------------------------------------------------- #"
rollappd add-genesis-account roluser "$rollapp_user_roluser_TOKEN_AMOUNT" --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}
rollappd add-genesis-account dao "$rollapp_user_dao_TOKEN_AMOUNT" --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}
rollappd add-genesis-account dev_operator "$rollapp_user_dev_operator_TOKEN_AMOUNT" --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}


rollappd gentx_seq --pubkey "$(rollappd dymint show-sequencer --home ${RAPP_NODE1_HOME})" --from roluser --home ${RAPP_NODE1_HOME}

rollappd gentx roluser "$rollapp_user_roluser_NODE_STAKING_AMOUNT" --chain-id "$ROLLAPP_CHAIN_ID" --keyring-backend ${KEYRING_BACKEND_NAME} --home ${RAPP_NODE1_HOME}

echo "# ---------------------------------------------------------------------------- #"
echo "# 添加5个手续费白名单                                                             #"
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

echo "# 5个手续费白名单已添加完成显示如下 #"
jq '.app_state.region.relayerList' "${ROLLAPP_GENESIS_FILE}"
echo "# 收集创世交易文件 #"
rollappd collect-gentxs --home ${RAPP_NODE1_HOME}

echo "# 验证创世交易 #"
rollappd validate-genesis --home ${RAPP_NODE1_HOME} 


echo "# ---------------------------------------------------------------------------- #"
echo "#      修改 DYMINT_FILE 文件                                                    #"
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

echo "# DYMINT_FILE 文件修改完成 #"
echo "# ---------------------------------------------------------------------------- #"
# 读取 JSON 数据
# base_url = da light addr
# timeout 参数在 DA（数据可用性）层客户端配置中用于设置请求的超时时间，以纳秒（nanoseconds）为单位，50秒（50,000,000,000 纳秒）
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
#设置区块证明生成的最大允许时间，防止证明生成过程无限期运行，影响区块生产的整体性能，生产环境建议根据实际测试调整，如果经常超时，可以适当增加该值
sed -i '/max_proof_time/s/=.*/= "5s"/' ${ROLLAPP_DYMINT_FILE}
#定义将交易批次提交到数据可用性层的超时时间，定义将交易批次提交到数据可用性层的超时时间，防止因网络问题导致的无限期等待
sed -i "/batch_submit_max_time/s/=.*/= \"${BATCH_SUBMIT_MAX_TIME}\"/" ${ROLLAPP_DYMINT_FILE}
#批次最大大小
sed -i '/^block_batch_max_size_bytes/s/500000/1874272/' ${ROLLAPP_DYMINT_FILE}
#定义 Rollup 使用的数据可用性层实现，确保交易数据对网络中的所有参与者可用。"mock"：模拟模式，仅用于测试，"me-da"：表示使用名为 "me-da" 的自定义 DA 层实现与da_config关联。
sed -i "s/da_layer = \"mock\"/da_layer = \"me-da\"/" ${ROLLAPP_DYMINT_FILE}
# 使用 sed 替换 TOML 文件中的 da_config 字段
sed -i "s|da_config = \"\"|da_config = \'$escaped_json\'|" ${ROLLAPP_DYMINT_FILE}

# 替换掉原来的 da light 地址
sed -i "/^da_config/s/http:\/\/.*26658/http\:\/\/${DA_RPC_ADDR}/"  ${ROLLAPP_DYMINT_FILE}

echo "--------------------------------------"
echo "create light token"
light_token=$(meda light auth --node.store ${DA_LIGHT_HOME} write)
echo "light1 token: ${light_token}"

sed -i "/^da_config/s/TOKEN123/${light_token}/" ${ROLLAPP_DYMINT_FILE}

echo "------------------------"
sed -n '/^keyring_home_dir/p' ${ROLLAPP_DYMINT_FILE}
echo "绑定容器内地址非物理机地址"
sed -i '/^keyring_home_dir/s#=.*#= "/root/.rollapp/sequencer_keys"#' ${ROLLAPP_DYMINT_FILE}
sed -n '/^keyring_home_dir/p' ${ROLLAPP_DYMINT_FILE}
#===================================================================

echo "设置DAO和devOperator地址"
echo "DAO地址: ${USER_ADDR_ROLLAPP_DAO}"
echo "devOperator地址: ${USER_ADDR_ROLLAPP_OPERATOR}"

jq --arg dao_address ${USER_ADDR_ROLLAPP_DAO} '.app_state.region.dao.address = $dao_address' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"

jq --arg devOperator ${USER_ADDR_ROLLAPP_OPERATOR} '.app_state.region.devOperator.address = $devOperator' "${ROLLAPP_GENESIS_FILE}" > "$tmp" && mv "$tmp" "${ROLLAPP_GENESIS_FILE}"
sed -i '/timeout_broadcast_tx_commit/s/10/360/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"
echo "开放 26657 + 26656端口"
sed -i '/^laddr/s/127.0.0.1/0.0.0.0/' "${ROLLAPP_TENDERMINT_CONFIG_FILE}"

echo "gas费 0.001 umec"
# sed -i'' -e "s/^minimum-gas-prices *= .*/minimum-gas-prices = \"0$DENOM\"/" "$APP_CONFIG_FILE"
# sed -i '/^minimum-gas-prices/s/=.*/= "0.001urax,0.001umec"/' "$APP_CONFIG_FILE"
sed -i '/^minimum-gas-prices/s/=.*/= "0.001umec"/' "${ROLLAPP_CONFIG_FILE}"


echo "设置API端口"
sed -i -e '/\[api\]/,+3 s/enable *= .*/enable = true/' "${ROLLAPP_CONFIG_FILE}"
sed -i -e '/\[api\]/,+10 s/swagger *= .*/swagger = true/' "${ROLLAPP_CONFIG_FILE}"
sed -i -e '/\[api\]/,+10 s/localhost/0.0.0.0/' "${ROLLAPP_CONFIG_FILE}"


echo "# ---------------------------------------------------------------------------- #"
echo "#      创建rollapp Node1 Docker Compose。                                             #"
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
      - ${RAPP_NODE1_HOME}:/root/.rollapp
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
    echo "从${ME_STAKING_AMOUNT_NODE1_REGION}提取SEQUENCER所需额度${ME_SEND_SEQUENCER1_AMOUNT}"
    echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_SEQUENCER}   ${ME_SEND_SEQUENCER1_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
    txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_SEQUENCER}   ${ME_SEND_SEQUENCER1_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
    echo "查询hash是否上链: ${txhash}"
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    echo "从${ME_STAKING_AMOUNT_NODE1_REGION}提取SEQUENCER2所需额度${ME_SEND_SEQUENCER2_AMOUNT}"
    echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_SEQUENCER2}   ${ME_SEND_SEQUENCER2_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
    txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_SEQUENCER2}   ${ME_SEND_SEQUENCER2_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
    echo "查询hash是否上链: ${txhash}"
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
    echo "从${ME_STAKING_AMOUNT_NODE1_REGION}提取SEQUENCER3所需额度${ME_SEND_SEQUENCER3_AMOUNT}"
    echo "shell:med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_SEQUENCER3}   ${ME_SEND_SEQUENCER3_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash'"
    txhash=$(med tx staking withdraw-from-region ${ME_STAKING_AMOUNT_NODE1_REGION}  ${USER_ADDR_ME_SEQUENCER3}   ${ME_SEND_SEQUENCER3_AMOUNT} --from global_dao --home ${ME_NODE1_HOME} -y -o json | jq -r '.txhash')
    echo "查询hash是否上链: ${txhash}"
    check_tx_status med ${txhash} "--home ${ME_NODE1_HOME}"
else
    echo "提交资管以下参数"
    echo "所属Rollapp:${ROLLAPP_CHAIN_ID}"
    echo "最大排数:${MAX_SEQUENCERS}"
    echo "sequence地址:${USER_ADDR_ME_SEQUENCER}、${USER_ADDR_ME_SEQUENCER2}、${USER_ADDR_ME_SEQUENCER3}"
    read -rp "等候资管执行rollapp create-rollapp后回车" 
    med q rollapp list --home ${ME_NODE1_HOME}
    read -rp "排序器列表是否一致,一致则回车,否则请资管重新执行rollapp create-rollapp"
    echo "${USER_ADDR_ME_SEQUENCER}账户信息如下"
    med query bank balances ${USER_ADDR_ME_SEQUENCER} --home ${ME_NODE1_HOME}-o json
    echo "${USER_ADDR_ME_SEQUENCER}账户信息如下"
    med query bank balances ${USER_ADDR_ME_SEQUENCER2} --home ${ME_NODE1_HOME}-o json
    echo "${USER_ADDR_ME_SEQUENCER}账户信息如下"
    med query bank balances ${USER_ADDR_ME_SEQUENCER3} --home ${ME_NODE1_HOME}-o json
    read -rp "请确认上面是否存在足够的余额,否则需要资管转账到对应账户,金额应该大于${SEQUENCER_AMOUNT},回车继续,余额不足请确认转账后按N" REPLY
    if [[ $REPLY =~ ^[Nn]$ ]]
    then
        echo "${USER_ADDR_ME_SEQUENCER}账户信息如下"
        med query bank balances ${USER_ADDR_ME_SEQUENCER} --home ${ME_NODE1_HOME}-o json
        echo "${USER_ADDR_ME_SEQUENCER2}账户信息如下"
        med query bank balances ${USER_ADDR_ME_SEQUENCER2} --home ${ME_NODE1_HOME}-o json
        echo "${USER_ADDR_ME_SEQUENCER3}账户信息如下"
        med query bank balances ${USER_ADDR_ME_SEQUENCER3} --home ${ME_NODE1_HOME}-o json
    fi
    read -rp "请确认上面是否存在足够的余额" 
fi
echo "me 绑定排序器bind sequencer"
echo "获取Node1排序器公钥"
echo "sequencer.info 内容:"
cat ${RAPP_NODE1_HOME}/sequencer.info
#正式环境 - 质押为 0
#注册排序器
#排序器地址需要在HUB上有umec,需要资管定期拨款 @资管,创世后先拨款1MEC,理论是可以发10000笔交易

echo "准备绑定排序器..."
create_sequencer ${SEQUENCER_MONIKER_NAME1} node1 sequencer

echo "排序器状态shell:med q sequencer list-sequencer --home ${ME_NODE1_HOME}"
med q sequencer list-sequencer --home ${ME_NODE1_HOME}
echo "-------------------------------------------------------------------------------"

docker compose  up -d rollapp-node1
echo "等待5s服务完全启动"
sleep 5
echo "# ---------------------------------------------------------------------------- #"
echo "#                  当前容器状态                                                  #"
echo "# ---------------------------------------------------------------------------- #"
docker compose  ps


setup_after_rollapp_node1

#测试语句
#提取金库的钱到global_dao
#./bin/med  tx staking withdraw-from-region me_earth  me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t  10mec --from global_dao --home nodes/hub-nodes/node1 -y -o json
#转账到rollapp ibc并查询rollapp ibc余额
#./bin/med query bank balances me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t --home nodes/hub-nodes/node1 -o json
#rollappd query bank balances me1qk7jwgdmq5lfraq5l4p3qa5vq9vwa84rzhmhhc  --home nodes/rollapp -o json
#./bin/med tx ibc-transfer transfer transfer channel-1 me1qk7jwgdmq5lfraq5l4p3qa5vq9vwa84rzhmhhc 1mec --from me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t --fees=100000umec  --chain-id mechain_400-1 --home /data/docker-dev/nodes/hub-nodes/node1 --keyring-backend test  -y -o json
#./bin/rollappd query bank balances me1qk7jwgdmq5lfraq5l4p3qa5vq9vwa84rzhmhhc --home nodes/rollapp -o json

#转账到ibc并查询ibc余额
#./bin/meda-appd query bank balances me1zps5zt4tvz8uu4ajuq59sdjj2cf48w0nrv0gvd --home nodes/da-nodes/node1
#./bin/med tx ibc-transfer transfer transfer channel-0 me1zps5zt4tvz8uu4ajuq59sdjj2cf48w0nrv0gvd 2mec --from me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t --fees=100000umec  --chain-id mechain_400-1  --keyring-backend test  -y -o json --home nodes/hub-nodes/node1
#./bin/med query tx 9EF49DA3EA7EAE9D1E2CF25515A5934292CF4289CA6B401E7F56065ED7A5BBD4  --home nodes/hub-nodes/node1 -o json
#./bin/meda-appd query bank balances me1zps5zt4tvz8uu4ajuq59sdjj2cf48w0nrv0gvd --home nodes/da-nodes/node1


#查询ibc状态
#trusting_period额度为511200s，此时unbonding_period为604800s
#./bin/med q ibc client state 07-tendermint-0 --home nodes/hub-nodes/node1 -o json
#./bin/med tx ibc client update 07-tendermint-0 --from <your-key> --chain-id mechain_400-1 --keyring-backend test -y
#查看 ibc 客户端有效期
#保持不过期的办法:
#1. 即将过期前发送一笔ibc交易进行续期；
#2. IBC客户端内置的刷新时间机制
#./bin/rly query clients-expiration hub-rollapp --home ./nodes/rly-rapp  -o json
#./bin/rly query clients-expiration hub-meda   --home ./nodes/rly-da -o json
#更新操作
#./bin/rly  tx update-clients hub-meda --home  ./nodes/rly-da
#./bin/rly  tx update-clients hub-rollapp --home ./nodes/rly-rapp

#Region查询余额并提现到指定账户
#ME_NODE1_HOME="nodes/hub-nodes/node1" 
#region_data=$(./bin/med  query staking regions --home ${ME_NODE1_HOME} -o json)
#echo "$region_data" | jq -r '.region[] | [.regionId, .region_treasure_addr] | @tsv' | while IFS=$'\t' read -r region_id treasure_addr; do
#    balance_json=$(./bin/med  query bank balances "$treasure_addr" --home ${ME_NODE1_HOME} -o json 2>/dev/null)
#    balance=$(echo "$balance_json" | jq -r '.balances[] | select(.denom=="umec") | .amount // "0"')
#    echo "$region_id 当前余额为 ${balance:-0} umec"
#
#done
#echo "me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t的余额为:"
#./bin/med query bank balances me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t --home nodes/hub-nodes/node1 -o json
#./bin/med tx staking withdraw-from-region me_earth me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t  200umec --from global_dao --home nodes/hub-nodes/node1 -y -o json | jq -r '.txhash'
#./bin/med query tx 65BF86C4AFD1E6C5D281A13D21B089E6620956F6EA6F6BD19A65419A1611F4C3 --home nodes/hub-nodes/node1 -o json


# 清理节点全部数据, 回到创世时的第0个块状态
# 基本没用, 只有在整个链需要从 1 开始运行，且不需要重新创世时
# ./bin/med tendermint unsafe-reset-all --home ./nodes/node1


# 查询验证者节点列表
# ./bin/med query staking validators --home ./nodes/node1

# 导出用户私钥
#./bin/med keys export user_name --unsafe --unarmored-hex --keyring-backend test --home nodes/hub-nodes/node1/ | xxd -r -p | base64
# 导入用户私钥
#./bin/med keys import <name> <keyfile> --home nodes/hub-nodes/node1/ --keyring-backend test


## 边界网络代理节点关键配置
# 开关 p2p 的交换机, 即不会再共享自己的 p2p 地址本，做到了网络隔离
# 执行 sed -i 's/pex = true/pex = false/' "${ME_NODE1_HOME}"/config/config.toml
# 外部区域节点连接此哨兵节点时，必须通过 persistent_peers 长连接才可以

#合约升级
# 1.  上传新版本代码
#med tx wasm store artifacts/matching_contract_v2.wasm --from admin --gas auto --gas-adjustment 1.3 --chain-id your-chain-id -y
# 记录新的 CODE_ID
# 2.  执行迁移
#med tx wasm migrate $CONTRACT_ADDR <NEW_CODE_ID> '{"new_version":"1.1.0"}' --from admin --chain-id your-chain-id -y

