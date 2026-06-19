package types

import "testing"

func TestDenomMetadataProposalTypes(t *testing.T) {
	createProposal := NewCreateMetadataProposal("create", "create metadata", nil)
	if got := createProposal.ProposalType(); got != ProposalTypeCreateDenomMetadata {
		t.Fatalf("create proposal type = %q, want %q", got, ProposalTypeCreateDenomMetadata)
	}

	updateProposal := NewUpdateDenomMetadataProposal("update", "update metadata", nil)
	if got := updateProposal.ProposalType(); got != ProposalTypeUpdateDenomMetadata {
		t.Fatalf("update proposal type = %q, want %q", got, ProposalTypeUpdateDenomMetadata)
	}
}
