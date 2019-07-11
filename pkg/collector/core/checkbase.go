package core

import (
	"time"
)

type CheckBase struct {
	name     string
	interval time.Duration
}

func NewCheckBase(name string, interval time.Duration) CheckBase {
	return CheckBase{
		name:     name,
		interval: interval,
	}
}

func (cb *CheckBase) Interval() time.Duration {
	return cb.interval
}
func (c *CheckBase) String() string {
	return c.checkName
}
