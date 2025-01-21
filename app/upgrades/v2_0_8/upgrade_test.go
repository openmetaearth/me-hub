package v2_0_8_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/st-chain/me-hub/app/upgrades/v2_0_8"
	"math/big"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	incentivestypes "github.com/st-chain/me-hub/x/incentives/types"
	"github.com/stretchr/testify/suite"

	"github.com/st-chain/me-hub/app"
	"github.com/st-chain/me-hub/app/apptesting"
)

// UpgradeTestSuite defines the structure for the upgrade test suite
type UpgradeTestSuite struct {
	suite.Suite
	Ctx sdk.Context
	App *app.App
}

// SetupTest initializes the necessary items for each test
func (s *UpgradeTestSuite) SetupTest(t *testing.T) {
	s.App = apptesting.Setup(t, false)
	s.Ctx = s.App.BaseApp.NewContext(false, cometbftproto.Header{Height: 1, ChainID: "dymension_100-1", Time: time.Now().UTC()})
}

// TestUpgradeTestSuite runs the suite of tests for the upgrade handler
func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

var (
	DYM = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

	// CreateGaugeFee is the fee required to create a new gauge.
	expectCreateGaugeFee = DYM.Mul(sdk.NewInt(10))
	// AddToGagugeFee is the fee required to add to gauge.
	expectAddToGaugeFee = sdk.ZeroInt()

	expectDelayedackEpochIdentifier = "hour"
	expectDelayedackBridgingFee     = sdk.NewDecWithPrec(1, 3)
)

const (
	dummyUpgradeHeight          = 5
	expectRollappsEnabled       = false
	expectDisputePeriodInBlocks = 120960
	expectMinBond               = "1000000000000000000000"
)

// TestUpgrade is a method of UpgradeTestSuite to test the upgrade process.
func (s *UpgradeTestSuite) TestUpgrade() {
	testCases := []struct {
		msg         string
		upgrade     func()
		postUpgrade func() error
		expPass     bool
	}{
		{
			"Test that upgrade does not panic and sets correct parameters",

			func() {
				// Run upgrade
				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
				plan := upgradetypes.Plan{Name: "v3", Height: dummyUpgradeHeight}
				err := s.App.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
				s.Require().NoError(err)
				_, exists := s.App.UpgradeKeeper.GetUpgradePlan(s.Ctx)
				s.Require().True(exists)

				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight)
				// simulate the upgrade process not panic.
				s.Require().NotPanics(func() {
					// simulate the upgrade process.
					s.App.BeginBlocker(s.Ctx, abci.RequestBeginBlock{})
				})
			},
			func() error {
				// Post-update validation to ensure parameters are correctly set

				// Check Delayedack parameters
				delayedackParams := s.App.DelayedAckKeeper.GetParams(s.Ctx)
				if delayedackParams.EpochIdentifier != expectDelayedackEpochIdentifier ||
					!delayedackParams.BridgingFee.Equal(expectDelayedackBridgingFee) {
					return fmt.Errorf("delayedack parameters not set correctly")
				}

				// Check Rollapp parameters
				rollappParams := s.App.RollappKeeper.GetParams(s.Ctx)
				if rollappParams.RollappsEnabled != expectRollappsEnabled || rollappParams.DisputePeriodInBlocks != expectDisputePeriodInBlocks {
					return fmt.Errorf("rollapp parameters not set correctly")
				}

				// Check Sequencer parameters
				seqParams := s.App.SequencerKeeper.GetParams(s.Ctx)
				if seqParams.MinBond.Amount.String() != expectMinBond {
					return fmt.Errorf("sequencer parameters not set correctly")
				}

				// Check Incentives parameters
				if !incentivestypes.CreateGaugeFee.Equal(expectCreateGaugeFee) || !incentivestypes.AddToGaugeFee.Equal(expectAddToGaugeFee) {
					return fmt.Errorf("incentives parameters not set correctly")
				}

				return nil
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			s.SetupTest(s.T()) // Reset for each case

			tc.upgrade()
			err := tc.postUpgrade()
			if tc.expPass {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *UpgradeTestSuite) TestDidData() {
	list := make(map[string]v2_0_8.DidData)
	list["me1ujufste3u23tpq3qhlq77u94nhw99emy3pr4p2"] = v2_0_8.DidData{
		Did:     "2405208027001",
		Uri:     "https://example.com/nft/metadata.json",
		UriHash: calculateUriHash("https://example.com/nft/metadata.json"),
	}
	list["me1phcakjkaf9vrn6jgttl3747dgnnpt88rt9440d"] = v2_0_8.DidData{
		Did:     "CHN2405204091002",
		Uri:     "https://example.com/nft1/metadata.json",
		UriHash: calculateUriHash("https://example.com/nft1/metadata.json"),
	}
	marshal, err := json.MarshalIndent(list, "", "  ")
	s.Require().NoError(err)
	s.T().Log(string(marshal))
}

func (s *UpgradeTestSuite) TestReadDidData() {
	list := `{
          "me1phcakjkaf9vrn6jgttl3747dgnnpt88rt9440d": {
            "did": "CHN2405204091002",
            "uri": "https://example.com/nft1/metadata.json",
            "uri_hash": "e00d196344dbd54550dadeab1167302ef39fade96eb211e302693a512ef131e1"
          },
          "me1ujufste3u23tpq3qhlq77u94nhw99emy3pr4p2": {
            "did": "2405208027001",
            "uri": "https://example.com/nft/metadata.json",
            "uri_hash": "9d4b6d9f5207db3461c26591fc4dd2a8e46ad4ce0d157f80b00bd57f52c90dec"
          }
        }`
	data := make(map[string]v2_0_8.DidData)
	err := json.Unmarshal([]byte(list), &data)
	s.Require().NoError(err)
	first, ok := data["me1phcakjkaf9vrn6jgttl3747dgnnpt88rt9440d"]
	s.Require().True(ok)
	s.Require().EqualValues("CHN2405204091002", first.Did)
}

func calculateUriHash(uri string) string {
	h := sha256.New()
	h.Write([]byte(uri))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
