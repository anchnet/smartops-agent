package plugin

import (
	"encoding/json"
	"time"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/collector/filter"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
)

func PluginCollect() []metric.MetricSample {
	var samples []metric.MetricSample
	pluginChecks := core.GetAllPluginsCheck()

	t := time.Now()
	for _, c := range pluginChecks {
		if s, err := c.PluginCollect(t); err != nil {
			_ = log.Warnf("Error while run collect %s, %v", c.PluginName(), err)
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
