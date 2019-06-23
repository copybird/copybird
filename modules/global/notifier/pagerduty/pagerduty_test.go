package pagerduty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	n := &GlobalNotifierPagerDuty{}
	require.Equal(t, "pagerduty", n.GetName())
}

func TestGetConfig(t *testing.T) {
	n := &GlobalNotifierPagerDuty{}
	require.Equal(t, &Config{}, n.GetConfig())
}

func TestInitModule(t *testing.T) {
	n := &GlobalNotifierPagerDuty{}
	err := n.InitModule(Config{})
	require.NoError(t, err, "should not be any error here")
}

func TestRun(t *testing.T) {
	n := &GlobalNotifierPagerDuty{}
	err := n.InitModule(Config{
		AuthToken: "insert auth token here",
		From:      "example@example.com",
	})
	require.NoError(t, err, "should not be any error here")
	err = n.Run()
	require.Error(t, err)

}
