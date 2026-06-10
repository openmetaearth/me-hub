package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxCredentialDataSize = 64 * 1024

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
	dids := make(map[string]struct{}, len(gs.Infos))
	addresses := make(map[string]struct{}, len(gs.Infos))
	for i, info := range gs.Infos {
		if len(info.Did) != DidLength {
			return fmt.Errorf("infos[%d]: DID length must be equal to %d", i, DidLength)
		}
		if _, ok := dids[info.Did]; ok {
			return fmt.Errorf("infos[%d]: duplicated DID %q", i, info.Did)
		}
		dids[info.Did] = struct{}{}

		if _, err := sdk.AccAddressFromBech32(info.Address); err != nil {
			return fmt.Errorf("infos[%d]: invalid address %q: %w", i, info.Address, err)
		}
		if _, ok := addresses[info.Address]; ok {
			return fmt.Errorf("infos[%d]: duplicated address %q", i, info.Address)
		}
		addresses[info.Address] = struct{}{}

		if info.Pubkey == "" {
			return fmt.Errorf("infos[%d]: pubkey must not be empty", i)
		}
		if _, ok := DidStatus_name[int32(info.Status)]; !ok {
			return fmt.Errorf("infos[%d]: DID status must be ACTIVE or INACTIVE", i)
		}
	}

	services := make(map[string]struct{}, len(gs.Svcs))
	for i, svc := range gs.Svcs {
		if len(svc.Sid) < 2 || len(svc.Sid) > 8 {
			return fmt.Errorf("svcs[%d]: sid length must be between 2 and 8", i)
		}
		if _, ok := services[svc.Sid]; ok {
			return fmt.Errorf("svcs[%d]: duplicated service id %q", i, svc.Sid)
		}
		services[svc.Sid] = struct{}{}

		if len(svc.Name) == 0 || len(svc.Name) > 20 {
			return fmt.Errorf("svcs[%d]: name length must be between 1 and 20", i)
		}
		if len(svc.Description) > 1024 {
			return fmt.Errorf("svcs[%d]: description length exceeds 1024", i)
		}
		if _, ok := ServiceStatus_name[int32(svc.Status)]; !ok {
			return fmt.Errorf("svcs[%d]: service status must be ACTIVE or INACTIVE", i)
		}
		for issuerIndex, issuer := range svc.Issuers {
			if len(issuer) != DidLength {
				return fmt.Errorf("svcs[%d].issuers[%d]: issuer DID length must be equal to %d", i, issuerIndex, DidLength)
			}
			if _, ok := dids[issuer]; !ok {
				return fmt.Errorf("svcs[%d].issuers[%d]: issuer DID %q not found", i, issuerIndex, issuer)
			}
		}
	}

	credentials := make(map[string]struct{}, len(gs.Vcs))
	for i, vc := range gs.Vcs {
		key, err := validateGenesisCredential(i, vc, dids, services)
		if err != nil {
			return err
		}
		if _, ok := credentials[key]; ok {
			return fmt.Errorf("vcs[%d]: duplicated credential %q", i, key)
		}
		credentials[key] = struct{}{}
	}

	filterLogs := make(map[string]struct{}, len(gs.Flogs))
	for i, flog := range gs.Flogs {
		key := credentialKey(flog.Did, flog.Sid)
		if _, ok := credentials[key]; !ok {
			return fmt.Errorf("flogs[%d]: credential %q not found", i, key)
		}
		if _, ok := filterLogs[key]; ok {
			return fmt.Errorf("flogs[%d]: duplicated filter logger %q", i, key)
		}
		filterLogs[key] = struct{}{}
		for filterIndex, filter := range flog.Filters {
			if len(filter) > 1024 {
				return fmt.Errorf("flogs[%d].filters[%d]: filter length exceeds 1024", i, filterIndex)
			}
		}
	}

	return nil
}

func validateGenesisCredential(index int, vc Credential, dids, services map[string]struct{}) (string, error) {
	if len(vc.Did) != DidLength {
		return "", fmt.Errorf("vcs[%d]: DID length must be equal to %d", index, DidLength)
	}
	if _, ok := dids[vc.Did]; !ok {
		return "", fmt.Errorf("vcs[%d]: DID %q not found", index, vc.Did)
	}
	if len(vc.Sid) < 2 || len(vc.Sid) > 8 {
		return "", fmt.Errorf("vcs[%d]: sid length must be between 2 and 8", index)
	}
	if _, ok := services[vc.Sid]; !ok {
		return "", fmt.Errorf("vcs[%d]: service %q not found", index, vc.Sid)
	}
	if len(vc.Hash) == 0 || len(vc.Hash) > 128 {
		return "", fmt.Errorf("vcs[%d]: hash length must be between 1 and 128", index)
	}
	if len(vc.Uri) > 1024 {
		return "", fmt.Errorf("vcs[%d]: uri length exceeds 1024", index)
	}
	if len(vc.Data) > MaxCredentialDataSize {
		return "", fmt.Errorf("vcs[%d]: data length exceeds %d", index, MaxCredentialDataSize)
	}

	return credentialKey(vc.Did, vc.Sid), nil
}

func credentialKey(did, sid string) string {
	return did + "/" + sid
}
