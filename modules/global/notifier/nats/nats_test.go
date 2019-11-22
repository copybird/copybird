package nats

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestNats_InvalidConn(t *testing.T) {
	conf := &Config{
		NATSURL: "0.0.0.0:4223",
		Topic:   "test.topic",
		Msg:     "Test",
	}

	n := &GlobalNotifierNats{}
	assert.Assert(t, n.GetConfig() != nil)
	err := n.InitModule(conf)
	assert.Error(t, err, "nats: no servers available for connection")
}

func TestNats_ValidConn(t *testing.T) {
	conf := &Config{
		NATSURL: "0.0.0.0:4222",
		Topic:   "test.topic",
		Msg:     "Test",
	}

	n := GlobalNotifierNats{}
	assert.Assert(t, n.GetConfig() != nil)
	err := n.InitModule(conf)
	if err != nil {
		t.Errorf("TestNats: %v", err)
	}

	err = n.Run(context.TODO())
	if err != nil {
		t.Errorf("TestNats: %v", err)
	}
}
