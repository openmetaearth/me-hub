package ante

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

type zkContractDIDTestKeeper struct {
	contractInfo   *didtypes.ZkContractAddressInfo
	allowed        bool
	allowErr       error
	checkedSenders *[]sdk.AccAddress
}

func (k zkContractDIDTestKeeper) IsZkContractAddress(_ sdk.Context, _ common.Address) bool {
	panic("IsZkContractAddress must not be called by ZkContractDIDDecorator")
}

func (k zkContractDIDTestKeeper) GetZkContractAddress(_ sdk.Context, contractType didtypes.ZkContractAddressType) *didtypes.ZkContractAddressInfo {
	if contractType != didtypes.PROXY_UNIVERSAL_VERIFIER {
		return nil
	}
	return k.contractInfo
}

func (k zkContractDIDTestKeeper) IsAllowToUseZkContractVerify(_ sdk.Context, sender sdk.AccAddress) (bool, error) {
	*k.checkedSenders = append(*k.checkedSenders, sender)
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

func newZkContractDIDTestMsg(target *common.Address, sender common.Address, data []byte) *evmtypes.MsgEthereumTx {
	msg := evmtypes.NewTx(
		big.NewInt(1),
		0,
		target,
		big.NewInt(0),
		100_000,
		big.NewInt(1),
		nil,
		nil,
		data,
		nil,
	)
	msg.From = sender.Hex()
	return msg
}

func submitResponseTestData(t *testing.T) []byte {
	t.Helper()
	// Encode the IVerifier ABI independently of the decorator's method selector.
	contractABI, err := abi.JSON(strings.NewReader(`[
		{"type":"function","name":"submitResponse","inputs":[
			{"name":"authResponse","type":"tuple","components":[
				{"name":"authMethod","type":"string"},{"name":"proof","type":"bytes"}]},
			{"name":"responses","type":"tuple[]","components":[
				{"name":"requestId","type":"uint256"},{"name":"proof","type":"bytes"},{"name":"metadata","type":"bytes"}]},
			{"name":"crossChainProofs","type":"bytes"}],"outputs":[]}
	]`))
	require.NoError(t, err)
	authResponse := struct {
		AuthMethod string
		Proof      []byte
	}{AuthMethod: "authV3", Proof: []byte{1, 2, 3}}
	responses := []struct {
		RequestId *big.Int
		Proof     []byte
		Metadata  []byte
	}{{RequestId: big.NewInt(1), Proof: []byte{4, 5, 6}}}
	data, err := contractABI.Pack("submitResponse", authResponse, responses, []byte{})
	require.NoError(t, err)
	return data
}

func TestZkContractDIDDecorator(t *testing.T) {
	protectedAddress := common.HexToAddress("0x000000000000000000000000000000000000AbCd")
	otherAddress := common.HexToAddress("0x0000000000000000000000000000000000002001")
	sender := common.HexToAddress("0x0000000000000000000000000000000000003001")
	permissionErr := errors.New("permission check failed")
	contractInfo := &didtypes.ZkContractAddressInfo{
		Type:            didtypes.PROXY_UNIVERSAL_VERIFIER,
		ContractAddress: strings.ToLower(protectedAddress.Hex()),
	}
	submitData := submitResponseTestData(t)
	otherMethodData := common.FromHex("0x8da5cb5b") // owner()

	tests := []struct {
		name              string
		target            *common.Address
		contractInfo      *didtypes.ZkContractAddressInfo
		data              []byte
		missingSender     bool
		allowed           bool
		allowErr          error
		wantError         bool
		wantErrorContains string
		wantAllowChecks   int
	}{
		{
			name:            "submitResponse with permission and lowercase configured address",
			target:          &protectedAddress,
			contractInfo:    contractInfo,
			data:            submitData,
			allowed:         true,
			wantAllowChecks: 1,
		},
		{
			name:            "submitResponse without permission",
			target:          &protectedAddress,
			contractInfo:    contractInfo,
			data:            submitData,
			wantError:       true,
			wantAllowChecks: 1,
		},
		{
			name:              "submitResponse permission check failure",
			target:            &protectedAddress,
			contractInfo:      contractInfo,
			data:              submitData,
			allowErr:          permissionErr,
			wantError:         true,
			wantErrorContains: permissionErr.Error(),
			wantAllowChecks:   1,
		},
		{
			name:         "submitResponse to another contract bypasses permission check",
			target:       &otherAddress,
			contractInfo: contractInfo,
			data:         submitData,
			allowErr:     permissionErr,
		},
		{
			name:         "contract creation bypasses permission check",
			contractInfo: contractInfo,
			data:         submitData,
			allowErr:     permissionErr,
		},
		{
			name:         "other method bypasses permission check",
			target:       &protectedAddress,
			contractInfo: contractInfo,
			data:         otherMethodData,
			allowErr:     permissionErr,
		},
		{
			name:         "unknown selector bypasses permission check",
			target:       &protectedAddress,
			contractInfo: contractInfo,
			data:         []byte{0xde, 0xad, 0xbe, 0xef},
			allowErr:     permissionErr,
		},
		{
			name:         "empty data bypasses permission check",
			target:       &protectedAddress,
			contractInfo: contractInfo,
			allowErr:     permissionErr,
		},
		{
			name:         "truncated selector bypasses permission check",
			target:       &protectedAddress,
			contractInfo: contractInfo,
			data:         submitData[:3],
			allowErr:     permissionErr,
		},
		{
			name:            "selector without arguments still requires permission",
			target:          &protectedAddress,
			contractInfo:    contractInfo,
			data:            submitData[:4],
			wantError:       true,
			wantAllowChecks: 1,
		},
		{
			name:     "missing verifier configuration bypasses permission check",
			target:   &protectedAddress,
			data:     submitData,
			allowErr: permissionErr,
		},
		{
			name:   "invalid configured address bypasses permission check",
			target: &common.Address{},
			contractInfo: &didtypes.ZkContractAddressInfo{
				Type: didtypes.PROXY_UNIVERSAL_VERIFIER,
			},
			data:     submitData,
			allowErr: permissionErr,
		},
		{
			name:              "submitResponse without sender is rejected",
			target:            &protectedAddress,
			contractInfo:      contractInfo,
			data:              submitData,
			missingSender:     true,
			wantError:         true,
			wantErrorContains: "sender is unavailable",
		},
		{
			name:          "other method does not require sender in this decorator",
			target:        &protectedAddress,
			contractInfo:  contractInfo,
			data:          otherMethodData,
			missingSender: true,
			allowErr:      permissionErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var checkedSenders []sdk.AccAddress
			didKeeper := zkContractDIDTestKeeper{
				contractInfo:   tc.contractInfo,
				allowed:        tc.allowed,
				allowErr:       tc.allowErr,
				checkedSenders: &checkedSenders,
			}
			decorator := NewZkContractDIDDecorator(didKeeper)
			msg := newZkContractDIDTestMsg(tc.target, sender, tc.data)
			if tc.missingSender {
				msg.From = ""
			}
			tx := zkContractDIDTestTx{
				msgs: []sdk.Msg{msg},
			}
			nextCalled := false
			next := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
				nextCalled = true
				return ctx, nil
			}

			_, err := decorator.AnteHandle(sdk.Context{}, tx, false, next)
			require.Len(t, checkedSenders, tc.wantAllowChecks)
			for _, checkedSender := range checkedSenders {
				require.Equal(t, sdk.AccAddress(sender.Bytes()), checkedSender)
			}

			if tc.wantError {
				require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
				if tc.wantErrorContains != "" {
					require.ErrorContains(t, err, tc.wantErrorContains)
				}
				require.False(t, nextCalled)
				return
			}
			require.NoError(t, err)
			require.True(t, nextCalled)
		})
	}
}
