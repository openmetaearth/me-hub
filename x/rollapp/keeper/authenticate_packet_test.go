package keeper

import (
	"fmt"
	"testing"

	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

type authenticatePacketChannelKeeper struct {
	chainIDsByEndpoint map[string]string
}

func (k authenticatePacketChannelKeeper) GetChannelClientState(
	_ sdk.Context,
	portID string,
	channelID string,
) (string, exported.ClientState, error) {
	chainID, ok := k.chainIDsByEndpoint[fmt.Sprintf("%s/%s", portID, channelID)]
	if !ok {
		return "", nil, fmt.Errorf("missing endpoint mapping for %s/%s", portID, channelID)
	}

	return "client-0", &ibctmtypes.ClientState{ChainId: chainID}, nil
}

func setupAuthenticatePacketKeeper(t *testing.T) (Keeper, sdk.Context) {
	t.Helper()

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	db := cometbftdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cometbftproto.Header{}, false, log.NewNopLogger())
	k := Keeper{
		cdc:      types.ModuleCdc,
		storeKey: storeKey,
		channelKeeper: authenticatePacketChannelKeeper{
			chainIDsByEndpoint: map[string]string{
				"transfer/channel-victim":   "victim_100-1",
				"transfer/channel-attacker": "attacker_100-1",
			},
		},
	}

	return k, ctx
}

func validAuthenticatePacketData() []byte {
	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "1",
		Sender:   "sender",
		Receiver: "receiver",
	}
	return transfertypes.ModuleCdc.MustMarshalJSON(&packetData)
}

func TestGetValidTransferFromReceivedPacketUsesDestinationEndpoint(t *testing.T) {
	k, ctx := setupAuthenticatePacketKeeper(t)
	k.SetRollapp(ctx, types.Rollapp{
		RollappId: "victim_100-1",
		ChannelId: "channel-victim",
	})

	transfer, err := k.GetValidTransferFromReceivedPacket(ctx, channeltypes.Packet{
		Data:               validAuthenticatePacketData(),
		SourcePort:         "transfer",
		SourceChannel:      "channel-attacker",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-victim",
	})

	require.NoError(t, err)
	require.True(t, transfer.IsRollapp())
	require.Equal(t, "victim_100-1", transfer.RollappId())
}

func TestGetValidTransferFromReceivedPacketIgnoresSpoofedSourceEndpoint(t *testing.T) {
	k, ctx := setupAuthenticatePacketKeeper(t)
	k.SetRollapp(ctx, types.Rollapp{
		RollappId: "victim_100-1",
		ChannelId: "channel-victim",
	})

	transfer, err := k.GetValidTransferFromReceivedPacket(ctx, channeltypes.Packet{
		Data:               validAuthenticatePacketData(),
		SourcePort:         "transfer",
		SourceChannel:      "channel-victim",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-attacker",
	})

	require.NoError(t, err)
	require.False(t, transfer.IsRollapp())
}

func TestGetValidTransferFromSentPacketUsesSourceEndpoint(t *testing.T) {
	k, ctx := setupAuthenticatePacketKeeper(t)
	k.SetRollapp(ctx, types.Rollapp{
		RollappId: "victim_100-1",
		ChannelId: "channel-victim",
	})

	transfer, err := k.GetValidTransferFromSentPacket(ctx, channeltypes.Packet{
		Data:               validAuthenticatePacketData(),
		SourcePort:         "transfer",
		SourceChannel:      "channel-victim",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-attacker",
	})

	require.NoError(t, err)
	require.True(t, transfer.IsRollapp())
	require.Equal(t, "victim_100-1", transfer.RollappId())
}
