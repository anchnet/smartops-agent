// +build linux, darwin

package util

import (
	"net"
	"net/http"
	"strings"
)

func CreateHttpTransport() *http.Transport {
	transport := &http.Transport{
		Proxy:                  nil,
		DialContext:            nil,
		Dial:                   nil,
		DialTLS:                nil,
		TLSClientConfig:        nil,
		TLSHandshakeTimeout:    0,
		DisableKeepAlives:      false,
		DisableCompression:     false,
		MaxIdleConns:           0,
		MaxIdleConnsPerHost:    0,
		MaxConnsPerHost:        0,
		IdleConnTimeout:        0,
		ResponseHeaderTimeout:  0,
		ExpectContinueTimeout:  0,
		TLSNextProto:           nil,
		ProxyConnectHeader:     nil,
		MaxResponseHeaderBytes: 0,
		WriteBufferSize:        0,
		ReadBufferSize:         0,
		ForceAttemptHTTP2:      false,
	}
	return transport
}

func LocalIPv4() ([]string, error) {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		pre1 := strings.HasPrefix(i.Name, "eth")
		pre2 := strings.HasPrefix(i.Name, "enp")
		if !pre1 && !pre2 {
			continue
		}
		addr, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addr {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() && ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}
