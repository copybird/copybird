package scp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/copybird/copybird/core"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Module Constants
const GROUP_NAME = "backup"
const TYPE_NAME = "output"
const MODULE_NAME = "scp"

type BackupOutputScp struct {
	core.Module
	reader io.Reader
	writer io.Writer
	config *Config
	sess   ssh.Session
	client *sftp.Client
}

func (m *BackupOutputScp) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *BackupOutputScp) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *BackupOutputScp) GetName() string {
	return MODULE_NAME
}

func (m *BackupOutputScp) GetConfig() interface{} {
	return &Config{}
}

func (m *BackupOutputScp) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *BackupOutputScp) InitModule(_config interface{}) error {
	m.config = _config.(*Config)

	// get host public key
	hostKey, err := getHostKey(m.config.Addr)
	if err != nil {
		return err
	}

	//TODO maybe also check for nil hostkey
	var clientConfig *ssh.ClientConfig

	if m.config.PathToKey != "" {
		priv, err := ioutil.ReadFile(m.config.PathToKey)
		if err != nil {
			return err
		}

		signer, err := ssh.ParsePrivateKey([]byte(priv))
		if err != nil && err.Error() != "ssh: cannot decode encrypted private keys" {
			return err
		}

		if err.Error() == "ssh: cannot decode encrypted private keys" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(priv, []byte(m.config.PrivateKeyPassword))
			if err != nil {
				return err
			}
		}

		clientConfig = &ssh.ClientConfig{
			User: m.config.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.FixedHostKey(hostKey),
		}

	} else {
		clientConfig = &ssh.ClientConfig{
			User: m.config.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(m.config.Password),
			},
			HostKeyCallback: ssh.FixedHostKey(hostKey),
		}
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", m.config.Addr, m.config.Port), clientConfig)
	if err != nil {
		return fmt.Errorf("Failed to dial: " + err.Error())
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (m *BackupOutputScp) Run(ctx context.Context) error {

	// create destination file
	dstFile, err := m.client.Create(m.config.FileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// copy bytes from reader to destination file
	_, err = io.Copy(dstFile, m.reader)
	return err
}

func (m *BackupOutputScp) Close() error {
	m.client.Close()
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
