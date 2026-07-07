package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	govQueryCmd := &cobra.Command{
		Use:                        govtypes.ModuleName,
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return govQueryCmd
}
