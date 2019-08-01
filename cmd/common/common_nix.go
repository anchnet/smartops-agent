// +build !windows

package common

const (
	// DefaultConfPath points to the folder containing smartops.yaml
	DefaultConfPath = "etc"
	// DefaultLogFile points to the log file that will be used if not configured
	DefaultLogFile = "var/log/agent.log"
)
