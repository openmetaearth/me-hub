#!/bin/bash

RelayerMnemonic="also dune road lumber jeans tiny float pulse escape must wheel gauge"
CHAIN_ID=${CHAIN_ID:-"mechain_100-1"}
KEY_NAME=${KEY_NAME:-"global_dao"}

function init-account {
  echo "$RelayerMnemonic" | med keys add r1 --recover --keyring-backend test --account 1
  echo "$RelayerMnemonic" | med keys add r2 --recover --keyring-backend test --account 2
  echo "$RelayerMnemonic" | med keys add r3 --recover --keyring-backend test --account 3
  echo "$RelayerMnemonic" | med keys add r4 --recover --keyring-backend test --account 4
  echo "$RelayerMnemonic" | med keys add r5 --recover --keyring-backend test --account 5
}

function proposal-relayers {
  r1_address=$(med keys show r1 -a --keyring-backend test)
  r2_address=$(med keys show r2 -a --keyring-backend test)
  r3_address=$(med keys show r3 -a --keyring-backend test)
  r4_address=$(med keys show r4 -a --keyring-backend test)
  r5_address=$(med keys show r5 -a --keyring-backend test)
   med tx bsc proposal-relayers --relayers "$r1_address,$r2_address,$r3_address,$r4_address,$r5_address" --from "$KEY_NAME" --chain-id "$CHAIN_ID" --keyring-backend test -y --gas-prices 0.02umec --gas auto --gas-adjustment 1.3
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  "$@" || (echo "exec $0 failed:" "$@" && exit 1)
fi