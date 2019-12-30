package plugin

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type NginxCheck struct {
	name string
}

type NgxData struct {
	connections int
	accepts     int
	dropPerSecd int
	request     int
	reading     int
	writing     int
	waiting     int
}

func (c *NginxCheck) PluginCollect(t time.Time) ([]metric.MetricSample, error) {
	ngxUrls := config.Nginx.GetStringSlice("instances.nginx_status_url")
	if ngxUrls == nil {
		return nil, nil
	}
	var samples []metric.MetricSample
	for _, url := range ngxUrls {
		tag := getTag(url)
		data, err := c.getNgxMonitorData(url)
		if err != nil {
			return nil, err
		}
		samples = append(samples, c.collectNginxMetrics(data, t, tag)...)
	}
	return samples, nil
}

func getTag(url string) string {
	re := regexp.MustCompile(`//(.*?):`)
	tag := re.FindAllStringSubmatch(url, -1)
	tmpIp := tag[0][1]
	realIp, err := net.LookupHost(tmpIp)
	if err == nil {
		return realIp[0]
	}
	return ""
}
func (c *NginxCheck) collectNginxMetrics(ngxData NgxData, time time.Time, tag string) []metric.MetricSample {
	tagMap := make(map[string]string, 1)
	tagMap["host"] = tag
	var samples []metric.MetricSample
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("connections"), float64(ngxData.connections), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("accepts"), float64(ngxData.accepts), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("conn_dropped_per_s"), float64(ngxData.dropPerSecd), metric.ReqPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("request_per_s"), float64(ngxData.request), metric.ReqPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("reading"), float64(ngxData.reading), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("writing"), float64(ngxData.writing), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("waiting"), float64(ngxData.waiting), metric.Conn, time, tagMap))
	return samples
}

func (c *NginxCheck) PluginName() string {
	return "nginx"
}

func (c *NginxCheck) getNgxMonitorData(url string) (NgxData, error) {
	resp, err := http.Get(url)
	var ngxData NgxData
	if err != nil {
		return ngxData, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ngxData, err
	}
	connections, accepts, dropPerSecd, request, read, writing, waiting := c.formatResponse(string(data))
	return NgxData{
		connections: connections,
		accepts:     accepts,
		dropPerSecd: dropPerSecd,
		request:     request,
		reading:     read,
		writing:     writing,
		waiting:     waiting,
	}, nil
}

func (c *NginxCheck) formatResponse(statData string) (int, int, int, int, int, int, int) {
	re := regexp.MustCompile(`Active connections:\s+(\d+)`)
	matched := re.FindAllStringSubmatch(statData, -1)
	var connections int
	if len(matched) > 0 {
		tmpConnections, _ := strconv.Atoi(matched[0][1])
		connections = tmpConnections
	}
	re = regexp.MustCompile(`(\d+)\s+(\d+)\s+(\d+)`)
	matched = re.FindAllStringSubmatch(statData, -1)
	accepts, _ := strconv.Atoi(matched[0][1])
	handled, _ := strconv.Atoi(matched[0][2])
	request, _ := strconv.Atoi(matched[0][3])
	dropPerSecd := accepts - handled
	re = regexp.MustCompile(`Reading: (\d+)\s+Writing: (\d+)\s+Waiting: (\d+)`)
	matched = re.FindAllStringSubmatch(statData, -1)
	read, _ := strconv.Atoi(matched[0][1])
	writing, _ := strconv.Atoi(matched[0][2])
	waiting, _ := strconv.Atoi(matched[0][3])
	return connections, accepts, dropPerSecd, request, read, writing, waiting
}

func (c *NginxCheck) formatMetric(metricName string) string {
	format := "nginx.net.%s"
	return fmt.Sprintf(format, metricName)
}

func init() {
	//config.SmartOps.GetString("nginx_port")
	core.RegisterPluginCheck(&NginxCheck{
		name: "nginx",
	})
}
