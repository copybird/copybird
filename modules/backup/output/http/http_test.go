package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	h := &BackupOutputHttp{}
	require.Equal(t, "http", h.GetName())
}

func TestGetConfig(t *testing.T) {
	h := &BackupOutputHttp{}
	require.Equal(t, &Config{}, h.GetConfig())
}

func TestInitModule(t *testing.T) {
	h := &BackupOutputHttp{}
	err := h.InitModule(&Config{TargetUrl: "https://test.com"})
	require.NoError(t, err, "should not be any error here")
}

func TestRun(t *testing.T) {
	h := &BackupOutputHttp{}
	conf := &Config{
		TargetUrl: "https://test.com",
	}
	err := h.InitModule(conf)
	require.NoError(t, err)
	err = h.Run(context.TODO())
	assert.NoError(t, err)
}
