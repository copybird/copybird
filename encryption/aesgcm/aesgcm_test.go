package aesgcm

import (
	"bytes"
	"testing"

	"gotest.tools/assert"
)

func EncryptionAESGCMTest(t *testing.T) {
	key := []byte("test")
	enc := EncryptionAESGCM{}

	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	assert.Assert(t, enc.GetConfig() != nil)
	assert.NilError(t, enc.Init(bufOutput, bufInput))
	assert.NilError(t, enc.InitModule(&Config{Key: key}))
	assert.NilError(t, enc.Run())
	assert.NilError(t, enc.Close())
}
