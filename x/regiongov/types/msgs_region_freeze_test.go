package types

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/testutil/sample"
)

func TestRegionFreezeMessages(t *testing.T) {
	authority := sample.AccAddress()
	for _, ids := range [][]string{{"usa"}, {"usa", "china"}} {
		freeze := &MsgFreezeRegion{Authority: authority, RegionIds: ids}
		unfreeze := &MsgUnfreezeRegion{Authority: authority, RegionIds: ids}
		require.NoError(t, freeze.ValidateBasic())
		require.NoError(t, unfreeze.ValidateBasic())
		for _, msg := range []proto.Message{freeze, unfreeze} {
			encoded, err := proto.Marshal(msg)
			require.NoError(t, err)
			decoded := proto.Clone(msg)
			decoded.Reset()
			require.NoError(t, proto.Unmarshal(encoded, decoded))
			require.Equal(t, msg, decoded)
		}
	}
	for _, ids := range [][]string{nil, {}, {""}, {"usa", ""}, {"usa", "usa"}, {"usa", "USA"}} {
		require.Error(t, (&MsgFreezeRegion{Authority: authority, RegionIds: ids}).ValidateBasic())
		require.Error(t, (&MsgUnfreezeRegion{Authority: authority, RegionIds: ids}).ValidateBasic())
	}
	require.Error(t, (&MsgFreezeRegion{Authority: "invalid", RegionIds: []string{"usa"}}).ValidateBasic())
	require.Error(t, (&MsgUnfreezeRegion{Authority: "invalid", RegionIds: []string{"usa"}}).ValidateBasic())
}
