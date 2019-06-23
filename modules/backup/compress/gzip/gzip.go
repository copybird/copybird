package gzip

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"

	"github.com/copybird/copybird/core"
)

const MODULE_NAME = "gzip"
const MODULE_GROUP = "backup"
const MODULE_TYPE = "compress"

type BackupCompressGzip struct {
	core.Module
	reader io.Reader
	writer io.Writer
	level  int
}

func (m *BackupCompressGzip) GetName() string {
	return MODULE_NAME
}

func (m *BackupCompressGzip) GetGroup() core.ModuleGroup {
	return MODULE_GROUP
}

func (m *BackupCompressGzip) GetType() core.ModuleType {
	return MODULE_TYPE
}

func (m *BackupCompressGzip) GetConfig() interface{} {
	return &Config{
		Level: 3,
	}
}

func (m *BackupCompressGzip) InitModule(_cfg interface{}) error {
	cfg := _cfg.(Config)
	level := cfg.Level
	if level < -1 || level > 9 {
		return errors.New("compression level must be between -1 and 9")
	}
	m.level = level
	return nil
}

func (m *BackupCompressGzip) Run() error {
	gw, err := gzip.NewWriterLevel(m.writer, m.level)
	if err != nil {
		return fmt.Errorf("cant start gzip writer with error: %s", err)
	}
	defer gw.Close()

	_, err = io.Copy(gw, m.reader)
	if err != nil {
		return fmt.Errorf("copy error: %s", err)
	}
	return nil
}

func (m *BackupCompressGzip) Close() error {
	return nil
}

func (m *BackupCompressGzip) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}
