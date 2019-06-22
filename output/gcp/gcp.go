package gcp

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/copybird/copybird/output"
)

type GCP struct {
	output.Output
	ctx        context.Context
	reader     io.Reader
	writer     io.Writer
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
	config     map[string]string
}

func (gcp *GCP) Init(w io.Writer, r io.Reader) error {
	gcp.reader = r
	gcp.writer = w
	return nil
}

func (gcp *GCP) InitOutput(config map[string]string) error {

	gcp.ctx = context.Background()

	client, err := storage.NewClient(gcp.ctx)
	if err != nil {
		return err
	}

	gcp.bucket = client.Bucket(gcp.bucketName)
	// check if the bucket exists
	if _, err := gcp.bucket.Attrs(gcp.ctx); err != nil {
		return err
	}

	return nil
}

func (gcp *GCP) Run() error {

	obj := gcp.bucket.Object(gcp.config["AWS_FILE_NAME"])
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
