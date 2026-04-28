package cli

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/openmetaearth/me-hub/x/megroup/types"
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

	cmd.AddCommand(CmdCreateGroup())
	//	cmd.AddCommand(CmdUpdateGroup())
	//	cmd.AddCommand(CmdDeleteGroup())
	cmd.AddCommand(CmdJoinGroup())
	cmd.AddCommand(CmdLeaveGroup())

	// this line is used by starport scaffolding # 1

	return cmd
}
