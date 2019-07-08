package system

import (
	"fmt"
	"smart-capture/pkg/collector/check"
	"smart-capture/pkg/collector/core"
)

const name = "memory"

type MemoryCheck struct {
	core.CheckBase
}

func (c *MemoryCheck) Run() error {
	fmt.Println(name)
	return nil
}

func MemoryFactory() check.Check {
	return &MemoryCheck{
		CheckBase: core.NewCheckBase(name),
	}
}

func init() {
	core.RegisterCheck(name, MemoryFactory)
}
