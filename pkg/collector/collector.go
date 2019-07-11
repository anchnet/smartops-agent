package collector

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/collector/defaults"
	"time"
)

var checks = core.LoadChecks()

func Collect() {
	ticker := time.NewTicker(defaults.CheckInterval)
	go func() {
		fmt.Println("Scheduling check: ", nil)
		for {
			select {
			case <-ticker.C:
				//check.Run()
			}
		}
	}()
}
