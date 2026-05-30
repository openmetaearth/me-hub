package keeper_test

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// RollappTestSuite defines the test suite for rollapp keeper tests.
// It validates that the system correctly normalizes whitespace-padded RollappId
// values to prevent namespace duplicate attacks and EIP155 index squatting.
type RollappTestSuite struct {
	suite.Suite
	Ctx       sdk.Context
	msgServer types.MsgServer
	App       *TestApp
}

// alice is a deterministic testing address.
var alice = sdk.AccAddress([]byte("alice"))

// bob is an alternative testing address.
var bob = sdk.AccAddress([]byte("bob"))

// SetupTest initializes the test environment before each test.
// It creates a fresh test application and message server instance.
func (s *RollappTestSuite) SetupTest() {
	s.App = newTestApp()
	s.Ctx = s.App.BaseApp.NewContext(false)
	s.msgServer = keeper.NewMsgServerImpl(s.App.RollappKeeper)
}

// TestWhiteSpaceRollappIDNormalization ensures that all whitespace-padded
// RollappId values are normalized to their trimmed, canonical form during creation.
// This prevents raw storage of padded IDs and guarantees that only the canonical
// key appears in the primary and EIP155 secondary indexes.
func (s *RollappTestSuite) TestWhiteSpaceRollappIDNormalization() {
	// Table‑driven sub‑tests covering a wide range of padding scenarios.
	tests := []struct {
		name        string
		rawID       string
		canonicalID string
	}{
		{name: "leading and trailing spaces", rawID: "  test_1-1  ", canonicalID: "test_1-1"},
		{name: "only leading tabs", rawID: "\t\tevil_1-1", canonicalID: "evil_1-1"},
		{name: "only trailing spaces", rawID: "dup ", canonicalID: "dup"},
		{name: "mixed whitespace characters", rawID: " \t mixed_view_1-1 \n", canonicalID: "mixed_view_1-1"},
		{name: "no padding (control)", rawID: "clean_1-1", canonicalID: "clean_1-1"},
		{name: "newlines and carriage returns", rawID: "\n\nrollapp_1-1\r\n", canonicalID: "rollapp_1-1"},
		{name: "only internal whitespace retained", rawID: "a b_1-1", canonicalID: "a b_1-1"}, // internal stays
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Use a fresh context for each sub‑test to avoid state leakage.
			s.SetupTest()
			goCtx := sdk.WrapSDKContext(s.Ctx)
			s.T().Logf("Processing raw ID: %q, expected canonical: %q", tt.rawID, tt.canonicalID)

			// Act: create rollapp with the raw (possibly padded) ID.
			createResp, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
				Creator:               alice,
				RollappId:             tt.rawID,
				MaxSequencers:         1,
				PermissionedAddresses: []string{},
			})
			require.NoError(s.T(), err, "creation with raw ID '%s' must succeed", tt.rawID)
			require.NotNil(s.T(), createResp, "response must not be nil")

			// Assert 1: The raw (padded) ID must NOT be stored as a key.
			_, found := s.App.RollappKeeper.GetRollapp(s.Ctx, tt.rawID)
			require.False(s.T(), found, "raw key '%s' must not be stored", tt.rawID)

			// Assert 2: The canonical (trimmed) ID MUST be stored.
			stored, found := s.App.RollappKeeper.GetRollapp(s.Ctx, tt.canonicalID)
			require.True(s.T(), found, "canonical key '%s' must exist", tt.canonicalID)
			require.Equal(s.T(), tt.canonicalID, stored.RollappId,
				"stored RollappId must be the trimmed canonical version")

			// Assert 3: The chain ID derived from NewChainID must match the canonical.
			chainID, err := types.NewChainID(tt.rawID)
			require.NoError(s.T(), err, "NewChainID should accept padded ID")
			canonicalChainID, err := types.NewChainID(tt.canonicalID)
			require.NoError(s.T(), err)
			require.Equal(s.T(), canonicalChainID, chainID,
				"chain ID must be the same regardless of padding")

			// Assert 4: If EIP155, the secondary index must reference the canonical rollapp.
			if chainID.IsEIP155() {
				eipRollapp, found := s.App.RollappKeeper.GetRollappByEIP155(s.Ctx, chainID.EIP155ChainID())
				require.True(s.T(), found, "EIP155 index must exist for '%s'", tt.canonicalID)
				require.Equal(s.T(), tt.canonicalID, eipRollapp.RollappId,
					"EIP155 index must reference the canonical rollapp")
			}
		})
	}
}

// TestCanonicalDuplicateRejection ensures that after a padded ID is stored,
// any subsequent attempt to create a rollapp with the same canonical ID is rejected.
// This prevents namespace duplication attacks via whitespace variation.
func (s *RollappTestSuite) TestCanonicalDuplicateRejection() {
	s.SetupTest()
	goCtx := sdk.WrapSDKContext(s.Ctx)
	s.T().Log("Testing duplicate canonical rollapp rejection after padded creation")

	// Create first rollapp with padding.
	_, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             "  dup  ",
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.NoError(s.T(), err, "first creation must succeed")

	// Attempt to create the same canonical ID without padding.
	_, err = s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               bob,
		RollappId:             "dup",
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.Error(s.T(), err, "duplicate canonical rollapp must be rejected")
	require.Contains(s.T(), err.Error(), "already exists",
		"error message must indicate duplicate")
}

// TestEIP155IndexSquattingPrevention verifies that a padded EIP155 RollappId
// correctly occupies the EIP155 index, blocking any other rollapp with the same
// EIP155 chain ID (even if the rollapp name differs) from being created.
// This prevents index squatting via whitespace manipulation.
func (s *RollappTestSuite) TestEIP155IndexSquattingPrevention() {
	s.SetupTest()
	goCtx := sdk.WrapSDKContext(s.Ctx)
	s.T().Log("Testing EIP155 index squatting prevention with padded ID")

	rawID := "  good_1-1  "
	canonicalID := "good_1-1"
	differentRollappSameEIP155 := "bad_1-1" // Different name, same chain ID "1-1"

	// Create a rollapp with a padded EIP155 ID.
	_, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             rawID,
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.NoError(s.T(), err, "creation with padded EIP155 ID must succeed")

	// Verify the EIP155 index points to the canonical ID.
	eipRollapp, found := s.App.RollappKeeper.GetRollappByEIP155(s.Ctx, "1-1")
	require.True(s.T(), found, "EIP155 index must exist")
	require.Equal(s.T(), canonicalID, eipRollapp.RollappId,
		"EIP155 index must reference the canonical (trimmed) rollapp")

	// Attempt to create a different rollapp with the same EIP155 chain ID.
	_, err = s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               bob,
		RollappId:             differentRollappSameEIP155,
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.Error(s.T(), err, "EIP155 duplicate must be rejected")
	require.Contains(s.T(), err.Error(), "already exists",
		"error must indicate index conflict")
}

// TestWhitespaceOnlyID verifies that a RollappId consisting solely of whitespace
// is rejected by the validation logic, preventing nonsensical entries.
func (s *RollappTestSuite) TestWhitespaceOnlyID() {
	s.SetupTest()
	goCtx := sdk.WrapSDKContext(s.Ctx)
	s.T().Log("Testing rejection of whitespace-only RollappId")

	_, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             "   ",
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.Error(s.T(), err, "whitespace-only RollappId must be rejected")
}

// TestEmptyID verifies that an empty string is rejected.
func (s *RollappTestSuite) TestEmptyID() {
	s.SetupTest()
	goCtx := sdk.WrapSDKContext(s.Ctx)
	s.T().Log("Testing rejection of empty RollappId")

	_, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             "",
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.Error(s.T(), err, "empty RollappId must be rejected")
}

// TestMaxLengthID verifies that overly long IDs are properly rejected.
func (s *RollappTestSuite) TestMaxLengthID() {
	s.SetupTest()
	goCtx := sdk.WrapSDKContext(s.Ctx)
	s.T().Log("Testing rejection of oversize RollappId")

	// Build an ID exceeding any reasonable length (e.g., 1000 characters).
	longID := strings.Repeat("a", 1000)
	_, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             longID,
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.Error(s.T(), err, "oversize RollappId must be rejected")
}

// TestInvalidCharactersID ensures IDs containing forbidden characters are rejected.
func (s *RollappTestSuite) TestInvalidCharactersID() {
	s.SetupTest()
	goCtx := sdk.WrapSDKContext(s.Ctx)
	s.T().Log("Testing rejection of RollappId with invalid characters")

	_, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             "invalid@id!",
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.Error(s.T(), err, "invalid characters must be rejected")
}

// TestMultiplePaddedIDsSameCanonical ensures that once a canonical ID is stored,
// no other padded variation can be created, even with different creators.
func (s *RollappTestSuite) TestMultiplePaddedIDsSameCanonical() {
	s.SetupTest()
	goCtx := sdk.WrapSDKContext(s.Ctx)
	s.T().Log("Testing that multiple padded variations of same canonical ID are rejected")

	// Create first variation.
	_, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               alice,
		RollappId:             "  testroll  ",
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.NoError(s.T(), err, "first creation must succeed")

	// Attempt second variation.
	_, err = s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
		Creator:               bob,
		RollappId:             "\t\ntestroll\n",
		MaxSequencers:         1,
		PermissionedAddresses: []string{},
	})
	require.Error(s.T(), err, "second variation must be rejected")
	require.Contains(s.T(), err.Error(), "already exists")
}

// TestGetAllRollappsAfterNormalization verifies that the GetAllRollapps query returns
// only canonical keys and no raw keys.
func (s *RollappTestSuite) TestGetAllRollappsAfterNormalization() {
	s.SetupTest()
	goCtx := sdk.WrapSDKContext(s.Ctx)
	s.T().Log("Testing GetAllRollapps returns only canonical keys")

	// Create several rollapps with varied padding.
	for i, raw := range []string{"  alpha_1-1  ", "beta_2-2", "\tgamma_3-3\n"} {
		_, err := s.msgServer.CreateRollapp(goCtx, &types.MsgCreateRollapp{
			Creator:               alice,
			RollappId:             raw,
			MaxSequencers:         1,
			PermissionedAddresses: []string{},
		})
		require.NoError(s.T(), err, "creation %d must succeed", i)
	}

	// GetAllRollapps should return exactly 3 items, all with trimmed IDs.
	rollapps := s.App.RollappKeeper.GetAllRollapps(s.Ctx)
	require.Len(s.T(), rollapps, 3, "should have exactly 3 rollapps")

	hasAlpha := false
	hasBeta := false
	hasGamma := false
	for _, r := range rollapps {
		switch r.RollappId {
		case "alpha_1-1":
			hasAlpha = true
		case "beta_2-2":
			hasBeta = true
		case "gamma_3-3":
			hasGamma = true
		default:
			s.T().Errorf("unexpected RollappId: %q", r.RollappId)
		}
	}
	require.True(s.T(), hasAlpha, "alpha_1-1 must be present")
	require.True(s.T(), hasBeta, "beta_2-2 must be present")
	require.True(s.T(), hasGamma, "gamma_3-3 must be present")
}

// TestRollappKeeperSuite runs the entire test suite.
func TestRollappKeeperSuite(t *testing.T) {
	suite.Run(t, new(RollappTestSuite))
}