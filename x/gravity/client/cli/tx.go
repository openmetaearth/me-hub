package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"strings"
)

func GetTxCmd(moduleName string, subNames ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      fmt.Sprintf("%s%s transaction subcommands", strings.ToUpper(moduleName[:1]), moduleName[1:]),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	for _, chainName := range subNames {
		cmd.AddCommand(GetTxCmd(chainName))
	}
	if len(subNames) == 0 {
		cmd.AddCommand(getTxSubCmds(moduleName)...)
	}
	return cmd
}

func getTxSubCmds(chainName string) []*cobra.Command {
	cmds := []*cobra.Command{}
	for _, command := range cmds {
		flags.AddTxFlagsToCmd(command)
	}
	return cmds
}
