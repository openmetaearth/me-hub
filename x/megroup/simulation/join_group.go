package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/st-chain/me-hub/x/megroup/keeper"
	"github.com/st-chain/me-hub/x/megroup/types"
)

func SimulateMsgJoinGroup(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgJoinGroup{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the JoinGroup simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "JoinGroup simulation not implemented"), nil, nil
	}
}
