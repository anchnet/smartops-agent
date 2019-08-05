// +build !windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"time"
)

const (
	memMetric  = "system.mem.%s"
	swapMetric = "system.swap.%s"
)

var runtimeOS = runtime.GOOS

func runMemCheck(time time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Errorf("Could not retrieve virtual memory diskStats: %s", err)
		return nil, err
	}
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "total"), float64(v.Total), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "free"), float64(v.Free), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "used"), float64(v.Used), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "used_pct"), v.UsedPercent, metric.UnitPercent, time, nil))

	switch runtimeOS {
	case "linux":
		samples = append(samples, linuxMemCheck(v, time)...)
	}
	return samples, nil
}
func linuxMemCheck(v *mem.VirtualMemoryStat, time time.Time) []metric.MetricSample {
	var samples []metric.MetricSample
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "cached"), float64(v.Cached), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "buffered"), float64(v.Buffers), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "shared"), float64(v.Shared), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "slab"), float64(v.Slab), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "page_tables"), float64(v.PageTables), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "commit_limit"), float64(v.CommitLimit), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(memMetric, "committed_as"), float64(v.CommittedAS), metric.UnitByte, time, nil))
	samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(swapMetric, "cached"), float64(v.SwapCached), metric.UnitByte, time, nil))
	return samples
}
