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
	t := cpuTimes[0]
	cycle := t.Total()
	toPercent := 100 / (cycle - lastCycle)
	user := ((t.User + t.Nice) - (lastCPUTimes.User + lastCPUTimes.Nice)) * toPercent
	system := ((t.System + t.Irq + t.Softirq) - (lastCPUTimes.System + lastCPUTimes.Irq + lastCPUTimes.Softirq)) * toPercent
	iowait := (t.Iowait - lastCPUTimes.Iowait) * toPercent
	idle := (t.Idle - lastCPUTimes.Idle) * toPercent
	steal := (t.Steal - lastCPUTimes.Steal) * toPercent
	guest := (t.Guest - lastCPUTimes.Guest) * toPercent

	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "user"), user, metric.UnitPercent, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "system"), system, metric.UnitPercent, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "iowait"), iowait, metric.UnitPercent, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "idle"), idle, metric.UnitPercent, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "steal"), steal, metric.UnitPercent, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(cpuMetric, "guest"), guest, metric.UnitPercent, time, nil))
	lastCycle = cycle
	return samples, nil
}

// 初始化 CPU Times
func init() {
	cpuTimes, _ := cpu.Times(false)
	t := cpuTimes[0]
	lastCPUTimes = t
	lastCycle = t.Total()
}
