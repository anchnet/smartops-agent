package receiver

import (
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/packet"
	log "github.com/cihub/seelog"
)

var (
	senderInstance *receiver
)

type receiver struct {
	forwardInstance *forwarder.Forwarder
}

func newReceiver() *receiver {
	return &receiver{
		forwardInstance: forwarder.GetForwarder(),
	}
}

func GetReceiver() *receiver {
	if senderInstance == nil {
		senderInstance = newReceiver()
	}
	return senderInstance
}

func (s *receiver) Connect() error {
	if err := s.forwardInstance.Connect(); err != nil {
		return err
	}
	return nil
}

func (s *receiver) Run(ch chan<- packet.Authorize) {
	for {
		str, err := s.forwardInstance.Receive(ch)
		if err != nil {
			log.Error(str)
		}
	}
}
