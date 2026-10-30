package types

import (
	"strings"
	"testing"
	"time"

	"github.com/cosmos/gogoproto/proto"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestBatchEquityMessages(t *testing.T) {
	proposer := sample.AccAddress()
	first, second := sample.AccAddress(), sample.AccAddress()
	for _, addresses := range [][]string{{first}, {first, second}, nil, {first, first}, {first, strings.ToUpper(first)}, {first, "invalid"}} {
		valid := len(addresses) > 0 && (len(addresses) == 1 || addresses[1] == second)
		freeze := &MsgSubmitFreezeEquity{Proposer: proposer, RegionId: "usa", Addresses: addresses}
		unfreeze := &MsgSubmitUnfreezeEquity{Proposer: proposer, RegionId: "usa", Addresses: addresses}
		for _, err := range []error{freeze.ValidateBasic(), unfreeze.ValidateBasic()} {
			if valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}
		for _, kind := range []PendingActionKind{PendingActionKindFreezeEquity, PendingActionKindUnfreezeEquity} {
			action := PendingAction{Id: 1, Kind: kind, Status: PendingActionStatusPending, RegionId: "usa", Proposer: proposer, Addresses: addresses, SubmitTime: time.Now()}
			err := validatePendingAction(action)
			if valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}
		if valid {
			for _, msg := range []proto.Message{freeze, unfreeze} {
				data, err := proto.Marshal(msg)
				require.NoError(t, err)
				decoded := proto.Clone(msg)
				decoded.Reset()
				require.NoError(t, proto.Unmarshal(data, decoded))
				require.Equal(t, msg, decoded)
			}
		}
	}
}
