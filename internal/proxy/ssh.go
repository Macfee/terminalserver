package proxy

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHProxy struct {
	conn     *ssh.Client
	session  *ssh.Session
	address  string
	username string
	password string
	keyFile  string
}

func NewSSHProxy(host string, port int, username, password, keyFile string) (*SSHProxy, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	return &SSHProxy{
		address:  address,
		username: username,
		password: password,
		keyFile:  keyFile,
	}, nil
}

func (p *SSHProxy) Connect() error {
	var authMethods []ssh.AuthMethod

	// 添加密码认证
	if p.password != "" {
		authMethods = append(authMethods, ssh.Password(p.password))
	}

	// 添加密钥认证
	if p.keyFile != "" {
		key, err := loadPrivateKey(p.keyFile)
		if err != nil {
			return fmt.Errorf("加载SSH密钥失败: %v", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(key))
	}

	config := &ssh.ClientConfig{
		User:            p.username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", p.address, config)
	if err != nil {
		return fmt.Errorf("SSH连接失败: %v", err)
	}
	p.conn = client

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	p.session = session

	// 请求伪终端
	if err := session.RequestPty("xterm", 40, 80, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		session.Close()
		client.Close()
		return fmt.Errorf("请求PTY失败: %v", err)
	}

	return nil
}

func (p *SSHProxy) Close() error {
	if p.session != nil {
		p.session.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *SSHProxy) Read(b []byte) (n int, err error) {
	if p.session == nil {
		return 0, fmt.Errorf("SSH会话未建立")
	}
	stdout, err := p.session.StdoutPipe()
	if err != nil {
		return 0, err
	}
	return stdout.Read(b)
}

func (p *SSHProxy) Write(b []byte) (n int, err error) {
	if p.session == nil {
		return 0, fmt.Errorf("SSH会话未建立")
	}
	stdin, err := p.session.StdinPipe()
	if err != nil {
		return 0, err
	}
	return stdin.Write(b)
}

func loadPrivateKey(keyFile string) (ssh.Signer, error) {
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(key)
}
