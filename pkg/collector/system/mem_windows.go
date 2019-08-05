// +build windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/mem"
	"time"
)

const (
	memMetric = "system.mem.%s"
)

func runMemCheck(time time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Errorf("Could not retrieve virtual memory diskStats: %s", err)
		return nil, err
	}
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "total"), float64(v.Total), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "free"), float64(v.Available), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "used"), float64(v.Used), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "used_pct"), v.UsedPercent, metric.UnitPercent, time, nil))

	return samples, nil
}
