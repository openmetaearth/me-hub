package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

const TypeMsgSetZkContractAddress = "set_zk_contract_address"

func NewMsgSetZkContractAddress(
	creator string,
	contractInfo *ZkContractAddressInfo,
) *MsgSetZkContractAddressRequest {
	return &MsgSetZkContractAddressRequest{Creator: creator, ContractInfo: contractInfo}
}

func (m *MsgSetZkContractAddressRequest) Route() string { return RouterKey }

func (m *MsgSetZkContractAddressRequest) Type() string { return TypeMsgSetZkContractAddress }

func (m *MsgSetZkContractAddressRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (m *MsgSetZkContractAddressRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetZkContractAddressRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "creator is not a valid bech32 address")
	}
	if m.ContractInfo == nil {
		return errors.Wrap(ErrParameter, "contract info is required")
	}
	if _, ok := ZkContractAddressType_name[int32(m.ContractInfo.Type)]; !ok {
		return errors.Wrap(ErrParameter, "unknown zk contract address type")
	}
	if !common.IsHexAddress(m.ContractInfo.ContractAddress) {
		return ErrInvalidZkContractAddress
	}
	if common.HexToAddress(m.ContractInfo.ContractAddress) == (common.Address{}) {
		return errors.Wrap(ErrInvalidZkContractAddress, "zero address is not allowed")
	}
	return nil
}
