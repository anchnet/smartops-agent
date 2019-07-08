package app

import (
	"github.com/spf13/cobra"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/check"
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
	loader := core.GoCheckLoader{}
	checkKeys := core.GetRegisteredFactoryKeys()
	var allChecks []check.Check
	for _, k := range checkKeys {
		if checker, _ := loader.Load(k); checker != nil {
			allChecks = append(allChecks, checker)
		}
	}
	for {
		for _, c := range allChecks {
			_ = c.Run()
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}
