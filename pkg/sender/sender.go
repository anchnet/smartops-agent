package sender

import (
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/metrics"
	"github.com/anchnet/smartops-agent/pkg/packet"
)

var (
	senderInstance *sender
	checkMetricIn  = make(chan []metrics.MetricSample)
)

type sender struct {
	smsOut          chan<- []metrics.MetricSample
	forwardInstance *forwarder.Forwarder
}

func newSender(smsOut chan<- []metrics.MetricSample) *sender {
	return &sender{
		smsOut:          smsOut,
		forwardInstance: forwarder.GetForwarder(),
	}
}

func GetSender() *sender {
	if senderInstance == nil {
		senderInstance = newSender(checkMetricIn)
	}
	return senderInstance
}

func (s *sender) Commit(metrics []metrics.MetricSample) {
	s.smsOut <- metrics
}

func (s *sender) Connect() error {
	if err := s.forwardInstance.Connect(); err != nil {
		return err
	}
	return nil
}

func (s *sender) Run() {
	for {
		select {
		case senderMetrics := <-checkMetricIn:
			s.forwardInstance.Send(packet.NewPacket(packet.MonitorData, senderMetrics))
		}
	}
}
