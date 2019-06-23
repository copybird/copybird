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

const MODULE_NAME = "ssh"
const MODULE_GROUP = "global"
const MODULE_TYPE = "connect"

type GlobalConnectSsh struct {
	reader io.Reader
	writer io.Writer
	config *Config
	tunnel *SSHtunnel
}

func (m *GlobalConnectSsh) GetName() string {
	return MODULE_NAME
}

func (m *GlobalConnectSsh) GetGroup() string {
	return MODULE_GROUP
}

func (m *GlobalConnectSsh) GetType() string {
	return MODULE_TYPE
}

func (m *GlobalConnectSsh) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalConnectSsh) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalConnectSsh) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)

	// Local machine tunnel output
	localEndpoint := &Endpoint{
		Host: m.config.LocalEndpointHost,
		Port: m.config.LocalEndpointPort,
	}

	// External server
	serverEndpoint := &Endpoint{
		Host: m.config.ServerEndpointHost,
		Port: m.config.ServerEndpointPort,
	}

	// External machine tunnel input
	remoteEndpoint := &Endpoint{
		Host: m.config.RemoteEndpointHost,
		Port: m.config.RemoteEndpointPort,
	}

	// get host public key
	hostKey, err := getHostKey(m.config.ServerEndpointHost)
	if err != nil {
		return err
	}

	key, err := ioutil.ReadFile(m.config.KeyPath)
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: m.config.RemoteUser,
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
	m.tunnel = tunnel

	return nil
}

func (m *GlobalConnectSsh) Run() error {
	return m.tunnel.Start()
}

func (m *GlobalConnectSsh) Close() error {
	return m.tunnel.Stop()
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
