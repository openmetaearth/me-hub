#!/bin/sh
set -x
# Common commands
genesis_config_cmds="$(dirname "$0")/src/genesis_config_commands.sh"

if [ -f "$genesis_config_cmds" ]; then
  . "$genesis_config_cmds"
else
  echo "Error: header file not found" >&2
  exit 1
fi

# Set parameters
DATA_DIRECTORY="$HOME/.mechain"
CONFIG_DIRECTORY="$DATA_DIRECTORY/config"
TENDERMINT_CONFIG_FILE="$CONFIG_DIRECTORY/config.toml"
CLIENT_CONFIG_FILE="$CONFIG_DIRECTORY/client.toml"
APP_CONFIG_FILE="$CONFIG_DIRECTORY/app.toml"
GENESIS_FILE="$CONFIG_DIRECTORY/genesis.json"
CHAIN_ID="mechain_100-1"
MONIKER_NAME="local"
KEY_NAME="hub-user"
#MNEMONIC="curtain hat remain song receive tower stereo hope frog cheap brown plate raccoon post reflect wool sail salmon game salon group glimpse adult shift"
MNEMONIC="short public silver age tent next need urge popular asthma rocket reward derive empty captain twenty calm hair flee broccoli plunge vital flavor fragile"
# Setting non-default ports to avoid port conflicts when running local rollapp
SETTLEMENT_ADDR=${SETTLEMENT_ADDR:-"0.0.0.0:36657"}
P2P_ADDRESS=${P2P_ADDRESS:-"0.0.0.0:36656"}
GRPC_ADDRESS=${GRPC_ADDRESS:-"0.0.0.0:8090"}
GRPC_WEB_ADDRESS=${GRPC_WEB_ADDRESS:-"0.0.0.0:8091"}
API_ADDRESS=${API_ADDRESS:-"0.0.0.0:1318"}
JSONRPC_ADDRESS=${JSONRPC_ADDRESS:-"0.0.0.0:9545"}
JSONRPC_WS_ADDRESS=${JSONRPC_WS_ADDRESS:-"0.0.0.0:9546"}

TOKEN_AMOUNT="1000000000000000umec" #10M MEC (1e6mec = 1e7 * 1e8 = 1e15umec )
STAKING_AMOUNT="670000000000000umec" #67% is staked (inflation goal)

# Validate mechain binary exists
export PATH=$PATH:$HOME/go/bin

if ! command -v med; then
  echo "mechain binary not found in $PATH"
  exit 1
fi

# Verify that a genesis file doesn't exists for the mechain chain
if [ -f "$GENESIS_FILE" ]; then
  printf "\n======================================================================================================\n"
  echo "A genesis file already exists. building the chain will delete all previous chain data. continue? (y/n)"
  read -r answer
  if [ "$answer" != "${answer#[Yy]}" ]; then
    rm -rf "$CONFIG_DIRECTORY"
  else
    exit 1
  fi
fi

# Create and init dymension chain
med init "$MONIKER_NAME" --chain-id="$CHAIN_ID"

# ---------------------------------------------------------------------------- #
#                              Set configurations                              #
# ---------------------------------------------------------------------------- #
sed -i'' -e "/\[rpc\]/,+3 s/laddr *= .*/laddr = \"tcp:\/\/$SETTLEMENT_ADDR\"/" "$TENDERMINT_CONFIG_FILE"
sed -i'' -e "/\[p2p\]/,+3 s/laddr *= .*/laddr = \"tcp:\/\/$P2P_ADDRESS\"/" "$TENDERMINT_CONFIG_FILE"

sed -i'' -e "/\[grpc\]/,+6 s/address *= .*/address = \"$GRPC_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[grpc-web\]/,+7 s/address *= .*/address = \"$GRPC_WEB_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[json-rpc\]/,+6 s/address *= .*/address = \"$JSONRPC_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[json-rpc\]/,+9 s/^ws-address *= .*/ws-address = \"$JSONRPC_WS_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e '/\[api\]/,+3 s/enable *= .*/enable = true/' "$APP_CONFIG_FILE"
sed -i'' -e "/\[api\]/,+9 s/address *= .*/address = \"tcp:\/\/$API_ADDRESS\"/" "$APP_CONFIG_FILE"

sed -i'' -e 's/^minimum-gas-prices *= .*/minimum-gas-prices = "0.02umec"/' "$APP_CONFIG_FILE"
sed -i'' -e "s/^chain-id *= .*/chain-id = \"$CHAIN_ID\"/" "$CLIENT_CONFIG_FILE"
sed -i'' -e "s/^keyring-backend *= .*/keyring-backend = \"test\"/" "$CLIENT_CONFIG_FILE"
sed -i'' -e "s/^node *= .*/node = \"tcp:\/\/$SETTLEMENT_ADDR\"/" "$CLIENT_CONFIG_FILE"

set_consenus_params
set_gov_params
set_hub_params
set_misc_params
set_EVM_params
set_bank_denom_metadata
set_epochs_params
set_incentives_params


echo "$MNEMONIC" | med keys add "$KEY_NAME" --recover --keyring-backend test
med add-genesis-account "$(med keys show "$KEY_NAME" -a --keyring-backend test)" "$TOKEN_AMOUNT"
med add-genesis-stake-pool
med add-genesis-m-account

jq '.app_state["dao"]["dao_addresses"]["global_dao"] = "me16sf2d3haq5g58chfd8z8ylveygjycmq8f67hw3"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
jq '.app_state["dao"]["dao_addresses"]["meid_dao"] = "me16sf2d3haq5g58chfd8z8ylveygjycmq8f67hw3"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
jq '.app_state["dao"]["dao_addresses"]["dev_operator"] = "me16sf2d3haq5g58chfd8z8ylveygjycmq8f67hw3"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
jq '.app_state["dao"]["dao_addresses"]["airdrop_address"] = "me16sf2d3haq5g58chfd8z8ylveygjycmq8f67hw3"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"

validator_address=$(med keys show "$KEY_NAME" -a --keyring-backend test)
echo $validator_address
med gentx "$KEY_NAME" "$STAKING_AMOUNT" --chain-id "$CHAIN_ID" --keyring-backend test --region-id me_earth --validator-address "$validator_address"
med collect-gentxs
set_authorised_deployer_account "$(med keys show "$KEY_NAME" -a --keyring-backend test)"
med validate-genesis
