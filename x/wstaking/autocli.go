package wstaking

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
)

// AutoCLIOptions returns the AutoCLI options for the wstaking module.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service:              "metaearth.wstaking.Query",
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// ========== Region Queries ==========
				{
					RpcMethod: "Region",
					Use:       "region [region-id]",
					Short:     "Query a specific region by its ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "regionId"},
					},
				},
				{
					RpcMethod: "AllRegion",
					Use:       "regions",
					Short:     "Query all regions with optional pagination",
				},
				// ========== Delegation Queries ==========
				{
					RpcMethod: "DelegationRewards",
					Use:       "delegation-rewards [delegator-address] [validator-address]",
					Short:     "Query rewards for a specific delegation",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator_address"},
						{ProtoField: "validator_address"},
					},
				},
				{
					RpcMethod: "Delegation",
					Use:       "delegation [delegator-addr] [validator-addr]",
					Short:     "Query delegation info for given validator delegator pair",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator_addr"},
						{ProtoField: "validator_addr"},
					},
				},
				{
					RpcMethod: "AllDelegations",
					Use:       "all-delegations",
					Short:     "Query all delegations with optional pagination",
				},
				{
					RpcMethod: "Stakes",
					Use:       "stakes",
					Short:     "Query all stakes",
				},
				// ========== Fixed Deposit Queries ==========
				{
					RpcMethod: "FixedDepositTotalAmount",
					Use:       "fixed-deposit-total-amount",
					Short:     "Query the total amount of all fixed deposits",
				},
				{
					RpcMethod: "FixedDepositAmountByMeid",
					Use:       "fixed-deposit-amount-by-meid [account]",
					Short:     "Query fixed deposit amount for a specific account",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "account"},
					},
				},
				{
					RpcMethod: "FixedDepositByAcct",
					Use:       "fixed-deposit-by-acct [account] [query-type]",
					Short:     "Query fixed deposits by account address",
					Long:      "Query fixed deposits by account. Query type: 0=all, 1=active, 2=withdrawn",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "account"},
						{ProtoField: "query_type"},
					},
				},
				{
					RpcMethod: "FixedDepositByRegion",
					Use:       "fixed-deposit-by-region [region-id]",
					Short:     "Query fixed deposits by region ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "region_id"},
					},
				},
				{
					RpcMethod: "FixedDeposit",
					Use:       "fixed-deposit [address] [id]",
					Short:     "Query a specific fixed deposit by address and ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "FixedDepositAll",
					Use:       "fixed-deposit-all",
					Short:     "Query all fixed deposits with optional pagination",
				},
				// ========== Fixed Deposit Config Queries ==========
				{
					RpcMethod: "FixedDepositCfg",
					Use:       "fixed-deposit-cfg",
					Short:     "Query fixed deposit configurations for regions",
				},
				{
					RpcMethod: "FixedDepositCfgByTerm",
					Use:       "fixed-deposit-cfg-by-term [region-id] [term]",
					Short:     "Query fixed deposit config by region ID and term",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "regionId"},
						{ProtoField: "term"},
					},
				},
				// ========== Record Queries ==========
				{
					RpcMethod: "QueryAllRecord",
					Use:       "query-all-record",
					Short:     "Query all records with optional pagination",
					Alias:     []string{"records"},
				},
				{
					RpcMethod: "QueryRecordByAddress",
					Use:       "query-record-by-address [account]",
					Short:     "Query records by account address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "account"},
					},
				},
				{
					RpcMethod: "QueryReviewRecordByID",
					Use:       "query-review-record-by-id [action-number]",
					Short:     "Query review record by action number",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "action_number"},
					},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: "metaearth.wstaking.Msg",
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// ========== Staking Transactions ==========
				{
					RpcMethod: "Stake",
					Use:       "stake [validator-address] [amount]",
					Short:     "Stake tokens to a validator",
					Long:      "Stake tokens to a validator. Amount should include denomination (e.g., 1000000ame)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "Unstake",
					Use:       "unstake [validator-address] [amount]",
					Short:     "Unstake tokens from a validator",
					Long:      "Unstake tokens from a validator. Amount should include denomination (e.g., 1000000ame)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "UpdateValidator",
					Use:       "update-validator",
					Short:     "Update an existing validator",
					Long:      "Update validator details including description, commission rate, and addresses",
				},
				{
					RpcMethod: "WithdrawDelegatorReward",
					Use:       "withdraw-rewards [validator-address]",
					Short:     "Withdraw delegation rewards from a validator",
					Alias:     []string{"withdraw-delegator-reward"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
					},
				},
				// ========== Region Transactions ==========
				{
					RpcMethod: "NewRegion",
					Use:       "new-region [name] [operator-address]",
					Short:     "Create a new region",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
						{ProtoField: "operator_address"},
					},
				},
				{
					RpcMethod: "RemoveRegion",
					Use:       "remove-region [region-id]",
					Short:     "Remove an existing region",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "region_id"},
					},
				},
				{
					RpcMethod: "WithdrawFromRegion",
					Use:       "withdraw-from-region [region-id] [receiver] [amount]",
					Short:     "Withdraw funds from a region",
					Long:      "Withdraw funds from region treasury to a receiver address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "region_id"},
						{ProtoField: "receiver"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "WithdrawFromGlobalDaoFeePool",
					Use:       "withdraw-from-global-dao-fee-pool [amount]",
					Short:     "Withdraw from global DAO fee pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "TransferRegion",
					Use:       "transfer-region [from-region] [to-region] [addresses...]",
					Short:     "Transfer addresses from one region to another",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "fromRegion"},
						{ProtoField: "toRegion"},
						{ProtoField: "address", Varargs: true},
					},
				},
				{
					RpcMethod: "IbcTransferFromRegionTreasure",
					Use:       "ibc-transfer-from-region [source-port] [source-channel] [region-id] [token] [timeout-height] [timeout-timestamp]",
					Short:     "IBC transfer from region treasury",
					Long:      "Send IBC transfer from region treasury to another chain",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "source_port"},
						{ProtoField: "source_channel"},
						{ProtoField: "regionId"},
						{ProtoField: "token"},
					},
				},
				// ========== Fixed Deposit Transactions ==========
				{
					RpcMethod: "NewFixedDepositCfg",
					Use:       "new-fixed-deposit-cfg [region-id] [term] [rate]",
					Short:     "Create a new fixed deposit configuration",
					Long:      "Create fixed deposit config. Term is in days, rate is decimal (e.g., 0.05 for 5%)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "regionId"},
						{ProtoField: "term"},
						{ProtoField: "rate"},
					},
				},
				{
					RpcMethod: "SetFixedDepositCfgStatus",
					Use:       "set-fixed-deposit-cfg-status [region-id] [term] [status]",
					Short:     "Set fixed deposit config status",
					Long:      "Set status of fixed deposit config (0=disabled, 1=enabled)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "regionId"},
						{ProtoField: "term"},
						{ProtoField: "status"},
					},
				},
				{
					RpcMethod: "SetFixedDepositCfgRate",
					Use:       "set-fixed-deposit-cfg-rate [region-id] [term] [rate]",
					Short:     "Update fixed deposit config rate",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "regionId"},
						{ProtoField: "term"},
						{ProtoField: "rate"},
					},
				},
				{
					RpcMethod: "RemoveFixedDepositCfg",
					Use:       "remove-fixed-deposit-cfg [region-id] [term]",
					Short:     "Remove a fixed deposit configuration",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "regionId"},
						{ProtoField: "term"},
					},
				},
				{
					RpcMethod: "DoFixedDeposit",
					Use:       "do-fixed-deposit [principal] [term]",
					Short:     "Create a new fixed deposit",
					Long:      "Lock tokens in a fixed deposit. Principal must include denomination (e.g., 1000000ame)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "principal"},
						{ProtoField: "term"},
					},
				},
				{
					RpcMethod: "WithdrawFixedDeposit",
					Use:       "withdraw-fixed-deposit [id]",
					Short:     "Withdraw a matured fixed deposit",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				// ========== Record Transactions ==========
				{
					RpcMethod: "NewRecord",
					Use:       "new-record [action-number] [action-url]",
					Short:     "Create a new record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "actionNumber"},
						{ProtoField: "actionUrl"},
					},
				},
				{
					RpcMethod: "ReviewRecord",
					Use:       "review-record [record-hash] [review-result] [action-number] [reviewed-address]",
					Short:     "Review an existing record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "recordHash"},
						{ProtoField: "reviewResult"},
						{ProtoField: "ActionNumber"},
						{ProtoField: "reviewedAddress"},
					},
				},
				// ========== Validator Management ==========
				{
					RpcMethod: "ReplaceConsensusPubKey",
					Use:       "replace-consensus-pubkey",
					Short:     "Replace validator consensus public key",
					Long:      "Replace the consensus public key for a validator (for key rotation)",
				},
				// ========== Module Transactions ==========
				{
					RpcMethod: "SendToModule",
					Use:       "send-to-module",
					Short:     "Send tokens to a module account",
				},
			},
		},
	}
}
