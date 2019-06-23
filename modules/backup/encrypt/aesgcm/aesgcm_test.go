package aesgcm

import (
	"bytes"
	"testing"

	"gotest.tools/assert"
)

func EncryptionAESGCMTest(t *testing.T) {
	key := []byte("test")
	enc := BackupEncryptAesgcm{}

	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	assert.Assert(t, GetConfig() != nil)
	assert.NilError(t, InitPipe(bufOutput, bufInput))
	assert.NilError(t, InitModule(&Config{Key: key}))
	assert.NilError(t, Run())
	assert.NilError(t, Close())
}
