package gcp

import (
	"context"
	output2 "github.com/copybird/copybird/backup/output"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const MODULE_NAME = "gcp"

type GCP struct {
	output2.Output
	ctx        context.Context
	reader     io.Reader
	writer     io.Writer
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
	config     *Config
}

func (gcp *GCP) GetName() string {
	return MODULE_NAME
}

func (gcp *GCP) GetConfig() interface{} {
	return &Config{}
}

func (gcp *GCP) InitPipe(w io.Writer, r io.Reader) error {
	gcp.reader = r
	gcp.writer = w
	return nil
}

func (gcp *GCP) InitModule(_config interface{}) error {
	config := _config.(Config)
	gcp.config = &config
	gcp.ctx = context.Background()

	switch {
	case config.CredentialsFilePath != "":
		client, err := storage.NewClient(gcp.ctx, option.WithCredentialsFile(config.CredentialsFilePath))
		if err != nil {
			return err
		}
		gcp.client = client
	default:
		client, err := storage.NewClient(gcp.ctx)
		if err != nil {
			return err
		}
		gcp.client = client
	}

	gcp.bucket = gcp.client.Bucket(gcp.bucketName)
	// check if the bucket exists
	if _, err := gcp.bucket.Attrs(gcp.ctx); err != nil {
		return err
	}

	return nil
}

func (gcp *GCP) Run() error {

	obj := gcp.bucket.Object(gcp.config.AWSFileName)
	w := obj.NewWriter(gcp.ctx)
	if _, err := io.Copy(w, gcp.reader); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	_, err := obj.Attrs(gcp.ctx)
	return err
}

func (gcp *GCP) Close() error {
	gcp.client.Close()
	return nil
}
