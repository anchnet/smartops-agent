package app

import (
	"github.com/spf13/cobra"
)

func init() {
	CaptureCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run collector",
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	defer func() {
		// Stop Collector
	}()
}
