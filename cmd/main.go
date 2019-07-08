package main

import (
	"os"
	"smartdog-agent/cmd/app"
	_ "smartdog-agent/pkg/collector/system"
)

func main() {
	if err := app.CaptureCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
