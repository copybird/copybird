package gcp

import (
	"context"
	"errors"
	"github.com/copybird/copybird/core"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Module Constants
const GROUP_NAME = "backup"
const TYPE_NAME = "output"
const MODULE_NAME = "gcp"

type BackupOutputGcp struct {
	core.Module
	ctx    context.Context
	reader io.Reader
	writer io.Writer
	client *storage.Client
	bucket *storage.BucketHandle
	config *Config
}

func (m *BackupOutputGcp) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *BackupOutputGcp) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *BackupOutputGcp) GetName() string {
	return MODULE_NAME
}

func (m *BackupOutputGcp) GetConfig() interface{} {
	return &Config{}
}

func (m *BackupOutputGcp) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *BackupOutputGcp) InitModule(_config interface{}) error {
	m.config = _config.(*Config)

	if m.config.AuthFile == "" {
		return errors.New("need auth_file")
	}
	if m.config.Bucket == "" {
		return errors.New("need bucker")
	}
	if m.config.File == "" {
		return errors.New("need file")
	}

	m.ctx = context.Background()

	switch {
	case m.config.AuthFile != "":
		client, err := storage.NewClient(m.ctx, option.WithCredentialsFile(m.config.AuthFile))
		if err != nil {
			return err
		}
		m.client = client
	default:
		client, err := storage.NewClient(m.ctx)
		if err != nil {
			return err
		}
		m.client = client
	}

	m.bucket = m.client.Bucket(m.config.Bucket)
	// check if the bucket exists
	if _, err := m.bucket.Attrs(m.ctx); err != nil {
		return err
	}

	return nil
}

func (m *BackupOutputGcp) Run(ctx context.Context) error {

	obj := m.bucket.Object(m.config.File)
	w := obj.NewWriter(m.ctx)
	if _, err := io.Copy(w, m.reader); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	_, err := obj.Attrs(m.ctx)
	return err
}

func (m *BackupOutputGcp) Close() error {
	m.client.Close()
	return nil
}
