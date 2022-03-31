package common

import (
	"github.com/anchnet/smartops-agent/pkg/util/executable"
)

const (
	Version = "2.0.7"
)

var (
	// utility variables
	_here, _ = executable.Folder()
	// Root directory
	RootDirectory = _here
)

const (

	// DefaultNgxConfPath points to the folder containing nginx.yaml
	DefaultNgxConfPath = "conf/plugins.d/nginx.d/"
	// DefaultMysqlConfPath points to the folder containing mysql.yaml
	DefaultMysqlConfPath = "conf/plugins.d/mysql.d/"
	// DefaultConfPath points to the folder containing smartops.yaml
	DefaultConfPath = "conf"
	// DefaultLogFile points to the log file that will be used if not configured
	DefaultLogFile = "var/log/agent.log"
	// Default systemd
	DefaultSystemdPath = "/lib/systemd/system/smartops-agent.service"
	// Default upstart
	DefaultUpstartPath = "/etc/init/smartops-agent.conf"
	// Default pid file
	DefaultPidFile = "var/run/agent.pid"
	// Cache dir
	DefaultCacheDir = "var/cache/"
)

var (
	LocalMetricListen    string = "127.0.0.1:48001"
	LocalMetricHandle    string = "/localmetric"
	PhysicalDeviceHandle string = "/physical_device"
)
