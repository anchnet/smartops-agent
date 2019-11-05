// +build !windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"time"
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
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("total"), float64(v.Total), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("free"), float64(v.Free), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("used"), float64(v.Used), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("used_pct"), v.UsedPercent, metric.UnitPercent, t, nil))

	switch runtime.GOOS {
	case "linux":
		samples = append(samples, c.linuxMemCheck(v, t)...)
	}
	return samples, nil
}

func (c MemCheck) linuxMemCheck(v *mem.VirtualMemoryStat, t time.Time) []metric.MetricSample {
	var samples []metric.MetricSample
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("cached"), float64(v.Cached), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("buffered"), float64(v.Buffers), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("shared"), float64(v.Shared), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("slab"), float64(v.Slab), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("page_tables"), float64(v.PageTables), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("commit_limit"), float64(v.CommitLimit), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatMemMetric("committed_as"), float64(v.CommittedAS), metric.UnitByte, t, nil))
	samples = append(samples, metric.NewServerMetricSample(c.formatSwapMetric("cached"), float64(v.SwapCached), metric.UnitByte, t, nil))
	return samples
}
func (c MemCheck) formatMemMetric(name string) string {
	format := "system.mem.%s"
	return fmt.Sprintf(format, name)
}

func (c MemCheck) formatSwapMetric(name string) string {
	format := "system.swap.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	c := &MemCheck{
		CheckBase: core.NewCheckBase("mem"),
	}
	core.RegisterCheck(c)
}
