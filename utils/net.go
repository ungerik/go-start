package utils

import (
	"net"
)

const LoopbackIP = "127.0.0.1"

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

