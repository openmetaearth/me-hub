package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNotImplemented        = sdkerrors.Register(ModuleName, 1, "regiongov message is not enabled yet")
	ErrUnauthorized          = sdkerrors.Register(ModuleName, 2, "unauthorized")
	ErrRegionFrozen          = sdkerrors.Register(ModuleName, 3, "region is frozen")
	ErrProposalNotFound      = sdkerrors.Register(ModuleName, 4, "proposal not found")
	ErrVotingClosed          = sdkerrors.Register(ModuleName, 5, "voting is closed")
	ErrAlreadyVoted          = sdkerrors.Register(ModuleName, 6, "voter has already voted")
	ErrEmptyElectorate       = sdkerrors.Register(ModuleName, 7, "electorate is empty")
	ErrInvalidVotingPeriod   = sdkerrors.Register(ModuleName, 8, "voting period is below the minimum")
	ErrUnknownRegion         = sdkerrors.Register(ModuleName, 9, "unknown region")
	ErrInvalidProposal       = sdkerrors.Register(ModuleName, 10, "invalid proposal")
	ErrPendingActionNotFound = sdkerrors.Register(ModuleName, 11, "pending action not found")
	ErrPendingActionClosed   = sdkerrors.Register(ModuleName, 12, "pending action is closed")
	ErrNotEquityAddress      = sdkerrors.Register(ModuleName, 13, "address is not an equity address of the region")
	ErrEquityNotFrozen       = sdkerrors.Register(ModuleName, 14, "equity address is not frozen")
	ErrEquityFrozen          = sdkerrors.Register(ModuleName, 15, "equity address is frozen")
)
