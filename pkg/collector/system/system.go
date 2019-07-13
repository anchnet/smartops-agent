package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/metrics"
)

type SystemCheck struct {
}

func (SystemCheck) Run() ([]metrics.MetricSample, error) {
	var samples = []metrics.MetricSample{}
	cpus, err := runCPUCheck()
	if err != nil {
		return nil, err
	}
	for _, c := range cpus {
		samples = append(samples, c)
	}
	mems, err := runMemCheck()
	if err != nil {
		return nil, err
	}
	for _, m := range mems {
		samples = append(samples, m)
	}

	fmt.Println(samples)
	return samples, nil
}
