#!/bin/bash

RelayerMnemonic="also dune road lumber jeans tiny float pulse escape must wheel gauge"
#CHAIN_ID=${CHAIN_ID:-"me-chain"}
CHAIN_ID=${CHAIN_ID:-"mechain_900-1"}
KEY_NAME=${KEY_NAME:-"global_dao"}
KEYRING="test"
CHAIN=${CHAIN:-"bsc"}
#NodeUrl=${NodeUrl:-"http://192.168.0.150:26657/"}
NodeUrl=${NodeUrl:-"http://118.175.0.244:26657/"}


if [ -z "$CHAIN" ]; then
  echo "Error: CHAIN environment variable is not set." >&2
  exit 1
fi

# cache: r1_address ... r5_address
get_relayers() {
  echo "$CHAIN"
  for i in 1 2 3 4 5; do
    # if not set, get from keyring
    eval "v=\${r${i}_address:-}"
    if [ -z "$v" ]; then
      eval "r${i}_address=\$(med keys show r${i} -a --keyring-backend ${KEYRING})"
    fi
  done
}

init_account() {
  echo "$RelayerMnemonic" | med keys add r1 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --index 1
  echo "$RelayerMnemonic" | med keys add r2 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --index 2
  echo "$RelayerMnemonic" | med keys add r3 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --index 3
  echo "$RelayerMnemonic" | med keys add r4 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --index 4
  echo "$RelayerMnemonic" | med keys add r5 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --index 5
  get_relayers
  for i in 1 2 3 4 5; do
    eval "addr=\$r${i}_address"
#    med tx bank send global_dao $addr 1000000000umec --from "$KEY_NAME" --keyring-backend=$KEYRING -y --output json --fees=100000umec --gas=300000 --node "$NodeUrl" --chain-id "$CHAIN_ID"
    med tx bank send $addr me1eff5px4606z48lwehyvapane9tc9lekn9c4d5t 10umec --from  r${i} --keyring-backend=$KEYRING -y --output json --fees=100000umec --gas=300000 --node "$NodeUrl" --chain-id "$CHAIN_ID"
#    sleep 5
  done
}

event_nonce() {
   get_relayers
   for i in 1 2 3 4 5; do
     eval "addr=\$r${i}_address"
     echo "relayer address: $addr"
     med q "$CHAIN" event-nonce "$addr" --node "$NodeUrl" --output json | jq .event_nonce
   done
}

proposal_relayers() {
  get_relayers
  relayers_csv=$(printf "%s," "$r1_address" "$r2_address" "$r3_address" "$r4_address" "$r5_address")
  relayers_csv=${relayers_csv%,}
  med tx "$CHAIN" proposal-relayers --relayers "$relayers_csv" --from "$KEY_NAME" --chain-id "$CHAIN_ID" --keyring-backend $KEYRING -y --gas-prices 0.02umec --gas auto --gas-adjustment 1.3 --node "$NodeUrl"
  sleep 5
  med q $"$CHAIN" proposal-relayers --node "$NodeUrl"
}

bonded_relayer() {
  get_relayers
  for i in 1 2 3 4 5; do
    eval "addr=\$r${i}_address"
    eval "hex_addr=\$(med me-debug addr \$addr | awk -F': *' '/^hex:/ {print \$2}')"
    eval "r${i}_hex=\$hex_addr"
  done
  for i in 1 2 3 4 5; do
    eval "hexv=\$r${i}_hex"
    med tx "$CHAIN" bonded-relayer "$hexv" 100000000umec --from r${i} --chain-id "$CHAIN_ID" --keyring-backend $KEYRING -y --gas-prices 0.02umec --gas 500000 --node "$NodeUrl"
#    sleep 500000
  done
}

bonded_relayer_tron() {
  get_relayers
  for i in 1 2 3 4 5; do
    eval "addr=\$r${i}_address"
    eval "hex_addr=\$(med me-debug addr \$addr | awk -F': *' '/^tron:/ {print \$2}')"
    eval "r${i}_hex=\$hex_addr"
  done
  for i in 1 2 3 4 5; do
    eval "hexv=\$r${i}_hex"
    echo "$hexv"
    med tx "$CHAIN" bonded-relayer "$hexv" 100000000umec --from r${i} --chain-id "$CHAIN_ID" --keyring-backend $KEYRING -y --gas-prices 0.02umec --gas 500000 --node "$NodeUrl"
  done
}

add_delegate() {
  get_relayers
    for i in 1 2 3 4 5; do
      eval "hexv=\$r${i}_hex"
    med tx "$CHAIN" add-delegate 1000000umec --from r${i} --chain-id "$CHAIN_ID" --keyring-backend $KEYRING -y --gas-prices 0.02umec --gas 500000 --node $NodeUrl
    done
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  cmd="${1//-/_}"
  shift
  "$cmd" "$@" || { echo "exec $0 failed: $cmd $*"; exit 1; }
fi
