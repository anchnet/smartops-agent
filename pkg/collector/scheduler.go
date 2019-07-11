package collector

import (
	"github.com/anchnet/smartops-agent/pkg/collector/check"
	log "github.com/cihub/seelog"
	"time"
)

func Schedule(check check.Check) {
	ticker := time.NewTicker(check.Interval())
	go func() {
		log.Infof("Scheduling check %s with an interval of %v", check.String(), check.Interval())
		for {
			select {
			case <-ticker.C:
				_ = check.Run()
			}
		}
	}()
}
