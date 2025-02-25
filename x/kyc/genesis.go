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

	var issuers []string

	for _, issuer := range genState.Issuers {
		address, err := k.MustAccAddressFromPubkeyString(issuer.Pubkey)
		if err != nil {
			panic(err)
		}
		if _, found := k.GetDID(ctx, address); found {
			panic(fmt.Errorf("issuer %s already exists", address))
		}

		k.SetDID(ctx, address, issuer.Did)
		k.SetDidInfo(ctx, issuer.Did, issuer)

		issuers = append(issuers, issuer.Did)
	}

	// Set if defined
	if genState.KycEventSeq != nil {
		k.SetKycEventSeq(ctx, *genState.KycEventSeq)
	}

	service := didtypes.Service{
		Sid:         types.ModuleName,
		Name:        types.ModuleName,
		Description: "The KYC verifiable credential issuer based The DID(Decentralized Identity).",
		Issuers:     issuers,
		Status:      didtypes.SERVICE_STATUS_ACTIVE,
	}
	k.SetService(ctx, service)

	// set SBT class
	if err := k.SetSbtClass(ctx); err != nil {
		panic(fmt.Errorf("set SBT class failed, err: %s", err.Error()))
	}
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

	var issuers []didtypes.DidInfo
	for _, issuer := range svc.Issuers {
		didInfo, found := k.GetDidInfo(ctx, issuer)
		if !found {
			panic(fmt.Errorf("get %s issuers info failed", types.ModuleName))
		}
		issuers = append(issuers, didInfo)
	}
	genesis.Issuers = issuers

	return genesis
}
