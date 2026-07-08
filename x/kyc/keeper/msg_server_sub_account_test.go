package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

// newSecp256k1UserAccount creates a user account with secp256k1 pubkey on chain.
func (s *KeeperTestSuite) newSecp256k1UserAccount() (sdk.AccAddress, string) {
	privKey := secp256k1.GenPrivKey()
	addr := sdk.AccAddress(privKey.PubKey().Address())
	pubkeyJSON, err := s.App.AppCodec().MarshalInterfaceJSON(privKey.PubKey())
	s.Require().NoError(err)

	acc := s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, addr)
	err = acc.SetPubKey(privKey.PubKey())
	s.Require().NoError(err)
	s.App.AccountKeeper.SetAccount(s.Ctx, acc)

	return addr, string(pubkeyJSON)
}

// newEthUserAccount creates a user account with ethsecp256k1 pubkey on chain.
func (s *KeeperTestSuite) newEthUserAccount() (sdk.AccAddress, string) {
	privKey, err := ethsecp256k1.GenerateKey()
	s.Require().NoError(err)
	addr := sdk.AccAddress(privKey.PubKey().Address())
	pubkeyJSON, err := s.App.AppCodec().MarshalInterfaceJSON(privKey.PubKey())
	s.Require().NoError(err)

	acc := s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, addr)
	err = acc.SetPubKey(privKey.PubKey())
	s.Require().NoError(err)
	s.App.AccountKeeper.SetAccount(s.Ctx, acc)

	return addr, string(pubkeyJSON)
}

// newAccountWithoutPubKey creates an on-chain account without a pubkey set.
func (s *KeeperTestSuite) newAccountWithoutPubKey() (sdk.AccAddress, string) {
	privKey := secp256k1.GenPrivKey()
	addr := sdk.AccAddress(privKey.PubKey().Address())
	pubkeyJSON, err := s.App.AppCodec().MarshalInterfaceJSON(privKey.PubKey())
	s.Require().NoError(err)

	acc := s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, addr)
	s.App.AccountKeeper.SetAccount(s.Ctx, acc)

	return addr, string(pubkeyJSON)
}

// newEthSubAccount generates a fresh ethsecp256k1 keypair and returns the
// cosmos address together with the pubkey encoded as proto-JSON, the format
// expected by Keeper.PubKeyFromString.
func (s *KeeperTestSuite) newEthSubAccount() (sdk.AccAddress, string) {
	privKey, err := ethsecp256k1.GenerateKey()
	s.Require().NoError(err)
	addr := sdk.AccAddress(privKey.PubKey().Address())
	pubkeyJSON, err := s.App.AppCodec().MarshalInterfaceJSON(privKey.PubKey())
	s.Require().NoError(err)
	return addr, string(pubkeyJSON)
}

// setupActiveKyc writes a DID + active DidInfo + KYC credential directly into
// the keeper stores, bypassing the Approve message handler (which calls
// KycReward and hits validator meid limits in tests).
func (s *KeeperTestSuite) setupActiveKyc(addr sdk.AccAddress, pubkey, did string) {
	s.Keeper().SetDID(s.Ctx, addr, did)
	s.Keeper().SetDidInfo(s.Ctx, did, didtypes.DidInfo{
		Did:     did,
		Address: addr.String(),
		Pubkey:  pubkey,
		Status:  didtypes.DID_STATUS_ACTIVE,
	})
	kyc := didtypes.NewCredential(did, types.ModuleName, "hash", "uri", []byte("meearth"))
	s.Keeper().SetKYC(s.Ctx, did, kyc)
}

func (s *KeeperTestSuite) TestCreateSubAccount() {
	daoCreator := s.Dao.GlobalDao

	s.Run("success", func() {
		const did = "1111111111111"
		userAddr, userPubkey := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr, userPubkey, did)

		subAddr, subPubkey := s.newEthSubAccount()

		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().NoError(err)

		info, found := s.Keeper().GetDidInfo(s.Ctx, did)
		s.Require().True(found)
		s.Require().Equal(subAddr.String(), info.SubAccount)

		subDid, ok := s.Keeper().GetSubAccountDidMap(s.Ctx, subAddr.String())
		s.Require().True(ok)
		s.Require().Equal(did, subDid)

		account := s.App.AccountKeeper.GetAccount(s.Ctx, subAddr)
		s.Require().NotNil(account)
		s.Require().NotNil(account.GetPubKey())
	})

	s.Run("non-dao creator unauthorized", func() {
		const did = "1000000000001"
		userAddr, userPubkey := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr, userPubkey, did)

		subAddr, subPubkey := s.newEthSubAccount()
		nonDaoAddr, _ := s.newSecp256k1UserAccount()

		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          nonDaoAddr.String(),
			Account:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().ErrorIs(err, types.ErrUnauthorized)
	})

	s.Run("account did not found", func() {
		userAddr, _ := s.newSecp256k1UserAccount()
		subAddr, subPubkey := s.newEthSubAccount()

		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().ErrorIs(err, didtypes.ErrDidNotFound)
	})

	s.Run("no kyc credential", func() {
		const didNoKyc = "2222222222222"
		userAddr, userPubkey := s.newSecp256k1UserAccount()
		s.Keeper().SetDID(s.Ctx, userAddr, didNoKyc)
		s.Keeper().SetDidInfo(s.Ctx, didNoKyc, didtypes.DidInfo{
			Did:     didNoKyc,
			Address: userAddr.String(),
			Pubkey:  userPubkey,
			Status:  didtypes.DID_STATUS_ACTIVE,
		})

		subAddr, subPubkey := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().ErrorIs(err, didtypes.ErrCredentialNotFound)
	})

	s.Run("sub_account already set on did", func() {
		const did2 = "3333333333333"
		userAddr, userPubkey := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr, userPubkey, did2)

		subAddr, subPubkey := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().NoError(err)

		subAddr2, subPubkey2 := s.newEthSubAccount()
		_, err = s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       subAddr2.String(),
			SubAccountPubkey: subPubkey2,
		})
		s.Require().ErrorIs(err, types.ErrSubAccountAlreadyExists)
	})

	s.Run("sub_account already registered by another did", func() {
		const did3 = "4444444444444"
		userAddr3, userPubkey3 := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr3, userPubkey3, did3)

		subAddr, subPubkey := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr3.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().NoError(err)

		const did4 = "5555555555555"
		userAddr4, userPubkey4 := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr4, userPubkey4, did4)

		_, err = s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr4.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().ErrorIs(err, types.ErrSubAccountAlreadyRegistered)
	})

	s.Run("invalid sub_account pubkey string", func() {
		const did5 = "6666666666666"
		userAddr, userPubkey := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr, userPubkey, did5)

		subAddr, _ := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: "not-valid-json",
		})
		s.Require().ErrorIs(err, sdkerrors.ErrInvalidPubKey)
	})

	s.Run("account pubkey not set", func() {
		const did8 = "1010101010101"
		userAddr, userPubkey := s.newAccountWithoutPubKey()
		s.setupActiveKyc(userAddr, userPubkey, did8)

		subAddr, subPubkey := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().ErrorIs(err, sdkerrors.ErrInvalidPubKey)
		s.Require().Contains(err.Error(), "unknown account pubkey type")
	})

	s.Run("eth account not allowed to create sub account", func() {
		const didEth = "9999999999999"
		userAddr, userPubkey := s.newEthUserAccount()
		s.setupActiveKyc(userAddr, userPubkey, didEth)

		subAddr, subPubkey := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       subAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().ErrorIs(err, types.ErrEthAccountNotAllowed)
	})

	s.Run("sub_account pubkey must be ethsecp256k1", func() {
		const did6 = "7777777777777"
		userAddr, userPubkey := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr, userPubkey, did6)

		secpPrivKey := secp256k1.GenPrivKey()
		secpAddr := sdk.AccAddress(secpPrivKey.PubKey().Address())
		secpPubkeyJSON, err := s.App.AppCodec().MarshalInterfaceJSON(secpPrivKey.PubKey())
		s.Require().NoError(err)

		_, err = s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       secpAddr.String(),
			SubAccountPubkey: string(secpPubkeyJSON),
		})
		s.Require().ErrorIs(err, sdkerrors.ErrInvalidPubKey)
	})

	s.Run("sub_account pubkey does not match sub_account address", func() {
		const did7 = "8888888888888"
		userAddr, userPubkey := s.newSecp256k1UserAccount()
		s.setupActiveKyc(userAddr, userPubkey, did7)

		_, subPubkey := s.newEthSubAccount()
		differentAddr, _ := s.newEthSubAccount()

		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:          daoCreator,
			Account:          userAddr.String(),
			SubAccount:       differentAddr.String(),
			SubAccountPubkey: subPubkey,
		})
		s.Require().ErrorIs(err, sdkerrors.ErrInvalidPubKey)
	})
}
