package system

import (
	"fmt"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/core"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/defaults"
)

const checkPrefix = "system."

type SystemCheck struct {
	core.CheckBase
}

func (SystemCheck) Run() error {
	//cpuMetrics =
}

func init() {
	core.RegisterCheck(checkPrefix, &SystemCheck{
		CheckBase: core.NewCheckBase(checkPrefix, defaults.CheckInterval),
	})
}
