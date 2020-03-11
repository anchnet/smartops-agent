package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"io/ioutil"
	"net/http"
)

const (
	apiKeyValidateEndpoint   = "/rundeck/run/agent/api/validate"
	agentHealthCheckEndpoint = "/agent/health_check"
)

func ValidateAPIKey() error {
	site := config.SmartOps.GetString("site")
	url := fmt.Sprintf("https://%s%s", site, apiKeyValidateEndpoint)
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
