package packet

import (
	"github.com/anchnet/smartops-agent/pkg/config"
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
	Message string `json:"message"`
}

func NewPacket(typ string, data interface{}) Packet {
	return Packet{Endpoint: config.SmartOps.GetString("endpoint"), Type: typ, Data: data, Time: time.Now()}
}
func NewServerPacket(typ string, data interface{}) Packet {
	return Packet{Endpoint: config.SmartOps.GetString("endpoint") + "_server", Type: typ, Data: data, Time: time.Now()}
}
