package metric

import (
	"time"
)

const (
	UnitPercent = "%"
	UnitByte    = "byte"
	Conn        = "conn"
	ReqPerSecd  = "req/s"
)

type MetricSample struct {
	Metric string            `json:"metric"`
	Value  float64           `json:"value"`
	Tag    map[string]string `json:"tag"`
	Unit   string            `json:"unit"`
	Time   time.Time         `json:"time"`
}

func NewServerMetricSample(metric string, value float64, unit string, time time.Time, tags map[string]string) MetricSample {
	return MetricSample{
		Metric: metric,
		Value:  value,
		Tag:    tags,
		Unit:   unit,
		Time:   time,
	}
}
