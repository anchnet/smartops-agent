package system

import (
	"github.com/anchnet/smartops-agent/pkg/metrics"
	"github.com/shirou/gopsutil/cpu"
)

var lastCycle float64
var lastCPUTimes cpu.TimesStat

func runCPUCheck() ([]metrics.MetricSample, error) {
	samples := []metrics.MetricSample{}
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

	samples = append(samples, metrics.NewServerMetricSample("system.cpu.user", user, nil))
	samples = append(samples, metrics.NewServerMetricSample("system.cpu.system", system, nil))
	samples = append(samples, metrics.NewServerMetricSample("system.cpu.iowait", iowait, nil))
	samples = append(samples, metrics.NewServerMetricSample("system.cpu.idle", idle, nil))
	samples = append(samples, metrics.NewServerMetricSample("system.cpu.steal", steal, nil))
	samples = append(samples, metrics.NewServerMetricSample("system.cpu.guest", guest, nil))
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
