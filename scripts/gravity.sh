#!/bin/bash

RelayerMnemonic="also dune road lumber jeans tiny float pulse escape must wheel gauge"
CHAIN_ID=${CHAIN_ID:-"mechain_100-1"}
KEY_NAME=${KEY_NAME:-"global_dao"}
KEYRING="test"
CHAIN=${CHAIN:-"tron"}
#CHAIN=${CHAIN:-"bsc"}

if [ -z "$CHAIN" ]; then
  echo "Error: CHAIN environment variable is not set." >&2
  exit 1
fi

# cache: r1_address ... r5_address
get_relayers() {
  for i in 1 2 3 4 5; do
    # if not set, get from keyring
    eval "v=\${r${i}_address:-}"
    if [ -z "$v" ]; then
      eval "r${i}_address=\$(med keys show r${i} -a --keyring-backend ${KEYRING})"
    fi
  done
}

init-account() {
  echo "$RelayerMnemonic" | med keys add r1 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --account 1
  echo "$RelayerMnemonic" | med keys add r2 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --account 2
  echo "$RelayerMnemonic" | med keys add r3 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --account 3
  echo "$RelayerMnemonic" | med keys add r4 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --account 4
  echo "$RelayerMnemonic" | med keys add r5 --recover --keyring-backend $KEYRING --key-type eth_secp256k1 --account 5
  get_relayers
  for i in 1 2 3 4 5; do
    eval "addr=\$r${i}_address"
    med tx bank send me139mq752delxv78jvtmwxhasyrycufsvr0mue6u $addr 1000000000mec --from "$KEY_NAME" --keyring-backend=$KEYRING -y --output json --fees=100000umec --gas=300000
    sleep 5
  done
}

proposal-relayers() {
  get_relayers
  relayers_csv=$(printf "%s," "$r1_address" "$r2_address" "$r3_address" "$r4_address" "$r5_address")
  relayers_csv=${relayers_csv%,}
  med tx "$CHAIN" proposal-relayers --relayers "$relayers_csv" --from "$KEY_NAME" --chain-id "$CHAIN_ID" --keyring-backend $KEYRING -y --gas-prices 0.02umec --gas auto --gas-adjustment 1.3
  sleep 5
  med q $"$CHAIN" proposal-relayers
}

bonded-relayer() {
  get_relayers
  for i in 1 2 3 4 5; do
    eval "addr=\$r${i}_address"
    eval "hex_addr=\$(med me-debug addr \$addr | awk -F': *' '/^hex:/ {print \$2}')"
    eval "r${i}_hex=\$hex_addr"
  done
  for i in 1 2 3 4 5; do
    eval "hexv=\$r${i}_hex"
    med tx "$CHAIN" bonded-relayer "$hexv" 100000000umec --from r${i} --chain-id "$CHAIN_ID" --keyring-backend $KEYRING -y --gas-prices 0.02umec --gas 500000
#    sleep 500000
  done
}

add-delegate() {
  get_relayers
  med tx "$CHAIN" add-delegate 100000000umec --from r1 --chain-id "$CHAIN_ID" --keyring-backend $KEYRING -y --gas-prices 0.02umec --gas 500000
  sleep 5
  med q "$CHAIN" relayer $r1_address
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  "$@" || { echo "exec $0 failed: $*"; exit 1; }
fi
