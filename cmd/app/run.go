package app

import (
	"fmt"
	"github.com/anchnet/smartops-agent/cmd/common"
	"github.com/anchnet/smartops-agent/pkg/collector"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/http"
	"github.com/anchnet/smartops-agent/pkg/pidfile"
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
		stopAgent()
	}()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	errorCh := make(chan error)
	go func() {
		sig := <-signalCh
		log.Infof("Receive signal '%s', shutting down...", sig)
		errorCh <- nil
	}()

	if err := startAgent(); err != nil {
		return err
	}

	err := <-errorCh
	return err
}

func startAgent() error {
	if err := common.SetupConfig(confFilePath); err != nil {
		_ = log.Errorf("Failed to setup config %v", err)
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
		return fmt.Errorf("error while setting up logging, exiting: %v", err)
	}

	log.Info("Starting SmartOps Agent...")

	err = pidfile.WritePID(common.DefaultPidFile)
	if err != nil {
		return log.Errorf("Error while writing PID file, exiting: %v", err)
	}
	log.Infof("pid '%d' written to pid file '%s'", os.Getpid(), common.DefaultPidFile)

	// cache dir
	cacheDir := filepath.Dir(common.DefaultCacheDir)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err = os.MkdirAll(cacheDir, os.ModePerm); err != nil {
			return fmt.Errorf("create cache dir error, %v", err)
		}
	}

	// validate api_key
	err = http.ValidateAPIKey()
	if err != nil {
		return log.Errorf("validate api_key error, %v", err)
	}
	log.Infof("API key validate success.")

	// setup the forwarder
	if err := forwarder.GetDefaultForwarder().Start(); err != nil {
		return log.Errorf("error start forwarder: %v", err)
	}

	// setup the collector
	go collector.Collect()

	// setup plugins

	log.Info("Start running...")
	return nil
}

func stopAgent() {
	log.Info("Stopping agent...")
	_ = os.Remove(common.DefaultPidFile)
	collector.Stop()
	if err := forwarder.GetDefaultForwarder().Stop(); err != nil {
		_ = log.Errorf("error while closing connection, %v", err)
	}
}

/*
func heartbeat() {
	url := fmt.Sprintf("%s%s?endpoint=%s", fh.domain, agentHealthCheckEndpoint, config.SmartOps.GetString("endpoint"))
	transport := util.CreateHttpTransport()
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		_ = log.Errorf("heartbeat network error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		_ = log.Errorf("heartbeat error: %v", string(bytes))
	}
}
*/
