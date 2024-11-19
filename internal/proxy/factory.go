package proxy

import (
	"fmt"
)

type Proxy interface {
	Connect() error
	Close() error
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

type ProxyFactory struct {
	// 配置信息可以在这里添加
}

func NewProxyFactory() *ProxyFactory {
	return &ProxyFactory{}
}

func (f *ProxyFactory) CreateProxy(protocol string, host string, port int, username, password string) (Proxy, error) {
	switch protocol {
	case "rdp":
		return NewRDPProxy(host, port, username, password)
	case "vnc":
		return NewVNCProxy(host, port, password)
	case "ssh":
		return NewSSHProxy(host, port, username, password, "")
	case "telnet":
		return NewTelnetProxy(host, port, username, password)
	default:
		return nil, fmt.Errorf("不支持的协议类型: %s", protocol)
	}
}
