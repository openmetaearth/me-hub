package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/openmetaearth/me-hub/app"
	"github.com/openmetaearth/me-hub/cmd/med/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
