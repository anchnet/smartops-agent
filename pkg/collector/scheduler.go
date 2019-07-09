package collector

import (
	"fmt"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/check"
	"time"
)

type Scheduler struct {
}

func (s *Scheduler) Schedule(check check.Check) error {
	ticker := time.NewTicker(check.Interval())
	go func() {
		fmt.Println("Scheduling check: ", check)
		for {
			select {
			case <-ticker.C:
				check.Run()
			}
		}
	}()
	return nil
}
