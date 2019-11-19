package heartbeat

import (
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"time"
)

var ticker *time.Ticker

func Run() {
	go func() {
		for range ticker.C {
			hb := packet.NewPacket(packet.Heartbeat, packet.HeartbeatPack{Message: "ping"})
			forwarder.Send(hb)
		}
	}()
}

func init() {
	ticker = time.NewTicker(10 * time.Second)
}
