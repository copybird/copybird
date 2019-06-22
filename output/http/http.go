package http

import (
	"io"
	"net/http"

	"github.com/copybird/copybird/output"
)

const MODULE_NAME = "http"

type Http struct {
	output.Output
	reader io.Reader
	writer io.Writer
	config *Config
}

func (h *Http) GetName() string {
	return MODULE_NAME
}

func (h *Http) GetConfig() interface{} {
	return Config{}
}

func (h *Http) InitPipe(w io.Writer, r io.Reader) error {
	h.reader = r
	h.writer = w
	return nil
}

func (h *Http) InitModule(_config interface{}) error {
	config := _config.(Config)
	h.config = &config

	return nil
}

func (h *Http) Run() error {
	resp, err := http.Post(h.config.targetUrl, "application/json", h.reader)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (h *Http) Close() error {
	return nil
}
	
