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
	for check := range r.pending {
		err := check.Run()
		if err != nil {
			fmt.Printf("Error running check %s: %s", check, err)
		}
	}
	fmt.Println("Finished processing checks.")
}
