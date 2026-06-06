package keepers

import "testing"

func TestFilterWasmCapabilitiesRemovesDisabledCapabilities(t *testing.T) {
	capabilities := []string{"iterator", "staking", "stargate", "cosmwasm_1_4"}

	filtered := filterWasmCapabilities(capabilities, "stargate")

	expected := []string{"iterator", "staking", "cosmwasm_1_4"}
	if len(filtered) != len(expected) {
		t.Fatalf("expected %d capabilities, got %d: %v", len(expected), len(filtered), filtered)
	}
	for i := range expected {
		if filtered[i] != expected[i] {
			t.Fatalf("capability %d: expected %q, got %q", i, expected[i], filtered[i])
		}
	}
}
