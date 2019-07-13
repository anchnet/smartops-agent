package main

import (
	"github.com/anchnet/smartops-agent/cmd/app"
	"os"
)

func main() {
	if err := app.Command.Execute(); err != nil {
		os.Exit(-1)
	}
}
