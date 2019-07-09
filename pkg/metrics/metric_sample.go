package metrics

type MetricSample struct {
	Endpoint string
	Metric   string
	Value    float64
	Tags     map[string]string
}

func NewMetricSample(metric string, value float64, tags map[string]string) *MetricSample {
	return &MetricSample{
		Endpoint: "111",
		Metric:   metric,
		Value:    value,
		Tags:     tags,
	}
}

type SenderMetrics struct {
	Metrics []*MetricSample
}

func NewSenderMetrics(metrics []*MetricSample) SenderMetrics {
	return SenderMetrics{
		Metrics: metrics,
	}
}
