// +build external

package rabbitmq

import (
	"context"
	"io"

	"github.com/copybird/copybird/core"
	"github.com/streadway/amqp"
)

const (
	GROUP_NAME  = "global"
	TYPE_NAME   = "notifier"
	MODULE_NAME = "rabbitmq"
)

type GlobalNotifierRabbitmq struct {
	core.Module
	config  *Config
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
	publ    amqp.Publishing
	reader  io.Reader
	writer  io.Writer
}

func (m *GlobalNotifierRabbitmq) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierRabbitmq) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *GlobalNotifierRabbitmq) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierRabbitmq) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierRabbitmq) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierRabbitmq) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)

	// connect to server
	conn, err := amqp.Dial(m.config.RabbitMQURL)
	if err != nil {
		return err
	}
	m.conn = conn

	// init channel
	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}
	m.channel = ch

	// declare queue
	q, err := m.channel.QueueDeclare(
		m.config.QueueName,       // name
		m.config.QueueDurable,    // durable
		m.config.QueueAutoDelete, // delete when unused
		m.config.QueueExclusive,  // exclusive
		m.config.QueueNoWait,     // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return err
	}
	m.queue = &q

	// init message to be published
	p := amqp.Publishing{
		ContentType: m.config.MsgContentType,
		Body:        []byte(m.config.MsgBody),
	}
	m.publ = p

	return nil
}

func (m *GlobalNotifierRabbitmq) Run(ctx context.Context) error {
	err := m.channel.Publish(
		m.config.PublishExchange,  // exchange
		m.config.PublishKey,       // routing key
		m.config.PublishMandatory, // mandatory
		m.config.PublishImmediate, // immediate
		m.publ,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *GlobalNotifierRabbitmq) Close() error {
	m.channel.Close()
	m.conn.Close()
	return nil
}
