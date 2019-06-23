package nats

import (
	"errors"
	"github.com/copybird/copybird/core"
	"io"

	"github.com/nats-io/go-nats"
)

const (
	GROUP_NAME = "global"
	TYPE_NAME = "notifier"
	MODULE_NAME = "nats"
)

var (
	errNats           = errors.New("NATS not connected")
	errNatsEmptyTopic = errors.New("NATS empty topic")
)

type GlobalNotifierNats struct {
	core.Module
	config *Config
	conn   *nats.Conn
	reader io.Reader
	writer io.Writer
}

func (m *GlobalNotifierNats) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierNats) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *GlobalNotifierNats) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierNats) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierNats) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *GlobalNotifierNats) InitModule(_cfg interface{}) error {
	m.config = _cfg.(*Config)

	if m.config.Topic == "" {
		return errNatsEmptyTopic
	}

	natsConn, err := nats.Connect(m.config.NATSURL)
	if err != nil {
		return err
	}
	m.conn = natsConn

	return nil
}

func (m *GlobalNotifierNats) Run() error {
	return m.conn.Publish(m.config.Topic, []byte(m.config.Msg))
}

// Close closes compressor
func (m *GlobalNotifierNats) Close() error {
	m.conn.Close()
	return nil
}
