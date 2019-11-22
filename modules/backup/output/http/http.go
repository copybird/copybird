package http

import (
	"context"
	"io"
	"net/http"

	"github.com/copybird/copybird/core"
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

func (m *BackupOutputHttp) InitModule(cfg interface{}) error {
	m.config = cfg.(*Config)

	return nil
}

func (m *BackupOutputHttp) Run(ctx context.Context) error {
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
