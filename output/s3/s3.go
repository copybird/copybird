package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/copybird/copybird/output"
)

type S3 struct {
	output.Output
	reader  io.Reader
	writer  io.Writer
	session *session.Session
	config  map[string]string
}

func (s *S3) Init(w io.Writer, r io.Reader) error {
	s.reader = r
	s.writer = w
	return nil
}

//InitOutput initializes S3 with session
func (s *S3) InitOutput(config map[string]string) error {

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

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return err
	}
	svc := s3manager.NewUploader(sess)
	input := &s3manager.UploadInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("whatever"),
		Body:   s.reader,
	}

	_, err = svc.Upload(input)
	if err != nil {
		return err
	}
	return nil
}

func (s *S3) Close() error {
	return nil
}
