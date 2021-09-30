package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/anchnet/smartops-agent/cmd/common"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/packet"
)

const (
	apiServerInfo = "/getinfo"
)

func GetServerInfo() (*ServerInfoData, error) {
	serverInfo := new(ServerInfo)
	site := config.SmartOps.GetString("site")
	url := fmt.Sprintf("https://%s%s", site, apiServerInfo)

	{ //process https
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			goto http
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			goto http
		}
		defer resp.Body.Close()
		if byts, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			if json.Unmarshal(byts, serverInfo); err != nil {
				return nil, err
			}
		}
		if resp.StatusCode == 200 {
			return &serverInfo.Data, nil
		}
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("response code %v, body: %s", resp.StatusCode, string(body))
	}

http:
	url = fmt.Sprintf("http://%s%s", site, apiServerInfo)
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
	if byts, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		if json.Unmarshal(byts, serverInfo); err != nil {
			return nil, err
		}
	}
	if resp.StatusCode == 200 {
		return &serverInfo.Data, nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return nil, fmt.Errorf("response code %v, body: %s", resp.StatusCode, string(body))
}

func ValidateAPIKey() error {
	url := serverInfoData.ApiKeyValidateEndpoint.URL
	reqBody, err := json.Marshal(packet.NewAPIKeyPacket())
	if err != nil {
		return err
	}
	req, err := http.NewRequest(serverInfoData.ApiKeyValidateEndpoint.Method, url, bytes.NewBuffer(reqBody))
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
	return fmt.Errorf("response code %v, body: %s", resp.StatusCode, string(body))
}

func UpsertPlugins(pluginCategory string, isExist bool) error {
	url := serverInfoData.PluginsUpdate.URL
	request, err := json.Marshal(packet.InitPluginPacket(pluginCategory, isExist))
	if err != nil {
		return err
	}
	req, err := http.NewRequest(serverInfoData.PluginsUpdate.Method, url, bytes.NewBuffer(request))
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
	return fmt.Errorf("response code %v, body: %s", resp.StatusCode, string(body))
}

func GetFilter() (byts []byte, err error) {
	url := serverInfoData.Metric.URL
	req, err := http.NewRequest(serverInfoData.Metric.Method, url, nil)
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
		return nil, errors.New("call get metric status not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ioutil read Error")
	}
	bodyM := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyM)
	if err != nil {
		return nil, errors.New("json Unmarshal Error")
	}
	return json.Marshal(bodyM["data"])
}

func LocalMetric(reqByts []byte) (byts []byte, err error) {
	url := fmt.Sprintf("http://%s%s", common.LocalMetricListen, common.LocalMetricHandle)
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
		return nil, errors.New("call local metric  server status not 200")
	}

	byts, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ioutil read Error")
	}
	return
}

func PhysicalDevice(reqByts []byte) (byts []byte, err error) {
	url := fmt.Sprintf("http://%s%s", common.LocalMetricListen, common.PhysicalDeviceHandle)
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
		return nil, errors.New("call local metric  server status not 20")
	}

	byts, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ioutil read Error")
	}
	return
}

func GetFilerIp() (ipPrefix []string, err error) {
	url := serverInfoData.Filter.URL
	req, err := http.NewRequest(serverInfoData.Filter.Method, url, nil)
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
		return nil, errors.New("call get metric status not 200")
	}
	bodySt := GetFilerIpBody{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ioutil read Error")
	}

	err = json.Unmarshal(body, &bodySt)
	if err != nil {
		return nil, errors.New("json Unmarshal Error")
	}

	return bodySt.Data.SwsIP, nil
}

type GetFilerIpBody struct {
	Data struct {
		SwsIP []string `json:"swsIp"`
	} `json:"data"`
}

type ServerInfo struct {
	Data ServerInfoData `json:"data"`
}

type ServerInfoData struct {
	ApiKeyValidateEndpoint struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"apiKeyValidateEndpoint"`
	AgentHealthCheckEndpoint struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"agentHealthCheckEndpoint"`
	PluginsUpdate struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"pluginsUpdate"`
	Metric struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"metric"`
	Filter struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"filter"`
	Transfer struct {
		URL string `json:"url"`
	} `json:"Transfer"`
	AutoUpdate struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"autoUpdate"`
	PutVersion struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"putVersion"`
}

func GetUpdateInfo(data interface{}) error {
	if serverInfoData.AutoUpdate.URL == "" {
		return errors.New("unable to upgrade currently")
	}

	url := serverInfoData.AutoUpdate.URL
	req, err := http.NewRequest(serverInfoData.AutoUpdate.Method, url, nil)
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
	if resp.StatusCode != 200 {
		return errors.New("call get metric status not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("ioutil read Error")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.New("json Unmarshal Error")
	}
	return nil
}

// func AutoUpdate() error {
// 	type Data struct {
// 		Data struct {
// 			Addr      string `json:"addr"`
// 			Sha256sum string `json:"sha256sum"`
// 			Domain    string `json:"domain"`
// 		} `json:"data"`
// 	}
// 	data := Data{}

// 	for {
// 		time.Sleep(10 * time.Second)
// 		log.Info("start try update agent")
// 		err := GetUpdateInfo(&data)
// 		if err != nil {
// 			log.Error("get update info err: ", err)
// 			continue
// 		}

// 		if data.Data.Addr == "" || data.Data.Sha256sum == "" {
// 			log.Info("no update required")
// 			continue
// 		}

// 		resp, err := http.Get(data.Data.Addr)
// 		if err != nil {
// 			log.Error(err)
// 			continue
// 		}

// 		if resp.StatusCode != 200 {
// 			log.Error("get agent binary code not 200")
// 			continue
// 		}

// 		err = update.Update(resp.Body, update.Options{
// 			Sha256Sum: data.Data.Sha256sum,
// 		})

// 		resp.Body.Close()
// 		if err != nil {
// 			log.Error(err)
// 			continue
// 		}
// 		//restrat agent by daemon
// 		log.Info("agent update success")
// 		if data.Data.Domain != "" {
// 			if err := conf.ChangeConfSite(data.Data.Domain, "./conf/smartops.yaml"); err != nil {
// 				log.Error("change smartops site error", err)
// 			}
// 		}
// 		time.Sleep(1 * time.Second)
// 		os.Exit(-1)
// 	}
// }

func SendAgentVersion() error {
	type Version struct {
		Version  string `json:"version"`
		Endpoint string `json:"endpoint"`
	}
	url := serverInfoData.PutVersion.URL
	reqBody, err := json.Marshal(Version{
		Endpoint: config.SmartOps.GetString("endpoint"),
		Version:  common.Version,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(serverInfoData.PutVersion.Method, url, bytes.NewBuffer(reqBody))
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
	return fmt.Errorf("response code %v, body: %s", resp.StatusCode, string(body))
}
