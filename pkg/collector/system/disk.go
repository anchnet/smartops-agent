package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/disk"
	"time"
)

const (
	diskMetric = "system.disk.%s"
)

func runDiskCheck(time time.Time) ([]metric.MetricSample, error) {

	var samples []metric.MetricSample
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}
	samples = append(samples, collectPartitionMetrics(partitions, time)...)
	return samples, nil
}
func excludeDisk(disk disk.PartitionStat) bool {
	if disk.Fstype == "devfs" {
		return true
	}
	return false
}

func collectPartitionMetrics(partitions []disk.PartitionStat, time time.Time) []metric.MetricSample {
	var samples []metric.MetricSample

	for _, partition := range partitions {
		if excludeDisk(partition) {
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

		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskMetric, "total"), float64(usage.Total), metric.UnitByte, time, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskMetric, "used"), float64(usage.Used), metric.UnitByte, time, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskMetric, "free"), float64(usage.Free), metric.UnitByte, time, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskMetric, "used.pct"), usage.UsedPercent, metric.UnitPercent, time, tag))
	}

	return samples
}
