package app

import (
	"github.com/spf13/cobra"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/core"
	"time"
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
	checks := core.LoadChecks()
	collector := collector.NewCollector()
	for _, c := range checks {
		collector.RunCheck(c)
	}
	time.Sleep(20 * time.Second)
	return nil
}
