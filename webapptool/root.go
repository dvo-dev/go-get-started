package webapptool

import (
	"github.com/spf13/cobra"
)

// This is the base command, i.e. called without subcommands
var RootCmd = &cobra.Command{
	Use:   "webapp-tool",
	Short: "webapp CLI tool",
	Long:  "webapp CLI tool called with no args...", // TODO: add list of args
}
