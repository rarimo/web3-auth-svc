package main

import (
	"os"

	"github.com/rarimo/web3-auth-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
