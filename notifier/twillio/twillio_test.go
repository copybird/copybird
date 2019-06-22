package twillio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var twillio Twillio
	name := twillio.GetName()
	require.Equal(t, "twillio", name)
}

func TestGetConfig(t *testing.T) {
	var twillio Twillio
	conf := twillio.GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitModule(t *testing.T) {
	var twillio Twillio
	err := twillio.InitModule(Config{})
	require.NoError(t, err, "should not be any error here")
}
