package packet

type Packet struct {
	Data interface{}
	Type string
}

func NewPacket(typ string, data interface{}) Packet {
	return Packet{Type: typ, Data: data}
}
