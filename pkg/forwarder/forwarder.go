package forwarder

import (
	"bytes"
	"encoding/json"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/packet"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
)

var (
	wsConn      *websocket.Conn
	isConnected bool
	authorized  bool
)

func Connect() error {
	var err error
	if wsConn != nil && isConnected {
		return nil
	}
	wsAddr := config.SmartOps.GetString("ws_site")
	wsConn, _, err = websocket.DefaultDialer.Dial(wsAddr, nil)
	if err != nil {
		return err
	}
	wsConn.EnableWriteCompression(true)
	isConnected = true
	Send(packet.NewPacket(packet.Auth, packet.AuthToken{Token: config.SmartOps.GetString("api_key")}))
	return nil
}

func Close() error {
	if wsConn != nil {
		if err := wsConn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func Send(packet packet.Packet) {
	if wsConn == nil || !isConnected {
		_ = log.Warn("Connection is closed, reconnecting...")
		err := Connect()
		if err != nil {
			_ = log.Errorf("Connecting error, %s.", err)
			return
		}
	}
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(packet)
	if wsConn != nil {
		err := wsConn.WriteJSON(packet)
		if err != nil {
			_ = log.Errorf("Sending message error, %s", err)
			isConnected = false
		} else {
			log.Infof("Sending message success, %d bytes.", len(buffer.Bytes()))
		}
	}
}
func Receive(ch chan<- packet.Authorize) (string, error) {
	if wsConn == nil || !isConnected {
		_ = log.Warn("Connection is closed, reconnecting...")
		err := Connect()
		if err != nil {
			_ = log.Errorf("Connecting error, %s.", err)
			return "", err
		}
	}
	_, p, err := wsConn.ReadMessage()
	if err != nil {
		_ = log.Errorf("Reading message error, %s", err)
		isConnected = false
		return "", err
	}
	msg := string(p)
	if !authorized {
		var au packet.Authorize
		_ = json.Unmarshal([]byte(msg), &au)
		ch <- au
		authorized = true
	}
	return msg, err
}
