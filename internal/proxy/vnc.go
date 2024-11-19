package proxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type VNCProxy struct {
	conn     net.Conn
	address  string
	password string
}

func NewVNCProxy(host string, port int, password string) (*VNCProxy, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	return &VNCProxy{
		address:  address,
		password: password,
	}, nil
}

func (p *VNCProxy) Connect() error {
	conn, err := net.Dial("tcp", p.address)
	if err != nil {
		return err
	}
	p.conn = conn

	// 执行VNC握手
	if err := p.handleVNCHandshake(); err != nil {
		p.conn.Close()
		return err
	}

	return nil
}

func (p *VNCProxy) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *VNCProxy) Read(b []byte) (n int, err error) {
	return p.conn.Read(b)
}

func (p *VNCProxy) Write(b []byte) (n int, err error) {
	return p.conn.Write(b)
}

func (p *VNCProxy) handleVNCHandshake() error {
	// 1. 读取服务器版本
	versionMsg := make([]byte, 12)
	if _, err := p.conn.Read(versionMsg); err != nil {
		return fmt.Errorf("读取服务器版本失败: %v", err)
	}

	// 检查版本号是否支持 (RFB 003.008)
	if !bytes.HasPrefix(versionMsg, []byte("RFB 003.008\n")) {
		return fmt.Errorf("不支持的VNC版本: %s", string(versionMsg))
	}

	// 2. 发送客户端版本
	clientVersion := []byte("RFB 003.008\n")
	if _, err := p.conn.Write(clientVersion); err != nil {
		return fmt.Errorf("发送客户端版本失败: %v", err)
	}

	// 3. 读取安全类型
	var numSecTypes uint8
	if err := binary.Read(p.conn, binary.BigEndian, &numSecTypes); err != nil {
		return fmt.Errorf("读取安全类型数量失败: %v", err)
	}

	if numSecTypes == 0 {
		// 读取错误原因
		var reasonLength uint32
		if err := binary.Read(p.conn, binary.BigEndian, &reasonLength); err != nil {
			return fmt.Errorf("读取错误信息长度失败: %v", err)
		}

		reason := make([]byte, reasonLength)
		if _, err := p.conn.Read(reason); err != nil {
			return fmt.Errorf("读取错误信息失败: %v", err)
		}

		return fmt.Errorf("服务器拒绝连接: %s", string(reason))
	}

	// 读取支持的安全类型列表
	secTypes := make([]uint8, numSecTypes)
	if _, err := p.conn.Read(secTypes); err != nil {
		return fmt.Errorf("读取安全类型列表失败: %v", err)
	}

	// 4. 选择安全类型 (VNC Authentication = 2)
	var selectedSecType uint8 = 2
	found := false
	for _, secType := range secTypes {
		if secType == selectedSecType {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("服务器不支持VNC认证")
	}

	// 发送选择的安全类型
	if err := binary.Write(p.conn, binary.BigEndian, selectedSecType); err != nil {
		return fmt.Errorf("发送安全类型选择失败: %v", err)
	}

	// 5. VNC认证
	challenge := make([]byte, 16)
	if _, err := p.conn.Read(challenge); err != nil {
		return fmt.Errorf("读取认证挑战失败: %v", err)
	}

	// 使用DES加密挑战
	response := vncEncryptChallenge(challenge, p.password)
	if _, err := p.conn.Write(response); err != nil {
		return fmt.Errorf("发送认证响应失败: %v", err)
	}

	// 6. 读取认证结果
	var authResult uint32
	if err := binary.Read(p.conn, binary.BigEndian, &authResult); err != nil {
		return fmt.Errorf("读取认证结果失败: %v", err)
	}

	if authResult != 0 {
		return fmt.Errorf("VNC认证失败")
	}

	// 7. 发送ClientInit消息
	if _, err := p.conn.Write([]byte{1}); err != nil { // 1表示共享连接
		return fmt.Errorf("发送ClientInit失败: %v", err)
	}

	// 8. 读取ServerInit消息
	return nil
}
