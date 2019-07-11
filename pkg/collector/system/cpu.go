package system

import (
	"github.com/anchnet/smartops-agent/pkg/collector/defaults"
	"github.com/anchnet/smartops-agent/pkg/metrics"
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

const cpuCheckPrefix = checkPrefix + "cpu."

type CPUCheck struct {
	SystemCheck
	lastCycle float64
	lastTimes cpu.TimesStat
}

func (c *CPUCheck) Run() ([]metrics.MetricSample, error) {
	var samples = make([]metrics.MetricSample, 6)
	cpuTimes, _ := cpu.Times(false)
	t := cpuTimes[0]
	if c.lastCycle == 0 {
		c.lastTimes = t
		c.lastCycle = t.Total()
		time.Sleep(defaults.CheckInterval)
	}
	cycle := t.Total()
	toPercent := 100 / (cycle - c.lastCycle)
	user := ((t.User + t.Nice) - (c.lastTimes.User + c.lastTimes.Nice)) * toPercent
	system := ((t.System + t.Irq + t.Softirq) - (c.lastTimes.System + c.lastTimes.Irq + c.lastTimes.Softirq)) * toPercent
	iowait := (t.Iowait - c.lastTimes.Iowait) * toPercent
	idle := (t.Idle - c.lastTimes.Idle) * toPercent
	steal := (t.Steal - c.lastTimes.Steal) * toPercent
	guest := (t.Guest - c.lastTimes.Guest) * toPercent

	samples = append(samples, *metrics.NewServerMetricSample(cpuCheckPrefix+"user", user, nil))
	samples = append(samples, *metrics.NewServerMetricSample(cpuCheckPrefix+"system", system, nil))
	samples = append(samples, *metrics.NewServerMetricSample(cpuCheckPrefix+"iowait", iowait, nil))
	samples = append(samples, *metrics.NewServerMetricSample(cpuCheckPrefix+"idle", idle, nil))
	samples = append(samples, *metrics.NewServerMetricSample(cpuCheckPrefix+"steal", steal, nil))
	samples = append(samples, *metrics.NewServerMetricSample(cpuCheckPrefix+"guest", guest, nil))
	c.lastCycle = cycle
	return samples, nil
}
