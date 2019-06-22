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
	return n.config
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
	qConf := c.config.QueueConfig
	q, err := c.channel.QueueDeclare(
		qConf.Name,       // name
		qConf.Durable,    // durable
		qConf.AutoDelete, // delete when unused
		qConf.Exclusive,  // exclusive
		qConf.NoWait,     // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}
	c.queue = &q

	// init message to be published
	mConf := c.config.MsgConfig
	p := amqp.Publishing{
		ContentType: mConf.ContentType,
		Body:        mConf.Body,
	}
	c.publ = p

	return nil
}

func (c *RabbitMQ) Run() error {

	pConf := c.config.PublishConfig
	err := c.channel.Publish(
		pConf.Exchange,  // exchange
		pConf.Key,       // routing key
		pConf.Mandatory, // mandatory
		pConf.Immediate, // immediate
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
