package packet

import "time"

type Packet struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
	Time time.Time   `json:"time"`
}

func NewPacket(typ string, data interface{}) Packet {
	return Packet{Type: typ, Data: data, Time: time.Now()}
}
