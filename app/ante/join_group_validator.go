package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	megrouptypes "github.com/openmetaearth/me-hub/x/megroup/types"
)

// JoinGroupValidateDecorator performs stateful validation on MsgJoinGroup messages
// before the fee-deduction step. This prevents attackers from exploiting the
// zero-fee exemption granted to MsgJoinGroup transactions by supplying invalid
// fields (non-existent group, missing KYC, unauthorized creator, etc.) that
// would otherwise only be caught inside the msg server after gas is consumed.
//
// The decorator is a no-op for transactions that contain no MsgJoinGroup messages.
type JoinGroupValidateDecorator struct {
	daoKeeper   DaoKeeper
	groupKeeper MeGroupKeeper
}

func NewJoinGroupValidateDecorator(dk DaoKeeper, gk MeGroupKeeper) JoinGroupValidateDecorator {
	if dk == nil {
		panic("JoinGroupValidateDecorator: daoKeeper must not be nil")
	}
	if gk == nil {
		panic("JoinGroupValidateDecorator: groupKeeper must not be nil")
	}
	return JoinGroupValidateDecorator{daoKeeper: dk, groupKeeper: gk}
}

func (d JoinGroupValidateDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		joinMsg, ok := msg.(*megrouptypes.MsgJoinGroup)
		if !ok {
			continue
		}
		if err := d.validateJoinGroup(ctx, joinMsg); err != nil {
			return ctx, err
		}
	}
	return next(ctx, tx, simulate)
}

func (d JoinGroupValidateDecorator) validateJoinGroup(ctx sdk.Context, msg *megrouptypes.MsgJoinGroup) error {
	// 1. GroupId must be non-zero (groups start from 1).
	if msg.GroupId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "group_id must be greater than 0")
	}

	// 2. ApplicantAddress is always required and must be a valid bech32 address.
	//    The msg server uses it unconditionally; an empty or malformed value always
	//    causes a failure — reject it here before we grant free gas.
	if msg.ApplicantAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "applicant_address is required")
	}
	applicant, err := sdk.AccAddressFromBech32(msg.ApplicantAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid applicant_address (%s)", err)
	}

	// 3. When the creator acts on behalf of someone else, the creator must be
	//    GlobalDao or MeidDao — same check as the msg server.
	if msg.ApplicantAddress != msg.Creator {
		if !d.daoKeeper.IsGlobalDao(ctx, msg.Creator) && d.daoKeeper.GetMeidDao(ctx) != msg.Creator {
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized,
				"creator is neither the applicant nor a DAO admin")
		}
	}

	// 4. The target group must already exist on-chain.
	groupInfo, found := d.groupKeeper.GetGroupInfo(ctx, msg.GroupId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress,
			"group %d does not exist", msg.GroupId)
	}

	// 5. Applicant must have an active level-2 KYC DID in the group's region.
	_, isKycActive := d.groupKeeper.GetDidAndKycActive(ctx, applicant, groupInfo.RegionID)
	if !isKycActive {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
			"applicant %s does not have active KYC in region %s",
			msg.ApplicantAddress, groupInfo.RegionID)
	}

	return nil
}
