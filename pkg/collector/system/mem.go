package system

import (
	"fmt"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/check"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/core"
)

const memCheckName = "memory"

type MemoryCheck struct {
	core.CheckBase
}

func (c *MemoryCheck) Run() error {
	fmt.Println(memCheckName)
	return nil
}

func MemoryFactory() check.Check {
	return &MemoryCheck{
		CheckBase: core.NewCheckBase(memCheckName),
	}
}

func init() {
	core.RegisterCheck(memCheckName, MemoryFactory)
}
