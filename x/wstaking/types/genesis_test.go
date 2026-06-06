package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
)

func TestGenesisStateMarshalPreservesReplaceConsensusPubKey(t *testing.T) {
	pendingReplace := []byte(`{"operator_address":"meta1operator","old_cons_address":"b2xkLWtleQ==","pub_key":"bmV3LWtleQ==","update_at_height":42}`)
	genesis := DefaultGenesisState()
	genesis.LastTotalPower = sdk.ZeroInt()
	genesis.ReplaceConsensusPubKey = pendingReplace

	bz, err := genesis.Marshal()
	require.NoError(t, err)

	var got GenesisState
	require.NoError(t, got.Unmarshal(bz))
	require.Equal(t, pendingReplace, got.GetReplaceConsensusPubKey())

	protoBz, err := proto.Marshal(genesis)
	require.NoError(t, err)

	var protoGot GenesisState
	require.NoError(t, proto.Unmarshal(protoBz, &protoGot))
	require.Equal(t, pendingReplace, protoGot.GetReplaceConsensusPubKey())
}
