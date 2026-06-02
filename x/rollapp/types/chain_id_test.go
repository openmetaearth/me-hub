package types

import "testing"

func TestNewChainIDRejectsEIP155AboveUint64(t *testing.T) {
	_, err := NewChainID("evil_18446744073709551616-1")
	if err == nil {
		t.Fatal("expected EIP155 chain id above uint64 max to be rejected")
	}
}

func TestNewChainIDAcceptsMaxUint64EIP155(t *testing.T) {
	chainID, err := NewChainID("max_18446744073709551615-1")
	if err != nil {
		t.Fatalf("expected max uint64 EIP155 chain id to be accepted: %v", err)
	}

	if !chainID.IsEIP155() {
		t.Fatal("expected chain id to be detected as EIP155")
	}

	if got := chainID.GetEIP155ID(); got != 18446744073709551615 {
		t.Fatalf("expected max uint64 EIP155 id, got %d", got)
	}
}
