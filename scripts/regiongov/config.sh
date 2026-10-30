#!/usr/bin/env bash
# Environment variables override these local defaults. Never set --home here.
MED_BIN=${MED_BIN:-med}
DAO_KEY=${DAO_KEY:-global_dao}
OPERATOR_KEY=${OPERATOR_KEY:-chn}
REGION_ID=${REGION_ID:-chn}
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
TX_TIMEOUT=${TX_TIMEOUT:-180}
POLL_INTERVAL=${POLL_INTERVAL:-2}
# Optional second region for testing batch region freeze (must initially be unfrozen).
SECOND_REGION_ID=${SECOND_REGION_ID:-}
# Existing validator operator address; defaults to the current operator for veto-only tests.
NEW_OPERATOR=${NEW_OPERATOR:-}
