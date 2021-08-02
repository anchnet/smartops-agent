// +build windows

package system

import (
	"fmt"
	"strconv"

	"time"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/v3/host"
)

type HostCheck struct {
	name string
}

func (c *HostCheck) Name() string {
	return "host"
}

func (c *HostCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample

	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	tags := make(map[string]string)
	tags["uptime"] = strconv.FormatUint(info.Uptime, 10)
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("info"), 0, "", t, tags))

	return samples, nil
}

func (c HostCheck) formatMetric(name string) string {
	format := "system.host.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	core.RegisterCheck(&HostCheck{
		name: "host",
	})
}
