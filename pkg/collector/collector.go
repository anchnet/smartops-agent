package collector

import (
	"github.com/anchnet/smartops-agent/pkg/collector/plugin"
	"github.com/anchnet/smartops-agent/pkg/collector/system"
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"time"
)

var stopCh = make(chan bool, 1)
var ticker = time.NewTicker(10 * time.Second)
var pluginTicker = time.NewTicker(10 * time.Second)

func Collect() {
	first := true
	for {
		select {
		case <-ticker.C:
			samples := system.Collect()
			//首次数据不发送
			if first {
				first = false
				continue
			}
			if forwarder.GetDefaultForwarder().Connected() {
				//log.Infof("Sending server metric samples: %d", cap(samples))
				forwarder.GetDefaultForwarder().SendMessage(packet.NewServerPacket(samples))
			}
		case <-pluginTicker.C:
			samples := plugin.PluginCollect()
			//首次数据不发送
			if first {
				first = false
				continue
			}
			if forwarder.GetDefaultForwarder().Connected() {
				//log.Infof("Sending server metric samples: %d", cap(samples))
				forwarder.GetDefaultForwarder().SendMessage(packet.NewServerPacket(samples))
			}
		case <-stopCh:
			return
		}
	}
}

func Stop() {
	ticker.Stop()
	stopCh <- true
}

func init() {
}
