# Set parameters
DATA_DIRECTORY="$HOME/.mechain_e/node/node1"
CONFIG_DIRECTORY="$DATA_DIRECTORY/config"
TENDERMINT_CONFIG_FILE="$CONFIG_DIRECTORY/config.toml"
CLIENT_CONFIG_FILE="$CONFIG_DIRECTORY/client.toml"
APP_CONFIG_FILE="$CONFIG_DIRECTORY/app.toml"
GENESIS_FILE="$CONFIG_DIRECTORY/genesis.json"
CHAIN_ID=${CHAIN_ID:-"mechain_100-1"}
MONIKER_NAME=${MONIKER_NAME:-"local"}

MAIN_DIRECTORY="$HOME/.mechain"
MAIN_GENESIS_FILE="$MAIN_DIRECTORY/config/genesis.json"

# Setting non-default ports to avoid port conflicts when running local rollapp
SETTLEMENT_ADDR=${SETTLEMENT_ADDR:-"0.0.0.0:36657"}
P2P_ADDRESS=${P2P_ADDRESS:-"0.0.0.0:36656"}
GRPC_ADDRESS=${GRPC_ADDRESS:-"0.0.0.0:9090"}
GRPC_WEB_ADDRESS=${GRPC_WEB_ADDRESS:-"0.0.0.0:9091"}
API_ADDRESS=${API_ADDRESS:-"0.0.0.0:2318"}
JSONRPC_ADDRESS=${JSONRPC_ADDRESS:-"0.0.0.0:10545"}
JSONRPC_WS_ADDRESS=${JSONRPC_WS_ADDRESS:-"0.0.0.0:10546"}

# Create and init chain (overwrite if exists)
med init "$MONIKER_NAME" --chain-id="$CHAIN_ID" --home "$DATA_DIRECTORY" --overwrite

# Copy genesis file from main directory
if [ -f "$MAIN_GENESIS_FILE" ]; then
  cp "$MAIN_GENESIS_FILE" "$GENESIS_FILE"
else
  echo "Error: genesis file not found at $MAIN_GENESIS_FILE" >&2
  exit 1
fi

# Query Main Chain Seed Configurations
GET_SEED=$(med tendermint show-node-id)
# Set config.toml configurations seeded from main config.toml
MAIN_CHAIN_SEED="${GET_SEED}@127.0.0.1:26656"
sed -i'' -e "s|^seeds =.*|seeds = \"$MAIN_CHAIN_SEED\"|g" "$TENDERMINT_CONFIG_FILE"

sed -i'' -e "/\[rpc\]/,+3 s/laddr *= .*/laddr = \"tcp:\/\/$SETTLEMENT_ADDR\"/" "$TENDERMINT_CONFIG_FILE"
sed -i'' -e "/\[p2p\]/,+3 s/laddr *= .*/laddr = \"tcp:\/\/$P2P_ADDRESS\"/" "$TENDERMINT_CONFIG_FILE"
sed -i'' -e "/\[grpc\]/,+6 s/address *= .*/address = \"$GRPC_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[grpc-web\]/,+7 s/address *= .*/address = \"$GRPC_WEB_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[json-rpc\]/,+6 s/address *= .*/address = \"$JSONRPC_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[json-rpc\]/,+9 s/^ws-address *= .*/ws-address = \"$JSONRPC_WS_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[api\]/,+9 s/address *= .*/address = \"tcp:\/\/$API_ADDRESS\"/" "$APP_CONFIG_FILE"
med start --home  "$DATA_DIRECTORY".