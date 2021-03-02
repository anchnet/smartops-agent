package system

import (
	"encoding/json"
	"time"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/collector/filter"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
)

func Collect() []metric.MetricSample {
	var samples []metric.MetricSample
	checks := core.GetAllChecks()

	t := time.Now()
	for _, c := range checks {
		if s, err := c.Collect(t); err != nil {
			_ = log.Warnf("Error while run collect %s, %v", c.Name(), err)
		} else {
			samples = filterAppend(samples, s...)
		}
	}
	jsonByte, _ := json.Marshal(samples)
	log.Debug(string(jsonByte))
	return samples
}

func filterAppend(sli []metric.MetricSample, elems ...metric.MetricSample) []metric.MetricSample {
	for _, elem := range elems {
		if filter.GetFilter().SubMetric(elem.Metric) {
			sli = append(sli, elem)
		}
	}
	return sli
}
