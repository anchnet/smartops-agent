package metric

import (
	"github.com/anchnet/smartops-agent/pkg/config"
	"time"
)

const (
	UnitPercent = "%"
	UnitByte    = "byte"
)

type MetricSample struct {
	Endpoint string            `json:"endpoint"`
	Metric   string            `json:"metric"`
	Value    float64           `json:"value"`
	Tag      map[string]string `json:"tag"`
	Unit     string            `json:"unit"`
	Time     time.Time         `json:"time"`
}

func NewServerMetricSample(metric string, value float64, unit string, time time.Time, tags map[string]string) MetricSample {
	return MetricSample{
		Endpoint: config.SmartOps.GetString("endpoint") + "_server",
		Metric:   metric,
		Value:    value,
		Tag:      tags,
		Unit:     unit,
		Time:     time,
	}
}
