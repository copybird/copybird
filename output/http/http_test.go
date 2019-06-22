package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var http Http
	name := http.GetName()
	require.Equal(t, "http", name)
}

func TestGetConfig(t *testing.T) {
	var http Http
	conf := http.GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitModule(t *testing.T) {
	var http Http
	err := http.InitModule(Config{TargetUrl: "test.com"})
	require.NoError(t, err, "should not be any error here")
}
