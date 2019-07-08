package core

import (
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/defaults"
	"time"
)

type CheckBase struct {
	checkName      string
	latestWarnings []error
	checkInterval  time.Duration
}

func NewCheckBase(name string) CheckBase {
	return CheckBase{
		checkName:     name,
		checkInterval: defaults.DefaultCheckInterval,
	}
}
