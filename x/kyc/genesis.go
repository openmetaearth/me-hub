package kyc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/keeper"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	var issuers []string

	for _, issuer := range genState.Issuers {
		addr := sdk.MustAccAddressFromBech32(issuer.Address)

		if existingDid, found := k.GetDID(ctx, addr); found {
			if existingDid != issuer.Did {
				panic(fmt.Errorf("issuer %s already exists with DID %s, expected %s", addr.String(), existingDid, issuer.Did))
			}
			existingInfo, found := k.GetDidInfo(ctx, issuer.Did)
			if !found {
				panic(fmt.Errorf("issuer %s DID info not found", addr.String()))
			}
			if existingInfo != issuer {
				panic(fmt.Errorf("issuer %s already exists with conflicting DID info", addr.String()))
			}
		} else {
			if existingInfo, found := k.GetDidInfo(ctx, issuer.Did); found && existingInfo != issuer {
				panic(fmt.Errorf("issuer DID %s already exists with conflicting DID info", issuer.Did))
			}
			k.SetDID(ctx, addr, issuer.Did)
			k.SetDidInfo(ctx, issuer.Did, issuer)
		}

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
