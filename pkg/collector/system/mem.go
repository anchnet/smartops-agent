package system

import (
	"github.com/shirou/gopsutil/mem"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/core"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/metrics"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/sender"
)

const memCheckName = "memory"

var virtualMemory = mem.VirtualMemory

type MemoryCheck struct {
	core.CheckBase
}

func (c *MemoryCheck) Run() error {
	var metricSamples []*metrics.MetricSample
	senderInstance := sender.GetSender()
	v, errVirt := virtualMemory()
	if errVirt == nil {
		metricSamples = append(metricSamples, metrics.NewServerMetricSample("system.mem.total", float64(v.Total), nil))
		metricSamples = append(metricSamples, metrics.NewServerMetricSample("system.mem.free", float64(v.Free), nil))
		metricSamples = append(metricSamples, metrics.NewServerMetricSample("system.mem.used", float64(v.Used), nil))
	}

	senderInstance.Commit(metrics.NewSenderMetrics(metricSamples))
	return nil
}

func init() {
	core.RegisterCheck(memCheckName, &MemoryCheck{
		CheckBase: core.NewCheckBase(memCheckName),
	})
}
