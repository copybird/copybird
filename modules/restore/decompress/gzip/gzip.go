package gzip

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"

	"github.com/copybird/copybird/core"
)

const GROUP_NAME = "restore"
const TYPE_NAME = "decompress"
const MODULE_NAME = "gzip"

type RestoreDecompressGzip struct {
	core.Module
	reader io.Reader
	writer io.Writer
}

func (m *RestoreDecompressGzip) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *RestoreDecompressGzip) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *RestoreDecompressGzip) GetName() string {
	return MODULE_NAME
}

func (m *RestoreDecompressGzip) GetConfig() interface{} {
	return nil
}

func (m *RestoreDecompressGzip) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *RestoreDecompressGzip) InitModule(_cfg interface{}) error {
	return nil
}

func (m *RestoreDecompressGzip) Run(ctx context.Context) error {
	gr, err := gzip.NewReader(m.reader)
	if err != nil {
		return fmt.Errorf("cant start gzip reader with error: %s", err)
	}
	defer gr.Close()
	_, err = io.Copy(m.writer, gr)
	if err != nil {
		return fmt.Errorf("copy error: %s", err)
	}
	return nil
}

func (m *RestoreDecompressGzip) Close() error {
	return nil
}
