package main

import (
	"fmt"
	"os"

	"github.com/dvo-dev/go-get-started/cmd/webapptool/subcommands"
)

func main() {
	if err := subcommands.RootCmd.Execute(); err != nil {
		fmt.Printf("error running RootCmd: %v", err)
		os.Exit(1)
	}
}
