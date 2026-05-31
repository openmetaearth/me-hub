package types

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateSequencer_ValidateBasic(t *testing.T) {
	pubkey := ed25519.GenPrivKey().PubKey()
	pkAny, err := codectypes.NewAnyWithValue(pubkey)
	require.NoError(t, err)

	invalidkey := "{\"@type\":\"/cosmos.crypto.ed25519.PubKey\",\"key\":\"OcEwSZhPfddSUr84dkfj6Sfsh6PDSkcBdySUFxPb0Fs=\"}"
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	codec := codec.NewProtoCodec(interfaceRegistry)
	var invalidpk cryptotypes.PubKey
	err = codec.UnmarshalInterfaceJSON([]byte(invalidkey), &invalidpk)
	require.NoError(t, err)
	pkInvalid, err := codectypes.NewAnyWithValue(invalidpk)
	require.NoError(t, err)

	bond := sdk.NewCoin("stake", sdk.NewInt(100))

	tests := []struct {
		name string
		msg  MsgCreateSequencer
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateSequencer{
				Creator:      "invalid_address",
				DymintPubKey: pkAny,
				Bond:         bond,
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkAny,
				Bond:         bond,
			},
		}, {
			name: "valid description",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkAny,
				Bond:         bond,
				Description: Description{
					Moniker:         strings.Repeat("a", MaxMonikerLength),
					Identity:        strings.Repeat("a", MaxIdentityLength),
					Website:         strings.Repeat("a", MaxWebsiteLength),
					SecurityContact: strings.Repeat("a", MaxSecurityContactLength),
					Details:         strings.Repeat("a", MaxDetailsLength),
				},
			},
		}, {
			name: "invalid moniker length",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkAny,
				Bond:         bond,
				Description: Description{
					Moniker: strings.Repeat("a", MaxMonikerLength+1),
				},
			},
			err: ErrInvalidRequest,
		}, {
			name: "invalid identity length",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkAny,
				Bond:         bond,
				Description: Description{
					Identity: strings.Repeat("a", MaxIdentityLength+1),
				},
			},
			err: ErrInvalidRequest,
		}, {
			name: "invalid website length",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkAny,
				Bond:         bond,
				Description: Description{
					Website: strings.Repeat("a", MaxWebsiteLength+1),
				},
			},
			err: ErrInvalidRequest,
		}, {
			name: "invalid security contact length",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkAny,
				Bond:         bond,
				Description: Description{
					SecurityContact: strings.Repeat("a", MaxSecurityContactLength+1),
				},
			},
			err: ErrInvalidRequest,
		}, {
			name: "invalid details length",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkAny,
				Bond:         bond,
				Description: Description{
					Details: strings.Repeat("a", MaxDetailsLength+1),
				},
			},
			err: ErrInvalidRequest,
		}, {
			name: "invalid bond",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkAny,
				Bond:         sdk.Coin{Denom: "k", Amount: sdk.NewInt(0)},
			},
			err: ErrInvalidCoins,
		}, {
			name: "invalid public key",
			msg: MsgCreateSequencer{
				Creator:      sample.AccAddress(),
				DymintPubKey: pkInvalid,
				Bond:         bond,
			},
			err: ErrInvalidPubKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgReplaceProposerRequestValidateBasic(t *testing.T) {
	creator := sample.AccAddress()
	oldProposer := sample.AccAddress()
	newProposer := sample.AccAddress()

	tests := []struct {
		name      string
		msg       MsgReplaceProposerRequest
		wantErr   error
		errSubstr string
	}{
		{
			name: "valid replace proposer",
			msg: MsgReplaceProposerRequest{
				Creator: creator,
				ReplaceProposer: &MsgRepalceProposer{
					RollappId:   "rollapp_1234-1",
					OldProposer: oldProposer,
					NewProposer: newProposer,
					BlockHeight: 10,
				},
			},
		},
		{
			name: "reject self replacement",
			msg: MsgReplaceProposerRequest{
				Creator: creator,
				ReplaceProposer: &MsgRepalceProposer{
					RollappId:   "rollapp_1234-1",
					OldProposer: oldProposer,
					NewProposer: oldProposer,
					BlockHeight: 10,
				},
			},
			wantErr:   sdkerrors.ErrInvalidRequest,
			errSubstr: "old proposer and new proposer cannot be the same address",
		},
		{
			name: "reject missing replace proposer",
			msg: MsgReplaceProposerRequest{
				Creator: creator,
			},
			wantErr:   sdkerrors.ErrInvalidRequest,
			errSubstr: "ReplaceProposer can not",
		},
		{
			name: "reject missing rollapp id",
			msg: MsgReplaceProposerRequest{
				Creator: creator,
				ReplaceProposer: &MsgRepalceProposer{
					OldProposer: oldProposer,
					NewProposer: newProposer,
					BlockHeight: 10,
				},
			},
			wantErr:   sdkerrors.ErrInvalidRequest,
			errSubstr: "rollapp id cannot be empty",
		},
		{
			name: "reject invalid old proposer",
			msg: MsgReplaceProposerRequest{
				Creator: creator,
				ReplaceProposer: &MsgRepalceProposer{
					RollappId:   "rollapp_1234-1",
					OldProposer: "not-a-bech32-address",
					NewProposer: newProposer,
					BlockHeight: 10,
				},
			},
			wantErr:   sdkerrors.ErrInvalidAddress,
			errSubstr: "invalid OldProposer address",
		},
		{
			name: "reject invalid new proposer",
			msg: MsgReplaceProposerRequest{
				Creator: creator,
				ReplaceProposer: &MsgRepalceProposer{
					RollappId:   "rollapp_1234-1",
					OldProposer: oldProposer,
					NewProposer: "not-a-bech32-address",
					BlockHeight: 10,
				},
			},
			wantErr:   sdkerrors.ErrInvalidAddress,
			errSubstr: "invalid NewProposer address",
		},
		{
			name: "reject invalid block height",
			msg: MsgReplaceProposerRequest{
				Creator: creator,
				ReplaceProposer: &MsgRepalceProposer{
					RollappId:   "rollapp_1234-1",
					OldProposer: oldProposer,
					NewProposer: newProposer,
				},
			},
			wantErr:   sdkerrors.ErrInvalidRequest,
			errSubstr: "invalid block number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.ErrorContains(t, err, tt.errSubstr)
				return
			}
			require.NoError(t, err)
		})
	}
}
