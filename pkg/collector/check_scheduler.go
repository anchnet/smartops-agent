package collector

import "gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/check"

type CheckScheduler struct {
	collector *Collector
}

func InitCheckScheduler(collector *Collector) *CheckScheduler {
	cs := &CheckScheduler{
		collector: collector,
	}
	return cs
}

func (s *CheckScheduler) Schedule(checks []check.Check) {
	//for _, v := range checks {
	//_,err = s.collector
	//}
}
