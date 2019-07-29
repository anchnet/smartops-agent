package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

var lastCycle float64
var lastCPUTimes cpu.TimesStat

const (
	cpuMetric = "system.cpu.%s"
)

func runCPUCheck(time time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	cpuTimes, _ := cpu.Times(false)
	timesStat := cpuTimes[0]
	cycle := timesStat.Total()

	if lastCycle != 0 {
		toPercent := 100 / (cycle - lastCycle)
		user := ((timesStat.User + timesStat.Nice) - (lastCPUTimes.User + lastCPUTimes.Nice)) * toPercent
		system := ((timesStat.System + timesStat.Irq + timesStat.Softirq) - (lastCPUTimes.System + lastCPUTimes.Irq + lastCPUTimes.Softirq)) * toPercent
		iowait := (timesStat.Iowait - lastCPUTimes.Iowait) * toPercent
		idle := (timesStat.Idle - lastCPUTimes.Idle) * toPercent
		steal := (timesStat.Steal - lastCPUTimes.Steal) * toPercent
		guest := (timesStat.Guest - lastCPUTimes.Guest) * toPercent

		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "user"), user, metric.UnitPercent, time, nil))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "system"), system, metric.UnitPercent, time, nil))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "iowait"), iowait, metric.UnitPercent, time, nil))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "idle"), idle, metric.UnitPercent, time, nil))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "steal"), steal, metric.UnitPercent, time, nil))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "guest"), guest, metric.UnitPercent, time, nil))
	}
	lastCycle = cycle
	lastCPUTimes = cpuTimes[0]
	return samples, nil
}
