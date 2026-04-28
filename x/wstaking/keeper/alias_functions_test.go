package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/stretchr/testify/require"
)

// perBlockUMEC returns the expected per-block umec reward for each halving period.
//
// Period 0 (year 1): 79_274_480_000
// Period 1 (year 2): 39_637_240_000
// Period 2 (year 3): 19_818_620_000
var perBlockUMEC = [3]int64{
	79_274_480_000,
	39_637_240_000,
	19_818_620_000,
}

// newTestKeeper returns a zero-value Keeper that is sufficient to test
// getRewardsByHeight, which does not access any keeper fields.
func newTestKeeper() Keeper {
	return Keeper{}
}

func TestGetRewardsByHeight(t *testing.T) {
	Y := int64(mintTypes.OneYearTotalBlocks)

	tests := []struct {
		name       string
		fromHeight int64
		toHeight   int64
		wantUMEC   int64
	}{
		{
			name:       "same period (year 1): zero blocks when fromHeight==toHeight",
			fromHeight: 100,
			toHeight:   100,
			// blockCount = 0
			wantUMEC: 0,
		},
		{
			name:       "same period (year 1): single block",
			fromHeight: 1,
			toHeight:   2,
			// blockCount = 2-1 = 1
			wantUMEC: perBlockUMEC[0] * 1,
		},
		{
			name:       "same period (year 1): 4 blocks",
			fromHeight: 1,
			toHeight:   5,
			// blockCount = 5-1 = 4
			wantUMEC: perBlockUMEC[0] * 4,
		},
		{
			name:       "same period (year 1): 100 blocks",
			fromHeight: 1,
			toHeight:   101,
			// blockCount = 100
			wantUMEC: perBlockUMEC[0] * 100,
		},
		{
			name:       "same period (year 2): 3 blocks",
			fromHeight: Y + 1,
			toHeight:   Y + 4,
			// lowMul=1, highMul=1, blockCount=3
			wantUMEC: perBlockUMEC[1] * 3,
		},
		{
			name:       "same period (year 3): 100 blocks",
			fromHeight: 2*Y + 1,
			toHeight:   2*Y + 101,
			// lowMul=2, highMul=2, blockCount=100
			wantUMEC: perBlockUMEC[2] * 100,
		},
		{
			// Cross-period: period 0 -> period 1
			// fromHeight=Y-5, toHeight=Y+5
			// i=0 (lowMul): blockCount = Y*1 - (Y-5) + 1 = 6  (includes fromHeight)
			// i=1 (highMul): blockCount = (Y+5) - Y*1 - 1 = 4  (misses last block of period 1)
			name:       "cross-period year1 to year2",
			fromHeight: Y - 5,
			toHeight:   Y + 5,
			wantUMEC:   perBlockUMEC[0]*6 + perBlockUMEC[1]*4,
		},
		{
			// Cross-period: period 1 -> period 2
			// fromHeight=2Y-3, toHeight=2Y+3
			// i=1 (lowMul): blockCount = Y*2 - (2Y-3) + 1 = 4
			// i=2 (highMul): blockCount = (2Y+3) - Y*2 - 1 = 2
			name:       "cross-period year2 to year3",
			fromHeight: 2*Y - 3,
			toHeight:   2*Y + 3,
			wantUMEC:   perBlockUMEC[1]*4 + perBlockUMEC[2]*2,
		},
		{
			// Cross three periods: period 0 -> period 1 -> period 2
			// fromHeight=Y-1, toHeight=2Y+1
			// lowMul=0, highMul=2
			// i=0: blockCount = Y*1 - (Y-1) + 1 = 2
			// i=1: (full middle period) blockCount = Y
			// i=2: blockCount = (2Y+1) - Y*2 - 1 = 0
			name:       "cross three periods year1 to year3",
			fromHeight: Y - 1,
			toHeight:   2*Y + 1,
			wantUMEC:   perBlockUMEC[0]*2 + perBlockUMEC[1]*Y + perBlockUMEC[2]*0,
		},
	}

	k := newTestKeeper()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := k.getRewardsByHeight(tc.fromHeight, tc.toHeight)
			want := sdk.NewDecFromInt(sdk.NewInt(tc.wantUMEC))
			require.True(t, want.Equal(got),
				"fromHeight=%d toHeight=%d: want %s, got %s",
				tc.fromHeight, tc.toHeight, want, got,
			)
		})
	}
}
