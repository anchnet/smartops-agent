package forward

import (
	"bytes"
	"encoding/json"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/metrics"
	//"github.com/anchnet/smartops-agent/pkg/util/log"
	log "github.com/cihub/seelog"
	"golang.org/x/net/websocket"
)

var (
	wsUrl  = config.SmartOps.GetString("ws_site")
	origin = config.SmartOps.GetString("site_ori")
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
	w.wsInstance, err = websocket.Dial(wsUrl, "", origin)
	if err != nil {
		return err
	}
	return nil
}

func (w *Forward) Send(metrics metrics.SenderMetrics) {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(metrics)
	_, _ = w.wsInstance.Write(buffer.Bytes())
	log.Infof("Successfully posted payload")
}
