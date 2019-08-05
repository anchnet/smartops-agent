package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/net"
	"time"
)

var (
	netTs    int64
	netStats map[string]net.IOCountersStat
)

const (
	netMetric = "system.net.if.%s"
)

func runNetworkCheck(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	ioByInterface, err := net.IOCounters(true)
	if err != nil {
		log.Errorf("system.NetworkCheck: could not retrieve io diskStats: %s", err)
		return samples, nil
	}
	interfaceMap := make(map[string]net.IOCountersStat)
	for _, s := range ioByInterface {
		interfaceMap[s.Name] = s
	}

	// timestamp
	now := time.Now().Unix()
	delta := float64(now - netTs)

	for _, interfaceIO := range ioByInterface {
		if netTs == 0 {
			continue
		}
		lastNetStats, ok := netStats[interfaceIO.Name]
		if !ok {
			log.Debug("New device netStats (possible hotplug) - full netStats unavailable this iteration.")
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
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(netMetric, "byte,recv"), rBytes, metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(netMetric, "byte.sent"), sBytes, metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(netMetric, "byte.recv.sec"), rBytes/delta, metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(netMetric, "byte.sent.sec"), sBytes/delta, metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(netMetric, "packet.in.count"), rCount, "", t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(netMetric, "packet.out.count"), sCount, "", t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(netMetric, "packet.in.err"), inErrCount, "", t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(netMetric, "packet.out.err"), outErrCount, "", t, tag))

	}
	netStats = interfaceMap
	netTs = now
	return samples, nil
}
