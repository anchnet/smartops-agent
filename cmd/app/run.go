package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.51idc.com/smartops/smartcat-agent/cmd/common"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/core"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/sender"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/util/log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	CaptureCmd.AddCommand(runCmd)
}

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run collector",
		RunE:  run,
	}
)

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
			log.Infof("Received signal '%s', shutting down...", sig)
			signalStop <- nil
		}
	}()

	if err := common.SetupConfig(confFilePath); err != nil {
		log.Errorf("Failed to setup config %v", err)
		return fmt.Errorf("ubable to set agent configuration: %v", err)
	}

	// setup the sender
	send := sender.GetSender()
	go func() {
		send.Run()
	}()

	checks := core.LoadChecks()
	coll := collector.NewCollector()
	for _, c := range checks {
		coll.RunCheck(c)
	}

	select {
	case err := <-signalStop:
		return err

	}
}
