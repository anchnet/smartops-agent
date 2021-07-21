// +build windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

const (
	memMetric = "system.mem.%s"
)

type MemCheck struct {
	name string
}

func (c *MemCheck) Name() string {
	return c.name
}

func (c *MemCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("total"), float64(v.Total), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("free"), float64(v.Available), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("used"), float64(v.Used), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("used_pct"), v.UsedPercent, metric.UnitPercent, t, nil))

	return samples, nil
}

func (c MemCheck) formatMetric(name string) string {
	format := "system.mem.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	core.RegisterCheck(&MemCheck{
		name: "mem",
	})
}
