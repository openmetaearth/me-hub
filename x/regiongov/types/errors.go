package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrNotImplemented        = errorsmod.Register(ModuleName, 1, "regiongov message is not enabled yet")
	ErrUnauthorized          = errorsmod.Register(ModuleName, 2, "unauthorized")
	ErrRegionFrozen          = errorsmod.Register(ModuleName, 3, "region is frozen")
	ErrProposalNotFound      = errorsmod.Register(ModuleName, 4, "proposal not found")
	ErrVotingClosed          = errorsmod.Register(ModuleName, 5, "voting is closed")
	ErrAlreadyVoted          = errorsmod.Register(ModuleName, 6, "voter has already voted")
	ErrEmptyElectorate       = errorsmod.Register(ModuleName, 7, "electorate is empty")
	ErrInvalidVotingPeriod   = errorsmod.Register(ModuleName, 8, "voting period is below the minimum")
	ErrUnknownRegion         = errorsmod.Register(ModuleName, 9, "unknown region")
	ErrInvalidProposal       = errorsmod.Register(ModuleName, 10, "invalid proposal")
	ErrPendingActionNotFound = errorsmod.Register(ModuleName, 11, "pending action not found")
	ErrPendingActionClosed   = errorsmod.Register(ModuleName, 12, "pending action is closed")
	ErrNotEquityAddress      = errorsmod.Register(ModuleName, 13, "address is not an equity address of the region")
	ErrEquityNotFrozen       = errorsmod.Register(ModuleName, 14, "equity address is not frozen")
	ErrEquityFrozen          = errorsmod.Register(ModuleName, 15, "equity address is frozen")
	ErrProposalExecution     = errorsmod.Register(ModuleName, 16, "proposal execution failed")
	ErrEquityWithdrawLimit   = errorsmod.Register(ModuleName, 17, "equity withdrawal exceeds the region treasury limit")
)
