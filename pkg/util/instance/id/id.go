package id

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Vendor string

const (
	AliYun     Vendor = "aliyun"
	TencentYun Vendor = "tencentyun"
	HuaweiYun  Vendor = "huaweiyun"
)

var (
	ErrHttpStatus = errors.New("http code status error")
)

func GetInstanceId(vendor Vendor) (string, error) {
	if vendor == AliYun {
		resp, err := http.Get("http://100.100.100.200/latest/meta-data/instance-id")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode > 299 {
			return "", ErrHttpStatus
		}

		byts, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(byts), nil
	}

	if vendor == TencentYun {
		resp, err := http.Get("http://metadata.tencentyun.com/latest/meta-data/instance-id")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode > 299 {
			return "", ErrHttpStatus
		}
		byts, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(byts), nil
	}

	if vendor == HuaweiYun {
		type Data struct {
			UUID string `uuid`
		}
		bod := Data{}
		resp, err := http.Get("http://169.254.169.254/openstack/latest/meta_data.json")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode > 299 {
			return "", ErrHttpStatus
		}
		byts, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(byts, &bod)
		if err != nil {
			return "", err
		}

		return bod.UUID, nil
	}
	return "", nil
}
