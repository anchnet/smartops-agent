package collector

import (
	"fmt"
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
	fmt.Println("Ready to process checks...")
	for c := range r.pending {
		err := c.Run()
		if err != nil {
			fmt.Printf("Error running c %s: %s", c, err)
		}
	}
	fmt.Println("Finished processing checks.")
}
