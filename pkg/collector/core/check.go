package core

import (
	"github.com/anchnet/smartops-agent/pkg/metric"
	"time"
)

// Check is an interface for types capable to run checks
type Check interface {
	Collect(t time.Time) ([]metric.MetricSample, error) // collect the metric
	String() string                                     // provide a printable version of the check name
}
