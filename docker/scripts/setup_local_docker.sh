#!/usr/bin/env bash
# Non-interactive setup script for Docker private network
# This script is based on setup_local.sh but removes all interactive prompts

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
CHAIN_ID=${CHAIN_ID:-"mechain_202404-1"}
MONIKER_NAME=${MONIKER_NAME:-"local"}
KEY_NAME=${KEY_NAME:-"global_dao"}
MNEMONIC="curtain hat remain song receive tower stereo hope frog cheap brown plate raccoon post reflect wool sail salmon game salon group glimpse adult shift"

# Setting non-default ports to avoid port conflicts when running local rollapp
SETTLEMENT_ADDR=${SETTLEMENT_ADDR:-"0.0.0.0:36657"}
P2P_ADDRESS=${P2P_ADDRESS:-"0.0.0.0:36656"}
GRPC_ADDRESS=${GRPC_ADDRESS:-"0.0.0.0:8090"}
GRPC_WEB_ADDRESS=${GRPC_WEB_ADDRESS:-"0.0.0.0:8091"}
API_ADDRESS=${API_ADDRESS:-"0.0.0.0:1318"}
JSONRPC_ADDRESS=${JSONRPC_ADDRESS:-"0.0.0.0:9545"}
JSONRPC_WS_ADDRESS=${JSONRPC_WS_ADDRESS:-"0.0.0.0:9546"}

TOKEN_AMOUNT=${TOKEN_AMOUNT:-"1000000000000000000000000umec"} #1M MEC (1e6mec = 1e6 * 1e18 = 1e24umec )
STAKING_AMOUNT=${STAKING_AMOUNT:-"10000000000000000umec"} #67% is staked (inflation goal)

KEY_NAME_SEQUENCER=${KEY_NAME_SEQUENCER:-"sequencer"}
SEQUENCER_TOKEN_AMOUNT=${SEQUENCER_TOKEN_AMOUNT:-"100000000000000umec"} #1M MEC (1e6mec = 1e6 * 1e18 = 1e24umec )

# Validate mechain binary exists
export PATH=$PATH:/usr/local/bin
if ! command -v med > /dev/null; then
  echo "mechain binary not found in $PATH"
  exit 1
fi

# Clean up existing data if present
if [ -f "$GENESIS_FILE" ]; then
  echo "Removing existing chain data..."
  rm -rf "$DATA_DIRECTORY"
fi

echo "Initializing chain..."
# Create and init dymension chain
med init "$MONIKER_NAME" --chain-id="$CHAIN_ID"

# ---------------------------------------------------------------------------- #
#                              Set configurations                              #
# ---------------------------------------------------------------------------- #
echo "Configuring chain parameters..."

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

# Enable monitoring by default in Docker
echo "Enabling monitoring..."
enable_monitoring

# Initialize AMM accounts by default in Docker
echo "Initializing AMM accounts..."
med keys add pools --keyring-backend test
med keys add user --keyring-backend test

# Add genesis accounts and provide coins to the accounts
med add-genesis-account $(med keys show pools --keyring-backend test -a) 1000000000000000000000000umec
# Give some uatom to the local-user as well
med add-genesis-account $(med keys show user --keyring-backend test -a) 1000000000000000000000umec

# ---------------------------------------------------------------------------- #
#                     Create additional genesis accounts                       #
# ---------------------------------------------------------------------------- #
# Support creating multiple accounts via environment variables
# Format for GENESIS_ACCOUNTS: "account1:amount1,account2:amount2,..."
# Example: GENESIS_ACCOUNTS="test1:1000000000000000000000umec,test2:500000000000000000000umec"
#
# Or use GENESIS_ACCOUNTS_JSON for more control:
# Example: GENESIS_ACCOUNTS_JSON='[{"name":"test1","amount":"1000000000000000000000umec","mnemonic":"..."},{"name":"test2","amount":"500000000000000000000umec"}]'

if [ -n "$GENESIS_ACCOUNTS" ]; then
  echo "Creating additional genesis accounts from GENESIS_ACCOUNTS..."
  echo "$GENESIS_ACCOUNTS" | tr ',' '\n' | while IFS=':' read -r account_name account_amount; do
    # Use default amount if not specified
    account_amount="${account_amount:-1000000000000000000000umec}"

    if [ -n "$account_name" ]; then
      echo "Creating account: $account_name with amount: $account_amount"
      med keys add "$account_name" --keyring-backend test
      med add-genesis-account $(med keys show "$account_name" --keyring-backend test -a) "$account_amount"
    fi
  done
fi

if [ -n "$GENESIS_ACCOUNTS_JSON" ]; then
  echo "Creating additional genesis accounts from GENESIS_ACCOUNTS_JSON..."
  account_count=$(echo "$GENESIS_ACCOUNTS_JSON" | jq '. | length')

  for i in $(seq 0 $((account_count - 1))); do
    account_name=$(echo "$GENESIS_ACCOUNTS_JSON" | jq -r ".[$i].name")
    account_amount=$(echo "$GENESIS_ACCOUNTS_JSON" | jq -r ".[$i].amount // \"1000000000000000000000umec\"")
    account_mnemonic=$(echo "$GENESIS_ACCOUNTS_JSON" | jq -r ".[$i].mnemonic // empty")

    if [ -n "$account_name" ]; then
      echo "Creating account: $account_name with amount: $account_amount"

      if [ -n "$account_mnemonic" ]; then
        echo "Using provided mnemonic for $account_name"
        echo "$account_mnemonic" | med keys add "$account_name" --recover --keyring-backend test
      else
        med keys add "$account_name" --keyring-backend test
      fi

      med add-genesis-account $(med keys show "$account_name" --keyring-backend test -a) "$account_amount"
    fi
  done
fi

echo "Setting up main account..."
echo "$MNEMONIC" | med keys add "$KEY_NAME" --recover --keyring-backend test
med add-genesis-account "$(med keys show "$KEY_NAME" -a --keyring-backend test)" "$TOKEN_AMOUNT"
med add-genesis-stake-pool
med add-genesis-m-accounts
med gentx_DAO --pubkey "$(med keys show "$KEY_NAME" -p)"

echo "Setting up sequencer account..."
med keys add "$KEY_NAME_SEQUENCER" --key-type secp256k1 --keyring-backend test

echo "Configuring DAO addresses..."
jq '.app_state["dao"]["dao_addresses"]["global_dao"] = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
jq '.app_state["dao"]["dao_addresses"]["meid_dao"] = "me1p7s6k4ecrm2kl0rs6399k99pyuk322dc78dcxq"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
jq '.app_state["dao"]["dao_addresses"]["dev_operator"] = "me16qle3emp70kr08wt5508t7gk7trst0zwclnscj"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
jq '.app_state["dao"]["dao_addresses"]["airdrop_address"] = "me1uzt6kk6ra9x0ap3au3xuqwp94l2rnw4zqscn2s"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
jq '.app_state["rollapp"]["params"]["dispute_period_in_blocks"] = "1"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
set_kyc_issuers

validator_address=$(med keys show "$KEY_NAME" -a --keyring-backend test)

echo "Creating genesis transaction..."
med gentx "$KEY_NAME" "$STAKING_AMOUNT" --chain-id "$CHAIN_ID" --keyring-backend test --region-id me_earth --validator-address "$validator_address"
med collect-gentxs

echo "Setting authorized deployer accounts..."
set_authorised_deployer_account "$(med keys show "$KEY_NAME" -a --keyring-backend test)"
set_authorised_deployer_account "$(med keys show "$KEY_NAME_SEQUENCER" -a --keyring-backend test)"

echo "Validating genesis..."
med validate-genesis

echo "Chain initialization completed successfully!"
