package app

import (
	"github.com/spf13/cobra"
	"smart-capture/pkg/collector/core"
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
	for {
		for _, k := range checkKeys {
			check, _ := loader.Load(k)
			if check != nil {
				_ = check.Run()
			}
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}
