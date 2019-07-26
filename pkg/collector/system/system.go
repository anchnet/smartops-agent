package system

import (
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"time"
)

type SystemCheck struct {
}

func (SystemCheck) Run() []metric.MetricSample {
	var samples []metric.MetricSample
	t := time.Now()

	//cpu
	if s, err := runCPUCheck(t); err != nil {
		log.Warn(err)
	} else {
		samples = append(samples, s...)
	}

	//mem
	if s, err := runMemCheck(t); err != nil {
		log.Warn(err)
	} else {
		samples = append(samples, s...)
	}

	//disk
	if s, err := runDiskCheck(t); err != nil {
		log.Warn(err)
	} else {
		samples = append(samples, s...)
	}

	//load
	if s, err := runLoadCheck(t); err != nil {
		log.Warn(err)
	} else {
		samples = append(samples, s...)
	}

	return samples
}
