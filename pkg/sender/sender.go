package sender

import (
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/anchnet/smartops-agent/pkg/packet"
)

var ms = make(chan []metric.MetricSample)

func Commit(metrics []metric.MetricSample) {
	ms <- metrics
}

func Run() {
	for {
		select {
		case senderMetrics := <-ms:
			forwarder.Send(packet.NewServerPacket(packet.Monitor, senderMetrics))
		}
	}
}
