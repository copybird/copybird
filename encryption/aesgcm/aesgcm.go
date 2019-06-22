package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/copybird/copybird/encryption"
)

const MODULE_NAME = "aesgcm"

type EncryptionAESGCM struct {
	encryption.Encryption
	reader io.Reader
	writer io.Writer
	gcm    cipher.AEAD
	nonce  []byte
}

func (e *EncryptionAESGCM) GetName() string {
	return MODULE_NAME
}

func (e *EncryptionAESGCM) GetConfig() interface{} {
	return &Config{}
}

func (e *EncryptionAESGCM) InitPipe(w io.Writer, r io.Reader, _cfg interface{}) error {
	e.reader = r
	e.writer = w
	return nil
}

func (e *EncryptionAESGCM) InitModule(_cfg interface{}) error {
	cfg := _cfg.(*Config)

	block, err := aes.NewCipher(cfg.Key)
	if err != nil {
		return fmt.Errorf("cipher init err: %s", err)
	}
	e.nonce = make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, e.nonce); err != nil {
		return fmt.Errorf("nonce generate err: %s", err)
	}
	e.gcm, err = cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("aes gcm init err: %s", err)
	}
	return nil
}

func (e *EncryptionAESGCM) Run() error {
	var err error
	buf := make([]byte, 12)
	for {
		_, err = e.reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read err: %s", err)
		}
		ciphertext := e.gcm.Seal(nil, e.nonce, buf, nil)
		_, err = e.writer.Write(ciphertext)
		if err != nil {
			return fmt.Errorf("write err: %s", err)
		}
	}
	return nil
}

func (e *EncryptionAESGCM) Close() error {
	return nil
}
