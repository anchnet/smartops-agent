// +build windows

package system

import (
	"github.com/anchnet/smartops-agent/pkg/metric"
	"time"
)

const (
	loadMetric = "system.load.%s"
)

func runLoadCheck(time time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	return samples, nil
}
