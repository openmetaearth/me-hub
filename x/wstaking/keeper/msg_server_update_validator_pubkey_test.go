package keeper_test

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestUpdateValidatorPubkey() {
	s.SetupTest()

	muvp := types.NewMsgUpdateValidatorPubkey(
		s.Dao.GlobalDao,
		s.meEarthValidator.OperatorAddress,
		s.meEarthValidator.ConsensusPubkey,
	)
	_, err := s.msgServer.UpdateValidatorPubkey(s.Ctx, muvp)
	s.Require().NoError(err)

	s.Require().NoError(err)
	tests := []struct {
		name            string
		staker          string
		operatorAddress string
		pubkey          *codectypes.Any
		expErr          error
		malleate        func()
	}{
		{
			name:            "Dao Permission",
			staker:          s.Dao.MeidDao,
			operatorAddress: s.meEarthValidator.OperatorAddress,
			pubkey:          s.meEarthValidator.ConsensusPubkey,
			expErr:          types.ErrCheckGlobalDao,
		}, {
			name:            "wrong validator address",
			staker:          s.Dao.GlobalDao,
			operatorAddress: "mevaloper139mq752delxv78jvtmwxhasyrycufsvr707ate",
			pubkey:          s.meEarthValidator.ConsensusPubkey,
			expErr:          stakingtypes.ErrNoValidatorFound,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			newMuvp := types.NewMsgUpdateValidatorPubkey(
				test.staker,
				test.operatorAddress,
				test.pubkey,
			)
			_, err := s.msgServer.UpdateValidatorPubkey(s.Ctx, newMuvp)
			if test.expErr != nil {
				s.Require().ErrorContains(err, test.expErr.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}

}

func TestXxx(t *testing.T) {
	s:="qeWJwMjAL2BDMd+1TnI6agfDa8IGLeVwRXfawHlxZgc="
	ss,err := base64.StdEncoding.DecodeString(s)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(len(ss))
	sss:=sha256.Sum256(ss)
	fmt.Println(hex.EncodeToString(sss[:20]))
}
	
