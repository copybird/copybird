package ssh

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

const MODULE_NAME = "connect"

type Ssh struct {
	reader io.Reader
	writer io.Writer
	config *Config
	tunnel *SSHtunnel
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

	localEndpoint := &Endpoint{
		Host: c.config.LocalEndpointHost,
		Port: c.config.LocalEndpointPort,
	}

	serverEndpoint := &Endpoint{
		Host: c.config.ServerEndpointHost,
		Port: c.config.ServerEndpointPort,
	}

	remoteEndpoint := &Endpoint{
		Host: c.config.RemoteEndpointHost,
		Port: c.config.RemoteEndpointPort,
	}

	// get host public key
	hostKey, err := getHostKey(c.config.ServerEndpointHost)
	if err != nil {
		return err
	}

	key, err := ioutil.ReadFile(c.config.KeyPath)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: c.config.RemoteUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}
	c.tunnel = tunnel

	return nil
}

func (c *Ssh) Run() error {
	return nil
}

func (c *Ssh) Close() error {
	return nil
}

func getHostKey(host string) (ssh.PublicKey, error) {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, fmt.Errorf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	return hostKey, nil
}
