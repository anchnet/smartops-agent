package system

import (
	"encoding/json"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"time"
)

type systemCheck struct {
	first bool
}

func NewSystemCheck() *systemCheck {
	return &systemCheck{first: true}
}

func (sys *systemCheck) Run() []metric.MetricSample {
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

	//disk io
	if s, err := runIOStatsCheck(t); err != nil {
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

	//net
	if s, err := runNetworkCheck(t); err != nil {
		log.Warn(err)
	} else {
		samples = append(samples, s...)
	}

	//proc
	if s, err := runProcCheck(t); err != nil {
		log.Warn(err)
	} else {
		samples = append(samples, s...)
	}

	jsonByte, _ := json.Marshal(samples)
	log.Debug(string(jsonByte))
	return samples
}
