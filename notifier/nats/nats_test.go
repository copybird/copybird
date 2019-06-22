package nats

import (
	"testing"

	"gotest.tools/assert"
)

func TestNats(t *testing.T) {
	conf := Config{
		ClientID:  "client",
		NATSURL:   "0.0.0.0:4222",
		ClusterID: "cluster",
	}

	nats := Nats{}
	assert.Assert(t, nats.GetConfig() != nil)
	err := nats.InitModule(conf)
	if err != nil {
		panic(err)
	}

	err = nats.Run()
	if err != nil {
		panic(err)
	}
}
