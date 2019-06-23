package ssh

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var gssh GlobalConnectSsh

func TestGlobalConnectSsh_GetName(t *testing.T) {
	name := gssh.GetName()
	require.Equal(t, "ssh", name)
}

func TestGlobalConnectSsh_GetGroup(t *testing.T) {
	group := gssh.GetGroup()
	require.Equal(t, "global", group)
}

func TestGlobalConnectSsh_GetType(t *testing.T) {
	moduleType := gssh.GetType()
	require.Equal(t, "connect", moduleType)
}

func TestGlobalConnectSsh_GetConfig(t *testing.T) {

}

func TestGlobalConnectSsh_InitModule(t *testing.T) {

}