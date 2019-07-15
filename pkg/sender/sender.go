package sender

import (
	"github.com/anchnet/smartops-agent/pkg/forward"
	"github.com/anchnet/smartops-agent/pkg/metrics"
)

var (
	senderInstance *checkSender
	checkMetricIn  = make(chan []metrics.MetricSample)
)

type checkSender struct {
	smsOut          chan<- []metrics.MetricSample
	forwardInstance *forward.Forward
}

func newCheckSender(smsOut chan<- []metrics.MetricSample) *checkSender {
	return &checkSender{
		smsOut:          smsOut,
		forwardInstance: forward.NewForward(),
	}
}

func GetSender() *checkSender {
	if senderInstance == nil {
		senderInstance = newCheckSender(checkMetricIn)
	}
	return senderInstance
}

func (s *checkSender) Commit(metrics []metrics.MetricSample) {
	s.smsOut <- metrics
}

func (s *checkSender) Run() {
	for {
		select {
		case senderMetrics := <-checkMetricIn:
			s.forwardInstance.Send(senderMetrics)
		}
	}
}
