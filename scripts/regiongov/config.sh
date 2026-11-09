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
# Optional: two or more test account addresses, space separated. If empty, run.sh
# derives two public test addresses at HD indices 21-22 (val1-val4 use 11-14).
# test_equity_treasury_proposal_pass needs at least three signing accounts; its
# default uses indices 21, 22, 24 in an isolated temporary test keyring. Explicit
# addresses must have signing keys in the configured keyring. It funds each with
# FUND_AMOUNT from DAO_KEY, withdraws 10umec from chn to DAO_KEY and restores the
# equity list after success. The other equities vote YES; OPERATOR_KEY (the chn
# operator) and DAO_KEY must also vote YES, outside the equity percentage count.
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
# Optional replacement validator operator address (mevaloper...). If empty, run.sh
# derives a public test address at HD index 23, separate from val1-val4/equity.
# Derivation is dry-run: no key is saved. A successful replacement changes the
# operator; recover its test key and update OPERATOR_KEY/VOTER_KEYS for later runs.
# Index 23 is deterministic, not fresh on every run. If already used on-chain,
# explicitly set NEW_OPERATOR to another unused address. Never use for real funds.
NEW_OPERATOR=${NEW_OPERATOR:-}
