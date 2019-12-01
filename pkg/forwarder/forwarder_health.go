package forwarder

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/anchnet/smartops-agent/pkg/util"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"time"
)

type forwarderHealth struct {
	f       *defaultForwarder
	stop    chan bool
	stopped chan struct{}
	apiKey  string
	domain  string
}

func (fh *forwarderHealth) Start() {
	log.Infof("Starting forwarder health check.")
	go fh.healthCheckLoop()
}

func (fh *forwarderHealth) Stop() {
	log.Infof("Stopping forwarder health check.")
	fh.stop <- true
	<-fh.stopped
}

func (fh *forwarderHealth) healthCheckLoop() {
	log.Info("Start health check loop...")

	healthCheckTicker := time.NewTicker(10 * time.Second)

	defer healthCheckTicker.Stop()
	defer close(fh.stopped)

	for {
		select {
		case <-fh.stop:
			return
		case <-healthCheckTicker.C:
			if fh.f.connected {
				//log.Info("Sending heart beat...")
				GetDefaultForwarder().SendMessage(packet.NewHeartbeatPacket())
			}
		}
	}
}

func (fh *forwarderHealth) validateAPIKey() (bool, error) {
	url := fmt.Sprintf("%s%s?api_key=%s", fh.domain, apiKeyValidateEndpoint, fh.apiKey)
	transport := util.CreateHttpTransport()
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, fmt.Errorf("unexpected response code: %v", resp.StatusCode)
}

func (fh *forwarderHealth) heartbeat() {
	url := fmt.Sprintf("%s%s?endpoint=%s", fh.domain, agentHealthCheckEndpoint, config.SmartOps.GetString("endpoint"))
	transport := util.CreateHttpTransport()
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		_ = log.Errorf("heartbeat network error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		_ = log.Errorf("heartbeat error: %v", string(bytes))
	}
}
