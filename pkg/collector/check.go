package collector

import "github.com/anchnet/smartops-agent/pkg/metric"

// Check is an interface for types capable to run checks
type Check interface {
	Run() []metric.MetricSample // run the check
	//Stop()                                               // stop the check if it's running
	//String() string // provide a printable version of the check name
	//Configure(config, initConfig integration.Data) error // configure the check from the outside
	//Interval() time.Duration // return the interval time for the check
	//GetWarnings() []error                                // return the last warning registered by the check
	//GetMetricStats() (map[string]int64, error)           // get metric stats from the sender
	//Version() string                                     // return the version of the check if available
}
