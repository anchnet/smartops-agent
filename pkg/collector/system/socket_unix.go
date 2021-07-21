// +build linux

package system

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
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

func (c SocketCheck) getTCPNum() (total float64) {
	total = 0

	//ipv4
	contents, err := ioutil.ReadFile("/proc/net/sockstat")
	if err == nil {
		lines := strings.Split(string(contents), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if fields[0] == "TCP:" {
					val, err := strconv.ParseFloat(fields[2], 10)
					if err != nil {
						val = 0
					}
					total += val
				}
			}
		}
	}

	//ipv6
	contents, err = ioutil.ReadFile("/proc/net/sockstat6")
	if err == nil {
		lines := strings.Split(string(contents), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if fields[0] == "TCP6:" {
					val, err := strconv.ParseFloat(fields[2], 10)
					if err != nil {
						val = 0
					}
					total += val
				}
			}
		}
	}
	return
}

func (c SocketCheck) collectPartitionMetrics(time time.Time) []metric.MetricSample {
	var samples []metric.MetricSample
	num := c.getTCPNum()
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("connection_count"), num, metric.UnitGe, time, nil))
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
