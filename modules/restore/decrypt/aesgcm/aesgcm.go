package aesgcm

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/copybird/copybird/core"
)

const GROUP_NAME = "restore"
const TYPE_NAME = "decrypt"
const MODULE_NAME = "aesgcm"

type RestoreDecryptAesgcm struct {
	core.Module
	reader io.Reader
	writer io.Writer
	gcm    cipher.AEAD
	nonce  []byte
}

func (m *RestoreDecryptAesgcm) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *RestoreDecryptAesgcm) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *RestoreDecryptAesgcm) GetName() string {
	return MODULE_NAME
}

func (m *RestoreDecryptAesgcm) GetConfig() interface{} {
	return &Config{}
}

func (m *RestoreDecryptAesgcm) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *RestoreDecryptAesgcm) InitModule(_cfg interface{}) error {
	cfg := _cfg.(*Config)

	if cfg.Key == "" {
		return errors.New("need key")
	}
	key, err := hex.DecodeString(cfg.Key)
	if err != nil {
		return fmt.Errorf("cipher key hex decode err: %s", err)
	}
	block, err := aes.NewCipher(key)
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

func (m *RestoreDecryptAesgcm) Run(ctx context.Context) error {
	var err error
	buf := make([]byte, 16)
	bufOut := make([]byte, 16)
	for {
		_, err = m.reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read err: %s", err)
		}
		_, err = m.gcm.Open(bufOut, m.nonce, buf, nil)
		if err != nil {
			return fmt.Errorf("decrypt err: %s", err)
		}
		_, err = m.writer.Write(bufOut)
		if err != nil {
			return fmt.Errorf("write err: %s", err)
		}
	}
	return nil
}

func (m *RestoreDecryptAesgcm) Close() error {
	return nil
}
