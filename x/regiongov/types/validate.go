package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func validateOperatorAddress(address string) error {
	if _, err := sdk.ValAddressFromBech32(address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator address %s: %s", address, err)
	}
	return nil
}

func validateAccAddress(address string) error {
	if _, err := sdk.AccAddressFromBech32(address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address %s: %s", address, err)
	}
	return nil
}

func validateRegionID(regionID string) error {
	if regionID == "" {
		return fmt.Errorf("region id is empty")
	}
	return nil
}

func validateCoins(amount sdk.Coins) error {
	if amount.Empty() {
		return fmt.Errorf("amount is empty")
	}
	if err := amount.Validate(); err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}
	return nil
}

func validateUniqueAddresses(addresses []string) error {
	seen := make(map[string]struct{}, len(addresses))
	for _, address := range addresses {
		if err := validateAccAddress(address); err != nil {
			return err
		}
		if _, ok := seen[address]; ok {
			return fmt.Errorf("duplicate address %s", address)
		}
		seen[address] = struct{}{}
	}
	return nil
}

func validateEquityFreezeAddresses(addresses []string) error {
	if len(addresses) == 0 {
		return fmt.Errorf("equity freeze action requires at least one address")
	}
	seen := make(map[string]struct{}, len(addresses))
	for _, address := range addresses {
		acc, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return err
		}
		canonical := acc.String()
		if _, ok := seen[canonical]; ok {
			return fmt.Errorf("duplicate address %s", address)
		}
		seen[canonical] = struct{}{}
	}
	return nil
}

func validateWithdrawals(withdrawals []Withdrawal) error {
	if len(withdrawals) == 0 {
		return fmt.Errorf("treasury withdrawal is empty")
	}
	seen := make(map[string]struct{}, len(withdrawals))
	for _, withdrawal := range withdrawals {
		if err := validateRegionID(withdrawal.RegionId); err != nil {
			return err
		}
		if _, ok := seen[withdrawal.RegionId]; ok {
			return fmt.Errorf("duplicate withdrawal region %s", withdrawal.RegionId)
		}
		seen[withdrawal.RegionId] = struct{}{}
		if err := validateAccAddress(withdrawal.Receiver); err != nil {
			return err
		}
		if err := validateCoins(withdrawal.Amount); err != nil {
			return err
		}
	}
	return nil
}

func validateReplaceOperator(payload *ReplaceRegionOperatorPayload) error {
	if payload == nil {
		return fmt.Errorf("replace operator payload is empty")
	}
	if err := validateRegionID(payload.RegionId); err != nil {
		return err
	}
	return validateOperatorAddress(payload.NewOperator)
}
