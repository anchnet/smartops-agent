// +build !windows

package system

import (
	"fmt"
	"time"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUCheck struct {
	name         string
	lastCycle    float64
	lastCPUTimes cpu.TimesStat
}

func (c *CPUCheck) Name() string {
	return c.name
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
		iowait := (timesStat.Iowait - c.lastCPUTimes.Iowait) * toPercent
		idle := (timesStat.Idle - c.lastCPUTimes.Idle) * toPercent
		steal := (timesStat.Steal - c.lastCPUTimes.Steal) * toPercent
		guest := (timesStat.Guest - c.lastCPUTimes.Guest) * toPercent
		usage := ((timesStat.Total() - timesStat.Guest - timesStat.GuestNice) - (c.lastCPUTimes.Total() - c.lastCPUTimes.Guest - c.lastCPUTimes.GuestNice) - (timesStat.Idle - c.lastCPUTimes.Idle) - (timesStat.Iowait - timesStat.Iowait)) * toPercent

		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("used"), usage, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("user"), user, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("system"), system, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("iowait"), iowait, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("idle"), idle, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("steal"), steal, metric.UnitPercent, t, nil))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("guest"), guest, metric.UnitPercent, t, nil))
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
	core.RegisterCheck(&CPUCheck{
		name: "cpu",
	})
}
