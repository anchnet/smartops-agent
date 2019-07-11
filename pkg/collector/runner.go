package collector

import (
	log "github.com/cihub/seelog"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/check"
)

type Runner struct {
	pending chan check.Check
	running bool
}

func NewRunner() *Runner {
	r := &Runner{
		pending: make(chan check.Check),
		running: true,
	}
	r.AddWorker()
	return r
}

func (r *Runner) AddWorker() {
	go r.work()
}

func (r *Runner) work() {
	log.Info("Ready to process checks...")
	for c := range r.pending {
		err := c.Run()
		if err != nil {
			log.Errorf("Error running c %s: %s", c, err)
		}
	}
	log.Info("Finished processing checks.")
}
