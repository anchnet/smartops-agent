package forwarder

import (
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/anchnet/smartops-agent/pkg/util"
	log "github.com/cihub/seelog"
	"time"
)

type forwarderHealth struct {
	f       *defaultForwarder
	stop    chan bool
	stopped chan struct{}
}

func (fh *forwarderHealth) Start() {
	log.Infof("Starting forwarder health check.")
	go fh.healthCheckLoop()
}

func (fh *forwarderHealth) Stop() {
	log.Infof("Stopping forwarder health check.")
	fh.stop <- true
	<-fh.stopped
}

func (fh *forwarderHealth) healthCheckLoop() {
	log.Info("Start health check loop...")

	healthCheckTicker := time.NewTicker(10 * time.Second)

	defer healthCheckTicker.Stop()
	defer close(fh.stopped)

	for {
		select {
		case <-fh.stop:
			return
		case <-healthCheckTicker.C:
			if fh.f.connected {
				// check local ip
				ipsv4, err := util.LocalIPv4()
				if err != nil {
					_ = log.Errorf("get local ipv4 error, %v", err)
					ipsv4 = nil
				}
				GetDefaultForwarder().SendMessage(packet.NewHeartbeatPacket(ipsv4))
			}
		}
	}
}
