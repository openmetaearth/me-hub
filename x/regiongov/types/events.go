package types

const (
	EventTypeSubmitProposal       = "submit_proposal"
	EventTypeVote                 = "vote"
	EventTypeVetoProposal         = "veto_proposal"
	EventTypeTallyProposal        = "tally_proposal"
	EventTypeSubmitEquityUpdate   = "submit_equity_update"
	EventTypeSubmitFreezeEquity   = "submit_freeze_equity"
	EventTypeSubmitUnfreezeEquity = "submit_unfreeze_equity"
	EventTypeApprovePendingAction = "approve_pending_action"
	EventTypeRejectPendingAction  = "reject_pending_action"
	EventTypeFreezeEquity         = "freeze_equity"
	EventTypeUnfreezeEquity       = "unfreeze_equity"
	EventTypeFreezeRegion         = "freeze_region"
	EventTypeUnfreezeRegion       = "unfreeze_region"

	AttributeKeyProposalID = "proposal_id"
	AttributeKeyActionID   = "action_id"
	AttributeKeyProposer   = "proposer"
	AttributeKeyVoter      = "voter"
	AttributeKeyOption     = "option"
	AttributeKeyStatus     = "status"
	AttributeKeyElectorate = "electorate"
	AttributeKeyRegionID   = "region_id"
	AttributeKeyAddress    = "address"
)
