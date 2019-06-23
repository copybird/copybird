package ssh

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var gssh GlobalConnectSsh
var config = Config {
	LocalEndpointHost: "127.0.0.1",
	LocalEndpointPort: 8080,
	ServerEndpointHost: "127.0.0.2",
	ServerEndpointPort: 22,
	RemoteEndpointHost: "127.0.0.1",
	RemoteEndpointPort: 8080,
	RemoteUser: "root",
	KeyPath: "",
}

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

// TODO: Это пока не работает, походу есть пробелма в глобальной логике получения конфигов.
//func TestGlobalConnectSsh_GetConfig(t *testing.T) {
//	conf := gssh.GetConfig()
//	assert.Equal(t, conf, &config)
//}

// TODO: Тут надо как-то замокать ключ.
//func TestGlobalConnectSsh_InitModule(t *testing.T) {
//	err := gssh.InitModule(&config)
//	assert.Equal(t, err, nil)
//	assert.Equal(t, gssh.tunnel.Config.User, config.RemoteUser)
//}
