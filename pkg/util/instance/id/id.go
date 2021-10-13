package id

import (
	"io"
	"net/http"
)

type Vendor string

const (
	AliYun     Vendor = "aliyun"
	TencentYun Vendor = "tencentyun"
)

func GetInstanceId(vendor Vendor) (string, error) {
	if vendor == AliYun {
		resp, err := http.Get("http://100.100.100.200/latest/meta-data/instance-id")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
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
		byts, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(byts), nil
	}

	return "", nil
}
