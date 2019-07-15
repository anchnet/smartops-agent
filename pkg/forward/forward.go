package forward

import (
	"bytes"
	"encoding/json"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/metrics"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"time"
)

var (
	wsUrl = config.SmartOps.GetString("ws_site")
)

type Forward struct {
	wsInstance *websocket.Conn
	connected  bool
}

func NewForward() *Forward {
	ws := &Forward{}
	return ws
}

func (w *Forward) Connect() error {
	var err error
	w.wsInstance, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}
	w.connected = true
	return nil
}

func (w *Forward) Send(metrics []metrics.MetricSample) {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(metrics)
	err := w.wsInstance.WriteJSON(metrics)
	if err != nil {
		w.connected = false
		w.reconnect()
		log.Error(err)
	} else {
		log.Infof("Successfully posted payload")
	}
}
func (w *Forward) reconnect() {
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
