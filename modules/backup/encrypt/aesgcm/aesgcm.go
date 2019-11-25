package aesgcm

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/copybird/copybird/core"
)

const GROUP_NAME = "backup"
const TYPE_NAME = "encrypt"
const MODULE_NAME = "aesgcm"
const BUF_SIZE = 4096

type BackupEncryptAesgcm struct {
	core.Module
	reader    io.Reader
	writer    io.Writer
	gcm       cipher.AEAD
	bufReader *bufio.Reader
}

func (m *BackupEncryptAesgcm) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *BackupEncryptAesgcm) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *BackupEncryptAesgcm) GetName() string {
	return MODULE_NAME
}

func (m *BackupEncryptAesgcm) GetConfig() interface{} {
	return &Config{}
}

func (m *BackupEncryptAesgcm) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	m.bufReader = bufio.NewReaderSize(m.reader, 4096)
	return nil
}

func (m *BackupEncryptAesgcm) InitModule(_cfg interface{}) error {
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

	m.gcm, err = cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("aes gcm init err: %s", err)
	}
	return nil
}

func (m *BackupEncryptAesgcm) Run(ctx context.Context) error {
	var err error
	var n int

	nonce := make([]byte, 12)

	originaData := make([]byte, BUF_SIZE)
	encryptedData := make([]byte, BUF_SIZE+m.gcm.Overhead())

	for {
		n, err = m.reader.Read(originaData)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read err: %s", err)
		}

		if err = binary.Write(m.writer, binary.LittleEndian, int32(n)); err != nil {
			return err
		}

		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return fmt.Errorf("nonce generate err: %s", err)
		}

		if _, err = m.writer.Write(nonce); err != nil {
			return fmt.Errorf("nonce write err: %s", err)
		}

		m.gcm.Seal(encryptedData, nonce, encryptedData, nil)
		_, err = m.writer.Write(encryptedData)
		if err != nil {
			return fmt.Errorf("write err: %s", err)
		}
	}
	return nil
}

func (m *BackupEncryptAesgcm) Close() error {
	return nil
}
