package collector

import (
	log "github.com/cihub/seelog"
	"gitlab.51idc.com/smartops/smartops-agent/pkg/collector/check"
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
