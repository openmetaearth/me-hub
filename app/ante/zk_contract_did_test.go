package ante

import (
	"errors"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

type zkContractDIDTestKeeper struct {
	protectedAddress common.Address
	allowed          bool
	allowErr         error
	allowChecks      *int
}

func (k zkContractDIDTestKeeper) IsZkContractAddress(_ sdk.Context, address common.Address) bool {
	return address == k.protectedAddress
}

func (k zkContractDIDTestKeeper) IsAllowToUseZkContractVerify(_ sdk.Context, _ sdk.AccAddress) (bool, error) {
	*k.allowChecks = *k.allowChecks + 1
	return k.allowed, k.allowErr
}

func (k zkContractDIDTestKeeper) GetDID(_ sdk.Context, _ sdk.AccAddress) (string, bool) {
	panic("GetDID must not be called by ZkContractDIDDecorator")
}

func (k zkContractDIDTestKeeper) GetDidInfo(_ sdk.Context, _ string) (didtypes.DidInfo, bool) {
	panic("GetDidInfo must not be called by ZkContractDIDDecorator")
}

type zkContractDIDTestTx struct {
	msgs []sdk.Msg
}

func (tx zkContractDIDTestTx) GetMsgs() []sdk.Msg {
	return tx.msgs
}

func (tx zkContractDIDTestTx) ValidateBasic() error {
	return nil
}

func newZkContractDIDTestMsg(target *common.Address, sender common.Address) *evmtypes.MsgEthereumTx {
	msg := evmtypes.NewTx(
		big.NewInt(1),
		0,
		target,
		big.NewInt(0),
		100_000,
		big.NewInt(1),
		nil,
		nil,
		nil,
		nil,
	)
	msg.From = sender.Hex()
	return msg
}

func TestZkContractDIDDecorator(t *testing.T) {
	protectedAddress := common.HexToAddress("0x0000000000000000000000000000000000001001")
	otherAddress := common.HexToAddress("0x0000000000000000000000000000000000002001")
	sender := common.HexToAddress("0x0000000000000000000000000000000000003001")
	permissionErr := errors.New("permission check failed")

	tests := []struct {
		name              string
		target            *common.Address
		allowed           bool
		allowErr          error
		wantError         bool
		wantErrorContains string
		wantAllowChecks   int
	}{
		{
			name:            "protected contract with permission",
			target:          &protectedAddress,
			allowed:         true,
			wantAllowChecks: 1,
		},
		{
			name:            "protected contract without permission",
			target:          &protectedAddress,
			wantError:       true,
			wantAllowChecks: 1,
		},
		{
			name:              "protected contract permission check failure",
			target:            &protectedAddress,
			allowErr:          permissionErr,
			wantError:         true,
			wantErrorContains: permissionErr.Error(),
			wantAllowChecks:   1,
		},
		{
			name:     "unprotected contract bypasses permission check",
			target:   &otherAddress,
			allowErr: permissionErr,
		},
		{
			name:     "contract creation bypasses permission check",
			allowErr: permissionErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			allowChecks := 0
			didKeeper := zkContractDIDTestKeeper{
				protectedAddress: protectedAddress,
				allowed:          tc.allowed,
				allowErr:         tc.allowErr,
				allowChecks:      &allowChecks,
			}
			decorator := NewZkContractDIDDecorator(didKeeper)
			tx := zkContractDIDTestTx{
				msgs: []sdk.Msg{newZkContractDIDTestMsg(tc.target, sender)},
			}
			nextCalled := false
			next := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
				nextCalled = true
				return ctx, nil
			}

			_, err := decorator.AnteHandle(sdk.Context{}, tx, false, next)

			if tc.wantError {
				require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
				if tc.wantErrorContains != "" {
					require.ErrorContains(t, err, tc.wantErrorContains)
				}
				require.False(t, nextCalled)
				require.Equal(t, tc.wantAllowChecks, allowChecks)
				return
			}
			require.NoError(t, err)
			require.True(t, nextCalled)
			require.Equal(t, tc.wantAllowChecks, allowChecks)
		})
	}
}
