package forwarder

import (
	"bytes"
	"encoding/json"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/packet"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	//	represent the state of an un started Forwarder.
	STOPPED uint32 = iota
	//	represent the state of an started Forwarder.
	STARTED
	RESPONSE_SUCCESS = 0
	RESPONSE_ERROR   = 1
)

const (
	apiKeyValidateEndpoint   = "/agent/validate"
	agentHealthCheckEndpoint = "/agent/health_check"
)

var forwarderInstance *defaultForwarder
var forwarderInit sync.Once

type defaultForwarder struct {
	wsAddr        string
	apiKey        string
	wsConn        *websocket.Conn
	healthChecker *forwarderHealth
	state         uint32
	connected     bool
	m             sync.Mutex
	messageCh     chan packet.Packet
	reconnect     chan bool
	stop          chan bool
	stopped       chan struct{}
	authenticated chan bool
	retryCount    int32
}

func newDefaultForwarder() *defaultForwarder {
	return &defaultForwarder{
		wsAddr: config.SmartOps.GetString("ws_url"),
		apiKey: config.SmartOps.GetString("api_key"),
		state:  STOPPED,
		healthChecker: &forwarderHealth{
			stop:    make(chan bool, 1),
			stopped: make(chan struct{}),
		},
		messageCh:     make(chan packet.Packet),
		reconnect:     make(chan bool, 1),
		stop:          make(chan bool, 1),
		stopped:       make(chan struct{}),
		authenticated: make(chan bool, 1),
	}
}

func GetDefaultForwarder() *defaultForwarder {
	forwarderInit.Do(func() {
		forwarderInstance = newDefaultForwarder()
		forwarderInstance.healthChecker.f = *forwarderInstance
	})
	return forwarderInstance
}

func (f *defaultForwarder) connect() error {
	if f.wsConn != nil && f.state == STARTED {
		err := f.wsConn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err == nil {
			return nil
		}
	}
	f.retryCount++
	log.Infof("Connecting to %v #%d", f.wsAddr, f.retryCount)
	conn, _, err := websocket.DefaultDialer.Dial(f.wsAddr, nil)
	if err != nil {
		f.connected = false
		return err
	}
	conn.EnableWriteCompression(true)
	f.wsConn = conn
	f.connected = true
	return nil
}

// Start initialize and run the forwarder.
func (f *defaultForwarder) Start() error {
	// Lock so we can't stop a forwarder while is starting
	f.m.Lock()
	defer f.m.Unlock()

	if f.state == STARTED {
		return log.Errorf("the forwarder is already started")
	}
	if err := f.connect(); err != nil {
		return log.Errorf("connect to server error: %v", err)
	}
	f.state = STARTED

	log.Info("Sending api key packet and waiting for response...")
	err := f.sendMessage(packet.NewAPIKeyPacket())
	if err != nil {
		return log.Errorf("sending api key packet error: %v", err)
	}
	go f.receiveLoop()
	if <-f.authenticated == false {
		return log.Errorf("api key authenticate error")
	}
	log.Info("API key authenticated success.")
	f.healthChecker.Start()
	go f.sendingLoop()

	return nil
}

func (f *defaultForwarder) Stop() error {
	f.m.Lock()
	defer f.m.Unlock()
	if f.state == STOPPED {
		return log.Errorf("the forwarder is already closed")
	}
	f.state = STOPPED
	f.healthChecker.Stop()
	if err := f.wsConn.Close(); err != nil {
		return log.Errorf("close connection error: %v", err)
	}
	f.stop <- true
	return nil
}

func (f *defaultForwarder) SendMessage(p packet.Packet) {
	f.messageCh <- p
}

func (f *defaultForwarder) sendMessage(p packet.Packet) error {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(p)
	err := f.wsConn.WriteJSON(p)
	if err != nil {
		return err
	}
	log.Infof("Sending message success, type: %v, length: %d.", p.Type, buffer.Len())
	return nil
}

func (f *defaultForwarder) sendingLoop() {
	log.Info("Waiting for message to send...")
	defer close(f.stopped)
	for {
		select {
		case <-f.reconnect:
			if err := f.connect(); err != nil {
				time.Sleep(30 * time.Second)
				f.reconnect <- true
			}
		case pack := <-f.messageCh:
			if f.connected {
				if err := f.sendMessage(pack); err != nil {
					f.connected = false
					f.reconnect <- true
					_ = log.Errorf("sending '%s' message error: %v", pack.Type, err)
				}
			}
		case <-f.stop:
			log.Info("Stopping sending loop.")
			return
		}
	}
}

func (f *defaultForwarder) receiveLoop() {
	log.Info("Waiting for message arriving...")
	for {
		response := new(packet.WsResponse)
		err := f.wsConn.ReadJSON(response)
		if err != nil && f.state == STARTED {
			f.connected = false
			f.reconnect <- true
			_ = log.Errorf("receiving message error: %v", err)
			var checkCount = 0
			for !f.connected {
				checkCount++
				log.Infof("Checking is connection recovery...#d", checkCount)
				time.Sleep(30 * time.Second)
			}
		} else {
			log.Infof("Message received: %s", response)
			if response.Type == "auth" {
				if response.Code == RESPONSE_SUCCESS {
					log.Info("Agent authenticate success.")
					f.authenticated <- true
				} else {
					f.authenticated <- false
					_ = log.Errorf("Agent authenticate error: %s", response.Content)
				}
			}
		}
	}
}
