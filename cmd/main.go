package main

import (
	"github.com/anchnet/smartops-agent/cmd/app"
	_ "github.com/anchnet/smartops-agent/pkg/collector/system"
	"os"
)

func main() {
	if err := app.Command.Execute(); err != nil {
		os.Exit(-1)
	}
}
