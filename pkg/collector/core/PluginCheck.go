package core

import (
	"github.com/anchnet/smartops-agent/pkg/metric"
	"time"
)

// Check is an interface for types capable to run checks
type PluginCheck interface {
	PluginCollect(t time.Time) ([]metric.MetricSample, error) // collect the metric
	PluginName() string                                       // provide a printable version of the check name
}
