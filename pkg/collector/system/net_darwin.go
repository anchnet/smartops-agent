//+build darwin

package system

import (
	"github.com/anchnet/smartops-agent/pkg/metric"
	"time"
)

func runNetworkCheck(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	return samples, nil
}
