package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var CaptureCmd = &cobra.Command{
	Use:   fmt.Sprint("%s [command]", os.Args[0]),
	Short: "SmartOps Collector at your service.",
}
