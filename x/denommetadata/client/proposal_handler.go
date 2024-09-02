package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/st-chain/me-hub/x/denommetadata/client/cli"
)

var (
	CreateDenomMetadataHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCreateDenomMetadataProposal)
	UpdateDenomMetadataHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateDenomMetadataProposal)
)
