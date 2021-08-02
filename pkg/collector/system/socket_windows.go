// +build windows

package system

import (
	"fmt"
	"time"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"

	"github.com/shirou/gopsutil/v3/net"
)

type SocketCheck struct {
	name string
}

func (c *SocketCheck) Name() string {
	return "tcp"
}

func (c *SocketCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	samples = append(samples, c.collectPartitionMetrics(t)...)
	return samples, nil
}

func (c SocketCheck) collectPartitionMetrics(time time.Time) []metric.MetricSample {
	var samples []metric.MetricSample
	conArray, err := net.Connections("tcp")
	if err != nil {
		return samples
	}
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("connection_count"), float64(len(conArray)), metric.UnitGe, time, nil))
	return samples
}

func (c SocketCheck) formatMetric(name string) string {
	format := "system.tcp.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	core.RegisterCheck(&SocketCheck{
		name: "tcp",
	})
}
