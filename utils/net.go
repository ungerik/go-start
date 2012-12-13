package utils

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

const LoopbackIP = "127.0.0.1"

// OwnIP returns the primary IP address of the system or an empty string.
func OwnIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ip := addr.String()
		if ip != LoopbackIP {
			return ip
		}
	}
	return ""
}

// ReadURL reads all data from an URL, including file:// URLs.
func ReadURL(url string) ([]byte, error) {
	if strings.Index(url, "file://") == 0 {
		return ioutil.ReadFile(url[len("file://"):])
	}
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}
