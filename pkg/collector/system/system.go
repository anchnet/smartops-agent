package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/collector/defaults"
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
