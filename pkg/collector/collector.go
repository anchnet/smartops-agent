package collector

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/system"
	"time"
)

const CHECK_INTERVAL = 10 * time.Second

var ticker *time.Ticker
var check Check

func Collect() {
	go func() {
		for _ = range ticker.C {
			samples, err := check.Run()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("samples:", len(samples))
		}
	}()
}

func init() {
	ticker = time.NewTicker(CHECK_INTERVAL)
	check = system.SystemCheck{}
}
