package types_test

import (
	"testing"

	"github.com/openmetaearth/me-hub/x/tron/types"
	"github.com/stretchr/testify/require"
)

func TestValidateTronAddress(t *testing.T) {
	testCases := []struct {
		testName   string
		value      string
		expectPass bool
		errStr     string
	}{
		{
			testName:   "empty address",
			value:      "",
			expectPass: false,
			errStr:     "empty",
		},
		{
			testName:   "address length not match",
			value:      "abcdddddd",
			expectPass: false,
			errStr:     "wrong length",
		},
		{
			testName:   "address length great than tron address",
			value:      "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t6666",
			expectPass: false,
			errStr:     "wrong length",
		},
		{
			testName:   "lowercase address",
			value:      "tr7nhqjekqxgtci8q8zy4pl8otszgjlj6t",
			expectPass: false,
			errStr:     "doesn't pass format validation",
		},
		{
			testName:   "uppercase address",
			value:      "TR7NHQJEKQXGTCI8Q8ZY4PL8OTSZGJLJ6T",
			expectPass: false,
			errStr:     "doesn't pass format validation",
		},
		{
			testName:   "normal address",
			value:      "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			expectPass: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err := types.ValidateTronAddress(testCase.value)
			if testCase.expectPass {
				require.NoError(t, err)
				return
			}
			require.EqualValues(t, testCase.errStr, err.Error(), testCase.value)
		})
	}
}
