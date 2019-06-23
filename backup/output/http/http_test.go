package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var http Http
	name := GetName()
	require.Equal(t, "http", name)
}

func TestGetConfig(t *testing.T) {
	var http Http
	conf := GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitModule(t *testing.T) {
	var http Http
	err := InitModule(Config{TargetUrl: "https://test.com"})
	require.NoError(t, err, "should not be any error here")
}

func TestRun(t *testing.T) {
	var http Http
	conf := Config{
		TargetUrl: "https://test.com",
	}
	err := InitModule(conf)
	require.NoError(t, err)
	err = Run()
	require.NoError(t, err)
}
