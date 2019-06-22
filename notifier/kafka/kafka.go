package kafka

import (
	"io"

	"github.com/Shopify/sarama"
)

const (
	// MODULE_NAME is name of a module
	MODULE_NAME = "kafka"
)

// Kafka represends ...
type Kafka struct {
	config *Config
	conn   sarama.SyncProducer
	reader io.Reader
	writer io.Writer
}

// GetName returns name of the module
func (c *Kafka) GetName() string {
	return MODULE_NAME
}

// GetConfig returns module config
func (c *Kafka) GetConfig() interface{} {
	return &Config{}
}

// InitPipe initializes pipe
func (c *Kafka) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

// InitModule initializes module
func (c *Kafka) InitModule(_cfg interface{}) error {
	cfg := _cfg.(*Config)
	c.config = cfg
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = c.config.MaxRetry
	config.Producer.Return.Successes = true
	conn, err := sarama.NewSyncProducer(c.config.BrokerList, config)
	if err != nil {
		return err
	}
	c.conn = conn

	return nil
}

// Run runs module
func (c *Kafka) Run() error {
	msg := &sarama.ProducerMessage{
		Topic: c.config.Topic,
		Value: sarama.StringEncoder(c.config.Message),
	}
	_, _, err := c.conn.SendMessage(msg)
	return err

}

// Close closes compressor
func (c *Kafka) Close() error {
	return nil
}
