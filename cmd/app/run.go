package app

import (
	"fmt"
	"github.com/anchnet/smartops-agent/cmd/common"
	"github.com/anchnet/smartops-agent/pkg/collector"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/sender"
	log "github.com/cihub/seelog"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func init() {
	Command.AddCommand(runCmd)
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

	if err := startAgent(); err != nil {
		return err
	}

	select {
	case err := <-signalStop:
		return err

	}
}
func startAgent() error {
	if err := common.SetupConfig(confFilePath); err != nil {
		log.Errorf("Failed to setup config %v", err)
		return fmt.Errorf("ubable to set agent configuration: %v", err)
	}
	logFile := config.SmartOps.GetString("log_file")
	if logFile == "" {
		logFile = common.DefaultLogFile
	}
	logFileDir := filepath.Dir(logFile)
	if _, err := os.Stat(logFileDir); os.IsNotExist(err) {
		if err = os.MkdirAll(logFileDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create log dir: %v", err)
		}
	}

	if config.SmartOps.GetBool("disable_file_logging") {
		// this will prevent any logging on file
		logFile = ""
	}

	err := config.SetupLogger(
		loggerName,
		config.SmartOps.GetString("log_level"),
		logFile,
		config.SmartOps.GetBool("log_to_console"),
		config.SmartOps.GetBool("log_format_json"),
	)
	if err != nil {
		return fmt.Errorf("Error while setting up logging, exiting: %v", err)
	}

	log.Info("Starting SmartOps Agent...")

	// setup the sender
	send := sender.GetSender()
	if err := send.Connect(); err != nil {
		return fmt.Errorf("Error while sender connect, %v", err)
	}
	go func() {
		send.Run()
	}()
	log.Info("Start running...")
	collector.Collect()
	return nil
}
