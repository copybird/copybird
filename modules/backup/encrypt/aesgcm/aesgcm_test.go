package aesgcm

import (
	"bytes"
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestEncryptionAESGCM(t *testing.T) {
	key := "testitnowpleasee"
	enc := &BackupEncryptAesgcm{}

	bufInput := bytes.NewBuffer([]byte("hello world"))
	bufOutput := &bytes.Buffer{}
	assert.Assert(t, enc.GetConfig() != nil)
	assert.NilError(t, enc.InitPipe(bufOutput, bufInput))
	assert.NilError(t, enc.InitModule(&Config{Key: key}))
	assert.NilError(t, enc.Run(context.TODO()))
	assert.NilError(t, enc.Close())
}
