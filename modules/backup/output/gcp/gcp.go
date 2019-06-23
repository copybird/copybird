package gcp

import (
	"context"
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
	ctx        context.Context
	reader     io.Reader
	writer     io.Writer
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
	config     *Config
}

func (m *BackupOutputGcp) GetGroup() string {
	return GROUP_NAME
}

func (m *BackupOutputGcp) GetType() string {
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
	m.ctx = context.Background()

	switch {
	case m.config.CredentialsFilePath != "":
		client, err := storage.NewClient(m.ctx, option.WithCredentialsFile(m.config.CredentialsFilePath))
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

	m.bucket = m.client.Bucket(m.bucketName)
	// check if the bucket exists
	if _, err := m.bucket.Attrs(m.ctx); err != nil {
		return err
	}

	return nil
}

func (m *BackupOutputGcp) Run() error {

	obj := m.bucket.Object(m.config.AWSFileName)
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
