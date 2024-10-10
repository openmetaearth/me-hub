package types

const (
	EvtStaking               = "EventStaking"
	EvtUnStaking             = "EventUnStaking"
	EvtRegisterRollappID     = "EventRegisterRollappID"
	EvtFirstElectionTime     = "EventFirstElectionTime"
	EvtElection              = "EventElection"
	EvtProcUnStake           = "EventProcUnStake"
	EvtProcUnStakeStatistics = "EventProcUnStakeStatistics"
	EvtPunishment            = "EventPunishment"
	EvtSequencerChange       = "EventSequencerChange"
	EvtRestStakeInfo         = "EventRestStakeInfo"
	//========da fraud
	EvtPunishDaChallengerFraud = "EventPunishDaChallengerFraud" //DA挑战者欺诈事件
	EvtPunishBlockDaSubmitter  = "EventPunishBlockDaSubmitter"  //提交BlockDa欺诈
	EvtAddBlackList            = "EventAddBlackList"
	EvtRewardDaFraudChallenger = "EventRewardDaFraudChallenger"
	EvtRewardDaFraudValidator  = "EventRewardDaFraudValidator"
)
