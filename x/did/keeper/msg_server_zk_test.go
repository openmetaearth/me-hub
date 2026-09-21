package keeper

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"testing"
	"time"

	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/x/did/types"
)

const (
	testCreator = "me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4"
	testDID     = "1234567890123"
)

var (
	testStateAddress = common.HexToAddress("0x0000000000000000000000000000000000001001")
	testAuthAddress  = common.HexToAddress("0x0000000000000000000000000000000000001002")
)

type testDAOKeeper struct {
	allowed string
}

func (k testDAOKeeper) IsGlobalDao(_ sdk.Context, address string) bool {
	return address == k.allowed
}

type testAccountKeeper struct {
	account authtypes.AccountI
}

func (k testAccountKeeper) HasAccount(_ sdk.Context, address sdk.AccAddress) bool {
	return k.account != nil && bytes.Equal(k.account.GetAddress(), address)
}

func (k testAccountKeeper) GetAccount(ctx sdk.Context, address sdk.AccAddress) authtypes.AccountI {
	if !k.HasAccount(ctx, address) {
		return nil
	}
	return k.account
}

type testEVMKeeper struct {
	stateExists bool
	userID      *big.Int
	challenge   *big.Int
	calls       []common.Address
	gasLimits   []uint64
}

func (k *testEVMKeeper) StaticCall(
	_ sdk.Context,
	_ common.Address,
	contract common.Address,
	_ []byte,
	gas uint64,
) ([]byte, error) {
	k.calls = append(k.calls, contract)
	k.gasLimits = append(k.gasLimits, gas)
	switch contract {
	case testStateAddress:
		return stateABI.Methods["idExists"].Outputs.Pack(k.stateExists)
	case testAuthAddress:
		return authV3ValidatorABI.Methods["verify"].Outputs.Pack(
			k.userID,
			[]authResponseField{{Name: "challenge", Value: k.challenge}},
		)
	default:
		return nil, types.ErrZkContractAddressNotSet
	}
}

func newZkTestKeeper(
	t *testing.T,
	dao types.DaoKeeper,
	evm types.EVMKeeper,
) (*Keeper, sdk.Context) {
	t.Helper()
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	db := cometbftdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	privateKey, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	account := authtypes.NewBaseAccount(
		sdk.MustAccAddressFromBech32(testCreator), privateKey.PubKey(), 0, 0,
	)
	k := NewKeeper(cdc, storeKey, dao, evm, testAccountKeeper{account: account})
	ctx := sdk.NewContext(
		stateStore,
		cometbftproto.Header{
			ChainID: "me-hub-2401",
			Time:    time.Unix(1_800_000_000, 0),
		},
		false,
		log.NewNopLogger(),
	)
	return k, ctx
}

func testCoreID() []byte {
	coreID := make([]byte, types.ZkCoreIDLength)
	coreID[0] = 0x01
	coreID[1] = 0xc1
	for i := 2; i < types.ZkCoreIDLength-2; i++ {
		coreID[i] = byte(i)
	}
	var checksum uint16
	for _, b := range coreID[:types.ZkCoreIDLength-2] {
		checksum += uint16(b)
	}
	binary.LittleEndian.PutUint16(coreID[types.ZkCoreIDLength-2:], checksum)
	return coreID
}

func configureZkContracts(k *Keeper, ctx sdk.Context) {
	k.SetZkContractAddressInfo(ctx, types.ZkContractAddressInfo{
		Type:            types.PROXY_STATE,
		ContractAddress: testStateAddress.Hex(),
	})
	k.SetZkContractAddressInfo(ctx, types.ZkContractAddressInfo{
		Type:            types.PROXY_AUTH_V3_VALIDATOR,
		ContractAddress: testAuthAddress.Hex(),
	})
}

func addTestDID(k *Keeper, ctx sdk.Context) {
	creator := sdk.MustAccAddressFromBech32(testCreator)
	k.SetDID(ctx, creator, testDID)
	k.SetDidInfo(ctx, testDID, types.DidInfo{
		Did:      testDID,
		Address:  testCreator,
		Pubkey:   "test-pubkey",
		Status:   types.DID_STATUS_ACTIVE,
		KycLevel: types.KYC_LEVEL_TWO,
	})
}

func TestSetZkContractAddress(t *testing.T) {
	k, ctx := newZkTestKeeper(t, testDAOKeeper{allowed: testCreator}, nil)
	server := msgServer{Keeper: k}

	_, err := server.SetZkContractAddress(sdk.WrapSDKContext(ctx), &types.MsgSetZkContractAddressRequest{
		Creator: testCreator,
		ContractInfo: &types.ZkContractAddressInfo{
			Type:            types.PROXY_STATE,
			ContractAddress: "0x0000000000000000000000000000000000001001",
		},
	})
	require.NoError(t, err)
	stored := k.GetZkContractAddress(ctx, types.PROXY_STATE)
	require.True(t, stored != nil)
	require.Equal(t, testStateAddress.Hex(), stored.ContractAddress)

	_, err = server.SetZkContractAddress(sdk.WrapSDKContext(ctx), &types.MsgSetZkContractAddressRequest{
		Creator: sdk.AccAddress(bytes.Repeat([]byte{2}, 20)).String(),
		ContractInfo: &types.ZkContractAddressInfo{
			Type:            types.PROXY_STATE,
			ContractAddress: testStateAddress.Hex(),
		},
	})
	require.ErrorIs(t, err, types.ErrPermissionDenial)
}

func TestBondZkCoreID(t *testing.T) {
	coreID := testCoreID()
	coreIDInt, err := types.ZkCoreIDToInt(coreID)
	require.NoError(t, err)
	expiresAt := int64(1_900_000_000)
	challenge, err := types.CalculateZkCoreIDBindingChallenge(
		"me-hub-2401", testCreator, testDID, coreID, 7, expiresAt,
	)
	require.NoError(t, err)

	evm := &testEVMKeeper{stateExists: true, userID: coreIDInt, challenge: challenge}
	k, ctx := newZkTestKeeper(t, nil, evm)
	configureZkContracts(k, ctx)
	addTestDID(k, ctx)
	server := msgServer{Keeper: k}
	request := &types.MsgBondZkCoreIDRequest{
		Creator: testCreator,
		CoreInfo: &types.ZkAuthClaimCoreIdInfo{
			CoreId:    coreID,
			AuthProof: []byte{1, 2, 3},
			Nonce:     7,
			ExpiresAt: expiresAt,
		},
	}

	_, err = server.BondZkCoreID(sdk.WrapSDKContext(ctx), request)
	require.NoError(t, err)
	require.Equal(t, []common.Address{testStateAddress, testAuthAddress}, evm.calls)
	require.Equal(t,
		[]uint64{stateIDExistsCallGas, authV3VerifyCallGas},
		evm.gasLimits,
	)
	stored, found := k.GetZkCoreIDInfo(ctx, testDID)
	require.True(t, found)
	require.Equal(t, request.CoreInfo.CoreId, stored.CoreId)
	require.Equal(t, request.CoreInfo.AuthProof, stored.AuthProof)
	boundDID, found := k.GetDIDByZkCoreID(ctx, coreID)
	require.True(t, found)
	require.Equal(t, testDID, boundDID)

	_, err = server.BondZkCoreID(sdk.WrapSDKContext(ctx), request)
	require.ErrorIs(t, err, types.ErrDidAlreadyBound)
}

func TestBondZkCoreIDRejectsInvalidContractResults(t *testing.T) {
	coreID := testCoreID()
	coreIDInt, err := types.ZkCoreIDToInt(coreID)
	require.NoError(t, err)
	expiresAt := int64(1_900_000_000)

	tests := []struct {
		name        string
		stateExists bool
		userID      *big.Int
		challenge   *big.Int
		expectedErr error
	}{
		{"missing state", false, coreIDInt, big.NewInt(1), types.ErrInvalidZkCoreID},
		{"wrong challenge", true, coreIDInt, big.NewInt(1), types.ErrZkChallengeMismatch},
		{"wrong core ID", true, big.NewInt(123), big.NewInt(1), types.ErrZkCoreIDMismatch},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			evm := &testEVMKeeper{
				stateExists: tc.stateExists,
				userID:      tc.userID,
				challenge:   tc.challenge,
			}
			k, ctx := newZkTestKeeper(t, nil, evm)
			configureZkContracts(k, ctx)
			addTestDID(k, ctx)
			server := msgServer{Keeper: k}
			_, err := server.BondZkCoreID(sdk.WrapSDKContext(ctx), &types.MsgBondZkCoreIDRequest{
				Creator: testCreator,
				CoreInfo: &types.ZkAuthClaimCoreIdInfo{
					CoreId: coreID, AuthProof: []byte{1}, Nonce: 7, ExpiresAt: expiresAt,
				},
			})
			require.ErrorIs(t, err, tc.expectedErr)
			_, found := k.GetZkCoreIDInfo(ctx, testDID)
			require.False(t, found)
		})
	}
}
