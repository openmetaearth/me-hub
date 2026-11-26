package ante_test

import (
	"bytes"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authantetestutil "github.com/cosmos/cosmos-sdk/x/auth/ante/testutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/ante"
	"github.com/openmetaearth/me-hub/app/ante/mock"
	"github.com/openmetaearth/me-hub/app/params"
)

type feeGrantErrorTx struct {
	mock.MockTx
	payer, granter []byte
	fees           sdk.Coins
}

func (tx feeGrantErrorTx) GetGas() uint64     { return 200000 }
func (tx feeGrantErrorTx) GetFee() sdk.Coins  { return tx.fees }
func (tx feeGrantErrorTx) FeePayer() []byte   { return tx.payer }
func (tx feeGrantErrorTx) FeeGranter() []byte { return tx.granter }

func TestFeeGrantErrorOutput(t *testing.T) {
	// These addresses contain non-text bytes, including DEL (0x7f), which
	// previously broke YAML output when formatted directly with %s.
	_, payer, err := bech32.DecodeAndConvert("me1vgkmtsg96f6nzp00q4lkscpxpfm4nfxyqns7yx")
	require.NoError(t, err)
	_, granter, err := bech32.DecodeAndConvert("me1jlruc674xjh9n7h522uktazzaek3qs0t4hhxxg")
	require.NoError(t, err)
	payerAddr, granterAddr := sdk.AccAddress(payer), sdk.AccAddress(granter)
	fees := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100000)))

	for _, tc := range []struct {
		name    string
		cause   error
		message string
	}{
		{
			name:    "missing allowance",
			cause:   sdkerrors.ErrNotFound.Wrap("fee-grant not found"),
			message: fmt.Sprintf("no fee allowance from granter %s to grantee %s; the granter must grant an allowance before --fee-granter can be used", granterAddr, payerAddr),
		},
		{
			name:    "other allowance failure",
			cause:   sdkerrors.ErrUnauthorized.Wrap("message type not allowed"),
			message: fmt.Sprintf("cannot use fee allowance from granter %s to grantee %s for fees %s", granterAddr, payerAddr, fees),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao := mock.NewMockDaoKeeper(ctrl)
			grants := authantetestutil.NewMockFeegrantKeeper(ctrl)
			dao.EXPECT().IsDao(gomock.Any(), payerAddr.String()).Return(false)
			dao.EXPECT().CheckFreeGasAccount(gomock.Any(), payerAddr.String()).Return(false)
			grants.EXPECT().UseGrantedFees(gomock.Any(), granterAddr, payerAddr, fees, gomock.Any()).Return(tc.cause)
			d := ante.NewDeductFeeDecorator(
				authantetestutil.NewMockAccountKeeper(ctrl), mock.NewMockBankKeeper(ctrl),
				grants, dao, mock.NewMockStakingKeeper(ctrl), mock.NewMockKycKeeper(ctrl),
				func(sdk.Context, sdk.Tx) (sdk.Coins, int64, error) { return fees, 0, nil },
				mock.NewMockWasmKeeper(ctrl),
			)
			tx := &feeGrantErrorTx{payer: payer, granter: granter, fees: fees}
			_, err := d.AnteHandle(sdk.Context{}, tx, false, func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
				t.Fatal("failed fee grant must stop ante processing")
				return ctx, nil
			})
			require.ErrorIs(t, err, tc.cause)
			require.EqualError(t, err, tc.message+": "+tc.cause.Error())

			// Exercise the same protobuf-to-YAML output path as the CLI.
			var output bytes.Buffer
			ctx := client.Context{}.
				WithCodec(codec.NewProtoCodec(codectypes.NewInterfaceRegistry())).
				WithOutput(&output).
				WithOutputFormat("text")
			require.NoError(t, ctx.PrintProto(&sdk.TxResponse{RawLog: err.Error()}))
			require.Contains(t, output.String(), payerAddr.String())
			require.Contains(t, output.String(), granterAddr.String())
		})
	}
}
