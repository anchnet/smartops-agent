package main

import (
	"capture/cmd/app"
	"os"
)

func main() {
	if err := app.CaptureCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
