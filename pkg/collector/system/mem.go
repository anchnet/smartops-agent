package system

import (
	"github.com/anchnet/smartops-agent/pkg/metrics"
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/mem"
)

func runMemCheck() ([]metrics.MetricSample, error) {
	samples := []metrics.MetricSample{}
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Errorf("Could not retrieve virtual memory stats: %s", err)
		return nil, err
	}
	samples = append(samples, metrics.NewServerMetricSample("system.mem.total", float64(v.Total), nil))
	samples = append(samples, metrics.NewServerMetricSample("system.mem.free", float64(v.Free), nil))
	samples = append(samples, metrics.NewServerMetricSample("system.mem.used", float64(v.Used), nil))
	return samples, nil
}
