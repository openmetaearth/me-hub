package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/st-chain/me-hub/x/streamer/client/cli"
)

var (
	CreateStreamHandler    = govclient.NewProposalHandler(cli.NewCmdSubmitCreateStreamProposal)
	TerminateStreamHandler = govclient.NewProposalHandler(cli.NewCmdSubmitTerminateStreamProposal)
	ReplaceStreamHandler   = govclient.NewProposalHandler(cli.NewCmdSubmitReplaceStreamDistributionProposal)
	UpdateStreamHandler    = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateStreamDistributionProposal)
)
