package http

import (
	"github.com/copybird/copybird/core"
	"io"
	"net/http"
)

// Module Constants
const GROUP_NAME = "backup"
const TYPE_NAME = "output"
const MODULE_NAME = "http"

type BackupOutputHttp struct {
	core.Module
	reader io.Reader
	writer io.Writer
	config *Config
}

func (m *BackupOutputHttp) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *BackupOutputHttp) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *BackupOutputHttp) GetName() string {
	return MODULE_NAME
}

func (m *BackupOutputHttp) GetConfig() interface{} {
	return &Config{}
}

func (m *BackupOutputHttp) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *BackupOutputHttp) InitModule(_config interface{}) error {
	config := _config.(Config)
	m.config = &config

	return nil
}

func (m *BackupOutputHttp) Run() error {
	resp, err := http.Post(m.config.TargetUrl, "application/json", m.reader)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (m *BackupOutputHttp) Close() error {
	return nil
}
