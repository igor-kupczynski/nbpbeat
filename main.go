package main

import (
	"os"

	"github.com/igor-kupczynski/nbpbeat/cmd"

	_ "github.com/igor-kupczynski/nbpbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
