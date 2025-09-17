package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/st-chain/me-hub/x/gravity/types"
)

func GetTxCmd(subCmd ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Crosschain transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(subCmd...)
	return cmd
}

func GetTxSubCmds(chainName string) []*cobra.Command {
	cmds := []*cobra.Command{}
	for _, command := range cmds {
		flags.AddTxFlagsToCmd(command)
	}
	return cmds
}
