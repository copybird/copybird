package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/copybird/copybird/core"
)

// Module Constants
const GROUP_NAME = "backup"
const TYPE_NAME = "output"
const MODULE_NAME = "s3"

type BackupOutputS3 struct {
	core.Module
	reader  io.Reader
	writer  io.Writer
	session *session.Session
	config  *Config
}

func (m *BackupOutputS3) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *BackupOutputS3) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *BackupOutputS3) GetName() string {
	return MODULE_NAME
}

func (m *BackupOutputS3) GetConfig() interface{} {
	return &Config{}
}

func (m *BackupOutputS3) InitPipe(w io.Writer, r io.Reader) error {
	m.reader = r
	m.writer = w
	return nil
}

func (m *BackupOutputS3) InitModule(_config interface{}) error {
	m.config = _config.(*Config)
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(m.config.Region),
		Credentials: credentials.NewStaticCredentials(m.config.AccessKeyID, m.config.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	m.session = session
	return nil
}

func (m *BackupOutputS3) Run(ctx context.Context) error {

	svc := s3manager.NewUploader(m.session)

	input := &s3manager.UploadInput{
		Bucket: aws.String(m.config.Bucket),
		Key:    aws.String(m.config.FileName),
		Body:   m.reader,
	}

	_, err := svc.Upload(input)
	if err != nil {
		return err
	}
	return nil
}

func (m *BackupOutputS3) Close() error {
	return nil
}
