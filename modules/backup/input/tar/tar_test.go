package tar

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalInput(t *testing.T) {
	wr := &bytes.Buffer{}
	b := &BackupInputTar{}
	assert.NoError(t, b.InitPipe(wr, nil))
	assert.NoError(t, b.InitModule(&Config{DirectoryPath: "target"}))
	assert.Equal(t, &Config{}, b.GetConfig())
	assert.NoError(t, b.Run(context.TODO()))
	assert.NotNil(t, wr.Bytes())

}
