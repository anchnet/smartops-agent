package heartbeat

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/packet"
	log "github.com/cihub/seelog"
	"time"
)

var ticker *time.Ticker

func Run(f *forwarder.Forwarder) {
	go func() {
		for range ticker.C {
			hb := packet.NewPacket(packet.Heartbeat, "ping")
			log.Info(fmt.Sprintf("Heartbeat: %s", hb))
			f.Send(hb)
		}
	}()
}

func init() {
	ticker = time.NewTicker(10 * time.Second)
}
