package aesgcm

import (
	"testing"
	"bytes"

	"gotest.tools/assert"
)

func EncryptionAESGCMTest(t *testing.T) {
	key := []byte("test")
	enc := EncryptionAESGCM{}

	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}

	assert.NilError(t, enc.Init(bufOutput, bufInput))
	assert.NilError(t, enc.InitEncryption(key))
	assert.NilError(t, enc.Run())
	assert.NilError(t, enc.Close())
}