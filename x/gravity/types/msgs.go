package types

import (
	"fmt"
	"sort"
)

type MsgValidateBasic interface {
	MsgRelayerSetConfirmValidate(m *MsgRelayerSetConfirm) (err error)
	MsgRelayerSetUpdatedClaimValidate(m *MsgRelayerSetUpdatedClaim) (err error)
	MsgBridgeTokenClaimValidate(m *MsgBridgeTokenClaim) (err error)
	MsgSendToExternalClaimValidate(m *MsgSendToExternalClaim) (err error)

	MsgSendToMeClaimValidate(m *MsgSendToMeClaim) (err error)
	MsgSendToExternalValidate(m *MsgSendToExternal) (err error)

	MsgRequestBatchValidate(m *MsgRequestBatch) (err error)
	MsgConfirmBatchValidate(m *MsgConfirmBatch) (err error)

	ValidateAddress(addr string) error
	AddressToBytes(addr string) ([]byte, error)
}

var msgValidateBasicRouter = make(map[string]MsgValidateBasic)

func MustGetMsgValidateBasic(chainName string) MsgValidateBasic {
	mvb, ok := msgValidateBasicRouter[chainName]
	if !ok {
		panic(fmt.Sprintf("chain %s validate basic not found", chainName))
	}
	return mvb
}

func GetValidateChains() []string {
	chains := make([]string, 0, len(msgValidateBasicRouter))
	for chainName := range msgValidateBasicRouter {
		chains = append(chains, chainName)
	}
	sort.SliceStable(chains, func(i, j int) bool {
		return chains[i] < chains[j]
	})
	return chains
}

func RegisterValidateBasic(chainName string, validate MsgValidateBasic) {
	if _, ok := msgValidateBasicRouter[chainName]; ok {
		panic(fmt.Sprintf("duplicate registry msg validateBasic! chainName: %s", chainName))
	}
	msgValidateBasicRouter[chainName] = validate
}
