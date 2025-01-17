package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/st-chain/me-hub/x/rollapp/client/cli"
)

var SubmitFraudHandler = govclient.NewProposalHandler(cli.NewCmdSubmitFraudProposal)
