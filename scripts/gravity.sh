#!/bin/bash

RelayerMnemonic="also dune road lumber jeans tiny float pulse escape must wheel gauge"

function proposal-relayers {
  echo "$RelayerMnemonic" | med keys add "$KEY_NAME" --recover --keyring-backend test
  validator_address=$(med keys show "$KEY_NAME" -a --keyring-backend test)

}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  "$@" || (echo "exec $0 failed:" "$@" && exit 1)
fi