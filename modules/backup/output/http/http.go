package http

import (
	"github.com/copybird/copybird/modules/backup/output"
	"io"
	"net/http"
)

const MODULE_NAME = "http"

type Http struct {
	output.Output
	reader io.Reader
	writer io.Writer
	config *Config
}

func (m *Http) GetName() string {
	return MODULE_NAME
}

func (m *Http) GetConfig() interface{} {
	return &Config{}
}

func (m *Http) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *Http) InitModule(_config interface{}) error {
	config := _config.(Config)
	m.config = &config

	return nil
}

func (m *Http) Run() error {
	resp, err := http.Post(m.config.TargetUrl, "application/json", m.reader)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (m *Http) Close() error {
	return nil
}
