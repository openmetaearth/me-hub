package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
)

var (
	_ sdk.Msg = &MsgSubmitProposal{}
	_ sdk.Msg = &MsgVote{}
	_ sdk.Msg = &MsgSubmitEquityUpdate{}
	_ sdk.Msg = &MsgSubmitFreezeEquity{}
	_ sdk.Msg = &MsgSubmitUnfreezeEquity{}
	_ sdk.Msg = &MsgApprovePendingAction{}
	_ sdk.Msg = &MsgRejectPendingAction{}
	_ sdk.Msg = &MsgFreezeRegion{}
	_ sdk.Msg = &MsgUnfreezeRegion{}
	_ sdk.Msg = &MsgFreezeEquity{}
	_ sdk.Msg = &MsgUnfreezeEquity{}
)

const (
	TypeMsgSubmitProposal       = "submit_proposal"
	TypeMsgVote                 = "vote"
	TypeMsgSubmitEquityUpdate   = "submit_equity_update"
	TypeMsgSubmitFreezeEquity   = "submit_freeze_equity"
	TypeMsgSubmitUnfreezeEquity = "submit_unfreeze_equity"
	TypeMsgApprovePendingAction = "approve_pending_action"
	TypeMsgRejectPendingAction  = "reject_pending_action"
	TypeMsgFreezeRegion         = "freeze_region"
	TypeMsgUnfreezeRegion       = "unfreeze_region"
	TypeMsgFreezeEquity         = "freeze_equity"
	TypeMsgUnfreezeEquity       = "unfreeze_equity"
)

func (msg *MsgSubmitProposal) Route() string { return RouterKey }
func (msg *MsgSubmitProposal) Type() string  { return TypeMsgSubmitProposal }

func (msg *MsgSubmitProposal) GetSigners() []sdk.AccAddress { return signers(msg.Proposer) }
func (msg *MsgSubmitProposal) GetSignBytes() []byte         { return signBytes(msg) }

func (msg *MsgSubmitProposal) ValidateBasic() error {
	if err := validateAccAddress(msg.Proposer); err != nil {
		return err
	}
	if msg.Title == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "title is empty")
	}
	if msg.VotingPeriod < 1 || msg.VotingPeriod > MaxVotingPeriod {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "voting_period must be a number of seconds in [1,%d]", MaxVotingPeriod)
	}
	switch msg.Kind {
	case ProposalKindTreasuryWithdraw:
		if msg.ReplaceOperator != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "treasury proposal cannot include a replace-operator payload")
		}
		if msg.Treasury == nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "treasury payload is empty")
		}
		if err := validateWithdrawals(msg.Treasury.Withdrawals); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	case ProposalKindReplaceRegionOperator:
		if msg.Treasury != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "replace-operator proposal cannot include a treasury payload")
		}
		if err := validateReplaceOperator(msg.ReplaceOperator); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	default:
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown proposal kind")
	}
	return nil
}

func (msg *MsgVote) Route() string { return RouterKey }
func (msg *MsgVote) Type() string  { return TypeMsgVote }

func (msg *MsgVote) GetSigners() []sdk.AccAddress { return signers(msg.Voter) }
func (msg *MsgVote) GetSignBytes() []byte         { return signBytes(msg) }

func (msg *MsgVote) ValidateBasic() error {
	if err := validateAccAddress(msg.Voter); err != nil {
		return err
	}
	if msg.ProposalId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "proposal id is empty")
	}
	if msg.Option != VoteOptionYes && msg.Option != VoteOptionNo {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "vote option must be yes or no")
	}
	return nil
}

func (msg *MsgSubmitEquityUpdate) Route() string { return RouterKey }
func (msg *MsgSubmitEquityUpdate) Type() string  { return TypeMsgSubmitEquityUpdate }

func (msg *MsgSubmitEquityUpdate) GetSigners() []sdk.AccAddress { return signers(msg.Proposer) }
func (msg *MsgSubmitEquityUpdate) GetSignBytes() []byte         { return signBytes(msg) }

func (msg *MsgSubmitEquityUpdate) ValidateBasic() error {
	if err := validateAccAddress(msg.Proposer); err != nil {
		return err
	}
	if err := validateRegionID(msg.RegionId); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if err := validateUniqueAddresses(msg.Addresses); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (msg *MsgSubmitFreezeEquity) Route() string { return RouterKey }
func (msg *MsgSubmitFreezeEquity) Type() string  { return TypeMsgSubmitFreezeEquity }

func (msg *MsgSubmitFreezeEquity) GetSigners() []sdk.AccAddress {
	return signers(msg.Proposer)
}
func (msg *MsgSubmitFreezeEquity) GetSignBytes() []byte { return signBytes(msg) }

func (msg *MsgSubmitFreezeEquity) ValidateBasic() error {
	return validateEquityTargets(msg.Proposer, msg.RegionId, msg.Addresses)
}

func (msg *MsgSubmitUnfreezeEquity) Route() string { return RouterKey }
func (msg *MsgSubmitUnfreezeEquity) Type() string  { return TypeMsgSubmitUnfreezeEquity }

func (msg *MsgSubmitUnfreezeEquity) GetSigners() []sdk.AccAddress {
	return signers(msg.Proposer)
}
func (msg *MsgSubmitUnfreezeEquity) GetSignBytes() []byte { return signBytes(msg) }

func (msg *MsgSubmitUnfreezeEquity) ValidateBasic() error {
	return validateEquityTargets(msg.Proposer, msg.RegionId, msg.Addresses)
}

func (msg *MsgApprovePendingAction) Route() string { return RouterKey }
func (msg *MsgApprovePendingAction) Type() string  { return TypeMsgApprovePendingAction }

func (msg *MsgApprovePendingAction) GetSigners() []sdk.AccAddress {
	return signers(msg.Authority)
}
func (msg *MsgApprovePendingAction) GetSignBytes() []byte { return signBytes(msg) }

func (msg *MsgApprovePendingAction) ValidateBasic() error {
	return validateAuthorityAction(msg.Authority, msg.ActionId)
}

func (msg *MsgRejectPendingAction) Route() string { return RouterKey }
func (msg *MsgRejectPendingAction) Type() string  { return TypeMsgRejectPendingAction }

func (msg *MsgRejectPendingAction) GetSigners() []sdk.AccAddress {
	return signers(msg.Authority)
}
func (msg *MsgRejectPendingAction) GetSignBytes() []byte { return signBytes(msg) }

func (msg *MsgRejectPendingAction) ValidateBasic() error {
	return validateAuthorityAction(msg.Authority, msg.ActionId)
}

func (msg *MsgFreezeRegion) Route() string { return RouterKey }
func (msg *MsgFreezeRegion) Type() string  { return TypeMsgFreezeRegion }

func (msg *MsgFreezeRegion) GetSigners() []sdk.AccAddress { return signers(msg.Authority) }
func (msg *MsgFreezeRegion) GetSignBytes() []byte         { return signBytes(msg) }

func (msg *MsgFreezeRegion) ValidateBasic() error {
	return validateAuthorityRegions(msg.Authority, msg.RegionIds)
}

func (msg *MsgUnfreezeRegion) Route() string { return RouterKey }
func (msg *MsgUnfreezeRegion) Type() string  { return TypeMsgUnfreezeRegion }

func (msg *MsgUnfreezeRegion) GetSigners() []sdk.AccAddress {
	return signers(msg.Authority)
}
func (msg *MsgUnfreezeRegion) GetSignBytes() []byte { return signBytes(msg) }

func (msg *MsgUnfreezeRegion) ValidateBasic() error {
	return validateAuthorityRegions(msg.Authority, msg.RegionIds)
}

func (msg *MsgFreezeEquity) Route() string { return RouterKey }
func (msg *MsgFreezeEquity) Type() string  { return TypeMsgFreezeEquity }

func (msg *MsgFreezeEquity) GetSigners() []sdk.AccAddress { return signers(msg.Authority) }
func (msg *MsgFreezeEquity) GetSignBytes() []byte         { return signBytes(msg) }

func (msg *MsgFreezeEquity) ValidateBasic() error {
	return validateEquityTarget(msg.Authority, msg.RegionId, msg.Address)
}

func (msg *MsgUnfreezeEquity) Route() string { return RouterKey }
func (msg *MsgUnfreezeEquity) Type() string  { return TypeMsgUnfreezeEquity }

func (msg *MsgUnfreezeEquity) GetSigners() []sdk.AccAddress {
	return signers(msg.Authority)
}
func (msg *MsgUnfreezeEquity) GetSignBytes() []byte { return signBytes(msg) }

func (msg *MsgUnfreezeEquity) ValidateBasic() error {
	return validateEquityTarget(msg.Authority, msg.RegionId, msg.Address)
}

func signers(address string) []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func signBytes(msg proto.Message) []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func validateAuthorityAction(authority string, actionID uint64) error {
	if err := validateAccAddress(authority); err != nil {
		return err
	}
	if actionID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "action id is empty")
	}
	return nil
}

func validateAuthorityRegions(authority string, regionIDs []string) error {
	if err := validateAccAddress(authority); err != nil {
		return err
	}
	if len(regionIDs) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "region ids are empty")
	}
	seen := make(map[string]struct{}, len(regionIDs))
	for _, regionID := range regionIDs {
		if err := validateRegionID(regionID); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		id := strings.ToLower(regionID)
		if _, ok := seen[id]; ok {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate region id %s", regionID)
		}
		seen[id] = struct{}{}
	}
	return nil
}

func validateEquityTargets(signer, regionID string, addresses []string) error {
	if err := validateAccAddress(signer); err != nil {
		return err
	}
	if err := validateRegionID(regionID); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if err := validateEquityFreezeAddresses(addresses); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func validateEquityTarget(signer, regionID, address string) error {
	if err := validateAccAddress(signer); err != nil {
		return err
	}
	if err := validateRegionID(regionID); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if err := validateAccAddress(address); err != nil {
		return err
	}
	return nil
}
