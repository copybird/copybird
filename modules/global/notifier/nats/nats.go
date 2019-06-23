package nats

import (
	"errors"
	"io"

	"github.com/nats-io/go-nats"
)

const (
	MODULE_NAME = "nats"
)

var (
	errNats           = errors.New("NATS not connected")
	errNatsEmptyTopic = errors.New("NATS empty topic")
)

type Nats struct {
	config *Config
	conn   *nats.Conn
	reader io.Reader
	writer io.Writer
}

func (n *Nats) GetName() string {
	return MODULE_NAME
}

func (n *Nats) GetConfig() interface{} {
	return &Config{}
}

func (c *Nats) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *Nats) InitModule(_cfg interface{}) error {
	c.config = _cfg.(*Config)

	if c.config.Topic == "" {
		return errNatsEmptyTopic
	}

	natsConn, err := nats.Connect(c.config.NATSURL)
	if err != nil {
		return err
	}
	c.conn = natsConn

	return nil
}

func (c *Nats) Run() error {
	return c.conn.Publish(c.config.Topic, []byte(c.config.Msg))
}

// Close closes compressor
func (c *Nats) Close() error {
	c.conn.Close()
	return nil
}
