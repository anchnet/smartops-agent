package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"time"

	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/disk"
)

type IOStatsCheck struct {
	name  string
	ts    int64
	stats map[string]disk.IOCountersStat
}

func (c *IOStatsCheck) Name() string {
	return c.name
}

func (c *IOStatsCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	ioMap, err := disk.IOCounters()
	if err != nil {
		log.Errorf("system.IOCheck: could not retrieve io diskStats: %s", err)
		return samples, err
	}
	// timestamp
	now := time.Now().Unix()
	delta := float64(now - c.ts)

	if c.ts != 0 {
		for device, ioStats := range ioMap {
			lastIOStats, ok := c.stats[device]
			if !ok {
				log.Debug("New device diskStats (possible hotplug) - full diskStats unavailable this iteration.")
				continue
			}
			if delta == 0 {
				log.Debug("No delta to compute - skipping.")
				continue
			}
			tag := make(map[string]string, 1)
			tag["device"] = device
			rBytes := float64(ioStats.ReadBytes - lastIOStats.ReadBytes)
			wBytes := float64(ioStats.WriteBytes - lastIOStats.WriteBytes)
			rCount := float64(ioStats.ReadCount - lastIOStats.ReadCount)
			wCount := float64(ioStats.WriteCount - lastIOStats.WriteCount)
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("byte.read"), rBytes, metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("byte.write"), wBytes, metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("byte.read.sec"), rBytes/delta, metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("byte.write.sec"), wBytes/delta, metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("read.count"), rCount, "", t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("write.count"), wCount, "", t, tag))
		}

	}
	c.stats = ioMap
	c.ts = now
	return samples, nil
}

func (c IOStatsCheck) formatMetric(name string) string {
	format := "system.disk.io.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	core.RegisterCheck(&IOStatsCheck{
		name: "iostats",
	})
}
