package packet

import (
	"time"

	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/anchnet/smartops-agent/pkg/util"
)

type Packet struct {
	Endpoint string      `json:"endpoint"`
	Data     interface{} `json:"data"`
	Type     string      `json:"type"`
	Time     time.Time   `json:"time"`
}

type APIKey struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"authToken"`
	Vendor   string `json:"vendor"`
}

type Plug struct {
	ResourceId     string `json:"resourceId"`
	PluginCategory string `json:"pluginCategory"`
	IsPluginExist  bool   `json:"isPluginExist"`
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

func NewAPIKeyPacket(vendor string) APIKey {
	return APIKey{
		Endpoint: config.SmartOps.GetString("endpoint"),
		Token:    config.SmartOps.GetString("api_key"),
		Vendor:   vendor,
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

func InitPluginPacket(plugCategory string, isExist bool) Plug {
	return Plug{
		ResourceId:     config.SmartOps.GetString("endpoint"),
		PluginCategory: plugCategory,
		IsPluginExist:  isExist,
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

func NewServerCustomPacket(id string, data interface{}) Packet {
	return Packet{
		Endpoint: id,
		Type:     "custom_monitor",
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
