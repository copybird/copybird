// +build external

package aesgcm

import (
	"bytes"
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestEncryptionAESGCM(t *testing.T) {
	key := "666f6f6261726b657973616d706c655f"
	enc := &RestoreDecryptAesgcm{}

	// TODO insert valid encrypted string
	bufInput := bytes.NewBuffer([]byte("invalid encrypted string"))
	bufOutput := &bytes.Buffer{}
	assert.Assert(t, enc.GetConfig() != nil)
	assert.NilError(t, enc.InitPipe(bufOutput, bufInput))
	assert.NilError(t, enc.InitModule(&Config{Key: key}))
	assert.NilError(t, enc.Run(context.TODO()))
	assert.NilError(t, enc.Close())
}
