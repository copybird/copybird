package local

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalInput(t *testing.T) {
	b := &BackupInputLocal{}
	wr := &bytes.Buffer{}
	assert.NoError(t, b.InitPipe(wr, nil))
	assert.NoError(t, b.InitModule(&Config{Filename: "test.file"}))
	assert.Equal(t, &Config{}, b.GetConfig())
	assert.NoError(t, b.Run(context.TODO()))
	assert.Equal(t, "test\n", wr.String())

}
