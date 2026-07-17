package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sequencerkeeper "github.com/openmetaearth/me-hub/x/sequencer/keeper"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

func migrateSequencers(ctx sdk.Context, k *sequencerkeeper.Keeper) {
	list := k.AllSequencers(ctx)
	for _, oldSequencer := range list {
		newSequencer := ConvertOldSequencerToNew(oldSequencer)
		k.SetSequencer(ctx, newSequencer)
	}
}

func ConvertOldSequencerToNew(old sequencertypes.Sequencer) sequencertypes.Sequencer {
	return sequencertypes.Sequencer{
		Address:      old.Address,
		DymintPubKey: old.DymintPubKey,
		RollappId:    old.RollappId,
		Status:       old.Status,
		Tokens:       old.Tokens,
		OptedIn:      true,
		Metadata: sequencertypes.SequencerMetadata{
			Moniker:     old.Metadata.Moniker,
			Details:     old.Metadata.Details,
			P2PSeeds:    nil,
			Rpcs:        nil,
			EvmRpcs:     nil,
			RestApiUrls: []string{},
			ExplorerUrl: "",
			GenesisUrls: []string{},
			ContactDetails: &sequencertypes.ContactDetails{
				Website:  "",
				Telegram: "",
				X:        "",
			},
			ExtraData: nil,
			Snapshots: []*sequencertypes.SnapshotInfo{},
			GasPrice:  "10000000000",
			FeeDenom:  nil,
		},
	}
}