package gotool

import (
	"net"
	"strings"
	"time"
)

//指定本机某IP 与remote IP连接
func Dial(network string, local string, remote string) (net.Conn, error) {
	network = strings.ToLower(network)
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	local = local + ":0" //端口0,系统会自动分配本机端口
	switch network {
	case "udp":
		addr, err := net.ResolveUDPAddr(network, local)
		if err != nil {
			return nil, err
		}
		dialer.LocalAddr = addr
	case "tcp":
		addr, err := net.ResolveTCPAddr(network, local)
		if err != nil {
			return nil, err
		}
		dialer.LocalAddr = addr
	}
	return dialer.Dial(network, remote)
}
