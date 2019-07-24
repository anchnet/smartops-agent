package forwarder

import (
	"bytes"
	"encoding/json"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/packet"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"time"
)

var (
	forwardInstance *Forwarder
)

type Forwarder struct {
	wsInstance *websocket.Conn
	connected  bool
	authorized bool
}

func NewForwarder() *Forwarder {
	ws := &Forwarder{}
	return ws
}
func GetForwarder() *Forwarder {
	if forwardInstance == nil {
		forwardInstance = NewForwarder()
	}
	return forwardInstance
}

func (w *Forwarder) Connect() error {
	var err error
	if w.wsInstance != nil && w.connected {
		return nil
	}
	wsUrl := config.SmartOps.GetString("ws_site")
	w.wsInstance, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}
	w.connected = true
	w.Send(packet.NewPacket(packet.Auth, config.SmartOps.GetString("api_key")))
	return nil
}

func (w *Forwarder) auth(msg string, ch chan<- packet.Authorize) {
	var au packet.Authorize
	_ = json.Unmarshal([]byte(msg), &au)
	ch <- au
	w.authorized = true
}
func (w *Forwarder) Stop() error {
	if w.wsInstance != nil {
		if err := w.wsInstance.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (w *Forwarder) Send(packet packet.Packet) {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(packet)
	err := w.wsInstance.WriteJSON(packet)
	if err != nil {
		log.Error(err)
		w.connected = false
		w.reconnect()
	} else {
		log.Infof("Successfully posted payload")
	}
}
func (w *Forwarder) Receive(ch chan<- packet.Authorize) (string, error) {
	_, p, err := w.wsInstance.ReadMessage()
	if err != nil {
		log.Error(err)
		w.connected = false
		w.reconnect()
	}
	msg := string(p)
	if !w.authorized {
		w.auth(msg, ch)
	}
	log.Infof("Receive message: %v", msg)
	return msg, err
}

func (w *Forwarder) reconnect() {
	num := 1
	for w.connected == false {
		if w.wsInstance != nil {
			w.wsInstance.Close()
		}
		log.Infof("Reconnect... #%d", num)
		w.Connect()
		num++
		time.Sleep(10 * time.Second)
	}
	log.Info("Reconnect successful")

}
