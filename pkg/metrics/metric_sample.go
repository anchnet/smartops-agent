package metrics

import "github.com/anchnet/smartops-agent/pkg/config"

type MetricSample struct {
	Endpoint string
	Metric   string
	Value    float64
	Tags     map[string]string
}

func NewServerMetricSample(metric string, value float64, tags map[string]string) MetricSample {
	return MetricSample{
		Endpoint: config.SmartOps.GetString("endpoint") + "_server",
		Metric:   metric,
		Value:    value,
		Tags:     tags,
	}
}
