package packet

import (
	"fmt"
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

type WsResponse struct {
	Type    string `json:"type"`
	Code    int32  `json:"code"`
	Content string `json:"content"`
}

func NewAPIKeyPacket() Packet {
	apiKey := config.SmartOps.GetString("api_key")
	return Packet{Endpoint: config.SmartOps.GetString("endpoint"), Type: APIKey, Data: &AuthToken{Token: apiKey}, Time: time.Now()}
}
func NewHeartbeatPacket() Packet {
	return Packet{Endpoint: config.SmartOps.GetString("endpoint"), Type: Heartbeat, Data: &HeartbeatPack{Message: "ping"}, Time: time.Now()}
}
func NewServerPacket(data interface{}) Packet {
	return Packet{Endpoint: config.SmartOps.GetString("endpoint") + "_server", Type: Monitor, Data: data, Time: time.Now()}
}

func (wr *WsResponse) ToString() string {
	return fmt.Sprintf("type: %s, code: %s, content: %s", wr.Type, wr.Code, wr.Content)
}
