package scp

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"

	"github.com/copybird/copybird/output"
	"golang.org/x/crypto/ssh"
)

const MODULE_NAME = "scp"

type SCP struct {
	output.Output
	reader io.Reader
	writer io.Writer
	config *Config
	sess   ssh.Session
	client *sftp.Client
}

func (scp *SCP) GetName() string {
	return MODULE_NAME
}

func (scp *SCP) GetConfig() interface{} {
	return Config{}
}

func (scp *SCP) InitPipe(w io.Writer, r io.Reader) error {
	scp.reader = r
	scp.writer = w
	return nil
}

func (scp *SCP) InitModule(_config interface{}) error {
	var clientConfig *ssh.ClientConfig
	conf := _config.(Config)
	scp.config = &conf

	// get host public key
	hostKey, err := getHostKey(conf.Addr)
	if err != nil {
		return err
	}

	//TODO maybe also check for nil hostkey

	if conf.PathToKey != "" {
		priv, err := ioutil.ReadFile(conf.PathToKey)
		if err != nil {
			return err
		}

		signer, err := ssh.ParsePrivateKey([]byte(priv))
		if err != nil && err.Error() != "ssh: cannot decode encrypted private keys" {
			return err
		}

		if err.Error() == "ssh: cannot decode encrypted private keys" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(priv, []byte(conf.PrivateKeyPassword))
			if err != nil {
				return err
			}
		}

		clientConfig = &ssh.ClientConfig{
			User: scp.config.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.FixedHostKey(hostKey),
		}

	} else {
		clientConfig = &ssh.ClientConfig{
			User: scp.config.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(scp.config.Password),
			},
			HostKeyCallback: ssh.FixedHostKey(hostKey),
		}
	}

	conn, err := ssh.Dial("tcp", conf.Addr+conf.Port, clientConfig)
	if err != nil {
		return fmt.Errorf("Failed to dial: " + err.Error())
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	scp.client = client
	return nil
}

func (scp *SCP) Run() error {

	// create destination file
	dstFile, err := scp.client.Create(scp.config.FileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// copy bytes from reader to destination file
	_, err = io.Copy(dstFile, scp.reader)
	return err
}

func (scp *SCP) Close() error {
	scp.client.Close()
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
