package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "med",
		Short: "Med command-line tool",
	}

	rootCmd.AddCommand(claimRewardCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}