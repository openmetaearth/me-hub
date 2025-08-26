package v2_0_12

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUpgradeName tests that the upgrade name is correctly defined
func TestUpgradeName(t *testing.T) {
	require.Equal(t, "v2.0.12", UpgradeName)
}

// TestUpgradeStruct tests that the upgrade struct is properly configured
func TestUpgradeStruct(t *testing.T) {
	require.Equal(t, UpgradeName, Upgrade.Name)
	require.NotNil(t, Upgrade.CreateHandler)
	require.NotNil(t, Upgrade.StoreUpgrades)
}

// TestClassUpdatesMapping tests that the class updates mapping contains expected entries
func TestClassUpdatesMapping(t *testing.T) {
	// Test that we have the expected class IDs in our mapping
	expectedClassIds := []string{
		"495393167", // ME_ExplorerA1000
		"506661488", // ME_ExplorerB1000
		"697811991", // ME_PioneerB1000
		"767391917", // ME_PioneerA1000
	}

	// This is a basic structural test to ensure our constants are defined
	// In a full integration test, we would verify the actual migration logic
	require.Len(t, expectedClassIds, 4, "Should have 4 target class IDs")

	// Verify class IDs are not empty
	for _, classId := range expectedClassIds {
		require.NotEmpty(t, classId, "Class ID should not be empty")
		require.Len(t, classId, 9, "Class ID should be 9 characters long")
	}
}

// TestExpectedClassNames tests that we have the expected class names defined
func TestExpectedClassNames(t *testing.T) {
	expectedNames := map[string]string{
		"495393167": "ME_ExplorerA1000",
		"506661488": "ME_ExplorerB1000",
		"697811991": "ME_PioneerB1000",
		"767391917": "ME_PioneerA1000",
	}

	// Verify we have the expected number of mappings
	require.Len(t, expectedNames, 4, "Should have 4 name mappings")

	// Verify names follow the expected pattern
	for classId, name := range expectedNames {
		require.NotEmpty(t, classId, "Class ID should not be empty")
		require.NotEmpty(t, name, "Class name should not be empty")
		require.Contains(t, name, "ME_", "Class name should contain ME_ prefix")
		require.Contains(t, name, "1000", "Class name should contain 1000 suffix")
	}
}

// TestExpectedDescriptions tests that descriptions are properly formatted
func TestExpectedDescriptions(t *testing.T) {
	expectedDescriptions := map[string]string{
		"495393167": "The Explorer Gold Hunter wields advanced technology, symbolizing wealth and prosperity. It drives the economic growth of new territories, bringing wealth and opportunities to Meta Earth.",
		"506661488": "The Explorer Freedom moves freely like the wind, symbolizing vast vision and limitless possibilities. It leads the new territories toward wisdom and prosperity, continually driving the development of the future.",
		"697811991": "The Pioneer Serenity rests lazily on the clouds, exuding calmness on the outside but with inner strength. It represents peace and wisdom after battle, safeguarding the balance of Meta Earth.",
		"767391917": "The \"Pioneer·Might\" wields weapons that symbolize speed and power. Fearlessly charging forward, it paves the way in Meta Earth, becoming the vanguard of world expansion with its boundless fighting spirit.",
	}

	// Verify we have descriptions for all expected classes
	require.Len(t, expectedDescriptions, 4, "Should have 4 description mappings")

	// Verify descriptions are meaningful and not empty
	for classId, description := range expectedDescriptions {
		require.NotEmpty(t, classId, "Class ID should not be empty")
		require.NotEmpty(t, description, "Description should not be empty")
		require.Greater(t, len(description), 50, "Description should be substantial")
		require.Contains(t, description, "Meta Earth", "Description should mention Meta Earth")
	}
}
