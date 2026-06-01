package types_test

import (
	"bytes"
	"testing"

	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	troncommon "github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/openmetaearth/me-hub/x/tron/types"
	"github.com/stretchr/testify/require"
)

func TestValidateTronAddress(t *testing.T) {
	nonTronPrefix := byte(0x00)
	nonTronPayload := append([]byte{nonTronPrefix}, bytes.Repeat([]byte{0x11}, tronaddress.AddressLength-1)...)
	shortTronPayload := append([]byte{tronaddress.TronBytePrefix}, bytes.Repeat([]byte{0x11}, tronaddress.AddressLength-2)...)

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
			errStr:     "invalid address length",
		},
		{
			testName:   "address length great than tron address",
			value:      "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t6666",
			expectPass: false,
			errStr:     "invalid address length",
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
			testName:   "base58check address with non-tron prefix",
			value:      troncommon.EncodeCheck(nonTronPayload),
			expectPass: false,
			errStr:     "invalid tron prefix",
		},
		{
			testName:   "base58check address with invalid decoded length",
			value:      troncommon.EncodeCheck(shortTronPayload),
			expectPass: false,
			errStr:     "invalid address length",
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
			require.Contains(t, err.Error(), testCase.errStr, testCase.value)
		})
	}
}
