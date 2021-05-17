package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/forwarder"
	"github.com/anchnet/smartops-agent/pkg/packet"
)

const (
	apiKeyValidateEndpoint   = "/rundeck/agent/api/validate"
	agentHealthCheckEndpoint = "/agent/health_check"
	pluginsUpdate            = "/rundeck/agent/api/create"
	getMetric                = "/monitor/sws/alert/agent/metric"
	localMetric              = "/localmetric"
)

func ValidateAPIKey() error {
	site := config.SmartOps.GetString("site")
	url := fmt.Sprintf("http://%s%s", site, apiKeyValidateEndpoint)
	//url := fmt.Sprintf("https://%s%s", site, apiKeyValidateEndpoint)
	reqBody, err := json.Marshal(packet.NewAPIKeyPacket())
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	//return fmt.Errorf("unexpected response code: %v", string(body))
	return errors.New(fmt.Sprintf("response code %v, body: %s", resp.StatusCode, string(body)))
}

func UpsertPlugins(pluginCategory string, isExist bool) error {
	site := config.SmartOps.GetString("site")
	url := fmt.Sprintf("http://%s%s", site, pluginsUpdate)
	//url := fmt.Sprintf("https://%s%s", site, pluginsUpdate)
	request, err := json.Marshal(packet.InitPluginPacket(pluginCategory, isExist))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(request))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 {
		return nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return errors.New(fmt.Sprintf("response code %v, body: %s", resp.StatusCode, string(body)))
}

func GetFilter() (byts []byte, err error) {
	site := config.SmartOps.GetString("site")
	url := fmt.Sprintf("http://%s%s", site, getMetric)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("Call get metric status not 200.")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Ioutil read Error")
	}
	bodyM := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyM)
	if err != nil {
		return nil, errors.New("Josn Unmarshal Error")
	}
	return json.Marshal(bodyM["data"])
}

func LocalMetric(reqByts []byte) (byts []byte, err error) {

	url := fmt.Sprintf("http://%s%s", forwarder.LocalMetricListen, forwarder.LocalMetricHandle)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqByts))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("Call local metric  server status not 200.")
	}

	byts, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Ioutil read Error")
	}
	return
}

func PhysicalDevice(reqByts []byte) (byts []byte, err error) {

	url := fmt.Sprintf("http://%s%s", forwarder.LocalMetricListen, forwarder.PhysicalDeviceHandle)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqByts))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("Call local metric  server status not 200.")
	}

	byts, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Ioutil read Error")
	}
	return
}
