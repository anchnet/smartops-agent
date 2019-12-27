package plugin

import (
	"encoding/json"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	log "github.com/cihub/seelog"
	"time"
)

func PluginCollect() []metric.MetricSample {
	var samples []metric.MetricSample
	pluginChecks := core.GetAllPluginsCheck()

	t := time.Now()
	for _, c := range pluginChecks {
		if s, err := c.PluginCollect(t); err != nil {
			_ = log.Warnf("Error while run collect %s, %v", c.PluginName(), err)
		} else {
			samples = append(samples, s...)
		}
	}
	jsonByte, _ := json.Marshal(samples)
	log.Debug(string(jsonByte))
	return samples
}
