package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

func (k Keeper) FreezeClientStateForTest(ctx sdk.Context, clientId string) error {
	return k.freezeClientState(ctx, clientId)
}

func (k Keeper) SetClientStateForTest(ctx sdk.Context, clientId string, clientState exported.ClientState) {
	k.ibcClientKeeper.SetClientState(ctx, clientId, clientState)
}

func (k Keeper) GetClientStateForTest(ctx sdk.Context, clientId string) (exported.ClientState, bool) {
	return k.ibcClientKeeper.GetClientState(ctx, clientId)
}
