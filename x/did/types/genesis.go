package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Infos: []DidInfo{},
		Svcs:  []Service{},
		Vcs:   []Credential{},
		Flogs: []FilterLogger{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	// 1. Validate Infos
	didMap := make(map[string]bool)
	addrMap := make(map[string]bool)
	for _, info := range gs.Infos {
		if len(info.Did) != DidLength {
			return fmt.Errorf("DID length must be equal to %d", DidLength)
		}
		if info.Pubkey == "" {
			return fmt.Errorf("pubkey must not be nil/empty")
		}
		if _, err := sdk.AccAddressFromBech32(info.Address); err != nil {
			return fmt.Errorf("invalid address %s: %s", info.Address, err)
		}
		if info.Status != DID_STATUS_INACTIVE && info.Status != DID_STATUS_ACTIVE {
			return fmt.Errorf("invalid DID status: %d", info.Status)
		}
		if didMap[info.Did] {
			return fmt.Errorf("duplicate DID info found: %s", info.Did)
		}
		didMap[info.Did] = true
		if addrMap[info.Address] {
			return fmt.Errorf("duplicate DID address found: %s", info.Address)
		}
		addrMap[info.Address] = true
	}

	// 2. Validate Svcs
	svcMap := make(map[string]bool)
	for _, svc := range gs.Svcs {
		if len(svc.Sid) < 2 || len(svc.Sid) > 8 {
			return fmt.Errorf("sid length must be between 2 and 8")
		}
		if len(svc.Name) == 0 || len(svc.Name) > 20 {
			return fmt.Errorf("name length must be between 1 and 20")
		}
		if len(svc.Description) > 1024 {
			return fmt.Errorf("description length exceeds 1024")
		}
		if svc.Status != SERVICE_STATUS_INACTIVE && svc.Status != SERVICE_STATUS_ACTIVE {
			return fmt.Errorf("invalid service status: %d", svc.Status)
		}
		if svcMap[svc.Sid] {
			return fmt.Errorf("duplicate service found: %s", svc.Sid)
		}
		svcMap[svc.Sid] = true

		for _, issuer := range svc.Issuers {
			if len(issuer) != DidLength {
				return fmt.Errorf("issuer length must be equal to %d", DidLength)
			}
			if !didMap[issuer] {
				return fmt.Errorf("issuer DID %s does not exist in DID Infos", issuer)
			}
		}
	}

	// 3. Validate Vcs
	type credKey struct {
		Did string
		Sid string
	}
	vcMap := make(map[credKey]bool)
	for _, vc := range gs.Vcs {
		if len(vc.Did) != DidLength {
			return fmt.Errorf("DID length in credential must be equal to %d", DidLength)
		}
		if len(vc.Sid) < 2 || len(vc.Sid) > 8 {
			return fmt.Errorf("sid length in credential must be between 2 and 8")
		}
		if len(vc.Hash) == 0 || len(vc.Hash) > 128 {
			return fmt.Errorf("hash length must be between 1 and 128")
		}
		if len(vc.Uri) > 1024 {
			return fmt.Errorf("uri length exceeds 1024")
		}
		if !didMap[vc.Did] {
			return fmt.Errorf("credential DID %s does not exist in DID Infos", vc.Did)
		}
		if !svcMap[vc.Sid] {
			return fmt.Errorf("credential Sid %s does not exist in Services", vc.Sid)
		}
		key := credKey{Did: vc.Did, Sid: vc.Sid}
		if vcMap[key] {
			return fmt.Errorf("duplicate credential found for did %s and sid %s", vc.Did, vc.Sid)
		}
		vcMap[key] = true
	}

	// 4. Validate Flogs
	flogMap := make(map[credKey]bool)
	for _, flog := range gs.Flogs {
		if len(flog.Did) != DidLength {
			return fmt.Errorf("DID length in filter logger must be equal to %d", DidLength)
		}
		if len(flog.Sid) < 2 || len(flog.Sid) > 8 {
			return fmt.Errorf("sid length in filter logger must be between 2 and 8")
		}
		for _, filter := range flog.Filters {
			if len(filter) > 1024 {
				return fmt.Errorf("filter length exceeds 1024")
			}
		}
		key := credKey{Did: flog.Did, Sid: flog.Sid}
		if !vcMap[key] {
			return fmt.Errorf("filter logger references non-existent credential [did: %s, sid: %s]", flog.Did, flog.Sid)
		}
		if flogMap[key] {
			return fmt.Errorf("duplicate filter logger found for did %s and sid %s", flog.Did, flog.Sid)
		}
		flogMap[key] = true
	}

	return nil
}
