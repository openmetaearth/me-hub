package kyc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	"github.com/st-chain/me-hub/x/kyc/keeper"
	"github.com/st-chain/me-hub/x/kyc/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	address := k.MustAccAddressFromPubkeyString(genState.Issuer.Pubkey)
	if _, found := k.GetDID(ctx, address); found {
		panic(fmt.Errorf("issuer %s already exists", address))
	}

	k.SetDID(ctx, address, genState.Issuer.Did)
	k.SetDidInfo(ctx, genState.Issuer.Did, genState.Issuer)
	// Set if defined
	if genState.KycEventSeq != nil {
		k.SetKycEventSeq(ctx, *genState.KycEventSeq)
	}
	service := didtypes.Service{
		Sid:         types.ModuleName,
		Name:        types.ModuleName,
		Description: "The KYC verifiable credential issuer based The DID(Decentralized Identity).",
		Issuer:      genState.Issuer.Did,
		Status:      didtypes.SERVICE_STATUS_ACTIVE,
	}
	k.SetService(ctx, service)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	svc, found := k.GetService(ctx)
	if !found {
		return genesis
	}
	// Get all kycEventSeq
	kycEventSeq, found := k.GetKycEventSeq(ctx)
	if found {
		genesis.KycEventSeq = &kycEventSeq
	}
	didInfo, found := k.GetDidInfo(ctx, svc.Issuer)
	genesis.Issuer = didInfo

	return genesis
}
