package wmint

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wmint/types"
)

func TestMintNum(t *testing.T) {
	//1,667,931,888
	coin := getMintCoinsByHeight(0, 2104000)
	fmt.Println(coin.String())
}

type KeeperForTest struct {
	CoinAmount big.Int
}

func NewKeeperForTest() KeeperForTest {
	return KeeperForTest{
		CoinAmount: *big.NewInt(0),
	}
}

func (k KeeperForTest) GetMintedCoinAmount(ctx sdk.Context) big.Int {
	return k.CoinAmount
}
func (k *KeeperForTest) SetMintedCoinAmount(ctx sdk.Context, amount big.Int) {
	k.CoinAmount = amount
}

// BeginBlocker mints new tokens for the previous block.
func (k *KeeperForTest) BeginBlocker(ctx sdk.Context, height int64) big.Int {

	mintedAmount := k.GetMintedCoinAmount(ctx)
	blockHeight := height
	mul := (blockHeight - 1) / types.OneYearTotalBlocks
	amount := types.InitOneYearMintAmount / types.OneYearTotalBlocks / math.Exp2(float64(mul))
	mintingMEAmount := RoundUpToFourDecimals(amount)
	mintingUMEAmount := mintingMEAmount * math.Pow(10, params.BaseDenomUnit)

	// Compare the currently mined coins with the total amount of coins
	// -1 means that the current accumulated amount of mined is smaller than the total amount
	result := mintedAmount.Cmp(big.NewInt(types.TotalMintCoinsAmount))
	if result == -1 {
		// Accumulate the mined coins
		mintedAmount.Add(&mintedAmount, big.NewInt(int64(mintingUMEAmount)))
		k.SetMintedCoinAmount(ctx, mintedAmount)
	} else {
		mintingUMEAmount = 0
	}
	return k.GetMintedCoinAmount(ctx)
}

func TestCoinMinterAlgorithm(t *testing.T) {
	keeper := NewKeeperForTest()
	coin := keeper.BeginBlocker(sdk.Context{}, 8)
	fmt.Println(coin.String())
	fmt.Println("what:", getMintCoinsByHeight(0, 1))
}

func TestGetMintConsByHeight(t *testing.T) {
	coin := getMintCoinsByHeight(0, 2115004)
	fmt.Println(coin.String())
}
