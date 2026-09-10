package app_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/app/upgrades"
	"github.com/openmetaearth/me-hub/app/upgrades/v3"
	legacydelayedack "github.com/openmetaearth/me-hub/app/upgrades/v3/types/delayedack"
	legacyeibc "github.com/openmetaearth/me-hub/app/upgrades/v3/types/eibc"
	legacyrollapp "github.com/openmetaearth/me-hub/app/upgrades/v3/types/rollapp"
	legacysequencer "github.com/openmetaearth/me-hub/app/upgrades/v3/types/sequencer"
	"github.com/openmetaearth/me-hub/testutil/sample"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

func TestV3UpgradeHandler(t *testing.T) {
	meApp := apptesting.Setup(t)
	ctx := meApp.NewContext(false)

	const rollappID = "mecheckin_100-1"
	owner := sample.AccAddress()
	legacyMinBond := sdk.NewCoin(params.BaseDenom, math.NewInt(50_000_000))
	legacyNotice := 48 * time.Hour
	legacyDispute := uint64(42)
	legacyTimeoutFee := math.LegacyMustNewDecFromStr("0.02")
	legacyBridgingFee := math.LegacyMustNewDecFromStr("0.005")

	seedLegacyParams(t, ctx, meApp, legacysequencer.Params{
		MinBond:       legacyMinBond,
		UnbondingTime: legacyNotice,
	}, legacyrollapp.Params{
		DisputePeriodInBlocks: legacyDispute,
		RollappsEnabled:       true,
	}, legacyeibc.Params{
		EpochIdentifier: "day",
		TimeoutFee:      legacyTimeoutFee,
		ErrackFee:       math.LegacyMustNewDecFromStr("0.03"),
	}, legacydelayedack.Params{
		EpochIdentifier:         "day",
		BridgingFee:             legacyBridgingFee,
		DeletePacketsEpochLimit: 123,
	})

	oldRollapp := rollapptypes.Rollapp{
		RollappId: rollappID,
		Owner:     owner,
		ChannelId: "channel-0",
	}
	meApp.RollappKeeper.SetRollapp(ctx, oldRollapp)
	appendLegacyRegisteredDenoms(t, ctx, meApp, oldRollapp, "ibc/ABC")

	pk := ed25519.GenPrivKey().PubKey()
	pkAny, err := codectypes.NewAnyWithValue(pk)
	require.NoError(t, err)
	seqAddr := sdk.AccAddress(pk.Address()).String()
	meApp.SequencerKeeper.SetSequencer(ctx, sequencertypes.Sequencer{
		Address:      seqAddr,
		DymintPubKey: pkAny,
		RollappId:    rollappID,
		Status:       sequencertypes.Bonded,
		Tokens:       sdk.NewCoins(rollapptypes.DefaultMinSequencerBondGlobalCoin),
		OptedIn:      false,
	})

	fromVM, err := meApp.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, fromVM)

	handler := v3.CreateUpgradeHandler(
		meApp.ModuleManager(),
		meApp.Configurator(),
		upgradeKeepers(meApp),
	)
	_, err = handler(ctx, upgradetypes.Plan{Name: v3.UpgradeName, Height: ctx.BlockHeight()}, fromVM)
	require.NoError(t, err)

	seqParams := meApp.SequencerKeeper.GetParams(ctx)
	require.Equal(t, legacyNotice, seqParams.NoticePeriod)

	raParams := meApp.RollappKeeper.GetParams(ctx)
	require.Equal(t, legacyDispute, raParams.DisputePeriodInBlocks)
	require.True(t, raParams.MinSequencerBondGlobal.Equal(legacyMinBond))

	eibcParams := meApp.EIBCKeeper.GetParams(ctx)
	require.Equal(t, "day", eibcParams.EpochIdentifier)
	require.True(t, eibcParams.TimeoutFee.Equal(legacyTimeoutFee))

	daParams := meApp.DelayedAckKeeper.GetParams(ctx)
	require.Equal(t, "day", daParams.EpochIdentifier)
	require.True(t, daParams.BridgingFee.Equal(legacyBridgingFee))
	require.Equal(t, int32(123), daParams.DeletePacketsEpochLimit)

	govParams, err := meApp.GovKeeper.Params.Get(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, govParams.MinDeposit)
	require.Equal(t, govParams.MinDeposit[0].Denom, govParams.ExpeditedMinDeposit[0].Denom)
	require.True(t, govParams.ExpeditedMinDeposit[0].Amount.Equal(govParams.MinDeposit[0].Amount.MulRaw(5)))
	if govParams.VotingPeriod != nil && govParams.ExpeditedVotingPeriod != nil {
		require.Less(t, govParams.ExpeditedVotingPeriod.Seconds(), govParams.VotingPeriod.Seconds())
	}

	migrated, found := meApp.RollappKeeper.GetRollapp(ctx, rollappID)
	require.True(t, found)
	require.True(t, migrated.Launched)
	require.NotZero(t, migrated.GenesisState.TransferProofHeight)
	require.True(t, migrated.GenesisInfo.Sealed)
	hasDenom, err := meApp.RollappKeeper.HasRegisteredDenom(ctx, rollappID, "ibc/ABC")
	require.NoError(t, err)
	require.True(t, hasDenom)

	seq := meApp.SequencerKeeper.GetSequencer(ctx, seqAddr)
	require.False(t, seq.Sentinel())
	require.True(t, seq.OptedIn)
	_, err = meApp.SequencerKeeper.SequencerByDymintAddr(ctx, pk.Address())
	require.NoError(t, err)
	require.Equal(t, seqAddr, meApp.SequencerKeeper.GetProposer(ctx, rollappID).Address)
}

func upgradeKeepers(meApp *app.App) *upgrades.UpgradeKeepers {
	return &upgrades.UpgradeKeepers{
		AccountKeeper:     &meApp.AccountKeeper,
		GovKeeper:         meApp.GovKeeper,
		RollappKeeper:     meApp.RollappKeeper,
		SequencerKeeper:   meApp.SequencerKeeper,
		ParamsKeeper:      &meApp.ParamsKeeper,
		DelayedAckKeeper:  &meApp.DelayedAckKeeper,
		EIBCKeeper:        &meApp.EIBCKeeper,
		LightClientKeeper: &meApp.LightClientKeeper,
		IBCKeeper:         meApp.IBCKeeper,
		MintKeeper:        &meApp.MintKeeper,
		SlashingKeeper:    &meApp.SlashingKeeper,
		ConsensusKeeper:   &meApp.ConsensusParamsKeeper,
		StakingKeeper:     meApp.StakingKeeper,
	}
}

func seedLegacyParams(
	t *testing.T,
	ctx sdk.Context,
	meApp *app.App,
	seqParams legacysequencer.Params,
	raParams legacyrollapp.Params,
	eibcParams legacyeibc.Params,
	daParams legacydelayedack.Params,
) {
	t.Helper()

	seqSS := getOrCreateSubspace(meApp, legacysequencer.ModuleName)
	if !seqSS.HasKeyTable() {
		seqSS = seqSS.WithKeyTable(legacysequencer.ParamKeyTable())
	}
	seqSS.SetParamSet(ctx, &seqParams)

	raSS := getOrCreateSubspace(meApp, legacyrollapp.ModuleName)
	if !raSS.HasKeyTable() {
		raSS = raSS.WithKeyTable(legacyrollapp.ParamKeyTable())
	}
	raSS.SetParamSet(ctx, &raParams)

	eibcSS := getOrCreateSubspace(meApp, legacyeibc.ModuleName)
	if !eibcSS.HasKeyTable() {
		eibcSS = eibcSS.WithKeyTable(legacyeibc.ParamKeyTable())
	}
	eibcSS.SetParamSet(ctx, &eibcParams)

	daSS := getOrCreateSubspace(meApp, legacydelayedack.ModuleName)
	if !daSS.HasKeyTable() {
		daSS = daSS.WithKeyTable(legacydelayedack.ParamKeyTable())
	}
	daSS.SetParamSet(ctx, &daParams)
}

func getOrCreateSubspace(meApp *app.App, name string) paramstypes.Subspace {
	ss, ok := meApp.ParamsKeeper.GetSubspace(name)
	if !ok {
		ss = meApp.ParamsKeeper.Subspace(name)
	}
	return ss
}

func appendLegacyRegisteredDenoms(t *testing.T, ctx sdk.Context, meApp *app.App, ra rollapptypes.Rollapp, denom string) {
	t.Helper()
	store := prefix.NewStore(ctx.KVStore(meApp.GetKey(rollapptypes.StoreKey)), []byte(rollapptypes.RollappKeyPrefix))
	key := rollapptypes.RollappKey(ra.RollappId)
	bz := store.Get(key)
	require.NotEmpty(t, bz)
	store.Set(key, append(bz, encodeStringField(10, denom)...))
}

func encodeStringField(fieldNum int, s string) []byte {
	tag := byte(fieldNum<<3 | 2)
	b := []byte{tag, byte(len(s))}
	return append(b, s...)
}
