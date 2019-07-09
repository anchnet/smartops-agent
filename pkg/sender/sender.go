package sender

import (
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/forward"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/metrics"
)

var (
	senderInstance *checkSender
	checkMetricIn  = make(chan metrics.SenderMetrics)
)

type checkSender struct {
	smsOut          chan<- metrics.SenderMetrics
	forwardInstance *forward.Forward
}

func newCheckSender(smsOut chan<- metrics.SenderMetrics) *checkSender {
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

func (s *checkSender) Commit(metrics metrics.SenderMetrics) {
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
