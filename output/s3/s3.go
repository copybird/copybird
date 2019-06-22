package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/copybird/copybird/output"
)

const MODULE_NAME = "s3"

type S3 struct {
	output.Output
	reader  io.Reader
	writer  io.Writer
	session *session.Session
	config  map[string]string
}

func (s *S3) GetName() string {
	return MODULE_NAME
}

func (s *S3) GetConfig() interface{} {
	return nil
}

func (s *S3) InitPipe(w io.Writer, r io.Reader) error {
	s.reader = r
	s.writer = w
	return nil
}

//InitOutput initializes S3 with session
func (s *S3) InitModule(_config interface{}) error {
	config := _config.(map[string]string)
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(config["AWS_REGION"]),
		Credentials: credentials.NewStaticCredentials(config["AWS_ACCESS_KEY_ID"], config["AWS_SECRET_ACCESS_KEY"], ""),
	})
	if err != nil {
		return err
	}

	s.session = session
	s.config = config
	return nil
}

func (s *S3) Run() error {

	svc := s3manager.NewUploader(s.session)

	input := &s3manager.UploadInput{
		Bucket: aws.String(s.config["AWS_BUCKET"]),
		Key:    aws.String(s.config["AWS_FILE_NAME"]),
		Body:   s.reader,
	}

	_, err := svc.Upload(input)
	if err != nil {
		return err
	}
	return nil
} 

func (s *S3) Close() error {
	return nil
}
