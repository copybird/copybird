package nats

import (
	"testing"

	"gotest.tools/assert"
)

func TestNats_InvalidConn(t *testing.T) {
	conf := Config{
		NATSURL: "0.0.0.0:4223",
		Topic:   "test.topic",
		Msg:     "Test",
	}

	nats := Nats{}
	assert.Assert(t, nats.GetConfig() != nil)
	err := nats.InitModule(conf)
	assert.Error(t, err, "nats: no servers available for connection")
}

func TestNats_ValidConn(t *testing.T) {
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
