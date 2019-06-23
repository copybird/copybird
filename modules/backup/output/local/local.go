package local

import (
	"github.com/copybird/copybird/modules/backup/output"
	"io"
	"os"
)

const MODULE_NAME = "local"

type Local struct {
	output.Output
	reader io.Reader
	writer io.Writer
	config *Config
}

func (m *Local) GetName() string {
	return MODULE_NAME
}

func (m *Local) GetConfig() interface{} {
	return &Config{
		FileName:    "output",
		DefaultMask: os.O_APPEND | os.O_CREATE | os.O_WRONLY,
	}
}

func (m *Local) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *Local) InitModule(_config interface{}) error {
	m.config = _config.(*Config)
	return nil
}

func (m *Local) Run() error {

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(m.config.FileName, m.config.DefaultMask, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, m.reader)
	return err
}

func (m *Local) Close() error {
	return nil
}
