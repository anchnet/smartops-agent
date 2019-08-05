// +build !windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/load"
	"time"
)

const (
	loadMetric = "system.load.%s"
)

func runLoadCheck(time time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	avg, err := load.Avg()
	if err != nil {
		return nil, err
	}
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(loadMetric, "1"), avg.Load1, "", time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(loadMetric, "1"), avg.Load1, "", time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(loadMetric, "1"), avg.Load1, "", time, nil))
	return samples, nil
}
