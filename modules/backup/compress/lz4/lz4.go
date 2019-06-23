package lz4

import (
	"errors"
	"fmt"
	"github.com/copybird/copybird/core"
	"io"

	"github.com/pierrec/lz4"
)

var (
	errCompLevel       = errors.New("compression level must be between -1 and 9")
	errNotCompressible = errors.New("is not compressible")
)

const GROUP_NAME = "backup"
const TYPE_NAME = "compress"
const MODULE_NAME = "lz4"

// BackupCompressLz4 represents ...
type BackupCompressLz4 struct {
	core.Module
	reader io.Reader
	writer io.Writer
	level  int
}

func (m *BackupCompressLz4) GetGroup() string {
	return GROUP_NAME
}

func (m *BackupCompressLz4) GetType() string {
	return TYPE_NAME
}

func (m *BackupCompressLz4) GetName() string {
	return MODULE_NAME
}

func (m *BackupCompressLz4) GetConfig() interface{} {
	return &Config{Level: 2}
}

func (m *BackupCompressLz4) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *BackupCompressLz4) InitModule(_cfg interface{}) error {
	cfg := _cfg.(*Config)
	if cfg.Level < -1 || cfg.Level > 9 {
		return errCompLevel
	}
	m.level = cfg.Level
	return nil
}

func (m *BackupCompressLz4) Run() error {
	lw := lz4.NewWriter(m.writer)
	lw.Header = lz4.Header{CompressionLevel: m.level}
	defer lw.Close()

	_, err := io.Copy(lw, m.reader)
	if err != nil {
		return fmt.Errorf("copy error: %s", err)
	}
	return nil
}

// Close closes compressor
func (m *BackupCompressLz4) Close() error {
	return nil
}
