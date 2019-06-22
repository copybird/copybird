package local

import (
	"io"
	"os"
	"fmt"
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
	return Config{}
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

	f, err := os.Create(loc.config.FileName)
    if err != nil {
        return err
	}
	
    defer f.Close()

	buf := make([]byte, 12)
	for {
		_, err = loc.reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read err: %s", err)
		}
		_, err = loc.writer.Write(buf)
		if err != nil {
			return fmt.Errorf("write err: %s", err)
		}
	}
	
	return nil
} 

func (loc *Local) Close() error {
	return nil
}
