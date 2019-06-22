package pagerduty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	var pd PagerDuty
	name := pd.GetName()
	require.Equal(t, "pagerduty", name)
}

func TestGetConfig(t *testing.T) {
	var pd PagerDuty
	conf := pd.GetConfig()
	require.Equal(t, Config{}, conf)
}

func TestInitModule(t *testing.T) {
	var pd PagerDuty
	err := pd.InitModule(Config{})
	require.NoError(t, err, "should not be any error here")
}
