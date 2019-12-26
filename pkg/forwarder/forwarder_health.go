package forwarder

import (
	"github.com/anchnet/smartops-agent/pkg/packet"
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
				//log.Info("Sending heart beat...")
				GetDefaultForwarder().SendMessage(packet.NewHeartbeatPacket())
			}
		}
	}
}
