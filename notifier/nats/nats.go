package nats

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/nats-io/go-nats"
)

const (
	MODULE_NAME = "nats"
)

var errNats = errors.New("NATS not connected")

type Nats struct {
	reqPtr interface{}
	topic  string
	config *Config
	conn   *nats.Conn
	reader io.Reader
	writer io.Writer
}

func (n *Nats) GetName() string {
	return MODULE_NAME
}

func (n *Nats) GetConfig() interface{} {
	return n.config
}

func (c *Nats) InitPipe(w io.Writer, r io.Reader) error {
	c.reader = r
	c.writer = w
	return nil
}

func (c *Nats) InitModule(_cfg interface{}) error {
	cfg := _cfg.(Config)
	c.config = &cfg

	natsConn, err := nats.Connect(
		cfg.NATSURL,
		nats.MaxReconnects(20),
		nats.Timeout(cfg.ConnectTimeout),
	)
	if err != nil {
		return err
	}
	c.conn = natsConn

	return nil
}

func (c *Nats) Run() error {
	b, err := json.Marshal(c.reqPtr)
	if err != nil {
		return err
	}

	c.conn.Publish(c.topic, b)
	return nil
}

// Close closes compressor
func (c *Nats) Close() error {
	return nil
}
