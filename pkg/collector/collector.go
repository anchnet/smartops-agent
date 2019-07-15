package collector

import (
	"github.com/anchnet/smartops-agent/pkg/collector/system"
	"github.com/anchnet/smartops-agent/pkg/sender"
	log "github.com/cihub/seelog"
	"time"
)

const CHECK_INTERVAL = 10 * time.Second

var ticker *time.Ticker
var check Check

func Collect() {
	send := sender.GetSender()
	go func() {
		for range ticker.C {
			if samples, err := check.Run(); err != nil {
				log.Warn(err)
			} else {
				send.Commit(samples)
				log.Infof("samples: %d", len(samples))
			}
		}
	}()
}

func init() {
	ticker = time.NewTicker(CHECK_INTERVAL)
	check = system.SystemCheck{}
}
