// +build windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/mem"
	"time"
)

const (
	memMetric = "system.mem.%s"
)

type MemCheck struct {
	name string
}

func (c *MemCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "total"), float64(v.Total), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "free"), float64(v.Available), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "used"), float64(v.Used), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "used_pct"), v.UsedPercent, metric.UnitPercent, t, nil))

	return samples, nil
}

func (c MemCheck) formatMetric(name string) string {
	format := "system.mem.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	c := &MemCheck{
		CheckBase: core.NewCheckBase("men"),
	}
	core.RegisterCheck(c)
}
