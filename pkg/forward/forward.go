package forward

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/metrics"
	"golang.org/x/net/websocket"
)

var (
	wsurl  = "ws://localhost:8100/monitor"
	origin = "http://localhost:8100/"
)

type Forward struct {
	wsInstance *websocket.Conn
}

func NewForward() *Forward {
	ws := &Forward{}
	if err := ws.connect(); err != nil {
		panic(err)
	}

	return ws
}

func (w *Forward) connect() error {
	var err error
	w.wsInstance, err = websocket.Dial(wsurl, "", origin)
	if err != nil {
		return err
	}
	return nil
}

func (w *Forward) Send(metrics metrics.SenderMetrics) {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(metrics)
	_, _ = w.wsInstance.Write(buffer.Bytes())
	fmt.Println("Send success")
}
