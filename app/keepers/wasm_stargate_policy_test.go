package keepers

import (
	"strings"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

func TestRejectStargateMsg(t *testing.T) {
	_, err := rejectStargateMsg(nil, &wasmvmtypes.StargateMsg{TypeURL: "/metaearth.gravity.v1.MsgSendToExternal"})
	if err == nil {
		t.Fatal("expected stargate message to be rejected")
	}
	if !strings.Contains(err.Error(), "stargate messages are disabled") {
		t.Fatalf("unexpected error: %v", err)
	}
}
