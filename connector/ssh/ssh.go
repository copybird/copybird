package ssh

import (
	"golang.org/x/crypto/ssh"
	"io"

	"github.com/copybird/copybird/compress"
)

const MODULE_NAME = "connector"

type Ssh struct {
	compress.Output
	reader io.Reader
	writer io.Writer
	config  *Config
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

	sshConfig := &ssh.ClientConfig{
		User: "vcap",
		Auth: []ssh.AuthMethod{
			Agent(),
		},
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
	return nil
}

func (c *Ssh) Close() error {
	return nil
}
