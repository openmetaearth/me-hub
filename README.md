# ME Hub 
![license](https://img.shields.io/github/license/openmetaearth/me-hub)
![Go](https://img.shields.io/badge/go-1.23-blue.svg)
![issues](https://img.shields.io/github/issues/openmetaearth/me-hub)

## Overview

Welcome to the ME Hub, the **Settlement Layer of the ME protocol**.

This guide will walk you through the steps required to set up and run a ME Hub full node.

## Table of Contents

- [ME Hub](#ME-hub)
  - [Region governance transactions](#region-governance-transactions)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Initializing `med`](#initializing-med)
  - [Running the Chain](#running-the-chain)
  - [Bootstrapping liquidity pools](#bootstrapping-liquidity-pools)
  - [Adding incentives](#adding-incentives)
    - [Creating incentives streams](#creating-incentives-streams)
    - [Locking tokens](#locking-tokens)
    - [check rewards](#check-rewards)
  - [Debugging Container](#debugging-container)
  - [Developer](#developer)
    - [Setup push hooks](#setup-push-hooks)

## Region governance transactions

### Binary integration tests

The scripts in `scripts/regiongov/` invoke the installed `med` binary directly;
none passes `--home`, so your default local home and client configuration are used.
They require Bash and `jq`. Edit `scripts/regiongov/config.sh` or override its
values through environment variables. Defaults are `DAO_KEY=global_dao`,
`OPERATOR_KEY=chn`, and `REGION_ID=chn`. Set `EQUITY_ADDRESSES` to at least two
space-separated test account addresses. `MED_BIN` can select another prebuilt
binary; scripts never compile or install it.

```bash
# Read-only queries and offline binary validation (no broadcast).
./scripts/regiongov/queries.sh
./scripts/regiongov/negative.sh
./scripts/regiongov/queries.sh proposal 1
./scripts/regiongov/queries.sh action 1
# Pass any other regiongov query and its flags directly.
./scripts/regiongov/queries.sh proposals --status rejected --kind treasury

# List independent test functions (also the default when no arguments are given).
./scripts/regiongov/run.sh --list
# Check or execute a single feature, without running other test cases.
./scripts/regiongov/run.sh --check test_freeze_region
./scripts/regiongov/run.sh --run test_freeze_region
EQUITY_ADDRESSES="me1... me1..." ./scripts/regiongov/run.sh --run test_submit_freeze_equity
# Execute selected cases in the supplied order, or explicitly run everything.
./scripts/regiongov/run.sh --run test_submit_treasury_proposal test_vote
EQUITY_ADDRESSES="me1... me1..." ./scripts/regiongov/run.sh --run all
```

Each `test_*` function is independently runnable and prepares its own prerequisites;
it does not depend on earlier tests. Setup/cleanup may send additional transactions
(for example, testing unfreeze first creates frozen state). Each successful equity
case restores the original membership, region cases unfreeze their targets, and
proposal cases veto their own proposal. Each case has an isolated result directory
and prints `PASS` or `FAIL`; failures stop the selected run without blind rollback.
Region freeze tests need only the DAO key and region IDs, proposal/vote tests do
not require equity addresses, and `test_queries` requires neither signing keys nor
equity configuration. `--run` requires an explicit test name or `all`.

The live suite exercises all 11 transaction commands: equity membership update
and approval, batch freeze rejection and approval, batch unfreeze approval,
immediate DAO equity freeze/unfreeze, region freeze/unfreeze, both proposal
kinds and voting. It checks committed transaction codes and queries resulting
state. Both proposal types are vetoed by the DAO so no withdrawal or operator
replacement executes. This does **not** test successful proposal execution at
voting expiry. Set `SECOND_REGION_ID` to exercise multi-region batches.

`NODE`, `CHAIN_ID`, `KEYRING_BACKEND`, `FEES` or `GAS_PRICES` can override local
settings. `GAS` defaults to `auto`. A broadcast timeout is never automatically
retried, because the transaction may already have committed.

Use a dedicated non-DAO country-region operator with funds for fees and at least
two test equity addresses. All target regions/addresses must initially be unfrozen,
and the test region must have no pending actions. Do not run concurrent writers.
On success the original equity membership is restored, but transaction fees,
proposal/action history, and explicit unfrozen records remain. On failure the
script stops without blind rollback; inspect the printed result directory (which
contains the original membership, transaction hashes and responses), resolve any
pending actions, and restore state manually before rerunning. The test does not
create regions, keys, validators, or fund accounts.


All region governance messages are available under `med tx regiongov`. Append
`--from <key> --chain-id <chain-id>` to each command. The signer supplies the
proposer, voter, or authority; authority operations still require DAO permission
on chain. Add `--generate-only` to inspect an unsigned transaction without broadcasting.

```sh
med tx regiongov submit-proposal proposal.json
med tx regiongov vote <proposal-id> yes # or no; submitting does not cast a vote
med tx regiongov submit-equity-update <region-id> <address1> <address2>
med tx regiongov submit-freeze-equity <region-id> <address> [address ...]
med tx regiongov submit-unfreeze-equity <region-id> <address> [address ...]
med tx regiongov approve-pending-action <action-id>
med tx regiongov reject-pending-action <action-id>
med tx regiongov freeze-region <region-id> [region-id ...]
med tx regiongov unfreeze-region <region-id> [region-id ...]
med tx regiongov freeze-equity <region-id> <address>
med tx regiongov unfreeze-equity <region-id> <address>
```

`freeze-region` and `unfreeze-region` accept one or more region IDs. Their messages
use `region_ids` (a JSON array) instead of the former `region_id` string. Empty
lists, empty IDs, and duplicate IDs (case-insensitive) are rejected. All regions
must exist before any state is changed; each affected region emits its own event.

`submit-equity-update` replaces the entire equity list after approval. Omit all
addresses to request clearing the list. The `submit-*-equity` commands request
approval; `freeze-equity` and `unfreeze-equity` are immediate authority operations.
`submit-freeze-equity` and `submit-unfreeze-equity` accept a nonempty list of
addresses in the same region and create one pending action for the entire batch.
Their message field is `addresses` (a JSON array), replacing `address`. Invalid,
duplicate, non-equity addresses or incompatible freeze states reject the whole
request without creating an action. DAO approval applies the entire batch;
rejection applies none. Each submitted address emits an event with the same action ID.
The immediate DAO commands still accept a single address.

A treasury proposal JSON file supports withdrawals from multiple regions:

```json
{
  "kind": "PROPOSAL_KIND_TREASURY_WITHDRAW",
  "title": "Regional budget",
  "summary": "Withdraw funds for regional operations",
  "voting_period": "3600",
  "treasury": {
    "withdrawals": [
      {
        "region_id": "<region-id>",
        "receiver": "<account-address>",
        "amount": [{"denom": "umec", "amount": "1000000"}]
      }
    ]
  }
}
```

For an operator replacement, use this format instead:

```json
{
  "kind": "PROPOSAL_KIND_REPLACE_REGION_OPERATOR",
  "title": "Replace regional operator",
  "summary": "Appoint a new operator",
  "voting_period": "3600",
  "replace_operator": {
    "region_id": "<region-id>",
    "new_operator": "<validator-operator-address>"
  }
}
```

Replace the placeholders with actual IDs and addresses. `voting_period` is in
seconds. Do not supply both payloads. Any `proposer` in the file is replaced by
`--from`.

## Prerequisites

- [Go (v1.23 or above)](https://go.dev/doc/install)

## Installation

Clone `me-hub`:

```sh
git clone https://github.com/openmetaearth/me-hub.git
cd me-hub
make install
```

Check that the med binaries have been successfully installed:

```sh
med version
```

If the med command is not found an error message is returned,
confirm that your [GOPATH](https://go.dev/doc/gopath_code#GOPATH) is correctly configured by running the following command:

```sh
export PATH=$PATH:$(go env GOPATH)/bin
```

## Initializing `med`

- Using the setup script:

  This method is preferred as it preconfigured to support [running rollapps locally](https://github.com/openmetaearth/rollapp)

  ```sh
  bash scripts/setup_local.sh
  ```

- Manually:

  First, set the following environment variables:

  ```sh
  export CHAIN_ID="mechain_2404-1"
  export KEY_NAME="hub-user"
  export MONIKER_NAME="local"
  ```

  Second, create genesis and init ME chain:

  ```sh
  med init "$MONIKER_NAME" --chain-id "$CHAIN_ID"
  ```

  Third, set parameters to ensure denom is umec:

  ```sh
  bash scripts/set_params.sh
  ```

  Then, add genesis account and provide token to the account:

  ```sh
  med keys add "$KEY_NAME" --keyring-backend test
  med add-genesis-account "$(med keys show "$KEY_NAME" -a --keyring-backend test)" 1000umec
  med gentx "$KEY_NAME" 670mec --chain-id "$CHAIN_ID" --keyring-backend test
  med collect-gentxs
  ```

## Running the Chain

Now start the chain!

```sh
med start
```

You should have a running local node!

## Bootstrapping liquidity pools

To bootstrap the `GAMM` module with pools:

```sh
sh scripts/pools/pools_bootstrap.sh
```

## Adding incentives

### Creating incentives streams

After creating the pools above, we create 2 incentive streams through gov:

```sh
sh scripts/incentives/fund_incentives.sh
```

Wait for the gov proposal to pass, and validate with:

```sh
med q streamer streams
```

### Locking tokens

To get incentives, we need to lock the LP tokens:

```sh
sh scripts/incentives/lockup_bootstrap.sh
```

validate with:

```sh
med q lockup module-balance
```

### check rewards

Every minute a share of the rewards will be distributed!

validate with:

```sh
med q incentives active-gauges

# alternatively, watch the outpup - you will see the "amount" change every minute
#  watch -n1 -d "med q incentives active-gauges --output json | jq '.data[] | { "id": .id, "coins": .coins } '"
```

## Debugging Container

Pre-requisite:
 Install [Docker](https://docs.docker.com/get-docker/)
 Install [VSCode](https://code.visualstudio.com/)
 Install [VSCode Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
 Install [Delve](https://github.com/go-delve/delve)

To debug, you can use the following command to run the debug container:

```sh
make docker-run-debug
```

Then you can run the debugger with the following config for `launch.json` in VSCode:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "ME Debug Container",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "port": 4000,
            "host": "127.0.0.1",
            "debugAdapter": "legacy" // To be remove in the future after https://github.com/golang/vscode-go/issues/3096 is fixed
        }
    ]
}
```

After that, you can run the debugger and set breakpoints in the code.

Example: 

Add breakpoint to `ctx` in `x/eibc/keeper/grpc_query.go` :

```go
func (q Querier) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}
```

Open your browser and go to `http://localhost:1318/openmetaearth/me-hub/eibc/params` and you will see debugger stop and print the value at the breakpoint.

### Setup push hooks

To setup push hooks, run the following command:

```sh
./scripts/setup_push_hooks.sh
```
