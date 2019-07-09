package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/core"
	"os"
	"os/signal"
	"syscall"
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
	signalOS := make(chan os.Signal, 1)
	signal.Notify(signalOS, os.Interrupt, syscall.SIGTERM)
	signalStop := make(chan error)
	go func() {
		select {
		case sig := <-signalOS:
			fmt.Println("Received signal '%s', shutting down...", sig)
			signalStop <- nil
		}
	}()

	checks := core.LoadChecks()
	collector := collector.NewCollector()
	for _, c := range checks {
		collector.RunCheck(c)
	}

	select {
	case err := <-signalStop:
		return err

	}
	return nil
}
