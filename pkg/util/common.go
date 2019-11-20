package util

import (
	"net/http"
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
