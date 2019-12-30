package plugin

import (
	"database/sql"
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/metric"
	_ "github.com/go-sql-driver/mysql"
	"github.com/json-iterator/go"
	"log"
	"strconv"
	"time"
)

type MysqlCheck struct {
	name     string
	UserName string
	Password string
	Port     int
	Host     string
}

type DatabaseStatus struct {
	Metadata  DatabaseMetadata
	Metrics   DatabaseMetrics
	Variables DatabaseVariables
}

type DatabaseMetadata struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DatabaseMetrics struct {
	CurrentConnections          int `json:"current_connections"`
	ConnectionsPerSecond        int `json:"connections_per_second"`
	AbortedConnectionsPerSecond int `json:"aborted_connections_per_second"`
	QueriesPerSecond            int `json:"queries_per_second"`
	ReadsPerSecond              int `json:"read_per_second"`
	WritesPerSecond             int `json:"write_per_second"`
	Uptime                      int
	connections                 int
	abortedConnections          int
	queries                     int
	reads                       int
	writes                      int
}
type DatabaseVariables struct {
	MaxConnections int `json:"max_connections"`
}

// format connect mysql url
func (c *MysqlCheck) formatMysqlUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema", c.UserName, c.Password, c.Host, c.Port)
}

func (c *MysqlCheck) PluginCollect(t time.Time) ([]metric.MetricSample, error) {
	connParams := config.Mysql.Get("instances")
	connBytes, err := jsoniter.Marshal(connParams)
	var samples []metric.MetricSample
	var params []map[string]map[string]interface{}
	if err != nil {
		fmt.Println(err)
	}
	if err := jsoniter.Unmarshal(connBytes, &params); err != nil {
		fmt.Println(err)
	}
	for _, outerVal := range params {
		innerVal := outerVal["mysql"]
		username := fmt.Sprintf("%v", innerVal["username"])
		password := fmt.Sprintf("%v", innerVal["password"])
		host := fmt.Sprintf("%v", innerVal["host"])
		tempPort, ok := innerVal["port"].(float64)
		if !ok {
			fmt.Println("convert origin port interface{} to float64 failed")
		}
		port := int(tempPort)
		msqCheck := &MysqlCheck{
			name:     "",
			UserName: username,
			Password: password,
			Port:     port,
			Host:     host,
		}
		// collect mysql monitor data
		databaseStatus, err := Status(*msqCheck, new(DatabaseStatus))
		if err != nil {
			fmt.Println(err)
		}
		samples = append(samples, c.collectNginxMetrics(*databaseStatus, t, c.Host)...)
		return samples, nil
	}

	return nil, nil
}

func (c *MysqlCheck) collectNginxMetrics(databaseStatus DatabaseStatus, time time.Time, tag string) []metric.MetricSample {
	var samples []metric.MetricSample
	data := databaseStatus.Metrics
	tagMap := make(map[string]string)
	tagMap["host"] = tag
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("connections"), float64(data.connections), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("connections_per_second"), float64(data.ConnectionsPerSecond), metric.ConnPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("queries"), float64(data.queries), metric.Req, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("queries_per_second"), float64(data.QueriesPerSecond), metric.ReqPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("aborted_connections"), float64(data.abortedConnections), metric.Conn, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("aborted_connections_per_second"), float64(data.AbortedConnectionsPerSecond), metric.ConnPerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("reads"), float64(data.reads), metric.UnitByte, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("reads_per_second"), float64(data.ReadsPerSecond), metric.BytePerSecond, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("writes"), float64(data.writes), metric.UnitByte, time, tagMap))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("writes_per_second"), float64(data.WritesPerSecond), metric.BytePerSecond, time, tagMap))
	//samples = append(samples, metric.NewServerMetricSample(c.formatMetric("uptime"), float64(data.Uptime), metric., time, tagMap))
	return samples
}

func Status(c MysqlCheck, previous *DatabaseStatus) (*DatabaseStatus, error) {
	status := &DatabaseStatus{
		Metadata: DatabaseMetadata{
			Name: c.UserName,
			Host: c.Host,
			Port: c.Port,
		},
		Metrics:   DatabaseMetrics{},
		Variables: DatabaseVariables{},
	}
	// Fetch the metrics
	err := execQuery(c, "metrics", previous, status)
	if err != nil {
		return nil, err
	}
	// Fetch the variables
	err = execQuery(c, "variables", previous, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (c *MysqlCheck) PluginName() string {
	return "mysql"
}

// Execute a query on the given database for looking up metrics/variables
func execQuery(c MysqlCheck, queryType string, previous *DatabaseStatus, status *DatabaseStatus) error {
	var (
		key   string
		value string
		table string
	)

	if queryType == "metrics" {
		table = "GLOBAL_STATUS"
	} else if queryType == "variables" {
		table = "GLOBAL_VARIABLES"
	} else {
		log.Fatal("Unknown queryType")
	}

	// Connect to the database
	conn, err := sql.Open("mysql", c.formatMysqlUrl())
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Query("SET GLOBAL show_compatibility_56 = ON")

	if err != nil {
		return err
	}
	rows, err := conn.Query(fmt.Sprintf("SELECT VARIABLE_NAME AS 'key', VARIABLE_VALUE AS 'value' FROM %s", table))
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			return err
		}
		if queryType == "metrics" {
			err = processMetric(previous, status, key, value)
		} else {
			err = processVariable(status, key, value)
		}
		if err != nil {
			return err
		}
	}
	err = postProcessMetrics(previous, status)
	if err != nil {
		return err
	}
	err = rows.Err()
	return err
}

func processMetric(previous *DatabaseStatus, status *DatabaseStatus, key string, value string) error {
	var (
		err                error
		currentConnections int
		connections        int
		diff               int
		abortedConnections int
		queries            int
		uptime             int
		readWriteValue     int
	)

	switch key {
	// Current connections
	case "THREADS_CONNECTED":
		currentConnections, err = strconv.Atoi(value)
		status.Metrics.CurrentConnections = currentConnections
	// Connections per second
	case "CONNECTIONS":
		connections, err = strconv.Atoi(value)
		if previous == nil || previous.Metrics.connections == 0 {
			status.Metrics.ConnectionsPerSecond = 0
			status.Metrics.connections = connections
		} else {
			diff = connections - previous.Metrics.connections
			if diff > 0 {
				status.Metrics.ConnectionsPerSecond = diff
			} else {
				status.Metrics.ConnectionsPerSecond = 0
			}
			status.Metrics.connections = connections
		}
	case "ABORTED_CONNECTS":
		abortedConnections, err = strconv.Atoi(value)

		if previous == nil || previous.Metrics.abortedConnections == 0 {
			status.Metrics.AbortedConnectionsPerSecond = 0
			status.Metrics.abortedConnections = abortedConnections
		} else {
			diff = abortedConnections - previous.Metrics.abortedConnections
			if diff > 0 {
				status.Metrics.AbortedConnectionsPerSecond = diff
			} else {
				status.Metrics.AbortedConnectionsPerSecond = 0
			}

			status.Metrics.abortedConnections = abortedConnections
		}
	// Queries per second
	case "QUERIES":
		queries, err = strconv.Atoi(value)
		if previous == nil || previous.Metrics.queries == 0 {
			status.Metrics.QueriesPerSecond = 0
			status.Metrics.queries = queries
		} else {
			diff = queries - previous.Metrics.queries
			if diff > 0 {
				status.Metrics.QueriesPerSecond = diff
			} else {
				status.Metrics.QueriesPerSecond = 0
			}

			status.Metrics.queries = queries
		}
	// Read/Writes per second
	case "COM_SELECT", "COM_INSERT_SELECT", "COM_REPLACE_SELECT", "COM_DELETE", "COM_INSERT", "COM_UPDATE", "COM_REPLACE":
		readWriteValue, err = strconv.Atoi(value)
		if key == "COM_SELECT" || key == "COM_INSERT_SELECT" || key == "COM_REPLACE_SELECT" {
			status.Metrics.reads += readWriteValue
			if key == "COM_INSERT_SELECT" || key == "COM_REPLACE_SELECT" {
				status.Metrics.writes += readWriteValue
			}
			// Writes
		} else {
			status.Metrics.writes += readWriteValue
		}
	// Uptime
	case "UPTIME":
		uptime, err = strconv.Atoi(value)
		status.Metrics.Uptime = uptime
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

func processVariable(status *DatabaseStatus, key string, value string) error {
	var (
		err            error
		maxConnections int
	)
	if key == "MAX_CONNECTIONS" {
		maxConnections, err = strconv.Atoi(value)
		status.Variables.MaxConnections = maxConnections
	}
	if err != nil {
		return err
	}
	return nil
}

func postProcessMetrics(previous *DatabaseStatus, status *DatabaseStatus) error {
	var diff int

	// If we don't have a previous value for the total reads
	if previous != nil {
		diff = status.Metrics.reads - previous.Metrics.reads
		if diff > 0 {
			status.Metrics.ReadsPerSecond = diff
		} else {
			status.Metrics.ReadsPerSecond = 0
		}
		diff = status.Metrics.writes - previous.Metrics.writes
		if diff > 0 {
			status.Metrics.WritesPerSecond = diff
		} else {
			status.Metrics.WritesPerSecond = 0
		}
	}

	return nil
}

func (c *MysqlCheck) formatMetric(metricName string) string {
	format := "mysql.%s"
	return fmt.Sprintf(format, metricName)
}

func init() {
	core.RegisterPluginCheck(&MysqlCheck{
		name: "nginx",
	})
}
