package plugin

import (
	"fmt"
	"github.com/anchnet/smartops-agent/cmd/common"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/anchnet/smartops-agent/pkg/util/file"
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
	currentTime          int64
	connections          int
	accepts              int
	drops                int
	request              int
	reading              int
	writing              int
	waiting              int
	connectionsPerSecond float64 "json:`connectionsPerSecond`"
	acceptsPerSecond     float64 "json:`acceptsPerSecond`"
	dropsPerSecond       float64 "json:`dropsPerSecond`"
	requestPerSecond     float64 "json:`requestsPerSecond`"
	readingPerSecond     float64 "json:`readingPerSecond`"
	writingPerSecond     float64 "json:`writingPerSecond`"
	waitingPerSecond     float64 "json:`waitingPerSecond`"
}

var previousNgxDate *NgxData

func (c *NginxCheck) PluginCollect(t time.Time) ([]metric.MetricSample, error) {
	// check nginx.yaml is exist
	if !file.IsExist(fmt.Sprintf("%s/nginx.yaml", common.DefaultNgxConfPath)) {
		return nil, nil
	}
	ngxUrls := config.Nginx.GetStringSlice("instances.nginx_status_url")
	if ngxUrls == nil {
		return nil, nil
	}
	var samples []metric.MetricSample
	for _, url := range ngxUrls {
		tag := getTag(url)
		data, err := c.calNgxMonitorData(previousNgxDate, url)
		if err != nil {
			return nil, err
		}
		if previousNgxDate == nil {
			// init previousNgxMonitorData
			previousNgxDate = data
			return nil, nil
		}

		samples = append(samples, c.collectNginxMetrics(*data, t, tag)...)
	}

	return samples, nil
}

func (c *NginxCheck) calNgxMonitorData(previous *NgxData, url string) (*NgxData, error) {
	var (
		connRate    float64
		acceptsRate float64
		dropsRate   float64
		requestRate float64
		readingRate float64
		writingRate float64
		waitingRate float64
	)

	data, err := c.getNgxMonitorData(url)
	if err != nil {
		return nil, nil
	}
	if previousNgxDate == nil {
		return data, nil
	}
	connRate = float64(data.connections-previous.connections) / float64(data.currentTime-previous.currentTime)
	if connRate <= 0 {
		connRate = 0
	}
	requestRate = float64(data.request-previous.request) / float64(data.currentTime-previous.currentTime)
	if requestRate <= 0 {
		requestRate = 0
	}
	dropsRate = float64(data.drops-previous.drops) / float64(data.currentTime-previous.currentTime)
	if dropsRate <= 0 {
		dropsRate = 0
	}
	acceptsRate = float64(data.accepts-previous.accepts) / float64(data.currentTime-previous.currentTime)
	if acceptsRate <= 0 {
		acceptsRate = 0
	}
	writingRate = float64(data.writing-previous.writing) / float64(data.currentTime-previous.currentTime)
	if writingRate <= 0 {
		writingRate = 0
	}
	readingRate = float64(data.reading-previous.reading) / float64(data.currentTime-previous.currentTime)
	if readingRate <= 0 {
		readingRate = 0
	}
	waitingRate = float64(data.waiting-previous.waiting) / float64(data.currentTime-previous.currentTime)
	if waitingRate <= 0 {
		waitingRate = 0
	}
	return &NgxData{
		currentTime:          0,
		connections:          0,
		accepts:              0,
		drops:                0,
		request:              0,
		reading:              0,
		writing:              0,
		waiting:              0,
		connectionsPerSecond: connRate,
		acceptsPerSecond:     acceptsRate,
		dropsPerSecond:       dropsRate,
		requestPerSecond:     requestRate,
		readingPerSecond:     readingRate,
		writingPerSecond:     writingRate,
		waitingPerSecond:     waitingRate,
	}, nil
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
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("drops"), float64(ngxData.drops), metric.ReqPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("requests"), float64(ngxData.request), metric.ReqPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("reading"), float64(ngxData.reading), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("writing"), float64(ngxData.writing), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("waiting"), float64(ngxData.waiting), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("connections_per_second"), ngxData.connectionsPerSecond, metric.ConnPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("request_per_second"), ngxData.requestPerSecond, metric.ReqPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("accept_per_second"), ngxData.acceptsPerSecond, metric.ReqPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("waiting_per_second"), ngxData.waitingPerSecond, metric.ReqPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("read_per_second"), ngxData.readingPerSecond, metric.BytePerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("writing_per_second"), ngxData.writingPerSecond, metric.BytePerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("drops_per_second"), ngxData.dropsPerSecond, metric.ConnPerSecond, time, tagMap))

	return samples
}

func (c *NginxCheck) PluginName() string {
	return "nginx"
}

func (c *NginxCheck) getNgxMonitorData(url string) (*NgxData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	connections, accepts, drops, request, read, writing, waiting := c.formatResponse(string(data))
	return &NgxData{
		currentTime:          time.Now().Unix(),
		connections:          connections,
		accepts:              accepts,
		drops:                drops,
		request:              request,
		reading:              read,
		writing:              writing,
		waiting:              waiting,
		connectionsPerSecond: 0,
		acceptsPerSecond:     0,
		dropsPerSecond:       0,
		requestPerSecond:     0,
		readingPerSecond:     0,
		writingPerSecond:     0,
		waitingPerSecond:     0,
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
