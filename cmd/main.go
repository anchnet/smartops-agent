package main

import (
	"os"
	"smart-capture/cmd/app"
	_ "smart-capture/pkg/collector/system"
)

func main() {
	if err := app.CaptureCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
