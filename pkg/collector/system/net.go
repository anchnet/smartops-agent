package system

import (
	"fmt"
	"time"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/net"
)

type NetCheck struct {
	name       string
	ts         int64
	stats      map[string]net.IOCountersStat
	interfaces map[string]net.InterfaceStat
}

func (c *NetCheck) Name() string {
	return c.name
}

func (c *NetCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	ioByInterface, err := net.IOCounters(true)
	if err != nil {
		return samples, err
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return samples, err
	}

	interfaceMap := make(map[string]net.IOCountersStat)
	for _, s := range ioByInterface {
		interfaceMap[s.Name] = s
	}

	// timestamp
	now := time.Now().Unix()
	delta := float64(now - c.ts)

	if c.ts != 0 {
		for _, interfaceIO := range ioByInterface {
			if c.exclude(interfaceIO) {
				continue
			}
			lastNetStats, ok := c.stats[interfaceIO.Name]
			if !ok {
				log.Debug("New device stats (possible hotplug) - full stats unavailable this iteration.")
				continue
			}
			if delta == 0 {
				log.Debug("No delta to compute - skipping.")
				continue
			}
			tag := make(map[string]string, 1)
			tag["device"] = interfaceIO.Name

			rBytes := float64(interfaceIO.BytesRecv - lastNetStats.BytesRecv)
			sBytes := float64(interfaceIO.BytesSent - lastNetStats.BytesSent)
			rCount := float64(interfaceIO.PacketsRecv - lastNetStats.PacketsRecv)
			sCount := float64(interfaceIO.PacketsSent - lastNetStats.PacketsSent)
			inErrCount := float64(interfaceIO.Errin - lastNetStats.Errin)
			outErrCount := float64(interfaceIO.Errout - lastNetStats.Errout)
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("byte.recv"), rBytes, metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("byte.sent"), sBytes, metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("byte.recv.sec"), rBytes/delta, metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("byte.sent.sec"), sBytes/delta, metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("packet.in.count"), rCount, "", t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("packet.out.count"), sCount, "", t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("packet.in.err"), inErrCount, "", t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("packet.out.err"), outErrCount, "", t, tag))
		}
	}
	c.stats = interfaceMap
	c.ts = now
	if c.interfaces == nil {
		c.interfaces = make(map[string]net.InterfaceStat, len(interfaces))
		for _, i := range interfaces {
			c.interfaces[i.Name] = i
		}
	}
	return samples, nil
}
func (c *NetCheck) exclude(stat net.IOCountersStat) bool {
	i, ok := c.interfaces[stat.Name]
	if ok {
		if i.Flags != nil {
			for _, v := range i.Flags {
				if v == "loopback" {
					return true
				}
			}
		}
	}
	return false
}

func (c NetCheck) formatMetric(name string) string {
	format := "system.net.if.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	core.RegisterCheck(&NetCheck{
		name: "net",
	})
}
