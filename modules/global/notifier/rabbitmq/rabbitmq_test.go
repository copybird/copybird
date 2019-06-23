package rabbitmq

import (
	"testing"

	"gotest.tools/assert"
)

func TestRabbitMQ_InvalidConn(t *testing.T) {
	qConf := QueueConfig{
		Name: "test.queue",
	}
	mConf := MsgConfig{
		ContentType: "text/plain",
		Body:        []byte("Test msg"),
	}
	pConf := PublishConfig{
		Key: "test.queue",
	}
	conf := Config{
		QueueConfig:   qConf,
		MsgConfig:     mConf,
		PublishConfig: pConf,
		RabbitMQURL:   "amqp://guest:guest@localhost:5679/",
	}

	rmq := GlobalNotifierRabbitmq{}
	assert.Assert(t, GetConfig() != nil)

	err := InitModule(conf)
	assert.Error(t, err, "dial tcp 127.0.0.1:5679: connect: connection refused")
}

func TestRabbitMQ_ValidConn(t *testing.T) {
	qConf := QueueConfig{
		Name: "test.queue",
	}
	mConf := MsgConfig{
		ContentType: "text/plain",
		Body:        []byte("Test msg"),
	}
	pConf := PublishConfig{
		Key: "test.queue",
	}
	conf := Config{
		QueueConfig:   qConf,
		MsgConfig:     mConf,
		PublishConfig: pConf,
		RabbitMQURL:   "amqp://guest:guest@localhost:5672/",
	}

	rmq := GlobalNotifierRabbitmq{}
	assert.Assert(t, GetConfig() != nil)

	err := InitModule(conf)
	if err != nil {
		t.Errorf("TestRabbitMQ: %v", err)
	}

	err = Run()
	if err != nil {
		t.Errorf("TestRabbitMQ: %v", err)
	}
}
