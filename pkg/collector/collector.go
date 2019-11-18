package collector

import (
	"github.com/anchnet/smartops-agent/pkg/collector/system"
	"github.com/anchnet/smartops-agent/pkg/sender"
	log "github.com/cihub/seelog"
	"time"
)

var ticker *time.Ticker

func Collect() {
	first := true
	go func() {
		for range ticker.C {
			samples := system.Collect()
			log.Infof("samples: %d.", len(samples))
			//首次数据不发送
			if first {
				first = false
				continue
			}
			sender.Commit(samples)
		}
	}()
}

func init() {
	ticker = time.NewTicker(10 * time.Second)
}
