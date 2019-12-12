package packet

import (
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/anchnet/smartops-agent/pkg/util"
	"time"
)

type Packet struct {
	Endpoint string      `json:"endpoint"`
	Data     interface{} `json:"data"`
	Type     string      `json:"type"`
	Time     time.Time   `json:"time"`
}

type AuthToken struct {
	Token string `json:"token"`
}

type HeartbeatPack struct {
	Message string   `json:"message"`
	IPsv4   []string `json:"ipsv4"`
}

type TaskResult struct {
	TaskId    string `json:"task_id"`
	Output    string `json:"output"`
	Code      int    `json:"code"`
	Completed bool   `json:"completed"`
}

func NewAPIKeyPacket() Packet {
	apiKey := config.SmartOps.GetString("api_key")
	return Packet{
		Endpoint: config.SmartOps.GetString("endpoint"),
		Type:     "auth",
		Data:     &AuthToken{Token: apiKey},
		Time:     time.Now(),
	}
}
func NewHeartbeatPacket() Packet {
	ipsv4, _ := util.LocalIPv4()
	return Packet{
		Endpoint: config.SmartOps.GetString("endpoint"),
		Type:     "heartbeat",
		Data:     &HeartbeatPack{Message: "ping", IPsv4: ipsv4},
		Time:     time.Now(),
	}
}
func NewServerPacket(data []metric.MetricSample) Packet {
	return Packet{
		Endpoint: config.SmartOps.GetString("endpoint") + "_server",
		Type:     "monitor",
		Data:     data,
		Time:     time.Now(),
	}
}

func NewTaskResultPacket(data TaskResult) Packet {
	return Packet{
		Endpoint: config.SmartOps.GetString("endpoint"),
		Type:     "task",
		Data:     data,
		Time:     time.Now(),
	}
}
