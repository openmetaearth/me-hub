package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/megroup/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group megroup queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListGroup())
	cmd.AddCommand(CmdShowGroup())
	cmd.AddCommand(CmdListGroupMember())
	cmd.AddCommand(CmdShowGroupMember())
	//	cmd.AddCommand(CmdListMemberJoined())
	cmd.AddCommand(CmdShowMemberJoinedGroup())
	//	cmd.AddCommand(CmdListGroupMemberCount())
	// this line is used by starport scaffolding # 1

	return cmd
}
