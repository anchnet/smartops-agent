package collector

import (
	"github.com/anchnet/smartops-agent/pkg/collector/system"
	"github.com/anchnet/smartops-agent/pkg/sender"
	log "github.com/cihub/seelog"
	"time"
)

const checkInterval = 10 * time.Second

var ticker *time.Ticker
var check Check

func Collect() {
	send := sender.GetSender()
	go func() {
		for range ticker.C {
			samples := check.Run()
			send.Commit(samples)
			log.Infof("samples: %d", len(samples))
		}
	}()
}

func init() {
	ticker = time.NewTicker(checkInterval)
	check = system.NewSystemCheck()
}
