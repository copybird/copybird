package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/copybird/copybird/core"
	"io"
)

const GROUP_NAME = "backup"
const TYPE_NAME = "encrypt"
const MODULE_NAME = "aesgcm"

type BackupEncryptAesgcm struct {
	core.Module
	reader io.Reader
	writer io.Writer
	gcm    cipher.AEAD
	nonce  []byte
}

func (m *BackupEncryptAesgcm) GetGroup() string {
	return GROUP_NAME
}

func (m *BackupEncryptAesgcm) GetType() string {
	return TYPE_NAME
}

func (m *BackupEncryptAesgcm) GetName() string {
	return MODULE_NAME
}

func (m *BackupEncryptAesgcm) GetConfig() interface{} {
	return &Config{
		Key: []byte("HELLOWORLD123"),
	}
}

func (m *BackupEncryptAesgcm) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *BackupEncryptAesgcm) InitModule(_cfg interface{}) error {
	cfg := _cfg.(*Config)

	block, err := aes.NewCipher(cfg.Key)
	if err != nil {
		return fmt.Errorf("cipher init err: %s", err)
	}
	m.nonce = make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, m.nonce); err != nil {
		return fmt.Errorf("nonce generate err: %s", err)
	}
	m.gcm, err = cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("aes gcm init err: %s", err)
	}
	return nil
}

func (m *BackupEncryptAesgcm) Run() error {
	var err error
	buf := make([]byte, 12)
	for {
		_, err = m.reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read err: %s", err)
		}
		ciphertext := m.gcm.Seal(nil, m.nonce, buf, nil)
		_, err = m.writer.Write(ciphertext)
		if err != nil {
			return fmt.Errorf("write err: %s", err)
		}
	}
	return nil
}

func (m *BackupEncryptAesgcm) Close() error {
	return nil
}
