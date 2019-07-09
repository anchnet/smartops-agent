package collector

import "gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/check"

type Collector struct {
	runner *Runner
	checks map[string]check.Check
}

func NewCollector() *Collector {
	run := NewRunner()

	c := &Collector{
		runner: run,
		checks: make(map[string]check.Check),
	}

	return c
}

func (c *Collector) RunCheck(check check.Check) {
	c.runner.pending <- check
}
