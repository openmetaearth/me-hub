package megroup

import (
	"math/rand"

	"me-hub/testutil/sample"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	megroupsimulation "github.com/st-chain/me-hub/x/megroup/simulation"
	"github.com/st-chain/me-hub/x/megroup/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = megroupsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgCreateGroup = "op_weight_msg_group"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateGroup int = 100

	opWeightMsgUpdateGroup = "op_weight_msg_group"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateGroup int = 100

	opWeightMsgDeleteGroup = "op_weight_msg_group"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteGroup int = 100

	opWeightMsgJoinGroup = "op_weight_msg_join_group"
	// TODO: Determine the simulation weight value
	defaultWeightMsgJoinGroup int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	megroupGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		GroupList: []types.Group{
			{
				Id:      0,
				Creator: sample.AccAddress(),
			},
			{
				Id:      1,
				Creator: sample.AccAddress(),
			},
		},
		GroupCount: 2,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&megroupGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateGroup int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateGroup, &weightMsgCreateGroup, nil,
		func(_ *rand.Rand) {
			weightMsgCreateGroup = defaultWeightMsgCreateGroup
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateGroup,
		megroupsimulation.SimulateMsgCreateGroup(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateGroup int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateGroup, &weightMsgUpdateGroup, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateGroup = defaultWeightMsgUpdateGroup
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateGroup,
		megroupsimulation.SimulateMsgUpdateGroup(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteGroup int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteGroup, &weightMsgDeleteGroup, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteGroup = defaultWeightMsgDeleteGroup
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteGroup,
		megroupsimulation.SimulateMsgDeleteGroup(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgJoinGroup int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgJoinGroup, &weightMsgJoinGroup, nil,
		func(_ *rand.Rand) {
			weightMsgJoinGroup = defaultWeightMsgJoinGroup
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgJoinGroup,
		megroupsimulation.SimulateMsgJoinGroup(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateGroup,
			defaultWeightMsgCreateGroup,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				megroupsimulation.SimulateMsgCreateGroup(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateGroup,
			defaultWeightMsgUpdateGroup,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				megroupsimulation.SimulateMsgUpdateGroup(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteGroup,
			defaultWeightMsgDeleteGroup,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				megroupsimulation.SimulateMsgDeleteGroup(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgJoinGroup,
			defaultWeightMsgJoinGroup,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				megroupsimulation.SimulateMsgJoinGroup(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
