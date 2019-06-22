package local

import (
	"io"
	"os"
	"github.com/copybird/copybird/output"
)

const MODULE_NAME = "local"

type Local struct {
	output.Output
	reader  io.Reader
	writer  io.Writer
	config  *Config
}

func (loc *Local) GetName() string {
	return MODULE_NAME
}

func (loc *Local) GetConfig() interface{} {
	return Config{
		DefaultMask: os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		FileName: "test.txt",
	}
}

func (loc *Local) InitPipe(w io.Writer, r io.Reader) error {
	loc.reader = r
	loc.writer = w
	return nil
}

func (loc *Local) InitModule(_config interface{}) error {
	config := _config.(Config)
	loc.config = &config

	return nil
}

func (loc *Local) Run() error {

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(loc.config.FileName, loc.config.DefaultMask, 0644)
	if err != nil {
		return err
	}
	
    defer f.Close()

	_, err = io.Copy(f, loc.reader)
	return err
} 

func (loc *Local) Close() error {
	return nil
}
