package rabbitmq

import (
	"io"

	"github.com/streadway/amqp"
)

const (
	MODULE_NAME = "rabbitmq"
)

type RabbitMQ struct {
	config  *Config
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
	publ    amqp.Publishing
	reader  io.Reader
	writer  io.Writer
}

func (n *RabbitMQ) GetName() string {
	return MODULE_NAME
}

func (n *RabbitMQ) GetConfig() interface{} {
	return &Config{}
}

func (c *RabbitMQ) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *RabbitMQ) InitModule(_cfg interface{}) error {
	cfg := _cfg.(Config)
	c.config = &cfg

	// connect to server
	conn, err := amqp.Dial(c.config.RabbitMQURL)
	if err != nil {
		return err
	}
	c.conn = conn

	// init channel
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	c.channel = ch

	// declare queue
	q, err := c.channel.QueueDeclare(
		c.config.QueueName,       // name
		c.config.QueueDurable,    // durable
		c.config.QueueAutoDelete, // delete when unused
		c.config.QueueExclusive,  // exclusive
		c.config.QueueNoWait,     // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return err
	}
	c.queue = &q

	// init message to be published
	p := amqp.Publishing{
		ContentType: c.config.MsgContentType,
		Body:        []byte(c.config.MsgBody),
	}
	c.publ = p

	return nil
}

func (c *RabbitMQ) Run() error {
	err := c.channel.Publish(
		c.config.PublishExchange,  // exchange
		c.config.PublishKey,       // routing key
		c.config.PublishMandatory, // mandatory
		c.config.PublishImmediate, // immediate
		c.publ,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *RabbitMQ) Close() error {
	c.channel.Close()
	c.conn.Close()
	return nil
}
