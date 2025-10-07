package main

import (
	"os"

	"github.com/opd-ai/asset-generator/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
