#!/usr/bin/env bash
# All accounts and queries use ./nodes/hub-nodes/node1.
MED_BIN=${MED_BIN:-med}
DAO_KEY=${DAO_KEY:-global_dao}
OPERATOR_KEY=${OPERATOR_KEY:-val4}
# Only country operators vote: val1 (me_earth) and val2 (experience_region) are excluded.
VOTER_KEYS=${VOTER_KEYS:-"val3 val4"}
REGION_ID=${REGION_ID:-chn}
# Public test identity only. Enter the test mnemonic interactively; never store it here.
EXPECTED_OPERATOR_ADDRESS=${EXPECTED_OPERATOR_ADDRESS:-me1u0s2z64yqdxpnxwjry95z5kqgl89kh90uhkcns}
FUND_AMOUNT=${FUND_AMOUNT:-1000000000umec}
# Required for live tests: two or more existing test account addresses, space separated.
EQUITY_ADDRESSES=${EQUITY_ADDRESSES:-}
# Optional overrides; empty values preserve med's local client configuration.
CHAIN_ID=${CHAIN_ID:-}
NODE=${NODE:-}
KEYRING_BACKEND=${KEYRING_BACKEND:-}
GAS=${GAS:-auto}
GAS_ADJUSTMENT=${GAS_ADJUSTMENT:-1.5}
FEES=${FEES:-}
GAS_PRICES=${GAS_PRICES:-}
# The chain requires at least 10000umec; keep explicit gas-price overrides usable.
if [[ -z $FEES && -z $GAS_PRICES ]]; then FEES=10000umec; fi
TX_TIMEOUT=${TX_TIMEOUT:-180}
POLL_INTERVAL=${POLL_INTERVAL:-2}
# Optional second region for testing batch region freeze (must initially be unfrozen).
SECOND_REGION_ID=${SECOND_REGION_ID:-}
# Existing validator operator address; defaults to the current operator for veto-only tests.
NEW_OPERATOR=${NEW_OPERATOR:-}
