// +build windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

type CPUCheck struct {
	core.CheckBase
	lastCycle    float64
	lastCPUTimes cpu.TimesStat
}

func (c *CPUCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	cpuTimes, _ := cpu.Times(false)
	timesStat := cpuTimes[0]
	cycle := timesStat.Total()

	if c.lastCycle != 0 {
		toPercent := 100 / (cycle - c.lastCycle)
		user := ((timesStat.User + timesStat.Nice) - (c.lastCPUTimes.User + c.lastCPUTimes.Nice)) * toPercent
		system := ((timesStat.System + timesStat.Irq + timesStat.Softirq) - (c.lastCPUTimes.System + c.lastCPUTimes.Irq + c.lastCPUTimes.Softirq)) * toPercent
		used := user + system
		idle := (timesStat.Idle - c.lastCPUTimes.Idle) * toPercent

		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("user"), user, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("system"), system, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("used"), used, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("idle"), idle, metric.UnitPercent, t, nil))
	}
	c.lastCycle = cycle
	c.lastCPUTimes = cpuTimes[0]
	return samples, nil
}

func (c CPUCheck) formatMetric(name string) string {
	format := "system.cpu.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	c := &CPUCheck{
		CheckBase: core.NewCheckBase("cpu"),
	}
	core.RegisterCheck(c.String(), c)
}
