package main

import (
	"fmt"
	"os"

	"github.com/dvo-dev/go-get-started/webapptool"
)

func main() {
	if err := webapptool.RootCmd.Execute(); err != nil {
		fmt.Printf("error running RootCmd: %v", err)
		os.Exit(1)
	}
}
