---
sidebar_position: 1
---

# `x/wdistri`


## Overview

This _simple_ distribution mechanism describes a functional way to passively
distribute rewards between validators and region delegators. Note that this mechanism does
not distribute funds in as precisely as active reward distribution mechanisms and
will therefore be upgraded in the future.

The mechanism operates as follows. Collected rewards are pooled around the world
Passively assigned to validators and delegators. Tokens are transferred out from the  `treasury_pool` which is a module address, due to the nature of passive statistics, whenever parameters affecting the reward distribution rate change
When this occurs, the allocating of the reward must also occur.

## Structure
The `wdistri` keeper inherits the keeper of the distribution module, but blocks most of the message services and only retains updateParams.


The `wdistri` keeper rewrites the ABCI interface methods implemented by the distribution module, remove `BeginBlocker`, and implement `EndBlocker` to allocate token.

## Effect on Staking

All transaction fees and mining rewards are stored in `treasury_pool`, and rewards are distributed to each region's `RegionTreasure` address every 24 hours. Rewards are distributed according to the staking ratio of regional validator.

## Contents

- [`x/wdistri`](#xwdistri)
  - [Overview](#overview)
  - [Structure](#structure)
  - [Effect on Staking](#effect-on-staking)
  - [Contents](#contents)
  - [Concepts](#concepts)
  - [State](#state)
    - [Params](#params)
  - [Messages](#messages)
    - [MsgUpdateParams](#msgupdateparams)
  - [Hooks](#hooks)
    - [Create or modify validator or region](#create-or-modify-validator-or-region)
      - [Before](#before)
  - [Events](#events)
    - [BeginBlocker](#beginblocker)
    - [EndBlocker](#endblocker)
  - [Client](#client)
    - [CLI](#cli)
      - [Query](#query)
        - [params](#params-1)
        - [rewards](#rewards)
        - [validator-distribution-info](#validator-distribution-info)
    - [gRPC](#grpc)
      - [Params](#params-2)
      - [ValidatorDistributionInfo](#validatordistributioninfo)

## Concepts

In Proof of Stake (PoS) blockchains, rewards gained from transaction fees are paid to validators. The fee distribution module fairly distributes the rewards to the validators' constituent delegators.

Rewards are calculated per period. The period is updated each time a validator's delegation changes, for example, when the validator receives a new delegation.
The rewards for a single validator can then be calculated by taking the total rewards for the period before the delegation started, minus the current total rewards.
To learn more, see the [F1 Fee Distribution paper](https://github.com/cosmos/cosmos-sdk/tree/main/docs/spec/fee_distribution/f1_fee_distr.pdf).

The rewards to a delegator are distributed when the delegation is changed or removed, or a withdrawal is requested.
Before rewards are distributed, all slashes to the validator that occurred during the current delegation are applied.

## State



### Params

The distribution module stores it's params in state with the prefix of `0x09`,
it can be updated with governance or the address with authority.

* Params: `0x09 | ProtocolBuffer(Params)`

```protobuf reference
https://github.com/cosmos/cosmos-sdk/blob/v0.47.0-rc1/proto/cosmos/distribution/v1beta1/distribution.proto#L12-L42
```




## Messages



### MsgUpdateParams

Distribution module params can be updated through `MsgUpdateParams`, which can be done using governance proposal and the signer will always be gov module account address.

```protobuf reference
https://github.com/cosmos/cosmos-sdk/blob/v0.47.0-rc1/proto/cosmos/distribution/v1beta1/tx.proto#L133-L147
```

The message handling can fail if:

* signer is not the gov module account address.

## Hooks

Available hooks that can be called by and from this module.

### Create or modify validator or region

* triggered-by: `staking.MsgNewRegion`, `staking.RemoveRegion`, `staking.MsgCreateValidator`,`staking.EditValidator`,`staking.MsgStake`,`staking.MsgUnStake`

#### Before

* the region create
* the region remove
* the validator create
* the validator was edited
* add stake
* remove stake



## Events

The distribution module emits the following events:

### BeginBlocker

| Type            | Attribute Key | Attribute Value    |
| --------------- | ------------- | ------------------ |
| proposer_reward | validator     | {validatorAddress} |
| proposer_reward | reward        | {proposerReward}   |
| commission      | amount        | {commissionAmount} |
| commission      | validator     | {validatorAddress} |
| rewards         | amount        | {rewardAmount}     |
| rewards         | validator     | {validatorAddress} |

### EndBlocker

| Type                   | Attribute Key           | Attribute Value         |
| ---------------------- | ----------------------- | ----------------------- |
| region_treasury_reword | region_treasury_address | {RegionTreasureAddress} |
| region_treasury_reword | region_id               | {region_id}             |
| region_treasury_reword | amount                  | {amount}                |

## Client

### CLI

A user can query and interact with the `distribution` module using the CLI.

#### Query

The `query` commands allow users to query `distribution` state.

```shell
med query distribution --help
```

##### params

The `params` command allows users to query the parameters of the `distribution` module.

```shell
med query distribution params [flags]
```

Example:

```shell
med query distribution params
```

Example Output:

```yml
base_proposer_reward: "0.000000000000000000"
bonus_proposer_reward: "0.000000000000000000"
community_tax: "0.020000000000000000"
withdraw_addr_enabled: true
```

##### rewards

The `rewards` command allows users to query delegator rewards. Users can optionally include the validator address to query rewards earned from a specific validator.

```shell
med query distribution rewards [delegator-addr] [validator-addr] [flags]
```

Example:

```shell
med query distribution rewards cosmos1...
```

Example Output:

```yml
rewards:
- reward:
  - amount: "1000000.000000000000000000"
    denom: stake
  validator_address: cosmosvaloper1..
total:
- amount: "1000000.000000000000000000"
  denom: stake
```

##### validator-distribution-info

The `validator-distribution-info` command allows users to query validator commission and self-delegation rewards for validator.

```shell
med query distribution validator-distribution-info cosmosvaloper1...
```

Example Output:

```yml
commission:
- amount: "100000.000000000000000000"
  denom: stake
operator_address: cosmosvaloper1...
self_bond_rewards:
- amount: "100000.000000000000000000"
  denom: stake
```

### gRPC

A user can query the `distribution` module using gRPC endpoints.

#### Params

The `Params` endpoint allows users to query parameters of the `distribution` module.

Example:

```shell
grpcurl -plaintext \
    localhost:9090 \
    cosmos.distribution.v1beta1.Query/Params
```

Example Output:

```json
{
  "params": {
    "communityTax": "20000000000000000",
    "baseProposerReward": "00000000000000000",
    "bonusProposerReward": "00000000000000000",
    "withdrawAddrEnabled": true
  }
}
```

#### ValidatorDistributionInfo

The `ValidatorDistributionInfo` queries validator commission and self-delegation rewards for validator.

Example:

```shell
grpcurl -plaintext \
    -d '{"validator_address":"cosmosvalop1..."}' \
    localhost:9090 \
    cosmos.distribution.v1beta1.Query/ValidatorDistributionInfo
```

Example Output:

```json
{
  "commission": {
    "commission": [
      {
        "denom": "stake",
        "amount": "1000000000000000"
      }
    ]
  },
  "self_bond_rewards": [
    {
      "denom": "stake",
      "amount": "1000000000000000"
    }
  ],
  "validator_address": "cosmosvalop1..."
}
```

