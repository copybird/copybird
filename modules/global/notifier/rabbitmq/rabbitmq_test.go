package rabbitmq

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestRabbitMQ_InvalidConn(t *testing.T) {
	conf := &Config{
		QueueName:      "test.queue",
		PublishKey:     "test.queue",
		RabbitMQURL:    "amqp://guest:guest@localhost:5679/",
		MsgContentType: "text/plain",
		MsgBody:        "Hello",
	}

	rmq := &GlobalNotifierRabbitmq{}
	assert.Assert(t, rmq.GetConfig() != nil)

	err := rmq.InitModule(conf)
	assert.Error(t, err, "dial tcp [::1]:5679: connect: connection refused")
}

func TestRabbitMQ_ValidConn(t *testing.T) {
	conf := &Config{
		QueueName:      "test.queue",
		MsgContentType: "text/plain",
		MsgBody:        "Hello",
		PublishKey:     "test.queue",
		RabbitMQURL:    "amqp://guest:guest@localhost:5672/",
	}

	rmq := &GlobalNotifierRabbitmq{}
	assert.Assert(t, rmq.GetConfig() != nil)

	err := rmq.InitModule(conf)
	if err != nil {
		t.Errorf("TestRabbitMQ: %v", err)
	}

	err = rmq.Run(context.TODO())
	if err != nil {
		t.Errorf("TestRabbitMQ: %v", err)
	}
}
