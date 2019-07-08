package main

import (
	"gitlab.51idc.com/smartops/smartcat-agent/cmd/app"
	_ "gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/system"
	"os"
)

func main() {
	if err := app.CaptureCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
