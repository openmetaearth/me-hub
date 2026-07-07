package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/wnft/keeper"
	"github.com/stretchr/testify/suite"
)

type ClassSupplyTestSuite struct {
	apptesting.KeeperTestHelper
	keeper   *keeper.Keeper
	TestAccs []sdk.AccAddress
}

func TestClassSupplyTestSuite(t *testing.T) {
	suite.Run(t, new(ClassSupplyTestSuite))
}

func (s *ClassSupplyTestSuite) SetupTest() {
	app := apptesting.Setup(s.T())
	ctx := app.GetBaseApp().NewContext(false)

	s.App = app
	s.Ctx = ctx
	s.keeper = app.WNFTKeeper

	s.TestAccs = s.NewAccounts(3)
}

// TestSetGetClassTotalSupplyCap tests setting and getting total supply cap
func (s *ClassSupplyTestSuite) TestSetGetClassTotalSupplyCap() {
	testCases := []struct {
		name     string
		classID  string
		supply   uint64
		expected uint64
	}{
		{
			name:     "set supply cap for class1",
			classID:  "class1",
			supply:   1000,
			expected: 1000,
		},
		{
			name:     "set supply cap for class2",
			classID:  "class2",
			supply:   500,
			expected: 500,
		},
		{
			name:     "set zero supply cap",
			classID:  "class-zero",
			supply:   0,
			expected: 0,
		},
		{
			name:     "set large supply cap",
			classID:  "class-large",
			supply:   1000000000,
			expected: 1000000000,
		},
		{
			name:     "set max uint64 supply cap",
			classID:  "class-max",
			supply:   ^uint64(0), // max uint64
			expected: ^uint64(0),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Set supply cap
			err := s.keeper.SetClassTotalSupplyCap(s.Ctx, tc.classID, tc.supply)
			s.Require().NoError(err)

			// Get and verify supply cap
			cap := s.keeper.GetClassTotalSupplyCap(s.Ctx, tc.classID)
			s.Require().Equal(tc.expected, cap)
		})
	}
}

// TestGetClassTotalSupplyCapNonExistent tests getting supply cap for non-existent class
func (s *ClassSupplyTestSuite) TestGetClassTotalSupplyCapNonExistent() {
	cap := s.keeper.GetClassTotalSupplyCap(s.Ctx, "non-existent-class")
	s.Require().Equal(uint64(0), cap, "non-existent class should return 0")
}

// TestUpdateClassTotalSupplyCap tests updating existing supply cap
func (s *ClassSupplyTestSuite) TestUpdateClassTotalSupplyCap() {
	classID := "update-test-class"

	// Set initial supply cap
	err := s.keeper.SetClassTotalSupplyCap(s.Ctx, classID, 100)
	s.Require().NoError(err)

	cap := s.keeper.GetClassTotalSupplyCap(s.Ctx, classID)
	s.Require().Equal(uint64(100), cap)

	// Update supply cap
	err = s.keeper.SetClassTotalSupplyCap(s.Ctx, classID, 200)
	s.Require().NoError(err)

	cap = s.keeper.GetClassTotalSupplyCap(s.Ctx, classID)
	s.Require().Equal(uint64(200), cap, "supply cap should be updated")

	// Update to zero
	err = s.keeper.SetClassTotalSupplyCap(s.Ctx, classID, 0)
	s.Require().NoError(err)

	cap = s.keeper.GetClassTotalSupplyCap(s.Ctx, classID)
	s.Require().Equal(uint64(0), cap, "supply cap should be updated to zero")
}

// TestClassTotalSupplyCapPersistence tests that supply cap persists across contexts
func (s *ClassSupplyTestSuite) TestClassTotalSupplyCapPersistence() {
	classID := "persistence-test-class"
	supply := uint64(999)

	// Set supply cap
	err := s.keeper.SetClassTotalSupplyCap(s.Ctx, classID, supply)
	s.Require().NoError(err)

	// Create new context (simulating new block)
	newCtx := s.App.GetBaseApp().NewContext(false)

	// Verify supply cap persists
	cap := s.keeper.GetClassTotalSupplyCap(newCtx, classID)
	s.Require().Equal(supply, cap, "supply cap should persist across contexts")
}

// TestMultipleClassSupplyCaps tests managing multiple classes simultaneously
func (s *ClassSupplyTestSuite) TestMultipleClassSupplyCaps() {
	classes := map[string]uint64{
		"class-a": 100,
		"class-b": 200,
		"class-c": 300,
		"class-d": 400,
		"class-e": 500,
	}

	// Set all supply caps
	for classID, supply := range classes {
		err := s.keeper.SetClassTotalSupplyCap(s.Ctx, classID, supply)
		s.Require().NoError(err)
	}

	// Verify all supply caps
	for classID, expectedSupply := range classes {
		cap := s.keeper.GetClassTotalSupplyCap(s.Ctx, classID)
		s.Require().Equal(expectedSupply, cap, "supply cap for %s should match", classID)
	}

	// Update one class and verify others remain unchanged
	err := s.keeper.SetClassTotalSupplyCap(s.Ctx, "class-c", 999)
	s.Require().NoError(err)

	s.Require().Equal(uint64(999), s.keeper.GetClassTotalSupplyCap(s.Ctx, "class-c"))
	s.Require().Equal(uint64(100), s.keeper.GetClassTotalSupplyCap(s.Ctx, "class-a"))
	s.Require().Equal(uint64(200), s.keeper.GetClassTotalSupplyCap(s.Ctx, "class-b"))
	s.Require().Equal(uint64(400), s.keeper.GetClassTotalSupplyCap(s.Ctx, "class-d"))
	s.Require().Equal(uint64(500), s.keeper.GetClassTotalSupplyCap(s.Ctx, "class-e"))
}

// TestClassSupplyCapWithSpecialCharacters tests class IDs with special characters
func (s *ClassSupplyTestSuite) TestClassSupplyCapWithSpecialCharacters() {
	testCases := []struct {
		name    string
		classID string
		supply  uint64
	}{
		{
			name:    "class with dash",
			classID: "test-class-123",
			supply:  100,
		},
		{
			name:    "class with underscore",
			classID: "test_class_456",
			supply:  200,
		},
		{
			name:    "class with numbers",
			classID: "class123456",
			supply:  300,
		},
		{
			name:    "lowercase class",
			classID: "lowercase",
			supply:  400,
		},
		{
			name:    "uppercase class",
			classID: "UPPERCASE",
			supply:  500,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.keeper.SetClassTotalSupplyCap(s.Ctx, tc.classID, tc.supply)
			s.Require().NoError(err)

			cap := s.keeper.GetClassTotalSupplyCap(s.Ctx, tc.classID)
			s.Require().Equal(tc.supply, cap)
		})
	}
}

// TestClassSupplyCapKeyCollision tests that different class IDs don't collide
func (s *ClassSupplyTestSuite) TestClassSupplyCapKeyCollision() {
	// Test potential collision scenarios
	testCases := []struct {
		classID string
		supply  uint64
	}{
		{"test", 100},
		{"test1", 200},
		{"test-1", 300},
		{"test_1", 400},
	}

	// Set all supplies
	for _, tc := range testCases {
		err := s.keeper.SetClassTotalSupplyCap(s.Ctx, tc.classID, tc.supply)
		s.Require().NoError(err)
	}

	// Verify all supplies are independent
	for _, tc := range testCases {
		cap := s.keeper.GetClassTotalSupplyCap(s.Ctx, tc.classID)
		s.Require().Equal(tc.supply, cap, "supply for %s should be independent", tc.classID)
	}
}

// TestClassSupplyCapBinaryEncoding tests that uint64 is properly encoded/decoded
func (s *ClassSupplyTestSuite) TestClassSupplyCapBinaryEncoding() {
	// Test boundary values for uint64
	testCases := []struct {
		name   string
		supply uint64
	}{
		{"zero", 0},
		{"one", 1},
		{"max uint8", 255},
		{"max uint16", 65535},
		{"max uint32", 4294967295},
		{"large value", 1234567890123456789},
		{"max uint64", ^uint64(0)},
	}

	classID := "encoding-test"
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.keeper.SetClassTotalSupplyCap(s.Ctx, classID, tc.supply)
			s.Require().NoError(err)

			cap := s.keeper.GetClassTotalSupplyCap(s.Ctx, classID)
			s.Require().Equal(tc.supply, cap, "encoding/decoding should preserve value")
		})
	}
}
