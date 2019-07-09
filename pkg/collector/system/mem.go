package system

import (
	"fmt"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/core"
	"time"
)

const memCheckName = "memory"

type MemoryCheck struct {
	core.CheckBase
}

func (c *MemoryCheck) Interval() time.Duration {
	return time.Duration(10 * time.Second)
}

func (c *MemoryCheck) Run() error {
	fmt.Println(memCheckName)
	return nil
}

func init() {
	core.RegisterCheck(memCheckName, &MemoryCheck{
		CheckBase: core.NewCheckBase(memCheckName),
	})
}
