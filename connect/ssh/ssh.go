package ssh

import (
	"io"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

const MODULE_NAME = "connect"

type Ssh struct {
	reader io.Reader
	writer io.Writer
	config *Config
}

func (c *Ssh) GetName() string {
	return MODULE_NAME
}

func (c *Ssh) GetConfig() interface{} {
	return &Config{}
}

func (c *Ssh) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *Ssh) InitModule(_cfg interface{}) error {
	c.config = _cfg.(*Config)
	return nil
}

func (c *Ssh) Run() error {
	localEndpoint := &Endpoint{
		Host: c.config.LocalEndpoint.Host,
		Port: c.config.LocalEndpoint.Port,
	}

	serverEndpoint := &Endpoint{
		Host: c.config.ServerEndpoint.Host,
		Port: c.config.ServerEndpoint.Port,
	}

	remoteEndpoint := &Endpoint{
		Host: c.config.RemoteEndpoint.Host,
		Port: c.config.RemoteEndpoint.Port,
	}

	key, err := ioutil.ReadFile(c.config.KeyPath)
	if err != nil {
		panic(err)
	}


	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}

	sshConfig := &ssh.ClientConfig{
		User: c.config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	err := tunnel.Start()
	if err != nil {
		return err
	}
	return nil
}

func (c *Ssh) Close() error {
	return nil
}
