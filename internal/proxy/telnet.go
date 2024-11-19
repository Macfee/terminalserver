package proxy

import (
	"bytes"
	"fmt"
	"net"
)

type TelnetProxy struct {
	conn     net.Conn
	address  string
	username string
	password string
}

func NewTelnetProxy(host string, port int, username, password string) (*TelnetProxy, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	return &TelnetProxy{
		address:  address,
		username: username,
		password: password,
	}, nil
}

func (p *TelnetProxy) Connect() error {
	conn, err := net.Dial("tcp", p.address)
	if err != nil {
		return err
	}
	p.conn = conn

	// 处理Telnet协议选项协商
	if err := p.handleTelnetNegotiation(); err != nil {
		p.conn.Close()
		return err
	}

	// 处理登录认证
	if err := p.handleAuthentication(); err != nil {
		p.conn.Close()
		return err
	}

	return nil
}

func (p *TelnetProxy) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *TelnetProxy) Read(b []byte) (n int, err error) {
	return p.conn.Read(b)
}

func (p *TelnetProxy) Write(b []byte) (n int, err error) {
	return p.conn.Write(b)
}

// Telnet协议常量
const (
	IAC  = 255 // 解释为命令
	DONT = 254 // 你不要使用选项
	DO   = 253 // 请你使用选项
	WONT = 252 // 我不使用选项
	WILL = 251 // 我将使用选项
)

func (p *TelnetProxy) handleTelnetNegotiation() error {
	buf := make([]byte, 256)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}

		// 处理Telnet命令
		i := 0
		for i < n {
			if buf[i] == IAC && i+2 < n {
				// 对所有选项都回复"不"
				response := []byte{IAC}
				switch buf[i+1] {
				case WILL:
					response = append(response, DONT)
				case DO:
					response = append(response, WONT)
				}
				response = append(response, buf[i+2])

				if _, err := p.conn.Write(response); err != nil {
					return err
				}
				i += 3
			} else {
				i++
			}
		}

		// 检查是否收到登录提示
		if bytes.Contains(buf[:n], []byte("login:")) || bytes.Contains(buf[:n], []byte("Password:")) {
			break
		}
	}

	return nil
}

func (p *TelnetProxy) handleAuthentication() error {
	// 在这里添加登录认证的逻辑
	return nil
}
