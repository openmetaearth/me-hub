# ME Hub 
![license](https://img.shields.io/github/license/openmetaearth/me-hub)
![Go](https://img.shields.io/badge/go-1.24.11-blue.svg)
![issues](https://img.shields.io/github/issues/openmetaearth/me-hub)

## Overview

Welcome to the ME Hub, the **Settlement Layer of the ME protocol**.

This guide will walk you through the steps required to set up and run a ME Hub full node.

For the v3 / v2.0.18 integration baseline, upgrade ordering, validation results,
and remaining release checks, see the [integration acceptance record](docs/v3-v2.0.18-integration.md).

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

The scripts in `scripts/regiongov/` invoke the installed `med` binary directly.
Account lookups, transactions and chain queries use `--home ./nodes/hub-nodes/node1`
(relative to the working directory); public test-key derivation uses an isolated
keyring under the case's result directory.
They require Bash and `jq`. Edit `scripts/regiongov/config.sh` or override its
values through environment variables. Defaults are `DAO_KEY=global_dao`,
`OPERATOR_KEY=val4`, and `REGION_ID=chn`. Ordinary equity cases derive two public
test addresses (HD indices 21–22) when `EQUITY_ADDRESSES` is empty, or accept at
least two space-separated test addresses. `MED_BIN` can select another prebuilt
binary; scripts never compile or install it.

```bash
# Read-only queries and offline binary validation (no broadcast).
./scripts/regiongov/queries.sh
./scripts/regiongov/negative.sh
# Offline signing-key recovery regressions (isolated public test keyrings only).
bash ./scripts/regiongov/run_test.sh
./scripts/regiongov/queries.sh proposal 1
./scripts/regiongov/queries.sh action 1
# Pass any other regiongov query and its flags directly.
./scripts/regiongov/queries.sh proposals --status rejected --kind treasury

# List commands (also the default when no arguments are given).
./scripts/regiongov/run.sh
# Validators and regions must already exist; setup_4validator.sh creates them.
# Run one case; explicitly unfreeze after the region freeze case.
./scripts/regiongov/run.sh test_freeze_region
./scripts/regiongov/run.sh test_unfreeze_region
EQUITY_ADDRESSES="me1... me1..." ./scripts/regiongov/run.sh test_submit_freeze_equity
./scripts/regiongov/run.sh test_submit_treasury_proposal
./scripts/regiongov/run.sh test_vote
# IND/CHN operators plus DAO vote YES; wait for the on-chain voting period.
./scripts/regiongov/run.sh test_treasury_proposal_pass
# Real operator ballots reach 51%; final DAO NO rejects only at voting expiry.
./scripts/regiongov/run.sh test_treasury_proposal_51pass_daoveto
./scripts/regiongov/run.sh test_treasury_proposal_51pass_daoveto_then_daoyes
./scripts/regiongov/run.sh test_replaceoperator_proposal_51pass_daoveto
# Run successful operator replacement LAST: it changes the key needed later.
./scripts/regiongov/run.sh test_replaceoperator_proposal_51pass_daoveto_then_daoyes
# Alternatively, explicitly opt into the ordered suite (stops on first failure).
# Do not run it immediately after the successful replacement above without
# updating signing keys, funding and the replacement address/index as described below.
# ./scripts/regiongov/run.sh all
```

Each case can be selected independently, provided the configured keys, balances
and starting state are correct. Setup/cleanup may send additional transactions.
Successful equity cases restore the original membership; region freeze/unfreeze
cases leave the requested state (the ordered suite pairs them). The basic submit
and vote cases cast DAO NO and **wait for voting expiry and REJECTED** as cleanup:
NO does not cancel a proposal immediately. Each case has an isolated result
directory and prints `PASS` on success; command errors stop execution without
blind rollback. Region freeze tests need only the DAO key and region IDs;
operator-electorate proposal/vote tests do not require equity addresses, and
`test_queries` requires neither signing keys nor equity configuration. Pass
exactly one command name; no arguments only lists commands.

The live suite exercises all 11 transaction commands: equity membership update
and approval, batch freeze rejection and approval, batch unfreeze approval,
immediate DAO equity freeze/unfreeze, region freeze/unfreeze, both proposal
kinds and voting. It checks committed transaction codes and queries resulting
state. Operator-electorate cases default to `VOTER_KEYS="val3 val4"` (IND and
CHN; me_earth/experience_region are excluded from voting and tally). All voters
must be eligible and funded for transaction fees. Set `SECOND_REGION_ID` to
exercise multi-region batches.

The `*_51pass_daoveto*` cases require `pass_threshold_percent=51` and cast the
smallest whole number of distinct, snapshotted operator YES ballots needed to
reach it; DAO approval does not count toward 51%. DAO YES→NO must leave
`VOTING` / `DAO_VOTE_NO`; `*_then_daoyes` restores `VOTING` / `DAO_VOTE_YES`.
At expiry final NO yields `REJECTED`, whereas final YES with the threshold met
yields `PASSED`. The operator cases also verify no replacement during voting,
then the unchanged operator on rejection or the exact new operator on passage.

`NEW_OPERATOR` must differ from the current on-chain operator. Replacement
proposals are submitted by DAO, so they do not require `OPERATOR_KEY` to match
the current operator; the actual voting keys must still belong to the proposal's
electorate snapshot. By default `recover_test_key` saves a signing key named
`regiongov-operator-23` under `RESULT_DIR/derived-keys` (backend `test`), and reads
its validator address with `keys show --bech val`. The script prints the keyring
path; `new-operator-key.json` contains the public key information. With `all`, the
key is retained in the printed preflight result directory. No node keys are
overwritten. An explicit `NEW_OPERATOR` uses that address without recovering a key.

A successful replacement is **not rolled back**. Retain and fund the new signing
key, update `OPERATOR_KEY` / `VOTER_KEYS` and choose an unused replacement before
another run. HD index 23 is deterministic, not fresh; if it is already used,
select another unused index, for example:

```bash
NEW_OPERATOR_INDEX=26 ./scripts/regiongov/run.sh test_replaceoperator_proposal_51pass_daoveto_then_daoyes
```

Changing the replacement index does not fix stale voting keys: `VOTER_KEYS` must
still contain eligible operators with signing keys in the configured keyring.
Never use these public test identities for real funds. `all` stops on the first
failure and puts successful DAO NO→YES replacement last to preserve earlier
cases' signing keys. The older `test_operator_proposal_pass` remains a standalone
command but is deliberately omitted from `all`: its successful replacement is
covered by that final case, avoiding a second replacement with obsolete keys.

Basic and DAO-override treasury payloads request **10umec for DAO**. The treasury
must already be funded; the runner never changes governance parameters, funds a
treasury or starts/restarts nodes. The ordinary treasury pass, treasury NO→YES
and equity treasury pass cases each withdraw 10umec (30umec total in `all`).
Provide enough spendable funds for sequential cases and the equity percentage
cap; rejected proposals transfer nothing, but treasury rewards may still accrue.
Existing equity treasury cases also fund each signing equity from DAO using
`FUND_AMOUNT` (default `1000000000umec`); those transfers and fees remain.

`NODE`, `CHAIN_ID`, `KEYRING_BACKEND`, `FEES` or `GAS_PRICES` can override local
settings. `GAS` defaults to `auto`; when neither `FEES` nor `GAS_PRICES` is set,
the script supplies `20000umec` in fees. Voting duration uses the current
`min_voting_period` parameter without modifying it; even cleanup waits may take
an hour with a 3600-second period. The wait budget is the voting period plus
`TX_TIMEOUT`, not `TX_TIMEOUT` alone. Allow time for all committed votes and
overwrites. Failed JSON assertions print the actual response,
including transaction error codes and logs. A broadcast timeout is never
automatically retried, because the transaction may already have committed.

Use a dedicated non-DAO country-region operator with funds for fees. All target
regions/addresses must initially be unfrozen, and the test region must have no
pending actions. Do not run concurrent writers. On success transaction fees,
funding transfers, treasury payouts, proposal/action history, explicit unfrozen
records and successful operator replacements remain. On failure inspect the
printed result directory (original state, transaction hashes and responses),
resolve pending actions and any still-voting proposals, and restore state manually
before rerunning. Validators and regions are created by
`scripts/regiongov/setup_4validator.sh`, not this runner. `all` checks the replacement
address/current operator up front, but does not globally preflight every balance
or pending action; configure the prerequisites before running a case.


All region governance messages are available under `med tx regiongov`. Append
`--from <key> --chain-id <chain-id>` to each command. The signer supplies the
proposer, voter, or authority; authority operations still require DAO permission
on chain. Add `--generate-only` to inspect an unsigned transaction without broadcasting.
For both `TreasuryWithdrawPayload` and `ReplaceRegionOperatorPayload`, the DAO
can repeatedly vote `yes` or `no` while chain time is **strictly before**
`voting_end_time`. Each DAO vote overwrites its previous ballot. Neither YES nor
NO ends voting early: the proposal remains `PROPOSAL_STATUS_VOTING`, with
`dao_vote` recording `DAO_VOTE_YES` or `DAO_VOTE_NO`. At the deadline, EndBlocker
uses the DAO's **last ballot**. Final NO (or no DAO approval) rejects the proposal
without executing it. Final YES permits passage only if the counted electorate
meets `pass_threshold_percent`; equity withdrawals also require the snapshotted
region operator's YES. A successful tally executes the payload at expiry; an
execution error yields `PROPOSAL_STATUS_FAILED`, not `PASSED`. Votes at or after
the deadline, or after finalization, are rejected. These overwrite rules apply to
the DAO only; ordinary electors cannot change an already cast ballot.

`electors[].operator`, `proposal.required_operator`, and `MsgVote.voter` use ordinary
account addresses (`me1…`). An operator's ballot account has the same address bytes
as the staking operator. `required_operator` is fixed when the proposal is submitted;
operator replacements do not change that approver. Legacy valoper proposal values
remain readable; proposal writes normalize the approver to an account address.
The v3 binary does not register a v2.0.18 or v2.0.18.rc1 upgrade handler: it expects
mainnet to have completed the v2.0.18 upgrade before upgrading to v3.
The staking `Region.OperatorAddress` and replacement `new_operator` retain
validator operator addresses (`mevaloper1…`).

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
seconds and must be at least the current `min_voting_period` parameter (query it
with `med query regiongov params`). Do not supply both payloads. Any `proposer`
in the file is replaced by `--from`.

## Prerequisites

- [Go (v1.24.11 or above)](https://go.dev/doc/install)

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
