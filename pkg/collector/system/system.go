package system

import (
	"encoding/json"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"time"
)

func Collect() []metric.MetricSample {
	var samples []metric.MetricSample
	checks := core.GetAllChecks()

	t := time.Now()
	for _, c := range checks {
		log.Infof("Running check %s", c.Name())
		if s, err := c.Collect(t); err != nil {
			log.Warnf("Error while run collect %s, %v", c.Name(), err)
		} else {
			samples = append(samples, s...)
		}
	}
	jsonByte, _ := json.Marshal(samples)
	log.Debug(string(jsonByte))
	return samples
}
