package http

var (
// apiKeyValidateEndpoint = "/rundeck/agent/api/validate"
// agentHealthCheckEndpoint = "/agent/health_check"
// pluginsUpdate = "/rundeck/agent/api/create"
// getMetric     = "/monitor/sws/alert/agent/metric"
// getFilter     = "/monitor/sws/alert/agent/filter"
)

var serverInfoData *ServerInfoData

func SetServerInfoData(data *ServerInfoData) {
	serverInfoData = data
}

func GetServerInfoData() *ServerInfoData {
	return serverInfoData
}
