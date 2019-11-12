// +build !windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/load"
	"time"
)

type LoadCheck struct {
	name string
}

func (c *LoadCheck) Name() string {
	return c.name
}

func (c *LoadCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	avg, err := load.Avg()
	if err != nil {
		return nil, err
	}
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("1"), avg.Load1, "", t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("5"), avg.Load5, "", t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("15"), avg.Load15, "", t, nil))
	return samples, nil
}

func (c LoadCheck) formatMetric(name string) string {
	format := "system.load.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	core.RegisterCheck(&LoadCheck{
		name: "load",
	})
}
