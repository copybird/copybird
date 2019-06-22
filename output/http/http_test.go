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
	err := http.InitModule(Config{TargetUrl: "https://test.com"})
	require.NoError(t, err, "should not be any error here")
}

func TestRun(t *testing.T) {
	var http Http
	conf := Config{
		TargetUrl: "https://test.com",
	}
	err := http.InitModule(conf)
	require.NoError(t, err)
	err = http.Run()
	require.NoError(t, err)
}
