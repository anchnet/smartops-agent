package metrics

import "github.com/anchnet/smartops-agent/pkg/config"

type MetricSample struct {
	Endpoint string
	Metric   string
	Value    float64
	Tags     map[string]string
}

func newMetricSample(category string, metric string, value float64, tags map[string]string) *MetricSample {
	return &MetricSample{
		Endpoint: config.SmartOps.GetString("endpoint") + "_" + category,
		Metric:   metric,
		Value:    value,
		Tags:     tags,
	}
}

func NewServerMetricSample(metric string, value float64, tags map[string]string) *MetricSample {
	return newMetricSample("server", metric, value, tags)
}

type SenderMetrics struct {
	Metrics []*MetricSample
}

func NewSenderMetrics(metrics []*MetricSample) SenderMetrics {
	return SenderMetrics{
		Metrics: metrics,
	}
}
