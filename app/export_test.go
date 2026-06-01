package app_test

import (
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplyZeroHeightValidatorStateWithJailAllowListDoesNotPanic(t *testing.T) {
	testApp := apptesting.Setup(t, false)
	ctx := testApp.NewContext(false, cometbftproto.Header{Height: testApp.LastBlockHeight()})

	validators := testApp.StakingKeeper.GetValidators(ctx, 10)
	require.GreaterOrEqual(t, len(validators), 2)

	allowedAddrsMap := map[string]bool{
		validators[0].GetOperator().String(): true,
	}

	require.NotPanics(t, func() {
		testApp.ApplyZeroHeightValidatorState(ctx, true, allowedAddrsMap)
	})

	jailedValidator, found := testApp.StakingKeeper.GetValidator(ctx, validators[1].GetOperator())
	require.True(t, found)
	require.True(t, jailedValidator.Jailed)

	allowedValidator, found := testApp.StakingKeeper.GetValidator(ctx, validators[0].GetOperator())
	require.True(t, found)
	require.False(t, allowedValidator.Jailed)
}
