package kafka

import (
	"github.com/copybird/copybird/core"
	"io"

	"github.com/Shopify/sarama"
)

const (
	// MODULE_NAME is name of a module
	MODULE_NAME = "kafka"
)

// GlobalNotifieKafka represends ...
type GlobalNotifieKafka struct {
	core.Module
	config *Config
	conn   sarama.SyncProducer
	reader io.Reader
	writer io.Writer
}

// GetName returns name of the module
func (m *GlobalNotifieKafka) GetName() string {
	return MODULE_NAME
}

// GetConfig returns module config
func (m *GlobalNotifieKafka) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initializes pipe
func (m *GlobalNotifieKafka) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

// InitModule initializes module
func (m *GlobalNotifieKafka) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = m.config.MaxRetry
	config.Producer.Return.Successes = true
	conn, err := sarama.NewSyncProducer(m.config.BrokerList, config)
	if err != nil {
		return err
	}
	m.conn = conn

	return nil
}

// Run runs module
func (m *GlobalNotifieKafka) Run() error {
	msg := &sarama.ProducerMessage{
		Topic: m.config.Topic,
		Value: sarama.StringEncoder(m.config.Message),
	}
	_, _, err := m.conn.SendMessage(msg)
	return err

}

// Close closes compressor
func (m *GlobalNotifieKafka) Close() error {
	return nil
}
