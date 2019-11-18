package receiver

import (
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/packet"
	log "github.com/cihub/seelog"
)

func Run(ch chan<- packet.Authorize) {
	for {
		msg, err := forwarder.Receive(ch)
		if err != nil {
			_ = log.Error("Receiving  error, %s", err)
			continue
		}
		log.Infof("Receiving success: %s", msg)
	}
}
