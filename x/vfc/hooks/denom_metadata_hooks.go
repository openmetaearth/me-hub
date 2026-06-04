package hooks

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	denommetadatamoduletypes "github.com/openmetaearth/me-hub/x/denommetadata/types"
)

var _ denommetadatamoduletypes.DenomMetadataHooks = VirtualFrontierBankContractRegistrationHook{}

const ibcDenomPrefix = "ibc/"

type VirtualFrontierBankContractRegistrationHook struct {
	evmKeeper evmkeeper.Keeper
}

// NewVirtualFrontierBankContractRegistrationHook returns the DenomMetadataHooks for VFBC registration
func NewVirtualFrontierBankContractRegistrationHook(evmKeeper evmkeeper.Keeper) VirtualFrontierBankContractRegistrationHook {
	return VirtualFrontierBankContractRegistrationHook{
		evmKeeper: evmKeeper,
	}
}

func (v VirtualFrontierBankContractRegistrationHook) AfterDenomMetadataCreation(ctx sdk.Context, newDenomMetadata banktypes.Metadata) error {
	if strings.HasPrefix(strings.ToLower(newDenomMetadata.Base), ibcDenomPrefix) { // only deploy for IBC denom.
		// Deploy the virtual frontier bank contract for the new IBC denom.
		// Error, if any, no state transition will be made.
		if err := v.evmKeeper.DeployVirtualFrontierBankContractForBankDenomMetadataRecord(ctx, newDenomMetadata.Base); err != nil {
			return fmt.Errorf("deploy virtual frontier bank contract for IBC denom %s: %w", newDenomMetadata.Base, err)
		}
	}

	return nil
}

func (v VirtualFrontierBankContractRegistrationHook) AfterDenomMetadataUpdate(ctx sdk.Context, updatedDenomMetadata banktypes.Metadata) error {
	if !strings.HasPrefix(strings.ToLower(updatedDenomMetadata.Base), ibcDenomPrefix) {
		return nil
	}

	_, found := v.evmKeeper.GetVirtualFrontierBankContractAddressByDenom(ctx, updatedDenomMetadata.Base)
	if !found {
		return nil
	}

	// Existing VFBC instances keep their constructor metadata, so allowing a
	// bank metadata update here would make the two records disagree.
	return fmt.Errorf("cannot update IBC denom metadata for %s with existing virtual frontier bank contract", updatedDenomMetadata.Base)
}
