package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/disk"
	"time"
)

type DiskCheck struct {
	name string
}

func (c *DiskCheck) Name() string {
	return "disk"
}

func (c *DiskCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}
	samples = append(samples, c.collectPartitionMetrics(partitions, t)...)
	return samples, nil
}

func (c DiskCheck) exclude(disk disk.PartitionStat) bool {
	switch disk.Fstype {
	case "devfs",
		"devtmpfs",
		"tmpfs":
		return true
	}
	return false
}

func (c DiskCheck) collectPartitionMetrics(partitions []disk.PartitionStat, time time.Time) []metric.MetricSample {
	var samples []metric.MetricSample

	for _, partition := range partitions {
		if c.exclude(partition) {
			continue
		}
		// Get disk metric here to be able to exclude on total usage
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Warnf("Unable to get disk metric of %s mount point: %s", partition.Mountpoint, err)
			continue
		}

		// Exclude disks with total disk size 0
		if usage.Total == 0 {
			continue
		}

		tag := make(map[string]string, 2)

		tag["filesystem"] = partition.Fstype
		tag["mountpoint"] = partition.Mountpoint

		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("total"), float64(usage.Total), metric.UnitByte, time, tag))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("used"), float64(usage.Used), metric.UnitByte, time, tag))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("free"), float64(usage.Free), metric.UnitByte, time, tag))
		samples = append(samples, metric.NewServerMetricSample(c.formatMetric("used.pct"), usage.UsedPercent, metric.UnitPercent, time, tag))
	}

	return samples
}
func (c DiskCheck) formatMetric(name string) string {
	format := "system.disk.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	core.RegisterCheck(&DiskCheck{
		name: "disk",
	})
}
