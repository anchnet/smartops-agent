package collector

import "time"

type Check interface {
	Run() error
	Stop()
	Configure()
	Interval() time.Duration
}
