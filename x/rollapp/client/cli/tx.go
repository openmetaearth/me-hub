package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreateRollapp())
	cmd.AddCommand(CmdUpdateRollapp())
	cmd.AddCommand(CmdTransferOwnership())
	cmd.AddCommand(CmdAddApp())
	cmd.AddCommand(CmdUpdateApp())
	cmd.AddCommand(CmdRemoveApp())
	cmd.AddCommand(CmdToggleTEE())
	cmd.AddCommand(CmdSkipDelayRollapp())

	return cmd
}
