package ante_test

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/ante"
	"github.com/openmetaearth/me-hub/app/ante/mock"
	megrouptypes "github.com/openmetaearth/me-hub/x/megroup/types"
)

// newJoinGroupCtx returns a minimal sdk.Context suitable for unit tests.
func newJoinGroupCtx() sdk.Context {
	return sdk.NewContext(nil, cometbftproto.Header{}, false, nil)
}

// stubNext is a trivial next-handler that simply returns the context it was given.
func stubNext(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }

// buildJoinGroupMsg returns a MsgJoinGroup with the given creator and applicant.
func buildJoinGroupMsg(creator, applicant string, groupID uint64) *megrouptypes.MsgJoinGroup {
	return &megrouptypes.MsgJoinGroup{
		Creator:          creator,
		ApplicantAddress: applicant,
		GroupId:          groupID,
	}
}

// newDecorator is a convenience constructor for the decorator under test.
func newDecorator(ctrl *gomock.Controller) (ante.JoinGroupValidateDecorator, *mock.MockDaoKeeper, *mock.MockMeGroupKeeper) {
	dk := mock.NewMockDaoKeeper(ctrl)
	gk := mock.NewMockMeGroupKeeper(ctrl)
	d := ante.NewJoinGroupValidateDecorator(dk, gk)
	return d, dk, gk
}

// ─── table-driven happy-path and error cases ────────────────────────────────

func TestJoinGroupValidateDecorator(t *testing.T) {
	regionID := "region-1"
	group := megrouptypes.GroupInfo{Id: 1, RegionID: regionID}

	creator := NewAccount()
	applicant := NewAccount()
	daoAdmin := NewAccount()

	tests := []struct {
		name      string
		msgs      []sdk.Msg
		setup     func(dk *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper)
		expectErr bool
		errText   string
	}{
		{
			name: "valid – creator == applicant, group exists, KYC active",
			msgs: []sdk.Msg{buildJoinGroupMsg(creator.Address, creator.Address, 1)},
			setup: func(dk *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(1)).Return(group, true)
				gk.EXPECT().GetDidAndKycActive(gomock.Any(), creator.GetAddress(), regionID).Return("did1", true)
			},
			expectErr: false,
		},
		{
			name: "valid – globalDAO creator acting for applicant",
			msgs: []sdk.Msg{buildJoinGroupMsg(daoAdmin.Address, applicant.Address, 1)},
			setup: func(dk *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				dk.EXPECT().IsGlobalDao(gomock.Any(), daoAdmin.Address).Return(true)
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(1)).Return(group, true)
				gk.EXPECT().GetDidAndKycActive(gomock.Any(), applicant.GetAddress(), regionID).Return("did2", true)
			},
			expectErr: false,
		},
		{
			name: "valid – meidDAO creator acting for applicant",
			msgs: []sdk.Msg{buildJoinGroupMsg(daoAdmin.Address, applicant.Address, 1)},
			setup: func(dk *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				dk.EXPECT().IsGlobalDao(gomock.Any(), daoAdmin.Address).Return(false)
				dk.EXPECT().GetMeidDao(gomock.Any()).Return(daoAdmin.Address)
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(1)).Return(group, true)
				gk.EXPECT().GetDidAndKycActive(gomock.Any(), applicant.GetAddress(), regionID).Return("did3", true)
			},
			expectErr: false,
		},
		{
			name: "non-MsgJoinGroup messages are skipped",
			msgs: []sdk.Msg{}, // empty tx – no MsgJoinGroup
			setup: func(_ *mock.MockDaoKeeper, _ *mock.MockMeGroupKeeper) {
				// no keeper calls expected
			},
			expectErr: false,
		},
		{
			name:      "group_id == 0 is rejected",
			msgs:      []sdk.Msg{buildJoinGroupMsg(creator.Address, creator.Address, 0)},
			setup:     func(_ *mock.MockDaoKeeper, _ *mock.MockMeGroupKeeper) {},
			expectErr: true,
			errText:   "group_id must be greater than 0",
		},
		{
			name:      "empty applicant_address is rejected",
			msgs:      []sdk.Msg{buildJoinGroupMsg(creator.Address, "", 1)},
			setup:     func(_ *mock.MockDaoKeeper, _ *mock.MockMeGroupKeeper) {},
			expectErr: true,
			errText:   "applicant_address is required",
		},
		{
			name:      "malformed applicant_address is rejected",
			msgs:      []sdk.Msg{buildJoinGroupMsg(creator.Address, "not-a-bech32", 1)},
			setup:     func(_ *mock.MockDaoKeeper, _ *mock.MockMeGroupKeeper) {},
			expectErr: true,
			errText:   "invalid applicant_address",
		},
		{
			name: "non-DAO creator acting for someone else is rejected",
			msgs: []sdk.Msg{buildJoinGroupMsg(daoAdmin.Address, applicant.Address, 1)},
			setup: func(dk *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				dk.EXPECT().IsGlobalDao(gomock.Any(), daoAdmin.Address).Return(false)
				dk.EXPECT().GetMeidDao(gomock.Any()).Return("some-other-address")
			},
			expectErr: true,
			errText:   "creator is neither the applicant nor a DAO admin",
		},
		{
			name: "non-existent group is rejected",
			msgs: []sdk.Msg{buildJoinGroupMsg(creator.Address, creator.Address, 99)},
			setup: func(_ *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(99)).Return(megrouptypes.GroupInfo{}, false)
			},
			expectErr: true,
			errText:   "group 99 does not exist",
		},
		{
			name: "applicant without active KYC is rejected",
			msgs: []sdk.Msg{buildJoinGroupMsg(creator.Address, creator.Address, 1)},
			setup: func(_ *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(1)).Return(group, true)
				gk.EXPECT().GetDidAndKycActive(gomock.Any(), creator.GetAddress(), regionID).Return("", false)
			},
			expectErr: true,
			errText:   "does not have active KYC",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, dk, gk := newDecorator(ctrl)
			tc.setup(dk, gk)

			ctx := newJoinGroupCtx()
			tx := &mock.MockTx{Msgs: tc.msgs}
			_, err := d.AnteHandle(ctx, tx, false, stubNext)

			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ─── nil-constructor panic guards ───────────────────────────────────────────

func TestNewJoinGroupValidateDecorator_PanicOnNilDaoKeeper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	gk := mock.NewMockMeGroupKeeper(ctrl)
	require.Panics(t, func() {
		ante.NewJoinGroupValidateDecorator(nil, gk)
	})
}

func TestNewJoinGroupValidateDecorator_PanicOnNilGroupKeeper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dk := mock.NewMockDaoKeeper(ctrl)
	require.Panics(t, func() {
		ante.NewJoinGroupValidateDecorator(dk, nil)
	})
}
