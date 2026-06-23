package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ethsecp256k1 "github.com/evmos/ethermint/crypto/ethsecp256k1"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

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
	const did = "1111111111111"
	kycAddr, kycPubkey := s.NewAccount()
	s.setupActiveKyc(kycAddr, kycPubkey, did)

	s.Run("success", func() {
		subAddr, subPubkey := s.newEthSubAccount()

		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr.String(),
			Did:        did,
			SubAccount: subAddr.String(),
			Pubkey:     subPubkey,
		})
		s.Require().NoError(err)

		info, found := s.Keeper().GetDidInfo(s.Ctx, did)
		s.Require().True(found)
		s.Require().Equal(subAddr.String(), info.SubAccount)
	})

	s.Run("did not found", func() {
		subAddr, subPubkey := s.newEthSubAccount()

		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr.String(),
			Did:        "9999999999999",
			SubAccount: subAddr.String(),
			Pubkey:     subPubkey,
		})
		s.Require().ErrorIs(err, didtypes.ErrHolderNotFound)
	})

	s.Run("no kyc credential", func() {
		// DidInfo exists and is active, but no KYC credential was stored.
		const didNoKyc = "2222222222222"
		noKycAddr, noKycPubkey := s.NewAccount()
		s.Keeper().SetDID(s.Ctx, noKycAddr, didNoKyc)
		s.Keeper().SetDidInfo(s.Ctx, didNoKyc, didtypes.DidInfo{
			Did:     didNoKyc,
			Address: noKycAddr.String(),
			Pubkey:  noKycPubkey,
			Status:  didtypes.DID_STATUS_ACTIVE,
		})

		subAddr, subPubkey := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    noKycAddr.String(),
			Did:        didNoKyc,
			SubAccount: subAddr.String(),
			Pubkey:     subPubkey,
		})
		s.Require().ErrorIs(err, didtypes.ErrCredentialNotFound)
	})

	s.Run("sub_account already set on did", func() {
		const did2 = "3333333333333"
		kycAddr2, kycPubkey2 := s.NewAccount()
		s.setupActiveKyc(kycAddr2, kycPubkey2, did2)

		subAddr, subPubkey := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr2.String(),
			Did:        did2,
			SubAccount: subAddr.String(),
			Pubkey:     subPubkey,
		})
		s.Require().NoError(err)

		// second attempt on the same DID must fail
		subAddr2, subPubkey2 := s.newEthSubAccount()
		_, err = s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr2.String(),
			Did:        did2,
			SubAccount: subAddr2.String(),
			Pubkey:     subPubkey2,
		})
		s.Require().ErrorIs(err, types.ErrSubAccountAlreadyExists)
	})

	s.Run("sub_account already registered by another did", func() {
		const did3 = "4444444444444"
		kycAddr3, kycPubkey3 := s.NewAccount()
		s.setupActiveKyc(kycAddr3, kycPubkey3, did3)

		subAddr, subPubkey := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr3.String(),
			Did:        did3,
			SubAccount: subAddr.String(),
			Pubkey:     subPubkey,
		})
		s.Require().NoError(err)

		// same sub_account address attempted for a different DID
		const did4 = "5555555555555"
		kycAddr4, kycPubkey4 := s.NewAccount()
		s.setupActiveKyc(kycAddr4, kycPubkey4, did4)

		_, err = s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr4.String(),
			Did:        did4,
			SubAccount: subAddr.String(),
			Pubkey:     subPubkey,
		})
		s.Require().ErrorIs(err, types.ErrSubAccountAlreadyRegistered)
	})

	s.Run("invalid pubkey string", func() {
		const did5 = "6666666666666"
		kycAddr5, kycPubkey5 := s.NewAccount()
		s.setupActiveKyc(kycAddr5, kycPubkey5, did5)

		subAddr, _ := s.newEthSubAccount()
		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr5.String(),
			Did:        did5,
			SubAccount: subAddr.String(),
			Pubkey:     "not-valid-json",
		})
		s.Require().ErrorIs(err, sdkerrors.ErrInvalidPubKey)
	})

	s.Run("wrong key type secp256k1 rejected", func() {
		const did6 = "7777777777777"
		kycAddr6, kycPubkey6 := s.NewAccount()
		s.setupActiveKyc(kycAddr6, kycPubkey6, did6)

		secpPrivKey := secp256k1.GenPrivKey()
		secpAddr := sdk.AccAddress(secpPrivKey.PubKey().Address())
		secpPubkeyJSON, err := s.App.AppCodec().MarshalInterfaceJSON(secpPrivKey.PubKey())
		s.Require().NoError(err)

		_, err = s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr6.String(),
			Did:        did6,
			SubAccount: secpAddr.String(),
			Pubkey:     string(secpPubkeyJSON),
		})
		s.Require().ErrorIs(err, sdkerrors.ErrInvalidPubKey)
	})

	s.Run("pubkey does not match sub_account address", func() {
		const did7 = "8888888888888"
		kycAddr7, kycPubkey7 := s.NewAccount()
		s.setupActiveKyc(kycAddr7, kycPubkey7, did7)

		_, subPubkey := s.newEthSubAccount()
		differentAddr, _ := s.newEthSubAccount()

		_, err := s.msgServer.CreateSubAccount(s.Ctx, &types.MsgCreateSubAccount{
			Creator:    kycAddr7.String(),
			Did:        did7,
			SubAccount: differentAddr.String(),
			Pubkey:     subPubkey,
		})
		s.Require().ErrorIs(err, sdkerrors.ErrInvalidPubKey)
	})
}
