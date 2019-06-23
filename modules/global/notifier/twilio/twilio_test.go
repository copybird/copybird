package twillio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	n := &GlobalNotifierTwilio{}
	require.Equal(t, "twillio", n.GetName())
}

func TestGetConfig(t *testing.T) {
	n := &GlobalNotifierTwilio{}
	require.Equal(t, &Config{}, n.GetConfig())
}

func TestInitModule(t *testing.T) {
	n := &GlobalNotifierTwilio{}
	require.NoError(t, n.InitModule(&Config{}), "should not be any error here")
}
