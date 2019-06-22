package nats

import (
	"testing"

	"gotest.tools/assert"
)

func TestNats(t *testing.T) {
	conf := Config{
		NATSURL: "0.0.0.0:4222",
		Topic:   "test.topic",
		Msg:     "Test",
	}

	nats := Nats{}
	assert.Assert(t, nats.GetConfig() != nil)
	err := nats.InitModule(conf)
	if err != nil {
		t.Errorf("TestNats: %v", err)
	}

	err = nats.Run()
	if err != nil {
		t.Errorf("TestNats: %v", err)
	}
}
