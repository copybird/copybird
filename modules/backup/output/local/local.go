package local

import (
	"github.com/copybird/copybird/core"
	"io"
	"os"
)

// Module Constants
const GROUP_NAME = "backup"
const TYPE_NAME = "output"
const MODULE_NAME = "local"

type BackupOutputLocal struct {
	core.Module
	reader io.Reader
	writer io.Writer
	config *Config
}

func (m *BackupOutputLocal) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *BackupOutputLocal) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *BackupOutputLocal) GetName() string {
	return MODULE_NAME
}

func (m *BackupOutputLocal) GetConfig() interface{} {
	return &Config{
		File:        "output",
		DefaultMask: os.O_APPEND | os.O_CREATE | os.O_WRONLY,
	}
}

func (m *BackupOutputLocal) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *BackupOutputLocal) InitModule(_config interface{}) error {
	m.config = _config.(*Config)
	return nil
}

func (m *BackupOutputLocal) Run() error {

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(m.config.File, m.config.DefaultMask, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, m.reader)
	return err
}

func (m *BackupOutputLocal) Close() error {
	return nil
}
