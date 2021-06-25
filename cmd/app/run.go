package app

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/anchnet/smartops-agent/cmd/common"
	"github.com/anchnet/smartops-agent/pkg/collector"
	"github.com/anchnet/smartops-agent/pkg/collector/filter"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/http"
	"github.com/anchnet/smartops-agent/pkg/pidfile"
	"github.com/cihub/seelog"
	log "github.com/cihub/seelog"
	"github.com/spf13/cobra"
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
		return fmt.Errorf("Error while setting up logging, exiting: %v", err)
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

	// init nginx configs
	if err := common.SetupNgxConfig(common.DefaultNgxConfPath); err != nil {
		// update plugin
		if err := http.UpsertPlugins("nginx", false); err != nil {
			log.Infof("Update online monitor plugins failed!")
		}
		_ = log.Errorf("Failed to setup config %v", err)
	} else {
		if err := http.UpsertPlugins("nginx", true); err != nil {
			log.Infof("Create nginx online  plugins failed!")
		}
	}
	if err := common.SetUpMysqlConfig(common.DefaultMysqlConfPath); err != nil {
		if err := http.UpsertPlugins("mysql", false); err != nil {
			log.Infof("Update online mysql plugins failed !")
		}
		_ = log.Errorf("Failed to setup config %v", err)

	} else {
		if err := http.UpsertPlugins("mysql", true); err != nil {
			log.Infof("Create online mysql plugin failed!")
		}
	}
	// setup filter
	// data := []byte(`{
	// 	"cpu": ["system.cpu.idle","system.cpu.used","system.cpu.system","system.cpu.iowait"],
	// 	"proc":["system.proc.cpu.usage","system.proc.cpu.system","system.proc.mem.rss","system.proc.thread.count","system.proc.mem.vms","system.proc.mem.pct"],
	// 	"mem":["system.mem.committed_as","system.mem.commit_limit","system.mem.page_tables","system.mem.slab","system.mem.shared","system.mem.buffered","system.mem.total","system.mem.free","system.mem.used","system.mem.used_pct", "system.mem.cached"],
	// 	"disk":["system.disk.total", "system.disk.used","system.disk.free","system.disk.used.pct"],
	//	"alarm": ["system.alarm.info","system.alarm.user_count"]
	// }`)
	// byts, err := http.GetFilter()
	// //FIXME:
	byts := []byte(`{
		"alarm": ["system.alarm.info","system.alarm.user_count"],
			"cpu": ["system.cpu.idle","system.cpu.used","system.cpu.system","system.cpu.iowait"],
			"mem":["system.mem.committed_as","system.mem.commit_limit","system.mem.page_tables","system.mem.slab","system.mem.shared","system.mem.buffered","system.mem.total","system.mem.free","system.mem.used","system.mem.used_pct", "system.mem.cached"],
			"disk":["system.disk.total", "system.disk.used","system.disk.free","system.disk.used.pct"]
	}`)
	seelog.Info("Filter data: ", string(byts))
	if err != nil {
		return seelog.Error(err)
	}
	if err := filter.SetFilter(byts); err != nil {
		return log.Errorf("error set filter: %v", err)
	}

	// setup the forwarder
	if err := forwarder.GetDefaultForwarder().Start(); err != nil {
		return log.Errorf("error start forwarder: %v", err)
	}

	// setup the collector
	go collector.Collect()
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
