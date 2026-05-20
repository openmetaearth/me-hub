package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// TestNewRecord tests creating new records
func (s *KeeperTestSuite) TestNewRecord() {
	testCases := []struct {
		name      string
		setup     func() *types.MsgNewRecord
		expectErr bool
		errMsg    string
	}{
		{
			name: "successful record creation",
			setup: func() *types.MsgNewRecord {
				account := s.TestAccs[0]
				return &types.MsgNewRecord{
					From:         account.String(),
					ActionNumber: "REC123456",
					ActionUrl:    "https://example.com/action/123456",
				}
			},
			expectErr: false,
		},
		{
			name: "record with alphanumeric number",
			setup: func() *types.MsgNewRecord {
				account := s.TestAccs[1]
				return &types.MsgNewRecord{
					From:         account.String(),
					ActionNumber: "ABC123XYZ789",
					ActionUrl:    "https://example.com/action/abc123",
				}
			},
			expectErr: false,
		},
		{
			name: "empty action number",
			setup: func() *types.MsgNewRecord {
				account := s.TestAccs[0]
				return &types.MsgNewRecord{
					From:         account.String(),
					ActionNumber: "",
					ActionUrl:    "https://example.com/action/123",
				}
			},
			expectErr: true,
			errMsg:    "invalid record number,is empty",
		},
		{
			name: "empty action url",
			setup: func() *types.MsgNewRecord {
				account := s.TestAccs[0]
				return &types.MsgNewRecord{
					From:         account.String(),
					ActionNumber: "REC123",
					ActionUrl:    "",
				}
			},
			expectErr: true,
			errMsg:    "url is empty",
		},
		{
			name: "action number with special characters",
			setup: func() *types.MsgNewRecord {
				account := s.TestAccs[0]
				return &types.MsgNewRecord{
					From:         account.String(),
					ActionNumber: "REC-123-456", // Contains hyphens
					ActionUrl:    "https://example.com/action/123",
				}
			},
			expectErr: true,
			errMsg:    "only letters and numbers are allowed",
		},
		{
			name: "action number with spaces",
			setup: func() *types.MsgNewRecord {
				account := s.TestAccs[0]
				return &types.MsgNewRecord{
					From:         account.String(),
					ActionNumber: "REC 123 456", // Contains spaces
					ActionUrl:    "https://example.com/action/123",
				}
			},
			expectErr: true,
			errMsg:    "only letters and numbers are allowed",
		},
		{
			name: "invalid from address",
			setup: func() *types.MsgNewRecord {
				return &types.MsgNewRecord{
					From:         "invalid_address",
					ActionNumber: "REC123",
					ActionUrl:    "https://example.com/action/123",
				}
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Use fresh context for each test
			s.SetupTest()
			msg := tc.setup()

			resp, err := s.msgServer.NewRecord(s.Ctx, msg)

			if tc.expectErr {
				s.Require().Error(err)
				if tc.errMsg != "" {
					s.Require().Contains(err.Error(), tc.errMsg)
				}
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				// Verify record was stored
				fromAddr, _ := sdk.AccAddressFromBech32(msg.From)
				records := s.App.StakingKeeper.GetRecordsByAddress(s.Ctx, fromAddr)
				s.Require().NotEmpty(records, "Record should be stored")
				// Find the specific record
				found := false
				for _, record := range records {
					if record.RecordNumber == msg.ActionNumber {
						found = true
						s.Require().Equal(msg.ActionUrl, record.Url)
						s.Require().Equal(msg.From, record.From)
						break
					}
				}
				s.Require().True(found, "Specific record should be found")
			}
		})
	}
}

// TestReviewRecord tests reviewing records
func (s *KeeperTestSuite) TestReviewRecord() {
	testCases := []struct {
		name      string
		setup     func() *types.MsgReviewRecord
		expectErr bool
		errMsg    string
	}{
		{
			name: "successful review by global dao",
			setup: func() *types.MsgReviewRecord {
				return &types.MsgReviewRecord{
					From:            s.Dao.GlobalDao,
					RecordHash:      "hash123456",
					ActionNumber:    "REC123456",
					ReviewResult:    "approved",
					ReviewedAddress: s.TestAccs[0].String(),
				}
			},
			expectErr: false,
		},
		{
			name: "successful review by meid dao",
			setup: func() *types.MsgReviewRecord {
				return &types.MsgReviewRecord{
					From:            s.Dao.MeidDao,
					RecordHash:      "hash789012",
					ActionNumber:    "REC789012",
					ReviewResult:    "rejected",
					ReviewedAddress: s.TestAccs[1].String(),
				}
			},
			expectErr: false,
		},
		{
			name: "unauthorized reviewer",
			setup: func() *types.MsgReviewRecord {
				account := s.TestAccs[0]
				return &types.MsgReviewRecord{
					From:            account.String(),
					RecordHash:      "hash123",
					ActionNumber:    "REC123",
					ReviewResult:    "approved",
					ReviewedAddress: s.TestAccs[1].String(),
				}
			},
			expectErr: true,
			errMsg:    "should  be global admin",
		},
		{
			name: "empty review result",
			setup: func() *types.MsgReviewRecord {
				return &types.MsgReviewRecord{
					From:            s.Dao.GlobalDao,
					RecordHash:      "hash456",
					ActionNumber:    "REC456",
					ReviewResult:    "",
					ReviewedAddress: s.TestAccs[0].String(),
				}
			},
			expectErr: true,
			errMsg:    "review result is empty",
		},
		{
			name: "empty record hash",
			setup: func() *types.MsgReviewRecord {
				return &types.MsgReviewRecord{
					From:            s.Dao.GlobalDao,
					RecordHash:      "",
					ActionNumber:    "REC789",
					ReviewResult:    "approved",
					ReviewedAddress: s.TestAccs[0].String(),
				}
			},
			expectErr: true,
			errMsg:    "invalid record hash,is empty",
		},
		{
			name: "empty action number",
			setup: func() *types.MsgReviewRecord {
				return &types.MsgReviewRecord{
					From:            s.Dao.GlobalDao,
					RecordHash:      "hash999",
					ActionNumber:    "",
					ReviewResult:    "approved",
					ReviewedAddress: s.TestAccs[0].String(),
				}
			},
			expectErr: true,
			errMsg:    "invalid record number,is empty",
		},
		{
			name: "invalid reviewer address",
			setup: func() *types.MsgReviewRecord {
				return &types.MsgReviewRecord{
					From:            "invalid_address",
					RecordHash:      "hash111",
					ActionNumber:    "REC111",
					ReviewResult:    "approved",
					ReviewedAddress: s.TestAccs[0].String(),
				}
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Use fresh context for each test
			s.SetupTest()
			msg := tc.setup()

			resp, err := s.msgServer.ReviewRecord(s.Ctx, msg)

			if tc.expectErr {
				s.Require().Error(err)
				if tc.errMsg != "" {
					s.Require().Contains(err.Error(), tc.errMsg)
				}
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				// Verify review record was stored
				reviewRecord := s.App.StakingKeeper.GetReviewRecordByID(s.Ctx, msg.ActionNumber)
				s.Require().Equal(msg.RecordHash, reviewRecord.RecordHash)
				s.Require().Equal(msg.ActionNumber, reviewRecord.ActionNumber)
				s.Require().Equal(msg.ReviewResult, reviewRecord.RecordResult)
				s.Require().Equal(msg.ReviewedAddress, reviewRecord.ReviewedAddress)
			}
		})
	}
}

// TestReviewRecordMultipleTimes tests multiple reviews
func (s *KeeperTestSuite) TestReviewRecordMultipleTimes() {
	// Create multiple review records
	reviews := []struct {
		hash   string
		number string
		result string
	}{
		{"hash001", "REC001", "approved"},
		{"hash002", "REC002", "rejected"},
		{"hash003", "REC003", "pending"},
	}

	for _, review := range reviews {
		msg := &types.MsgReviewRecord{
			From:            s.Dao.GlobalDao,
			RecordHash:      review.hash,
			ActionNumber:    review.number,
			ReviewResult:    review.result,
			ReviewedAddress: s.TestAccs[0].String(),
		}

		resp, err := s.msgServer.ReviewRecord(s.Ctx, msg)
		s.Require().NoError(err)
		s.Require().NotNil(resp)

		// Verify each review was stored correctly
		reviewRecord := s.App.StakingKeeper.GetReviewRecordByID(s.Ctx, review.number)
		s.Require().Equal(review.result, reviewRecord.RecordResult)
	}
}
