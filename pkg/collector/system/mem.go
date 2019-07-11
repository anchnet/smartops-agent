package system

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"gitlab.51idc.com/smartops/smartops-agent/pkg/metrics"
)

const memCheckPrefix = checkPrefix + "mem."

type MemoryCheck struct {
	SystemCheck
}

func (c *MemoryCheck) Run() ([]metrics.MetricSample, error) {
	var samples = make([]metrics.MetricSample, 3)
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("...")
		return nil, err
	}
	samples = append(samples, *metrics.NewServerMetricSample(memCheckPrefix+"total", float64(v.Total), nil))
	samples = append(samples, *metrics.NewServerMetricSample(memCheckPrefix+"free", float64(v.Free), nil))
	samples = append(samples, *metrics.NewServerMetricSample(memCheckPrefix+"used", float64(v.Used), nil))
	return samples, nil
}
