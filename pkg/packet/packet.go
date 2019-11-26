package packet

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/config"
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
	ipsv4, _ := util.LocalIPv4()
	return Packet{Endpoint: config.SmartOps.GetString("endpoint"), Type: Heartbeat, Data: &HeartbeatPack{Message: "ping", IPsv4: ipsv4}, Time: time.Now()}
}
func NewServerPacket(data interface{}) Packet {
	return Packet{Endpoint: config.SmartOps.GetString("endpoint") + "_server", Type: Monitor, Data: data, Time: time.Now()}
}

func (wr *WsResponse) String() string {
	return fmt.Sprintf("type: %s, code: %s, content: %s", wr.Type, fmt.Sprint(wr.Code), wr.Content)
}
