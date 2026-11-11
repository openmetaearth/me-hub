package ante_test

import (
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/evmos/ethermint/testutil"
	"github.com/golang/mock/gomock"

	"github.com/openmetaearth/me-hub/app"
	"github.com/openmetaearth/me-hub/app/params"
	regiongovtypes "github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (suite *AnteTestSuite) TestRegionGovVoteAccountSignature() {
	for _, tc := range []struct {
		name           string
		wrongPublicKey bool
		wrongSignature bool
	}{
		{name: "voter signs with their own key"},
		{name: "another account cannot claim the voter", wrongPublicKey: true},
		{name: "voter public key with another private key", wrongSignature: true},
	} {
		suite.Run(tc.name, func() {
			suite.SetupTest(false)
			voterKey := secp256k1.GenPrivKey()
			otherKey := secp256k1.GenPrivKey()
			voter := sdk.AccAddress(voterKey.PubKey().Address())
			acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, voter)
			suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			suite.Require().NoError(testutil.FundAccount(suite.app.BankKeeper, suite.ctx, voter,
				sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100))))

			feeReceiver := NewAccount()
			suite.mockDaoKeeper.EXPECT().IsDao(gomock.Any(), voter.String()).Return(false)
			suite.mockDaoKeeper.EXPECT().CheckFreeGasAccount(gomock.Any(), voter.String()).Return(false)
			suite.mockDaoKeeper.EXPECT().GetDevOperator(gomock.Any()).Return(feeReceiver.Address)
			suite.mockStakingKeeper.EXPECT().GetProposerOwnerAddress(gomock.Any()).Return(feeReceiver.Address, nil)
			suite.mockDaoKeeper.EXPECT().GetGlobalDaoFeePoolAddr(gomock.Any()).Return(feeReceiver.GetAddress())

			msg := &regiongovtypes.MsgVote{Voter: voter.String(), ProposalId: 1, Option: regiongovtypes.VoteOptionYes}
			suite.Require().NoError(msg.ValidateBasic())
			suite.Require().Equal([]sdk.AccAddress{voter}, msg.GetSigners())
			// GetTxConfig is an SDK test helper with only SDK messages registered;
			// use the application's encoding config for regiongov transaction decoding.
			txConfig := app.MakeEncodingConfig().TxConfig
			builder := txConfig.NewTxBuilder()
			suite.Require().NoError(builder.SetMsgs(msg))
			builder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 20)))
			builder.SetGasLimit(200_000)
			publicKey := voterKey.PubKey()
			if tc.wrongPublicKey {
				publicKey = otherKey.PubKey()
			}
			suite.Require().NoError(builder.SetSignatures(signing.SignatureV2{
				PubKey: publicKey, Sequence: acc.GetSequence(),
				Data: &signing.SingleSignatureData{SignMode: signing.SignMode_SIGN_MODE_DIRECT},
			}))
			privateKey := voterKey
			if tc.wrongPublicKey || tc.wrongSignature {
				privateKey = otherKey
			}
			signature, err := clienttx.SignWithPrivKey(signing.SignMode_SIGN_MODE_DIRECT, authsigning.SignerData{
				Address: voter.String(), ChainID: suite.ctx.ChainID(),
				AccountNumber: acc.GetAccountNumber(), Sequence: acc.GetSequence(),
			}, builder, privateKey, txConfig, acc.GetSequence())
			suite.Require().NoError(err)
			signature.PubKey = publicKey
			suite.Require().NoError(builder.SetSignatures(signature))
			wire, err := txConfig.TxEncoder()(builder.GetTx())
			suite.Require().NoError(err)
			decoded, err := txConfig.TxDecoder()(wire)
			suite.Require().NoError(err)

			// Non-simulated DeliverTx ante verification checks the real signature.
			// Keeper tests separately check electorate membership and freeze rules.
			_, err = suite.anteHandler(suite.ctx, decoded, false)
			switch {
			case tc.wrongPublicKey:
				suite.Require().ErrorIs(err, sdkerrors.ErrInvalidPubKey)
			case tc.wrongSignature:
				suite.Require().ErrorIs(err, sdkerrors.ErrUnauthorized)
			default:
				suite.Require().NoError(err)
				suite.Require().Equal(acc.GetSequence()+1, suite.app.AccountKeeper.GetAccount(suite.ctx, voter).GetSequence())
			}
		})
	}
}
