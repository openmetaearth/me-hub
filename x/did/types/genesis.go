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
	// --- validate Infos ---
	didSet := make(map[string]struct{}, len(gs.Infos))
	addrSet := make(map[string]struct{}, len(gs.Infos))
	for i, info := range gs.Infos {
		if len(info.Did) != DidLength {
			return fmt.Errorf("infos[%d]: DID length must be %d, got %d", i, DidLength, len(info.Did))
		}
		if info.Pubkey == "" {
			return fmt.Errorf("infos[%d]: pubkey must not be empty", i)
		}
		if _, err := sdk.AccAddressFromBech32(info.Address); err != nil {
			return fmt.Errorf("infos[%d]: invalid address %q: %w", i, info.Address, err)
		}
		if _, ok := DidStatus_name[int32(info.Status)]; !ok {
			return fmt.Errorf("infos[%d]: invalid DID status %d", i, info.Status)
		}
		if _, dup := didSet[info.Did]; dup {
			return fmt.Errorf("infos[%d]: duplicate DID %q", i, info.Did)
		}
		didSet[info.Did] = struct{}{}
		if _, dup := addrSet[info.Address]; dup {
			return fmt.Errorf("infos[%d]: duplicate address %q", i, info.Address)
		}
		addrSet[info.Address] = struct{}{}
	}

	// --- validate Svcs ---
	svcSet := make(map[string]struct{}, len(gs.Svcs))
	for i, svc := range gs.Svcs {
		if len(svc.Sid) < 2 || len(svc.Sid) > 8 {
			return fmt.Errorf("svcs[%d]: sid length must be between 2 and 8, got %d", i, len(svc.Sid))
		}
		if _, ok := ServiceStatus_name[int32(svc.Status)]; !ok {
			return fmt.Errorf("svcs[%d]: invalid service status %d", i, svc.Status)
		}
		if _, dup := svcSet[svc.Sid]; dup {
			return fmt.Errorf("svcs[%d]: duplicate service sid %q", i, svc.Sid)
		}
		svcSet[svc.Sid] = struct{}{}
	}

	// --- validate Vcs ---
	vcSet := make(map[string]struct{}, len(gs.Vcs))
	for i, vc := range gs.Vcs {
		if len(vc.Did) != DidLength {
			return fmt.Errorf("vcs[%d]: DID length must be %d, got %d", i, DidLength, len(vc.Did))
		}
		if len(vc.Sid) < 2 || len(vc.Sid) > 8 {
			return fmt.Errorf("vcs[%d]: sid length must be between 2 and 8, got %d", i, len(vc.Sid))
		}
		if len(vc.Hash) == 0 || len(vc.Hash) > 128 {
			return fmt.Errorf("vcs[%d]: hash length must be between 1 and 128, got %d", i, len(vc.Hash))
		}
		if len(vc.Data) > maxCredentialDataLength {
			return fmt.Errorf("vcs[%d]: data length exceeds %d", i, maxCredentialDataLength)
		}
		key := vc.Did + "/" + vc.Sid
		if _, dup := vcSet[key]; dup {
			return fmt.Errorf("vcs[%d]: duplicate credential (did=%s, sid=%s)", i, vc.Did, vc.Sid)
		}
		vcSet[key] = struct{}{}
	}

	// --- validate Flogs ---
	// Each FilterLogger must reference an existing Credential in Vcs to prevent
	// a panic in InitGenesis when GetCredential returns not-found.
	flogSet := make(map[string]struct{}, len(gs.Flogs))
	for i, flog := range gs.Flogs {
		if len(flog.Did) != DidLength {
			return fmt.Errorf("flogs[%d]: DID length must be %d, got %d", i, DidLength, len(flog.Did))
		}
		if len(flog.Sid) < 2 || len(flog.Sid) > 8 {
			return fmt.Errorf("flogs[%d]: sid length must be between 2 and 8, got %d", i, len(flog.Sid))
		}
		vcKey := flog.Did + "/" + flog.Sid
		if _, found := vcSet[vcKey]; !found {
			return fmt.Errorf("flogs[%d]: referenced credential (did=%s, sid=%s) not found in vcs", i, flog.Did, flog.Sid)
		}
		if _, dup := flogSet[vcKey]; dup {
			return fmt.Errorf("flogs[%d]: duplicate filter logger (did=%s, sid=%s)", i, flog.Did, flog.Sid)
		}
		flogSet[vcKey] = struct{}{}
	}

	return nil
}
